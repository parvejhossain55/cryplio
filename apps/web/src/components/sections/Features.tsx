"use client";

import React from "react";
import { motion } from "framer-motion";
import { Shield, Zap, HeartHandshake, UserCheck, Smartphone, Globe2 } from "lucide-react";

const features = [
    {
        title: "Secure Escrow",
        description: "Funds are locked in a smart escrow until both parties confirm fulfillment.",
        icon: Shield,
        color: "text-blue-500",
        className: "md:col-span-2 lg:col-span-2 h-[300px]",
        glow: "from-blue-500/20 to-transparent"
    },
    {
        title: "0% Maker Fees",
        description: "Create ads and trade for free. We don't charge makers.",
        icon: UserCheck,
        color: "text-emerald-500",
        className: "md:col-span-1 lg:col-span-1 h-[300px]",
        glow: "from-emerald-500/20 to-transparent"
    },
    {
        title: "Instant Settlement",
        description: "Average trade completion time is under 15 minutes.",
        icon: Zap,
        color: "text-amber-500",
        className: "md:col-span-1 lg:col-span-1 h-[300px]",
        glow: "from-amber-500/20 to-transparent"
    },
    {
        title: "Global Liquidity",
        description: "Connect with thousands of verified traders from 50+ countries.",
        icon: Globe2,
        color: "text-cyan-500",
        className: "md:col-span-2 lg:col-span-2 h-[300px]",
        glow: "from-cyan-500/20 to-transparent"
    },
];

const Features = () => {
    return (
        <section className="py-32 px-4 md:px-6 relative overflow-hidden">
            {/* Background Glows */}
            <div className="absolute top-1/2 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-[120px] -z-10" />
            <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-accent/10 rounded-full blur-[120px] -z-10" />

            <div className="container mx-auto">
                <div className="text-center max-w-3xl mx-auto mb-20">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        whileInView={{ opacity: 1, y: 0 }}
                        viewport={{ once: true }}
                    >
                        <h2 className="text-4xl md:text-6xl font-black mb-8 leading-tight">
                            Engineered for <br />
                            <span className="gradient-text">Trust & Performance</span>
                        </h2>
                        <p className="text-text-dim text-lg md:text-xl leading-relaxed">
                            We've built the most secure p2p protocol to ensure your assets are always protected and your trades are always fast.
                        </p>
                    </motion.div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {features.map((feature, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, scale: 0.95 }}
                            whileInView={{ opacity: 1, scale: 1 }}
                            viewport={{ once: true }}
                            transition={{ delay: i * 0.1 }}
                            className={`group relative glass rounded-[40px] border-border p-10 overflow-hidden flex flex-col justify-between hover:border-primary/50 transition-all duration-500 ${feature.className}`}
                        >
                            {/* Inner Glow */}
                            <div className={`absolute -bottom-20 -right-20 w-40 h-40 bg-gradient-to-br ${feature.glow} blur-[60px] opacity-0 group-hover:opacity-100 transition-opacity duration-500`} />

                            <div className="relative z-10">
                                <div className={`w-16 h-16 rounded-[24px] bg-surface-light border border-border flex items-center justify-center mb-10 group-hover:bg-background transition-colors`}>
                                    <feature.icon className={`w-8 h-8 ${feature.color}`} />
                                </div>
                                <h3 className="text-2xl md:text-3xl font-black mb-4 tracking-tight">{feature.title}</h3>
                                <p className="text-text-dim text-lg leading-relaxed max-w-[280px]">
                                    {feature.description}
                                </p>
                            </div>

                            <div className="relative z-10 mt-6 flex justify-end">
                                <div className="w-10 h-10 rounded-full border border-border flex items-center justify-center group-hover:bg-primary group-hover:border-primary transition-all duration-300">
                                    <Zap className="w-4 h-4 group-hover:text-white transition-colors" />
                                </div>
                            </div>
                        </motion.div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default Features;
