export const fetchWithRefresh = async (url: string, options: RequestInit = {}): Promise<Response> => {
    let response = await fetch(url, options);

    if (response.status === 401) {
        // Import authService dynamically to avoid circular dependency
        const { authService } = await import("./authService");
        
        try {
            await authService.refreshToken();
            // Retry original request
            response = await fetch(url, options);
        } catch {
            console.error("Refresh failed, user must re-login");
            if (typeof window !== 'undefined') {
                localStorage.clear();
                sessionStorage.clear();
                window.location.href = "/login";
            }
            throw new Error("Session expired. Please login again.");
        }
    }

    return response;
};

export class ApiClient {
    static async get<T>(url: string, options: RequestInit = {}): Promise<T> {
        const response = await fetchWithRefresh(url, {
            ...options,
            method: "GET",
            credentials: "include",
        });
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Failed to fetch ${url}`);
        }
        return response.json();
    }

    static async post<T>(url: string, body?: any, options: RequestInit = {}): Promise<T> {
        const response = await fetchWithRefresh(url, {
            ...options,
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            body: body ? JSON.stringify(body) : undefined,
            credentials: "include",
        });
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Failed to post to ${url}`);
        }
        return response.json();
    }

    static async put<T>(url: string, body?: any, options: RequestInit = {}): Promise<T> {
        const response = await fetchWithRefresh(url, {
            ...options,
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            body: body ? JSON.stringify(body) : undefined,
            credentials: "include",
        });
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Failed to put to ${url}`);
        }
        return response.json();
    }

    static async patch<T>(url: string, body?: any, options: RequestInit = {}): Promise<T> {
        const response = await fetchWithRefresh(url, {
            ...options,
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                ...options.headers,
            },
            body: body ? JSON.stringify(body) : undefined,
            credentials: "include",
        });
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Failed to patch ${url}`);
        }
        return response.json();
    }

    static async delete<T>(url: string, options: RequestInit = {}): Promise<T> {
        const response = await fetchWithRefresh(url, {
            ...options,
            method: "DELETE",
            credentials: "include",
        });
        if (!response.ok) {
            const error = await response.json().catch(() => ({}));
            throw new Error(error.error || `Failed to delete ${url}`);
        }
        return response.json();
    }
}
