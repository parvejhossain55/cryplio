import { NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function GET() {
    return NextResponse.redirect(`${API_BASE_URL}/api/v1/auth/oauth/google`);
}
