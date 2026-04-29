import { ReactNode } from "react";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const API_BASE_URL = process.env.API_GATEWAY_URL || "http://localhost:8080";
const AUTH_COOKIE_NAME = "auth_token";

async function validateDashboardSession() {
    const cookieStore = await cookies();
    const authToken = cookieStore.get(AUTH_COOKIE_NAME)?.value;

    if (!authToken) {
        redirect("/login");
    }

    const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
        method: "GET",
        headers: {
            cookie: `${AUTH_COOKIE_NAME}=${authToken}`,
        },
        cache: "no-store",
    });

    if (!response.ok) {
        redirect("/login");
    }
}

export default async function DashboardRouteLayout({
    children,
}: {
    children: ReactNode;
}) {
    await validateDashboardSession();

    return <>{children}</>;
}
