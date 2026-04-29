"use client";

import React, { useState } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    Store,
    Shield,
    Zap,
    Globe,
    MessageSquare,
    Lock,
    Smartphone,
    ChevronRight,
    LogOut,
    Settings,
    CreditCard,
    CheckCircle2,
    AlertTriangle
} from "lucide-react";
import { motion } from "framer-motion";

const MerchantSettings = () => {
    const [activeTab, setActiveTab] = useState("store");

    const tabs = [
        { id: "store", label: "Store Portal", icon: Store },
        { id: "verification", label: "Verification", icon: Shield },
        { id: "payments", label: "Payments", icon: CreditCard },
        { id: "security", label: "Security", icon: Lock },
        { id: "automation", label: "Automation", icon: Zap },
    ];

    return (
        <DashboardLayout title="Merchant Config" role="merchant">
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                {/* Tabs Sidebar */}
                <div className="lg:col-span-1 space-y-4">
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-4 space-y-1">
                        {tabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`w-full flex items-center px-6 py-4 rounded-2xl transition-all group ${activeTab === tab.id
                                        ? "bg-primary text-white shadow-xl shadow-primary/20"
                                        : "text-text-dim hover:bg-white/5 hover:text-white"
                                    }`}
                            >
                                <tab.icon className={`w-5 h-5 mr-4 ${activeTab === tab.id ? "text-white" : "text-text-dim group-hover:text-white"}`} />
                                <span className="text-[10px] font-black uppercase tracking-[0.15em]">{tab.label}</span>
                                {activeTab === tab.id && <ChevronRight className="ml-auto w-4 h-4" />}
                            </button>
                        ))}
                    </div>

                    <div className="bg-accent/10 border border-accent/20 rounded-[2.5rem] p-6">
                        <div className="flex items-center space-x-3 mb-4">
                            <CheckCircle2 className="w-5 h-5 text-accent" />
                            <h4 className="text-[10px] font-black text-white uppercase tracking-widest">Trust Index</h4>
                        </div>
                        <p className="text-[10px] text-text-dim font-medium uppercase tracking-[0.15em] leading-relaxed mb-4">You are currently ranked in the top 5% of all merchants.</p>
                        <div className="w-full bg-white/5 h-1.5 rounded-full overflow-hidden">
                            <div className="w-[98%] bg-accent h-full rounded-full" />
                        </div>
                    </div>

                    <button className="w-full flex items-center px-10 py-5 text-primary hover:bg-primary/10 rounded-[2.5rem] transition-all group font-black uppercase tracking-[0.2em] text-[10px]">
                        <LogOut className="w-5 h-5 mr-4" />
                        Terminate Session
                    </button>
                </div>

                {/* Content Area */}
                <div className="lg:col-span-3 space-y-8">
                    {activeTab === "store" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <h3 className="text-2xl font-black text-white mb-8 uppercase tracking-tight">Merchant Identity</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Display Store Name</label>
                                        <input
                                            type="text"
                                            defaultValue="AlexPro Liquidity"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all font-mono"
                                        />
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Auto-Response Message</label>
                                        <input
                                            type="text"
                                            defaultValue="Fast release. Online 24/7."
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all"
                                        />
                                    </div>
                                    <div className="md:col-span-2 space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Trade Terms Protocol</label>
                                        <textarea
                                            className="w-full bg-white/5 border border-white/5 py-6 px-8 rounded-3xl text-sm font-bold outline-none focus:border-primary/50 transition-all min-h-[150px]"
                                            defaultValue="1. Only accept payments from verified accounts. 2. Reference must be provided. 3. Fast Release."
                                        />
                                    </div>
                                </div>
                                <div className="pt-10 flex justify-end">
                                    <button className="bg-primary text-white px-10 py-4 rounded-2xl font-black uppercase tracking-[0.2em] text-[10px] hover:scale-105 active:scale-95 transition-all shadow-xl shadow-primary/20">
                                        Save Protocol
                                    </button>
                                </div>
                            </div>
                        </motion.div>
                    )}

                    {activeTab === "verification" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <div className="flex items-center justify-between mb-10">
                                    <h3 className="text-2xl font-black text-white uppercase tracking-tight">Trust Verification</h3>
                                    <Shield className="w-8 h-8 text-accent" />
                                </div>
                                <div className="space-y-6">
                                    {[
                                        { label: "Merchant KYC", status: "Verified", desc: "Allows unlimited trading volume.", color: "text-accent" },
                                        { label: "Security Deposit", status: "Active ($2,500)", desc: "Escrow guarantee for clients.", color: "text-accent" },
                                        { label: "Business Credentials", status: "Pending", desc: "Corporate account status.", color: "text-primary" },
                                    ].map((item, i) => (
                                        <div key={i} className="flex items-center justify-between p-6 rounded-3xl bg-white/5 border border-white/5 group hover:bg-white/[0.08] transition-all">
                                            <div>
                                                <h4 className="text-xs font-black text-white uppercase tracking-widest">{item.label}</h4>
                                                <p className="text-[10px] text-text-dim font-medium mt-1">{item.desc}</p>
                                            </div>
                                            <span className={`text-[10px] font-black uppercase tracking-widest ${item.color}`}>{item.status}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10 overflow-hidden relative">
                                <div className="absolute top-0 right-0 p-10 opacity-5">
                                    <AlertTriangle className="w-32 h-32 text-primary" />
                                </div>
                                <h3 className="text-xl font-black text-white mb-2 uppercase tracking-tight relative z-10">Level 5 Upgrade</h3>
                                <p className="text-xs text-text-dim font-medium mb-8 max-w-lg relative z-10">
                                    Complete 50 more trades with 100% positive feedback to unlock the "Certified Liquidator" badge and 0.1% fee reduction.
                                </p>
                                <div className="flex items-center space-x-4 relative z-10">
                                    <div className="flex-1 bg-white/5 h-2 rounded-full overflow-hidden border border-white/5">
                                        <div className="w-[85%] bg-primary h-full rounded-full" />
                                    </div>
                                    <span className="text-[10px] font-black text-primary uppercase tracking-widest">85% Progress</span>
                                </div>
                            </div>
                        </motion.div>
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
};

export default MerchantSettings;
