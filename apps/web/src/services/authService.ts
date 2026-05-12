import { AuthResponse, BackendUser } from "@/types/api";
import { 
    fetchWithRefresh, 
    handleResponse, 
    rememberAuthSession, 
    forgetAuthSession,
    refreshToken
} from "./apiClient";

export const authService = {
    register: async (email: string, username: string, password: string): Promise<BackendUser> => {
        const response = await fetch("/api/v1/auth/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, username, password }),
            credentials: "include",
        });

        const data: AuthResponse = await handleResponse<AuthResponse>(response);
        rememberAuthSession();
        return data.user;
    },

    login: async (email: string, password: string): Promise<AuthResponse> => {
        const response = await fetch("/api/v1/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
            credentials: "include",
        });

        const data: AuthResponse = await handleResponse<AuthResponse>(response);
        if (!data.requires_2fa) {
            rememberAuthSession();
        }
        return data;
    },

    logout: async (): Promise<void> => {
        try {
            await fetch("/api/v1/auth/logout", {
                method: "POST",
                credentials: "include",
            });
        } finally {
            forgetAuthSession();
            if (typeof window !== "undefined") {
                localStorage.clear();
                sessionStorage.clear();
            }
        }
    },

    refreshToken,

    loginWithGoogle: (): void => {
        if (typeof window !== "undefined") {
            window.location.href = "/api/v1/auth/google";
        }
    },

    // Password reset
    requestPasswordReset: async (email: string): Promise<void> => {
        const response = await fetch("/api/v1/auth/password/reset-request", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email }),
            credentials: "include",
        });
        await handleResponse(response);
    },

    resetPassword: async (token: string, newPassword: string): Promise<void> => {
        const response = await fetch("/api/v1/auth/password/reset", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ token, password: newPassword }),
            credentials: "include",
        });
        await handleResponse(response);
    },

    // 2FA methods
    setup2FA: async (): Promise<{ secret: string; provisioning_uri: string }> => {
        const response = await fetchWithRefresh("/api/v1/auth/2fa/setup", {
            method: "POST",
            credentials: "include",
        });
        return handleResponse(response);
    },

    verify2FA: async (code: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/auth/2fa/verify", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ code }),
            credentials: "include",
        });
        await handleResponse(response);
    },

    disable2FA: async (password: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/auth/2fa/disable", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ password }),
            credentials: "include",
        });
        await handleResponse(response);
    },

    complete2FALogin: async (tempToken: string, code: string): Promise<BackendUser> => {
        const response = await fetch("/api/v1/auth/2fa/complete-login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ temp_token: tempToken, code }),
            credentials: "include",
        });

        const data: AuthResponse = await handleResponse<AuthResponse>(response);
        rememberAuthSession();
        return data.user;
    },

    // Email verification
    requestEmailVerification: async (userId: string): Promise<void> => {
        const response = await fetchWithRefresh("/api/v1/auth/email/request", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ user_id: userId }),
            credentials: "include",
        });
        await handleResponse(response);
    },

    verifyEmail: async (token: string): Promise<void> => {
        const response = await fetch("/api/v1/auth/email/verify", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ token }),
            credentials: "include",
        });
        await handleResponse(response);
    },
};
