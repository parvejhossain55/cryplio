"use client";

import React from "react";
import MarketOverview from "@/components/sections/MarketOverview";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import { motion } from "framer-motion";
import { Search, Shield, Zap, Globe } from "lucide-react";

const MarketplacePage = () => {
    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Header Section */}
            <section className="pt-32 pb-12 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[500px] bg-primary/5 blur-[120px] rounded-full -z-10" />

                <div className="container mx-auto px-4 md:px-6">
                    <div className="w-full text-left">
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            className="space-y-6"
                        >
                            <div className="inline-flex items-center space-x-2 bg-surface border border-border px-4 py-2 rounded-2xl">
                                <Globe className="w-4 h-4 text-primary animate-pulse" />
                                <span className="text-xs font-black uppercase tracking-widest text-text-dim">Global Trading Network</span>
                            </div>
                            <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-[0.9]">
                                Peer-to-Peer <br />
                                <span className="gradient-text">Exchange Hub</span>
                            </h1>
                            <p className="text-xl text-text-dim max-w-2xl leading-relaxed font-medium">
                                Browse thousands of secure offers from verified merchants worldwide. Use our advanced filters to find the perfect trade partner in seconds.
                            </p>
                        </motion.div>

                        {/* Stats Info */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mt-12 pt-12 border-t border-border">
                            {[
                                { label: "24h Volume", value: "$4.2M+", icon: Zap },
                                { label: "Active Offers", value: "12,450", icon: Search },
                                { label: "Verified Merch", value: "3,200", icon: Shield },
                                { label: "Avg. Match", value: "< 2min", icon: Globe },
                            ].map((stat, i) => (
                                <motion.div
                                    key={i}
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    transition={{ delay: 0.2 + i * 0.1 }}
                                >
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1 flex items-center gap-1">
                                        <stat.icon className="w-3 h-3" />
                                        {stat.label}
                                    </p>
                                    <p className="text-xl font-black text-white">{stat.value}</p>
                                </motion.div>
                            ))}
                        </div>
                    </div>
                </div>
            </section>

            {/* Main Marketplace Content */}
            <MarketOverview hideViewAll={true} />

            {/* Trust Section */}
            <section className="py-24 bg-surface/30">
                <div className="container mx-auto px-4 md:px-6">
                    <div className="glass rounded-[40px] border-border p-8 md:p-16 flex flex-col md:flex-row items-center justify-between gap-12">
                        <div className="space-y-6 max-w-xl">
                            <h2 className="text-3xl md:text-4xl font-black tracking-tight">Trade with 100% <br />Peace of Mind</h2>
                            <p className="text-lg text-text-dim leading-relaxed">
                                Our institutional-grade escrow protection ensures that your funds never leave your wallet until the trade is successfully completed.
                            </p>
                            <div className="flex flex-wrap gap-4">
                                <div className="flex items-center gap-2 bg-background/50 px-4 py-2 rounded-xl border border-border">
                                    <Shield className="w-5 h-5 text-accent" />
                                    <span className="font-bold text-sm">Escrow Secure</span>
                                </div>
                                <div className="flex items-center gap-2 bg-background/50 px-4 py-2 rounded-xl border border-border">
                                    <Zap className="w-5 h-5 text-primary" />
                                    <span className="font-bold text-sm">Instant Release</span>
                                </div>
                            </div>
                        </div>
                        <div className="w-full max-w-[400px] aspect-square bg-gradient-to-br from-primary/10 to-accent/10 rounded-[40px] border border-border flex items-center justify-center p-12 relative">
                            <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20" />
                            <Shield className="w-32 h-32 text-white/20" />
                            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-center">
                                <p className="text-4xl font-black text-white">100%</p>
                                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">Escrow Protected</p>
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
