"use client";

import React, { Suspense, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";
import { useAuth } from "@/context/AuthContext";
import { authService } from "@/services/authService";

const VerifyEmailContent = () => {
    const searchParams = useSearchParams();
    const token = searchParams.get("token");
    const { refreshUser } = useAuth();
    const hasStartedVerification = useRef(false);
    const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
    const [message, setMessage] = useState("Verifying your email...");

    useEffect(() => {
        const verifyEmail = async () => {
            if (hasStartedVerification.current) {
                return;
            }
            hasStartedVerification.current = true;

            if (!token) {
                setStatus("error");
                setMessage("Invalid or missing verification token.");
                return;
            }

            try {
                await authService.verifyEmail(token);
                await refreshUser();
                setStatus("success");
                setMessage("Your email has been verified successfully.");
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : "Email verification failed.";
                if (errorMessage.toLowerCase().includes("already verified")) {
                    await refreshUser();
                    setStatus("success");
                    setMessage("Your email is already verified.");
                    return;
                }

                setStatus("error");
                setMessage(errorMessage);
            }
        };

        void verifyEmail();
    }, [refreshUser, token]);

    return (
        <AuthLayout
            title={status === "success" ? "Email Verified" : "Verify Email"}
            subtitle={message}
        >
            <div className="py-10 text-center">
                <div className="w-20 h-20 rounded-full flex items-center justify-center mb-6 mx-auto bg-white/5">
                    {status === "loading" && <Loader2 className="w-10 h-10 text-primary animate-spin" />}
                    {status === "success" && <CheckCircle className="w-10 h-10 text-green-500" />}
                    {status === "error" && <AlertCircle className="w-10 h-10 text-red-500" />}
                </div>

                <Link
                    href={status === "success" ? "/user/dashboard" : "/login"}
                    className="inline-flex bg-white text-background px-10 py-4 rounded-2xl font-black text-lg transition-all hover:scale-105 active:scale-95 shadow-xl"
                >
                    {status === "success" ? "Go to Dashboard" : "Go to Login"}
                </Link>
            </div>
        </AuthLayout>
    );
};

const VerifyEmailPage = () => {
    return (
        <Suspense
            fallback={
                <AuthLayout title="Verify Email" subtitle="Verifying your email...">
                    <div className="py-10 text-center">
                        <Loader2 className="w-10 h-10 text-primary animate-spin mx-auto" />
                    </div>
                </AuthLayout>
            }
        >
            <VerifyEmailContent />
        </Suspense>
    );
};

export default VerifyEmailPage;
