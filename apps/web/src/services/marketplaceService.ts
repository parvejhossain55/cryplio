import { ApiClient } from "./apiClient";
import { AdResponse } from "./authService";

export class MarketplaceService {
    static async getAds(params: Record<string, string>): Promise<{ ads: AdResponse[]; total: number }> {
        const query = new URLSearchParams(params).toString();
        return ApiClient.get(`/api/v1/marketplace/ads?${query}`);
    }

    static async createAd(adData: any): Promise<any> {
        return ApiClient.post(`/api/v1/marketplace/ads`, adData);
    }

    static async updateAd(adId: string, adData: any): Promise<any> {
        return ApiClient.put(`/api/v1/marketplace/ads/${adId}`, adData);
    }

    static async deleteAd(adId: string): Promise<any> {
        return ApiClient.delete(`/api/v1/marketplace/ads/${adId}`);
    }

    static async getMyAds(): Promise<AdResponse[]> {
        const data = await ApiClient.get<{ ads: AdResponse[] }>(`/api/v1/marketplace/my-ads`);
        return data.ads || [];
    }

    static async toggleAdStatus(adId: string): Promise<any> {
        return ApiClient.patch(`/api/v1/marketplace/ads/${adId}/status`);
    }
}
