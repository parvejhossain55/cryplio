import { AdResponse } from "@/types/api";
import { fetchWithRefresh, handleResponse } from "./apiClient";

export const tradeService = {
    // Advertisements
    getAds: async (params: { 
        type?: string; 
        cryptoId?: number; 
        fiatId?: number;
        limit?: number;
        offset?: number;
        fiat_currency?: string;
        payment_method?: string;
    } = {}): Promise<{ ads: AdResponse[]; total: number }> => {
        const queryParams = new URLSearchParams();
        if (params.type) queryParams.append("type", params.type);
        if (params.cryptoId) queryParams.append("crypto_id", params.cryptoId.toString());
        if (params.fiatId) queryParams.append("fiat_id", params.fiatId.toString());
        if (params.limit) queryParams.append("limit", params.limit.toString());
        if (params.offset) queryParams.append("offset", params.offset.toString());
        if (params.fiat_currency) queryParams.append("fiat_currency", params.fiat_currency);
        if (params.payment_method) queryParams.append("payment_method", params.payment_method);

        const response = await fetch(`/api/v1/marketplace/ads?${queryParams.toString()}`, {
            credentials: "include",
        });
        const data = await handleResponse<{ ads: AdResponse[]; total: number }>(response);
        return {
            ads: data.ads || [],
            total: data.total || 0,
        };
    },

    getMyAds: async (): Promise<AdResponse[]> => {
        const response = await fetchWithRefresh("/api/v1/marketplace/my-ads", {
            credentials: "include",
        });
        const data = await handleResponse<{ ads: AdResponse[] }>(response);
        return data.ads || [];
    },

    createAd: async (adData: any): Promise<any> => {
        const response = await fetchWithRefresh("/api/v1/marketplace/ads", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(adData),
            credentials: "include",
        });
        return handleResponse(response);
    },

    updateAd: async (adId: string, adData: any): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/ads/${adId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(adData),
            credentials: "include",
        });
        return handleResponse(response);
    },

    deleteAd: async (adId: string): Promise<void> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/ads/${adId}`, {
            method: "DELETE",
            credentials: "include",
        });
        await handleResponse(response);
    },

    toggleAdStatus: async (adId: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/ads/${adId}/status`, {
            method: "PATCH",
            credentials: "include",
        });
        return handleResponse(response);
    },

    // Trades
    initiateTrade: async (adId: string, amount: number): Promise<{ trade_id: string; message: string }> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/ads/${adId}/trades`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ amount }),
            credentials: "include",
        });
        return handleResponse(response);
    },

    getMyTrades: async (role?: string): Promise<any[]> => {
        const url = role ? `/api/v1/marketplace/trades?role=${role}` : `/api/v1/marketplace/trades`;
        const response = await fetchWithRefresh(url, {
            credentials: "include",
        });
        const data = await handleResponse<{ trades: any[] }>(response);
        return data?.trades || [];
    },

    getTradeDetails: async (tradeId: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}`, {
            credentials: "include",
        });
        return handleResponse(response);
    },

    updateTradeStatus: async (tradeId: string, action: string, reason?: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}/status`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ action, reason }),
            credentials: "include",
        });
        return handleResponse(response);
    },

    disputeTrade: async (tradeId: string, reasonCode: string, reasonText: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}/dispute`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ reason_code: reasonCode, reason_text: reasonText }),
            credentials: "include",
        });
        return handleResponse(response);
    },

    // Chat
    getChatHistory: async (tradeId: string): Promise<any[]> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}/messages`, {
            credentials: "include",
        });
        return handleResponse(response);
    },

    sendMessage: async (tradeId: string, content: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}/messages`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content }),
            credentials: "include",
        });
        return handleResponse(response);
    },

    leaveFeedback: async (tradeId: string, rating: number, comment: string): Promise<any> => {
        const response = await fetchWithRefresh(`/api/v1/marketplace/trades/${tradeId}/feedback`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ rating, comment }),
            credentials: "include",
        });
        return handleResponse(response);
    },
};
