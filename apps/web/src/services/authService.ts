// Types for backend API responses
export interface BackendUser {
    id: string;
    email: string;
    username: string;
    kyc_level: number;
}

export const authService = {
    login: async (email: string, password: string): Promise<BackendUser> => {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password }),
            credentials: "include", // Important for cookies
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Login failed");
        }

        const data = await response.json();
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

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Registration failed");
        }

        const data = await response.json();
        return data.user;
    },

    logout: async (): Promise<void> => {
        const response = await fetch("/api/auth/logout", {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || "Logout failed");
        }
    },

    getCurrentUser: async (): Promise<BackendUser | null> => {
        try {
            const response = await fetch("/api/users/me", {
                credentials: "include",
            });

            if (!response.ok) {
                return null;
            }

            const data = await response.json();
            return data.user;
        } catch {
            return null;
        }
    },

    loginWithGoogle: (): void => {
        window.location.href = "/api/auth/google";
    },
};
