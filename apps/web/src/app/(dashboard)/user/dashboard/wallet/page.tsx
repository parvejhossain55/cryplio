"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    Plus,
    ArrowUpRight,
    ArrowDownLeft,
    Settings,
    Shield,
    Eye,
    EyeOff,
    History,
    TrendingUp,
    CreditCard,
    Wallet as WalletIcon,
    ChevronRight
} from "lucide-react";
import { motion } from "framer-motion";

const UserWallet = () => {
    return (
        <DashboardLayout title="Universal Wallet" role="user">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Balanced Section */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Portfolio Glow Card */}
                    <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-10 group">
                        <div className="absolute top-0 right-0 w-[60%] h-[150%] bg-gradient-to-br from-primary/10 via-secondary/5 to-transparent -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

                        <div className="relative z-10">
                            <div className="flex items-center justify-between mb-8">
                                <div className="flex items-center space-x-3">
                                    <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center shadow-2xl shadow-primary/30">
                                        <WalletIcon className="text-white w-6 h-6" />
                                    </div>
                                    <div>
                                        <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Total Balance</p>
                                        <div className="flex items-center space-x-2">
                                            <h2 className="text-4xl font-black text-white">$42,850.24</h2>
                                            <Eye className="w-4 h-4 text-text-dim cursor-pointer hover:text-white transition-colors" />
                                        </div>
                                    </div>
                                </div>
                                <div className="text-right">
                                    <span className="text-accent font-bold bg-accent/10 px-3 py-1 rounded-xl text-xs flex items-center border border-accent/20">
                                        <TrendingUp className="w-3 h-3 mr-1" /> +$1,240.50 (24h)
                                    </span>
                                </div>
                            </div>

                            <div className="flex items-center space-x-4 pt-4 border-t border-white/5">
                                <button className="flex-1 bg-white text-background py-5 rounded-3xl font-black text-sm flex items-center justify-center transition-all hover:scale-[1.02] active:scale-[0.98] shadow-2xl shadow-white/5">
                                    <Plus className="w-5 h-5 mr-2" /> Add Funds
                                </button>
                                <button className="flex-1 bg-surface-light border border-white/10 text-white py-5 rounded-3xl font-black text-sm flex items-center justify-center transition-all hover:bg-white/5 active:scale-[0.98]">
                                    <ArrowUpRight className="w-5 h-5 mr-2" /> Withdraw
                                </button>
                            </div>
                        </div>
                    </div>

                    {/* Assets Grid */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <div className="flex items-center justify-between mb-8">
                            <h3 className="text-xl font-black text-white uppercase tracking-tight">Your Assets</h3>
                            <div className="flex items-center space-x-2">
                                <button className="p-2 hover:bg-white/5 rounded-lg">
                                    <Settings className="w-4 h-4 text-text-dim" />
                                </button>
                            </div>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {[
                                { name: "Bitcoin", symbol: "BTC", balance: "0.4502", value: "$28,450.00", change: "+4.2%", color: "border-orange-500/20" },
                                { name: "Ethereum", symbol: "ETH", balance: "12.50", value: "$12,400.00", change: "-1.5%", color: "border-blue-500/20" },
                                { name: "USDT", symbol: "USDT", balance: "2,000.00", value: "$2,000.00", change: "+0.01%", color: "border-green-500/20" },
                                { name: "Solana", symbol: "SOL", balance: "450.0", value: "$125.00", change: "+12.4%", color: "border-purple-500/20" },
                            ].map((asset) => (
                                <div key={asset.symbol} className={`p-6 rounded-[2rem] bg-white/5 border ${asset.color} hover:bg-white/[0.08] transition-all group cursor-pointer`}>
                                    <div className="flex items-center justify-between mb-4">
                                        <div className="flex items-center space-x-3">
                                            <div className="w-10 h-10 rounded-xl bg-surface flex items-center justify-center font-black text-sm">
                                                {asset.symbol[0]}
                                            </div>
                                            <div>
                                                <p className="text-xs font-black text-white">{asset.name}</p>
                                                <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest">{asset.symbol}</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className={`text-[10px] font-black ${asset.change.startsWith('+') ? 'text-accent' : 'text-primary'}`}>
                                                {asset.change}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="flex items-end justify-between">
                                        <div>
                                            <p className="text-xl font-black text-white">{asset.balance}</p>
                                            <p className="text-[10px] font-medium text-text-dim">{asset.value}</p>
                                        </div>
                                        <ChevronRight className="w-5 h-5 text-text-dim opacity-0 group-hover:opacity-100 transition-all" />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Sidebar: Card & Security */}
                <div className="space-y-8">
                    {/* Physical Card Preview */}
                    <div className="bg-gradient-to-br from-[#121212] to-primary/20 p-8 rounded-[2.5rem] border border-white/10 h-64 relative overflow-hidden group shadow-2xl">
                        <div className="absolute top-0 right-0 p-8 opacity-20">
                            <CreditCard className="w-20 h-20 text-white" />
                        </div>
                        <div className="relative z-10 flex flex-col h-full justify-between">
                            <div className="flex justify-between items-start">
                                <WalletIcon className="w-8 h-8 text-white" />
                                <span className="text-[10px] font-black text-white uppercase tracking-[0.2em] bg-white/10 px-2 py-1 rounded">Titanium</span>
                            </div>
                            <div>
                                <p className="text-xl font-black text-white tracking-widest mb-1">•••• •••• •••• 4820</p>
                                <div className="flex items-center space-x-4">
                                    <div>
                                        <p className="text-[8px] font-black text-text-dim uppercase tracking-widest">Exp</p>
                                        <p className="text-xs font-bold text-white">08/28</p>
                                    </div>
                                    <div>
                                        <p className="text-[8px] font-black text-text-dim uppercase tracking-widest">Holder</p>
                                        <p className="text-xs font-bold text-white uppercase">Alex Morgan</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className="absolute top-0 left-0 w-full h-full bg-primary/5 group-hover:bg-primary/10 transition-colors pointer-events-none" />
                    </div>

                    {/* Quick Tools */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <h3 className="text-lg font-black text-white mb-6 uppercase tracking-tight">Wallet Security</h3>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between p-4 rounded-2xl bg-white/5 border border-white/5">
                                <div className="flex items-center space-x-3">
                                    <Shield className="w-5 h-5 text-accent" />
                                    <span className="text-xs font-bold text-white">Cold Storage</span>
                                </div>
                                <div className="w-10 h-5 bg-accent/20 rounded-full relative">
                                    <div className="absolute right-1 top-1 w-3 h-3 bg-accent rounded-full" />
                                </div>
                            </div>
                            <div className="flex items-center justify-between p-4 rounded-2xl bg-white/5 border border-white/5">
                                <div className="flex items-center space-x-3">
                                    <Shield className="w-5 h-5 text-text-dim" />
                                    <span className="text-xs font-bold text-white">Whitelist Only</span>
                                </div>
                                <div className="w-10 h-5 bg-white/10 rounded-full relative">
                                    <div className="absolute left-1 top-1 w-3 h-3 bg-white/40 rounded-full" />
                                </div>
                            </div>
                        </div>
                        <button className="w-full mt-6 py-4 bg-white/5 border border-white/5 rounded-2xl text-[10px] font-black text-white uppercase tracking-[0.2em] hover:bg-white/10 transition-all">
                            Manage Multi-sig
                        </button>
                    </div>

                    {/* Transaction Feed */}
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-lg font-black text-white uppercase tracking-tight">Recent Activity</h3>
                            <History className="w-4 h-4 text-text-dim" />
                        </div>
                        <div className="space-y-4">
                            {[1, 2, 3].map((i) => (
                                <div key={i} className="flex items-center justify-between group cursor-pointer">
                                    <div className="flex items-center space-x-3">
                                        <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center">
                                            <ArrowDownLeft className="w-4 h-4 text-accent" />
                                        </div>
                                        <div>
                                            <p className="text-[10px] font-black text-white uppercase">Received BTC</p>
                                            <p className="text-[8px] text-text-dim">10 min ago</p>
                                        </div>
                                    </div>
                                    <p className="text-xs font-black text-accent">+$120.00</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserWallet;
