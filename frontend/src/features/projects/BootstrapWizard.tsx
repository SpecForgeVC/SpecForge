import { useState, useEffect, useRef } from 'react';
import {
    Copy,
    Check,
    Loader2,
    ChevronRight,
    Terminal,
    Zap,
    ShieldCheck,
    ArrowLeft,
    PlusCircle,
    Activity,
    RefreshCw,
    Lock
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useGenerateBootstrapPrompt, useIngestBootstrap } from '@/hooks/use-bootstrap';
import { useMCPTokens, useGenerateMCPToken } from '@/hooks/use-mcp';
import type { BootstrapIngestResponse, ImportSession } from '@/api/bootstrap';
import { bootstrapApi } from '@/api/bootstrap';
import { useMutation } from '@tanstack/react-query';

type WizardStep = 'selection' | 'mcp_connect' | 'waiting_snapshot' | 'generate' | 'copy' | 'paste' | 'complete';
type IDEType = 'cursor' | 'anti-gravity' | 'claude';

interface BootstrapWizardProps {
    projectId: string;
    onComplete?: (result: BootstrapIngestResponse) => void;
}

export function BootstrapWizard({ projectId: initialProjectId, onComplete }: BootstrapWizardProps) {
    const [projectId, setProjectId] = useState(initialProjectId);
    const [projectName, setProjectName] = useState('');
    const [activeSubStep, setActiveSubStep] = useState<'config' | 'instructions'>('config');
    const [step, setStep] = useState<WizardStep>('selection');
    const [prompt, setPrompt] = useState('');
    const [copied, setCopied] = useState(false);
    const [pastedJson, setPastedJson] = useState('');
    const [parseError, setParseError] = useState('');
    const [selectedIDE, setSelectedIDE] = useState<IDEType>('cursor');
    const [result, setResult] = useState<BootstrapIngestResponse | null>(null);
    const [importSession, setImportSession] = useState<ImportSession | null>(null);
    const [instructionText, setInstructionText] = useState('Loading instructions...');

    const generatePrompt = useGenerateBootstrapPrompt();
    const ingestBootstrap = useIngestBootstrap(projectId);
    const hasAttemptedTokenGen = useRef(false);
    const [rawToken, setRawToken] = useState<string | null>(null);

    const { data: tokens, isLoading: isLoadingTokens } = useMCPTokens(projectId);
    const generateToken = useGenerateMCPToken(projectId);

    // Dynamic token injection
    const activeToken = rawToken || (tokens && tokens.length > 0 ? tokens[0].token_prefix : null) || 'YOUR_API_TOKEN';

    // Fetch instructions on mount
    useEffect(() => {
        fetch('/IPCP_INSTRUCTIONS.md')
            .then(res => res.text())
            .then(text => setInstructionText(text))
            .catch(() => setInstructionText('Failed to load instructions. Please ensure IPCP_INSTRUCTIONS.md exists in the public folder.'));
    }, []);

    const createProject = useMutation({
        mutationFn: (data: { name: string }) => bootstrapApi.createProject(data),
        onSuccess: (data: { id: string }) => {
            setProjectId(data.id);
        }
    });

    useEffect(() => {
        if (step === 'mcp_connect' && !isLoadingTokens && !rawToken && !generateToken.isPending && !hasAttemptedTokenGen.current) {
            if (tokens && tokens.length > 0) {
                hasAttemptedTokenGen.current = true;
                return;
            }

            hasAttemptedTokenGen.current = true;
            generateToken.mutate(undefined, {
                onSuccess: (data) => {
                    setRawToken(data.token_raw);
                }
            });
        }
    }, [step, tokens, isLoadingTokens, rawToken, generateToken, projectId]);

    useEffect(() => {
        let interval: number;
        if (projectId && (step === 'mcp_connect' || step === 'waiting_snapshot')) {
            interval = window.setInterval(() => {
                bootstrapApi.getLatestImportSession(projectId).then((session: ImportSession) => {
                    if (session) {
                        setImportSession(session);
                        if (session.status === 'complete') {
                            bootstrapApi.getLatestSnapshot(projectId).then((data: any) => {
                                if (data && data.snapshot_json) {
                                    setResult({
                                        snapshot: data,
                                        scores: data.scores,
                                        confidence: data.confidence,
                                        warnings: []
                                    });
                                    setStep('complete');
                                    onComplete?.(data as any);
                                }
                            });
                        }
                    }
                }).catch(() => { });
            }, 3000);
        }
        return () => {
            if (interval) window.clearInterval(interval);
        };
    }, [step, projectId, onComplete]);

    // Auto-redirect to dashboard when session is completed and locked
    useEffect(() => {
        if (importSession?.status === 'complete' && importSession?.locked) {
            // we could fetch the project/workspace but for now we'll use a relative path
            // the dashboard expects /workspaces/:workspaceId/projects/:projectId/dashboard
            // for simplicity since we don't have workspaceId here, we'll try to redirect
            // to the project page which will handle the layout
            window.location.href = `/projects/${projectId}`;
        }
    }, [importSession, projectId]);

    const handleGenerateToken = () => {
        if (!generateToken.isPending) {
            generateToken.mutate(undefined, {
                onSuccess: (data) => {
                    setRawToken(data.token_raw);
                }
            });
        }
    };

    const handleGeneratePrompt = async () => {
        try {
            const res = await generatePrompt.mutateAsync(projectId);
            setPrompt(res.prompt);
            setStep('copy');
        } catch (err) {
            console.error(err);
        }
    };

    const handleCopy = async (text: string) => {
        await navigator.clipboard.writeText(text);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const handleIngest = async () => {
        setParseError('');
        try {
            const parsed = JSON.parse(pastedJson);
            const res = await ingestBootstrap.mutateAsync(parsed);
            setResult(res);
            setStep('complete');
            onComplete?.(res);
        } catch (err: any) {
            setParseError(err.message || 'Failed to ingest bootstrap data.');
        }
    };

    if (!projectId) {
        return (
            <div className="space-y-6">
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <PlusCircle className="w-5 h-5 text-blue-400" />
                            Create New Project
                        </CardTitle>
                        <CardDescription>
                            Initialize a new project intelligence container.
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Project Name</label>
                            <Input
                                placeholder="e.g. My Awesome App"
                                value={projectName}
                                onChange={(e) => setProjectName(e.target.value)}
                            />
                        </div>
                        <Button
                            className="w-full bg-blue-600 hover:bg-blue-700"
                            onClick={() => createProject.mutate({ name: projectName })}
                            disabled={createProject.isPending || !projectName.trim()}
                        >
                            {createProject.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Create Project'}
                        </Button>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {step === 'selection' && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <Card
                        className="bg-card border-border hover:border-violet-500/50 cursor-pointer transition-all group shadow-sm"
                        onClick={() => setStep('mcp_connect')}
                    >
                        <CardHeader>
                            <div className="w-10 h-10 rounded-full bg-violet-500/10 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform">
                                <Zap className="w-5 h-5 text-violet-400" />
                            </div>
                            <CardTitle>Automated Onboarding</CardTitle>
                            <CardDescription>
                                Connect via MCP Server for real-time artifact synchronization.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Button className="w-full mt-6 bg-violet-600 hover:bg-violet-700">Connect MCP Server</Button>
                        </CardContent>
                    </Card>

                    <Card
                        className="bg-card border-border hover:border-blue-500/50 cursor-pointer transition-all group shadow-sm"
                        onClick={() => setStep('generate')}
                    >
                        <CardHeader>
                            <div className="w-10 h-10 rounded-full bg-blue-500/10 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform">
                                <Terminal className="w-5 h-5 text-blue-400" />
                            </div>
                            <CardTitle>Manual IDE Analysis</CardTitle>
                            <CardDescription>
                                Run a custom analysis prompt in your IDE and paste the results.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <Button variant="outline" className="w-full mt-6">Generate Prompt</Button>
                        </CardContent>
                    </Card>
                </div>
            )}

            {(step === 'mcp_connect' || step === 'waiting_snapshot') && (
                <Card className="bg-card border-border shadow-sm border-violet-500/20">
                    <CardHeader className="border-b border-border/50 pb-4">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Activity className="w-5 h-5 text-violet-400" />
                                <CardTitle>Collecting Project Intelligence</CardTitle>
                            </div>
                            <Button variant="ghost" size="sm" onClick={() => { setStep('selection'); setActiveSubStep('config'); }}>
                                <ArrowLeft className="w-4 h-4 mr-1" /> Back
                            </Button>
                        </div>
                        <CardDescription>
                            {activeSubStep === 'config' ? 'Step 1: Connect your IDE to SpecForge via MCP.' : 'Step 2: Instruct your IDE to catalogue the project.'}
                        </CardDescription>
                    </CardHeader>
                    <CardContent className="pt-6 space-y-6">
                        {importSession && (
                            <div className="space-y-4">
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="bg-muted/30 p-3 rounded-md border border-border/50">
                                        <div className="text-[10px] text-muted-foreground uppercase mb-1">Completeness Score</div>
                                        <div className="flex items-end gap-2">
                                            <span className="text-2xl font-bold text-violet-400">
                                                {importSession.completeness_score}%
                                            </span>
                                        </div>
                                    </div>
                                    <div className="bg-muted/30 p-3 rounded-md border border-border/50">
                                        <div className="text-[10px] text-muted-foreground uppercase mb-1">Session Status</div>
                                        <div className="flex items-center gap-2 text-xs font-semibold">
                                            <RefreshCw className="w-3 h-3 text-violet-400" />
                                            Iteration {importSession.iteration_count}
                                        </div>
                                    </div>
                                </div>
                                <div className="p-4 bg-violet-500/10 border border-violet-500/20 rounded-lg text-sm text-foreground">
                                    <p><strong>Iterative Flow Active:</strong> SpecForge will prompt the IDE assistant to gather missing sections iteratively.</p>
                                </div>
                            </div>
                        )}

                        {!importSession && (
                            <div className="space-y-6">
                                {activeSubStep === 'config' ? (
                                    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-300">
                                        <div className="space-y-3">
                                            <div className="flex items-center gap-2">
                                                <div className="w-5 h-5 rounded-full bg-violet-500/20 flex items-center justify-center text-[10px] font-bold text-violet-400">1</div>
                                                <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Select Your IDE</h4>
                                            </div>
                                            <div className="grid grid-cols-3 gap-2">
                                                {(['cursor', 'anti-gravity', 'claude'] as IDEType[]).map((id) => (
                                                    <button
                                                        key={id}
                                                        onClick={() => setSelectedIDE(id)}
                                                        className={`p-2 rounded-lg border-2 transition-all flex flex-col items-center gap-1.5 ${selectedIDE === id ? 'border-violet-600 bg-violet-500/10' : 'border-border bg-muted/30 hover:border-violet-500/30'}`}
                                                    >
                                                        <span className="text-[9px] font-bold uppercase tracking-wider">{id.replace('-', ' ')}</span>
                                                    </button>
                                                ))}
                                            </div>
                                        </div>

                                        <div className="space-y-3">
                                            <div className="flex items-center gap-2">
                                                <div className="w-5 h-5 rounded-full bg-violet-500/20 flex items-center justify-center text-[10px] font-bold text-violet-400">2</div>
                                                <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Copy MCP Configuration</h4>
                                            </div>
                                            <div className="relative group">
                                                <div className="bg-muted/50 p-2 flex items-center justify-between border-b border-input rounded-t-md">
                                                    <span className="text-[10px] uppercase font-bold text-muted-foreground px-2">Configuration (JSON)</span>
                                                    <div className="flex items-center gap-2">
                                                        <Button
                                                            size="sm"
                                                            variant="ghost"
                                                            className="h-6 text-[10px] text-violet-400"
                                                            onClick={handleGenerateToken}
                                                        >
                                                            Regenerate Token
                                                        </Button>
                                                    </div>
                                                </div>
                                                <pre className="p-4 bg-background border border-input rounded-b-md text-[10px] font-mono text-violet-400 overflow-x-auto leading-relaxed max-h-48">
                                                    {JSON.stringify(generatePreviewConfig(selectedIDE, projectId, activeToken, 8081), null, 2)}
                                                </pre>
                                                <Button
                                                    size="icon"
                                                    variant="outline"
                                                    className="absolute top-2 right-2 bg-background/50 backdrop-blur"
                                                    onClick={() => handleCopy(JSON.stringify(generatePreviewConfig(selectedIDE, projectId, activeToken, 8081), null, 2))}
                                                >
                                                    {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                                                </Button>
                                            </div>
                                        </div>

                                        <div className="flex justify-end pt-2">
                                            <Button
                                                className="bg-violet-600 hover:bg-violet-700"
                                                onClick={() => setActiveSubStep('instructions')}
                                            >
                                                Next: Prompt & Instructions <ChevronRight className="w-4 h-4 ml-1" />
                                            </Button>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="space-y-6 animate-in fade-in slide-in-from-right-2 duration-300">
                                        <div className="space-y-4">
                                            <div className="border border-input rounded-md overflow-hidden">
                                                <div className="bg-muted/50 p-2 flex items-center justify-between border-b border-input">
                                                    <div className="flex items-center gap-2 px-2">
                                                        <ShieldCheck className="w-3 h-3 text-violet-400" />
                                                        <span className="text-xs font-semibold">IPCP_INSTRUCTIONS.md</span>
                                                    </div>
                                                    <Button
                                                        size="sm"
                                                        variant="ghost"
                                                        onClick={() => handleCopy(instructionText)}
                                                    >
                                                        {copied ? <Check className="h-3 w-3 mr-1" /> : <Copy className="h-3 w-3 mr-1" />}
                                                        Copy Instructions
                                                    </Button>
                                                </div>
                                                <div className="p-3 bg-background max-h-48 overflow-y-auto text-[10px] font-mono whitespace-pre-wrap text-muted-foreground text-left">
                                                    {instructionText}
                                                </div>
                                            </div>

                                            <div className="relative group">
                                                <div className="bg-muted/50 p-2 flex items-center justify-between border-b border-border rounded-t-md">
                                                    <div className="flex items-center gap-2 px-2">
                                                        <Lock className="w-3 h-3 text-violet-400" />
                                                        <span className="text-[10px] uppercase font-bold text-muted-foreground">MANDATORY AI Prompt</span>
                                                    </div>
                                                </div>
                                                <pre className="bg-background border border-input p-4 rounded-b-md text-xs font-mono text-violet-300 break-words whitespace-pre-wrap text-left">
                                                    {`MANDATORY: You must follow the SpecForge IPCP protocol. \n1. Read the IPCP_INSTRUCTIONS.md file in the project root.\n2. Initialize project ${projectId} using mcp_specforge_init_project_import.\n3. Scaffold .specforge/ and submit snapshots iteratively until 100% completeness.\nDO NOT STOP until the tool returns 'status: complete'.`}
                                                </pre>
                                                <Button
                                                    size="icon"
                                                    variant="outline"
                                                    className="absolute top-2 right-2 bg-background/50 backdrop-blur"
                                                    onClick={() => handleCopy(`MANDATORY: You must follow the SpecForge IPCP protocol.\n1. Read the IPCP_INSTRUCTIONS.md file in the project root.\n2. Initialize project ${projectId} using mcp_specforge_init_project_import.\n3. Scaffold .specforge/ and submit snapshots iteratively until 100% completeness.\nDO NOT STOP until the tool returns 'status: complete'.`)}
                                                >
                                                    {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                                                </Button>
                                            </div>
                                        </div>

                                        <div className="flex items-center justify-between pt-2">
                                            <Button variant="ghost" onClick={() => setActiveSubStep('config')}>
                                                <ArrowLeft className="w-4 h-4 mr-1" /> Back to Config
                                            </Button>
                                            <p className="text-xs text-muted-foreground flex items-center gap-2">
                                                <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                                                Waiting for first snapshot submission...
                                            </p>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </CardContent>
                </Card>
            )}

            {step === 'generate' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <Button variant="ghost" size="sm" className="w-fit -ml-2 mb-2 text-muted-foreground hover:text-foreground" onClick={() => setStep('selection')}>
                            <ArrowLeft className="w-4 h-4 mr-1" /> Back
                        </Button>
                        <CardTitle>Manual IDE Analysis</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <Button
                            onClick={handleGeneratePrompt}
                            disabled={generatePrompt.isPending}
                            className="w-full bg-blue-600 hover:bg-blue-700"
                        >
                            {generatePrompt.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Generate Prompt'}
                        </Button>
                    </CardContent>
                </Card>
            )}

            {step === 'copy' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle>Copy Analysis Prompt</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="relative">
                            <pre className="bg-muted border border-input rounded-lg p-4 text-[10px] text-muted-foreground overflow-auto max-h-80 font-mono whitespace-pre-wrap">
                                {prompt}
                            </pre>
                            <Button size="sm" variant="outline" className="absolute top-2 right-2" onClick={() => handleCopy(prompt)}>
                                {copied ? <Check className="w-3 h-3 mr-1" /> : <Copy className="w-3 h-3 mr-1" />}
                                {copied ? 'Copied' : 'Copy'}
                            </Button>
                        </div>
                        <div className="flex justify-between">
                            <Button variant="outline" size="sm" onClick={() => setStep('generate')}>Back</Button>
                            <Button size="sm" onClick={() => setStep('paste')}>Next <ChevronRight className="w-4 h-4 ml-1" /></Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {step === 'paste' && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle>Paste Results</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <textarea
                            className="w-full h-48 bg-background border border-input rounded-lg p-4 text-xs font-mono"
                            placeholder='Paste JSON here...'
                            value={pastedJson}
                            onChange={(e) => setPastedJson(e.target.value)}
                        />
                        {parseError && <p className="text-xs text-red-500">{parseError}</p>}
                        <div className="flex justify-between">
                            <Button variant="outline" size="sm" onClick={() => setStep('copy')}>Back</Button>
                            <Button
                                size="sm"
                                onClick={handleIngest}
                                className="bg-emerald-600 hover:bg-emerald-700"
                                disabled={ingestBootstrap.isPending || !pastedJson.trim()}
                            >
                                {ingestBootstrap.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : 'Ingest Intelligence'}
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {step === 'complete' && result && (
                <Card className="bg-card border-border shadow-sm">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-3 text-emerald-500">
                            <ShieldCheck className="w-6 h-6" />
                            <span>Harvesting Successful</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <ScoreTile label="Architecture" value={result.scores.architecture_score} />
                            <ScoreTile label="Contract Density" value={result.scores.contract_density} />
                            <ScoreTile label="Risk" value={result.scores.risk_score} />
                            <ScoreTile label="Alignment" value={result.scores.alignment_score} />
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}

function ScoreTile({ label, value }: { label: string; value: number }) {
    return (
        <div className="rounded-lg p-3 bg-muted border border-border text-center">
            <div className="text-2xl font-bold text-violet-400">{value.toFixed(0)}</div>
            <div className="text-[10px] text-muted-foreground uppercase mt-1">{label}</div>
        </div>
    );
}

function generatePreviewConfig(ide: IDEType, projectId: string, token: string, port: number) {
    const mcpUrl = `http://localhost:${port}`;
    if (ide === 'anti-gravity') {
        return {
            "mcpServers": {
                "specforge": {
                    "serverUrl": mcpUrl,
                    "headers": {
                        "Authorization": `Bearer ${token}`,
                        "Content-Type": "application/json"
                    }
                }
            }
        };
    }

    return {
        "mcpServers": {
            "specforge": {
                "command": "specforge-mcp",
                "args": ["serve"],
                "env": {
                    "SF_SERVER_URL": mcpUrl,
                    "SF_PROJECT_ID": projectId,
                    "SF_API_TOKEN": token
                }
            }
        }
    };
}
