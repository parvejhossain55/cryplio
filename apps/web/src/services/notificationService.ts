import { ApiClient } from "./apiClient";

export class NotificationService {
    static async getNotifications(): Promise<any[]> {
        const data = await ApiClient.get<{ notifications: any[] }>(`/api/v1/notifications`);
        return data.notifications || [];
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
