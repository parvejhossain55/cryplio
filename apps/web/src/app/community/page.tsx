"use client";

import React from "react";
import { motion } from "framer-motion";
import {
    Users,
    MessageCircle,
    Globe,
    Code2,
    MessagesSquare,
    ShieldCheck,
    Zap,
    Send,
    ExternalLink,
    Bot
} from "lucide-react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";

const CommunityPage = () => {
    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Hero Section */}
            <section className="pt-32 pb-20 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-primary/5 blur-[120px] rounded-full -z-10" />

                <div className="container mx-auto px-4 text-center">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="space-y-6"
                    >
                        <div className="inline-flex items-center space-x-2 bg-surface border border-border px-4 py-2 rounded-2xl mx-auto">
                            <Users className="w-4 h-4 text-primary" />
                            <span className="text-xs font-black uppercase tracking-widest text-text-dim">Worldwide Trader Network</span>
                        </div>
                        <h1 className="text-5xl md:text-8xl font-black tracking-tight leading-[0.85] uppercase">
                            Join the <br />
                            <span className="gradient-text">Cryplio Pulse.</span>
                        </h1>
                        <p className="max-w-2xl mx-auto text-text-dim text-lg md:text-xl font-medium pt-4">
                            Connect with over 1M+ traders, developers, and blockchain enthusiasts building the future of P2P exchange.
                        </p>
                    </motion.div>
                </div>
            </section>

            {/* Social Channels Grid */}
            <section className="py-24 border-t border-white/5">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                        {[
                            {
                                icon: Send,
                                title: "Telegram",
                                stats: "450K+ Members",
                                desc: "Get real-time updates, price alerts, and chat with the core team instantly.",
                                color: "bg-[#0088cc]/10 text-[#0088cc]",
                                link: "@CryplioOfficial"
                            },
                            {
                                icon: MessagesSquare,
                                title: "Discord",
                                stats: "120K+ Traders",
                                desc: "The official hub for merchant support, developer discussion, and market calls.",
                                color: "bg-[#5865F2]/10 text-[#5865F2]",
                                link: "discord.gg/cryplio"
                            },
                            {
                                icon: Globe,
                                title: "X (Twitter)",
                                stats: "2.1M+ Followers",
                                desc: "Stay informed on platform news, major listings, and community contests.",
                                color: "bg-white/10 text-white",
                                link: "@Cryplio"
                            },
                            {
                                icon: MessageCircle,
                                title: "Reddit",
                                stats: "90K+ Traders",
                                desc: "Deeper discussions on platform features, market trends, and general crypto news.",
                                color: "bg-[#ff4500]/10 text-[#ff4500]",
                                link: "r/cryplio"
                            },
                            {
                                icon: Code2,
                                title: "GitHub",
                                stats: "5K+ Stars",
                                desc: "Explore our open-source tools, liquidity SDKs, and security scripts.",
                                color: "bg-white/10 text-white",
                                link: "github.com/cryplio"
                            },
                            {
                                icon: Bot,
                                title: "Status Bot",
                                stats: "24/7 Monitoring",
                                desc: "Follow our automated status bot for real-time network and liquidity reports.",
                                color: "bg-primary/10 text-primary",
                                link: "@CryplioStatus"
                            }
                        ].map((channel, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 30 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                transition={{ delay: i * 0.1 }}
                                className="glass p-10 rounded-[40px] border-border hover:border-primary/30 transition-all group relative overflow-hidden"
                            >
                                <div className={`w-14 h-14 ${channel.color} rounded-2xl flex items-center justify-center mb-8 group-hover:scale-110 transition-transform shadow-2xl`}>
                                    <channel.icon className="w-7 h-7" />
                                </div>
                                <div className="flex items-center justify-between mb-4">
                                    <h3 className="text-2xl font-black uppercase tracking-tight">{channel.title}</h3>
                                    <span className="text-[10px] font-black uppercase tracking-widest text-text-dim bg-white/5 px-2 py-1 rounded-lg">{channel.stats}</span>
                                </div>
                                <p className="text-text-dim leading-relaxed font-medium mb-8">
                                    {channel.desc}
                                </p>
                                <div className="flex items-center justify-between border-t border-white/5 pt-6">
                                    <span className="text-sm font-bold text-white/50">{channel.link}</span>
                                    <ExternalLink className="w-5 h-5 text-primary group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
                                </div>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Newsletter Section */}
            <section className="py-24 relative overflow-hidden">
                <div className="container mx-auto px-4">
                    <div className="bg-surface border border-border rounded-[48px] p-12 md:p-20 text-center relative overflow-hidden">
                        <div className="absolute top-0 left-0 w-full h-full bg-primary/5 blur-[120px] pointer-events-none" />
                        <h2 className="text-4xl md:text-6xl font-black mb-8 leading-tight">
                            NEVER MISS <br />
                            A <span className="text-primary">PULSE.</span>
                        </h2>
                        <p className="max-w-xl mx-auto text-text-dim text-lg font-medium mb-10">
                            Weekly insights on P2P market trends, platform updates, and exclusive merchant offers. Join 500k+ subscribers.
                        </p>
                        <form className="max-w-md mx-auto flex flex-col sm:flex-row gap-4">
                            <input
                                type="email"
                                placeholder="Enter your email"
                                className="flex-1 bg-background border border-border rounded-2xl px-6 py-4 outline-none focus:border-primary transition-all font-medium"
                            />
                            <button className="bg-white text-background px-8 py-4 rounded-2xl font-black transition-all hover:scale-[1.02] active:scale-95 shadow-2xl">
                                Subscribe
                            </button>
                        </form>
                    </div>
                </div>
            </section>

            <Footer />
        </main>
    );
};

export default CommunityPage;
