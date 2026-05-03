import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";
const AUTH_COOKIE_NAME = process.env.COOKIE_NAME || "auth_token";

export async function GET(request: NextRequest) {
    try {
        const cookieHeader = request.headers.get("cookie") || "";
        const hasAuthCookie = cookieHeader
            .split(";")
            .some((cookie) => {
                const name = cookie.trim().split("=")[0];
                return name === AUTH_COOKIE_NAME || name === `${AUTH_COOKIE_NAME}_refresh`;
            });

        if (!hasAuthCookie) {
            return NextResponse.json({ user: null });
        }

        const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
            method: "GET",
            headers: {
                cookie: cookieHeader,
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

export async function PUT(request: NextRequest) {
    try {
        const body = await request.json();
        const cookieHeader = request.headers.get("cookie") || "";

        const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                cookie: cookieHeader,
            },
            body: JSON.stringify(body),
            credentials: "include",
        });

        const data = await response.json();

        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to update user" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("Update user error:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
