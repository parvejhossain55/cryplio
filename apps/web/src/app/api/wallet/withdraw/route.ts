import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/wallet/withdraw`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                cookie: cookieHeader,
            },
            body: JSON.stringify(body),
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to request withdrawal" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Wallet withdrawal error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
