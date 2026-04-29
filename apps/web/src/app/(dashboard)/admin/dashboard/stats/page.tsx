"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    BarChart3,
    TrendingUp,
    Users,
    ShieldCheck,
    Activity,
    Globe,
    ArrowUpRight,
    ArrowDownLeft,
    Server,
    Database,
    Cpu
} from "lucide-react";
import { motion } from "framer-motion";

const AdminStats = () => {
    return (
        <DashboardLayout title="Network Analytics" role="admin">
            <div className="space-y-8">
                {/* Real-time Ticker */}
                <div className="bg-red-500/10 border border-red-500/20 rounded-2xl p-4 flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <Activity className="w-4 h-4 text-red-500 animate-pulse" />
                        <span className="text-[10px] font-black text-white uppercase tracking-[0.2em]">Platform Pulse: 100% Operational • 42,012 Active Connections</span>
                    </div>
                    <span className="text-[10px] font-black text-white/40 uppercase tracking-widest">Last Kernel Sync: 0.2s ago</span>
                </div>

                {/* Major Metrics */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    {[
                        { label: "Gross Platform Volume", value: "$12.4M", change: "+12.4%", icon: Globe },
                        { label: "Aggregated Fees", value: "$62.1K", change: "+4.2%", icon: TrendingUp },
                        { label: "Active Nodes", value: "24", change: "Stable", icon: Server },
                        { label: "Network Latency", value: "14ms", change: "-2ms", icon: Cpu },
                    ].map((stat, i) => (
                        <div key={i} className="bg-surface border border-white/10 p-8 rounded-[2rem] hover:border-white/20 transition-all group">
                            <div className="flex items-center justify-between mb-4">
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">{stat.label}</p>
                                <stat.icon className="w-5 h-5 text-red-500 group-hover:rotate-12 transition-transform" />
                            </div>
                            <h3 className="text-3xl font-black text-white mb-2">{stat.value}</h3>
                            <p className={`text-[10px] font-bold ${stat.change.startsWith('+') ? 'text-accent' : 'text-primary'}`}>
                                {stat.change} <span className="text-text-dim ml-1">since last epoch</span>
                            </p>
                        </div>
                    ))}
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    {/* Traffic Source */}
                    <div className="lg:col-span-2 bg-surface border border-white/10 rounded-[2.5rem] p-10">
                        <div className="flex items-center justify-between mb-10">
                            <div>
                                <h3 className="text-xl font-black text-white uppercase tracking-tight">Geographical Traffic</h3>
                                <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Global node distribution</p>
                            </div>
                            <button className="text-[10px] font-black text-red-500 uppercase border border-red-500/20 px-4 py-2 rounded-xl hover:bg-red-500/10 transition-all">Export Geodata</button>
                        </div>

                        <div className="space-y-8">
                            {[
                                { region: "North America", traffic: "42%", volume: "$5.2M", color: "bg-primary" },
                                { region: "European Union", traffic: "28%", volume: "$3.4M", color: "bg-secondary" },
                                { region: "Asia Pacific", traffic: "20%", volume: "$2.5M", color: "bg-accent" },
                                { region: "Other", traffic: "10%", volume: "$1.3M", color: "bg-white/20" },
                            ].map((region, i) => (
                                <div key={i}>
                                    <div className="flex items-center justify-between mb-2 px-1">
                                        <span className="text-xs font-black text-white uppercase tracking-widest">{region.region}</span>
                                        <span className="text-xs font-black text-text-dim uppercase">{region.volume} • {region.traffic}</span>
                                    </div>
                                    <div className="w-full bg-white/5 h-3 rounded-full overflow-hidden border border-white/5">
                                        <div className={`h-full ${region.color} rounded-full`} style={{ width: region.traffic }} />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* System Health */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                        <h3 className="text-xl font-black text-white mb-10 uppercase tracking-tight">Kernel Health</h3>
                        <div className="space-y-6">
                            {[
                                { label: "Ingress Flow", status: "Optimal", color: "text-accent" },
                                { label: "Transaction Relay", status: "Active", color: "text-accent" },
                                { label: "Blockchain Sync", status: "Lagging (2s)", color: "text-primary" },
                                { label: "Storage Delta", status: "Balanced", color: "text-accent" },
                                { label: "Security Mesh", status: "Shield Active", color: "text-secondary" },
                            ].map((item, i) => (
                                <div key={i} className="flex items-center justify-between group">
                                    <span className="text-[10px] font-black text-text-dim uppercase tracking-widest">{item.label}</span>
                                    <span className={`text-[10px] font-black uppercase tracking-widest px-3 py-1 rounded-lg bg-white/5 border border-white/5 ${item.color}`}>{item.status}</span>
                                </div>
                            ))}
                        </div>
                        <div className="mt-10 pt-10 border-t border-white/5">
                            <button className="w-full py-5 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl font-black uppercase tracking-[0.2em] text-[10px] hover:bg-red-500/20 transition-all">
                                Run Global Diagnostic
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminStats;
