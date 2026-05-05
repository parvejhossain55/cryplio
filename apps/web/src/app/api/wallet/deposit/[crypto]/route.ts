import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

type Params = { params: Promise<{ crypto: string }> };

export async function GET(request: NextRequest, { params }: Params) {
    try {
        const { crypto } = await params;
        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/api/v1/wallet/deposit/${crypto}`, {
            method: "GET",
            headers: {
                cookie: cookieHeader,
            },
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to fetch deposit address" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Deposit address error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
