import { useParams } from "react-router-dom";
import { useRoadmapItem, useUpdateRoadmapItem } from "@/hooks/use-roadmap-item";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { AlertTriangle, Shield, Sparkles, History, Info, Variable, Wand2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { RecommendationModal } from "./components/RecommendationModal";
import { IntelligencePanel } from "../intelligence/components/IntelligencePanel";
import { useState, useEffect } from 'react';
import Editor from '@monaco-editor/react';
import { contractsApi } from '@/api/contracts';
import { variablesApi } from '@/api/variables';
import { intelligenceApi } from '@/api/intelligence';
import { BuildArtifactPanel } from "./components/BuildArtifactPanel";
import type { ContractDefinition } from '@/api/contracts';
import ReactMarkdown from 'react-markdown';

export default function RoadmapItemPage() {
    const { roadmapItemId } = useParams<{ roadmapItemId: string }>();
    const { data: item, isLoading } = useRoadmapItem(roadmapItemId);
    const updateItem = useUpdateRoadmapItem();
    const queryClient = useQueryClient();

    const { data: contracts = [] } = useQuery({
        queryKey: ["contracts", roadmapItemId],
        queryFn: () => contractsApi.listContracts(roadmapItemId!),
        enabled: !!roadmapItemId
    });

    const { data: allVariables = [] } = useQuery({
        queryKey: ["variables", item?.project_id],
        queryFn: () => {
            if (!item?.project_id) return Promise.resolve([]);
            return variablesApi.listVariables(item.project_id);
        },
        enabled: !!item?.project_id
    });

    // Filter variables relevant to this roadmap item (via contracts)
    const variables = allVariables.filter(v =>
        contracts.some(c => c.id === v.contract_id)
    );

    const { data: driftHistory = [] } = useQuery({
        queryKey: ["drift-history"],
        queryFn: intelligenceApi.getDriftHistory
    });

    const relevantReports = driftHistory.filter(log =>
        contracts.some(c => c.id === log.entity_id) && log.new_data?.drift_report
    );

    const [selectedContract, setSelectedContract] = useState<ContractDefinition | null>(null);

    // Set default selected contract
    useEffect(() => {
        if (contracts.length > 0 && !selectedContract) {
            setSelectedContract(contracts[0]);
        }
    }, [contracts, selectedContract]);

    const [recModal, setRecModal] = useState<{
        isOpen: boolean;
        title: string;
        description: string;
        targetType: string;
        refineContent?: any;
    }>({
        isOpen: false,
        title: "",
        description: "",
        targetType: "context"
    });

    if (isLoading) {
        return <div className="p-8 text-center italic text-muted-foreground">Loading item details...</div>;
    }

    if (!item) {
        return <div className="p-8 text-center text-red-600">Roadmap item not found.</div>;
    }

    return (
        <div className="flex h-[calc(100vh-4rem)] overflow-hidden">
            {/* Main Content Area (Center Panel) */}
            <div className="flex-1 overflow-y-auto p-6 space-y-6">
                <div className="flex items-center justify-between">
                    <div className="space-y-1">
                        <div className="flex items-center gap-2">
                            <h2 className="text-3xl font-bold tracking-tight">{item.title}</h2>
                            <Badge variant={item.risk_level === "HIGH" ? "destructive" : "secondary"}>
                                {item.risk_level} RISK
                            </Badge>
                        </div>
                        <div className="text-muted-foreground text-sm [&_p]:mb-2 [&_strong]:text-foreground [&_ul]:list-disc [&_ul]:pl-4 [&_ol]:list-decimal [&_ol]:pl-4 [&_li]:mb-1 [&_code]:bg-muted [&_code]:px-1 [&_code]:py-0.5 [&_code]:rounded [&_code]:text-xs [&_code]:font-mono [&_h1]:text-lg [&_h1]:font-semibold [&_h1]:mb-2 [&_h1]:text-foreground [&_h2]:text-base [&_h2]:font-semibold [&_h2]:mb-1.5 [&_h2]:text-foreground [&_h3]:text-sm [&_h3]:font-medium [&_h3]:mb-1 [&_h3]:text-foreground">
                            <ReactMarkdown>{item.description}</ReactMarkdown>
                        </div>
                    </div>
                    <Badge variant="outline" className="px-3 py-1 font-mono">
                        {item.status}
                    </Badge>
                </div>

                <Tabs defaultValue="overview" className="w-full">
                    <TabsList className="grid w-full grid-cols-5 lg:w-[750px]">
                        <TabsTrigger value="overview" className="flex gap-2">
                            <Info className="h-4 w-4" /> Overview
                        </TabsTrigger>
                        <TabsTrigger value="contracts" className="flex gap-2">
                            <Shield className="h-4 w-4" /> Contracts
                        </TabsTrigger>
                        <TabsTrigger value="ai-proposals" className="flex gap-2">
                            <Sparkles className="h-4 w-4" /> AI Proposals
                        </TabsTrigger>
                        <TabsTrigger value="variables" className="flex gap-2">
                            <Variable className="h-4 w-4" /> Variables
                        </TabsTrigger>
                        <TabsTrigger value="reports" className="flex gap-2">
                            <History className="h-4 w-4" /> Reports
                        </TabsTrigger>
                    </TabsList>

                    <TabsContent value="overview" className="space-y-4 pt-4">
                        <div className="grid gap-4 md:grid-cols-2">
                            <Card>
                                <CardHeader className="flex flex-row items-center justify-between">
                                    <CardTitle className="text-sm font-medium">Business Context</CardTitle>
                                    <Button variant="ghost" size="icon" onClick={() => setRecModal({
                                        isOpen: true,
                                        title: "Refine Context",
                                        description: "Analyze and improve business and technical context.",
                                        targetType: "context"
                                    })}>
                                        <Sparkles className="h-4 w-4 text-muted-foreground hover:text-primary" />
                                    </Button>
                                </CardHeader>
                                <CardContent className="max-h-[400px] overflow-y-auto">
                                    {item.business_context ? (
                                        <div className="text-sm text-slate-700 [&_p]:mb-2 [&_p]:leading-relaxed [&_strong]:text-slate-900 [&_ul]:list-disc [&_ul]:pl-4 [&_ul]:space-y-1 [&_ol]:list-decimal [&_ol]:pl-4 [&_ol]:space-y-1 [&_li]:leading-relaxed [&_code]:bg-slate-100 [&_code]:px-1 [&_code]:py-0.5 [&_code]:rounded [&_code]:text-xs [&_code]:font-mono [&_h1]:text-base [&_h1]:font-semibold [&_h1]:mb-2 [&_h1]:text-slate-900 [&_h2]:text-sm [&_h2]:font-semibold [&_h2]:mb-1.5 [&_h2]:text-slate-900 [&_h3]:text-sm [&_h3]:font-medium [&_h3]:mb-1 [&_h3]:text-slate-800">
                                            <ReactMarkdown>{item.business_context}</ReactMarkdown>
                                        </div>
                                    ) : (
                                        <p className="text-sm text-muted-foreground italic">No business context provided.</p>
                                    )}
                                </CardContent>
                            </Card>
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-sm font-medium">Technical Context</CardTitle>
                                </CardHeader>
                                <CardContent className="max-h-[400px] overflow-y-auto">
                                    {item.technical_context ? (
                                        <div className="text-sm text-slate-700 [&_p]:mb-2 [&_p]:leading-relaxed [&_strong]:text-slate-900 [&_ul]:list-disc [&_ul]:pl-4 [&_ul]:space-y-1 [&_ol]:list-decimal [&_ol]:pl-4 [&_ol]:space-y-1 [&_li]:leading-relaxed [&_code]:bg-slate-100 [&_code]:px-1 [&_code]:py-0.5 [&_code]:rounded [&_code]:text-xs [&_code]:font-mono [&_h1]:text-base [&_h1]:font-semibold [&_h1]:mb-2 [&_h1]:text-slate-900 [&_h2]:text-sm [&_h2]:font-semibold [&_h2]:mb-1.5 [&_h2]:text-slate-900 [&_h3]:text-sm [&_h3]:font-medium [&_h3]:mb-1 [&_h3]:text-slate-800">
                                            <ReactMarkdown>{item.technical_context}</ReactMarkdown>
                                        </div>
                                    ) : (
                                        <p className="text-sm text-muted-foreground italic">No technical context provided.</p>
                                    )}
                                </CardContent>
                            </Card>
                        </div>
                        {item.breaking_change && (
                            <div className="flex items-center gap-3 p-4 bg-red-50 border border-red-200 rounded-lg text-red-800">
                                <AlertTriangle className="h-5 w-5" />
                                <div className="text-sm font-medium">
                                    This item is marked as a BREAKING CHANGE. Ensure all downstream consumers are notified.
                                </div>
                            </div>
                        )}

                        <div className="pt-4">
                            <BuildArtifactPanel
                                roadmapItemId={roadmapItemId!}
                                completenessScore={item.readiness_level === "READY" ? 95 : 45}
                            />
                        </div>
                    </TabsContent>

                    <TabsContent value="contracts" className="pt-4">
                        <div className="grid gap-4">
                            {contracts.length === 0 ? (
                                <div className="p-8 text-center border-2 border-dashed rounded-md space-y-3">
                                    <p className="text-muted-foreground">No contracts defined for this feature yet.</p>
                                    <Button size="sm" variant="outline" onClick={() => setRecModal({
                                        isOpen: true,
                                        title: "Recommend Contract",
                                        description: "Generate an OpenAPI schema based on requirements.",
                                        targetType: "contract"
                                    })}>
                                        <Sparkles className="h-4 w-4 mr-2" />
                                        Recommend Contract
                                    </Button>
                                </div>
                            ) : (
                                <Card>
                                    <CardHeader className="flex flex-row items-center justify-between">
                                        <CardTitle>API Contract Editor</CardTitle>
                                        <div className="flex gap-2">
                                            <Button size="sm" variant="outline" onClick={() => setRecModal({
                                                isOpen: true,
                                                title: "Recommend Contract",
                                                description: "Generate an OpenAPI schema based on requirements.",
                                                targetType: "contract"
                                            })}>
                                                Recommend
                                            </Button>
                                            <Button size="sm" variant="outline" onClick={() => setRecModal({
                                                isOpen: true,
                                                title: "Refine Contract",
                                                description: "Improve existing contract with AI suggestions.",
                                                targetType: "contract",
                                                refineContent: selectedContract
                                            })}>
                                                <Wand2 className="h-4 w-4 mr-2" />
                                                Refine
                                            </Button>
                                            {contracts.map((c, idx) => {
                                                const colorClass = ["bg-purple-100 text-purple-800 border-purple-200",
                                                    "bg-blue-100 text-blue-800 border-blue-200",
                                                    "bg-orange-100 text-orange-800 border-orange-200",
                                                    "bg-pink-100 text-pink-800 border-pink-200"][idx % 4];

                                                return (
                                                    <Button
                                                        key={c.id}
                                                        variant={selectedContract?.id === c.id ? "default" : "outline"}
                                                        size="sm"
                                                        onClick={() => setSelectedContract(c)}
                                                        className="relative"
                                                    >
                                                        {c.contract_type} {c.version}
                                                        <span className={`absolute -top-2 -right-2 text-[10px] px-1.5 py-0 rounded-full border ${colorClass} bg-white shadow-sm`}>
                                                            Link #{idx + 1}
                                                        </span>
                                                    </Button>
                                                );
                                            })}
                                        </div>
                                    </CardHeader>
                                    <CardContent className="h-[500px] border rounded-md p-0 overflow-hidden relative">
                                        {selectedContract && (
                                            <Editor
                                                height="100%"
                                                defaultLanguage="json"
                                                value={JSON.stringify(selectedContract.output_schema, null, 2)}
                                                options={{
                                                    minimap: { enabled: false },
                                                    scrollBeyondLastLine: false,
                                                    fontSize: 14,
                                                }}
                                            />
                                        )}
                                    </CardContent>
                                </Card>
                            )}
                        </div>
                    </TabsContent>

                    <TabsContent value="ai-proposals" className="pt-4">
                        <div className="text-center py-12 border-2 border-dashed rounded-lg text-muted-foreground">
                            No AI proposals found for this item.
                        </div>
                    </TabsContent>

                    <TabsContent value="variables" className="pt-4">
                        <div className="flex justify-end mb-4">
                            <Button size="sm" onClick={() => setRecModal({
                                isOpen: true,
                                title: "Recommend Variables",
                                description: "Identify necessary environment variables and secrets.",
                                targetType: "variable"
                            })}>
                                <Sparkles className="h-4 w-4 mr-2" />
                                Recommend Variables
                            </Button>
                        </div>
                        {variables.length === 0 ? (
                            <div className="text-center py-12 border-2 border-dashed rounded-lg text-muted-foreground">
                                No variables defined.
                            </div>
                        ) : (
                            <div className="grid gap-4">
                                {variables.map((v: any) => {
                                    const contractIndex = contracts.findIndex(c => c.id === v.contract_id);
                                    const colorClass = contractIndex >= 0
                                        ? ["bg-purple-100 text-purple-800 border-purple-200",
                                            "bg-blue-100 text-blue-800 border-blue-200",
                                            "bg-orange-100 text-orange-800 border-orange-200",
                                            "bg-pink-100 text-pink-800 border-pink-200"][contractIndex % 4]
                                        : "bg-gray-100 text-gray-800 border-gray-200";

                                    return (
                                        <Card key={v.id}>
                                            <CardHeader className="py-3">
                                                <div className="flex items-center justify-between">
                                                    <div className="flex items-center gap-3">
                                                        <Variable className="h-4 w-4 text-muted-foreground" />
                                                        <span className="font-mono font-medium">{v.name}</span>
                                                    </div>
                                                    <div className="flex items-center gap-2">
                                                        {contractIndex >= 0 && (
                                                            <span className={`text-xs px-2 py-0.5 rounded-full border ${colorClass}`}>
                                                                Link #{contractIndex + 1}
                                                            </span>
                                                        )}
                                                        <Badge variant="outline">{v.type}</Badge>
                                                    </div>
                                                </div>
                                            </CardHeader>
                                            <CardContent className="py-3 text-sm text-muted-foreground">
                                                {v.description}
                                                {v.default_value && (
                                                    <div className="mt-2 text-xs bg-muted p-2 rounded font-mono">
                                                        Default: {v.default_value}
                                                    </div>
                                                )}
                                            </CardContent>
                                        </Card>
                                    );
                                })}
                            </div>
                        )}
                    </TabsContent>

                    <TabsContent value="reports" className="pt-4">
                        {relevantReports.length === 0 ? (
                            <div className="text-center py-12 border-2 border-dashed rounded-lg text-muted-foreground">
                                No drift reports or snapshot history available for linked contracts.
                            </div>
                        ) : (
                            <div className="grid gap-4">
                                {relevantReports.map((log: any) => {
                                    const report = log.new_data.drift_report;
                                    const isBreaking = report?.breaking_changes && report.breaking_changes.length > 0;
                                    const riskScore = report?.risk_score || 0;
                                    const severity = riskScore > 0.7 ? 'HIGH' : riskScore > 0.3 ? 'MEDIUM' : 'LOW';
                                    const date = new Date(log.created_at).toLocaleString();

                                    return (
                                        <Card key={log.id}>
                                            <CardHeader className="py-3">
                                                <div className="flex justify-between items-start">
                                                    <div>
                                                        <CardTitle className="text-base font-medium flex items-center gap-2">
                                                            {log.entity_id}
                                                            <span className="text-xs font-normal text-muted-foreground">
                                                                {log.action}
                                                            </span>
                                                        </CardTitle>
                                                        <div className="text-xs text-muted-foreground">{date}</div>
                                                    </div>
                                                    <Badge variant={severity === 'HIGH' ? "destructive" : "outline"}>
                                                        {severity} RISK
                                                    </Badge>
                                                </div>
                                            </CardHeader>
                                            <CardContent className="py-3">
                                                {isBreaking ? (
                                                    <div className="text-sm text-red-600 flex items-center gap-2">
                                                        <AlertTriangle className="h-4 w-4" />
                                                        Breaking changes detected
                                                    </div>
                                                ) : (
                                                    <div className="text-sm text-muted-foreground">
                                                        No breaking changes detected. Risk Score: {riskScore.toFixed(2)}
                                                    </div>
                                                )}
                                                {report?.breaking_changes?.length > 0 && (
                                                    <div className="mt-2 bg-red-50 p-2 rounded text-xs text-red-800 font-mono">
                                                        {report.breaking_changes[0].issue}
                                                        {report.breaking_changes.length > 1 && ` (+${report.breaking_changes.length - 1} more)`}
                                                    </div>
                                                )}
                                            </CardContent>
                                        </Card>
                                    );
                                })}
                            </div>
                        )}
                    </TabsContent>
                </Tabs>
            </div >

            {/* Intelligence Panel (Right Panel) */}
            < div className="w-80 border-l bg-slate-50/50 p-4 overflow-y-auto" >
                <IntelligencePanel roadmapItemId={roadmapItemId!} />
            </div >

            <RecommendationModal
                isOpen={recModal.isOpen}
                onClose={() => setRecModal({ ...recModal, isOpen: false })}
                title={recModal.title}
                description={recModal.description}
                targetType={recModal.targetType}
                contextData={item}
                contracts={contracts}
                refineContent={recModal.refineContent}
                onApply={async (result) => {
                    if (!roadmapItemId || !item || !item.project_id) return;

                    try {
                        switch (recModal.targetType) {
                            case "context":
                                await updateItem.mutateAsync({
                                    id: roadmapItemId,
                                    data: {
                                        title: item.title,
                                        status: item.status,
                                        description: item.description,
                                        business_context: result.business_context,
                                        technical_context: result.technical_context
                                    }
                                });
                                break;

                            case "variable":
                                // Result might be the raw array (legacy) or an object with variables + metadata
                                const vars = Array.isArray(result) ? result : (result.variables || []);
                                const selectedId = !Array.isArray(result) ? result.selectedContractId : undefined;

                                if (contracts.length === 0) {
                                    alert("No contracts found to attach variables to. Please create a contract first.");
                                    break;
                                }

                                // Use selected contract from modal, or fallback to first one
                                const targetContract = contracts.find(c => c.id === selectedId) || contracts[0];

                                for (const v of vars) {
                                    await variablesApi.createVariable(item.project_id, {
                                        ...v,
                                        project_id: item.project_id,
                                        contract_id: targetContract.id
                                    });
                                }
                                await queryClient.invalidateQueries({ queryKey: ["variables"] });
                                alert(`Added ${vars.length} variables to contract "${targetContract.contract_type} v${targetContract.version}".`);
                                break;

                            case "contract":
                                // Result is { contract: object, variables: array }
                                // or fallback to just contract object if older prompt used
                                const contractContent = result.contract || result;
                                const contractVars = result.variables || [];

                                // Store the full OpenAPI spec as input_schema so it can
                                // be exported directly as openapi.spec.yaml
                                const newContract = await contractsApi.createContract(item.project_id, {
                                    roadmap_item_id: roadmapItemId,
                                    contract_type: contractContent.contract_type || "REST",
                                    version: contractContent.info?.version || contractContent.version || "1.0.0",
                                    input_schema: contractContent,
                                    output_schema: contractContent.components?.schemas || {},
                                    error_schema: contractContent.components?.responses || {},
                                });
                                await queryClient.invalidateQueries({ queryKey: ["contracts"] });

                                if (contractVars.length > 0) {
                                    for (const v of contractVars) {
                                        await variablesApi.createVariable(item.project_id, {
                                            ...v,
                                            project_id: item.project_id,
                                            contract_id: newContract.id
                                        });
                                    }
                                    await queryClient.invalidateQueries({ queryKey: ["variables"] });
                                    alert(`Contract created and ${contractVars.length} variables added.`);
                                } else {
                                    alert("Contract created.");
                                }
                                break;
                        }
                        setRecModal({ ...recModal, isOpen: false });
                    } catch (err) {
                        console.error("Failed to apply recommendation:", err);
                        alert("Failed to apply changes. See console for details.");
                    }
                }}
            />
        </div >
    );
}

