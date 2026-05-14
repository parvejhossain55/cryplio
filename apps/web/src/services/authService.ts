import { ApiClient } from "./apiClient";

// Types for backend API responses
export interface BackendUser {
    id: string;
    email: string;
    username: string;
    role: string;
    bio?: string;
    email_verified: boolean;
    two_fa_enabled: boolean;
    avatar_url?: string | null;
    is_online: boolean;
}

export interface BackendSession {
    id?: string;
    token_id: string;
    device_type?: string;
    ip_address?: string;
    last_used_at: string;
}

export interface UserStats {
    total_trades: number;
    successful_trades: number;
    dispute_rate: number;
    avg_rating?: number;
    positive_feedback_count: number;
    neutral_feedback_count: number;
    negative_feedback_count: number;
    total_volume_usd: number;
    last_trade_at?: string;
    success_rate: number;
}

export interface UserPaymentMethod {
    id: string;
    user_id: string;
    payment_method_code: string;
    display_name: string;
    account_name?: string;
    account_number?: string;
    bank_name?: string;
    is_active: boolean;
    is_default: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreatePaymentMethodRequest {
    payment_method_code: string;
    display_name: string;
    account_name?: string;
    account_number?: string;
    bank_name?: string;
}

// Auth response types remain here for now as they are core to identity

// Wallet related interfaces moved to walletService.ts

export interface AuthResponse {
    user?: BackendUser;
    token?: string;
    requires_2fa?: boolean;
    temp_token?: string;
    error?: string;
}

export interface LoginData {
    email: string;
    password?: string;
}

export interface RegisterData {
    email: string;
    username: string;
    password?: string;
}

let isRefreshing = false;
let refreshSubscribers: any[] = [];

// Helper to remember session in local storage for AuthContext
const rememberAuthSession = () => {
    if (typeof window !== 'undefined') {
        localStorage.setItem('auth_session', 'active');
        localStorage.setItem('last_active', Date.now().toString());
    }
};

const clearAuthSession = () => {
    if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_session');
        localStorage.removeItem('last_active');
        localStorage.removeItem('user_id');
    }
};

