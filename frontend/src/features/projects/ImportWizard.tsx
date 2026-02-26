import { useState } from 'react';
import {
    Copy,
    Check,
    Loader2,
    Monitor,
    ShieldCheck,
    Terminal,
    ChevronRight,
    ArrowLeft,
    RefreshCw,
    Download,
    Cpu,
    ExternalLink,
    AlertCircle,
    Power
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/api/client';

type WizardStep = 'status' | 'token' | 'ide' | 'config' | 'download';

interface MCPStatus {
    enabled: boolean;
    port: number;
    bind_address: string;
    is_running: boolean;
    auth_required: boolean;
}

interface MCPToken {
    id: string;
    token_prefix: string;
    revoked: boolean;
    last_used_at?: string;
}

interface MCPTokenRaw {
    token_raw: string;
}

export default function ImportWizard({ projectId }: { projectId: string }) {
    const [step, setStep] = useState<WizardStep>('status');
    const [selectedIDE, setSelectedIDE] = useState<'cursor' | 'anti-gravity' | 'claude'>('cursor');
    const [rawToken, setRawToken] = useState<string | null>(null);
    const queryClient = useQueryClient();

    // Queries
    const { data: status, isLoading: loadingStatus, refetch: refetchStatus } = useQuery<MCPStatus>({
        queryKey: ['mcp-status'],
        queryFn: () => apiClient.get('mcp/status').then(res => res.data.data),
        refetchInterval: 5000
    });

    const { data: tokens } = useQuery<MCPToken[]>({
        queryKey: ['mcp-tokens', projectId],
        queryFn: () => apiClient.get(`mcp/tokens?project_id=${projectId}`).then(res => res.data.data)
    });

    // Mutations
    const generateToken = useMutation({
        mutationFn: () => apiClient.post(`mcp/token?project_id=${projectId}`).then(res => res.data.data),
        onSuccess: (data: MCPTokenRaw) => {
            setRawToken(data.token_raw);
            queryClient.invalidateQueries({ queryKey: ['mcp-tokens', projectId] });
        }
    });

    const activeToken = tokens?.find((t: MCPToken) => !t.revoked);

    const handleDownloadConfig = () => {
        const tokenParam = rawToken ? `&token=${rawToken}` : '';
        // Note: Download link still needs full URL or /api/v1/mcp/... for browser download
        window.location.href = `/api/v1/mcp/config/download?ide=${selectedIDE}&project_id=${projectId}${tokenParam}`;
    };

    return (
        <div className="max-w-4xl mx-auto py-8 px-4">
            <div className="mb-8 space-y-2">
                <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-violet-400 to-indigo-500 bg-clip-text text-transparent">Project Import Wizard</h1>
                <p className="text-muted-foreground">Follow these steps to connect your IDE to SpecForge via MCP.</p>
            </div>

            {/* Stepper */}
            <div className="flex gap-4 mb-8 overflow-x-auto pb-2">
                {[
                    { id: 'status', label: 'Server Status', icon: Monitor },
                    { id: 'token', label: 'API Token', icon: ShieldCheck },
                    { id: 'ide', label: 'Choose IDE', icon: Cpu },
                    { id: 'config', label: 'Configure', icon: Terminal },
                    { id: 'download', label: 'Download', icon: Download }
                ].map((s, idx) => (
                    <div key={s.id} className="flex items-center gap-2 flex-shrink-0">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold transition-all ${step === s.id ? 'bg-violet-600 text-white shadow-[0_0_15px_rgba(139,92,246,0.5)]' :
                            ['status', 'token', 'ide', 'config', 'download'].indexOf(step) > idx ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/50' :
                                'bg-muted text-muted-foreground'
                            }`}>
                            {['status', 'token', 'ide', 'config', 'download'].indexOf(step) > idx ? <Check className="w-4 h-4" /> : idx + 1}
                        </div>
                        <span className={`text-sm font-medium ${step === s.id ? 'text-foreground' : 'text-muted-foreground'}`}>{s.label}</span>
                        {idx < 4 && <div className="w-12 h-[1px] bg-border mx-2" />}
                    </div>
                ))}
            </div>

            {/* Step Content */}
            <div className="grid gap-6">
                {step === 'status' && (
                    <Card className="bg-card/50 backdrop-blur border-border/50 shadow-xl border-t-2 border-t-violet-500/20">
                        <CardHeader>
                            <CardTitle className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Monitor className="w-5 h-5 text-violet-400" />
                                    MCP Server Status
                                </div>
                                {status && (
                                    <div className={`px-2 py-1 rounded text-[10px] font-bold uppercase tracking-widest ${status.is_running ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/30' : 'bg-red-500/10 text-red-400 border border-red-500/30'}`}>
                                        {status.is_running ? 'Online' : 'Offline'}
                                    </div>
                                )}
                            </CardTitle>
                            <CardDescription>Verify the Reality Anchor Engine (RAE) is running correctly.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            {loadingStatus ? (
                                <div className="flex items-center gap-2 text-sm text-muted-foreground"><Loader2 className="w-4 h-4 animate-spin" /> Verifying status...</div>
                            ) : status ? (
                                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                    <StatusBit label="Role" value="RAE Primary" icon={Cpu} />
                                    <StatusBit label="Port" value={status.port.toString()} icon={Terminal} />
                                    <StatusBit label="Host" value={status.bind_address} icon={Monitor} />
                                    <StatusBit label="State" value={status.is_running ? 'ACTIVE' : 'READY'} icon={Power} success={status.is_running} />
                                </div>
                            ) : (
                                <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/30 text-red-500 text-sm flex items-center gap-2">
                                    <AlertCircle className="w-4 h-4" /> Failed to connect to backend management API.
                                </div>
                            )}

                            <div className="pt-4 flex justify-end gap-2">
                                <Button variant="outline" size="sm" onClick={() => refetchStatus()}><RefreshCw className="w-3 h-3 mr-2" /> Refresh</Button>
                                <Button className="bg-violet-600 hover:bg-violet-700" onClick={() => setStep('token')} disabled={!status?.is_running}>
                                    Next Step <ChevronRight className="w-4 h-4 ml-1" />
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                )}

                {step === 'token' && (
                    <Card className="bg-card/50 backdrop-blur border-border/50 shadow-xl border-t-2 border-t-indigo-500/20">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <ShieldCheck className="w-5 h-5 text-indigo-400" />
                                Generate API Token
                            </CardTitle>
                            <CardDescription>This token authorizes your IDE to connect to the MCP server securely.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            {rawToken ? (
                                <div className="space-y-4">
                                    <div className="p-4 rounded-lg bg-emerald-500/5 border border-emerald-500/20 space-y-3">
                                        <div className="flex items-center justify-between">
                                            <span className="text-xs font-semibold text-emerald-400 uppercase tracking-wider">New Token Generated</span>
                                            <div className="flex items-center gap-1 text-[10px] text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded">
                                                <AlertCircle className="w-3 h-3" /> Save this now. It will not be shown again.
                                            </div>
                                        </div>
                                        <div className="relative">
                                            <pre className="bg-background/50 p-4 rounded border border-emerald-500/20 text-sm font-mono text-emerald-300 break-all whitespace-pre-wrap">
                                                {rawToken}
                                            </pre>
                                            <Button size="icon" variant="ghost" className="absolute right-2 top-2 h-8 w-8 text-emerald-400" onClick={() => navigator.clipboard.writeText(rawToken)}>
                                                <Copy className="w-4 h-4" />
                                            </Button>
                                        </div>
                                    </div>
                                    <Button className="w-full bg-emerald-600 hover:bg-emerald-700 text-white" onClick={() => setStep('ide')}>
                                        I have saved my token <ChevronRight className="w-4 h-4 ml-1" />
                                    </Button>
                                </div>
                            ) : activeToken ? (
                                <div className="space-y-4">
                                    <div className="p-4 rounded-lg bg-muted/50 border border-border flex items-center justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 rounded-full bg-indigo-500/10 flex items-center justify-center">
                                                <ShieldCheck className="w-5 h-5 text-indigo-400" />
                                            </div>
                                            <div>
                                                <h4 className="text-sm font-medium">Active Token: {activeToken.token_prefix}...</h4>
                                                <p className="text-[10px] text-muted-foreground uppercase tracking-wider">SECURE & ROTATABLE</p>
                                            </div>
                                        </div>
                                        <Button variant="outline" size="sm" className="text-xs" onClick={() => generateToken.mutate()}>Regenerate</Button>
                                    </div>
                                    <Button className="w-full bg-indigo-600 hover:bg-indigo-700" onClick={() => setStep('ide')}>
                                        Continue <ChevronRight className="w-4 h-4 ml-1" />
                                    </Button>
                                </div>
                            ) : (
                                <Button className="w-full h-24 border-dashed border-2 hover:border-indigo-500/50 hover:bg-indigo-500/5 transition-all text-indigo-400 font-medium" onClick={() => generateToken.mutate()} disabled={generateToken.isPending}>
                                    {generateToken.isPending ? <><Loader2 className="w-5 h-5 mr-3 animate-spin" /> Generating...</> : 'Generate Secure API Token'}
                                </Button>
                            )}

                            <div className="pt-2">
                                <Button variant="ghost" size="sm" onClick={() => setStep('status')} className="text-muted-foreground"><ArrowLeft className="w-3 h-3 mr-2" /> Back</Button>
                            </div>
                        </CardContent>
                    </Card>
                )}

                {step === 'ide' && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <IDECard id="cursor" label="Cursor" description="The AI Code Editor" icon="/cursor-logo.png" selected={selectedIDE === 'cursor'} onClick={() => { setSelectedIDE('cursor'); setStep('config'); }} />
                        <IDECard id="anti-gravity" label="Anti-Gravity" description="Advanced Agentic AI" icon="/logo.png" selected={selectedIDE === 'anti-gravity'} onClick={() => { setSelectedIDE('anti-gravity'); setStep('config'); }} />
                        <IDECard id="claude" label="Claude" description="Anthropic Desktop" icon="/claude-logo.png" selected={selectedIDE === 'claude'} onClick={() => { setSelectedIDE('claude'); setStep('config'); }} />
                    </div>
                )}

                {step === 'config' && (
                    <Card className="bg-card/50 backdrop-blur border-border/50 shadow-xl border-t-2 border-t-emerald-500/20">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Terminal className="w-5 h-5 text-emerald-400" />
                                Configuration Preview
                            </CardTitle>
                            <CardDescription>Verify the generated configuration for {selectedIDE}.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-6">
                            <div className="relative">
                                <pre className="bg-muted p-4 rounded-lg border border-border text-xs font-mono text-muted-foreground overflow-auto max-h-64">
                                    {JSON.stringify(generatePreviewConfig(selectedIDE, projectId, activeToken?.token_prefix || '...', status?.port || 8081), null, 2)}
                                </pre>
                                <div className="absolute top-2 right-2 flex gap-2">
                                    <Button size="icon" variant="ghost" className="h-8 w-8" onClick={() => navigator.clipboard.writeText(JSON.stringify(generatePreviewConfig(selectedIDE, projectId, '...', 8081), null, 2))}>
                                        <Copy className="w-4 h-4" />
                                    </Button>
                                </div>
                            </div>

                            <div className="pt-4 flex justify-between">
                                <Button variant="ghost" onClick={() => setStep('ide')}><ArrowLeft className="w-4 h-4 mr-2" /> Change IDE</Button>
                                <Button className="bg-emerald-600 hover:bg-emerald-700 text-white" onClick={() => setStep('download')}>
                                    Looks Good <ChevronRight className="w-4 h-4 ml-1" />
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                )}

                {step === 'download' && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                        <Card className="bg-card/50 backdrop-blur border-border/50 shadow-xl overflow-hidden">
                            <div className="h-1 bg-gradient-to-r from-emerald-400 to-teal-500" />
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Download className="w-5 h-5 text-emerald-400" />
                                    Final Assets
                                </CardTitle>
                                <CardDescription>Download your config and install the connector.</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-6">
                                <Button className="w-full bg-emerald-600 hover:bg-emerald-700 text-white shadow-lg shadow-emerald-500/20 h-12" onClick={handleDownloadConfig}>
                                    <Download className="w-4 h-4 mr-2" /> Download {selectedIDE.toUpperCase()} Config
                                </Button>

                                <div className="space-y-4">
                                    <h4 className="text-sm font-semibold text-foreground border-b pb-2">IDE Connection Steps</h4>
                                    <div className="space-y-4">
                                        <div className="space-y-2">
                                            <p className="text-xs text-muted-foreground font-medium uppercase tracking-wider">1. Open Settings</p>
                                            <p className="text-[11px] text-muted-foreground">Go to your IDE's MCP or Tool settings panel.</p>
                                        </div>
                                        <div className="space-y-2">
                                            <p className="text-xs text-muted-foreground font-medium uppercase tracking-wider">2. Paste Configuration</p>
                                            <p className="text-[11px] text-muted-foreground">Copy the JSON block from the previous step and paste it into the manual configuration field.</p>
                                        </div>
                                        <div className="space-y-2">
                                            <p className="text-xs text-muted-foreground font-medium uppercase tracking-wider">3. Connect binary</p>
                                            <p className="text-[11px] text-muted-foreground">Ensure the <code className="text-emerald-400">specforge-mcp</code> binary is in your PATH or update the configuration with its absolute path.</p>
                                        </div>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <div className="space-y-6">
                            <Card className="bg-card/50 backdrop-blur border-border/50 shadow-lg">
                                <CardHeader>
                                    <CardTitle className="text-sm font-bold flex items-center gap-2">
                                        <AlertCircle className="w-4 h-4 text-amber-400" />
                                        Important Instructions
                                    </CardTitle>
                                </CardHeader>
                                <CardContent className="text-xs space-y-3 text-muted-foreground">
                                    <p>• The configuration file contains your <span className="text-emerald-400 font-bold">API Token</span>. Do not commit it to public repositories.</p>
                                    <p>• Ensure your IDE's MCP client is pointed to <span className="text-foreground">localhost:8081</span>.</p>
                                    <p>• If the connection fails, regenerate your token and update the config file.</p>
                                    <div className="pt-2">
                                        <a href="#" className="text-violet-400 hover:text-violet-300 flex items-center gap-1 font-medium transition-colors">
                                            View Documentation <ExternalLink className="w-3 h-3" />
                                        </a>
                                    </div>
                                </CardContent>
                            </Card>

                            <div className="p-6 rounded-xl bg-violet-500/10 border border-violet-500/20 text-center space-y-4 shadow-[0_0_20px_rgba(139,92,246,0.1)]">
                                <div className="w-12 h-12 rounded-full bg-violet-600 mx-auto flex items-center justify-center text-white shadow-[0_0_20px_rgba(139,92,246,0.5)]">
                                    <Sparkles className="w-6 h-6" />
                                </div>
                                <h3 className="font-bold text-foreground">Ready to Build?</h3>
                                <p className="text-xs text-muted-foreground">Once connected, your IDE will automatically synchronize every change with SpecForge.</p>
                                <Button variant="outline" className="w-full border-violet-500/30 text-violet-400 hover:bg-violet-500/10" onClick={() => window.location.href = `/projects/${projectId}`}>Return to Dashboard</Button>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}

function StatusBit({ label, value, icon: Icon, success }: { label: string; value: string; icon: any; success?: boolean }) {
    return (
        <div className="p-3 rounded-lg bg-muted/30 border border-border flex flex-col gap-1 transition-all hover:bg-muted/50">
            <div className="flex items-center gap-1 text-[10px] text-muted-foreground font-medium uppercase tracking-tight">
                <Icon className="w-2.5 h-2.5" />
                {label}
            </div>
            <div className={`text-sm font-bold truncate ${success ? 'text-emerald-400' : 'text-foreground'}`}>
                {value}
            </div>
        </div>
    );
}

function IDECard({ label, description, icon, selected, onClick }: { label: string; description: string; icon: string; selected: boolean; onClick: () => void; id: string }) {
    return (
        <div
            onClick={onClick}
            className={`cursor-pointer group p-6 rounded-2xl border-2 transition-all flex flex-col items-center text-center gap-4 ${selected ? 'border-violet-600 bg-violet-600/5 shadow-lg shadow-violet-600/10 scale-[1.02]' : 'border-border bg-card/50 hover:border-violet-600/50 hover:bg-violet-600/5'
                }`}
        >
            <div className={`w-16 h-16 rounded-2xl flex items-center justify-center transition-all ${selected ? 'bg-violet-600 text-white' : 'bg-muted group-hover:bg-violet-600/10'}`}>
                <img src={icon} alt={label} className="w-10 h-10 object-contain brightness-0 invert opacity-80 group-hover:opacity-100" />
            </div>
            <div>
                <h3 className="font-bold text-foreground transition-colors group-hover:text-violet-400">{label}</h3>
                <p className="text-xs text-muted-foreground mt-1">{description}</p>
            </div>
            <div className={`mt-auto px-4 py-1.5 rounded-full text-[10px] font-bold uppercase transition-all ${selected ? 'bg-violet-600 text-white' : 'bg-muted text-muted-foreground group-hover:bg-violet-600/20 group-hover:text-violet-400'}`}>
                {selected ? 'Selected' : 'Select'}
            </div>
        </div>
    );
}

function Sparkles(props: any) {
    return (
        <svg {...props} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-sparkles"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z" /><path d="M5 3v4" /><path d="M19 17v4" /><path d="M3 5h4" /><path d="M17 19h4" /></svg>
    )
}

function generatePreviewConfig(ide: string, projectId: string, token: string, port: number) {
    const mcpUrl = `http://localhost:${port}`;
    if (ide === 'cursor') {
        return {
            "mcpServers": {
                "specforge": {
                    "command": "specforge-mcp",
                    "args": ["serve"],
                    "env": {
                        "SF_SERVER_URL": mcpUrl,
                        "SF_PROJECT_ID": projectId,
                        "SF_API_TOKEN": "sf_live_" + token + "..."
                    }
                }
            }
        };
    }
    if (ide === 'anti-gravity') {
        return {
            "mcpServers": {
                "specforge": {
                    "serverUrl": mcpUrl,
                    "headers": {
                        "Authorization": "Bearer sf_live_" + token + "...",
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
                    "SF_API_TOKEN": "sf_live_" + token + "..."
                }
            }
        }
    };
}
