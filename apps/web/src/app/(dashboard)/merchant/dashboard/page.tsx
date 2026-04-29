"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    BarChart3,
    TrendingUp,
    Users,
    Store,
    ArrowUpRight,
    ArrowDownLeft,
    ChevronRight,
    Plus,
    Filter,
    Search
} from "lucide-react";
import { motion } from "framer-motion";

const MerchantDashboard = () => {
    return (
        <DashboardLayout title="Merchant Hub" role="merchant">
            {/* Top Header Stats */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {[
                    { label: "Total Revenue", value: "$124,500.00", change: "+18.2%", icon: BarChart3, color: "text-primary" },
                    { label: "Active Orders", value: "48", change: "+4", icon: Store, color: "text-secondary" },
                    { label: "Completed P2P", value: "1,240", change: "+124", icon: TrendingUp, color: "text-accent" },
                    { label: "Customer Satisfaction", value: "99.8%", change: "High", icon: Users, color: "text-primary" },
                ].map((stat, i) => (
                    <div key={i} className="bg-surface border border-white/10 p-6 rounded-[2rem] hover:border-white/20 transition-all group">
                        <div className="flex items-center justify-between mb-4">
                            <div className={`p-3 rounded-2xl bg-white/5 ${stat.color} group-hover:scale-110 transition-transform`}>
                                <stat.icon className="w-5 h-5" />
                            </div>
                            <span className={`text-[10px] font-black uppercase tracking-widest ${stat.change.startsWith('+') ? 'text-accent' : 'text-primary'}`}>
                                {stat.change}
                            </span>
                        </div>
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-1">{stat.label}</p>
                        <h3 className="text-2xl font-black text-white">{stat.value}</h3>
                    </div>
                ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Main Orders Section */}
                <div className="lg:col-span-2 space-y-8">
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] overflow-hidden">
                        <div className="p-8 border-b border-white/5 flex flex-col md:flex-row md:items-center justify-between gap-4">
                            <div>
                                <h3 className="text-xl font-black text-white">Active Order Stream</h3>
                                <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Real-time p2p transactions</p>
                            </div>
                            <div className="flex items-center space-x-2">
                                <button className="p-3 bg-white/5 border border-white/5 rounded-xl hover:bg-white/10 transition-all">
                                    <Filter className="w-4 h-4 text-text-dim" />
                                </button>
                                <button className="px-6 py-3 bg-primary text-white rounded-xl font-black text-xs uppercase tracking-widest flex items-center hover:scale-105 active:scale-95 transition-all">
                                    <Plus className="w-4 h-4 mr-2" /> New Ad
                                </button>
                            </div>
                        </div>

                        <div className="overflow-x-auto">
                            <table className="w-full text-left">
                                <thead>
                                    <tr className="border-b border-white/5">
                                        <th className="px-8 py-4 text-[10px] font-black text-text-dim uppercase tracking-widest">Customer</th>
                                        <th className="px-8 py-4 text-[10px] font-black text-text-dim uppercase tracking-widest">Type / Asset</th>
                                        <th className="px-8 py-4 text-[10px] font-black text-text-dim uppercase tracking-widest">Amount</th>
                                        <th className="px-8 py-4 text-[10px] font-black text-text-dim uppercase tracking-widest">Status</th>
                                        <th className="px-8 py-4"></th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-white/5">
                                    {[
                                        { user: "User_4920", type: "BUY", asset: "USDT", amount: "$500.00", status: "In Escrow", time: "2m ago" },
                                        { user: "CryptoKing", type: "SELL", asset: "BTC", amount: "0.45 BTC", status: "Confirming", time: "5m ago" },
                                        { user: "WhaleTraders", type: "BUY", asset: "ETH", amount: "1.2 ETH", status: "Pending", time: "12m ago" },
                                        { user: "TraderMax", type: "BUY", asset: "USDT", amount: "$1,200.00", status: "In Escrow", time: "18m ago" },
                                        { user: "MoonShot", type: "SELL", asset: "SOL", amount: "50 SOL", status: "Completed", time: "45m ago" },
                                    ].map((order, i) => (
                                        <tr key={i} className="hover:bg-white/[0.02] transition-all group">
                                            <td className="px-8 py-6">
                                                <div className="flex items-center space-x-3">
                                                    <div className="w-8 h-8 rounded-lg bg-white/5 border border-white/5 flex items-center justify-center font-bold text-xs uppercase">
                                                        {order.user.charAt(0)}
                                                    </div>
                                                    <div>
                                                        <p className="text-xs font-black text-white">{order.user}</p>
                                                        <p className="text-[10px] text-text-dim font-medium">{order.time}</p>
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-8 py-6">
                                                <div className="flex items-center space-x-2">
                                                    <span className={`text-[10px] font-black px-2 py-0.5 rounded ${order.type === 'BUY' ? 'bg-accent/10 text-accent' : 'bg-primary/10 text-primary'}`}>
                                                        {order.type}
                                                    </span>
                                                    <span className="text-xs font-bold text-white">{order.asset}</span>
                                                </div>
                                            </td>
                                            <td className="px-8 py-6">
                                                <p className="text-xs font-black text-white">{order.amount}</p>
                                            </td>
                                            <td className="px-8 py-6">
                                                <div className="flex items-center space-x-2">
                                                    <div className={`w-1.5 h-1.5 rounded-full ${order.status === 'In Escrow' ? 'bg-secondary' : 'bg-primary'} animate-pulse`} />
                                                    <span className="text-[10px] font-bold text-white uppercase tracking-widest">{order.status}</span>
                                                </div>
                                            </td>
                                            <td className="px-8 py-6 text-right">
                                                <button className="p-2 hover:bg-white/5 rounded-lg transition-all opacity-0 group-hover:opacity-100">
                                                    <ChevronRight className="w-4 h-4 text-text-dim" />
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                {/* Right Sidebar: Analytics & Ad Performance */}
                <div className="space-y-8">
                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8">
                        <h3 className="text-lg font-black text-white mb-6">Ad Performance</h3>
                        <div className="space-y-6">
                            {[
                                { label: "USDT Sales (High Liquidity)", views: "2.4k", conversion: "12%", status: "Active" },
                                { label: "BTC/Cash (Premium Layer)", views: "850", conversion: "4%", status: "Paused" },
                                { label: "ETH Fast Swap", views: "1.2k", conversion: "18%", status: "Active" },
                            ].map((ad, i) => (
                                <div key={i} className="p-4 rounded-2xl bg-white/5 border border-white/5">
                                    <div className="flex items-center justify-between mb-3">
                                        <span className="text-[10px] font-black text-white uppercase tracking-widest">{ad.label}</span>
                                        <span className={`text-[8px] font-black uppercase tracking-tighter px-1.5 py-0.5 rounded ${ad.status === 'Active' ? 'bg-accent/10 text-accent' : 'bg-white/10 text-text-dim'}`}>
                                            {ad.status}
                                        </span>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="text-[8px] font-bold text-text-dim uppercase tracking-widest">Views</p>
                                            <p className="text-sm font-black text-white">{ad.views}</p>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-[8px] font-bold text-text-dim uppercase tracking-widest">Conv.</p>
                                            <p className="text-sm font-black text-accent">{ad.conversion}</p>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="bg-gradient-to-br from-primary/10 to-transparent border border-primary/20 rounded-[2.5rem] p-8">
                        <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center mb-6 shadow-xl shadow-primary/20">
                            <TrendingUp className="text-white w-6 h-6" />
                        </div>
                        <h3 className="text-lg font-black text-white mb-2">Merchant Level 4</h3>
                        <p className="text-xs text-text-dim font-medium mb-6 leading-relaxed">
                            You're in the top 5% of merchants. Unlock lower escrow fees by completing 10 more trades.
                        </p>
                        <div className="w-full bg-white/5 h-2 rounded-full overflow-hidden mb-2">
                            <div className="w-[85%] bg-primary h-full rounded-full" />
                        </div>
                        <div className="flex justify-between items-center text-[10px] font-black uppercase tracking-widest text-text-dim">
                            <span>Level 4</span>
                            <span>Level 5</span>
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default MerchantDashboard;
