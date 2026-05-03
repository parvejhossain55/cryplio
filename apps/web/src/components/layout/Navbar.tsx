"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Wallet, Mail, LogOut, Menu, X, User } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { authService } from "@/services/authService";

const Navbar = () => {
    const { user, logout, isLoading } = useAuth();
    const router = useRouter();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const [isResending, setIsResending] = useState(false);

    const handleLogout = async () => {
        await logout();
        router.push("/login");
    };

    const handleResendVerification = async () => {
        if (!user) {
            return;
        }

        setIsResending(true);
        try {
            await authService.requestEmailVerification(user.id);
        } catch (error) {
            console.error("Failed to resend:", error);
        } finally {
            setIsResending(false);
        }
    };

    return (
        <nav className="fixed top-0 w-full bg-surface/80 backdrop-blur-xl border-b border-white/5 z-50">
            <div className="container mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    <Link href="/" className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center">
                            <Wallet className="w-6 h-6 text-white" />
                        </div>
                        <span className="text-2xl font-black tracking-tight">
                            Cryp<span className="gradient-text">lio</span>
                        </span>
                    </Link>

                    {/* Desktop Navigation */}
                    <div className="hidden md:flex items-center space-x-8">
                        <Link href="/marketplace" className="text-sm font-bold text-text-dim hover:text-white transition-colors">
                            Marketplace
                        </Link>
                        <Link href="/swap" className="text-sm font-bold text-text-dim hover:text-white transition-colors">
                            Swap
                        </Link>
                        <Link href="/support" className="text-sm font-bold text-text-dim hover:text-white transition-colors">
                            Support
                        </Link>

                        {user ? (
                            <>
                                {/* Email verification warning */}
                                {!user.emailVerified && (
                                    <div className="flex items-center space-x-2 bg-yellow-500/10 text-yellow-500 text-xs px-3 py-1.5 rounded-full">
                                        <Mail className="w-3.5 h-3.5" />
                                        <span>Verify email</span>
                                    </div>
                                )}

                                <div className="flex items-center space-x-4">
                                    <Link href="/user/dashboard" className="text-sm font-bold text-text-dim hover:text-white transition-colors">
                                        Dashboard
                                    </Link>
                                    <div className="relative group">
                                        <button className="flex items-center space-x-2 text-sm font-bold text-text-dim hover:text-white transition-colors">
                                            <User className="w-4 h-4" />
                                            <span>{user.username}</span>
                                        </button>
                                        {/* Dropdown menu could be added here */}
                                    </div>
                                    <button
                                        onClick={handleLogout}
                                        disabled={isLoading}
                                        className="flex items-center space-x-2 text-sm font-bold text-red-400 hover:text-red-300 transition-colors"
                                    >
                                        <LogOut className="w-4 h-4" />
                                        <span>Logout</span>
                                    </button>
                                </div>
                            </>
                        ) : (
                            <div className="flex items-center space-x-4">
                                <Link href="/login" className="text-sm font-bold text-white hover:text-primary transition-colors">
                                    Login
                                </Link>
                                <Link href="/register" className="bg-white text-background px-6 py-2.5 rounded-xl text-sm font-black hover:scale-105 transition-transform">
                                    Sign Up
                                </Link>
                            </div>
                        )}
                    </div>

                    {/* Mobile menu button */}
                    <button
                        onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                        className="md:hidden p-2 text-text-dim hover:text-white"
                    >
                        {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                    </button>
                </div>
            </div>

            {/* Mobile Menu */}
            <AnimatePresence>
                {isMobileMenuOpen && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: "auto" }}
                        exit={{ opacity: 0, height: 0 }}
                        className="md:hidden bg-surface border-t border-white/5 overflow-hidden"
                    >
                        <div className="px-4 py-6 space-y-4">
                            <Link href="/marketplace" className="block text-sm font-bold text-text-dim hover:text-white">
                                Marketplace
                            </Link>
                            <Link href="/swap" className="block text-sm font-bold text-text-dim hover:text-white">
                                Swap
                            </Link>
                            <Link href="/support" className="block text-sm font-bold text-text-dim hover:text-white">
                                Support
                            </Link>

                            {user ? (
                                <>
                                    <div className="pt-4 border-t border-white/5">
                                        <p className="text-xs font-black text-text-dim uppercase tracking-widest mb-3">
                                            Account
                                        </p>
                                        <Link href="/user/dashboard" className="block py-2 text-sm font-medium">
                                            Dashboard
                                        </Link>

                                        {/* Email verification prompt */}
                                        {!user.emailVerified && (
                                            <button
                                                onClick={handleResendVerification}
                                                disabled={isResending}
                                                className="flex items-center space-x-2 py-2 text-sm text-yellow-500"
                                            >
                                                <Mail className="w-4 h-4" />
                                                <span>{isResending ? "Sending..." : "Resend verification email"}</span>
                                            </button>
                                        )}

                                        <button
                                            onClick={handleLogout}
                                            className="flex items-center space-x-2 py-2 text-sm text-red-400"
                                        >
                                            <LogOut className="w-4 h-4" />
                                            <span>Logout</span>
                                        </button>
                                    </div>
                                </>
                            ) : (
                                <div className="pt-4 border-t border-white/5 space-y-3">
                                    <Link href="/login" className="block w-full text-center py-3 border border-white/10 rounded-xl text-sm font-bold">
                                        Login
                                    </Link>
                                    <Link href="/register" className="block w-full text-center bg-white text-background py-3 rounded-xl text-sm font-black">
                                        Sign Up
                                    </Link>
                                </div>
                            )}
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </nav>
    );
};

export default Navbar;
