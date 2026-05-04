"use client";

import React from "react";
import MarketOverview from "@/components/sections/MarketOverview";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import { motion } from "framer-motion";
import Link from "next/link";
import { Search, Shield, Zap, Globe, ArrowRight, Terminal, BarChart3 } from "lucide-react";

const MarketplacePage = () => {
    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Header Section */}
            <section className="pt-40 pb-20 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-primary/5 blur-[150px] rounded-full -z-10" />
                <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-10 pointer-events-none" />

                <div className="container mx-auto px-4 md:px-6">
                    <div className="grid grid-cols-1 lg:grid-cols-12 gap-16 items-end">
                        <div className="lg:col-span-8">
                            <motion.div
                                initial={{ opacity: 0, x: -30 }}
                                animate={{ opacity: 1, x: 0 }}
                                className="space-y-8"
                            >
                                <div className="flex items-center gap-3">
                                    <Terminal className="w-4 h-4 text-primary" />
                                    <span className="text-[10px] font-black uppercase tracking-widest text-text-dim">EXCHANGE_MESH_v4.0</span>
                                </div>

                                <h1 className="text-6xl md:text-9xl font-black italic uppercase tracking-tighter leading-[0.8] text-white">
                                    P2P <br />
                                    <span className="gradient-text">LIQUIDITY.</span>
                                </h1>

                                <p className="text-xl text-text-dim max-w-2xl font-medium leading-tight uppercase italic tracking-tight">
                                    Access the global clearing layer for digital assets. Connect with institutional liquidity providers and verified peers.
                                </p>
                            </motion.div>
                        </div>

                        <div className="lg:col-span-4 lg:text-right">
                            <Link
                                href="/marketplace/create"
                                className="inline-flex items-center px-12 py-5 bg-white text-background rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-105 active:scale-95 transition-all shadow-2xl shadow-white/10 group"
                            >
                                Post Order
                                <ArrowRight className="w-5 h-5 ml-4 group-hover:translate-x-1 transition-transform" />
                            </Link>
                        </div>
                    </div>

                    {/* Stats Bar */}
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-12 mt-20 pt-12 border-t border-white/5">
                        {[
                            { label: "Aggregate Volume", value: "$14.2B", icon: BarChart3 },
                            { label: "Active Nodes", value: "8,402", icon: Globe },
                            { label: "Compliance Score", value: "99.8%", icon: Shield },
                            { label: "Settle Delta", value: "< 90s", icon: Zap },
                        ].map((stat, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                transition={{ delay: 0.2 + i * 0.1 }}
                            >
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-2 flex items-center gap-2 lg:justify-start">
                                    <stat.icon className="w-3.5 h-3.5 text-primary" />
                                    {stat.label}
                                </p>
                                <p className="text-3xl font-black text-white italic tracking-tighter">{stat.value}</p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Main Marketplace Content */}
            <div className="relative z-10">
                <MarketOverview hideViewAll={true} />
            </div>

            {/* Support/Trust CTA */}
            <section className="py-32 px-4 md:px-6">
                <div className="container mx-auto">
                    <div className="glass rounded-[3rem] border-white/5 p-8 md:p-20 relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-96 h-96 bg-primary/10 blur-[150px] -z-10" />

                        <div className="flex flex-col lg:flex-row items-center justify-between gap-16">
                            <div className="max-w-2xl space-y-8">
                                <h2 className="text-5xl md:text-7xl font-black italic uppercase tracking-tighter leading-[0.9]">
                                    SECURE <br />
                                    <span className="text-white/20">CLEARANCE.</span>
                                </h2>
                                <p className="text-xl text-text-dim font-bold uppercase tracking-widest leading-loose">
                                    Every interaction is protected by our proprietary escrow execution layer. Funds are cryptographically locked until absolute consensus is reached.
                                </p>
                                <div className="flex flex-wrap gap-4">
                                    <div className="px-6 py-3 bg-white/5 border border-white/10 rounded-xl flex items-center gap-3">
                                        <Shield className="w-5 h-5 text-primary" />
                                        <span className="text-[10px] font-black uppercase tracking-widest">ISO 27001</span>
                                    </div>
                                    <div className="px-6 py-3 bg-white/5 border border-white/10 rounded-xl flex items-center gap-3">
                                        <Globe className="w-5 h-5 text-primary" />
                                        <span className="text-[10px] font-black uppercase tracking-widest">Global Escrow</span>
                                    </div>
                                </div>
                            </div>

                            <div className="w-full lg:w-[400px] aspect-square relative group">
                                <div className="absolute inset-0 bg-primary/20 rounded-full blur-[80px] group-hover:bg-primary/30 transition-all duration-700" />
                                <div className="absolute inset-0 border-2 border-primary/20 rounded-full animate-[spin_20s_linear_infinite] border-dashed" />
                                <div className="absolute inset-4 border border-white/10 rounded-full animate-[spin_15s_linear_infinite_reverse]" />

                                <div className="absolute inset-0 flex flex-col items-center justify-center text-center">
                                    <Shield className="w-20 h-20 text-white mb-4" />
                                    <p className="text-4xl font-black italic text-white">$0.00</p>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">LOSS RECORD</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <Footer />
        </main>
    );
};

export default MarketplacePage;
