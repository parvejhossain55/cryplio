import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        // Forward cookies from incoming request to backend
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
            method: "POST",
            headers: {
                cookie: cookieHeader,
            },
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Token refresh failed" },
                { status: response.status }
            );
        }

        // Forward new cookies
        const forwardCookies = response.headers.getSetCookie();
        const headers = new Headers();
        forwardCookies.forEach((cookie) => {
            headers.append("Set-Cookie", cookie);
        });

        return NextResponse.json(data, { headers });
    } catch (error) {
        console.error("Token refresh error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
