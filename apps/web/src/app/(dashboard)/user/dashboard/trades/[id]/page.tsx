"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    ArrowLeft,
    Clock,
    CheckCircle2,
    XCircle,
    AlertTriangle,
    Loader2,
    User,
    DollarSign,
    Shield,
    FileText,
    Download
} from "lucide-react";
import { useParams, useRouter } from "next/navigation";
import { authService } from "@/services/authService";
import { TradeService } from "@/services/tradeService";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import TradeChat from "@/components/trade/TradeChat";
import { toast } from "sonner";

interface TradeDetail {
    trade_id: string;
    ad_id: string;
    buyer_id: string;
    seller_id: string;
    crypto_amount: number;
    fiat_amount: number;
    crypto_symbol: string;
    fiat_symbol: string;
    exchange_rate: number;
    status: string;
    type: "buy" | "sell";
    created_at: string;
    updated_at: string;
    buyer_username?: string;
    seller_username?: string;
    payment_method?: number;
    payment_method_name?: string;
    payment_details?: any;
    escrow_id?: string;
    dispute_id?: string;
    payment_window_minutes?: number;
    timer_expires_at?: string;
}

const TradeDetailPage = () => {
    const params = useParams();
    const router = useRouter();
    const tradeId = params.id as string;

    // Validate tradeId exists
    if (!tradeId) {
        return (
            <DashboardLayout title="Trade Details" role="user">
                <div className="text-center py-12">
                    <p className="text-text-dim">Invalid trade ID</p>
                </div>
            </DashboardLayout>
        );
    }

    const [trade, setTrade] = useState<TradeDetail | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [currentUser, setCurrentUser] = useState<any>(null);
    const [timeRemaining, setTimeRemaining] = useState<number>(0);
    const [isActionLoading, setIsActionLoading] = useState<string | null>(null);

    useEffect(() => {
        fetchTradeDetail();
        fetchCurrentUser();
    }, [tradeId]);

    // Payment timer countdown
    useEffect(() => {
        if (!trade?.timer_expires_at || trade.status !== "active") {
            setTimeRemaining(0);
            return;
        }

        const calculateTimeRemaining = () => {
            const now = new Date().getTime();
            const expiryTime = new Date(trade.timer_expires_at as string).getTime();
            const remaining = Math.max(0, expiryTime - now);
            setTimeRemaining(remaining);
        };

        calculateTimeRemaining();
        const interval = setInterval(calculateTimeRemaining, 1000);

        return () => clearInterval(interval);
    }, [trade?.timer_expires_at, trade?.status]);

    const fetchTradeDetail = async () => {
        setIsLoading(true);
        try {
            const data = await TradeService.getTradeDetails(tradeId);
            setTrade(data);
        } catch (error: any) {
            console.error("Error fetching trade detail:", error);
            toast.error("Failed to load trade details");
        } finally {
            setIsLoading(false);
        }
    };

    const fetchCurrentUser = async () => {
        try {
            const user = await authService.getCurrentUser();
            setCurrentUser(user);
        } catch (error) {
            console.error("Error fetching current user:", error);
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
            case "disputed": return "text-orange-500 bg-orange-500/10 border-orange-500/20";
            default: return "text-text-dim bg-white/10 border-white/20";
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "completed": return <CheckCircle2 className="w-4 h-4" />;
            case "cancelled": return <XCircle className="w-4 h-4" />;
            case "disputed": return <AlertTriangle className="w-4 h-4" />;
            default: return <Clock className="w-4 h-4" />;
        }
    };

    const formatTimeRemaining = (milliseconds: number) => {
        if (milliseconds <= 0) return "Expired";

        const hours = Math.floor(milliseconds / (1000 * 60 * 60));
        const minutes = Math.floor((milliseconds % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((milliseconds % (1000 * 60)) / 1000);

        if (hours > 0) {
            return `${hours}h ${minutes}m ${seconds}s`;
        } else if (minutes > 0) {
            return `${minutes}m ${seconds}s`;
        } else {
            return `${seconds}s`;
        }
    };

    const handleMarkAsPaid = async () => {
        setIsActionLoading("mark_paid");
        try {
            await TradeService.updateTradeStatus(tradeId, "pay");
            toast.success("Trade marked as paid");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to mark as paid");
        } finally {
            setIsActionLoading(null);
        }
    };

    const handleReleaseEscrow = async () => {
        setIsActionLoading("release");
        try {
            await TradeService.updateTradeStatus(tradeId, "release");
            toast.success("Escrow released successfully");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to release escrow");
        } finally {
            setIsActionLoading(null);
        }
    };

    const handleDisputeTrade = async () => {
        const reason = prompt("Please enter dispute reason:");
        if (!reason) return;

        setIsActionLoading("dispute");
        try {
            await TradeService.disputeTrade(tradeId, "OTHER", reason);
            toast.success("Dispute created successfully");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to create dispute");
        } finally {
            setIsActionLoading(null);
        }
    };

    const handleCancelTrade = async () => {
        if (!confirm("Are you sure you want to cancel this trade?")) return;

        setIsActionLoading("cancel");
        try {
            await TradeService.updateTradeStatus(tradeId, "cancel");
            toast.success("Trade cancelled successfully");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to cancel trade");
        } finally {
            setIsActionLoading(null);
        }
    };

    const getCounterpartUsername = () => {
        if (!trade || !currentUser) return "";

        if (trade.buyer_id === currentUser.id) {
            return trade.seller_username;
        } else {
            return trade.buyer_username;
        }
    };

    const isBuyer = () => {
        return trade && currentUser && trade.buyer_id === currentUser.id;
    };

    const canMarkAsPaid = () => {
        return isBuyer() && trade?.status === "active" && timeRemaining > 0;
    };

    const canReleaseEscrow = () => {
        return !isBuyer() && trade?.status === "paid";
    };

    const canDispute = () => {
        return (trade?.status === "active" || trade?.status === "paid") && !trade?.dispute_id;
    };

    const canCancel = () => {
        return trade?.status === "active" && timeRemaining > 0;
    };

    const isTimerExpired = () => {
        return trade?.status === "active" && timeRemaining <= 0;
    };

    if (isLoading) {
        return (
            <DashboardLayout title="Trade Details" role="user">
                <div className="flex items-center justify-center h-64">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
            </DashboardLayout>
        );
    }

    if (!trade) {
        return (
            <DashboardLayout title="Trade Details" role="user">
                <div className="text-center py-12">
                    <p className="text-text-dim">Trade not found</p>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Trade Details" role="user">
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <button
                            onClick={() => router.back()}
                            className="p-2 text-text-dim hover:text-white transition-colors"
                        >
                            <ArrowLeft className="w-5 h-5" />
                        </button>
                        <div>
                            <h1 className="text-2xl font-black text-white italic uppercase tracking-tighter shadow-sm">Trade Details</h1>
                            <p className="text-text-dim text-[10px] uppercase font-bold tracking-widest">Trade ID: #{trade.trade_id.slice(0, 8)}</p>
                        </div>
                    </div>
                    <div className={`flex items-center gap-2 px-4 py-2 rounded-xl border-2 shadow-lg ${getStatusColor(trade.status)}`}>
                        {getStatusIcon(trade.status)}
                        <span className="text-xs font-black uppercase tracking-widest">
                            {trade.status}
                        </span>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Trade Information */}
                    <div className="lg:col-span-2 space-y-6 transition-all">
                        {/* Trade Summary */}
                        <div className="bg-surface border border-white/10 rounded-[2rem] p-8 relative overflow-hidden group shadow-2xl">
                            <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 blur-3xl -z-10 group-hover:bg-primary/10 transition-colors" />
                            <h2 className="text-lg font-black text-white mb-6 uppercase italic tracking-tight">Trade Summary</h2>

                            <div className="grid grid-cols-2 gap-8">
                                <div className="space-y-1">
                                    <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">You are</p>
                                    <p className="text-2xl font-black text-white italic">
                                        {isBuyer() ? "Buying" : "Selling"}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">Amount</p>
                                    <p className="text-2xl font-black text-white italic">
                                        {trade.crypto_amount?.toFixed(4) || '0.0000'} {trade.crypto_symbol || 'USDT'}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">Rate</p>
                                    <p className="text-2xl font-black text-white italic">
                                        {trade.exchange_rate?.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'} {trade.fiat_symbol || 'USD'}
                                    </p>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">{isBuyer() ? "Total Pay" : "Total Receive"}</p>
                                    <p className="text-3xl font-black text-primary italic">
                                        {trade.fiat_amount?.toLocaleString() || '0'} {trade.fiat_symbol || 'USD'}
                                    </p>
                                </div>
                            </div>

                            <div className="mt-8 pt-8 border-t border-white/5 grid grid-cols-2 gap-4">
                                <div className="flex items-center gap-4 bg-white/5 p-4 rounded-2xl border border-white/5">
                                    <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center">
                                        <User className="w-5 h-5 text-primary" />
                                    </div>
                                    <div>
                                        <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">Counterpart</p>
                                        <p className="text-white font-bold">{getCounterpartUsername()}</p>
                                    </div>
                                </div>
                                <div className="flex items-center gap-4 bg-white/5 p-4 rounded-2xl border border-white/5">
                                    <div className="w-10 h-10 bg-accent/10 rounded-xl flex items-center justify-center">
                                        <Clock className="w-5 h-5 text-accent" />
                                    </div>
                                    <div>
                                        <p className="text-text-dim text-[10px] uppercase font-black tracking-widest">Started</p>
                                        <p className="text-white font-bold">
                                            {trade.created_at ? new Date(trade.created_at as string).toLocaleDateString() : 'Unknown'}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Payment Information */}
                        {(trade.payment_method_name || trade.payment_method) && (
                            <div className="bg-surface border border-white/10 rounded-[2rem] p-8 shadow-2xl relative overflow-hidden group">
                                <div className="absolute -bottom-4 -left-4 w-24 h-24 bg-accent/5 blur-2xl -z-10 group-hover:bg-accent/10 transition-colors" />
                                <h2 className="text-lg font-black text-white mb-6 uppercase italic tracking-tight">Payment Information</h2>

                                <div className="space-y-6">
                                    <div className="flex items-center gap-4">
                                        <div className="w-12 h-12 bg-white/5 rounded-2xl flex items-center justify-center border border-white/10">
                                            <DollarSign className="w-6 h-6 text-primary" />
                                        </div>
                                        <div>
                                            <p className="text-text-dim text-[10px] uppercase font-black tracking-widest mb-1">Payment Method</p>
                                            <p className="text-xl font-bold text-white capitalize">{trade.payment_method_name || "Bank Transfer"}</p>
                                        </div>
                                    </div>

                                    {trade.payment_details ? (
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
                                            {trade.payment_details.bank_name && (
                                                <div className="p-4 rounded-2xl bg-white/5 border border-white/5">
                                                    <p className="text-text-dim text-[8px] uppercase font-black tracking-widest mb-1">Bank Name</p>
                                                    <p className="text-white font-bold">{trade.payment_details.bank_name}</p>
                                                </div>
                                            )}
                                            {trade.payment_details.account_number && (
                                                <div className="p-4 rounded-2xl bg-white/5 border border-white/5 group/copy relative">
                                                    <p className="text-text-dim text-[8px] uppercase font-black tracking-widest mb-1">Account Number</p>
                                                    <p className="text-white font-bold font-mono tracking-wider">{trade.payment_details.account_number}</p>
                                                </div>
                                            )}
                                            {trade.payment_details.account_name && (
                                                <div className="p-4 rounded-2xl bg-white/5 border border-white/5 col-span-1 md:col-span-2">
                                                    <p className="text-text-dim text-[8px] uppercase font-black tracking-widest mb-1">Account Name</p>
                                                    <p className="text-white font-bold">{trade.payment_details.account_name}</p>
                                                </div>
                                            )}
                                            {trade.payment_details.message && (
                                                <div className="p-4 rounded-2xl bg-primary/5 border border-primary/10 col-span-1 md:col-span-2">
                                                    <p className="text-primary text-[10px] font-bold">{trade.payment_details.message}</p>
                                                </div>
                                            )}
                                        </div>
                                    ) : (
                                        <div className="p-4 rounded-2xl bg-white/5 border border-white/5 text-center">
                                            <p className="text-text-dim text-[10px] font-bold uppercase tracking-widest italic opacity-50">Fetching payment details...</p>
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}
                        {/* Payment Timer */}
                        {trade.status === "active" && (
                            <div className={`bg-surface border rounded-[2rem] p-8 shadow-2xl ${isTimerExpired()
                                ? "border-red-500/20 bg-red-500/5 shadow-red-500/5"
                                : "border-white/10"
                                }`}>
                                <h2 className="text-lg font-black text-white mb-6 uppercase italic tracking-tight">Payment Timer</h2>
                                <div className="flex items-center justify-between">
                                    <div>
                                        <p className="text-text-dim text-[10px] uppercase font-black tracking-widest mb-1 leading-none shadow-sm">
                                            {isBuyer() ? "Time to complete payment" : "Waiting for buyer payment"}
                                        </p>
                                        <p className={`text-3xl font-black italic tracking-tighter ${isTimerExpired() ? "text-red-500" : "text-white text-glow-primary"
                                            }`}>
                                            {formatTimeRemaining(timeRemaining)}
                                        </p>
                                    </div>
                                    <div className={`w-14 h-14 rounded-2xl flex items-center justify-center border-2 ${isTimerExpired()
                                        ? "bg-red-500/20 text-red-500 border-red-500/30"
                                        : "bg-primary/20 text-primary border-primary/30"
                                        }`}>
                                        <Clock className="w-8 h-8" />
                                    </div>
                                </div>
                                {isTimerExpired() && (
                                    <div className="mt-6 p-4 bg-red-500/10 border border-red-500/20 rounded-2xl animate-pulse">
                                        <p className="text-red-500 text-xs font-black uppercase tracking-widest text-center">
                                            Payment time expired!
                                        </p>
                                    </div>
                                )}
                            </div>
                        )}

                        {/* Actions */}
                        <div className="bg-surface border border-white/10 rounded-[2rem] p-8 shadow-2xl relative overflow-hidden">
                            <div className="absolute top-0 right-0 w-24 h-24 bg-white/5 blur-3xl -z-10" />
                            <h2 className="text-lg font-black text-white mb-6 uppercase italic tracking-tight">Controls</h2>
                            <div className="flex flex-wrap gap-4">
                                {canMarkAsPaid() && (
                                    <button
                                        onClick={handleMarkAsPaid}
                                        disabled={isActionLoading === "pay"}
                                        className="h-12 px-8 bg-primary text-white rounded-2xl hover:bg-primary/90 transition-all active:scale-95 text-xs font-black uppercase tracking-widest disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3 shadow-xl shadow-primary/20 border-b-4 border-primary-dark"
                                    >
                                        {isActionLoading === "pay" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : <DollarSign className="w-4 h-4" />}
                                        Mark as Paid
                                    </button>
                                )}
                                {canReleaseEscrow() && (
                                    <button
                                        onClick={handleReleaseEscrow}
                                        disabled={isActionLoading === "release"}
                                        className="h-12 px-8 bg-accent text-white rounded-2xl hover:bg-accent/90 transition-all active:scale-95 text-xs font-black uppercase tracking-widest disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3 shadow-xl shadow-accent/20 border-b-4 border-accent-dark"
                                    >
                                        {isActionLoading === "release" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : <CheckCircle2 className="w-4 h-4" />}
                                        Release Escrow
                                    </button>
                                )}
                                {canCancel() && (
                                    <button
                                        onClick={handleCancelTrade}
                                        disabled={isActionLoading === "cancel"}
                                        className="h-12 px-6 bg-white/5 text-white border border-white/10 rounded-2xl hover:bg-white/10 transition-all active:scale-95 text-xs font-black uppercase tracking-widest disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3"
                                    >
                                        {isActionLoading === "cancel" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : <XCircle className="w-4 h-4" />}
                                        Cancel Trade
                                    </button>
                                )}
                                {canDispute() && (
                                    <button
                                        onClick={handleDisputeTrade}
                                        disabled={isActionLoading === "dispute"}
                                        className="h-12 px-6 bg-red-500/10 text-red-500 border border-red-500/20 rounded-2xl hover:bg-red-500/20 transition-all active:scale-95 text-xs font-black uppercase tracking-widest disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-3"
                                    >
                                        {isActionLoading === "dispute" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : <AlertTriangle className="w-4 h-4" />}
                                        Open Dispute
                                    </button>
                                )}
                                <button className="h-12 px-6 bg-white/5 text-text-dim border border-white/10 rounded-2xl hover:bg-white/10 hover:text-white transition-all active:scale-95 text-xs font-black uppercase tracking-widest flex items-center gap-3">
                                    <Download className="w-4 h-4" />
                                    Receipt
                                </button>
                            </div>
                        </div>
                    </div>

                    {/* Chat */}
                    <div className="lg:col-span-1 h-full">
                        <div className="bg-surface border border-white/10 rounded-[2rem] overflow-hidden shadow-2xl h-[700px] flex flex-col sticky top-6">
                            <div className="p-6 border-b border-white/5 bg-white/5">
                                <h3 className="text-sm font-black text-white italic uppercase tracking-widest flex items-center gap-2">
                                    <Shield className="w-4 h-4 text-primary" />
                                    Secure Trade Chat
                                </h3>
                            </div>
                            <div className="flex-1 overflow-hidden relative">
                                <TradeChat
                                    tradeId={trade.trade_id}
                                    currentUserId={currentUser?.id || ""}
                                    counterpartUsername={getCounterpartUsername()}
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default TradeDetailPage;
