import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
    try {
        const { searchParams } = new URL(request.url);
        const limit = searchParams.get("limit") || "20";
        const offset = searchParams.get("offset") || "0";
        const status = searchParams.get("status");

        let url = `${API_BASE_URL}/admin/trades?limit=${limit}&offset=${offset}`;
        if (status) {
            url += `&status=${status}`;
        }

        const response = await fetch(url, {
            method: "GET",
            headers: {
                cookie: request.headers.get("cookie") || "",
            },
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to fetch trades" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Admin trades error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
