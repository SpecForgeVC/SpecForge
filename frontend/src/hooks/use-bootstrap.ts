import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
    bootstrapApi,
    type BootstrapPayload,
    type BootstrapIngestResponse,
    type BootstrapPromptResponse,
    type BootstrapDiffResult,
} from "@/api/bootstrap";

// --- Query Keys ---

export const bootstrapKeys = {
    all: ["bootstrap"] as const,
    snapshots: (projectId: string) => [...bootstrapKeys.all, "snapshots", projectId] as const,
    snapshot: (projectId: string, snapshotId: string) => [...bootstrapKeys.all, "snapshot", projectId, snapshotId] as const,
    latest: (projectId: string) => [...bootstrapKeys.all, "latest", projectId] as const,
};

// --- Queries ---

export function useBootstrapSnapshots(projectId?: string) {
    return useQuery({
        queryKey: bootstrapKeys.snapshots(projectId || ""),
        queryFn: () => bootstrapApi.listSnapshots(projectId!),
        enabled: !!projectId,
        retry: (failureCount, error: any) => {
            if (error?.response?.status === 404) return false;
            return failureCount < 1;
        }
    });
}

export function useBootstrapSnapshot(projectId?: string, snapshotId?: string) {
    return useQuery({
        queryKey: bootstrapKeys.snapshot(projectId || "", snapshotId || ""),
        queryFn: () => bootstrapApi.getSnapshot(projectId!, snapshotId!),
        enabled: !!projectId && !!snapshotId,
    });
}

export function useLatestBootstrapSnapshot(projectId?: string) {
    return useQuery({
        queryKey: bootstrapKeys.latest(projectId || ""),
        queryFn: () => bootstrapApi.getLatestSnapshot(projectId!),
        enabled: !!projectId,
        retry: (failureCount, error: any) => {
            // Don't retry if the snapshot doesn't exist (404)
            if (error?.response?.status === 404) {
                return false;
            }
            return failureCount < 1;
        }
    });
}

// --- Mutations ---

export function useGenerateBootstrapPrompt() {
    return useMutation<BootstrapPromptResponse, Error, string>({
        mutationFn: (projectId: string) => bootstrapApi.generatePrompt(projectId),
    });
}

export function useIngestBootstrap(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation<BootstrapIngestResponse, Error, BootstrapPayload>({
        mutationFn: (payload: BootstrapPayload) => bootstrapApi.ingestBootstrap(projectId, payload),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: bootstrapKeys.snapshots(projectId) });
            queryClient.invalidateQueries({ queryKey: bootstrapKeys.latest(projectId) });
        },
    });
}

export function useDiffBootstrapSnapshots(projectId: string) {
    return useMutation<BootstrapDiffResult, Error, { fromId?: string; toId?: string }>({
        mutationFn: ({ fromId, toId }) => bootstrapApi.diffSnapshots(projectId, fromId, toId),
    });
}
