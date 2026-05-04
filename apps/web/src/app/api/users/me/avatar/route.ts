import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";
const AUTH_COOKIE_NAME = process.env.COOKIE_NAME || "auth_token";

export const runtime = 'nodejs';

// Increase body size limit for file uploads (max 2MB)
export const config = {
    api: {
        bodyParser: {
            sizeLimit: '5mb',
        },
    },
} as const;

export async function POST(request: NextRequest) {
    try {
        const cookieHeader = request.headers.get("cookie") || "";
        console.log("[Avatar Upload] Received request, cookies present:", !!cookieHeader);

        const formData = await request.formData();
        const file = formData.get("avatar") as File | null;
        if (!file) {
            console.warn("[Avatar Upload] No file in form data");
            return NextResponse.json(
                { error: "Avatar file is required" },
                { status: 400 }
            );
        }

        console.log("[Avatar Upload] File received:", file.name, "size:", file.size);

        const response = await fetch(`${API_BASE_URL}/api/v1/users/me/avatar`, {
            method: "POST",
            headers: {
                cookie: cookieHeader,
            },
            body: formData,
        });

        console.log("[Avatar Upload] Backend response status:", response.status);

        const data = await response.json();

        if (!response.ok) {
            console.error("[Avatar Upload] Backend error:", response.status, data);
            return NextResponse.json(
                { error: data.error || "Failed to upload avatar" },
                { status: response.status }
            );
        }

        console.log("[Avatar Upload] Success");
        return NextResponse.json(data);
    } catch (error) {
        console.error("[Avatar Upload] Exception:", error);
        return NextResponse.json(
            { error: "Internal server error", details: String(error) },
            { status: 500 }
        );
    }
}
