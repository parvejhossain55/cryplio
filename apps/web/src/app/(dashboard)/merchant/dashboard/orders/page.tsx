"use client";

import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ArrowLeftRight,
    MessageSquare,
    Clock,
    CheckCircle2,
    AlertCircle,
    ExternalLink,
    ChevronRight,
    Filter,
    Search,
    Zap
} from "lucide-react";
import { motion } from "framer-motion";

const MerchantOrders = () => {
    const orders = [
        { id: "ORD_7720", user: "CryptoLover_99", type: "BUY", asset: "USDT", amount: "1,200.00 USDT", fiat: "$1,200.00", status: "In Escrow", time: "2m ago", platformFee: "$12.00" },
        { id: "ORD_7721", user: "WolfOfWallSt", type: "SELL", asset: "BTC", amount: "0.24 BTC", fiat: "$15,400.00", status: "Confirming", time: "15m ago", platformFee: "$154.00" },
        { id: "ORD_7722", user: "MoonWalker", type: "BUY", asset: "ETH", amount: "4.5 ETH", status: "Completed", time: "2h ago", platformFee: "$45.00" },
        { id: "ORD_7723", user: "BitKing", type: "BUY", asset: "USDT", amount: "5,000.00 USDT", status: "Disputed", time: "1d ago", platformFee: "$50.00" },
        { id: "ORD_7724", user: "Satoshi_N", type: "SELL", asset: "SOL", amount: "100.00 SOL", status: "Completed", time: "2d ago", platformFee: "$10.00" },
    ];

    return (
        <DashboardLayout title="Client Orders" role="merchant">
            <div className="space-y-8">
                {/* Stats Summary */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    {[
                        { label: "Active Orders", value: "12", icon: Zap, color: "text-primary" },
                        { label: "Completion Rate", value: "99.2%", icon: CheckCircle2, color: "text-accent" },
                        { label: "Avg. Resolution", value: "8.4m", icon: Clock, color: "text-secondary" },
                        { label: "Disputed Orders", value: "1", icon: AlertCircle, color: "text-primary" },
                    ].map((stat, i) => (
                        <div key={i} className="bg-surface border border-white/10 p-6 rounded-3xl hover:border-white/20 transition-all">
                            <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-3">{stat.label}</p>
                            <div className="flex items-center justify-between">
                                <h3 className="text-2xl font-black text-white">{stat.value}</h3>
                                <stat.icon className={`w-5 h-5 ${stat.color}`} />
                            </div>
                        </div>
                    ))}
                </div>

                {/* Table Section */}
                <div className="bg-surface border border-white/10 rounded-[2.5rem] overflow-hidden">
                    <div className="p-8 border-b border-white/5 flex flex-col md:flex-row md:items-center justify-between gap-4">
                        <div>
                            <h3 className="text-xl font-black text-white uppercase tracking-tight">Real-time Order Stream</h3>
                            <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">Direct peer-to-peer interactions</p>
                        </div>
                        <div className="flex items-center space-x-3">
                            <div className="relative">
                                <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-4 h-4" />
                                <input
                                    type="text"
                                    placeholder="Search Orders..."
                                    className="bg-white/5 border border-white/10 py-3 pl-12 pr-6 rounded-xl text-xs w-64 outline-none focus:border-primary/50 transition-all"
                                />
                            </div>
                            <button className="p-3 bg-white/5 border border-white/10 rounded-xl hover:bg-white/15 transition-all">
                                <Filter className="w-4 h-4 text-white" />
                            </button>
                        </div>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full text-left">
                            <thead>
                                <tr className="border-b border-white/5">
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Customer / ID</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Asset / Type</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Order Amount</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Status</th>
                                    <th className="px-8 py-5 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {orders.map((order) => (
                                    <tr key={order.id} className="hover:bg-white/[0.02] transition-all group">
                                        <td className="px-8 py-6">
                                            <div className="flex items-center space-x-4">
                                                <div className="w-10 h-10 rounded-xl bg-white/5 border border-white/5 flex items-center justify-center font-black text-xs text-primary">
                                                    {order.user.charAt(0)}
                                                </div>
                                                <div>
                                                    <p className="text-xs font-black text-white">{order.user}</p>
                                                    <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-0.5">{order.id} • {order.time}</p>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-8 py-6">
                                            <div className="flex items-center space-x-2">
                                                <span className={`text-[9px] font-black px-2 py-0.5 rounded ${order.type === 'BUY' ? 'bg-accent/10 text-accent' : 'bg-primary/10 text-primary'}`}>
                                                    {order.type}
                                                </span>
                                                <span className="text-xs font-bold text-white uppercase tracking-widest">{order.asset}</span>
                                            </div>
                                        </td>
                                        <td className="px-8 py-6">
                                            <p className="text-xs font-black text-white">{order.amount}</p>
                                            <p className="text-[10px] font-medium text-text-dim mt-0.5">Fee: {order.platformFee}</p>
                                        </td>
                                        <td className="px-8 py-6">
                                            <div className="flex items-center space-x-2">
                                                <div className={`w-1.5 h-1.5 rounded-full ${order.status === 'In Escrow' ? 'bg-secondary' :
                                                        order.status === 'Completed' ? 'bg-accent' :
                                                            'bg-primary'
                                                    } ${order.status !== 'Completed' ? 'animate-pulse' : ''}`} />
                                                <span className="text-[10px] font-black text-white uppercase tracking-widest">{order.status}</span>
                                            </div>
                                        </td>
                                        <td className="px-8 py-6 text-right">
                                            <div className="flex items-center space-x-2 justify-end">
                                                <button className="p-2.5 bg-primary/10 border border-primary/20 rounded-xl hover:bg-primary/20 transition-all text-primary">
                                                    <MessageSquare className="w-4 h-4" />
                                                </button>
                                                <button className="p-2.5 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-all text-white">
                                                    <ChevronRight className="w-4 h-4" />
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Dispute Alert Section */}
                <div className="bg-red-500/5 border border-red-500/20 rounded-[2.5rem] p-8 flex items-center justify-between">
                    <div className="flex items-center space-x-6">
                        <div className="w-12 h-12 bg-red-500/20 rounded-2xl flex items-center justify-center">
                            <AlertCircle className="w-6 h-6 text-red-500" />
                        </div>
                        <div>
                            <h4 className="font-black text-white uppercase tracking-widest leading-tight">1 Active Dispute Detected</h4>
                            <p className="text-xs text-text-dim font-medium mt-1">Resolution Protocol #T723 requires your immediate attention.</p>
                        </div>
                    </div>
                    <button className="bg-red-500 text-white px-8 py-3 rounded-xl font-black uppercase text-[10px] tracking-widest hover:scale-105 transition-all">
                        Resolve Now
                    </button>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default MerchantOrders;
