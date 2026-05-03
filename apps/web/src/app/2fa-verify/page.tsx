"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Shield, Loader2, ArrowLeft, Smartphone } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";
import { useAuth } from "@/context/AuthContext";

const TwoFactorVerifyPage = () => {
    const { complete2FALogin, isLoading } = useAuth();
    const [code, setCode] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [countdown, setCountdown] = useState(5);

    const router = useRouter();

    useEffect(() => {
        // Check if we have a temp token in sessionStorage
        const tempToken = sessionStorage.getItem("2fa_temp_token");
        if (!tempToken) {
            router.push("/login");
        }

        // Auto-redirect after 5 minutes if code not entered
        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer);
                    sessionStorage.removeItem("2fa_temp_token");
                    sessionStorage.removeItem("2fa_user_id");
                    router.push("/login");
                    return 0;
                }
                return prev - 1;
            });
        }, 60000); // Update every minute

        return () => clearInterval(timer);
    }, [router]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        if (code.length !== 6 || !/^\d+$/.test(code)) {
            setError("Please enter a valid 6-digit verification code");
            return;
        }

        try {
            await complete2FALogin(code);
            // AuthContext will handle redirect
        } catch (err: any) {
            setError(err.message || "Verification failed. Please try again.");
        }
    };

    return (
        <AuthLayout
            title="Two-Factor Verification"
            subtitle="Enter the 6-digit code from your authenticator app."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                {error && (
                    <div className="p-4 bg-red-500/10 border border-red-500/50 rounded-xl text-red-500 text-sm">
                        {error}
                    </div>
                )}

                <div className="bg-surface/50 border border-white/5 rounded-3xl p-8 text-center space-y-6">
                    <div className="w-20 h-20 mx-auto bg-primary/10 rounded-full flex items-center justify-center">
                        <Smartphone className="w-10 h-10 text-primary" />
                    </div>
                    <div>
                        <h3 className="text-lg font-black text-white mb-2">Verify Your Identity</h3>
                        <p className="text-sm text-text-dim">
                            Open your authenticator app and enter the 6-digit code.
                        </p>
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">
                        Verification Code
                    </label>
                    <input
                        type="text"
                        inputMode="numeric"
                        autoComplete="one-time-code"
                        maxLength={6}
                        placeholder="000000"
                        value={code}
                        onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                        className="w-full bg-transparent border-b border-border py-4 pl-4 text-3xl text-center font-mono tracking-widest outline-none focus:border-primary transition-all placeholder:text-text-dim/30"
                        autoFocus
                    />
                </div>

                <div className="pt-2">
                    <button
                        type="submit"
                        disabled={isLoading || code.length !== 6}
                        className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl shadow-white/5 flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-50 group"
                    >
                        {isLoading ? (
                            <Loader2 className="w-6 h-6 animate-spin" />
                        ) : (
                            <>
                                <span>Verify</span>
                                <Shield className="w-5 h-5 group-hover:scale-110 transition-transform" />
                            </>
                        )}
                    </button>
                </div>

                <div className="pt-6 border-t border-white/5">
                    <p className="text-xs font-medium text-text-dim text-center">
                        Having trouble?{" "}
                        <button
                            type="button"
                            onClick={() => {
                                if (confirm("Cancel 2FA verification and return to login?")) {
                                    sessionStorage.removeItem("2fa_temp_token");
                                    sessionStorage.removeItem("2fa_user_id");
                                    router.push("/login");
                                }
                            }}
                            className="text-primary hover:underline"
                        >
                            Cancel
                        </button>
                    </p>
                </div>

                {/* Countdown indicator */}
                <div className="text-center">
                    <p className="text-[10px] text-text-dim">
                        This verification session expires in <span className="text-primary font-bold">{countdown}</span> minutes
                    </p>
                </div>
            </form>
        </AuthLayout>
    );
};

export default TwoFactorVerifyPage;
