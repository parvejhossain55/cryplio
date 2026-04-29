"use client";

import React, { useState } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ShieldAlert,
    Server,
    Globe,
    Lock,
    Database,
    Zap,
    Cpu,
    Activity,
    Terminal,
    Settings,
    ShieldCheck,
    Smartphone,
    ChevronRight,
    LogOut
} from "lucide-react";
import { motion } from "framer-motion";

const AdminSettings = () => {
    const [activeTab, setActiveTab] = useState("system");

    const tabs = [
        { id: "system", label: "Infrastructure", icon: Server },
        { id: "security", label: "Network Security", icon: ShieldCheck },
        { id: "api", label: "API Protocol", icon: Terminal },
        { id: "database", label: "Storage Delta", icon: Database },
        { id: "logs", label: "Kernel Logs", icon: Activity },
    ];

    return (
        <DashboardLayout title="System Kernel" role="admin">
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                {/* Tabs Sidebar */}
                <div className="lg:col-span-1 space-y-4">
                    <div className="p-6 bg-red-500/10 border border-red-500/20 rounded-[2.5rem] mb-6">
                        <div className="flex flex-col items-center">
                            <div className="w-20 h-20 rounded-full bg-red-500/20 border-2 border-red-500/30 flex items-center justify-center">
                                <ShieldAlert className="w-10 h-10 text-red-500" />
                            </div>
                            <h3 className="mt-4 text-xl font-black text-white">Root Admin</h3>
                            <p className="text-[10px] font-medium text-red-500 uppercase tracking-widest mt-1">Full System Access</p>
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
                                <tab.icon className={`w-5 h-5 mr-4 ${activeTab === tab.id ? "text-white" : "text-text-dim group-hover:text-white"}`} />
                                <span className="text-[10px] font-black uppercase tracking-widest">{tab.label}</span>
                                {activeTab === tab.id && <ChevronRight className="ml-auto w-4 h-4" />}
                            </button>
                        ))}
                    </div>

                    <button className="w-full flex items-center px-10 py-5 text-red-500 hover:bg-red-500/10 rounded-[2.5rem] transition-all group font-black uppercase tracking-[0.2em] text-[10px]">
                        <LogOut className="w-5 h-5 mr-4" />
                        Hard Shutdown Session
                    </button>
                </div>

                {/* Content Area */}
                <div className="lg:col-span-3 space-y-8">
                    {activeTab === "system" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <h3 className="text-2xl font-black text-white mb-8 uppercase tracking-tight">Mainframe Controller</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Global Trade Status</label>
                                        <div className="flex items-center justify-between p-4 bg-white/5 border border-white/5 rounded-2xl">
                                            <span className="text-xs font-bold text-white uppercase">Trading Engine</span>
                                            <div className="w-12 h-6 bg-accent rounded-full relative cursor-pointer">
                                                <div className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full" />
                                            </div>
                                        </div>
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Maintenance Overlay</label>
                                        <div className="flex items-center justify-between p-4 bg-white/5 border border-white/5 rounded-2xl">
                                            <span className="text-xs font-bold text-white uppercase">Global Banner</span>
                                            <div className="w-12 h-6 bg-white/10 rounded-full relative cursor-pointer">
                                                <div className="absolute left-1 top-1 w-4 h-4 bg-white/40 rounded-full" />
                                            </div>
                                        </div>
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Escrow Timeout (ms)</label>
                                        <input
                                            type="number"
                                            defaultValue="1800000"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-red-500/50 transition-all font-mono"
                                        />
                                    </div>
                                    <div className="space-y-4">
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">System Fee Protocol (%)</label>
                                        <input
                                            type="text"
                                            defaultValue="0.005"
                                            className="w-full bg-white/5 border border-white/5 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-red-500/50 transition-all font-mono"
                                        />
                                    </div>
                                </div>
                                <div className="pt-10 flex justify-end">
                                    <button className="bg-red-500 text-white px-10 py-4 rounded-2xl font-black uppercase tracking-[0.2em] text-[10px] hover:scale-105 active:scale-95 transition-all shadow-xl shadow-red-500/20">
                                        Override System Config
                                    </button>
                                </div>
                            </div>

                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <div className="flex items-center justify-between mb-10">
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Resource Management</h3>
                                    <Cpu className="w-6 h-6 text-red-500" />
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="p-6 rounded-3xl bg-white/5 border border-white/5">
                                        <div className="flex items-center justify-between mb-4">
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Master Database</p>
                                            <span className="text-[10px] font-black text-accent uppercase">Operational</span>
                                        </div>
                                        <div className="w-full bg-white/5 h-2 rounded-full overflow-hidden mb-2">
                                            <div className="w-[12%] bg-accent h-full" />
                                        </div>
                                        <p className="text-[10px] text-text-dim font-medium uppercase tracking-widest">12% Storage Capability Used</p>
                                    </div>
                                    <div className="p-6 rounded-3xl bg-white/5 border border-white/5">
                                        <div className="flex items-center justify-between mb-4">
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Redis Cache</p>
                                            <span className="text-[10px] font-black text-accent uppercase">Optimized</span>
                                        </div>
                                        <div className="w-full bg-white/5 h-2 rounded-full overflow-hidden mb-2">
                                            <div className="w-[42%] bg-primary h-full" />
                                        </div>
                                        <p className="text-[10px] text-text-dim font-medium uppercase tracking-widest">Latency: 2.4ms</p>
                                    </div>
                                </div>
                            </div>
                        </motion.div>
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminSettings;
