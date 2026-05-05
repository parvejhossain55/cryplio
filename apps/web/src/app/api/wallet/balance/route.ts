import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
    try {
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/wallet/balance`, {
            method: "GET",
            headers: {
                cookie: cookieHeader,
            },
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to fetch wallet balance" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Wallet balance error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
