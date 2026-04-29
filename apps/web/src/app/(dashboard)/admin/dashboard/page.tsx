"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ShieldAlert,
    Users,
    BarChart3,
    Server,
    ShieldCheck,
    AlertTriangle,
    ChevronRight,
    Activity,
    Globe,
    Database
} from "lucide-react";

const AdminDashboard = () => {
    return (
        <DashboardLayout title="Control Center" role="admin">
            {/* System Health Overview */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
                <div className="md:col-span-2 bg-surface border border-white/10 rounded-[2.5rem] p-8 overflow-hidden relative">
                    <div className="absolute top-0 right-0 p-8 opacity-10">
                        <Activity className="w-40 h-40 text-primary" />
                    </div>
                    <div className="relative z-10">
                        <h3 className="text-xl font-black text-white mb-8">System Infrastructure</h3>
                        <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
                            {[
                                { label: "API Latency", value: "24ms", icon: Server, status: "healthy" },
                                { label: "Blockchain Sync", value: "99.9%", icon: Database, status: "healthy" },
                                { label: "Active Sessions", value: "14.2k", icon: Globe, status: "healthy" },
                                { label: "Pending KYC", value: "420", icon: Users, status: "warning" },
                            ].map((sys, i) => (
                                <div key={i} className="p-4 rounded-2xl bg-white/5 border border-white/5">
                                    <div className="flex items-center space-x-2 mb-3">
                                        <sys.icon className="w-4 h-4 text-text-dim" />
                                        <div className={`w-2 h-2 rounded-full ${sys.status === 'healthy' ? 'bg-accent' : 'bg-primary'}`} />
                                    </div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{sys.label}</p>
                                    <p className="text-xl font-black text-white mt-1">{sys.value}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                <div className="bg-primary/10 border border-primary/20 rounded-[2.5rem] p-8 flex flex-col justify-between">
                    <div>
                        <h3 className="text-xl font-black text-white mb-2">Security Protocol</h3>
                        <p className="text-xs text-text-dim font-medium mb-8">All systems operating within normal parameters.</p>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/5">
                                <span className="text-[10px] font-black uppercase tracking-widest text-white">Cold Storage</span>
                                <span className="text-[10px] font-black uppercase text-accent tracking-widest">Secured</span>
                            </div>
                            <div className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/5">
                                <span className="text-[10px] font-black uppercase tracking-widest text-white">Hot Wallet</span>
                                <span className="text-[10px] font-black uppercase text-accent tracking-widest">Normal</span>
                            </div>
                        </div>
                    </div>
                    <button className="w-full mt-8 py-4 bg-primary text-white rounded-2xl font-black text-xs uppercase tracking-[0.2em] hover:scale-105 active:scale-95 transition-all">
                        Emergency Lockdown
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* User Moderation Queue */}
                <div className="lg:col-span-2 bg-surface border border-white/10 rounded-[2.5rem] p-8">
                    <div className="flex items-center justify-between mb-8">
                        <div>
                            <h3 className="text-xl font-black text-white uppercase tracking-tight">Compliance Queue</h3>
                            <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Users awaiting manual review</p>
                        </div>
                        <button className="text-xs font-bold text-primary hover:text-white transition-colors">
                            View Full Audit Log
                        </button>
                    </div>

                    <div className="space-y-4">
                        {[
                            { user: "Merchant_X", issue: "Large Transaction Alert", amount: "$150,000", time: "10 min ago", severity: "high" },
                            { user: "CryptoDave", issue: "Identity Verification (KYC)", amount: "-", time: "25 min ago", severity: "medium" },
                            { user: "QuickTrade", issue: "Dispute: Order #8291", amount: "$4,200", time: "45 min ago", severity: "high" },
                            { user: "Z_Traders", issue: "Address Change Verification", amount: "-", time: "2h ago", severity: "low" },
                        ].map((item, i) => (
                            <div key={i} className="flex items-center justify-between p-4 rounded-2x border border-white/5 hover:border-white/20 hover:bg-white/[0.02] transition-all group">
                                <div className="flex items-center space-x-4">
                                    <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${item.severity === 'high' ? 'bg-primary/20 text-primary' :
                                            item.severity === 'medium' ? 'bg-secondary/20 text-secondary' :
                                                'bg-accent/20 text-accent'
                                        }`}>
                                        {item.severity === 'high' ? <ShieldAlert className="w-5 h-5" /> : <ShieldCheck className="w-5 h-5" />}
                                    </div>
                                    <div>
                                        <h4 className="text-xs font-black text-white uppercase tracking-widest">{item.user}</h4>
                                        <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-0.5">{item.issue}</p>
                                    </div>
                                </div>
                                <div className="text-right flex items-center space-x-6">
                                    <div className="hidden md:block">
                                        <p className="text-[10px] font-black text-white uppercase tracking-widest">Timestamp</p>
                                        <p className="text-[10px] font-medium text-text-dim mt-0.5">{item.time}</p>
                                    </div>
                                    <button className="px-5 py-2.5 bg-white/5 border border-white/5 rounded-xl text-[10px] font-black text-white uppercase tracking-widest group-hover:bg-white/10 transition-all">
                                        Review
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Financial Summary */}
                <div className="space-y-8">
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <h3 className="text-lg font-black text-white mb-6 uppercase tracking-tight">Platform Financials</h3>
                        <div className="space-y-6">
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Volume (24h)</p>
                                <h4 className="text-3xl font-black text-white">$14,850,290.00</h4>
                            </div>
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-2">Platform Fees Revenue</p>
                                <h4 className="text-3xl font-black text-accent">$242,900.00</h4>
                            </div>
                            <div className="pt-6 border-t border-white/5">
                                <div className="flex items-center justify-between mb-4">
                                    <span className="text-xs font-bold text-text-dim">Withdrawal Queue</span>
                                    <span className="text-xs font-black text-white">42 Pending</span>
                                </div>
                                <div className="w-full bg-white/5 h-1.5 rounded-full overflow-hidden">
                                    <div className="w-1/3 bg-primary h-full rounded-full" />
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 bg-[url('/grid.svg')] bg-center bg-cover">
                        <div className="flex items-center space-x-3 mb-6">
                            <div className="w-10 h-10 rounded-xl bg-primary flex items-center justify-center">
                                <AlertTriangle className="w-5 h-5 text-white" />
                            </div>
                            <h4 className="font-black text-white uppercase tracking-widest text-xs">Major Incident Log</h4>
                        </div>
                        <p className="text-[10px] text-text-dim font-bold mb-6 italic">No critical incidents reported in the last 7 days.</p>
                        <button className="w-full py-4 border border-white/10 rounded-2xl text-[10px] font-black text-white uppercase tracking-[0.2em] hover:bg-white/5 transition-all">
                            Open Status Page
                        </button>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminDashboard;
