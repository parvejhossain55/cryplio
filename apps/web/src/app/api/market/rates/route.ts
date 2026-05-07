import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
    try {
        const { searchParams } = new URL(request.url);
        const crypto = searchParams.get("crypto");
        const fiat = searchParams.get("fiat");

        let url = `${API_BASE_URL}/market/rates`;
        if (crypto || fiat) {
            const params = new URLSearchParams();
            if (crypto) params.append("crypto", crypto);
            if (fiat) params.append("fiat", fiat);
            url += `?${params.toString()}`;
        }

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to fetch market rates" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Market rates error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
