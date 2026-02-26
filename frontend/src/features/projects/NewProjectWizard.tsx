import React, { useState } from 'react';
import {
    Check,
    Loader2,
    ChevronRight,
    ChevronLeft,
    Rocket,
    Wand2,
    Settings2,
    CheckCircle2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useCreateProject, useTechStackRecommendation, useUpdateProject } from '@/hooks/use-projects';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

type WizardStep = 'purpose' | 'tech_stack' | 'confirm' | 'creating';

interface NewProjectWizardProps {
    workspaceId: string;
    onComplete?: (projectId: string) => void;
    onCancel?: () => void;
}

export function NewProjectWizard({ workspaceId, onComplete, onCancel }: NewProjectWizardProps) {
    const [step, setStep] = useState<WizardStep>('purpose');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [purpose, setPurpose] = useState('');
    const [techStack, setTechStack] = useState<any>({});
    const [reasoning, setReasoning] = useState('');
    const [error, setError] = useState('');
    const [createdProjectId, setCreatedProjectId] = useState<string | null>(null);

    const recommendStack = useTechStackRecommendation();
    const createProject = useCreateProject(workspaceId);

    const handlePurposeSubmit = async () => {
        if (!name.trim()) {
            setError('Project name is required');
            return;
        }
        if (!purpose.trim()) {
            setError('Please describe the purpose of your application');
            return;
        }

        setError('');
        try {
            // Create the project early so we have a valid ID for MCP/Onboarding
            const projectResult = await createProject.mutateAsync({
                name,
                description,
                // @ts-ignore
                purpose,
                tech_stack: {},
                settings: {
                    enable_self_evaluation: true,
                },
                project_type: 'NEW' as any,
            });
            setCreatedProjectId(projectResult.data.id);

            const res = await recommendStack.mutateAsync({ purpose });
            setTechStack(res.recommended_stack);
            setReasoning(res.reasoning);
            setStep('tech_stack');
        } catch (err: any) {
            setError(err.response?.data?.error?.message || 'Failed to get recommendations');
        }
    };

    const updateProject = useUpdateProject(createdProjectId || '');

    const handleCreateProject = async () => {
        if (!createdProjectId) {
            setError('Project ID not found. Please restart the wizard.');
            return;
        }

        setError('');
        setStep('creating');
        try {
            await updateProject.mutateAsync({
                name,
                description,
                // @ts-ignore
                tech_stack: techStack,
                settings: {
                    enable_self_evaluation: true,
                },
            });
            onComplete?.(createdProjectId);
        } catch (err: any) {
            setStep('confirm');
            setError(err.response?.data?.error?.message || 'Failed to finalize project');
        }
    };

    const steps: { key: WizardStep; label: string; number: number }[] = [
        { key: 'purpose', label: 'Context & Purpose', number: 1 },
        { key: 'tech_stack', label: 'Tech Stack', number: 2 },
        { key: 'confirm', label: 'Confirmation', number: 3 },
    ];

    const currentStepIndex = steps.findIndex(s => s.key === step);

    return (
        <div className="space-y-6 max-w-2xl mx-auto">
            {/* Step indicator */}
            <div className="flex items-center justify-between mb-8 overflow-x-auto pb-2">
                {steps.map((s, idx) => (
                    <React.Fragment key={s.key}>
                        <div className="flex items-center gap-2 flex-shrink-0">
                            <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors ${idx < currentStepIndex
                                ? 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/20'
                                : idx === currentStepIndex
                                    ? 'bg-violet-600/10 text-violet-600 border border-violet-600/20'
                                    : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700'
                                }`}>
                                {idx < currentStepIndex ? <Check className="w-4 h-4" /> : s.number}
                            </div>
                            <span className={`text-sm whitespace-nowrap ${idx === currentStepIndex
                                ? 'text-zinc-900 dark:text-white font-bold'
                                : 'text-zinc-500 dark:text-zinc-400'
                                }`}>{s.label}</span>
                        </div>
                        {idx < steps.length - 1 && (
                            <ChevronRight className="w-4 h-4 text-zinc-400 dark:text-zinc-600 mx-2 flex-shrink-0" />
                        )}
                    </React.Fragment>
                ))}
            </div>

            {/* Error Message */}
            {error && (
                <Alert variant="destructive" className="bg-red-500/10 border-red-500/50 text-red-400">
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            )}

            {/* Step Content */}
            {step === 'purpose' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Rocket className="w-5 h-5 text-violet-400" />
                            Define Project Purpose
                        </CardTitle>
                        <CardDescription>
                            Tell us what you want to build. Our AI will recommend a modern tech stack.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="wiz-name">Project Name</Label>
                            <Input
                                id="wiz-name"
                                placeholder="e.g. My Awesome App"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="bg-background border-input"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="wiz-description">Tagline (Optional)</Label>
                            <Input
                                id="wiz-description"
                                placeholder="Short description"
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                className="bg-background border-input"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="wiz-purpose">What are you building?</Label>
                            <Textarea
                                id="wiz-purpose"
                                placeholder="Describe the app's functionality, scale, and specific requirements..."
                                className="min-h-[120px] bg-background border-input"
                                value={purpose}
                                onChange={(e) => setPurpose(e.target.value)}
                            />
                            <p className="text-xs text-muted-foreground">
                                Be specific about scale, security needs, and user interface preferences.
                            </p>
                        </div>
                        <div className="flex justify-between pt-4">
                            <Button variant="ghost" onClick={onCancel}>Cancel</Button>
                            <Button
                                onClick={handlePurposeSubmit}
                                disabled={recommendStack.isPending}
                                className="bg-violet-600 hover:bg-violet-700"
                            >
                                {recommendStack.isPending ? (
                                    <><Loader2 className="w-4 h-4 mr-2 animate-spin" /> Analyzing...</>
                                ) : (
                                    <><Wand2 className="w-4 h-4 mr-2" /> Recommend Stack</>
                                )}
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {step === 'tech_stack' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Settings2 className="w-5 h-5 text-emerald-400" />
                            Recommended Tech Stack
                        </CardTitle>
                        <CardDescription>
                            AI-suggested foundation based on your project goals.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                        <div className="bg-violet-500/5 border border-violet-500/20 rounded-lg p-4">
                            <h4 className="text-sm font-semibold text-violet-600 dark:text-violet-300 mb-2">AI Reasoning</h4>
                            <p className="text-xs text-muted-foreground italic">"{reasoning}"</p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {Object.entries(techStack).map(([category, details]: [string, any]) => (
                                <div key={category} className="p-3 rounded-lg bg-muted/30 border border-border">
                                    <h5 className="text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-2">{category}</h5>
                                    <div className="text-sm text-foreground font-medium">{typeof details === 'string' ? details : JSON.stringify(details)}</div>
                                </div>
                            ))}
                        </div>

                        <div className="flex justify-between pt-4">
                            <Button variant="outline" onClick={() => setStep('purpose')}>
                                <ChevronLeft className="w-4 h-4 mr-1" /> Back
                            </Button>
                            <Button
                                onClick={() => setStep('confirm')}
                                className="bg-emerald-600 hover:bg-emerald-700"
                            >
                                Looks Good <ChevronRight className="w-4 h-4 ml-1" />
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {step === 'confirm' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle2 className="w-5 h-5 text-blue-400" />
                            Ready to Orchestrate?
                        </CardTitle>
                        <CardDescription>
                            Finalize your project settings.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="rounded-lg border border-border divide-y divide-border overflow-hidden">
                            <div className="p-3 bg-muted/30 grid grid-cols-3">
                                <span className="text-xs text-muted-foreground">Project Name</span>
                                <span className="text-sm font-medium text-foreground col-span-2">{name}</span>
                            </div>
                            <div className="p-3 bg-muted/30 grid grid-cols-3">
                                <span className="text-xs text-muted-foreground">Framework</span>
                                <span className="text-sm font-medium text-foreground col-span-2">{techStack.Frontend || techStack.Framework || 'Default'}</span>
                            </div>
                            <div className="p-3 bg-muted/30 grid grid-cols-3">
                                <span className="text-xs text-muted-foreground">Database</span>
                                <span className="text-sm font-medium text-foreground col-span-2">{techStack.Database || 'N/A'}</span>
                            </div>
                        </div>

                        <div className="flex justify-between pt-4">
                            <Button variant="outline" onClick={() => setStep('tech_stack')}>
                                <ChevronLeft className="w-4 h-4 mr-1" /> Back
                            </Button>
                            <Button
                                onClick={handleCreateProject}
                                disabled={createProject.isPending}
                                className="bg-violet-600 hover:bg-violet-700"
                            >
                                {createProject.isPending ? (
                                    <><Loader2 className="w-4 h-4 mr-2 animate-spin" /> Creating...</>
                                ) : (
                                    'Create & Initialize Project'
                                )}
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {step === 'creating' && (
                <div className="flex flex-col items-center justify-center space-y-4 min-h-[300px]">
                    <div className="relative">
                        <div className="w-20 h-20 rounded-full border-2 border-violet-500/20 animate-pulse"></div>
                        <Loader2 className="w-10 h-10 text-violet-400 animate-spin absolute top-1/2 left-1/2 -mt-5 -ml-5" />
                    </div>
                    <div className="text-center">
                        <h3 className="text-lg font-bold text-foreground">Initializing Universe</h3>
                        <p className="text-sm text-muted-foreground mt-1">Deploying governance contracts and setting up workspace...</p>
                    </div>
                </div>
            )}
        </div>
    );
}
