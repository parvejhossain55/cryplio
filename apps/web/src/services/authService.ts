// Types for backend API responses
export interface BackendUser {
    id: string;
    email: string;
    username: string;
    bio?: string;
    email_verified: boolean;
    kyc_level: number;
    is_merchant: boolean;
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

export interface UserBlock {
    id: string;
    blocker_id: string;
    blocked_id: string;
    reason?: string;
    is_permanent: boolean;
    expires_at?: string;
    created_at: string;
}

export interface AdResponse {
    ad_id: string;
    user_id: string;
    username: string;
    user_avatar?: string;
    user_rating: number;
    user_trades: number;
    type: "buy" | "sell";
    crypto_symbol: string;
    fiat_symbol: string;
    price_type: string;
    price: number;
    min_amount: number;
    max_amount: number;
    payment_methods: string[];
    payment_window_minutes: number;
    is_online: boolean;
}

export interface AuthResponse {
    user: BackendUser;
    token?: string;
    refresh_token?: string;
    requires_2fa?: boolean;
    temp_token?: string;
    error?: string;
}

// Interceptor for auto-refresh
const AUTH_SESSION_KEY = "cryplio_has_auth_session";
let isRefreshing = false;
let refreshSubscribers: Array<{
    resolve: () => void;
    reject: (error: unknown) => void;
}> = [];

const rememberAuthSession = () => {
    localStorage.setItem(AUTH_SESSION_KEY, "true");
};

const forgetAuthSession = () => {
    localStorage.removeItem(AUTH_SESSION_KEY);
};

const hasRememberedAuthSession = () => {
    return localStorage.getItem(AUTH_SESSION_KEY) === "true";
};

export const authService = {
    login: async (email: string, password: string): Promise<BackendUser> => {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password }),
            credentials: "include",
        });

        const data: AuthResponse = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Login failed");
        }

        // Check for 2FA requirement
        if (data.requires_2fa) {
            if (data.temp_token) {
                sessionStorage.setItem("2fa_temp_token", data.temp_token);
                sessionStorage.setItem("2fa_user_id", data.user.id);
            }
            const error = Object.assign(new Error("2FA_REQUIRED"), {
                requires2FA: true,
                tempToken: data.temp_token,
                user: data.user,
            });
            throw error;
        }

        rememberAuthSession();
        return data.user;
    },

    register: async (
        email: string,
        username: string,
        password: string
    ): Promise<BackendUser> => {
        const response = await fetch("/api/auth/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, username, password }),
            credentials: "include",
        });

        const data: AuthResponse = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Registration failed");
        }

        rememberAuthSession();
        return data.user;
    },

    logout: async (): Promise<void> => {
        try {
            await fetch("/api/auth/logout", {
                method: "POST",
                credentials: "include",
            });
        } finally {
            // Always clear local state
            forgetAuthSession();
            localStorage.clear();
            sessionStorage.clear();
        }
    },

    getCurrentUser: async (): Promise<BackendUser | null> => {
        try {
            const response = await fetch("/api/users/me", {
                credentials: "include",
            });

            if (!response.ok) {
                if (response.status === 401) {
                    if (!hasRememberedAuthSession()) {
                        return null;
                    }

                    // Try to refresh once
                    try {
                        await authService.refreshToken();
                        // Retry getting user
                        const retryResponse = await fetch("/api/users/me", {
                            credentials: "include",
                        });
                        if (retryResponse.ok) {
                            const data = await retryResponse.json();
                            rememberAuthSession();
                            return data.user;
                        }
                    } catch {
                        // No refresh cookie is a normal guest state on first load.
                        forgetAuthSession();
                    }
                }
                return null;
            }

            const data = await response.json();
            if (data.user) {
                rememberAuthSession();
            }
            return data.user;
        } catch {
            return null;
        }
    },

    updateCurrentUser: async (updates: {
        username?: string;
        bio?: string;
        avatarUrl?: string | null;
    }): Promise<BackendUser> => {
        const response = await fetch("/api/users/me", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(updates),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Failed to update user");
        }

        rememberAuthSession();
        return data.user;
    },

    loginWithGoogle: (): void => {
        window.location.href = "/api/auth/google";
    },

    // Password reset
    requestPasswordReset: async (email: string): Promise<void> => {
        const response = await fetch("/api/auth/password/reset-request", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Password reset request failed");
        }
    },

    resetPassword: async (token: string, newPassword: string): Promise<void> => {
        const response = await fetch("/api/auth/password/reset", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ token, password: newPassword }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Password reset failed");
        }
    },

    // 2FA methods
    setup2FA: async (): Promise<{ secret: string; provisioning_uri: string }> => {
        const response = await fetch("/api/auth/2fa/setup", {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "2FA setup failed");
        }

        return response.json();
    },

    verify2FA: async (code: string): Promise<void> => {
        const response = await fetch("/api/auth/2fa/verify", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ code }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "2FA verification failed");
        }
    },

    disable2FA: async (password: string): Promise<void> => {
        const response = await fetch("/api/auth/2fa/disable", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ password }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "2FA disable failed");
        }
    },

    complete2FALogin: async (tempToken: string, code: string): Promise<BackendUser> => {
        const response = await fetch("/api/auth/2fa/complete-login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ temp_token: tempToken, code }),
            credentials: "include",
        });

        const data: AuthResponse = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "2FA login failed");
        }

        rememberAuthSession();
        return data.user;
    },

    // Email verification
    requestEmailVerification: async (userId: string): Promise<void> => {
        const response = await fetch("/api/auth/email/request", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ user_id: userId }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Verification request failed");
        }
    },

    verifyEmail: async (token: string): Promise<void> => {
        const response = await fetch("/api/auth/email/verify", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ token }),
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Email verification failed");
        }
    },

    // Sessions
    getSessions: async (): Promise<BackendSession[]> => {
        try {
            const response = await fetch("/api/sessions", {
                credentials: "include",
            });

            if (!response.ok) {
                // Try to refresh once
                try {
                    await authService.refreshToken();
                    const retryResponse = await fetch("/api/sessions", {
                        credentials: "include",
                    });
                    if (retryResponse.ok) {
                        const data = await retryResponse.json();
                        return data.sessions || [];
                    }
                } catch (refreshError) {
                    console.error("Token refresh failed:", refreshError);
                }
                throw new Error("Failed to get sessions");
            }

            const data = await response.json();
            return data.sessions || [];
        } catch (error) {
            throw error;
        }
    },

    revokeSession: async (tokenId: string): Promise<void> => {
        const response = await fetch(`/api/sessions/${tokenId}`, {
            method: "DELETE",
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Failed to revoke session");
        }
    },

    // Token refresh
    refreshToken: async (): Promise<void> => {
        if (isRefreshing) {
            // Wait for refresh to complete
            return new Promise((resolve, reject) => {
                refreshSubscribers.push({ resolve, reject });
            });
        }

        isRefreshing = true;
        let refreshError: unknown;
        try {
            const response = await fetch("/api/auth/refresh", {
                method: "POST",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Token refresh failed");
            }

            // Refresh token automatically saved by browser via Set-Cookie header
        } catch (error) {
            refreshError = error;
            throw error;
        } finally {
            isRefreshing = false;
            refreshSubscribers.forEach(({ resolve, reject }) => {
                if (refreshError) {
                    reject(refreshError);
                    return;
                }
                resolve();
            });
            refreshSubscribers = [];
        }
    },

    // User Profile & Stats
    getUserByUsername: async (username: string): Promise<{ user: BackendUser; stats: UserStats }> => {
        const response = await fetch(`/api/users/username/${username}`, {
            credentials: "include",
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || "User not found");
        }

        const data = await response.json();
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

    // User Blocking
    blockUser: async (blockedId: string, reason: string = "No reason provided", isPermanent: boolean = true): Promise<void> => {
        const response = await fetchWithRefresh("/api/users/me/block", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                blocked_id: blockedId,
                reason: reason,
                is_permanent: isPermanent
            }),
            credentials: "include",
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || "Failed to block user");
        }
    },

    unblockUser: async (blockedId: string): Promise<void> => {
        const response = await fetchWithRefresh(`/api/users/me/block/${blockedId}`, {
            method: "DELETE",
            credentials: "include",
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || "Failed to unblock user");
        }
    },

    getBlocks: async (): Promise<UserBlock[]> => {
        const response = await fetchWithRefresh("/api/users/me/block", {
            credentials: "include",
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || "Failed to fetch blocked users");
        }

        const data = await response.json();
        return data.blocks || [];
    },

    // Marketplace
    getAds: async (params: { type?: string; cryptoId?: number; fiatId?: number } = {}): Promise<AdResponse[]> => {
        const queryParams = new URLSearchParams();
        if (params.type) queryParams.append("type", params.type);
        if (params.cryptoId) queryParams.append("crypto_id", params.cryptoId.toString());
        if (params.fiatId) queryParams.append("fiat_id", params.fiatId.toString());

        const response = await fetch(`/api/v1/marketplace/ads?${queryParams.toString()}`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch advertisements");
        }

        const data = await response.json();
        return data.ads || [];
    },

    initiateTrade: async (adId: string, amount: number): Promise<{ trade_id: string; message: string }> => {
        const response = await fetch(`/api/v1/marketplace/ads/${adId}/trades`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ amount }),
            credentials: "include",
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to initiate trade");
        }

        return await response.json();
    },

    createAd: async (adData: any): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/ads`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(adData),
            credentials: "include",
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to create advertisement");
        }

        return await response.json();
    },

    getMyAds: async (): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/my-ads`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch your advertisements");
        }

        return await response.json();
    },

    toggleAdStatus: async (adId: string): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/ads/${adId}/status`, {
            method: "PATCH",
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to update advertisement status");
        }

        return await response.json();
    },

    getMyTrades: async (role?: string): Promise<any> => {
        const url = role ? `/api/v1/marketplace/trades?role=${role}` : `/api/v1/marketplace/trades`;
        const response = await fetch(url, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch your trades");
        }

        return await response.json();
    },

    getTradeDetails: async (tradeId: string): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/trades/${tradeId}`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch trade details");
        }

        return await response.json();
    },

    updateTradeStatus: async (tradeId: string, action: string, reason?: string): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/trades/${tradeId}/status`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ action, reason }),
            credentials: "include",
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || "Failed to update trade status");
        }

        return await response.json();
    },

    getTradeMessages: async (tradeId: string): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/trades/${tradeId}/messages`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch messages");
        }

        return await response.json();
    },

    sendTradeMessage: async (tradeId: string, content: string): Promise<any> => {
        const response = await fetch(`/api/v1/marketplace/trades/${tradeId}/messages`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ content }),
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to send message");
        }

        return await response.json();
    },
};

// Wrap fetch to automatically refresh on 401
export const fetchWithRefresh = async (url: string, options: RequestInit = {}): Promise<Response> => {
    let response = await fetch(url, options);

    if (response.status === 401) {
        // Try to refresh token once
        const originalBody = options.body;
        try {
            await authService.refreshToken();
            // Retry original request with same body
            response = await fetch(url, {
                ...options,
                body: originalBody,
            });
        } catch {
            console.error("Refresh failed, user must re-login");
            // Logout user on refresh failure
            forgetAuthSession();
            localStorage.clear();
            sessionStorage.clear();
            window.location.href = "/login";
            throw new Error("Session expired. Please login again.");
        }
    }

    return response;
};
