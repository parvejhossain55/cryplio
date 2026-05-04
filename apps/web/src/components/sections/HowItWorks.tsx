"use client";

import React from "react";
import { UserPlus, ShieldCheck, Wallet, ArrowRight } from "lucide-react";
import { motion } from "framer-motion";

const HowItWorks = () => {
    const steps = [
        {
            icon: UserPlus,
            title: "INITIALIZE",
            desc: "Onboard your identity through our zero-knowledge verification protocol in < 120s.",
            color: "bg-primary/10 text-primary",
        },
        {
            icon: ShieldCheck,
            title: "EXECUTE",
            desc: "Select an institutional provider or deploy your own liquidity offer to the mesh.",
            color: "bg-surface text-white border border-white/5",
        },
        {
            icon: Wallet,
            title: "SETTLE",
            desc: "Assets are cryptographically verified and released from escrow to your secure vault.",
            color: "bg-primary text-background",
        },
    ];

    return (
        <section className="py-32 bg-background relative overflow-hidden">
            <div className="container mx-auto px-4 md:px-6">
                <div className="flex flex-col md:flex-row md:items-end justify-between mb-20 gap-8">
                    <div className="max-w-3xl">
                        <p className="text-[10px] font-black uppercase tracking-widest text-primary mb-4">PROTOCOL FLOW</p>
                        <h2 className="text-5xl md:text-8xl font-black italic uppercase tracking-tighter leading-[0.8]">
                            CLEARING <br />
                            <span className="text-white/20">PROCEDURE.</span>
                        </h2>
                    </div>
                    <p className="text-text-dim max-w-sm text-sm font-medium leading-relaxed uppercase tracking-widest">
                        Cryplio utilizes a distributed settlement layer to ensure absolute security for every P2P interaction.
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-12">
                    {steps.map((step, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, x: -20 }}
                            whileInView={{ opacity: 1, x: 0 }}
                            viewport={{ once: true }}
                            transition={{ delay: i * 0.1 }}
                            className="group"
                        >
                            <div className="flex items-center gap-6 mb-8">
                                <div className={`w-16 h-16 ${step.color} rounded-2xl flex items-center justify-center font-black group-hover:rotate-12 transition-transform duration-500`}>
                                    <step.icon className="w-8 h-8" />
                                </div>
                                <div className="h-px flex-1 bg-white/5 group-last:hidden" />
                            </div>

                            <p className="text-[10px] font-black text-text-dim mb-3">STEP 0{i + 1}</p>
                            <h3 className="text-2xl font-black italic uppercase tracking-tight mb-4 group-hover:text-primary transition-colors">{step.title}</h3>
                            <p className="text-text-dim text-sm font-bold uppercase tracking-widest leading-loose">
                                {step.desc}
                            </p>
                        </motion.div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default HowItWorks;
