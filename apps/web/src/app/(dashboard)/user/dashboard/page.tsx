"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    TrendingUp,
    Wallet,
    ArrowUpRight,
    ArrowDownLeft,
    CreditCard,
    Clock,
    ExternalLink,
    ChevronRight,
    ArrowLeftRight,
    Store,
    ShieldCheck
} from "lucide-react";
import { motion } from "framer-motion";

const UserDashboard = () => {
    return (
        <DashboardLayout title="Overview" role="user">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column: Stats & Wallet */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Main Wallet Card */}
                    <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-8 group">
                        <div className="absolute top-0 right-0 w-[50%] h-[150%] bg-primary/10 -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

                        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                            <div>
                                <p className="text-xs font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Balance</p>
                                <div className="flex items-baseline space-x-3">
                                    <h2 className="text-5xl font-black text-white">$42,850.24</h2>
                                    <span className="text-accent font-bold bg-accent/10 px-2 py-0.5 rounded-lg text-sm flex items-center">
                                        <TrendingUp className="w-3 h-3 mr-1" /> +12.5%
                                    </span>
                                </div>
                            </div>

                            <div className="flex items-center space-x-3">
                                <button className="flex-1 md:flex-none px-6 py-4 bg-white text-background rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:scale-105 active:scale-95">
                                    <ArrowUpRight className="w-4 h-4 mr-2" /> Send
                                </button>
                                <button className="flex-1 md:flex-none px-6 py-4 bg-surface-light border border-white/5 text-white rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:bg-white/5 active:scale-95">
                                    <ArrowDownLeft className="w-4 h-4 mr-2" /> Receive
                                </button>
                            </div>
                        </div>

                        <div className="relative z-10 mt-10 grid grid-cols-2 sm:grid-cols-4 gap-4">
                            {[
                                { name: "Bitcoin", symbol: "BTC", amount: "0.45", price: "$28,450", icon: "₿" },
                                { name: "Ethereum", symbol: "ETH", amount: "12.5", price: "$12,400", icon: "Ξ" },
                                { name: "USDT", symbol: "USDT", amount: "2,000", price: "$2,000", icon: "₮" },
                                { name: "Solana", symbol: "SOL", amount: "450.0", price: "$125", icon: "S" },
                            ].map((coin) => (
                                <div key={coin.symbol} className="p-4 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 transition-all cursor-pointer group/card">
                                    <div className="w-10 h-10 rounded-2xl bg-surface flex items-center justify-center text-lg font-black mb-3 group-hover/card:scale-110 transition-transform">
                                        {coin.icon}
                                    </div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{coin.symbol}</p>
                                    <p className="text-sm font-bold text-white mt-1">{coin.amount}</p>
                                    <p className="text-[10px] font-medium text-text-dim mt-0.5">{coin.price}</p>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Recent Activity */}
                    <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                        <div className="flex items-center justify-between mb-8">
                            <h3 className="text-xl font-black text-white">Recent Activity</h3>
                            <button className="text-xs font-bold text-primary hover:text-white transition-colors flex items-center">
                                View All <ChevronRight className="w-4 h-4 ml-1" />
                            </button>
                        </div>

                        <div className="space-y-4">
                            {[
                                { type: "buy", asset: "USDT", amount: "500", status: "completed", date: "2 mins ago", price: "$1.00" },
                                { type: "sell", asset: "BTC", amount: "0.02", status: "pending", date: "15 mins ago", price: "$64,200" },
                                { type: "swap", asset: "ETH to USDT", amount: "1.2", status: "completed", date: "2 hours ago", price: "$3,200" },
                                { type: "deposit", asset: "USD", amount: "1,200", status: "failed", date: "1 day ago", price: "-" },
                            ].map((tx, i) => (
                                <div key={i} className="flex items-center justify-between p-4 rounded-2xl hover:bg-white/5 transition-all group">
                                    <div className="flex items-center space-x-4">
                                        <div className={`w-12 h-12 rounded-2xl flex items-center justify-center border border-white/5 ${tx.type === 'buy' ? 'bg-accent/10 text-accent' :
                                            tx.type === 'sell' ? 'bg-primary/10 text-primary' :
                                                'bg-secondary/10 text-secondary'
                                            }`}>
                                            {tx.type === 'buy' ? <ArrowDownLeft className="w-6 h-6" /> :
                                                tx.type === 'sell' ? <ArrowUpRight className="w-6 h-6" /> :
                                                    <ArrowLeftRight className="w-5 h-5" />}
                                        </div>
                                        <div>
                                            <h4 className="font-bold text-white tracking-tight flex items-center uppercase">
                                                {tx.type} {tx.asset}
                                                <span className={`ml-2 text-[8px] px-1.5 py-0.5 rounded font-black uppercase tracking-widest ${tx.status === 'completed' ? 'bg-accent/10 text-accent' :
                                                    tx.status === 'pending' ? 'bg-primary/10 text-primary' :
                                                        'bg-red-500/10 text-red-500'
                                                    }`}>
                                                    {tx.status}
                                                </span>
                                            </h4>
                                            <p className="text-[10px] font-medium text-text-dim mt-1">{tx.date} • Price: {tx.price}</p>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <p className="font-black text-white">{tx.amount} {tx.asset.split(' ')[0]}</p>
                                        <ExternalLink className="w-3 h-3 ml-auto mt-1 text-text-dim opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer hover:text-white" />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Right Column: Cards & Offers */}
                <div className="space-y-8">
                    {/* Cryplio Card Preview */}
                    <div className="bg-gradient-to-br from-primary to-secondary p-8 rounded-[2.5rem] relative overflow-hidden group shadow-2xl shadow-primary/20 h-64">
                        <div className="absolute top-0 right-0 p-8">
                            <Wallet className="w-10 h-10 text-white/50 group-hover:rotate-12 transition-transform duration-500" />
                        </div>
                        <div className="absolute bottom-0 left-0 p-8 w-full">
                            <div className="flex justify-between items-end">
                                <div>
                                    <p className="text-[10px] font-black text-white/60 uppercase tracking-[0.2em] mb-1">Virtual Card</p>
                                    <p className="text-2xl font-black text-white tracking-widest">•••• 4820</p>
                                </div>
                                <div className="text-right">
                                    <p className="text-[10px] font-black text-white/60 uppercase tracking-[0.2em] mb-1">Exp</p>
                                    <p className="text-lg font-black text-white italic">08/28</p>
                                </div>
                            </div>
                        </div>
                        {/* Card visual elements */}
                        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[150%] h-[150%] border border-white/20 rounded-full opacity-20 group-hover:scale-110 transition-transform duration-700" />
                    </div>

                    {/* Quick Actions */}
                    <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                        <h3 className="text-lg font-black text-white mb-6">Quick Actions</h3>
                        <div className="grid grid-cols-2 gap-4">
                            {[
                                { name: "P2P Marketplace", icon: Store, color: "text-primary", bg: "bg-primary/10" },
                                { name: "Instant Swap", icon: ArrowLeftRight, color: "text-secondary", bg: "bg-secondary/10" },
                                { name: "Card Services", icon: CreditCard, color: "text-accent", bg: "bg-accent/10" },
                                { name: "Refer & Earn", icon: TrendingUp, color: "text-primary", bg: "bg-primary/10" },
                            ].map((action, i) => (
                                <button key={i} className="flex flex-col items-center justify-center p-6 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 hover:bg-white/[0.08] transition-all group">
                                    <div className={`w-12 h-12 rounded-2xl ${action.bg} ${action.color} flex items-center justify-center mb-3 group-hover:scale-110 transition-transform`}>
                                        <action.icon className="w-6 h-6" />
                                    </div>
                                    <span className="text-[10px] font-black text-white text-center uppercase tracking-widest leading-tight">{action.name}</span>
                                </button>
                            ))}
                        </div>
                    </div>

                    {/* Security Status */}
                    <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                        <div className="flex items-center space-x-4 mb-6">
                            <div className="w-12 h-12 rounded-2xl bg-accent/20 flex items-center justify-center">
                                <ShieldCheck className="w-6 h-6 text-accent" />
                            </div>
                            <div>
                                <h4 className="font-black text-white uppercase text-xs tracking-widest">Account Health</h4>
                                <p className="text-[10px] font-bold text-accent uppercase tracking-widest">Excellent</p>
                            </div>
                        </div>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-text-dim font-medium">Identity Verification</span>
                                <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Level 2</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-text-dim font-medium">2FA Status</span>
                                <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Enabled</span>
                            </div>
                            <div className="flex items-center justify-between">
                                <span className="text-xs text-text-dim font-medium">Login Notifications</span>
                                <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Active</span>
                            </div>
                        </div>
                        <button className="w-full mt-8 py-4 bg-white/5 border border-white/5 rounded-2xl text-[10px] font-black text-white uppercase tracking-[0.2em] hover:bg-white/10 transition-all">
                            Security Hub
                        </button>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserDashboard;
