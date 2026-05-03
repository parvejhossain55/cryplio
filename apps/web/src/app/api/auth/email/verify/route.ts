import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const { token } = body;

        if (!token) {
            return NextResponse.json(
                { error: "Token is required" },
                { status: 400 }
            );
        }
        
        // Forward cookies from incoming request to backend
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/auth/email/verify`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                cookie: cookieHeader,
            },
            body: JSON.stringify({ token }),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Email verification failed" },
                { status: response.status }
            );
        }

        // Forward auth cookies
        const forwardCookies = response.headers.getSetCookie();
        const headers = new Headers();
        forwardCookies.forEach((cookie) => {
            headers.append("Set-Cookie", cookie);
        });

        return NextResponse.json(data, { headers });
    } catch (error) {
        console.error("Email verify error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
