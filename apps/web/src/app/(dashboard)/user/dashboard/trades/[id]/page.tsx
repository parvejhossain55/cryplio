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

    useEffect(() => {
        fetchTradeDetail();
        fetchCurrentUser();
    }, [tradeId]);

    const fetchTradeDetail = async () => {
        setIsLoading(true);
        try {
            const data = await authService.getTradeDetails(tradeId);
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

    const handleMarkAsPaid = async () => {
        try {
            await authService.updateTradeStatus(tradeId, "mark_paid");
            toast.success("Trade marked as paid");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to mark as paid");
        }
    };

    const handleReleaseEscrow = async () => {
        try {
            await authService.updateTradeStatus(tradeId, "release");
            toast.success("Escrow released successfully");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to release escrow");
        }
    };

    const handleDisputeTrade = async () => {
        const reason = prompt("Please enter dispute reason:");
        if (!reason) return;

        try {
            await authService.disputeTrade(tradeId, "OTHER", reason);
            toast.success("Dispute created successfully");
            fetchTradeDetail();
        } catch (error: any) {
            toast.error(error.message || "Failed to create dispute");
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
        return isBuyer() && trade?.status === "active";
    };

    const canReleaseEscrow = () => {
        return !isBuyer() && trade?.status === "paid";
    };

    const canDispute = () => {
        return trade?.status === "active" || trade?.status === "paid";
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
                                                {trade.created_at ? new Date(trade.created_at).toLocaleDateString() : 'Unknown'}
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

                        {/* Actions */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h2 className="text-lg font-black text-white mb-4">Actions</h2>
                            <div className="flex flex-wrap gap-3">
                                {canMarkAsPaid() && (
                                    <button
                                        onClick={handleMarkAsPaid}
                                        className="px-4 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors text-sm font-black uppercase tracking-wider"
                                    >
                                        Mark as Paid
                                    </button>
                                )}
                                {canReleaseEscrow() && (
                                    <button
                                        onClick={handleReleaseEscrow}
                                        className="px-4 py-2 bg-accent text-white rounded-xl hover:bg-accent/90 transition-colors text-sm font-black uppercase tracking-wider"
                                    >
                                        Release Escrow
                                    </button>
                                )}
                                {canDispute() && (
                                    <button
                                        onClick={handleDisputeTrade}
                                        className="px-4 py-2 bg-red-500 text-white rounded-xl hover:bg-red-500/90 transition-colors text-sm font-black uppercase tracking-wider"
                                    >
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
