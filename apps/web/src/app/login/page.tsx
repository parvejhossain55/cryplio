"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Mail, Lock, Eye, EyeOff, Loader2, ArrowRight } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";
import { useAuth } from "@/context/AuthContext";

const LoginPage = () => {
    const { login, loginWithGoogle } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            await login(email, password);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Login failed");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <AuthLayout
            title="Login"
            subtitle="Enter your email and password to login."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                {error && (
                    <div className="p-4 bg-red-500/10 border border-red-500/50 rounded-xl text-red-500 text-sm">
                        {error}
                    </div>
                )}
                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Email</label>
                    <div className="relative group">
                        <Mail className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="email"
                            required
                            placeholder="Enter your email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/30"
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <div className="flex items-center justify-between px-1">
                        <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block">Password</label>
                        <Link href="/forgot-password" title="Recover your password" className="text-xs font-bold text-primary hover:text-white transition-colors">Forgot Password?</Link>
                    </div>
                    <div className="relative group">
                        <Lock className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type={showPassword ? "text" : "password"}
                            required
                            placeholder="Enter your password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
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

                <div className="pt-4">
                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl shadow-white/5 flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70 group"
                    >
                        {isLoading ? <Loader2 className="w-6 h-6 animate-spin" /> : (
                            <>
                                <span>Login</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </>
                        )}
                    </button>
                </div>

                <div className="relative flex items-center py-4">
                    <div className="flex-grow border-t border-white/5"></div>
                    <span className="flex-shrink mx-4 text-[10px] font-black text-text-dim uppercase tracking-widest">Global SSO</span>
                    <div className="flex-grow border-t border-white/5"></div>
                </div>

                <button
                    type="button"
                    onClick={loginWithGoogle}
                    className="flex items-center justify-center space-x-3 w-full bg-surface border border-border py-4 rounded-2xl hover:bg-surface-light transition-all active:scale-[0.98] font-bold text-sm"
                >
                    <svg className="w-5 h-5" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                        <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                        <path fill="currentColor" d="M5.84 14.1c-.22-.66-.35-1.36-.35-2.1s.13-1.44.35-2.1V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.83z" />
                        <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.83c.87-2.6 3.3-4.52 6.16-4.52z" />
                    </svg>
                    <span>Continue with Google Cloud</span>
                </button>

                <div className="pt-6 text-center">
                    <p className="text-xs font-medium text-text-dim">
                        Don&apos;t have an account? <Link href="/register" className="text-white font-black hover:underline underline-offset-4">Create one</Link>
                    </p>
                </div>
            </form>
        </AuthLayout>
    );
};

export default LoginPage;
