import { useState } from "react";
import { type components } from "@/api/generated/schema";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
    ExternalLink,
    Cpu,
    Check,
    Copy,
    Monitor,
    Code2,
    Loader2,
    Zap
} from "lucide-react";
import { useMCPTokens, useGenerateMCPToken } from "@/hooks/use-mcp";

interface MCPConnectionGuideProps {
    project: components["schemas"]["Project"];
}

type IDEType = 'cursor' | 'anti-gravity' | 'claude';

export function MCPConnectionGuide({ project }: MCPConnectionGuideProps) {
    const projectId = project.id || "";
    const mcpPort = project.mcp_settings?.port || 8081;
    const { data: tokens } = useMCPTokens(projectId);
    const generateToken = useGenerateMCPToken(projectId);
    const [selectedIDE, setSelectedIDE] = useState<IDEType>('cursor');
    const [rawToken, setRawToken] = useState<string | null>(null);

    const handleGenerateToken = () => {
        if (!generateToken.isPending) {
            generateToken.mutate(undefined, {
                onSuccess: (data) => {
                    setRawToken(data.token_raw);
                }
            });
        }
    };

    const activeToken = rawToken || 'YOUR_API_TOKEN_IS_REDACTED_CLICK_REGENERATE';

    const generatePreviewConfig = (ide: string, projectId: string, token: string, port: number) => {
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
                            "SF_API_TOKEN": token
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
    };

    return (
        <Card className="border-indigo-100 shadow-sm overflow-hidden">
            <CardHeader className="bg-indigo-50/30">
                <div className="flex items-center gap-2">
                    <Monitor className="h-5 w-5 text-indigo-600" />
                    <div>
                        <CardTitle className="text-lg">IDE Integration Guide</CardTitle>
                        <CardDescription>
                            Connect your tools to the Reality Anchor Engine.
                        </CardDescription>
                    </div>
                </div>
            </CardHeader>
            <CardContent className="pt-6 space-y-6">
                {/* Step 1: Select IDE */}
                <div className="space-y-3">
                    <div className="flex items-center gap-2">
                        <div className="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center text-[10px] font-bold text-slate-500">1</div>
                        <h4 className="text-xs font-semibold uppercase tracking-wider text-slate-600">Select Your IDE</h4>
                    </div>
                    <div className="grid grid-cols-3 gap-2">
                        {(['cursor', 'anti-gravity', 'claude'] as IDEType[]).map((id) => (
                            <button
                                key={id}
                                onClick={() => setSelectedIDE(id)}
                                className={`p-2 rounded-lg border-2 transition-all flex flex-col items-center gap-1.5 ${selectedIDE === id ? 'border-indigo-600 bg-indigo-50/50' : 'border-slate-100 bg-slate-50/30 hover:border-indigo-300'}`}
                            >
                                <span className="text-[9px] font-bold uppercase tracking-wider">{id.replace('-', ' ')}</span>
                            </button>
                        ))}
                    </div>
                </div>

                {/* Step 2: Config Block */}
                <div className="space-y-3">
                    <div className="flex items-center gap-2">
                        <div className="w-5 h-5 rounded-full bg-slate-100 flex items-center justify-center text-[10px] font-bold text-slate-500">2</div>
                        <h4 className="text-xs font-semibold uppercase tracking-wider text-slate-600">Copy Configuration</h4>
                    </div>
                    <div className="relative group">
                        {(!tokens || tokens.length === 0) && !rawToken ? (
                            <div className="p-8 bg-slate-900 rounded-lg flex flex-col items-center justify-center gap-3 border border-indigo-500/30">
                                <Code2 className="h-8 w-8 text-indigo-400 opacity-50" />
                                <div className="text-center">
                                    <p className="text-xs text-indigo-300 font-medium">No API Token Found</p>
                                    <p className="text-[10px] text-slate-500 mt-1">Generate a secure token to connect your IDE.</p>
                                </div>
                                <Button
                                    size="sm"
                                    className="h-8 bg-indigo-600 hover:bg-indigo-700 text-white text-[10px] gap-2 px-4 mt-1"
                                    onClick={handleGenerateToken}
                                    disabled={generateToken.isPending}
                                >
                                    {generateToken.isPending ? <Loader2 className="h-3 w-3 animate-spin" /> : <Zap className="h-3 w-3" />}
                                    Generate API Token
                                </Button>
                            </div>
                        ) : (
                            <>
                                <pre className="p-4 bg-slate-900 text-indigo-300 rounded-lg text-[10px] font-mono overflow-x-auto leading-relaxed max-h-48 border border-slate-800">
                                    {generateToken.isPending ? 'Generating token...' : JSON.stringify(generatePreviewConfig(selectedIDE, projectId, activeToken, mcpPort), null, 2)}
                                </pre>
                                <div className="absolute right-2 top-2 flex gap-2">
                                    {!rawToken && (
                                        <Button
                                            size="sm"
                                            variant="secondary"
                                            className="h-7 text-[10px] gap-1 px-2 text-indigo-100 bg-indigo-600 hover:bg-indigo-700"
                                            onClick={handleGenerateToken}
                                            disabled={generateToken.isPending}
                                        >
                                            {generateToken.isPending ? <Loader2 className="h-3 w-3 animate-spin" /> : <Zap className="h-3 w-3" />}
                                            Regenerate
                                        </Button>
                                    )}
                                    <Button
                                        size="icon"
                                        variant="ghost"
                                        className="h-7 w-7 text-slate-400 hover:text-white"
                                        onClick={() => navigator.clipboard.writeText(JSON.stringify(generatePreviewConfig(selectedIDE, projectId, activeToken, mcpPort), null, 2))}
                                    >
                                        <Copy className="h-3.5 w-3.5" />
                                    </Button>
                                </div>
                                {rawToken && (
                                    <p className="text-[9px] text-emerald-600 mt-2 flex items-center gap-1 font-medium">
                                        <Check className="h-3 w-3" /> New token generated and added to config.
                                    </p>
                                )}
                            </>
                        )}
                    </div>
                </div>

                <div className="pt-2 border-t">
                    <div className="flex items-center gap-2 text-indigo-700 font-medium mb-1.5">
                        <Cpu className="h-3.5 w-3.5" />
                        <h4 className="text-[11px] uppercase tracking-wider">Need the CLI?</h4>
                    </div>
                    <p className="text-[10px] text-slate-500 leading-relaxed mb-3">
                        If your IDE requires the binary, you can download the <code>specforge-mcp</code> CLI for your platform from the Import Wizard.
                    </p>
                    <Button variant="outline" size="sm" className="h-8 text-[10px] gap-1.5 w-full border-indigo-100 text-indigo-600 hover:bg-indigo-50" asChild>
                        <a href="https://modelcontextprotocol.io" target="_blank" rel="noreferrer">
                            Learn more about MCP <ExternalLink className="h-3 w-3" />
                        </a>
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}
