"use client";

import React from "react";
import { motion } from "framer-motion";
import { Smartphone, Download, Apple, Play } from "lucide-react";

const MobileApp = () => {
    return (
        <section className="py-24 bg-background overflow-hidden">
            <div className="container mx-auto px-4 md:px-6">
                <div className="flex flex-col lg:flex-row items-center gap-16">
                    <div className="flex-1 text-left">
                        <motion.div
                            initial={{ opacity: 0, x: -30 }}
                            whileInView={{ opacity: 1, x: 0 }}
                            viewport={{ once: true }}
                        >
                            <h2 className="text-4xl md:text-6xl font-black mb-8 leading-tight">
                                Trade on the go <br />
                                with the <span className="gradient-text">Cryplio App</span>
                            </h2>
                            <p className="text-text-dim text-lg mb-10 max-w-xl leading-relaxed">
                                Experience the power of the world's most advanced P2P exchange in your pocket.
                                Manage your wallet, chat with traders, and execute orders anywhere, anytime.
                            </p>

                            <div className="flex flex-wrap gap-4">
                                <button className="flex items-center space-x-3 bg-surface border border-border px-6 py-3 rounded-2xl hover:border-primary/50 transition-colors group">
                                    <Apple className="w-8 h-8 group-hover:scale-110 transition-transform" />
                                    <div className="text-left">
                                        <p className="text-[10px] font-bold text-text-dim uppercase tracking-widest leading-none">Download on</p>
                                        <p className="text-lg font-black leading-none mt-1">App Store</p>
                                    </div>
                                </button>
                                <button className="flex items-center space-x-3 bg-surface border border-border px-6 py-3 rounded-2xl hover:border-primary/50 transition-colors group">
                                    <Play className="w-8 h-8 group-hover:scale-110 transition-transform" />
                                    <div className="text-left">
                                        <p className="text-[10px] font-bold text-text-dim uppercase tracking-widest leading-none">Get it on</p>
                                        <p className="text-lg font-black leading-none mt-1">Google Play</p>
                                    </div>
                                </button>
                            </div>
                        </motion.div>
                    </div>

                    <div className="flex-1 relative flex justify-center">
                        {/* Background Glow */}
                        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-primary/20 rounded-full blur-[120px] pointer-events-none" />

                        {/* Phone Mockup Frame */}
                        <motion.div
                            initial={{ opacity: 0, scale: 0.8, y: 50 }}
                            whileInView={{ opacity: 1, scale: 1, y: 0 }}
                            viewport={{ once: true }}
                            className="relative w-[280px] h-[580px] bg-[#0f1115] border-[8px] border-[#25282c] rounded-[48px] shadow-2xl overflow-hidden p-4"
                        >
                            {/* Screen Content Mock */}
                            <div className="w-full h-full rounded-[32px] bg-background flex flex-col p-4">
                                <div className="flex justify-between items-center mb-6 px-2">
                                    <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                                        <Smartphone className="w-4 h-4 text-white" />
                                    </div>
                                    <div className="w-8 h-8 rounded-full border border-border" />
                                </div>

                                <p className="text-xs font-bold text-text-dim mb-4">Total Balance</p>
                                <p className="text-2xl font-black mb-8">$12,450.80</p>

                                <div className="grid grid-cols-2 gap-3 mb-8">
                                    <div className="bg-surface p-3 rounded-2xl border border-border">
                                        <div className="w-8 h-8 bg-primary/20 rounded-lg mb-2 flex items-center justify-center">
                                            <Download className="w-4 h-4 text-primary" />
                                        </div>
                                        <p className="text-[10px] font-bold">Buy</p>
                                    </div>
                                    <div className="bg-surface p-3 rounded-2xl border border-border">
                                        <div className="w-8 h-8 bg-accent/20 rounded-lg mb-2 flex items-center justify-center">
                                            <Download className="w-4 h-4 text-accent rotate-180" />
                                        </div>
                                        <p className="text-[10px] font-bold">Sell</p>
                                    </div>
                                </div>

                                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim mb-4">Recent Trades</p>
                                <div className="space-y-3">
                                    {[1, 2, 3].map((i) => (
                                        <div key={i} className="flex justify-between items-center p-3 bg-surface rounded-xl border border-border">
                                            <div className="flex items-center space-x-2">
                                                <div className="w-6 h-6 bg-white/5 rounded-full" />
                                                <div className="space-y-1">
                                                    <p className="text-[10px] font-bold">ID: #432{i}</p>
                                                    <p className="text-[8px] text-text-dim">Completed</p>
                                                </div>
                                            </div>
                                            <p className="text-[10px] font-bold text-accent">+$120.00</p>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </motion.div>

                        {/* Decoration Floating Elements */}
                        <motion.div
                            animate={{ y: [0, -20, 0] }}
                            transition={{ duration: 4, repeat: Infinity }}
                            className="absolute top-20 -right-4 bg-surface p-4 rounded-2xl border border-border shadow-2xl z-20"
                        >
                            <p className="text-xs font-black text-accent">+2.4% BTC</p>
                        </motion.div>
                        <motion.div
                            animate={{ y: [0, 20, 0] }}
                            transition={{ duration: 5, repeat: Infinity, delay: 1 }}
                            className="absolute bottom-40 -left-10 bg-surface p-4 rounded-2xl border border-border shadow-2xl z-20"
                        >
                            <p className="text-xs font-black text-primary">Trade Confirmed</p>
                        </motion.div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default MobileApp;
