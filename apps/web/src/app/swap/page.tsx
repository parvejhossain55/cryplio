"use client";

import React from "react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import SwapSection from "@/components/sections/SwapSection";
import { motion } from "framer-motion";
import { Zap, Shield, Globe, ArrowRightLeft } from "lucide-react";

const SwapPage = () => {
    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Header Section */}
            <section className="pt-32 pb-4 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[500px] bg-primary/5 blur-[120px] rounded-full -z-10" />

                <div className="container mx-auto px-4 md:px-6">
                    <div className="w-full text-center">
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            className="space-y-6"
                        >
                            <div className="inline-flex items-center space-x-2 bg-surface border border-border px-4 py-2 rounded-2xl mx-auto">
                                <ArrowRightLeft className="w-4 h-4 text-primary animate-pulse" />
                                <span className="text-xs font-black uppercase tracking-widest text-text-dim">Institutional Liquidity</span>
                            </div>
                            <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-[0.9]">
                                Zero-Fee <br />
                                <span className="gradient-text">Asset Swaps</span>
                            </h1>
                            <p className="text-xl text-text-dim max-w-2xl mx-auto leading-relaxed font-medium">
                                Swap between hundreds of crypto assets instantly. Best rates guaranteed by our smart liquidity aggregator.
                            </p>
                        </motion.div>
                    </div>
                </div>
            </section>

            {/* Main Swap Component */}
            <SwapSection />

            {/* Features Stats */}
            <section className="pb-24">
                <div className="container mx-auto px-4 md:px-6">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        {[
                            {
                                icon: Zap,
                                title: "Instant Execution",
                                desc: "Your trades are processed in milliseconds with our high-frequency matching engine.",
                                color: "text-primary"
                            },
                            {
                                icon: Shield,
                                title: "Slippage Protection",
                                desc: "Advanced algorithms protect your trades from front-running and high price impact.",
                                color: "text-accent"
                            },
                            {
                                icon: Globe,
                                title: "Multi-Chain Access",
                                desc: "Access liquidity across Ethereum, BSC, Solana, and more from a single interface.",
                                color: "text-secondary"
                            }
                        ].map((item, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 20 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                transition={{ delay: i * 0.1 }}
                                className="glass p-8 rounded-[32px] border-border hover:border-primary/30 transition-all group"
                            >
                                <div className={`w-14 h-14 rounded-2xl bg-surface border border-border flex items-center justify-center mb-6 group-hover:scale-110 transition-transform`}>
                                    <item.icon className={`w-7 h-7 ${item.color}`} />
                                </div>
                                <h3 className="text-xl font-black mb-3">{item.title}</h3>
                                <p className="text-text-dim text-sm leading-relaxed font-medium">
                                    {item.desc}
                                </p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            <Footer />
        </main>
    );
};

export default SwapPage;
