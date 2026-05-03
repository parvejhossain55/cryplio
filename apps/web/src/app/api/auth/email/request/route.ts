import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const { user_id } = body;

        if (!user_id) {
            return NextResponse.json(
                { error: "User ID is required" },
                { status: 400 }
            );
        }
        
        // Forward cookies from incoming request to backend
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/auth/email/request`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                cookie: cookieHeader,
            },
            body: JSON.stringify({ user_id }),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Verification request failed" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Email verification request error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
