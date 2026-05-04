import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export const runtime = 'nodejs';

export async function GET(
    request: NextRequest,
    { params }: { params: { userId: string } }
) {
    try {
        const cookieHeader = request.headers.get("cookie") || "";
        console.log("[Avatar Proxy] Request for user:", params.userId, "cookies:", !!cookieHeader);

        // Fetch current user from backend using forwarded cookies
        const backendResponse = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
            headers: {
                cookie: cookieHeader,
            },
        });

        if (!backendResponse.ok) {
            console.log("[Avatar Proxy] Backend auth failed:", backendResponse.status);
            return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
        }

        const data = await backendResponse.json();
        const currentUser = data.user;

        if (!currentUser) {
            return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
        }

        // Verify the requested userId matches the authenticated user
        if (currentUser.id !== params.userId) {
            console.log("[Avatar Proxy] User ID mismatch:", currentUser.id, "vs", params.userId);
            return NextResponse.json({ error: "Forbidden" }, { status: 403 });
        }

        const avatarUrl = currentUser.avatar_url;
        if (!avatarUrl) {
            console.log("[Avatar Proxy] No avatar for user:", params.userId);
            return NextResponse.json({ error: "No avatar found" }, { status: 404 });
        }

        console.log("[Avatar Proxy] Fetching image from:", avatarUrl);

        // Fetch the image from storage (MinIO/S3)
        const imageResponse = await fetch(avatarUrl);
        if (!imageResponse.ok) {
            console.error("[Avatar Proxy] Storage fetch error:", imageResponse.status);
            return NextResponse.json({ error: "Failed to fetch avatar from storage" }, { status: 502 });
        }

        const contentType = imageResponse.headers.get("content-type") || "image/jpeg";
        const imageBuffer = await imageResponse.arrayBuffer();

        // Return image with caching and CORS headers
        return new Response(imageBuffer, {
            status: 200,
            headers: {
                "Content-Type": contentType,
                "Cache-Control": "public, max-age=86400", // 24 hours
                "Access-Control-Allow-Origin": "*",
            },
        });
    } catch (error) {
        console.error("[Avatar Proxy] Exception:", error);
        return NextResponse.json(
            { error: "Internal server error" },
            { status: 500 }
        );
    }
}
