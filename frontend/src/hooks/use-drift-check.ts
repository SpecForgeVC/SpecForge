import { useMutation } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useDriftCheck(contractId: string) {
    return useMutation({
        mutationFn: async (againstVersion: string) => {
            const response = await apiClient.post<components["schemas"]["DriftReport"]>(
                `/contracts/${contractId}/drift-check`,
                { against_version: againstVersion }
            );
            return response.data;
        },
    });
}
