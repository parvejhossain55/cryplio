"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { authService, BackendUser } from "@/services/authService";

export interface User {
    id: string;
    name: string;
    email: string;
    username: string;
    bio?: string;
    avatarUrl?: string;
    role: "user" | "admin" | null;
    emailVerified: boolean;
    twoFAEnabled: boolean;
    lastSeenAt?: string;
    isOnline: boolean;
}

interface AuthContextType {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    loginWithGoogle: () => void;
    register: (
        email: string,
        username: string,
        password: string
    ) => Promise<void>;
    logout: () => Promise<void>;
    isLoading: boolean;
    requires2FA: boolean;
    setRequires2FA: (requires: boolean) => void;
    temp2FAToken: string | null;
    setTemp2FAToken: (token: string | null) => void;
    complete2FALogin: (code: string) => Promise<void>;
    refreshUser: () => Promise<void>;
}

interface TwoFactorLoginError extends Error {
    requires2FA?: boolean;
    tempToken?: string;
    user?: BackendUser;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [requires2FA, setRequires2FA] = useState(false);
    const [temp2FAToken, setTemp2FAToken] = useState<string | null>(null);
    const router = useRouter();

    useEffect(() => {
        const checkSession = async () => {
            try {
                const currentUser: BackendUser | null = await authService.getCurrentUser();
                if (currentUser) {
                    setUser(mapBackendUser(currentUser));
                    localStorage.setItem("user_id", currentUser.id);
                } else {
                    localStorage.removeItem("user_id");
                }
            } catch (error) {
                console.error("Failed to check session:", error);
                localStorage.removeItem("user_id");
            } finally {
                setIsLoading(false);
            }
        };

        checkSession();
    }, []);

    const mapBackendUser = (backendUser: BackendUser): User => ({
        id: backendUser.id,
        name: backendUser.username || backendUser.email.split("@")[0],
        email: backendUser.email,
        username: backendUser.username || backendUser.email.split("@")[0],
        role: backendUser.username === "admin" ? "admin" : "user",
        emailVerified: backendUser.email_verified ?? false,
        twoFAEnabled: backendUser.two_fa_enabled ?? false,
        bio: backendUser.bio ?? "",
        avatarUrl: backendUser.avatar_url ?? undefined,
        isOnline: backendUser.is_online ?? false,
    });

    const login = async (email: string, password: string) => {
        setIsLoading(true);
        try {
            const backendUser: BackendUser = await authService.login(email, password);
            setUser(mapBackendUser(backendUser));
            localStorage.setItem("user_id", backendUser.id);
            router.push("/user/dashboard");
        } catch (error) {
            const loginError = error as TwoFactorLoginError;
            // Check for 2FA requirement
            if (loginError.requires2FA || loginError.message === "2FA_REQUIRED") {
                setRequires2FA(true);
                setTemp2FAToken(loginError.tempToken || sessionStorage.getItem("2fa_temp_token") || null);
                // Store in sessionStorage as backup
                if (loginError.tempToken && loginError.user) {
                    sessionStorage.setItem("2fa_temp_token", loginError.tempToken);
                    sessionStorage.setItem("2fa_user_id", loginError.user.id);
                    localStorage.setItem("user_id", loginError.user.id);
                }
                // Redirect to 2FA verification page
                router.push("/2fa-verify");
                return;
            }
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const complete2FALogin = async (code: string) => {
        setIsLoading(true);
        try {
            const tempToken = temp2FAToken || sessionStorage.getItem("2fa_temp_token");
            if (!tempToken) {
                throw new Error("2FA session expired. Please login again.");
            }

            const backendUser: BackendUser = await authService.complete2FALogin(tempToken, code);
            setUser(mapBackendUser(backendUser));
            localStorage.setItem("user_id", backendUser.id);
            setRequires2FA(false);
            setTemp2FAToken(null);
            sessionStorage.removeItem("2fa_temp_token");
            sessionStorage.removeItem("2fa_user_id");
            router.push("/user/dashboard");
        } finally {
            setIsLoading(false);
        }
    };

    const loginWithGoogle = () => {
        authService.loginWithGoogle();
    };

    const register = async (
        email: string,
        username: string,
        password: string
    ) => {
        setIsLoading(true);
        try {
            const backendUser: BackendUser = await authService.register(
                email,
                username,
                password
            );
            setUser(mapBackendUser(backendUser));
            router.push("/user/dashboard");
        } finally {
            setIsLoading(false);
        }
    };

    const logout = async () => {
        setIsLoading(true);
        try {
            await authService.logout();
        } catch (error) {
            console.error("Logout failed:", error);
        } finally {
            setUser(null);
            setRequires2FA(false);
            setTemp2FAToken(null);
            sessionStorage.removeItem("2fa_temp_token");
            sessionStorage.removeItem("2fa_user_id");
            setIsLoading(false);
            router.push("/login");
        }
    };

    const refreshUser = async () => {
        try {
            const currentUser: BackendUser | null = await authService.getCurrentUser();
            if (currentUser) {
                setUser(mapBackendUser(currentUser));
            } else {
                setUser(null);
            }
        } catch {
            setUser(null);
        }
    };

    return (
        <AuthContext.Provider value={{
            user,
            login,
            loginWithGoogle,
            register,
            logout,
            isLoading,
            requires2FA,
            setRequires2FA,
            temp2FAToken,
            setTemp2FAToken,
            complete2FALogin,
            refreshUser
        }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
