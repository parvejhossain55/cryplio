"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Wallet, Mail, LogOut, Menu, X, User } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { authService } from "@/services/authService";
import NotificationBell from "./NotificationBell";

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
        <nav className="fixed top-0 w-full bg-background/80 backdrop-blur-xl border-b border-white/5 z-50">
            <div className="container mx-auto px-4 md:px-6">
                <div className="flex items-center justify-between h-20">
                    {/* Logo */}
                    <Link href="/" className="flex items-center space-x-3 group">
                        <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20 group-hover:rotate-12 transition-transform duration-500">
                            <Wallet className="w-6 h-6 text-background" />
                        </div>
                        <span className="text-2xl font-black italic uppercase tracking-tighter text-white">
                            CRYP<span className="text-primary truncate">LIO</span>
                        </span>
                    </Link>

                    {/* Desktop Navigation */}
                    <div className="hidden md:flex items-center space-x-10">
                        {["Marketplace", "Swap", "Support"].map((item) => (
                            <Link
                                key={item}
                                href={`/${item.toLowerCase()}`}
                                className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] hover:text-primary transition-all relative group"
                            >
                                {item}
                                <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-primary transition-all group-hover:w-full" />
                            </Link>
                        ))}

                        <div className="h-4 w-px bg-white/10" />

                        {user ? (
                            <div className="flex items-center space-x-6">
                                <NotificationBell />
                                <Link
                                    href="/user/dashboard"
                                    className="text-[10px] font-black text-white px-5 py-2.5 bg-white/5 border border-white/5 rounded-xl hover:bg-white/10 transition-all uppercase tracking-widest"
                                >
                                    Portal
                                </Link>
                                <button
                                    onClick={handleLogout}
                                    disabled={isLoading}
                                    className="p-2 text-red-500 hover:bg-red-500/10 rounded-xl transition-all"
                                    title="Logout"
                                >
                                    <LogOut className="w-5 h-5" />
                                </button>
                            </div>
                        ) : (
                            <div className="flex items-center space-x-6">
                                <Link href="/login" className="text-[10px] font-black text-text-dim hover:text-white uppercase tracking-widest transition-colors">
                                    Login
                                </Link>
                                <Link href="/register" className="bg-white text-background px-8 py-3 rounded-xl text-[10px] font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-xl shadow-white/5">
                                    Initialize
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
