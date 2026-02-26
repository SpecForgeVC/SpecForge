import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { llmApi, type LLMConfiguration } from "@/api/llm";

export type LLMConfig = LLMConfiguration;

export function useLLMSettings() {
    const queryClient = useQueryClient();
    const [testResult, setTestResult] = useState<{ success: boolean; message: string } | null>(null);

    const { data: config, isLoading } = useQuery({
        queryKey: ["llm-settings"],
        queryFn: async (): Promise<LLMConfig> => {
            const data = await llmApi.getConfig();
            return data || {
                provider: "openai",
                api_key: "",
                model: "gpt-4o",
                is_active: false,
                base_url: "",
            };
        },
    });

    const updateMutation = useMutation({
        mutationFn: async (newConfig: LLMConfig) => {
            return await llmApi.updateConfig(newConfig);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["llm-settings"] });
        },
    });

    const testConnectionMutation = useMutation({
        mutationFn: async (config: LLMConfig) => {
            setTestResult(null);
            try {
                const res = await llmApi.testConnection(config);
                return { success: true, message: res.message };
            } catch (error: any) {
                console.error("Test connection failed:", error);
                throw new Error(error.response?.data?.error || error.message || "Failed to connect to provider.");
            }
        },
        onSuccess: (data) => {
            setTestResult(data);
        },
        onError: (error: any) => {
            setTestResult({ success: false, message: error.message });
        },
    });

    const listModels = async (config: LLMConfig) => {
        try {
            return await llmApi.listModels(config);
        } catch (e: any) {
            console.error("Failed to list models", e);
            throw e;
        }
    };

    return {
        config,
        isLoading,
        updateSettings: updateMutation.mutate,
        isUpdating: updateMutation.isPending,
        testConnection: testConnectionMutation.mutate,
        isTesting: testConnectionMutation.isPending,
        testResult,
        getWarmupEndpoint: llmApi.getWarmupEndpoint,
        listModels
    };
}
