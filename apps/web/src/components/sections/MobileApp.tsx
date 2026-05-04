"use client";

import React from "react";
import { motion } from "framer-motion";
import { Smartphone, Download, Apple, Play, Shield, Globe } from "lucide-react";

const MobileApp = () => {
    return (
        <section className="py-32 bg-background relative overflow-hidden">
            {/* Background Texture */}
            <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-10 pointer-events-none" />

            <div className="container mx-auto px-4 md:px-6 relative z-10">
                <div className="grid grid-cols-1 lg:grid-cols-12 gap-16 items-center">
                    <div className="lg:col-span-7">
                        <motion.div
                            initial={{ opacity: 0, x: -30 }}
                            whileInView={{ opacity: 1, x: 0 }}
                            viewport={{ once: true }}
                            className="space-y-8"
                        >
                            <div className="flex items-center gap-3">
                                <div className="h-px w-12 bg-primary" />
                                <span className="text-[10px] font-black uppercase tracking-widest text-primary">MOBILE INTERFACE</span>
                            </div>

                            <h2 className="text-5xl md:text-8xl font-black italic uppercase tracking-tighter leading-[0.8]">
                                COMMAND <br />
                                <span className="text-white/20">OPERATIONS.</span>
                            </h2>
                            <p className="text-text-dim text-xl font-medium max-w-xl leading-relaxed uppercase tracking-tight italic">
                                Execute large-scale p2p trades, manage secure escrows, and coordinate with verified partners through our mobile terminal.
                            </p>

                            <div className="flex flex-wrap gap-4 pt-4">
                                <button className="flex items-center gap-4 bg-white text-background px-8 py-4 rounded-2xl group transition-all hover:scale-105 active:scale-95 shadow-2xl shadow-white/5">
                                    <Apple className="w-6 h-6 " />
                                    <span className="text-xs font-black uppercase tracking-widest">DEPLOY ON IOS</span>
                                </button>
                                <button className="flex items-center gap-4 bg-surface border border-white/10 text-white px-8 py-4 rounded-2xl group transition-all hover:bg-white/5">
                                    <Play className="w-6 h-6 text-primary" />
                                    <span className="text-xs font-black uppercase tracking-widest">DEPLOY ON ANDROID</span>
                                </button>
                            </div>

                            <div className="grid grid-cols-2 gap-8 pt-12 border-t border-white/5">
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-2 flex items-center gap-2">
                                        <Shield className="w-3 h-3" /> Biometric Link
                                    </p>
                                    <p className="text-xs font-bold text-white uppercase tracking-tight">ENCRYPTED AT REST</p>
                                </div>
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-2 flex items-center gap-2">
                                        <Globe className="w-3 h-3" /> Peer Mesh
                                    </p>
                                    <p className="text-xs font-bold text-white uppercase tracking-tight">GLOBAL CLOUD SYNC</p>
                                </div>
                            </div>
                        </motion.div>
                    </div>

                    <div className="lg:col-span-5 relative flex justify-center">
                        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full bg-primary/20 blur-[150px] -z-10 rounded-full" />

                        <motion.div
                            initial={{ opacity: 0, rotate: -10, y: 50 }}
                            whileInView={{ opacity: 1, rotate: -5, y: 0 }}
                            viewport={{ once: true }}
                            className="relative w-[300px] h-[600px] bg-surface-light border-[1px] border-white/10 rounded-[48px] shadow-[0_50px_100px_-20px_rgba(0,0,0,0.8)] overflow-hidden"
                        >
                            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-6 bg-background rounded-b-2xl z-20" />

                            <div className="p-8 space-y-10">
                                <div className="flex justify-between items-center">
                                    <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center">
                                        <Smartphone className="w-5 h-5 text-primary" />
                                    </div>
                                    <p className="text-[10px] font-black text-text-dim uppercase">TRM-882</p>
                                </div>

                                <div className="space-y-2">
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">ASSET OVERVIEW</p>
                                    <h3 className="text-3xl font-black italic text-white">$142,500.00</h3>
                                </div>

                                <div className="space-y-4">
                                    {[1, 2, 3].map((i) => (
                                        <div key={i} className="p-4 bg-background border border-white/5 rounded-2xl flex justify-between items-center group hover:border-primary/30 transition-all">
                                            <div className="w-2 h-2 rounded-full bg-primary" />
                                            <div className="flex-1 px-4">
                                                <p className="text-[10px] font-black text-white">XCH-40{i}</p>
                                                <p className="text-[8px] text-text-dim">COMPLETED</p>
                                            </div>
                                            <p className="text-[10px] font-black text-primary">+$4.2K</p>
                                        </div>
                                    ))}
                                </div>

                                <div className="pt-4">
                                    <div className="h-1 w-full bg-white/5 rounded-full overflow-hidden">
                                        <div className="h-full w-2/3 bg-primary" />
                                    </div>
                                    <p className="mt-2 text-[8px] font-black text-text-dim uppercase">SYSTEM LOAD: 42%</p>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </div>
            </div>
        </section>
    );
};

export default MobileApp;
