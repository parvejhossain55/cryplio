import { fetchWithRefresh, handleResponse } from "./apiClient";

export const notificationService = {
    getPreferences: async (): Promise<any> => {
        const response = await fetchWithRefresh("/api/v1/notifications/preferences", {
            credentials: "include",
        });
        return handleResponse(response);
    },

    savePreferences: async (prefs: { email?: any; push?: any; sms?: any }): Promise<any> => {
        const response = await fetchWithRefresh("/api/v1/notifications/preferences", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(prefs),
            credentials: "include",
        });
        return handleResponse(response);
    },

    getNotifications: async (): Promise<any[]> => {
        const response = await fetchWithRefresh("/api/v1/notifications", {
            credentials: "include",
        });
        const data = await handleResponse<{ notifications: any[] }>(response);
        return data.notifications || [];
    },

    markRead: async (id: string): Promise<void> => {
        const response = await fetchWithRefresh(`/api/v1/notifications/${id}/read`, {
            method: "PATCH",
            credentials: "include",
        });
        await handleResponse(response);
    },
};
