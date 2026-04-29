import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(body),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Login failed" },
                { status: response.status }
            );
        }

        // Forward the auth cookie
        const forwardCookies = response.headers.getSetCookie();
        const headers = new Headers();
        forwardCookies.forEach((cookie) => {
            headers.append("Set-Cookie", cookie);
        });

        return NextResponse.json(data, { headers });
    } catch (error) {
        console.error("Login error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