export const authService = {
    // Auth methods
    register: async (data: RegisterData): Promise<BackendUser> => {
        const res = await ApiClient.post<AuthResponse>("/api/v1/auth/register", data);
        rememberAuthSession();
        return res.user!;
    },

    login: async (data: LoginData): Promise<BackendUser | { temp_token: string }> => {
        const res = await ApiClient.post<AuthResponse>("/api/v1/auth/login", data);
        if (res.user) {
            rememberAuthSession();
            if (res.user.id) {
                localStorage.setItem('user_id', res.user.id);
            }
            return res.user;
        }
        return { temp_token: res.temp_token! };
    },

    logout: async (): Promise<void> => {
        await ApiClient.post("/api/v1/auth/logout");
        clearAuthSession();
    },

    getCurrentUser: async (): Promise<BackendUser | null> => {
        try {
            const res = await ApiClient.get<{ user: BackendUser }>("/api/v1/users/me");
            if (res.user && res.user.id) {
                localStorage.setItem('user_id', res.user.id);
            }
            return res.user;
        } catch (error: any) {
            // Return null for "user not found" or unauthorized errors
            if (error.message?.includes("user not found") || error.message?.includes("Unauthorized")) {
                return null;
            }
            throw error;
        }
    },

    updateCurrentUser: async (updates: any): Promise<BackendUser> => {
        const res = await ApiClient.put<{ user: BackendUser }>("/api/v1/users/me", updates);
        rememberAuthSession();
        return res.user;
    },

    loginWithGoogle: (): void => {
        window.location.href = "/api/v1/auth/google";
    },

    // Password reset
    requestPasswordReset: async (email: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/password/reset-request", { email });
    },

    resetPassword: async (token: string, newPassword: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/password/reset", { token, password: newPassword });
    },

    // 2FA methods
    setup2FA: async (): Promise<{ secret: string; provisioning_uri: string }> => {
        return ApiClient.post("/api/v1/auth/2fa/setup");
    },

    verify2FA: async (code: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/2fa/verify", { code });
    },

    disable2FA: async (password: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/2fa/disable", { password });
    },

    complete2FALogin: async (tempToken: string, code: string): Promise<BackendUser> => {
        const res = await ApiClient.post<AuthResponse>("/api/v1/auth/2fa/complete-login", {
            temp_token: tempToken,
            code
        });
        rememberAuthSession();
        return res.user!;
    },

    // Email verification
    requestEmailVerification: async (userId: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/email/request", { user_id: userId });
    },

    verifyEmail: async (token: string): Promise<void> => {
        await ApiClient.post("/api/v1/auth/email/verify", { token });
    },

    // Sessions
    getSessions: async (): Promise<BackendSession[]> => {
        const res = await ApiClient.get<{ sessions: BackendSession[] }>("/api/v1/sessions");
        return res.sessions || [];
    },

    revokeSession: async (tokenId: string): Promise<void> => {
        await ApiClient.delete(`/api/v1/sessions/${tokenId}`);
    },

    // Token refresh
    refreshToken: async (): Promise<void> => {
        if (isRefreshing) {
            return new Promise((resolve, reject) => {
                refreshSubscribers.push({ resolve, reject });
            });
        }

        isRefreshing = true;
        let refreshError: any = null;
        try {
            const res = await fetch("/api/v1/auth/refresh", {
                method: "POST",
                credentials: "include",
            });
            if (!res.ok) {
                throw new Error("Refresh token invalid or expired");
            }
        } catch (error) {
            refreshError = error;
            throw error;
        } finally {
            isRefreshing = false;
            refreshSubscribers.forEach(({ resolve, reject }) => {
                if (refreshError) reject(refreshError);
                else resolve();
            });
            refreshSubscribers = [];
        }
    },

    // Profile & Stats
    getUserByUsername: async (username: string): Promise<{ user: BackendUser; stats: UserStats }> => {
        return ApiClient.get(`/api/v1/users/username/${username}`);
    },

    getUserStats: async (userId: string): Promise<UserStats> => {
        return ApiClient.get(`/api/v1/users/${userId}/stats`);
    },

    // Payment Methods
    getUserPaymentMethods: async (userId: string): Promise<UserPaymentMethod[]> => {
        const res = await ApiClient.get<{ payment_methods: UserPaymentMethod[] }>(`/api/v1/users/${userId}/payment-methods`);
        return res.payment_methods || [];
    },

    getMyPaymentMethods: async (): Promise<UserPaymentMethod[]> => {
        const res = await ApiClient.get<{ payment_methods: UserPaymentMethod[] }>("/api/v1/users/me/payment-methods");
        return res.payment_methods || [];
    },

    createPaymentMethod: async (data: CreatePaymentMethodRequest): Promise<UserPaymentMethod> => {
        return ApiClient.post("/api/v1/users/me/payment-methods", data);
    },

    updatePaymentMethod: async (id: string, data: Partial<CreatePaymentMethodRequest>): Promise<UserPaymentMethod> => {
        return ApiClient.put(`/api/v1/users/me/payment-methods/${id}`, data);
    },

    deletePaymentMethod: async (id: string): Promise<void> => {
        await ApiClient.delete(`/api/v1/users/me/payment-methods/${id}`);
    },

    setDefaultPaymentMethod: async (id: string): Promise<void> => {
        await ApiClient.patch(`/api/v1/users/me/payment-methods/${id}/default`);
    },

    // Block/Report features removed per previous architectural decisions

    // Admin Methods
    getAdminDisputes: async (): Promise<any[]> => {
        const res = await ApiClient.get<{ disputes: any[] }>("/api/v1/admin/disputes");
        return res.disputes || [];
    },

    assignDispute: async (disputeId: string): Promise<void> => {
        await ApiClient.post(`/api/v1/admin/disputes/${disputeId}/assign`);
    },

    resolveDispute: async (disputeId: string, resolution: string, note: string): Promise<void> => {
        await ApiClient.post(`/api/v1/admin/disputes/${disputeId}/resolve`, {
            resolution,
            note
        });
    },
};
