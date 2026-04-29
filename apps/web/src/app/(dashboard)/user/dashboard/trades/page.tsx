"use client";

import React, { useState } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    Search,
    Filter,
    ArrowDownLeft,
    ArrowUpRight,
    ArrowLeftRight,
    Clock,
    ChevronRight,
    ExternalLink,
    Calendar,
    MoreVertical
} from "lucide-react";
import { motion } from "framer-motion";

const UserTrades = () => {
    const [filter, setFilter] = useState("all");

    const trades = [
        { id: "T3920", type: "buy", asset: "USDT", amount: "500", fiat: "$500.00", status: "completed", date: "Apr 26, 2026", time: "14:20" },
        { id: "T3921", type: "sell", asset: "BTC", amount: "0.024", fiat: "$1,540.00", status: "pending", date: "Apr 25, 2026", time: "09:45" },
        { id: "T3922", type: "swap", asset: "ETH/USDT", amount: "1.2", fiat: "$3,840.00", status: "completed", date: "Apr 24, 2026", time: "11:15" },
        { id: "T3923", type: "buy", asset: "USDT", amount: "1,200", fiat: "$1,200.00", status: "failed", date: "Apr 23, 2026", time: "18:30" },
        { id: "T3924", type: "deposit", asset: "USD", amount: "2,000", fiat: "$2,000.00", status: "completed", date: "Apr 22, 2026", time: "10:00" },
    ];

    return (
        <DashboardLayout title="My Trades" role="user">
            <div className="space-y-8">
                {/* Header Actions */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div className="relative flex-1 max-w-md">
                        <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-4 h-4" />
                        <input
                            type="text"
                            placeholder="Search by Order ID or Asset..."
                            className="w-full bg-surface border border-white/5 py-3 pl-12 pr-6 rounded-2xl text-sm outline-none focus:border-primary/50 transition-all font-medium"
                        />
                    </div>
                    <div className="flex items-center space-x-3">
                        <button className="flex items-center space-x-2 px-5 py-3 bg-surface border border-white/5 rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/5 transition-all">
                            <Calendar className="w-4 h-4 text-text-dim" />
                            <span>Date Range</span>
                        </button>
                        <button className="flex items-center space-x-2 px-5 py-3 bg-surface border border-white/5 rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/5 transition-all">
                            <Filter className="w-4 h-4 text-text-dim" />
                            <span>Filters</span>
                        </button>
                    </div>
                </div>

                {/* Filters Tabs */}
                <div className="flex items-center space-x-2 overflow-x-auto pb-4 scrollbar-hide">
                    {["All Trades", "Buying", "Selling", "Swaps", "Completed", "Pending"].map((tab) => (
                        <button
                            key={tab}
                            onClick={() => setFilter(tab.toLowerCase())}
                            className={`px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest whitespace-nowrap transition-all border ${(tab.toLowerCase().includes(filter) || (tab === "All Trades" && filter === "all"))
                                    ? "bg-primary text-white border-primary shadow-lg shadow-primary/20"
                                    : "bg-white/5 text-text-dim border-transparent hover:border-white/10"
                                }`}
                        >
                            {tab}
                        </button>
                    ))}
                </div>

                {/* Trades List */}
                <div className="bg-surface border border-white/10 rounded-[2.5rem] overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead>
                                <tr className="border-b border-white/5">
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Transaction</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Asset / Pair</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Amount</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Price / Fiat</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Status</th>
                                    <th className="px-8 py-5"></th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {trades.map((trade) => (
                                    <tr key={trade.id} className="hover:bg-white/[0.02] transition-all group">
                                        <td className="px-8 py-6">
                                            <div className="flex items-center space-x-4">
                                                <div className={`w-10 h-10 rounded-xl flex items-center justify-center border border-white/5 ${trade.type === 'buy' ? 'bg-accent/10 text-accent' :
                                                        trade.type === 'sell' ? 'bg-primary/10 text-primary' :
                                                            'bg-secondary/10 text-secondary'
                                                    }`}>
                                                    {trade.type === 'buy' ? <ArrowDownLeft className="w-5 h-5" /> :
                                                        trade.type === 'sell' ? <ArrowUpRight className="w-5 h-5" /> :
                                                            <ArrowLeftRight className="w-5 h-5" />}
                                                </div>
                                                <div>
                                                    <p className="text-xs font-black text-white uppercase tracking-widest">{trade.type}</p>
                                                    <p className="text-[10px] font-medium text-text-dim mt-0.5">{trade.date} • {trade.time}</p>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-8 py-6">
                                            <p className="text-xs font-black text-white">{trade.asset}</p>
                                            <p className="text-[10px] font-medium text-text-dim mt-0.5">Order ID: {trade.id}</p>
                                        </td>
                                        <td className="px-8 py-6">
                                            <p className="text-xs font-black text-white">{trade.amount} {trade.asset.split('/')[0]}</p>
                                        </td>
                                        <td className="px-8 py-6">
                                            <p className="text-xs font-black text-white">{trade.fiat}</p>
                                        </td>
                                        <td className="px-8 py-6">
                                            <span className={`text-[9px] font-black px-2.5 py-1 rounded-lg uppercase tracking-widest ${trade.status === 'completed' ? 'bg-accent/10 text-accent border border-accent/20' :
                                                    trade.status === 'pending' ? 'bg-primary/10 text-primary border border-primary/20 animate-pulse' :
                                                        'bg-red-500/10 text-red-500 border border-red-500/20'
                                                }`}>
                                                {trade.status}
                                            </span>
                                        </td>
                                        <td className="px-8 py-6 text-right">
                                            <div className="flex items-center justify-end space-x-2">
                                                <button className="p-2.5 hover:bg-white/5 rounded-xl transition-all text-text-dim hover:text-white">
                                                    <ExternalLink className="w-4 h-4" />
                                                </button>
                                                <button className="p-2.5 hover:bg-white/5 rounded-xl transition-all text-text-dim hover:text-white">
                                                    <MoreVertical className="w-4 h-4" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Footer Message */}
                <div className="flex flex-col items-center justify-center py-10 opacity-40">
                    <Clock className="w-10 h-10 mb-4" />
                    <p className="text-sm font-bold">Only trades from the last 90 days are shown.</p>
                    <button className="text-[10px] font-black text-primary uppercase tracking-[0.2em] mt-2 hover:underline">Request Archive Extract</button>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserTrades;
