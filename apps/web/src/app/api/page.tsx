"use client";

import React from "react";
import { motion } from "framer-motion";
import {
    Code2,
    Terminal,
    Cpu,
    Globe,
    ShieldCheck,
    Zap,
    BookOpen,
    Copy,
    ArrowRight,
    ExternalLink
} from "lucide-react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";

const ApiDocsPage = () => {
    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Hero Section */}
            <section className="pt-32 pb-20 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-primary/5 blur-[120px] rounded-full -z-10" />

                <div className="container mx-auto px-4">
                    <div className="flex flex-col lg:flex-row items-center gap-16">
                        <motion.div
                            initial={{ opacity: 0, x: -30 }}
                            animate={{ opacity: 1, x: 0 }}
                            className="flex-1 space-y-6"
                        >
                            <div className="inline-flex items-center space-x-2 bg-surface border border-border px-4 py-2 rounded-2xl">
                                <Terminal className="w-4 h-4 text-primary" />
                                <span className="text-xs font-black uppercase tracking-widest text-text-dim">Developer Protocol v2.4</span>
                            </div>
                            <h1 className="text-5xl md:text-8xl font-black tracking-tight leading-[0.85] uppercase">
                                Build on <br />
                                <span className="gradient-text">Cryplio Hub.</span>
                            </h1>
                            <p className="max-w-xl text-text-dim text-lg md:text-xl font-medium pt-4">
                                Strategic APIs for institutional liquidity, automated P2P trading, and secure escrow management.
                            </p>
                            <div className="flex flex-wrap gap-4 pt-6">
                                <button className="bg-white text-background px-8 py-4 rounded-2xl font-black transition-all hover:scale-[1.02] active:scale-95 shadow-2xl">
                                    Get API Keys
                                </button>
                                <button className="bg-surface border border-border px-8 py-4 rounded-2xl font-bold transition-all hover:bg-surface-light">
                                    Read SDK Docs
                                </button>
                            </div>
                        </motion.div>

                        <motion.div
                            initial={{ opacity: 0, scale: 0.9 }}
                            animate={{ opacity: 1, scale: 1 }}
                            className="flex-1 w-full"
                        >
                            <div className="glass p-2 rounded-[32px] border-border overflow-hidden shadow-2xl">
                                <div className="bg-surface-light p-4 flex items-center justify-between border-b border-white/5 mx-[-8px] mt-[-8px] rounded-t-[32px]">
                                    <div className="flex space-x-2 px-2">
                                        <div className="w-3 h-3 rounded-full bg-red-500/50" />
                                        <div className="w-3 h-3 rounded-full bg-amber-500/50" />
                                        <div className="w-3 h-3 rounded-full bg-emerald-500/50" />
                                    </div>
                                    <div className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Authentication Example</div>
                                    <div className="w-12" />
                                </div>
                                <div className="p-8 font-mono text-sm leading-relaxed overflow-x-auto">
                                    <div className="flex items-start space-x-4">
                                        <span className="text-white/20 select-none">01</span>
                                        <span className="text-primary">curl</span>
                                        <span>-X POST https://api.cryplio.com/v2/auth</span>
                                    </div>
                                    <div className="flex items-start space-x-4 mt-2">
                                        <span className="text-white/20 select-none">02</span>
                                        <span className="text-secondary">-H</span>
                                        <span className="text-text-dim">"Content-Type: application/json"</span>
                                    </div>
                                    <div className="flex items-start space-x-4 mt-2">
                                        <span className="text-white/20 select-none">03</span>
                                        <span className="text-secondary">-d</span>
                                        <span className="text-text-dim">{'{'}</span>
                                    </div>
                                    <div className="flex items-start space-x-4 mt-1">
                                        <span className="text-white/20 select-none">04</span>
                                        <span className="pl-4 text-emerald-400">"api_key"</span>
                                        <span>:</span>
                                        <span className="text-amber-300">"CRYP_LIVE_8d1..."</span>
                                    </div>
                                    <div className="flex items-start space-x-4 mt-1">
                                        <span className="text-white/20 select-none">05</span>
                                        <span className="pl-4 text-emerald-400">"secret"</span>
                                        <span>:</span>
                                        <span className="text-amber-300">"••••••••••••••••"</span>
                                    </div>
                                    <div className="flex items-start space-x-4 mt-1">
                                        <span className="text-white/20 select-none">06</span>
                                        <span className="text-text-dim">{'}'}</span>
                                    </div>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </div>
            </section>

            {/* Features Grid */}
            <section className="py-24 border-t border-white/5">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        {[
                            {
                                icon: Zap,
                                title: "Real-time Webhooks",
                                desc: "Get notified instantly about trade status changes, payment confirmations, and escrow releases."
                            },
                            {
                                icon: ShieldCheck,
                                title: "Secure Escrow API",
                                desc: "Programmatically manage escrow locks and releases with institutional-grade security signatures."
                            },
                            {
                                icon: Cpu,
                                title: "High Throughput",
                                desc: "Infrastructure designed for low-latency market making and automated high-frequency P2P trading."
                            },
                            {
                                icon: Globe,
                                title: "Global Liquidity",
                                desc: "Access verified merchants and payment methods across 50+ countries through a single unified interface."
                            },
                            {
                                icon: BookOpen,
                                title: "Interactive Docs",
                                desc: "Try out API calls directly in our interactive playground with real sandboxed environments."
                            },
                            {
                                icon: Code2,
                                title: "SDKs & Tools",
                                desc: "Official libraries for Go, Python, Node.js, and Rust to jumpstart your integration process."
                            }
                        ].map((feature, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 30 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                transition={{ delay: i * 0.1 }}
                                className="glass p-10 rounded-[40px] border-border hover:border-primary/30 transition-all group"
                            >
                                <div className="w-14 h-14 bg-surface rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 transition-transform shadow-2xl">
                                    <feature.icon className="w-7 h-7 text-primary" />
                                </div>
                                <h3 className="text-2xl font-black mb-4 uppercase tracking-tight">{feature.title}</h3>
                                <p className="text-text-dim leading-relaxed font-medium">
                                    {feature.desc}
                                </p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Resource Links */}
            <section className="py-24 relative overflow-hidden bg-surface/30">
                <div className="container mx-auto px-4 text-center">
                    <h2 className="text-4xl md:text-5xl font-black mb-16 uppercase tracking-tight">Developer <span className="text-primary">Resources</span></h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                        {[
                            { title: "API Reference", icon: Terminal },
                            { title: "Postman Gallery", icon: Globe },
                            { title: "GitHub Repo", icon: Code2 },
                            { title: "Status Page", icon: Zap },
                        ].map((item, i) => (
                            <div key={i} className="glass p-8 rounded-[32px] border-border flex flex-col items-center justify-center space-y-4 hover:border-primary/50 transition-all cursor-pointer group">
                                <item.icon className="w-8 h-8 text-text-dim group-hover:text-primary transition-colors" />
                                <span className="font-black text-sm uppercase tracking-widest">{item.title}</span>
                                <ExternalLink className="w-4 h-4 text-text-dim opacity-0 group-hover:opacity-100 transition-all" />
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            <Footer />
        </main>
    );
};

export default ApiDocsPage;
