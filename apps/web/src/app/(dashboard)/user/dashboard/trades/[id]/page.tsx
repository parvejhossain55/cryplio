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
import { tradeService } from "@/services/tradeService";
import { userService } from "@/services/userService";
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
    rate: number;
    status: string;
    type: "buy" | "sell";
    created_at: string;
    updated_at: string;
    buyer_username?: string;
    seller_username?: string;
    payment_method?: string;
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
            const data = await tradeService.getTradeDetails(tradeId);
            setTrade(data);
        } catch (error: any) {
            toast.error("Failed to load trade details");
            router.push("/user/dashboard/trades");
        } finally {
            setIsLoading(false);
        }
    };

    const fetchCurrentUser = async () => {
        try {
            const user = await userService.getCurrentUser();
            setCurrentUser(user);
        } catch (error) {
            console.error("Failed to fetch current user:", error);
        }
    };

    const handleAction = async (action: string) => {
        setIsActionLoading(action);
        try {
            await tradeService.updateTradeStatus(tradeId, action);
            toast.success(`Trade ${action === "pay" ? "marked as paid" : (action === "release" ? "released" : "cancelled")}`);
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || `Failed to ${action} trade`);
        } finally {
            setIsActionLoading(null);
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
            await tradeService.updateTradeStatus(tradeId, "mark_paid");
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
            await tradeService.updateTradeStatus(tradeId, "release");
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
            await tradeService.disputeTrade(tradeId, "OTHER", reason);
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
            await tradeService.updateTradeStatus(tradeId, "cancel");
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
        
        if (trade.buyer_id === currentUser.user_id) {
            return trade.seller_username;
        } else {
            return trade.buyer_username;
        }
    };

    const isBuyer = () => {
        return trade && currentUser && trade.buyer_id === currentUser.user_id;
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
                            <h1 className="text-2xl font-black text-white">Trade Details</h1>
                            <p className="text-text-dim text-sm">Trade ID: #{trade.trade_id.slice(0, 8)}</p>
                        </div>
                    </div>
                    <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full border ${getStatusColor(trade.status)}`}>
                        {getStatusIcon(trade.status)}
                        <span className="text-xs font-black uppercase tracking-wider">
                            {trade.status}
                        </span>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Trade Information */}
                    <div className="lg:col-span-2 space-y-6">
                        {/* Trade Summary */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h2 className="text-lg font-black text-white mb-4">Trade Summary</h2>
                            
                            <div className="grid grid-cols-2 gap-6">
                                <div>
                                    <p className="text-text-dim text-sm mb-1">You are</p>
                                    <p className="text-xl font-black text-white">
                                        {isBuyer() ? "Buying" : "Selling"}
                                    </p>
                                </div>
                                <div>
                                    <p className="text-text-dim text-sm mb-1">Amount</p>
                                    <p className="text-xl font-black text-white">
                                        {trade.crypto_amount?.toFixed(4) || '0.0000'} {trade.crypto_symbol || 'USDT'}
                                    </p>
                                </div>
                                <div>
                                    <p className="text-text-dim text-sm mb-1">Rate</p>
                                    <p className="text-xl font-black text-white">
                                        {trade.rate?.toFixed(2) || '0.00'} {trade.fiat_symbol || 'USD'}
                                    </p>
                                </div>
                                <div>
                                    <p className="text-text-dim text-sm mb-1">Total</p>
                                    <p className="text-xl font-black text-white">
                                        {trade.fiat_amount?.toLocaleString() || '0'} {trade.fiat_symbol || 'USD'}
                                    </p>
                                </div>
                            </div>

                            <div className="mt-6 pt-6 border-t border-white/10">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-3">
                                        <User className="w-4 h-4 text-text-dim" />
                                        <div>
                                            <p className="text-text-dim text-sm">Counterpart</p>
                                            <p className="text-white font-medium">{getCounterpartUsername()}</p>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <Clock className="w-4 h-4 text-text-dim" />
                                        <div>
                                            <p className="text-text-dim text-sm">Created</p>
                                            <p className="text-white font-medium">
                                                {trade.created_at ? new Date(trade.created_at as string).toLocaleDateString() : 'Unknown'}
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Payment Information */}
                        {trade.payment_method && (
                            <div className="bg-surface border border-white/10 rounded-2xl p-6">
                                <h2 className="text-lg font-black text-white mb-4">Payment Information</h2>
                                <div className="space-y-3">
                                    <div className="flex items-center gap-3">
                                        <DollarSign className="w-4 h-4 text-text-dim" />
                                        <p className="text-text-dim text-sm">Payment Method</p>
                                        <p className="text-white font-medium">{trade.payment_method}</p>
                                    </div>
                                    {trade.payment_details && (
                                        <div className="flex items-center gap-3">
                                            <FileText className="w-4 h-4 text-text-dim" />
                                            <p className="text-text-dim text-sm">Payment Details</p>
                                            <button className="text-primary hover:text-primary/80 text-sm font-medium">
                                                View Details
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Payment Timer */}
                        {trade.status === "active" && (
                            <div className={`bg-surface border rounded-2xl p-6 ${
                                isTimerExpired() 
                                    ? "border-red-500/20 bg-red-500/5" 
                                    : "border-white/10"
                            }`}>
                                <h2 className="text-lg font-black text-white mb-4">Payment Timer</h2>
                                <div className="flex items-center justify-between">
                                    <div>
                                        <p className="text-text-dim text-sm mb-1">
                                            {isBuyer() ? "Time to complete payment" : "Waiting for buyer payment"}
                                        </p>
                                        <p className={`text-2xl font-black ${
                                            isTimerExpired() ? "text-red-500" : "text-white"
                                        }`}>
                                            {formatTimeRemaining(timeRemaining)}
                                        </p>
                                    </div>
                                    <div className={`p-3 rounded-full ${
                                        isTimerExpired() 
                                            ? "bg-red-500/20 text-red-500" 
                                            : "bg-primary/20 text-primary"
                                    }`}>
                                        <Clock className="w-6 h-6" />
                                    </div>
                                </div>
                                {isTimerExpired() && (
                                    <div className="mt-4 p-3 bg-red-500/10 border border-red-500/20 rounded-xl">
                                        <p className="text-red-500 text-sm font-medium">
                                            Payment time expired. Trade will be auto-cancelled or disputed.
                                        </p>
                                    </div>
                                )}
                            </div>
                        )}

                        {/* Actions */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h2 className="text-lg font-black text-white mb-4">Actions</h2>
                            <div className="flex flex-wrap gap-3">
                                {canMarkAsPaid() && (
                                    <button
                                        onClick={handleMarkAsPaid}
                                        disabled={isActionLoading === "mark_paid"}
                                        className="px-4 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors text-sm font-black uppercase tracking-wider disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                                    >
                                        {isActionLoading === "mark_paid" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : null}
                                        Mark as Paid
                                    </button>
                                )}
                                {canReleaseEscrow() && (
                                    <button
                                        onClick={handleReleaseEscrow}
                                        disabled={isActionLoading === "release"}
                                        className="px-4 py-2 bg-accent text-white rounded-xl hover:bg-accent/90 transition-colors text-sm font-black uppercase tracking-wider disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                                    >
                                        {isActionLoading === "release" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : null}
                                        Release Escrow
                                    </button>
                                )}
                                {canCancel() && (
                                    <button
                                        onClick={handleCancelTrade}
                                        disabled={isActionLoading === "cancel"}
                                        className="px-4 py-2 bg-orange-500 text-white rounded-xl hover:bg-orange-500/90 transition-colors text-sm font-black uppercase tracking-wider disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                                    >
                                        {isActionLoading === "cancel" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : null}
                                        Cancel Trade
                                    </button>
                                )}
                                {canDispute() && (
                                    <button
                                        onClick={handleDisputeTrade}
                                        disabled={isActionLoading === "dispute"}
                                        className="px-4 py-2 bg-red-500 text-white rounded-xl hover:bg-red-500/90 transition-colors text-sm font-black uppercase tracking-wider disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                                    >
                                        {isActionLoading === "dispute" ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : null}
                                        Open Dispute
                                    </button>
                                )}
                                <button className="px-4 py-2 bg-white/10 text-white rounded-xl hover:bg-white/20 transition-colors text-sm font-black uppercase tracking-wider">
                                    <Download className="w-4 h-4 inline mr-2" />
                                    Download Receipt
                                </button>
                            </div>
                        </div>
                    </div>

                    {/* Chat */}
                    <div className="lg:col-span-1">
                        <div className="h-[600px]">
                            <TradeChat
                                tradeId={trade.trade_id}
                                currentUserId={currentUser?.user_id || ""}
                                counterpartUsername={getCounterpartUsername()}
                            />
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default TradeDetailPage;
