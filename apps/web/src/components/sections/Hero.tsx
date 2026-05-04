"use client";

import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, ShieldCheck, Zap, Globe, Activity, Terminal } from "lucide-react";
import Link from "next/link";

const Hero = () => {
    return (
        <section className="relative pt-32 pb-24 lg:pt-48 lg:pb-40 overflow-hidden">
            {/* Massive Industrial Background Text */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-[20vw] font-black text-white/[0.02] uppercase tracking-tighter select-none pointer-events-none italic">
                PROTOCOL
            </div>

            {/* Grid Pattern Overlay */}
            <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20 pointer-events-none" />
            <div className="absolute inset-0 bg-radial-[circle_at_50%_50%] from-primary/5 via-transparent to-transparent pointer-events-none" />

            <div className="container mx-auto px-4 md:px-6 relative z-10">
                <div className="max-w-5xl">
                    {/* Status Badge */}
                    <motion.div
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        className="inline-flex items-center space-x-3 bg-surface border border-white/5 pl-2 pr-4 py-1.5 rounded-full mb-8"
                    >
                        <div className="bg-primary/20 p-1.5 rounded-full">
                            <Activity className="w-3.5 h-3.5 text-primary animate-pulse" />
                        </div>
                        <span className="text-[10px] font-black uppercase tracking-widest text-text-dim">
                            <span className="text-white">SYSTEM ONLINE:</span> 12,402 Active Nodes
                        </span>
                    </motion.div>

                    {/* Main Headline */}
                    <motion.h1
                        initial={{ opacity: 0, y: 30 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
                        className="text-6xl md:text-8xl lg:text-9xl font-black italic uppercase tracking-tighter leading-[0.8] mb-10"
                    >
                        TRADE WITHOUT <br />
                        <span className="gradient-text">BOUNDARIES.</span>
                    </motion.h1>

                    <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-end">
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: 0.2 }}
                            className="lg:col-span-12"
                        >
                            <p className="text-xl md:text-2xl text-text-dim max-w-3xl leading-tight font-medium mb-12">
                                Cryplio is an institutional-grade P2P clearing layer. Secure your assets with military-grade escrow and settle trades across 50+ fiat jurisdictions instantly.
                            </p>

                            <div className="flex flex-wrap gap-4">
                                <Link
                                    href="/marketplace"
                                    className="px-12 py-5 bg-white text-background rounded-2xl font-black uppercase tracking-widest text-sm hover:scale-105 active:scale-95 transition-all shadow-2xl shadow-white/10 group flex items-center gap-3"
                                >
                                    Access Exchange
                                    <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                                </Link>
                                <Link
                                    href="/register"
                                    className="px-12 py-5 bg-surface border border-white/10 text-white rounded-2xl font-black uppercase tracking-widest text-sm hover:bg-white/5 transition-all flex items-center gap-3"
                                >
                                    Initialize Account
                                    <Terminal className="w-5 h-5 text-text-dim" />
                                </Link>
                            </div>
                        </motion.div>
                    </div>

                    {/* Industrial Stats Bar */}
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-12 mt-24 pt-12 border-t border-white/5">
                        {[
                            { label: "Network Volume", value: "$4.20B+", color: "text-primary" },
                            { label: "Settle Time", value: "< 120s", color: "text-white" },
                            { label: "Global Reach", value: "142 COUNTRIES", color: "text-white" },
                            { label: "Security Protocol", value: "AES-256", color: "text-primary" },
                        ].map((stat, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                transition={{ delay: 0.4 + i * 0.1 }}
                                className="space-y-1"
                            >
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{stat.label}</p>
                                <p className={`text-2xl font-black ${stat.color} italic uppercase`}>{stat.value}</p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Floating Industrial Elements */}
            <div className="absolute right-0 bottom-0 top-0 w-1/3 overflow-hidden pointer-events-none hidden lg:block">
                <div className="absolute top-1/4 right-0 w-[500px] h-[500px] border border-white/[0.03] rounded-full rotate-45 transform translate-x-1/2" />
                <div className="absolute bottom-1/4 right-0 w-[800px] h-[800px] border border-white/[0.02] rounded-full -rotate-12 transform translate-x-1/3" />
            </div>
        </section>
    );
};

export default Hero;
