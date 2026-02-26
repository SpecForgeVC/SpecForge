import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { type UIRoadmapItem, uiRoadmapApi } from '../api/ui_roadmap';

export const useUIRoadmapItems = (projectId: string | undefined) => {
    return useQuery({
        queryKey: ['ui-roadmap', projectId],
        queryFn: () => uiRoadmapApi.list(projectId!),
        enabled: !!projectId,
    });
};

export const useUIRoadmapItem = (id: string | undefined) => {
    return useQuery({
        queryKey: ['ui-roadmap-item', id],
        queryFn: () => uiRoadmapApi.get(id!),
        enabled: !!id,
    });
};

export const useSaveUIRoadmapItem = (projectId: string | undefined) => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (item: Partial<UIRoadmapItem>) => uiRoadmapApi.save(projectId!, item),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['ui-roadmap', projectId] });
        },
    });
};

export const useUpdateUIRoadmapItem = (projectId: string | undefined) => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, item }: { id: string; item: Partial<UIRoadmapItem> }) => uiRoadmapApi.update(id, item),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['ui-roadmap', projectId] });
            queryClient.invalidateQueries({ queryKey: ['ui-roadmap-item', variables.id] });
        },
    });
};

export const useDeleteUIRoadmapItem = (projectId: string | undefined) => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (id: string) => uiRoadmapApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['ui-roadmap', projectId] });
        },
    });
};

export const useExportUIRoadmapItem = () => {
    return useMutation({
        mutationFn: (id: string) => uiRoadmapApi.export(id),
    });
};

export const useFigmaPluginAssets = (id: string | undefined) => {
    return useQuery({
        queryKey: ['ui-roadmap', id, 'plugin-assets'],
        queryFn: () => uiRoadmapApi.getPluginAssets(id!),
        enabled: !!id,
    });
};
