"use client";

import React from "react";
import Link from "next/link";
import { Wallet, ArrowLeft, ShieldCheck, Zap, Globe } from "lucide-react";
import { motion } from "framer-motion";

interface AuthLayoutProps {
    children: React.ReactNode;
    title: string;
    subtitle: string;
}

const AuthLayout = ({ children, title, subtitle }: AuthLayoutProps) => {
    return (
        <div className="min-h-screen bg-background flex flex-col lg:flex-row overflow-hidden">
            {/* Left Side: Marketing & Visuals */}
            <div className="hidden lg:flex lg:w-[45%] bg-surface relative flex-col justify-between p-12 overflow-hidden border-r border-border">
                {/* Abstract Background Patterns */}
                <div className="absolute top-0 left-0 w-full h-full opacity-30 pointer-events-none">
                    <div className="absolute top-[-10%] left-[-10%] w-[60%] h-[60%] bg-primary/20 rounded-full blur-[120px]" />
                    <div className="absolute bottom-[-10%] right-[-10%] w-[60%] h-[60%] bg-secondary/20 rounded-full blur-[120px]" />
                    <div className="absolute top-[20%] right-[-20%] w-[40%] h-[40%] bg-accent/10 rounded-full blur-[100px]" />
                </div>

                <div className="relative z-10">
                    <Link href="/" className="flex items-center space-x-3 group">
                        <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center shadow-2xl shadow-primary/40 group-hover:rotate-12 transition-transform duration-500">
                            <Wallet className="text-white w-7 h-7" />
                        </div>
                        <span className="text-3xl font-black tracking-tighter text-white">
                            Cryp<span className="gradient-text">lio</span>
                        </span>
                    </Link>
                </div>

                <div className="relative z-10 max-w-lg mb-20">
                    <motion.div
                        initial={{ opacity: 0, x: -30 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ duration: 0.8, delay: 0.2 }}
                    >
                        <h2 className="text-5xl font-black text-white leading-tight mb-8">
                            The best way <br />
                            to trade <br />
                            crypto.
                        </h2>

                        <div className="space-y-8">
                            {[
                                { icon: ShieldCheck, title: "Safest platform", desc: "Your trades are protected by our secure escrow." },
                                { icon: Zap, title: "Very fast", desc: "Get your crypto or cash in minutes." },
                                { icon: Globe, title: "Available worldwide", desc: "Trade with people from anywhere in the world." },
                            ].map((feature, i) => (
                                <div key={i} className="flex items-start space-x-4 group">
                                    <div className="w-10 h-10 rounded-xl bg-surface-light border border-white/5 flex items-center justify-center text-primary group-hover:scale-110 transition-transform">
                                        <feature.icon className="w-5 h-5" />
                                    </div>
                                    <div>
                                        <h3 className="font-bold text-white text-lg">{feature.title}</h3>
                                        <p className="text-text-dim text-sm leading-relaxed">{feature.desc}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </motion.div>
                </div>

                <div className="relative z-10 flex items-center justify-between border-t border-white/5 pt-8">
                    <div className="flex -space-x-3">
                        {[1, 2, 3, 4].map((i) => (
                            <div key={i} className="w-10 h-10 rounded-full border-2 border-surface bg-surface-light flex items-center justify-center text-xs font-bold text-text-dim">
                                U{i}
                            </div>
                        ))}
                        <div className="w-10 h-10 rounded-full border-2 border-surface bg-primary flex items-center justify-center text-[10px] font-black text-white">
                            +1M
                        </div>
                    </div>
                    <p className="text-xs font-medium text-text-dim uppercase tracking-widest">Trusted by traders worldwide</p>
                </div>
            </div>

            {/* Right Side: Form Content */}
            <div className="flex-1 flex flex-col bg-background relative overflow-y-auto">
                <div className="lg:hidden p-6 flex justify-between items-center">
                    <Link href="/" className="flex items-center space-x-2">
                        <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                            <Wallet className="text-white w-5 h-5" />
                        </div>
                        <span className="text-xl font-bold text-white">Cryplio</span>
                    </Link>
                    <Link href="/" className="text-sm font-bold text-text-dim flex items-center space-x-1">
                        <ArrowLeft className="w-4 h-4" />
                        <span>Home</span>
                    </Link>
                </div>

                <div className="flex-1 flex items-center justify-center p-6 md:p-12 lg:p-20">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="w-full max-w-md"
                    >
                        <div className="mb-10">
                            <h1 className="text-3xl font-black text-white mb-3">{title}</h1>
                            <p className="text-text-dim font-medium">{subtitle}</p>
                        </div>

                        <div className="relative">
                            {children}
                        </div>
                    </motion.div>
                </div>

                <div className="p-8 text-center border-t border-white/5">
                    <p className="text-[10px] font-bold text-text-dim uppercase tracking-[0.2em]">
                        Cryplio Protocol • Secure SSL • 2026 Version 2.4.0
                    </p>
                </div>
            </div>
        </div>
    );
};

export default AuthLayout;
