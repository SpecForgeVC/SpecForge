import { useState } from "react";
import { type components } from "@/api/generated/schema";
import { useUpdateProject } from "@/hooks/use-projects";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
    ShieldCheck,
    Settings2,
    Lock,
    Globe,
    Terminal,
    Activity,
    CheckCircle2,
    AlertCircle,
    Copy,
    Check,
    RefreshCcw,
    Layers,
    Box
} from "lucide-react";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

interface MCPSettingsCardProps {
    project: components["schemas"]["Project"];
}

export function MCPSettingsCard({ project }: MCPSettingsCardProps) {
    const updateProject = useUpdateProject(project.id || "");
    const [copied, setCopied] = useState(false);

    // Local state for inputs to avoid jumping on every keystroke
    const [port, setPort] = useState(project.mcp_settings?.port || 8081);
    const [bindAddress, setBindAddress] = useState(project.mcp_settings?.bind_address || "0.0.0.0");
    const [tokenRequired, setTokenRequired] = useState(project.mcp_settings?.token_required ?? true);
    const [importMode, setImportMode] = useState(project.mcp_settings?.import_mode || "light");
    const [repositoryType, setRepositoryType] = useState(project.mcp_settings?.repository_type || "single");

    const isEnabled = !!project.mcp_settings?.enabled;
    const healthStatus = project.mcp_settings?.health_status || "unknown";
    const mcpToken = project.mcp_settings?.token || "default-rae-token-change-me";

    const handleSave = () => {
        updateProject.mutate({
            mcp_settings: {
                enabled: isEnabled,
                port: Number(port),
                bind_address: bindAddress,
                token_required: tokenRequired,
                import_mode: importMode,
                repository_type: repositoryType
            }
        });
    };

    const toggleEnabled = (checked: boolean) => {
        updateProject.mutate({
            mcp_settings: {
                ...project.mcp_settings,
                enabled: checked
            }
        });
    };

    const copyToken = () => {
        navigator.clipboard.writeText(mcpToken);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <Card className="border-indigo-100 shadow-sm overflow-hidden">
            <CardHeader className="bg-gradient-to-r from-indigo-50/50 to-transparent border-b">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <div className="p-2 bg-indigo-100 rounded-lg">
                            <ShieldCheck className="h-5 w-5 text-indigo-600" />
                        </div>
                        <div>
                            <CardTitle className="text-lg">Reality Anchor Engine (MCP)</CardTitle>
                            <CardDescription>
                                Synchronize your IDE with SpecForge's intelligence alignment.
                            </CardDescription>
                        </div>
                    </div>
                    <div className="flex items-center gap-4">
                        {isEnabled && (
                            <div className="flex items-center gap-1.5 px-2.5 py-1 bg-green-50 text-green-700 rounded-full border border-green-100 text-xs font-medium">
                                <Activity className="h-3 w-3 animate-pulse" />
                                {healthStatus.toUpperCase()}
                            </div>
                        )}
                        <Switch
                            checked={isEnabled}
                            onCheckedChange={toggleEnabled}
                            disabled={updateProject.isPending}
                        />
                    </div>
                </div>
            </CardHeader>
            <CardContent className="pt-6 space-y-6">
                {!isEnabled ? (
                    <div className="py-8 text-center space-y-4">
                        <div className="mx-auto w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center">
                            <Terminal className="h-6 w-6 text-slate-400" />
                        </div>
                        <div className="max-w-[400px] mx-auto">
                            <p className="text-sm text-muted-foreground">
                                The Reality Anchor Engine allows SpecForge to see your actual code, migrations, and API routes in real-time via the Model Context Protocol.
                            </p>
                        </div>
                        <Button
                            variant="outline"
                            onClick={() => toggleEnabled(true)}
                            disabled={updateProject.isPending}
                        >
                            Enable Reality Anchor
                        </Button>
                    </div>
                ) : (
                    <>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div className="space-y-4">
                                <div className="space-y-2">
                                    <Label className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                                        Network Configuration
                                    </Label>
                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="space-y-1.5">
                                            <div className="flex items-center gap-1.5 text-sm">
                                                <Globe className="h-3.5 w-3.5 text-slate-400" />
                                                <span>Bind Address</span>
                                            </div>
                                            <Input
                                                value={bindAddress}
                                                onChange={(e) => setBindAddress(e.target.value)}
                                                placeholder="0.0.0.0"
                                                className="h-9"
                                            />
                                        </div>
                                        <div className="space-y-1.5">
                                            <div className="flex items-center gap-1.5 text-sm">
                                                <Settings2 className="h-3.5 w-3.5 text-slate-400" />
                                                <span>Port</span>
                                            </div>
                                            <Input
                                                type="number"
                                                value={port}
                                                onChange={(e) => setPort(Number(e.target.value))}
                                                placeholder="8081"
                                                className="h-9"
                                            />
                                        </div>
                                    </div>
                                </div>

                                <div className="flex items-center justify-between p-3 border rounded-lg bg-slate-50/50">
                                    <div className="space-y-0.5">
                                        <div className="flex items-center gap-1.5 text-sm font-medium">
                                            <Lock className="h-3.5 w-3.5 text-slate-400" />
                                            <span>Token Authentication</span>
                                        </div>
                                        <p className="text-xs text-muted-foreground">Require bearer token for MCP clients.</p>
                                    </div>
                                    <Switch
                                        checked={tokenRequired}
                                        onCheckedChange={setTokenRequired}
                                    />
                                </div>

                                <div className="space-y-4 pt-2 border-t mt-4">
                                    <Label className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                                        Ingestion Controls
                                    </Label>
                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="space-y-1.5">
                                            <div className="flex items-center gap-1.5 text-sm">
                                                <Layers className="h-3.5 w-3.5 text-slate-400" />
                                                <span>Import Mode</span>
                                            </div>
                                            <Select value={importMode} onValueChange={(val) => setImportMode(val as "light" | "full")}>
                                                <SelectTrigger className="h-9">
                                                    <SelectValue placeholder="Select mode" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="light">Light (Metadata Only)</SelectItem>
                                                    <SelectItem value="full">Full (Deep Structural)</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        </div>
                                        <div className="space-y-1.5">
                                            <div className="flex items-center gap-1.5 text-sm">
                                                <Box className="h-3.5 w-3.5 text-slate-400" />
                                                <span>Repo Type</span>
                                            </div>
                                            <Select value={repositoryType} onValueChange={(val) => setRepositoryType(val as "monorepo" | "polyrepo" | "single")}>
                                                <SelectTrigger className="h-9">
                                                    <SelectValue placeholder="Select type" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="single">Single Component</SelectItem>
                                                    <SelectItem value="monorepo">Monorepo</SelectItem>
                                                    <SelectItem value="polyrepo">Polyrepo (Distributed)</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="space-y-4">
                                <div className="space-y-2">
                                    <Label className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                                        Authentication Token
                                    </Label>
                                    <div className="relative group">
                                        <pre className="p-3 bg-slate-900 text-slate-50 rounded-lg text-xs font-mono overflow-x-auto pr-10">
                                            {mcpToken}
                                        </pre>
                                        <Button
                                            size="icon"
                                            variant="ghost"
                                            className="absolute right-1.5 top-1.5 h-7 w-7 text-slate-400 hover:text-white"
                                            onClick={copyToken}
                                        >
                                            {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
                                        </Button>
                                    </div>
                                    <p className="text-[10px] text-muted-foreground italic">
                                        Use this token in your IDE settings to authorize the RAE client.
                                    </p>
                                </div>

                                <Alert className="bg-blue-50/50 border-blue-100">
                                    <AlertCircle className="h-4 w-4 text-blue-600" />
                                    <AlertTitle className="text-xs font-semibold text-blue-800">Connection string</AlertTitle>
                                    <AlertDescription className="text-[11px] text-blue-700 font-mono mt-1 break-all">
                                        http://{bindAddress === "0.0.0.0" ? "localhost" : bindAddress}:{port}/mcp
                                    </AlertDescription>
                                </Alert>
                            </div>
                        </div>

                        <div className="flex items-center justify-end gap-3 pt-2 border-t mt-4">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="text-xs"
                                onClick={() => {
                                    setPort(project.mcp_settings?.port || 8081);
                                    setBindAddress(project.mcp_settings?.bind_address || "0.0.0.0");
                                    setTokenRequired(project.mcp_settings?.token_required ?? true);
                                    setImportMode(project.mcp_settings?.import_mode || "light");
                                    setRepositoryType(project.mcp_settings?.repository_type || "single");
                                }}
                            >
                                Discard Changes
                            </Button>
                            <Button
                                size="sm"
                                className="text-xs bg-indigo-600 hover:bg-indigo-700"
                                onClick={handleSave}
                                disabled={updateProject.isPending}
                            >
                                {updateProject.isPending ? (
                                    <RefreshCcw className="h-3 w-3 animate-spin mr-2" />
                                ) : (
                                    <CheckCircle2 className="h-3 w-3 mr-2" />
                                )}
                                Save Configuration
                            </Button>
                        </div>
                    </>
                )}
            </CardContent>
        </Card>
    );
}
