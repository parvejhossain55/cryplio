import { ApiClient } from "./apiClient";
export interface WalletBalance {
    wallet_id: string;
    user_id: string;
    crypto_id: number | null;
    crypto_symbol: string | null;
    address: string;
    balance: number;
    locked_balance: number;
    is_active: boolean;
}

export interface WalletTransaction {
    tx_id: string;
    wallet_id: string;
    type: "deposit" | "withdrawal";
    amount: number;
    status: string;
    created_at: string;
}

export class WalletService {
    static async getBalances(): Promise<WalletBalance[]> {
        const data = await ApiClient.get<{ wallets: WalletBalance[] }>(`/api/v1/wallet/balance`);
        return data.wallets || [];
    }

    static async getTransactions(params: { limit: number; offset: number }): Promise<any> {
        const query = new URLSearchParams(params as any).toString();
        return ApiClient.get(`/api/v1/wallet/transactions?${query}`);
    }

    static async getDepositAddress(crypto: string): Promise<{ address: string }> {
        return ApiClient.get(`/api/v1/wallet/deposit/${crypto}`);
    }

    static async withdraw(data: { crypto: string; amount: number; address: string }): Promise<any> {
        return ApiClient.post(`/api/v1/wallet/withdraw`, data);
    }
}
