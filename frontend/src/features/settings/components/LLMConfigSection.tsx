import { useState, useEffect, useRef } from "react";
import { useLLMSettings, type LLMConfig } from "@/hooks/use-llm-settings";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Sparkles, Key, Globe, CheckCircle2, AlertCircle, RefreshCw, Flame, Terminal } from "lucide-react";
import { getAccessToken, API_BASE_URL } from "@/api/client";

const PROVIDERS = [
    { id: "openai", name: "OpenAI" },
    { id: "anthropic", name: "Anthropic" },
    { id: "gemini", name: "Google Gemini" },
    { id: "ollama", name: "Local / Ollama" },
];

export function LLMConfigSection() {
    const { config, isLoading, updateSettings, isUpdating, testConnection, isTesting, testResult, getWarmupEndpoint, listModels } = useLLMSettings();
    const [formData, setFormData] = useState<LLMConfig | null>(null);
    const [isWarmingUp, setIsWarmingUp] = useState(false);
    const [warmupLogs, setWarmupLogs] = useState<string[]>([]);
    const [availableModels, setAvailableModels] = useState<string[]>([]);
    const [modelsLoading, setModelsLoading] = useState(false);
    const logsEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (config) setFormData(config);
    }, [config]);

    useEffect(() => {
        if (logsEndRef.current) {
            logsEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [warmupLogs]);



    // Fetch models when provider or api key changes (debounced ideally, but on blur/action for now)
    // Actually, let's fetch on mount if we have creds, and maybe add a refresh button
    useEffect(() => {
        if (formData?.provider && formData?.api_key) {
            fetchModels(formData);
        }
    }, [formData?.provider]); // Only on provider change automatically? Or maybe we need a manual refresh button? 
    // The user request said: "Once the user has selected their llm provider and set the apikey then the system must call..."
    // So maybe on blur of API key or selection of provider.

    const fetchModels = async (cfg: LLMConfig) => {
        if (!cfg.api_key) return;
        setModelsLoading(true);
        try {
            const models = await listModels(cfg);
            setAvailableModels(models);
            // If current model is not in list, select first
            if (models.length > 0 && !models.includes(cfg.model)) {
                setFormData(prev => prev ? { ...prev, model: models[0] } : null);
            }
        } catch (e) {
            console.error(e);
            // Fallback or show error? For now just keep empty or existing
        } finally {
            setModelsLoading(false);
        }
    };

    const handleProviderChange = (val: any) => {
        if (!formData) return;
        const newConfig = { ...formData, provider: val, model: "" };
        setFormData(newConfig);
        fetchModels(newConfig);
    };

    const handleSave = () => {
        if (formData) updateSettings(formData);
    };

    const handleTest = () => {
        if (formData) testConnection(formData);
    };

    const handleWarmup = () => {
        if (isWarmingUp) return;
        setIsWarmingUp(true);
        setWarmupLogs(["Starting warmup sequence...", "Connecting to LLM provider..."]);

        // Use the endpoint from the hook/api
        // We need to resolve the full URL since fetch needs it, but getWarmupEndpoint returns relative path
        const endpoint = getWarmupEndpoint();
        // const baseUrl = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"; 
        // Note: API_BASE_URL usually includes /api/v1, checking if endpoint has it
        // llmApi.getWarmupEndpoint returns '/api/v1/settings/llm/warmup'
        // So we need origin + endpoint
        // If API_BASE_URL is 'http://localhost:8080/api/v1', we should extract origin

        // Simpler: assume API_BASE_URL is the prefix.
        // If API_BASE_URL is 'http://localhost:8080/api/v1', and endpoint is '/api/v1...', we might double up.
        // Let's rely on constructing it carefully.
        // Actually, let's just use the relative path if on same domain, or full if CORS.
        // Ideally apiClient structure handles this but for fetch we do it manually.

        const url = `${API_BASE_URL}${endpoint}`;

        fetchWarmup(url);
    };

    const fetchWarmup = async (url: string) => {
        try {
            const token = getAccessToken();
            const headers: HeadersInit = {};
            if (token) {
                headers["Authorization"] = `Bearer ${token}`;
            }

            const response = await fetch(url, { headers });

            if (!response.body) throw new Error("No response body");
            const reader = response.body.getReader();
            const decoder = new TextDecoder();

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;
                const chunk = decoder.decode(value);
                const lines = chunk.split("\n\n");
                for (const line of lines) {
                    if (line.startsWith("data: ")) {
                        const msg = line.replace("data: ", "");
                        if (msg === "{}") continue;
                        setWarmupLogs(prev => [...prev, msg]);
                    }
                    if (line.startsWith("event: done")) {
                        setIsWarmingUp(false);
                        return;
                    }
                    if (line.startsWith("event: error")) {
                        setWarmupLogs(prev => [...prev, `Error: ${line}`]);
                        setIsWarmingUp(false);
                        return;
                    }
                }
            }
        } catch (e: any) {
            setWarmupLogs(prev => [...prev, `Connection failed: ${e.message}`]);
            setIsWarmingUp(false);
        }
    };

    if (isLoading || !formData) return <div>Loading settings...</div>;

    return (
        <div className="space-y-6">
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <div className="space-y-1">
                            <CardTitle>LLM Provider Configuration</CardTitle>
                            <CardDescription>
                                Configure the primary intelligence engine for SpecForge.
                            </CardDescription>
                        </div>
                        <div className="flex items-center gap-2">
                            <Label htmlFor="llm-enabled" className="text-xs font-mono uppercase text-muted-foreground">Enabled</Label>
                            <Switch
                                id="llm-enabled"
                                checked={formData.is_active}
                                onCheckedChange={(val) => setFormData({ ...formData, is_active: val })}
                            />
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label>Provider</Label>
                            <Select
                                value={formData.provider}
                                onValueChange={handleProviderChange}
                            >
                                <SelectTrigger>
                                    <SelectValue placeholder="Select Provider" />
                                </SelectTrigger>
                                <SelectContent>
                                    {PROVIDERS.map((p) => (
                                        <SelectItem key={p.id} value={p.id}>{p.name}</SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label className="flex justify-between">
                                Model
                                {modelsLoading && <span className="text-xs text-muted-foreground animate-pulse">Loading...</span>}
                            </Label>
                            <Select
                                value={formData.model}
                                onValueChange={(val) => setFormData({ ...formData, model: val })}
                                disabled={modelsLoading}
                            >
                                <SelectTrigger>
                                    <SelectValue placeholder={modelsLoading ? "Loading models..." : "Select Model"} />
                                </SelectTrigger>
                                <SelectContent>
                                    {availableModels.length > 0 ? (
                                        availableModels.map((m) => (
                                            <SelectItem key={m} value={m}>{m}</SelectItem>
                                        ))
                                    ) : (
                                        <SelectItem value="custom" disabled>No models found</SelectItem>
                                    )}
                                </SelectContent>
                            </Select>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label className="flex items-center gap-2">
                            <Key className="h-3.5 w-3.5" /> API Key
                        </Label>
                        <Input
                            type="password"
                            placeholder="sk-..."
                            value={formData.api_key}
                            onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                        />
                        <p className="text-[10px] text-muted-foreground px-1">
                            Your API key is stored locally and used only for requests from your browser.
                        </p>
                        <Button variant="ghost" size="sm" onClick={() => formData && fetchModels(formData)} className="h-6 text-xs mt-1">
                            <RefreshCw className="mr-2 h-3 w-3" /> Refresh Models
                        </Button>
                    </div>

                    {formData.provider === "ollama" && (
                        <div className="space-y-2">
                            <Label className="flex items-center gap-2">
                                <Globe className="h-3.5 w-3.5" /> Base URL (Optional)
                            </Label>
                            <Input
                                placeholder="http://localhost:11434/v1"
                                value={formData.base_url || ""}
                                onChange={(e) => setFormData({ ...formData, base_url: e.target.value })}
                            />
                        </div>
                    )}

                    {testResult && (
                        <div className={`p-3 rounded-md flex items-start gap-3 text-sm ${testResult.success ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'}`}>
                            {testResult.success ? <CheckCircle2 className="h-4 w-4 mt-0.5" /> : <AlertCircle className="h-4 w-4 mt-0.5" />}
                            <div>
                                <span className="font-bold">{testResult.success ? "Success: " : "Error: "}</span>
                                {testResult.message}
                            </div>
                        </div>
                    )}

                    {isWarmingUp && (
                        <div className="bg-black text-green-400 font-mono text-xs p-4 rounded-md h-32 overflow-y-auto border border-gray-800 shadow-inner">
                            <div className="flex items-center gap-2 border-b border-gray-800 pb-2 mb-2">
                                <Terminal className="h-3 w-3" />
                                <span>Warmup Console</span>
                            </div>
                            <div className="space-y-1">
                                {warmupLogs.map((log, i) => (
                                    <div key={i}>{log}</div>
                                ))}
                                <div ref={logsEndRef} />
                            </div>
                        </div>
                    )}

                    <div className="flex gap-3 pt-2">
                        <Button
                            variant="outline"
                            className="flex-1"
                            onClick={handleTest}
                            disabled={isTesting || !formData.api_key}
                        >
                            {isTesting && <RefreshCw className="mr-2 h-4 w-4 animate-spin" />}
                            <Sparkles className="mr-2 h-4 w-4" /> Test Connection
                        </Button>
                        <Button
                            variant="outline"
                            className="flex-1"
                            onClick={handleWarmup}
                            disabled={isWarmingUp || !formData.api_key}
                        >
                            {isWarmingUp ? <RefreshCw className="mr-2 h-4 w-4 animate-spin" /> : <Flame className="mr-2 h-4 w-4" />}
                            Warmup Model
                        </Button>
                        <Button
                            className="flex-1"
                            onClick={handleSave}
                            disabled={isUpdating}
                        >
                            Save Configuration
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
