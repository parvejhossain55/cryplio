"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { authService, BackendUser } from "@/services/authService";

interface User {
    id: string;
    name: string;
    email: string;
    username: string;
    role: "user" | "merchant" | "admin" | null;
    kycLevel: number;
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
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        const checkSession = async () => {
            try {
                const currentUser: BackendUser | null = await authService.getCurrentUser();
                if (currentUser) {
                    // Map the backend user to our frontend user format
                    const mappedUser: User = {
                        id: currentUser.id,
                        name: currentUser.username || currentUser.email.split("@")[0],
                        email: currentUser.email,
                        username: currentUser.username || currentUser.email.split("@")[0],
                        role: "user",
                        kycLevel: currentUser.kyc_level ?? 0,
                    };
                    setUser(mappedUser);
                }
            } catch (error) {
                console.error("Failed to check session:", error);
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
        role: "user",
        kycLevel: backendUser.kyc_level ?? 0,
    });

    const login = async (email: string, password: string) => {
        setIsLoading(true);
        try {
            const backendUser: BackendUser = await authService.login(email, password);
            // Map the backend user to our frontend user format
            setUser(mapBackendUser(backendUser));

            // Redirect to dashboard
            router.push("/user/dashboard");
        } catch (error) {
            throw error;
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
            router.push("user/dashboard");
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
            setIsLoading(false);
            router.push("/login");
        }
    };

    return (
        <AuthContext.Provider value={{ user, login, loginWithGoogle, register, logout, isLoading }}>
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
