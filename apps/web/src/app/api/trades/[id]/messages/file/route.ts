import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";

export async function POST(
    request: NextRequest,
    { params }: { params: { id: string } }
) {
    try {
        const formData = await request.formData();
        const file = formData.get('file') as File;
        
        if (!file) {
            return NextResponse.json(
                { error: "No file provided" },
                { status: 400 }
            );
        }

        // Forward the file to the backend
        const backendFormData = new FormData();
        backendFormData.append('file', file);

        const cookieHeader = request.headers.get("cookie") || "";
        const response = await fetch(`${API_BASE_URL}/trades/${params.id}/messages/file`, {
            method: "POST",
            headers: {
                cookie: cookieHeader,
            },
            body: backendFormData,
            credentials: "include",
        });

        const data = await response.json();
        if (!response.ok) {
            return NextResponse.json(
                { error: data.error || "Failed to upload file" },
                { status: response.status }
            );
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error("File upload error:", error);
        return NextResponse.json({ error: "Internal server error" }, { status: 500 });
    }
}
