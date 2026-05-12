// Base API Client for Cryplio Frontend

let isRefreshing = false;
let refreshSubscribers: Array<{
    resolve: () => void;
    reject: (error: unknown) => void;
}> = [];

const AUTH_SESSION_KEY = "cryplio_has_auth_session";

export const hasRememberedAuthSession = (): boolean => {
    if (typeof window === "undefined") return false;
    return localStorage.getItem(AUTH_SESSION_KEY) === "true";
};

export const rememberAuthSession = (): void => {
    if (typeof window === "undefined") return;
    localStorage.setItem(AUTH_SESSION_KEY, "true");
};

export const forgetAuthSession = (): void => {
    if (typeof window === "undefined") return;
    localStorage.removeItem(AUTH_SESSION_KEY);
};

export const refreshToken = async (): Promise<void> => {
    if (isRefreshing) {
        return new Promise((resolve, reject) => {
            refreshSubscribers.push({ resolve, reject });
        });
    }

    isRefreshing = true;
    let refreshError: unknown;
    try {
        const response = await fetch("/api/v1/auth/refresh", {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Token refresh failed");
        }
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
};

export const fetchWithRefresh = async (url: string, options: RequestInit = {}): Promise<Response> => {
    let response = await fetch(url, options);

    if (response.status === 401) {
        // Try to refresh token once
        const originalBody = options.body;
        try {
            await refreshToken();
            // Retry original request
            response = await fetch(url, {
                ...options,
                body: originalBody,
            });
        } catch {
            console.error("Refresh failed, user must re-login");
            forgetAuthSession();
            if (typeof window !== "undefined" && window.location.pathname !== "/login") {
                localStorage.clear();
                sessionStorage.clear();
                window.location.href = "/login";
            }
            throw new Error("Session expired. Please login again.");
        }
    }

    return response;
};

export const handleResponse = async <T>(response: Response): Promise<T> => {
    const data = await response.json();
    if (!response.ok) {
        throw new Error(data.error || "Request failed");
    }
    return data;
};
