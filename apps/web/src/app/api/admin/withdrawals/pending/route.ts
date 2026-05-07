import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
    try {
        const { searchParams } = new URL(request.url);
        const limit = searchParams.get("limit") || "50";
        const offset = searchParams.get("offset") || "0";

        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/admin/withdrawals/pending?limit=${limit}&offset=${offset}`, {
            method: "GET",
            headers: {
                cookie: cookieHeader,
            },
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to fetch pending withdrawals" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Pending withdrawals error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
