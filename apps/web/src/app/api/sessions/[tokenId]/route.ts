import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function DELETE(
    request: NextRequest,
    context: { params: Promise<{ tokenId: string }> }
) {
    try {
        const { tokenId } = await context.params;

        if (!tokenId) {
            return NextResponse.json(
                { error: "Token ID is required" },
                { status: 400 }
            );
        }

        const response = await fetch(`${API_BASE_URL}/api/v1/sessions/${tokenId}`, {
            method: "DELETE",
            headers: {
                cookie: request.headers.get("cookie") || "",
            },
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to revoke session" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Revoke session error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
