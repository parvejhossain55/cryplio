import { BackendUser, BackendSession, UserStats, UserPaymentMethod, CreatePaymentMethodRequest } from "@/types/api";
import { fetchWithRefresh, handleResponse, rememberAuthSession } from "./apiClient";

export const userService = {
    getCurrentUser: async (): Promise<BackendUser> => {
        const response = await fetchWithRefresh("/api/v1/users/me", {
            credentials: "include",
        });
        const data = await handleResponse<{ user: BackendUser }>(response);
        rememberAuthSession();
        return data.user;
    },

    updateCurrentUser: async (updates: {
        username?: string;
        bio?: string;
        avatarUrl?: string | null;
    }): Promise<BackendUser> => {
        const response = await fetchWithRefresh("/api/v1/users/me", {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(updates),
            credentials: "include",
        });
        const data = await handleResponse<{ user: BackendUser }>(response);
        rememberAuthSession();
        return data.user;
    },

    getUserByUsername: async (username: string): Promise<{ user: BackendUser; stats: UserStats }> => {
        const response = await fetch("/api/v1/users/username/" + username, {
            credentials: "include",
        });
        const data = await handleResponse<{ user: BackendUser; stats: any }>(response);
        
        return {
            user: data.user,
            stats: {
                ...data.stats,
                success_rate: data.stats?.total_trades > 0
                    ? (data.stats.successful_trades / data.stats.total_trades) * 100
                    : 100
            }
        };
    },

    getSessions: async (): Promise<BackendSession[]> => {
        const response = await fetchWithRefresh("/api/v1/sessions", {
            credentials: "include",
        });
        const data = await handleResponse<{ sessions: BackendSession[] }>(response);
        return data.sessions || [];
    },

    revokeSession: async (tokenId: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/sessions/" + tokenId, {
            method: "DELETE",
            credentials: "include",
        });
        await handleResponse(response);
    },

    // Payment Methods
    getPaymentMethods: async (): Promise<UserPaymentMethod[]> => {
        const response = await fetchWithRefresh("/api/v1/users/me/payment-methods", {
            credentials: "include",
        });
        const data = await handleResponse<{ payment_methods: UserPaymentMethod[] }>(response);
        return data.payment_methods || [];
    },

    createPaymentMethod: async (req: CreatePaymentMethodRequest): Promise<UserPaymentMethod> => {
        const response = await fetchWithRefresh("/api/v1/users/me/payment-methods", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(req),
            credentials: "include",
        });
        const data = await handleResponse<{ payment_method: UserPaymentMethod }>(response);
        return data.payment_method;
    },

    updatePaymentMethod: async (id: string, req: Partial<CreatePaymentMethodRequest>): Promise<UserPaymentMethod> => {
        const response = await fetchWithRefresh("/api/v1/users/me/payment-methods/" + id, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(req),
            credentials: "include",
        });
        const data = await handleResponse<{ payment_method: UserPaymentMethod }>(response);
        return data.payment_method;
    },

    deletePaymentMethod: async (id: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/users/me/payment-methods/" + id, {
            method: "DELETE",
            credentials: "include",
        });
        await handleResponse(response);
    },

    setDefaultPaymentMethod: async (id: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/users/me/payment-methods/" + id + "/default", {
            method: "PATCH",
            credentials: "include",
        });
        await handleResponse(response);
    },
};
