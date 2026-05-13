import { ApiClient } from "./apiClient";
import { WalletBalance } from "./authService";

export class WalletService {
    static async getBalances(): Promise<WalletBalance[]> {
        return ApiClient.get(`/api/v1/wallet/balance`);
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
