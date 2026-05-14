"use client";

import React, { useState, useEffect } from "react";
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
    ShieldCheck,
    Loader2
} from "lucide-react";
import { motion } from "framer-motion";
import { WalletService, WalletBalance, WalletTransaction } from "@/services/walletService";
import { TradeService } from "@/services/tradeService";
import { useAuth } from "@/context/AuthContext";
import Link from "next/link";

const UserDashboard = () => {
    const { user } = useAuth();
    const [balances, setBalances] = useState<WalletBalance[]>([]);
    const [activities, setActivities] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fetchDashboardData = async () => {
        setIsLoading(true);
        try {
            const [balanceData, transactionData, tradeData] = await Promise.all([
                WalletService.getBalances(),
                WalletService.getTransactions({ limit: 10, offset: 0 }),
                TradeService.getMyTrades()
            ]);

            setBalances(balanceData || []);

            // Merge and sort activities
            const walletActivities = (transactionData.transactions || []).map((tx: WalletTransaction) => ({
                id: tx.tx_id,
                type: tx.type, // 'deposit' or 'withdrawal'
                asset: "USDT", // Should ideally come from backend
                amount: tx.amount,
                status: tx.status,
                date: new Date(tx.created_at).toLocaleString(),
                raw_date: tx.created_at,
                price: "-"
            }));

            const tradeActivities = (tradeData || []).map((t: any) => ({
                id: t.trade_id,
                type: t.buyer_id === user?.id ? 'buy' : 'sell',
                asset: t.crypto_symbol || "USDT",
                amount: t.crypto_amount,
                status: t.status,
                date: new Date(t.created_at).toLocaleString(),
                raw_date: t.created_at,
                price: `$${t.exchange_rate}`
            }));

            const merged = [...walletActivities, ...tradeActivities]
                .sort((a, b) => new Date(b.raw_date).getTime() - new Date(a.raw_date).getTime())
                .slice(0, 8);

            setActivities(merged);
        } catch (error) {
            console.error("Failed to fetch dashboard data:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const totalBalance = balances.reduce((acc, curr) => acc + (curr.balance || 0), 0);

    if (isLoading) {
        return (
            <DashboardLayout title="Overview" role="user">
                <div className="flex items-center justify-center min-h-[60vh]">
                    <Loader2 className="w-12 h-12 text-primary animate-spin" />
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Overview" role="user">
            <div className="grid grid-cols-1">
                {/* Left Column: Stats & Wallet */}
                <div className="lg:col-span-2 space-y-8">
                    {/* Main Wallet Card */}
                    <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-8 group">
                        <div className="absolute top-0 right-0 w-[50%] h-[150%] bg-primary/10 -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

                        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                            <div>
                                <p className="text-xs font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Balance (USDT)</p>
                                <div className="flex items-baseline space-x-3">
                                    <h2 className="text-5xl font-black text-white italic">₮{totalBalance.toLocaleString(undefined, { minimumFractionDigits: 2 })}</h2>
                                    <span className="text-accent font-bold bg-accent/10 px-2 py-0.5 rounded-lg text-sm flex items-center">
                                        <TrendingUp className="w-3 h-3 mr-1" /> +0.0%
                                    </span>
                                </div>
                            </div>

                            <div className="flex items-center space-x-3">
                                <Link href="/user/dashboard/wallet" className="flex-1 md:flex-none px-6 py-4 bg-white text-background rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:scale-105 active:scale-95 shadow-xl shadow-white/5">
                                    <ArrowUpRight className="w-4 h-4 mr-2" /> Send
                                </Link>
                                <Link href="/user/dashboard/wallet" className="flex-1 md:flex-none px-6 py-4 bg-surface-light border border-white/5 text-white rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:bg-white/5 active:scale-95">
                                    <ArrowDownLeft className="w-4 h-4 mr-2" /> Receive
                                </Link>
                            </div>
                        </div>

                        <div className="relative z-10 mt-10 grid grid-cols-2 sm:grid-cols-4 gap-4">
                            {balances.length > 0 ? balances.map((coin) => (
                                <div key={coin.crypto_symbol} className="p-4 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 transition-all cursor-pointer group/card active:scale-95">
                                    <div className="w-10 h-10 rounded-2xl bg-surface border border-white/10 flex items-center justify-center text-lg font-black mb-3 group-hover/card:scale-110 group-hover/card:text-primary transition-all">
                                        {coin.crypto_symbol?.[0] || 'C'}
                                    </div>
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{coin.crypto_symbol}</p>
                                    <p className="text-sm font-bold text-white mt-1">{coin.balance.toFixed(4)}</p>
                                    <p className="text-[10px] font-medium text-text-dim mt-0.5">${(coin.balance * 1).toFixed(2)}</p>
                                </div>
                            )) : (
                                <div className="col-span-4 p-8 text-center border border-dashed border-white/10 rounded-3xl">
                                    <p className="text-xs font-bold text-text-dim uppercase tracking-widest">No assets found in your portfolio</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Recent Activity */}
                    <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8 relative overflow-hidden">
                        <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 blur-3xl -z-10" />

                        <div className="flex items-center justify-between mb-8">
                            <h3 className="text-xl font-black text-white italic uppercase tracking-tight">Recent Activity</h3>
                            <Link href="/user/dashboard/trades" className="text-xs font-black text-primary hover:text-white transition-all flex items-center uppercase tracking-widest">
                                View All <ChevronRight className="w-4 h-4 ml-1" />
                            </Link>
                        </div>

                        <div className="space-y-4">
                            {activities.length > 0 ? activities.map((tx, i) => (
                                <Link
                                    key={i}
                                    href={tx.type === 'buy' || tx.type === 'sell' ? `/user/dashboard/trades/${tx.id}` : '/user/dashboard/wallet'}
                                    className="flex items-center justify-between p-4 rounded-2xl hover:bg-white/5 border border-transparent hover:border-white/5 transition-all group"
                                >
                                    <div className="flex items-center space-x-4">
                                        <div className={`w-12 h-12 rounded-2xl flex items-center justify-center border border-white/5 shadow-inner transition-transform group-hover:scale-110 ${tx.type === 'buy' || tx.type === 'deposit' ? 'bg-accent/10 text-accent' :
                                            tx.type === 'sell' || tx.type === 'withdrawal' ? 'bg-primary/10 text-primary' :
                                                'bg-secondary/10 text-secondary'
                                            }`}>
                                            {tx.type === 'buy' || tx.type === 'deposit' ? <ArrowDownLeft className="w-6 h-6" /> :
                                                tx.type === 'sell' || tx.type === 'withdrawal' ? <ArrowUpRight className="w-6 h-6" /> :
                                                    <ArrowLeftRight className="w-5 h-5" />}
                                        </div>
                                        <div>
                                            <h4 className="font-black text-white tracking-tight flex items-center uppercase italic">
                                                {tx.type} {tx.asset}
                                                <span className={`ml-3 text-[8px] px-2 py-0.5 rounded-lg font-black uppercase tracking-widest border ${tx.status === 'completed' || tx.status === 'active' || tx.status === 'success' ? 'bg-accent/10 text-accent border-accent/20' :
                                                    tx.status === 'pending' ? 'bg-primary/10 text-primary border-primary/20' :
                                                        'bg-red-500/10 text-red-500 border-red-500/20'
                                                    }`}>
                                                    {tx.status}
                                                </span>
                                            </h4>
                                            <p className="text-[10px] font-bold text-text-dim mt-1 uppercase tracking-widest opacity-60">
                                                {tx.date} {tx.price !== '-' && `• Price: ${tx.price}`}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <p className="font-black text-white text-lg italic">{tx.amount.toString()} {tx.asset.split(' ')[0]}</p>
                                        <ExternalLink className="w-4 h-4 ml-auto mt-1 text-text-dim opacity-0 group-hover:opacity-100 transition-all cursor-pointer hover:text-white" />
                                    </div>
                                </Link>
                            )) : (
                                <div className="py-12 text-center border border-dashed border-white/10 rounded-3xl">
                                    <Clock className="w-8 h-8 text-text-dim mx-auto mb-3 opacity-20" />
                                    <p className="text-[10px] font-bold text-text-dim uppercase tracking-[0.2em]">No recent activity to display</p>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserDashboard;
