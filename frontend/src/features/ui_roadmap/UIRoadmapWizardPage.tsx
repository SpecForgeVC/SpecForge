import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useSaveUIRoadmapItem, useUIRoadmapItem, useUpdateUIRoadmapItem } from "@/hooks/use-ui-roadmap";
import { uiRoadmapApi } from "@/api/ui_roadmap";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ChevronLeft, ChevronRight, Save, Info, Layers, Zap, Sparkles, AlertTriangle, ShieldAlert, Loader2, CheckCircle2 } from "lucide-react";
import { Progress } from "@/components/ui/progress";
import { ComponentTreeEditor } from "./components/ComponentTreeEditor";
import { StateMachineEditor } from "./components/StateMachineEditor";
import { Badge } from "@/components/ui/badge";

export function UIRoadmapWizardPage() {
    const { projectId, id } = useParams<{ projectId: string; id?: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const { data: existingItem, isLoading: isLoadingItem } = useUIRoadmapItem(id);
    const saveMutation = useSaveUIRoadmapItem(projectId);
    const updateMutation = useUpdateUIRoadmapItem(projectId);

    const [step, setStep] = useState(1);
    const [isGeneratingAI, setIsGeneratingAI] = useState(false);
    const [isCheckingCompliance, setIsCheckingCompliance] = useState(false);
    const [complianceIssues, setComplianceIssues] = useState<any[]>([]);
    const totalSteps = 4;

    const [formData, setFormData] = useState<any>({
        name: "",
        description: "",
        user_persona: "",
        use_case: "",
        screen_type: "page",
        layout_definition: { type: "flex", direction: "column", gap: "1rem" },
        component_tree: { type: "Root", children: [] },
        state_machine: {
            states: {
                idle: { visual_changes: "Default", interaction_changes: "Enabled", messaging: "Ready" },
                loading: { visual_changes: "Skeleton", interaction_changes: "Disabled", messaging: "Loading..." },
                success: { visual_changes: "Success UI", interaction_changes: "Enabled", messaging: "Success" },
                error: { visual_changes: "Error UI", interaction_changes: "Retry enabled", messaging: "Error" },
                empty: { visual_changes: "Empty UI", interaction_changes: "Enabled", messaging: "No data" },
                disabled: { visual_changes: "Greyed out", interaction_changes: "Disabled", messaging: "N/A" }
            }
        },
        backend_bindings: [{ endpoint: "/api/v1/context", method: "GET" }],
        accessibility_spec: { role: "main", keyboard_nav: "Tab to interactive", focus_management: "Auto", screen_reader_text: "Page content", contrast_compliance: true },
        responsive_spec: { mobile: { columns: 1 }, tablet: { columns: 1 }, desktop: { columns: 2 } },
        validation_rules: { global: "Standard compliance required" },
        animation_rules: {},
        design_tokens_used: [],
        edge_cases: {},
        test_scenarios: { smoke_test: "Verify basic page load" },
    });

    useEffect(() => {
        if (existingItem) {
            setFormData(existingItem);
        }
    }, [existingItem]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setFormData((prev: any) => ({ ...prev, [name]: value }));
    };

    const handleSelectChange = (name: string, value: string) => {
        setFormData((prev: any) => ({ ...prev, [name]: value }));
    };

    const handleGenerateTree = async () => {
        if (!formData.name || !formData.description) {
            alert("Please provide at least a name and description in Step 1 first.");
            return;
        }
        setIsGeneratingAI(true);
        try {
            const tree = await uiRoadmapApi.recommendTree(projectId!, formData);
            setFormData((prev: any) => ({ ...prev, component_tree: tree }));
        } catch (error) {
            console.error("Failed to generate AI tree", error);
        } finally {
            setIsGeneratingAI(false);
        }
    };


    const handleGenerateStateMachine = async () => {
        if (!formData.name || !formData.description) {
            alert("Please provide at least a name and description in Step 1 first.");
            return;
        }
        setIsGeneratingAI(true);
        try {
            const sm = await uiRoadmapApi.recommendStateMachine(projectId!, formData);
            setFormData((prev: any) => ({ ...prev, state_machine: sm }));
        } catch (error) {
            console.error("Failed to generate AI state machine", error);
        } finally {
            setIsGeneratingAI(false);
        }
    };

    const handleCheckCompliance = async () => {
        setIsCheckingCompliance(true);
        try {
            const issues = await uiRoadmapApi.checkCompliance(projectId!, formData);
            setComplianceIssues(issues);
        } catch (error) {
            console.error("Failed to check compliance", error);
        } finally {
            setIsCheckingCompliance(false);
        }
    };

    const handleRepairDrift = async () => {
        if (complianceIssues.length === 0) return;
        setIsGeneratingAI(true);
        try {
            const fixedItem = await uiRoadmapApi.recommendFix(projectId!, formData, complianceIssues);
            setFormData(fixedItem);
            setComplianceIssues([]); // Reset issues after repair
            alert("Specification repaired by AI!");
        } catch (error) {
            console.error("Failed to repair drift", error);
        } finally {
            setIsGeneratingAI(false);
        }
    };

    const handleSave = async () => {
        try {
            if (isEditing) {
                await updateMutation.mutateAsync({ id: id!, item: formData });
            } else {
                await saveMutation.mutateAsync(formData);
            }
            navigate(`/projects/${projectId}/ui-roadmap`);
        } catch (error: any) {
            console.error(error.response?.data?.error || "Failed to save UI Roadmap Item");
        }
    };

    if (isLoadingItem) return <div className="p-8 text-center mt-20">Loading specification...</div>;

    return (
        <div className="p-8 max-w-4xl mx-auto">
            <div className="mb-8 space-y-4">
                <div className="flex items-center justify-between">
                    <Button variant="ghost" onClick={() => navigate(-1)} className="gap-2 -ml-2">
                        <ChevronLeft className="h-4 w-4" /> Back to List
                    </Button>
                    <div className="text-sm font-medium text-muted-foreground uppercase tracking-widest">
                        Step {step} of {totalSteps}
                    </div>
                </div>
                <Card className="border-none shadow-none bg-muted/50">
                    <CardContent className="p-0">
                        <Progress value={(step / totalSteps) * 100} className="h-1 rounded-none bg-primary/10" />
                    </CardContent>
                </Card>
            </div>

            <Card className="shadow-lg border-primary/10">
                <CardHeader className="border-b bg-muted/30 pb-6">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-primary rounded-lg text-primary-foreground">
                            {step === 1 && <Info className="h-5 w-5" />}
                            {step === 2 && <Layers className="h-5 w-5" />}
                            {step === 3 && <Zap className="h-5 w-5" />}
                            {step === 4 && <Save className="h-5 w-5" />}
                        </div>
                        <div>
                            <CardTitle>{getStepTitle(step)}</CardTitle>
                            <p className="text-sm text-muted-foreground mt-1">{getStepDescription(step)}</p>
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="pt-8 min-h-[400px]">
                    {step === 1 && (
                        <div className="space-y-6">
                            <div className="grid gap-2">
                                <Label htmlFor="name">Feature Name</Label>
                                <Input id="name" name="name" value={formData.name} onChange={handleInputChange} placeholder="e.g. User Profile Dashboard" />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="screen_type">Screen Type</Label>
                                <Select value={formData.screen_type} onValueChange={(v) => handleSelectChange("screen_type", v)}>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="page">Full Page</SelectItem>
                                        <SelectItem value="modal">Modal / Overlay</SelectItem>
                                        <SelectItem value="component">Reusable Component</SelectItem>
                                        <SelectItem value="layout">Layout Shell</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="description">Business Context / Description</Label>
                                <Textarea id="description" name="description" value={formData.description} onChange={handleInputChange} rows={3} placeholder="Describe the purpose and goals of this UI item..." />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div className="grid gap-2">
                                    <Label htmlFor="user_persona">User Persona</Label>
                                    <Input id="user_persona" name="user_persona" value={formData.user_persona} onChange={handleInputChange} placeholder="e.g. System Admin" />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="use_case">Primary Use Case</Label>
                                    <Input id="use_case" name="use_case" value={formData.use_case} onChange={handleInputChange} placeholder="e.g. Viewing audit logs" />
                                </div>
                            </div>
                        </div>
                    )}

                    {step === 2 && (
                        <div className="space-y-4">
                            <div className="flex justify-between items-center bg-primary/5 p-4 rounded-xl border border-primary/20 shadow-sm border-dashed">
                                <div className="flex items-center gap-3">
                                    <div className="p-2 bg-primary/10 rounded-full">
                                        <Sparkles className="h-4 w-4 text-primary animate-pulse" />
                                    </div>
                                    <div>
                                        <h4 className="text-sm font-bold tracking-tight">AI Accelerator</h4>
                                        <p className="text-[11px] text-muted-foreground leading-tight">
                                            Build initial hierarchy based on Step 1 context.
                                        </p>
                                    </div>
                                </div>
                                <Button
                                    size="sm"
                                    onClick={handleGenerateTree}
                                    disabled={isGeneratingAI}
                                    variant="outline"
                                    className="h-8 border-primary/30 hover:bg-primary/10 hover:text-primary transition-all duration-300"
                                >
                                    {isGeneratingAI ? (
                                        <>
                                            <Zap className="mr-2 h-3 w-3 animate-spin" /> Analyzing...
                                        </>
                                    ) : (
                                        <>
                                            <Sparkles className="mr-2 h-3 w-3" /> Generate with AI
                                        </>
                                    )}
                                </Button>
                            </div>
                            <ComponentTreeEditor
                                data={formData.component_tree}
                                onChange={(tree) => setFormData((prev: any) => ({ ...prev, component_tree: tree }))}
                            />
                        </div>
                    )}

                    {step === 3 && (
                        <div className="space-y-4">
                            <div className="flex justify-between items-center bg-primary/5 p-4 rounded-xl border border-primary/20 shadow-sm border-dashed">
                                <div className="flex items-center gap-3">
                                    <div className="p-2 bg-primary/10 rounded-full">
                                        <Sparkles className="h-4 w-4 text-primary animate-pulse" />
                                    </div>
                                    <div>
                                        <h4 className="text-sm font-bold tracking-tight">Logic Accelerator</h4>
                                        <p className="text-[11px] text-muted-foreground leading-tight">
                                            Generate states and transitions based on the components and persona.
                                        </p>
                                    </div>
                                </div>
                                <Button
                                    size="sm"
                                    onClick={handleGenerateStateMachine}
                                    disabled={isGeneratingAI}
                                    variant="outline"
                                    className="h-8 border-primary/30 hover:bg-primary/10 hover:text-primary transition-all duration-300"
                                >
                                    {isGeneratingAI ? (
                                        <>
                                            <Zap className="mr-2 h-3 w-3 animate-spin" /> Analyzing...
                                        </>
                                    ) : (
                                        <>
                                            <Sparkles className="mr-2 h-3 w-3" /> Generate with AI
                                        </>
                                    )}
                                </Button>
                            </div>
                            <StateMachineEditor
                                data={formData.state_machine}
                                onChange={(sm) => setFormData((prev: any) => ({ ...prev, state_machine: sm }))}
                            />
                        </div>
                    )}

                    {step === 4 && (
                        <div className="space-y-6">
                            <div className="bg-blue-500/5 p-4 rounded-lg border border-blue-500/20 flex gap-4">
                                <Info className="h-5 w-5 text-blue-500 shrink-0 mt-0.5" />
                                <div className="space-y-1">
                                    <p className="text-sm text-blue-700 font-bold">Governance & Compliance Engine</p>
                                    <p className="text-xs text-blue-600 leading-relaxed">
                                        Final verification. Our governance engine will check the state machine for completeness, design token adherence, and backend contract binding.
                                    </p>
                                </div>
                            </div>

                            <div className="flex justify-between items-center bg-muted/30 p-4 rounded-lg border">
                                <div className="space-y-0.5">
                                    <h4 className="text-sm font-semibold">Compliance Status</h4>
                                    <p className="text-[10px] text-muted-foreground">Run a deep check against backend contracts and rules.</p>
                                </div>
                                <Button
                                    size="sm"
                                    onClick={handleCheckCompliance}
                                    disabled={isCheckingCompliance}
                                    variant="outline"
                                    className="gap-2"
                                >
                                    {isCheckingCompliance ? (
                                        <>
                                            <Loader2 className="h-3 w-3 animate-spin" /> Checking...
                                        </>
                                    ) : (
                                        <>
                                            <ShieldAlert className="h-3 w-3" /> Check Compliance
                                        </>
                                    )}
                                </Button>
                            </div>

                            {complianceIssues.length > 0 && (
                                <div className="space-y-4 animate-in slide-in-from-top duration-300">
                                    <div className="rounded-xl border border-red-200 bg-red-500/5 overflow-hidden">
                                        <div className="bg-red-500/10 p-3 flex justify-between items-center border-b border-red-200">
                                            <div className="flex items-center gap-2 text-red-700 font-bold text-xs uppercase tracking-wider">
                                                <AlertTriangle className="h-3.5 w-3.5" /> {complianceIssues.length} Problems Identified
                                            </div>
                                            <Button
                                                size="sm"
                                                variant="destructive"
                                                className="h-7 text-[10px] gap-2 shadow-sm"
                                                onClick={handleRepairDrift}
                                                disabled={isGeneratingAI}
                                            >
                                                {isGeneratingAI ? (
                                                    <Loader2 className="h-3 w-3 animate-spin" />
                                                ) : (
                                                    <Sparkles className="h-3 w-3" />
                                                )}
                                                Repair with AI
                                            </Button>
                                        </div>
                                        <div className="p-3 space-y-2 max-h-[200px] overflow-y-auto">
                                            {complianceIssues.map((issue, idx) => (
                                                <div key={idx} className="flex items-start gap-3 p-2 bg-white rounded border border-red-100 shadow-sm">
                                                    <Badge variant="outline" className="text-[9px] uppercase h-5 px-1 bg-red-50 text-red-600 border-red-200">
                                                        {issue.type}
                                                    </Badge>
                                                    <div className="space-y-0.5">
                                                        <p className="text-[10px] font-bold text-slate-800">{issue.field}</p>
                                                        <p className="text-[10px] text-slate-600 italic">"{issue.description}"</p>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </div>
                            )}

                            {complianceIssues.length === 0 && !isCheckingCompliance && (
                                <div className="grid grid-cols-2 gap-4">
                                    <Card className="bg-card">
                                        <CardHeader className="pb-2"><CardTitle className="text-sm text-primary flex items-center gap-2"><CheckCircle2 className="h-4 w-4" /> Accessibility Spec</CardTitle></CardHeader>
                                        <CardContent className="text-xs text-muted-foreground">
                                            ARIA-Role: {formData.accessibility_spec.role}<br />
                                            Keyboard: {formData.accessibility_spec.keyboard_nav}
                                        </CardContent>
                                    </Card>
                                    <Card className="bg-card">
                                        <CardHeader className="pb-2"><CardTitle className="text-sm text-primary flex items-center gap-2"><CheckCircle2 className="h-4 w-4" /> Responsive Engine</CardTitle></CardHeader>
                                        <CardContent className="text-xs text-muted-foreground">
                                            Mobile: {formData.responsive_spec.mobile.columns} col<br />
                                            Desktop: {formData.responsive_spec.desktop.columns} col
                                        </CardContent>
                                    </Card>
                                </div>
                            )}
                        </div>
                    )}
                </CardContent>
                <CardFooter className="flex justify-between border-t bg-muted/30 py-6">
                    <Button
                        variant="ghost"
                        onClick={() => setStep((s) => Math.max(1, s - 1))}
                        disabled={step === 1}
                    >
                        <ChevronLeft className="mr-2 h-4 w-4" /> Previous
                    </Button>
                    {step < totalSteps ? (
                        <Button onClick={() => setStep((s) => Math.min(totalSteps, s + 1))}>
                            Next Step <ChevronRight className="ml-2 h-4 w-4" />
                        </Button>
                    ) : (
                        <Button onClick={handleSave} className="bg-primary hover:bg-primary/90 shadow-md">
                            <Save className="mr-2 h-4 w-4" /> Complete Specification
                        </Button>
                    )}
                </CardFooter>
            </Card>
        </div>
    );
}

function getStepTitle(step: number) {
    switch (step) {
        case 1: return "Identity & Context";
        case 2: return "Visual Hierarchy";
        case 3: return "Interactive States";
        case 4: return "Governance & Compliance";
        default: return "";
    }
}

function getStepDescription(step: number) {
    switch (step) {
        case 1: return "Define the persona, use case, and type of UI feature.";
        case 2: return "Construct the deterministic component tree structure.";
        case 3: return "Define transitions and state-based visual mutations.";
        case 4: return "Enforce accessibility rules and verify backend bindings.";
        default: return "";
    }
}

