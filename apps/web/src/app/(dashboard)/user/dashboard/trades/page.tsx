"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Clock,
    CheckCircle2,
    XCircle,
    AlertTriangle,
    ChevronRight,
    Search,
    Filter,
    ArrowUpRight,
    ArrowDownLeft,
    Loader2,
    MessageSquare,
    ExternalLink
} from "lucide-react";
import Link from "next/link";
import { tradeService } from "@/services/tradeService";
import DashboardLayout from "@/components/dashboard/DashboardLayout";

const UserTradesPage = () => {
    const [trades, setTrades] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState<"all" | "active" | "completed" | "cancelled">("all");

    useEffect(() => {
        fetchTrades();
    }, []);

    const fetchTrades = async () => {
        setIsLoading(true);
        try {
            const data = await tradeService.getMyTrades();
            setTrades(data || []);
        } catch (err: any) {
            console.error(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "pending": return "text-yellow-500 bg-yellow-500/10 border-yellow-500/20";
            case "active": return "text-primary bg-primary/10 border-primary/20";
            case "paid": return "text-blue-500 bg-blue-500/10 border-blue-500/20";
            case "released": return "text-accent bg-accent/10 border-accent/20";
            case "completed": return "text-accent bg-accent/10 border-accent/20";
            case "cancelled": return "text-red-500 bg-red-500/10 border-red-500/20";
            case "expired": return "text-text-dim bg-white/5 border-white/10";
            default: return "text-text-dim bg-white/5 border-white/10";
        }
    };

    const filteredTrades = trades.filter(t => {
        if (filter === "all") return true;
        if (filter === "active") return ["pending", "active", "paid", "disputed"].includes(t.status);
        if (filter === "completed") return ["released", "completed"].includes(t.status);
        if (filter === "cancelled") return ["cancelled", "expired"].includes(t.status);
        return true;
    });

    return (
        <DashboardLayout title="My Trades" role="user">
            <div className="space-y-8 pb-10">
                {/* Header & Filters */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                    <div>
                        <h1 className="text-3xl font-black text-white italic uppercase tracking-tight">Active <span className="gradient-text">Operations</span></h1>
                        <p className="text-[10px] text-text-dim mt-1 font-medium italic">Track your peer-to-peer trade executions in real-time.</p>
                    </div>
                    <div className="flex items-center gap-2 bg-surface/50 p-1 rounded-xl border border-white/5">
                        {["all", "active", "completed", "cancelled"].map((f) => (
                            <button
                                key={f}
                                onClick={() => setFilter(f as any)}
                                className={`px-4 py-2 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all ${filter === f ? "bg-white text-background" : "text-text-dim hover:text-white"}`}
                            >
                                {f}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Trades List */}
                <div className="space-y-4">
                    <AnimatePresence mode="wait">
                        {isLoading ? (
                            <motion.div
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                className="flex flex-col items-center justify-center py-24 space-y-4"
                            >
                                <Loader2 className="w-10 h-10 text-primary animate-spin" />
                                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">Syncing Ledger...</p>
                            </motion.div>
                        ) : filteredTrades.length === 0 ? (
                            <motion.div
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="bg-surface border border-white/5 rounded-[2.5rem] p-20 text-center space-y-4"
                            >
                                <div className="w-16 h-16 bg-white/5 rounded-full flex items-center justify-center mx-auto">
                                    <Clock className="w-8 h-8 text-text-dim/20" />
                                </div>
                                <h3 className="text-xl font-black text-white italic uppercase tracking-tight">No Transactions Found</h3>
                                <p className="text-xs text-text-dim max-w-xs mx-auto">Your trade history will appear here once you initiate or respond to an offer.</p>
                                <Link
                                    href="/marketplace"
                                    className="inline-flex items-center px-6 py-3 bg-white text-background rounded-xl font-black uppercase tracking-widest text-[10px] hover:scale-105 transition-all mt-4"
                                >
                                    Explore Marketplace
                                </Link>
                            </motion.div>
                        ) : (
                            filteredTrades.map((trade, i) => (
                                <motion.div
                                    key={trade.trade_id}
                                    initial={{ opacity: 0, x: -20 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    transition={{ delay: i * 0.05 }}
                                    className="group relative bg-surface border border-white/5 rounded-[1.5rem] p-6 hover:border-white/10 transition-all cursor-pointer"
                                >
                                    <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                                        <div className="flex items-center gap-6">
                                            {/* Type Indicator */}
                                            <div className={`p-4 rounded-2xl bg-white/5 ${trade.type === 'buy' ? "text-primary" : "text-accent"}`}>
                                                {trade.type === 'buy' ? <ArrowUpRight className="w-6 h-6" /> : <ArrowDownLeft className="w-6 h-6" />}
                                            </div>

                                            <div className="space-y-1">
                                                <div className="flex items-center gap-2">
                                                    <span className="text-[10px] font-black uppercase tracking-widest text-text-dim/50">#{trade.trade_id.slice(0, 8)}</span>
                                                    <span className={`px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border ${getStatusColor(trade.status)}`}>
                                                        {trade.status}
                                                    </span>
                                                </div>
                                                <h3 className="text-lg font-black text-white">
                                                    {trade.crypto_amount.toFixed(4)} <span className="text-text-dim text-xs font-medium uppercase tracking-tight">USDT</span>
                                                    <span className="text-white/20 mx-2">/</span>
                                                    {trade.fiat_amount.toLocaleString()} <span className="text-text-dim text-xs font-medium uppercase tracking-tight">USD</span>
                                                </h3>
                                                <div className="flex items-center gap-3">
                                                    <p className="text-[9px] text-text-dim/60 font-bold uppercase tracking-widest flex items-center gap-1.5">
                                                        <Clock className="w-3 h-3" />
                                                        {new Date(trade.created_at).toLocaleString()}
                                                    </p>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-2">
                                            <button className="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white/5 text-text-dim hover:text-white transition-all text-[10px] font-black uppercase tracking-widest border border-transparent hover:border-white/10">
                                                <MessageSquare className="w-3.5 h-3.5" />
                                                Support
                                            </button>
                                            <Link
                                                href={`/marketplace/trade/${trade.trade_id}`}
                                                className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary text-white transition-all text-[10px] font-black uppercase tracking-widest shadow-lg shadow-primary/10 hover:scale-105 active:scale-95"
                                            >
                                                Details
                                                <ChevronRight className="w-3.5 h-3.5" />
                                            </Link>
                                        </div>
                                    </div>
                                </motion.div>
                            ))
                        )}
                    </AnimatePresence>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserTradesPage;
