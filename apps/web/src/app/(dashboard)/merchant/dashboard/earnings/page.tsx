"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    BarChart3,
    TrendingUp,
    ArrowUpRight,
    ArrowDownLeft,
    DollarSign,
    Wallet,
    Calendar,
    ChevronRight,
    Target,
    Zap,
    ArrowRight
} from "lucide-react";
import { motion } from "framer-motion";

const MerchantEarnings = () => {
    return (
        <DashboardLayout title="Financial Protocol" role="merchant">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Stats Column */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Revenue Overview Card */}
                    <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-10 group">
                        <div className="absolute top-0 right-0 w-full h-full bg-[url('/grid.svg')] bg-center opacity-5 pointer-events-none" />
                        <div className="absolute top-0 right-0 w-[50%] h-[150%] bg-primary/10 -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

                        <div className="relative z-10">
                            <div className="flex items-center justify-between mb-10">
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Managed Volume</p>
                                    <h2 className="text-5xl font-black text-white tracking-tighter">$1,240,850.24</h2>
                                </div>
                                <div className="p-4 bg-primary/10 rounded-[2rem] border border-primary/20 group-hover:scale-110 transition-transform">
                                    <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center shadow-xl shadow-primary/30">
                                        <TrendingUp className="text-white w-6 h-6" />
                                    </div>
                                </div>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 pt-8 border-t border-white/5">
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1">Gross Profit</p>
                                    <p className="text-xl font-black text-accent">$42,290.00</p>
                                </div>
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1">Fee Exposure</p>
                                    <p className="text-xl font-black text-white">$12,400.00</p>
                                </div>
                                <div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1">Net Earnings</p>
                                    <p className="text-xl font-black text-secondary">$29,890.00</p>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Chart Replacement / Volume Analysis */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <div className="flex items-center justify-between mb-10">
                            <div>
                                <h3 className="text-xl font-black text-white uppercase tracking-tight">Performance Stream</h3>
                                <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Daily trading activity analytics</p>
                            </div>
                            <div className="flex items-center bg-white/5 p-1 rounded-xl border border-white/5">
                                <button className="px-4 py-2 text-[10px] font-black uppercase tracking-widest rounded-lg bg-surface text-white border border-white/10">1W</button>
                                <button className="px-4 py-2 text-[10px] font-black uppercase tracking-widest rounded-lg text-text-dim hover:text-white transition-colors">1M</button>
                                <button className="px-4 py-2 text-[10px] font-black uppercase tracking-widest rounded-lg text-text-dim hover:text-white transition-colors">1Y</button>
                            </div>
                        </div>

                        <div className="space-y-6">
                            {[
                                { day: "Monday", volume: "$142k", profit: "+$4,200", percent: 85 },
                                { day: "Tuesday", volume: "$98k", profit: "+$2,800", percent: 60 },
                                { day: "Wednesday", volume: "$210k", profit: "+$6,100", percent: 100 },
                                { day: "Thursday", volume: "$125k", profit: "+$3,900", percent: 75 },
                                { day: "Friday (Today)", volume: "$45k", profit: "+$1,400", percent: 35 },
                            ].map((item, i) => (
                                <div key={i} className="group">
                                    <div className="flex items-center justify-between mb-3 px-1">
                                        <div className="flex items-center space-x-3">
                                            <span className="text-[10px] font-black text-text-dim uppercase tracking-widest">{item.day}</span>
                                            <span className="text-[10px] font-bold text-white uppercase">{item.volume} Volume</span>
                                        </div>
                                        <span className="text-[10px] font-black text-accent tracking-widest">{item.profit} Profit</span>
                                    </div>
                                    <div className="w-full bg-white/5 h-2 rounded-full overflow-hidden border border-white/5">
                                        <div
                                            className="h-full bg-gradient-to-r from-primary to-secondary rounded-full group-hover:opacity-80 transition-opacity"
                                            style={{ width: `${item.percent}%` }}
                                        />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Right Sidebar */}
                <div className="space-y-8">
                    {/* Settlement Card */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <div className="w-12 h-12 bg-secondary/20 rounded-2xl flex items-center justify-center mb-6">
                            <Wallet className="text-secondary w-6 h-6" />
                        </div>
                        <h3 className="text-lg font-black text-white uppercase tracking-tight mb-2">Payout Protocol</h3>
                        <p className="text-xs text-text-dim font-medium mb-8 leading-relaxed">
                            Earnings are automatically processed every 24 hours. Minimal threshold is $100.
                        </p>
                        <div className="space-y-4 mb-8">
                            <div className="flex items-center justify-between p-4 rounded-3xl bg-white/5 border border-white/5">
                                <span className="text-[10px] font-black text-text-dim uppercase tracking-widest">Next Payout</span>
                                <span className="text-xs font-black text-white tracking-widest">4h 20m</span>
                            </div>
                            <div className="flex items-center justify-between p-4 rounded-3xl bg-white/5 border border-white/5">
                                <span className="text-[10px] font-black text-text-dim uppercase tracking-widest">Amount</span>
                                <span className="text-xs font-black text-accent tracking-widest">$4,250.00</span>
                            </div>
                        </div>
                        <button className="w-full py-5 bg-white text-background rounded-3xl font-black text-sm uppercase tracking-[0.2em] hover:scale-[1.02] active:scale-[0.98] transition-all flex items-center justify-center">
                            Fast Release <ArrowRight className="w-4 h-4 ml-2" />
                        </button>
                    </div>

                    {/* Goals / Targets */}
                    <div className="bg-gradient-to-br from-primary/20 to-transparent border border-primary/30 rounded-[2.5rem] p-8">
                        <div className="flex items-center space-x-3 mb-6">
                            <Target className="w-6 h-6 text-primary" />
                            <h4 className="text-xs font-black text-white uppercase tracking-widest">Merchant Goal</h4>
                        </div>
                        <div className="mb-6">
                            <div className="flex items-baseline justify-between mb-2 px-1">
                                <h5 className="text-2xl font-black text-white">$1M</h5>
                                <span className="text-[10px] font-black text-primary uppercase tracking-widest">Threshold</span>
                            </div>
                            <p className="text-[10px] text-text-dim font-medium uppercase tracking-[0.15em]">Reach $1M monthly volume for Zero Fees.</p>
                        </div>
                        <div className="w-full bg-white/10 h-3 rounded-full overflow-hidden border border-white/10 mb-2">
                            <div className="w-[65%] bg-primary h-full rounded-full shadow-[0_0_20px_rgba(99,102,241,0.5)]" />
                        </div>
                        <p className="text-[10px] font-black text-primary text-center uppercase tracking-widest mt-4">65% of Target Achieved</p>
                    </div>

                    {/* Quick Analytics Feed */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <div className="flex items-center justify-between mb-8">
                            <h3 className="text-lg font-black text-white uppercase tracking-tight">Fee Reports</h3>
                            <button className="text-[10px] font-black text-primary uppercase hover:underline tracking-widest">All</button>
                        </div>
                        <div className="space-y-6">
                            {[
                                { label: "Network Gas", amount: "$1,240", icon: Zap, color: "text-secondary" },
                                { label: "Escrow Fee", amount: "$8,500", icon: Shield, color: "text-accent" },
                                { label: "Conversion Fee", amount: "$2,660", icon: ArrowLeftRight, color: "text-primary" },
                            ].map((fee, i) => (
                                <div key={i} className="flex items-center justify-between group cursor-pointer hover:translate-x-1 transition-transform">
                                    <div className="flex items-center space-x-4">
                                        <div className={`w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center ${fee.color}`}>
                                            <fee.icon className="w-4 h-4" />
                                        </div>
                                        <span className="text-[10px] font-black text-white uppercase tracking-widest">{fee.label}</span>
                                    </div>
                                    <span className="text-[10px] font-black text-white/50">{fee.amount}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

const Shield = ({ className }: { className?: string }) => (
    <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10" /><path d="m9 12 2 2 4-4" /></svg>
);

const ArrowLeftRight = ({ className }: { className?: string }) => (
    <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="m2 17 8-8 4 4 8-8" /><path d="M18 7h4v4" /></svg>
)

export default MerchantEarnings;
