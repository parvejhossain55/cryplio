"use client";

import React, { useState } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    User,
    Shield,
    Bell,
    CreditCard,
    Globe,
    Zap,
    Mail,
    Lock,
    Smartphone,
    ChevronRight,
    LogOut,
    Camera,
    Activity
} from "lucide-react";
import { motion } from "framer-motion";

const UserSettings = () => {
    const [activeTab, setActiveTab] = useState("profile");

    const tabs = [
        { id: "profile", label: "Profile", icon: User },
        { id: "security", label: "Security", icon: Shield },
        { id: "notifications", label: "Alerts", icon: Bell },
        { id: "billing", label: "Payments", icon: CreditCard },
        { id: "preferences", label: "Protocol", icon: Zap },
    ];

    return (
        <DashboardLayout title="System Settings" role="user">
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                {/* Tabs Sidebar */}
                <div className="lg:col-span-1 space-y-2">
                    <div className="p-6 bg-surface border border-white/10 rounded-[2.5rem] mb-6">
                        <div className="flex flex-col items-center">
                            <div className="relative group cursor-pointer">
                                <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center overflow-hidden transition-transform group-hover:scale-105">
                                    <User className="w-10 h-10 text-white/50" />
                                </div>
                                <div className="absolute bottom-0 right-0 w-8 h-8 bg-primary rounded-full border-4 border-surface flex items-center justify-center">
                                    <Camera className="w-3.5 h-3.5 text-white" />
                                </div>
                            </div>
                            <h3 className="mt-4 text-xl font-black text-white">Alex Morgan</h3>
                            <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Pro Trader • Level 2</p>
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-4 flex flex-col space-y-1">
                        {tabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`flex items-center px-6 py-4 rounded-2xl transition-all group ${activeTab === tab.id
                                        ? "bg-primary text-white shadow-xl shadow-primary/20"
                                        : "text-text-dim hover:bg-white/5 hover:text-white"
                                    }`}
                            >
                                <tab.icon className={`w-5 h-5 mr-4 transition-colors ${activeTab === tab.id ? "text-white" : "text-text-dim group-hover:text-white"}`} />
                                <span className="text-sm font-black uppercase tracking-widest">{tab.label}</span>
                                {activeTab === tab.id && <ChevronRight className="ml-auto w-4 h-4" />}
                            </button>
                        ))}
                    </div>

                    <button className="w-full flex items-center px-10 py-5 text-primary hover:bg-primary/10 rounded-[2.5rem] transition-all group border border-transparent mt-4 font-black uppercase tracking-[0.2em] text-xs">
                        <LogOut className="w-5 h-5 mr-4" />
                        Logout Session
                    </button>
                </div>

                {/* Content Area */}
                <div className="lg:col-span-3 space-y-8">
                    {activeTab === "profile" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <h3 className="text-2xl font-black text-white mb-8 uppercase tracking-tight">Identity Profile</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Full Identity Name</label>
                                        <input
                                            type="text"
                                            defaultValue="Alex Morgan"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all"
                                        />
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Network Email</label>
                                        <input
                                            type="email"
                                            defaultValue="alex.morgan@cryplio.com"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all opacity-60"
                                            disabled
                                        />
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Trade Nickname</label>
                                        <input
                                            type="text"
                                            defaultValue="AlexPro_92"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all"
                                        />
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Phone Protocol</label>
                                        <input
                                            type="text"
                                            defaultValue="+1 (555) 482-0192"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary/50 transition-all"
                                        />
                                    </div>
                                </div>
                                <div className="pt-10 flex justify-end">
                                    <button className="bg-primary text-white px-10 py-4 rounded-2xl font-black uppercase tracking-[0.2em] text-xs hover:scale-105 active:scale-95 transition-all shadow-xl shadow-primary/20">
                                        Commit Changes
                                    </button>
                                </div>
                            </div>

                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <div className="flex items-center justify-between mb-8">
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Public Profile Analytics</h3>
                                    <Activity className="w-6 h-6 text-primary" />
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                    {[
                                        { label: "Positive Feedback", value: "99.8%", color: "text-accent" },
                                        { label: "Completion Rate", value: "94.2%", color: "text-white" },
                                        { label: "Avg. Release Time", value: "4.2m", color: "text-secondary" },
                                    ].map((stat, i) => (
                                        <div key={i} className="p-6 rounded-3xl bg-white/5 border border-white/5">
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-2">{stat.label}</p>
                                            <p className={`text-2xl font-black ${stat.color}`}>{stat.value}</p>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </motion.div>
                    )}

                    {activeTab === "security" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <h3 className="text-2xl font-black text-white mb-8 uppercase tracking-tight">Security Protocol</h3>
                                <div className="space-y-6">
                                    {[
                                        { icon: Mail, label: "Email Verification", desc: "Secures account access and transactions.", status: "Active", color: "bg-accent/20 text-accent" },
                                        { icon: Smartphone, label: "SMS Authentication", desc: "Adds an extra layer of security via mobile.", status: "Inactive", color: "bg-white/10 text-text-dim" },
                                        { icon: Lock, label: "Identity Verification (KYC)", desc: "Mandatory for high-volume fiat trading.", status: "Level 2 Completed", color: "bg-accent/20 text-accent" },
                                        { icon: Zap, label: "Anti-Phishing Code", desc: "Verification code in every official email.", status: "Enabled", color: "bg-accent/20 text-accent" },
                                    ].map((item, i) => (
                                        <div key={i} className="flex items-center justify-between p-6 rounded-3xl bg-white/5 border border-white/5 group hover:border-white/20 transition-all">
                                            <div className="flex items-start space-x-6">
                                                <div className="p-4 rounded-2xl bg-surface border border-white/10 group-hover:scale-110 transition-transform">
                                                    <item.icon className="w-6 h-6 text-white" />
                                                </div>
                                                <div>
                                                    <h4 className="font-black text-white uppercase text-sm tracking-widest">{item.label}</h4>
                                                    <p className="text-xs text-text-dim font-medium mt-1">{item.desc}</p>
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <span className={`text-[8px] font-black uppercase tracking-widest px-3 py-1.5 rounded-lg ${item.color}`}>
                                                    {item.status}
                                                </span>
                                                <button className="block text-[8px] font-bold text-primary uppercase mt-2 hover:underline">Configure</button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="bg-primary/10 border border-primary/20 rounded-[2.5rem] p-10 flex items-center justify-between">
                                <div>
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Advanced Entropy Keys</h3>
                                    <p className="text-xs text-text-dim font-medium mt-1">Download your backup security phrase to prevent account lockout.</p>
                                </div>
                                <button className="bg-primary text-white px-8 py-4 rounded-2xl font-black uppercase tracking-[0.2em] text-[10px] hover:scale-105 transition-all shadow-xl shadow-primary/20 whitespace-nowrap">
                                    Export Keys
                                </button>
                            </div>
                        </motion.div>
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserSettings;
