"use client";

import React from "react";
import { UserPlus, ShieldCheck, Wallet } from "lucide-react";
import { motion } from "framer-motion";

const HowItWorks = () => {
    const steps = [
        {
            icon: UserPlus,
            title: "Create Account",
            desc: "Sign up and verify your identity in less than 2 minutes to start trading securely.",
            color: "bg-blue-500/10 text-blue-500",
        },
        {
            icon: ShieldCheck,
            title: "Place an Order",
            desc: "Choose an offer from a verified merchant or create your own custom trade offer.",
            color: "bg-primary/10 text-primary",
        },
        {
            icon: Wallet,
            title: "Get Crypto/Cash",
            desc: "Once payment is verified, the crypto is released from escrow to your wallet instantly.",
            color: "bg-accent/10 text-accent",
        },
    ];

    return (
        <section className="py-24 bg-surface/30">
            <div className="container mx-auto px-4 md:px-6 text-center">
                <div className="mb-16">
                    <h2 className="text-3xl md:text-5xl font-black mb-4">How it works</h2>
                    <p className="text-text-dim max-w-2xl mx-auto">Trading on Cryplio is simple, secure, and fast. Follow these three steps to complete your first trade.</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8 relative">
                    {/* Connecting Line (Desktop) */}
                    <div className="hidden md:block absolute top-1/2 left-0 w-full h-[2px] bg-border -translate-y-1/2 z-0" />

                    {steps.map((step, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, y: 20 }}
                            whileInView={{ opacity: 1, y: 0 }}
                            viewport={{ once: true }}
                            transition={{ delay: i * 0.2 }}
                            className="relative z-10 bg-surface border border-border p-8 rounded-[32px] hover:border-primary/50 transition-colors group"
                        >
                            <div className={`w-16 h-16 ${step.color} rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform`}>
                                <step.icon className="w-8 h-8" />
                            </div>
                            <h3 className="text-xl font-bold mb-4">{step.title}</h3>
                            <p className="text-text-dim text-sm leading-relaxed">{step.desc}</p>
                            <div className="absolute -top-4 -right-4 w-10 h-10 bg-surface border border-border rounded-full flex items-center justify-center font-black text-primary shadow-lg">
                                {i + 1}
                            </div>
                        </motion.div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default HowItWorks;
