"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Lock, Eye, EyeOff, Loader2, ArrowRight, CheckCircle } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";
import { authService } from "@/services/authService";

const ResetPasswordPage = () => {
    const router = useRouter();
    // Get token from URL using useSearchParams in a separate component
    const [token, setToken] = useState<string | null>(null);
    
    useEffect(() => {
        // This is a workaround for Next.js navigation hooks in client components
        // In practice, we'd access search params differently or handle this in layout
        const urlParams = new URLSearchParams(window.location.search);
        setToken(urlParams.get("token"));
    }, []);

    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isSuccess, setIsSuccess] = useState(false);

    useEffect(() => {
        if (!token) {
            setError("Invalid or missing reset token");
        }
    }, [token]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        if (!token) {
            setError("Invalid reset token");
            return;
        }

        if (password !== confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        if (password.length < 8) {
            setError("Password must be at least 8 characters");
            return;
        }

        setIsLoading(true);
        try {
            await authService.resetPassword(token, password);
            setIsSuccess(true);
        } catch (err: any) {
            setError(err.message || "Failed to reset password");
        } finally {
            setIsLoading(false);
        }
    };

    if (isSuccess) {
        return (
            <AuthLayout
                title="Password Reset Successful"
                subtitle="Your password has been changed."
            >
                <div className="py-10 text-center">
                    <div className="w-20 h-20 bg-green-500/10 rounded-full flex items-center justify-center mb-6 mx-auto">
                        <CheckCircle className="w-10 h-10 text-green-500" />
                    </div>
                    <p className="text-text-dim text-lg leading-relaxed mb-8">
                        You can now log in with your new password.
                    </p>
                    <Link
                        href="/login"
                        className="inline-flex bg-white text-background px-10 py-4 rounded-2xl font-black text-lg transition-all hover:scale-105 active:scale-95 shadow-xl"
                    >
                        Go to Login
                    </Link>
                </div>
            </AuthLayout>
        );
    }

    return (
        <AuthLayout
            title="Reset Password"
            subtitle="Enter your new password."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                {error && (
                    <div className="p-4 bg-red-500/10 border border-red-500/50 rounded-xl text-red-500 text-sm">
                        {error}
                    </div>
                )}

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">
                        New Password
                    </label>
                    <div className="relative group">
                        <Lock className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type={showPassword ? "text" : "password"}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            minLength={8}
                            placeholder="Enter new password"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-10 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/30"
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute right-0 top-1/2 -translate-y-1/2 text-text-dim hover:text-white transition-colors"
                        >
                            {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                        </button>
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">
                        Confirm Password
                    </label>
                    <div className="relative group">
                        <Lock className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type={showPassword ? "text" : "password"}
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                            minLength={8}
                            placeholder="Confirm new password"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/30"
                        />
                    </div>
                </div>

                <div className="pt-4">
                    <button
                        type="submit"
                        disabled={isLoading || !password || !confirmPassword}
                        className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl shadow-white/5 flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70 group"
                    >
                        {isLoading ? (
                            <Loader2 className="w-6 h-6 animate-spin" />
                        ) : (
                            <>
                                <span>Reset Password</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </>
                        )}
                    </button>
                </div>
            </form>
        </AuthLayout>
    );
};

export default ResetPasswordPage;
