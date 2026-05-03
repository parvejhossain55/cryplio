import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const { code } = body;

        if (!code) {
            return NextResponse.json(
                { error: "Verification code is required" },
                { status: 400 }
            );
        }

        const response = await fetch(`${API_BASE_URL}/api/v1/auth/2fa/verify`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                cookie: request.headers.get("cookie") || "",
            },
            body: JSON.stringify({ code }),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "2FA verification failed" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("2FA verify error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
