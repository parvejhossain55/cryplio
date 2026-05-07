import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(
    request: NextRequest,
    { params }: { params: { id: string } }
) {
    try {
        const body = await request.json();
        const { tx_hash } = body;

        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/admin/withdrawals/${params.id}/approve`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                cookie: cookieHeader,
            },
            body: JSON.stringify({ tx_hash }),
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to approve withdrawal" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Approve withdrawal error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
