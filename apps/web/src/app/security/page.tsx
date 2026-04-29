"use client";

import React from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
    ShieldCheck,
    Lock,
    Key,
    Smartphone,
    ShieldAlert,
    Server,
    UserCheck,
    Eye,
    Zap,
    Globe,
    Wallet
} from "lucide-react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";

const SecurityPage = () => {
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
                            <ShieldCheck className="w-4 h-4 text-primary" />
                            <span className="text-xs font-black uppercase tracking-widest text-text-dim">Institutional Grade Security</span>
                        </div>
                        <h1 className="text-5xl md:text-8xl font-black tracking-tight leading-[0.85] uppercase">
                            Your Assets <br />
                            <span className="gradient-text">Always Protected.</span>
                        </h1>
                        <p className="max-w-2xl mx-auto text-text-dim text-lg md:text-xl font-medium pt-4">
                            Cryplio employs a multi-layered security architecture to ensure your funds and data remain safe from any threat.
                        </p>
                    </motion.div>
                </div>
            </section>

            {/* Security Features Grid */}
            <section className="py-24 border-t border-white/5">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        {[
                            {
                                icon: Lock,
                                title: "Escrow Protection",
                                desc: "Every P2P trade is backed by our secure escrow protocol, ensuring funds are released only when both parties fulfill their obligations."
                            },
                            {
                                icon: Smartphone,
                                title: "2FA Everywhere",
                                desc: "Mandatory multi-factor authentication for withdrawals, critical settings changes, and logins from new devices."
                            },
                            {
                                icon: Server,
                                title: "Cold Storage",
                                desc: "98% of digital assets are stored in geographically distributed offline cold wallets, protected by multi-signature technology."
                            },
                            {
                                icon: Eye,
                                title: "Real-time Monitoring",
                                desc: "Our AI-driven security systems monitor all transactions 24/7 to identify and block suspicious activities instantly."
                            },
                            {
                                icon: ShieldAlert,
                                title: "Anti-Phishing",
                                desc: "Personalized anti-phishing codes in every email ensure you're always communicating with the genuine Cryplio platform."
                            },
                            {
                                icon: UserCheck,
                                title: "Advanced KYC",
                                desc: "Rigorous identity verification processes to maintain a clean ecosystem and prevent fraudulent actors from joining."
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

            {/* Bottom CTA / Stats */}
            <section className="py-24 relative overflow-hidden">
                <div className="container mx-auto px-4">
                    <div className="bg-surface border border-border rounded-[48px] p-12 md:p-20 relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-1/2 h-full bg-primary/5 blur-[100px] pointer-events-none" />

                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
                            <div>
                                <h2 className="text-4xl md:text-6xl font-black leading-tight mb-8">
                                    STAY <span className="text-primary">SECURE</span> <br />
                                    WITH THE BEST.
                                </h2>
                                <div className="space-y-6">
                                    {[
                                        "ISO/IEC 27001 Certified Infrastructure",
                                        "Bi-annual Comprehensive Security Audits",
                                        "Global Bounty Programs for Vulnerabilities",
                                        "Insurance Fund for Protocol Security"
                                    ].map((item, i) => (
                                        <div key={i} className="flex items-center space-x-4">
                                            <div className="w-6 h-6 rounded-full bg-primary/20 flex items-center justify-center">
                                                <div className="w-2 h-2 bg-primary rounded-full" />
                                            </div>
                                            <span className="font-bold text-lg text-white/90">{item}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="relative">
                                <div className="glass p-8 rounded-[32px] border-border text-center space-y-6">
                                    <Key className="w-12 h-12 text-primary mx-auto" />
                                    <h4 className="text-2xl font-black">Ready to start trading?</h4>
                                    <p className="text-text-dim font-medium">Create your secure account in less than 2 minutes.</p>
                                    <div className="flex flex-col space-y-4 pt-4">
                                        <Link href="/register" className="w-full bg-white text-background py-5 rounded-2xl font-black text-lg transition-all hover:scale-[1.02] active:scale-95 flex items-center justify-center">
                                            Get Started Now
                                        </Link>
                                    </div>
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

export default SecurityPage;
