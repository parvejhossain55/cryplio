import { ApiClient } from "./apiClient";

export class NotificationService {
    static async getNotifications(): Promise<any[]> {
        try {
            const data = await ApiClient.get<{ notifications: any[] }>(`/api/v1/notifications`);
            return data?.notifications || [];
        } catch (error: any) {
            // Return empty array for unauthorized/not found during initial load
            if (error.message?.includes("user not found") || error.message?.includes("Unauthorized") || error.message?.includes("Session expired")) {
                return [];
            }
            throw error;
        }
    }

    static async markAsRead(notificationId: string): Promise<any> {
        return ApiClient.patch(`/api/v1/notifications/${notificationId}/read`);
    }

    static async getPreferences(): Promise<any> {
        return ApiClient.get(`/api/v1/notifications/preferences`);
    }

    static async savePreferences(prefs: any): Promise<any> {
        return ApiClient.post(`/api/v1/notifications/preferences`, prefs);
    }
}
