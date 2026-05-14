import { ApiClient } from "./apiClient";
export interface AdResponse {
    ad_id: string;
    user_id: string;
    username: string;
    user_avatar?: string;
    is_online?: boolean;
    type: "buy" | "sell";
    crypto_symbol: string;
    fiat_symbol: string;
    price_type: "fixed" | "floating";
    price: number;
    min_amount: number;
    max_amount: number;
    payment_methods: string[];
    payment_method_ids: number[];
    payment_window_minutes: number;
    trade_terms?: string;
    status: "active" | "paused" | "closed";
    created_at: string;
    user_trades?: number;
    user_rating?: number;
}

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
