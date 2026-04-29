import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET(_request: NextRequest) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
            method: "GET",
            headers: {
                cookie: _request.headers.get("cookie") || "",
            },
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to get user" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Get user error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
