import { useState, useEffect } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { useRefinement } from "@/hooks/use-refinement";
import { RefinementProgress } from "./RefinementProgress";
import { Badge } from "@/components/ui/badge";
import { Sparkles, Copy, Check, Loader2, Wand2, FileText, Variable, ChevronDown, ChevronUp, AlertTriangle, ShieldCheck } from "lucide-react";
import Editor from "@monaco-editor/react";
import { refinementApi } from "@/api/refinement";
import { getAccessToken, API_BASE_URL } from "@/api/client";
import type { ContractDefinition } from "@/api/contracts";

interface RecommendationModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    description: string;
    targetType: string; // 'context', 'contract', 'variable'
    contextData: any;
    contracts?: ContractDefinition[];
    refineContent?: any;
    onApply?: (result: any) => void;
}

const TARGET_TYPE_LABELS: Record<string, string> = {
    context: "business and technical context",
    contract: "API contracts and variables",
    variable: "environment variables and secrets",
    requirement: "technical requirements and variables",
    validation_rule: "validation rules and constraints",
};

export function RecommendationModal({
    isOpen,
    onClose,
    title,
    description,
    targetType,
    contextData,
    contracts = [],
    refineContent,
    onApply
}: RecommendationModalProps) {
    const { session, events, startSession, reset } = useRefinement();
    const [prompt, setPrompt] = useState("");
    const [isCopied, setIsCopied] = useState(false);
    const [isRefining, setIsRefining] = useState(false);
    const [selectedContractId, setSelectedContractId] = useState<string>("");
    const [evaluationOpen, setEvaluationOpen] = useState(false);

    const getScoreColor = (score: number) => {
        if (score >= 8) return "bg-green-100 text-green-800 border-green-200";
        if (score >= 6) return "bg-yellow-100 text-yellow-800 border-yellow-200";
        return "bg-red-100 text-red-800 border-red-200";
    };

    // Reset state when modal opens
    useEffect(() => {
        if (isOpen) {
            reset();
            // If refining content, default prompt to explain the goal
            if (refineContent) {
                setPrompt("Review and improve the existing content for completeness, security, and best practices.");
            } else {
                setPrompt("");
            }
            // Default to first contract if available and not set
            if (targetType === "variable" && contracts.length > 0) {
                setSelectedContractId(contracts[0].id);
            } else {
                setSelectedContractId("");
            }
        }
    }, [isOpen, targetType, contracts, refineContent]);

    const handleRefineInstructions = async () => {
        setIsRefining(true);
        try {
            const targetLabel = TARGET_TYPE_LABELS[targetType] || targetType;
            const itemTitle = contextData?.title || "this roadmap item";

            // If a contract is selected, include its info in the refinement context
            let contractContext = "";
            if (targetType === "variable" && selectedContractId) {
                const contract = contracts.find(c => c.id === selectedContractId);
                if (contract) {
                    // Try to use parsed schema if available, or fall back to description/type
                    // We assume input_schema/output_schema might be populated from the list
                    const schemaSummary = JSON.stringify({
                        type: contract.contract_type,
                        version: contract.version,
                        input: contract.input_schema,
                        output: contract.output_schema
                    }, null, 2);
                    contractContext = `\nThe user has selected the existing "${contract.contract_type}" contract (v${contract.version}). Please ensure variables match this contract's requirements:\n${schemaSummary.substring(0, 1000)}...`; // Truncate to avoid context limit
                }
            }

            // If refining exist content
            let refineContext = "";
            if (refineContent) {
                refineContext = `\nExisting Content to Refine:\n${JSON.stringify(refineContent, null, 2).substring(0, 2000)}...`;
            }

            const refinePrompt = [
                `You are an expert at writing precise instructions for AI-powered code generation.`,
                `The user wants to generate ${targetLabel} for the roadmap item: "${itemTitle}".`,
                contextData?.description ? `Item description: ${contextData.description}` : "",
                contextData?.business_context ? `Business context: ${contextData.business_context.substring(0, 500)}` : "",
                contextData?.technical_context ? `Technical context: ${contextData.technical_context.substring(0, 500)}` : "",
                contractContext,
                refineContext,
                ``,
                prompt.trim()
                    ? `The user has written these draft instructions:\n"${prompt}"\n\nPlease refine, expand, and improve these instructions to be more specific, actionable, and complete.`
                    : `The user has not provided any instructions yet. Please generate clear, specific, and actionable instructions that would produce excellent ${targetLabel} for this item.`,
                ``,
                `You MUST return a JSON object with a single field "instructions" containing the improved instruction text as a string.`,
                `Example: {"instructions": "Your improved instructions here..."}`,
                `Return ONLY the JSON object, no markdown formatting.`,
            ].filter(Boolean).join("\n");

            // Use refinement API for instruction refinement
            const newSession = await refinementApi.startSession(
                "instruction",
                targetType,
                refinePrompt,
                contextData,
                3
            );

            if (!newSession.id) throw new Error("Session ID missing");

            // Stream the SSE response and extract the refined text
            const url = `${API_BASE_URL}${refinementApi.getEventsEndpoint(newSession.id)}`;
            const token = getAccessToken();
            const headers: HeadersInit = {};
            if (token) headers["Authorization"] = `Bearer ${token}`;

            const response = await fetch(url, { headers });
            if (!response.body) throw new Error("No response body");

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let foundResult = false;

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split("\n\n");

                for (const line of lines) {
                    if (line.startsWith("data: ")) {
                        const dataStr = line.replace("data: ", "");
                        if (dataStr === "{}") continue;

                        try {
                            const event = JSON.parse(dataStr);
                            if (event.type === "SUCCESS" && event.payload?.artifact) {
                                const result = event.payload.artifact;
                                const refinedText = typeof result === "string"
                                    ? result
                                    : result.instructions || result.prompt || result.text || JSON.stringify(result, null, 2);
                                setPrompt(refinedText);
                                foundResult = true;
                            }
                            if (event.type === "ERROR" && event.message) {
                                console.warn("Refine error event:", event.message);
                            }
                        } catch {
                            // Skip parse errors
                        }
                    }
                }
            }

            if (!foundResult) {
                console.warn("Refine completed without producing a result");
            }
        } catch (err: any) {
            console.error("Failed to refine instructions:", err);
        } finally {
            setIsRefining(false);
        }
    };

    const handleStart = () => {
        // Prepare context data - include selected contract if applicable
        const enhancedContext = { ...contextData };
        if (targetType === "variable" && selectedContractId) {
            const contract = contracts.find(c => c.id === selectedContractId);
            if (contract) {
                enhancedContext.selected_contract = {
                    type: contract.contract_type,
                    version: contract.version,
                    input: contract.input_schema,
                    output: contract.output_schema
                };
            }
        }

        // Include existing content if refining
        if (refineContent) {
            enhancedContext.existing_artifact = refineContent;
        }

        startSession(
            "recommendation", // ArtifactType generic
            targetType,
            prompt || "Generate the requested artifact based on the context.",
            enhancedContext,
            3 // Default max iterations for recommendations
        );
    };

    const handleCopy = () => {
        if (session?.result) {
            navigator.clipboard.writeText(JSON.stringify(session.result, null, 2));
            setIsCopied(true);
            setTimeout(() => setIsCopied(false), 2000);
        }
    };

    const handleApply = () => {
        if (session?.result && onApply) {
            // Pass the selected contract ID along with the result if needed
            const resultToApply = targetType === "variable"
                ? { ...session.result, selectedContractId }
                : session.result;

            onApply(resultToApply);
            onClose();
        }
    };

    const renderPreview = (type: string, data: any) => {
        if (!data) return null;

        switch (type) {
            case "context":
                const renderContextValue = (value: any) => {
                    if (typeof value === 'string') return value;
                    if (typeof value === 'object' && value !== null) {
                        return JSON.stringify(value, null, 2);
                    }
                    return String(value || "Not generated.");
                };

                return (
                    <div className="space-y-6 text-sm">
                        <div className="space-y-2">
                            <h4 className="font-semibold text-indigo-600 flex items-center gap-2">
                                <Sparkles className="h-4 w-4" /> Business Context
                            </h4>
                            <p className="text-muted-foreground whitespace-pre-wrap leading-relaxed font-mono text-xs">
                                {renderContextValue(data.business_context)}
                            </p>
                        </div>
                        <div className="space-y-2">
                            <h4 className="font-semibold text-emerald-600 flex items-center gap-2">
                                <Sparkles className="h-4 w-4" /> Technical Context
                            </h4>
                            <p className="text-muted-foreground whitespace-pre-wrap leading-relaxed font-mono text-xs">
                                {renderContextValue(data.technical_context)}
                            </p>
                        </div>
                    </div>
                );

            case "variable":
                // Expecting array of { name, description, required }
                const vars = Array.isArray(data) ? data : (data.variables || []);
                if (vars.length === 0) return <div className="text-muted-foreground italic">No variables identified.</div>;

                return (
                    <div className="border rounded-md overflow-hidden">
                        <table className="w-full text-sm">
                            <thead className="bg-muted/50">
                                <tr>
                                    <th className="px-3 py-2 text-left font-medium">Name</th>
                                    <th className="px-3 py-2 text-left font-medium">Description</th>
                                    <th className="px-3 py-2 text-center font-medium">Required</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y">
                                {vars.map((v: any, i: number) => (
                                    <tr key={i} className="hover:bg-muted/20">
                                        <td className="px-3 py-2 font-mono text-indigo-600">{v.name}</td>
                                        <td className="px-3 py-2 text-muted-foreground">{v.description}</td>
                                        <td className="px-3 py-2 text-center">
                                            {v.required ? <Check className="h-4 w-4 mx-auto text-green-500" /> : <span className="text-muted-foreground">-</span>}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                );

            case "validation_rule":
                return (
                    <div className="space-y-4">
                        {data.rules?.map((rule: any, i: number) => (
                            <div key={i} className="p-3 border rounded-md bg-white dark:bg-slate-800 space-y-2">
                                <div className="flex items-center justify-between">
                                    <div className="font-semibold flex items-center gap-2">
                                        <ShieldCheck className="h-4 w-4 text-emerald-500" />
                                        {rule.name}
                                    </div>
                                    <Badge variant="outline">{rule.rule_type}</Badge>
                                </div>
                                <p className="text-xs text-muted-foreground">{rule.description}</p>
                                {rule.rule_config && (
                                    <div className="text-[10px] font-mono bg-slate-50 dark:bg-slate-900 p-2 rounded border">
                                        {JSON.stringify(rule.rule_config, null, 2)}
                                    </div>
                                )}
                            </div>
                        ))}
                        {!data.rules?.length && <div className="text-center text-muted-foreground py-4">No rules generated.</div>}
                    </div>
                );
            case "requirement":
                const reqs = data.requirements || [];
                const reqVars = data.variables || [];
                return (
                    <div className="space-y-6">
                        <div className="space-y-2">
                            <h4 className="font-semibold text-indigo-600 flex items-center gap-2">
                                <FileText className="h-4 w-4" /> Recommended Requirements
                            </h4>
                            {reqs.length === 0 ? (
                                <p className="text-muted-foreground italic">No requirements generated.</p>
                            ) : (
                                <div className="border rounded-md overflow-hidden">
                                    <table className="w-full text-sm">
                                        <thead className="bg-muted/50">
                                            <tr>
                                                <th className="px-3 py-2 text-left font-medium">Title</th>
                                                <th className="px-3 py-2 text-left font-medium">Priority</th>
                                                <th className="px-3 py-2 text-left font-medium">Testable</th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y">
                                            {reqs.map((r: any, i: number) => (
                                                <tr key={i} className="hover:bg-muted/20">
                                                    <td className="px-3 py-2 font-medium">{r.title}</td>
                                                    <td className="px-3 py-2">
                                                        <span className={`text-xs px-2 py-0.5 rounded-full border ${r.priority === 'HIGH' ? 'bg-red-100 text-red-800 border-red-200' :
                                                            r.priority === 'MEDIUM' ? 'bg-yellow-100 text-yellow-800 border-yellow-200' :
                                                                'bg-slate-100 text-slate-800 border-slate-200'
                                                            }`}>
                                                            {r.priority}
                                                        </span>
                                                    </td>
                                                    <td className="px-3 py-2 text-center">
                                                        {r.testable ? <Check className="h-4 w-4 mx-auto text-green-500" /> : <span className="text-muted-foreground">-</span>}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </div>

                        {reqVars.length > 0 && (
                            <div className="space-y-2">
                                <h4 className="font-semibold text-emerald-600 flex items-center gap-2">
                                    <Variable className="h-4 w-4" /> Suggested Variables
                                </h4>
                                <div className="border rounded-md overflow-hidden">
                                    <table className="w-full text-sm">
                                        <thead className="bg-muted/50">
                                            <tr>
                                                <th className="px-3 py-2 text-left font-medium">Name</th>
                                                <th className="px-3 py-2 text-left font-medium">Default</th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y">
                                            {reqVars.map((v: any, i: number) => (
                                                <tr key={i} className="hover:bg-muted/20">
                                                    <td className="px-3 py-2 font-mono text-indigo-600">{v.name}</td>
                                                    <td className="px-3 py-2 text-muted-foreground font-mono text-xs">{v.default_value || "-"}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        )}
                    </div>
                );

            case "contract":
            default:
                return (
                    <div className="h-96">
                        <Editor
                            height="100%"
                            defaultLanguage="json"
                            value={JSON.stringify(data, null, 2)}
                            options={{ readOnly: true, minimap: { enabled: false }, scrollBeyondLastLine: false }}
                        />
                    </div>
                );
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
            <DialogContent className="sm:max-w-[800px] h-[80vh] flex flex-col">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        <Sparkles className="h-5 w-5 text-indigo-500" />
                        {title}
                    </DialogTitle>
                    <DialogDescription>
                        {description}
                    </DialogDescription>
                </DialogHeader>

                <div className="flex-1 overflow-y-auto py-4 space-y-4">
                    {!session ? (
                        <div className="space-y-4">
                            {/* Contract Selection for Variable Generation */}
                            {targetType === "variable" && contracts.length > 0 && (
                                <div className="space-y-2">
                                    <Label>Target Contract</Label>
                                    <Select
                                        value={selectedContractId}
                                        onValueChange={setSelectedContractId}
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select a contract to analyze..." />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {contracts.map(c => (
                                                <SelectItem key={c.id} value={c.id}>
                                                    {c.contract_type} (v{c.version})
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    <p className="text-xs text-muted-foreground">
                                        Select the contract this variable will belong to. The AI will analyze it to suggest relevant variables.
                                    </p>
                                </div>
                            )}

                            <div className="space-y-2">
                                <div className="flex items-center justify-between">
                                    <Label>Instructions (Optional)</Label>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-7 gap-1.5 text-xs text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50"
                                        onClick={handleRefineInstructions}
                                        disabled={isRefining}
                                    >
                                        {isRefining ? (
                                            <>
                                                <Loader2 className="h-3 w-3 animate-spin" />
                                                Refining...
                                            </>
                                        ) : (
                                            <>
                                                <Wand2 className="h-3 w-3" />
                                                Refine with AI
                                            </>
                                        )}
                                    </Button>
                                </div>
                                <Textarea
                                    placeholder="e.g., 'Focus on security constraints' or 'Include these specific env vars...'"
                                    value={prompt}
                                    onChange={(e) => setPrompt(e.target.value)}
                                    className="resize-y"
                                    rows={4}
                                />
                                <p className="text-xs text-muted-foreground">
                                    The AI will analyze the current Roadmap Item details to generate the recommendation.
                                    Use <strong>Refine with AI</strong> to improve your instructions before generating.
                                </p>
                            </div>
                        </div>
                    ) : (
                        <div className="flex-1 flex flex-col space-y-4 h-full">
                            <RefinementProgress events={events} status={session.status} />

                            {session.status === 'VALIDATED' && (
                                <div className="flex-1 flex flex-col space-y-4 overflow-hidden">
                                    {/* @ts-ignore */}
                                    {session.evaluation && (
                                        <div className="border rounded-lg overflow-hidden border-indigo-100 bg-indigo-50/30">
                                            <button
                                                onClick={() => setEvaluationOpen(!evaluationOpen)}
                                                className="w-full flex items-center justify-between px-4 py-3 hover:bg-indigo-50/50 transition-colors"
                                            >
                                                <div className="flex items-center gap-3">
                                                    <ShieldCheck className="h-5 w-5 text-indigo-600" />
                                                    <span className="font-semibold text-indigo-900">AI Quality Score</span>
                                                    <Badge className={getScoreColor((session as any).evaluation.score)}>
                                                        {(session as any).evaluation.score}/10
                                                    </Badge>
                                                </div>
                                                {evaluationOpen ? <ChevronUp className="h-4 w-4 text-indigo-400" /> : <ChevronDown className="h-4 w-4 text-indigo-400" />}
                                            </button>

                                            {evaluationOpen && (
                                                <div className="px-4 pb-4 pt-1 space-y-4 text-sm border-t border-indigo-100">
                                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                                        {(session as any).evaluation.security_concerns?.length > 0 && (
                                                            <div className="space-y-1">
                                                                <h5 className="font-semibold text-red-700 flex items-center gap-1.5">
                                                                    <AlertTriangle className="h-3.5 w-3.5" /> Security
                                                                </h5>
                                                                <ul className="list-disc list-inside text-slate-600 pl-1">
                                                                    {(session as any).evaluation.security_concerns.map((item: string, i: number) => (
                                                                        <li key={i}>{item}</li>
                                                                    ))}
                                                                </ul>
                                                            </div>
                                                        )}
                                                        {(session as any).evaluation.improvement_suggestions?.length > 0 && (
                                                            <div className="space-y-1">
                                                                <h5 className="font-semibold text-indigo-700">Suggestions</h5>
                                                                <ul className="list-disc list-inside text-slate-600 pl-1">
                                                                    {(session as any).evaluation.improvement_suggestions.map((item: string, i: number) => (
                                                                        <li key={i}>{item}</li>
                                                                    ))}
                                                                </ul>
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    )}

                                    <div className="flex-1 flex flex-col space-y-2 overflow-hidden">
                                        <Label>Result Preview</Label>
                                        <div className="flex-1 border rounded-md overflow-y-auto bg-slate-50 dark:bg-slate-900 p-4">
                                            {renderPreview(targetType, session.result)}
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                <DialogFooter>
                    <Button variant="outline" onClick={onClose}>Cancel</Button>

                    {!session ? (
                        <Button onClick={handleStart} className="gap-2" disabled={targetType === "variable" && contracts.length > 0 && !selectedContractId}>
                            <Sparkles className="h-4 w-4" />
                            Generate Recommendation
                        </Button>
                    ) : (
                        <>
                            {session.status === 'VALIDATED' && (
                                <>
                                    <Button variant="secondary" onClick={handleCopy} disabled={!session.result}>
                                        {isCopied ? <Check className="h-4 w-4 mr-2" /> : <Copy className="h-4 w-4 mr-2" />}
                                        {isCopied ? "Copied" : "Copy JSON"}
                                    </Button>
                                    {onApply && (
                                        <Button onClick={handleApply}>
                                            Apply to Item
                                        </Button>
                                    )}
                                </>
                            )}
                            {(session.status === 'FAILED' || session.status === 'VALIDATED') && (
                                <Button variant="ghost" onClick={reset}>Start Over</Button>
                            )}
                        </>
                    )}
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
