import { fetchWithRefresh, handleResponse } from "./apiClient";

export const adminService = {
    // Disputes
    getDisputes: async (): Promise<any[]> => {
        const response = await fetchWithRefresh("/api/v1/admin/disputes", {
            credentials: "include",
        });
        const data = await handleResponse<{ disputes: any[] }>(response);
        return data.disputes || [];
    },

    assignDispute: async (disputeId: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/admin/disputes/${disputeId}/assign`, {
            method: "POST",
            credentials: "include",
        });
        return handleResponse(response);
    },

    resolveDispute: async (disputeId: string, resolution: string, winnerId: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/admin/disputes/${disputeId}/resolve`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ resolution, winner_id: winnerId }),
            credentials: "include",
        });
        return handleResponse(response);
    },

    // Dashboard Stats
    getStats: async (): Promise<any> => {
        const response = await fetchWithRefresh("/api/v1/admin/dashboard/stats", {
            credentials: "include",
        });
        return handleResponse(response);
    },
};
