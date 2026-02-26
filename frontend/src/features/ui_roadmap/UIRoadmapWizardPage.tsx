import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useSaveUIRoadmapItem, useUIRoadmapItem, useUpdateUIRoadmapItem } from "@/hooks/use-ui-roadmap";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ChevronLeft, ChevronRight, Save, Info, Layers, Zap } from "lucide-react";
import { Progress } from "@/components/ui/progress";
import { ComponentTreeEditor } from "./components/ComponentTreeEditor";
import { StateMachineEditor } from "./components/StateMachineEditor";

export function UIRoadmapWizardPage() {
    const { projectId, id } = useParams<{ projectId: string; id?: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const { data: existingItem, isLoading: isLoadingItem } = useUIRoadmapItem(id);
    const saveMutation = useSaveUIRoadmapItem(projectId);
    const updateMutation = useUpdateUIRoadmapItem(projectId);

    const [step, setStep] = useState(1);
    const totalSteps = 4;

    const [formData, setFormData] = useState<any>({
        name: "",
        description: "",
        user_persona: "",
        use_case: "",
        screen_type: "page",
        layout_definition: {},
        component_tree: { type: "Root", children: [] },
        state_machine: { states: { idle: { visual_changes: "Default", interaction_changes: "Enabled", messaging: "Ready" } } },
        backend_bindings: [],
        accessibility_spec: { role: "main", keyboard_nav: "Tab to interactive", focus_management: "Auto", screen_reader_text: "Page content", contrast_compliance: true },
        responsive_spec: { mobile: { columns: 1 }, tablet: { columns: 1 }, desktop: { columns: 2 } },
        validation_rules: {},
        animation_rules: {},
        design_tokens_used: [],
        edge_cases: {},
        test_scenarios: {},
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
                            <ComponentTreeEditor
                                data={formData.component_tree}
                            />
                        </div>
                    )}

                    {step === 3 && (
                        <div className="space-y-4">
                            <StateMachineEditor
                                data={formData.state_machine}
                            />
                        </div>
                    )}

                    {step === 4 && (
                        <div className="space-y-6">
                            <div className="bg-blue-500/5 p-4 rounded-lg border border-blue-500/20 flex gap-4">
                                <Info className="h-5 w-5 text-blue-500 shrink-0 mt-0.5" />
                                <p className="text-sm text-blue-700 leading-relaxed">
                                    Final verification. Our governance engine will check the state machine for completeness, design token adherence, and backend contract binding.
                                </p>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <Card className="bg-card">
                                    <CardHeader className="pb-2"><CardTitle className="text-sm">Accessibility Spec</CardTitle></CardHeader>
                                    <CardContent className="text-xs text-muted-foreground">
                                        ARIA: {formData.accessibility_spec.role}<br />
                                        Keyboard: {formData.accessibility_spec.keyboard_nav}
                                    </CardContent>
                                </Card>
                                <Card className="bg-card">
                                    <CardHeader className="pb-2"><CardTitle className="text-sm">Responsive Engine</CardTitle></CardHeader>
                                    <CardContent className="text-xs text-muted-foreground">
                                        Mobile: {formData.responsive_spec.mobile.columns} col<br />
                                        Desktop: {formData.responsive_spec.desktop.columns} col
                                    </CardContent>
                                </Card>
                            </div>
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

