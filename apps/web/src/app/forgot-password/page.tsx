"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Mail, Loader2, ArrowRight, ArrowLeft, KeyRound } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";

const ForgotPasswordPage = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [isSubmitted, setIsSubmitted] = useState(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        // Simulate password reset request
        setTimeout(() => {
            setIsLoading(false);
            setIsSubmitted(true);
        }, 2000);
    };

    if (isSubmitted) {
        return (
            <AuthLayout
                title="Check Email"
                subtitle="We've sent password reset instructions to your email."
            >
                <div className="py-10">
                    <div className="w-24 h-24 bg-primary/10 rounded-3xl flex items-center justify-center mb-10 shadow-2xl shadow-primary/5">
                        <KeyRound className="text-primary w-12 h-12" />
                    </div>
                    <p className="text-text-dim text-lg leading-relaxed mb-10">
                        Please follow the link in your email to reset your password. If you didn't receive it, check your spam folder.
                    </p>
                    <div className="space-y-4">
                        <button
                            type="button"
                            className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl active:scale-[0.98]"
                            onClick={() => setIsSubmitted(false)}
                        >
                            Resend Email
                        </button>
                        <Link
                            href="/login"
                            className="block w-full text-center py-4 text-sm font-bold text-text-dim hover:text-white transition-colors flex items-center justify-center space-x-2"
                        >
                            <ArrowLeft className="w-4 h-4" />
                            <span>Back to Login</span>
                        </Link>
                    </div>
                </div>
            </AuthLayout>
        );
    }

    return (
        <AuthLayout
            title="Reset Password"
            subtitle="Enter your email to receive recovery instructions."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Email Address</label>
                    <div className="relative group">
                        <Mail className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="email"
                            required
                            placeholder="Enter your registered email"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/30"
                        />
                    </div>
                </div>

                <div className="pt-4">
                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl shadow-white/5 flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70 group"
                    >
                        {isLoading ? <Loader2 className="w-6 h-6 animate-spin" /> : (
                            <>
                                <span>Send Instructions</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </>
                        )}
                    </button>
                </div>

                <div className="pt-6 text-center">
                    <p className="text-xs font-medium text-text-dim">
                        Remembered your password? <Link href="/login" className="text-white font-black hover:underline underline-offset-4">Login</Link>
                    </p>
                </div>
            </form>
        </AuthLayout>
    );
};

export default ForgotPasswordPage;
