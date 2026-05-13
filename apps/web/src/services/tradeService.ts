import { ApiClient } from "./apiClient";

export class TradeService {
    static async initiateTrade(adId: string, amount: number): Promise<{ trade_id: string; message: string }> {
        return ApiClient.post(`/api/v1/marketplace/ads/${adId}/trades`, { amount });
    }

    static async getMyTrades(role?: string): Promise<any[]> {
        const url = role ? `/api/v1/marketplace/trades?role=${role}` : `/api/v1/marketplace/trades`;
        const data = await ApiClient.get<{ trades: any[] }>(url);
        return data.trades || [];
    }

    static async getTradeDetails(tradeId: string): Promise<any> {
        const data = await ApiClient.get<{ trade: any }>(`/api/v1/marketplace/trades/${tradeId}`);
        return data.trade || null;
    }

    static async updateTradeStatus(tradeId: string, action: string, reason?: string): Promise<any> {
        return ApiClient.patch(`/api/v1/marketplace/trades/${tradeId}/status`, { action, reason });
    }

    static async getTradeMessages(tradeId: string): Promise<any[]> {
        const data = await ApiClient.get<{ messages: any[] }>(`/api/v1/marketplace/trades/${tradeId}/messages`);
        return data.messages || [];
    }

    static async sendTradeMessage(tradeId: string, content: string): Promise<any> {
        return ApiClient.post(`/api/v1/marketplace/trades/${tradeId}/messages`, { content });
    }

    static async disputeTrade(tradeId: string, reasonCode: string, reasonText: string): Promise<any> {
        return ApiClient.post(`/api/v1/marketplace/trades/${tradeId}/dispute`, { 
            reason_code: reasonCode, 
            reason_text: reasonText 
        });
    }

    static async leaveFeedback(tradeId: string, rating: string, comment: string): Promise<any> {
        return ApiClient.post(`/api/v1/marketplace/trades/${tradeId}/feedback`, { rating, comment });
    }
}
