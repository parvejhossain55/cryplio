"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Mail, Lock, User, Eye, EyeOff, Loader2, ArrowRight } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";
import { useAuth } from "@/context/AuthContext";

const RegisterPage = () => {
    const { register } = useAuth();
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            const formData = new FormData(e.currentTarget as HTMLFormElement);
            const fullName = formData.get("username") as string;
            const email = formData.get("email") as string;
            const password = formData.get("password") as string;

            await register(email, fullName, password);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Registration failed");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <AuthLayout
            title="Register"
            subtitle="Create an account to start trading crypto."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                {error && (
                    <div className="p-4 bg-red-500/10 border border-red-500/50 rounded-xl text-red-500 text-sm">
                        {error}
                    </div>
                )}
                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Full Name</label>
                    <div className="relative group">
                        <User className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="text"
                            name="username"
                            required
                            placeholder="Enter your full name"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium"
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Email Address</label>
                    <div className="relative group">
                        <Mail className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="email"
                            name="email"
                            required
                            placeholder="Enter your email"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium"
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Password</label>
                    <div className="relative group">
                        <Lock className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                                type={showPassword ? "text" : "password"}
                                name="password"
                                required
                                placeholder="Create a password"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-12 text-base outline-none focus:border-primary transition-all font-medium"
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

                <div className="space-y-4 pt-4">
                    <div className="flex items-start space-x-3 group">
                        <div className="relative flex items-center h-5">
                            <input
                                type="checkbox"
                                id="terms"
                                required
                                className="w-5 h-5 rounded-lg border-border bg-surface text-primary focus:ring-primary/20 accent-primary cursor-pointer"
                            />
                        </div>
                        <label htmlFor="terms" className="text-xs text-text-dim leading-relaxed cursor-pointer select-none group-hover:text-white transition-colors">
                            I agree to the <Link href="/terms" className="text-white font-black hover:underline underline-offset-4">Terms and Conditions</Link>.
                        </label>
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
                                <span>Register</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </>
                        )}
                    </button>
                </div>

                <div className="pt-6 text-center border-t border-white/5">
                    <p className="text-xs font-medium text-text-dim">
                        Already have an account? <Link href="/login" className="text-white font-black hover:underline underline-offset-4">Login</Link>
                    </p>
                </div>
            </form>
        </AuthLayout>
    );
};

export default RegisterPage;
