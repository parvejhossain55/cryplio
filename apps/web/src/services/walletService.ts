import { WalletBalance, WalletTransaction } from "@/types/api";
import { fetchWithRefresh, handleResponse } from "./apiClient";

export const walletService = {
    getBalances: async (): Promise<WalletBalance[]> => {
        const response = await fetchWithRefresh("/api/v1/wallet/balance", {
            credentials: "include",
        });
        const data = await handleResponse<{ wallets: WalletBalance[] }>(response);
        return data.wallets || [];
    },

    getDepositAddress: async (cryptoSymbol: string): Promise<{ wallet_id: string; crypto_id: number; address: string }> => {
        const response = await fetchWithRefresh(`/api/v1/wallet/deposit/${encodeURIComponent(cryptoSymbol)}`, {
            credentials: "include",
        });
        return handleResponse(response);
    },

    withdraw: async (payload: { 
        crypto_symbol: string; 
        address: string; 
        amount: number; 
        two_fa_code: string; 
        fee?: number; 
        memo?: string 
    }): Promise<any> => {
        const response = await fetchWithRefresh("/api/v1/wallet/withdraw", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
            credentials: "include",
        });
        return handleResponse(response);
    },

    getTransactions: async (params: { limit?: number; offset?: number } = {}): Promise<{ transactions: WalletTransaction[]; total: number }> => {
        const query = new URLSearchParams();
        if (params.limit !== undefined) query.append("limit", String(params.limit));
        if (params.offset !== undefined) query.append("offset", String(params.offset));
        const suffix = query.toString() ? `?${query.toString()}` : "";

        const response = await fetchWithRefresh(`/api/v1/wallet/transactions${suffix}`, {
            credentials: "include",
        });
        const data = await handleResponse<{ transactions: WalletTransaction[]; total: number }>(response);
        return {
            transactions: data.transactions || [],
            total: data.total || 0,
        };
    },
};
