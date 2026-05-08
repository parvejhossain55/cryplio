"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    AlertTriangle,
    Clock,
    CheckCircle2,
    XCircle,
    Eye,
    MessageSquare,
    FileText,
    Download,
    Plus
} from "lucide-react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

interface Dispute {
    dispute_id: string;
    trade_id: string;
    raiser_id: string;
    respondent_id: string;
    reason: string;
    description: string;
    status: "open" | "investigating" | "resolved" | "closed";
    resolution?: string;
    resolution_notes?: string;
    resolved_by?: string;
    resolved_at?: string;
    evidence_files: string[];
    created_at: string;
    updated_at: string;
    trade_details?: {
        crypto_amount: number;
        fiat_amount: number;
        crypto_symbol: string;
        fiat_symbol: string;
        counterpart_username: string;
    };
}

const UserDisputes = () => {
    const router = useRouter();
    const [disputes, setDisputes] = useState<Dispute[]>([]);
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState<"all" | "open" | "resolved">("all");

    useEffect(() => {
        fetchDisputes();
    }, [filter]);

    const fetchDisputes = async () => {
        setLoading(true);
        try {
            const response = await fetch(`/api/disputes?status=${filter}`, {
                credentials: "include",
            });
            
            if (!response.ok) {
                throw new Error("Failed to fetch disputes");
            }
            
            const data = await response.json();
            setDisputes(data.disputes || []);
        } catch (error) {
            console.error("Error fetching disputes:", error);
            toast.error("Failed to load disputes");
        } finally {
            setLoading(false);
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "open":
                return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
            case "investigating":
                return <Clock className="w-4 h-4 text-blue-500" />;
            case "resolved":
                return <CheckCircle2 className="w-4 h-4 text-green-500" />;
            case "closed":
                return <XCircle className="w-4 h-4 text-gray-500" />;
            default:
                return <AlertTriangle className="w-4 h-4 text-gray-500" />;
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "open":
                return "bg-yellow-500/20 text-yellow-500 border-yellow-500/20";
            case "investigating":
                return "bg-blue-500/20 text-blue-500 border-blue-500/20";
            case "resolved":
                return "bg-green-500/20 text-green-500 border-green-500/20";
            case "closed":
                return "bg-gray-500/20 text-gray-500 border-gray-500/20";
            default:
                return "bg-gray-500/20 text-gray-500 border-gray-500/20";
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString();
    };

    const handleCreateDispute = () => {
        router.push("/dashboard/disputes/create");
    };

    if (loading) {
        return (
            <DashboardLayout title="Disputes" role="user">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Disputes" role="user">
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-white">Disputes</h1>
                        <p className="text-text-dim">Manage your trade disputes</p>
                    </div>
                    <button
                        onClick={handleCreateDispute}
                        className="px-4 py-2 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-colors flex items-center gap-2"
                    >
                        <Plus className="w-4 h-4" />
                        New Dispute
                    </button>
                </div>

                {/* Filter Tabs */}
                <div className="flex gap-2 p-1 bg-white/5 rounded-xl border border-white/10">
                    {[
                        { key: "all", label: "All Disputes" },
                        { key: "open", label: "Open" },
                        { key: "resolved", label: "Resolved" }
                    ].map((tab) => (
                        <button
                            key={tab.key}
                            onClick={() => setFilter(tab.key as any)}
                            className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-all ${
                                filter === tab.key
                                    ? "bg-primary text-white"
                                    : "text-text-dim hover:text-white"
                            }`}
                        >
                            {tab.label}
                        </button>
                    ))}
                </div>

                {/* Disputes List */}
                {disputes.length === 0 ? (
                    <div className="text-center py-12">
                        <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-white mb-2">No disputes found</h3>
                        <p className="text-text-dim">
                            {filter === "all" 
                                ? "You haven't been involved in any disputes yet."
                                : `No ${filter} disputes found.`
                            }
                        </p>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {disputes.map((dispute) => (
                            <motion.div
                                key={dispute.dispute_id}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="bg-surface border border-white/10 rounded-2xl p-6"
                            >
                                {/* Header */}
                                <div className="flex items-start justify-between mb-4">
                                    <div className="flex items-center gap-3">
                                        {getStatusIcon(dispute.status)}
                                        <div>
                                            <h3 className="text-lg font-semibold text-white">
                                                Dispute #{dispute.dispute_id.slice(0, 8)}
                                            </h3>
                                            <p className="text-text-dim text-sm">
                                                Trade #{dispute.trade_id.slice(0, 8)}
                                            </p>
                                        </div>
                                    </div>
                                    <div className={`px-3 py-1 rounded-full text-xs font-black uppercase border ${getStatusColor(dispute.status)}`}>
                                        {dispute.status}
                                    </div>
                                </div>

                                {/* Trade Details */}
                                {dispute.trade_details && (
                                    <div className="bg-white/5 rounded-xl p-4 mb-4">
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <p className="text-sm text-text-dim">Trade Amount</p>
                                                <p className="text-white font-semibold">
                                                    {dispute.trade_details.crypto_amount} {dispute.trade_details.crypto_symbol} / 
                                                    {dispute.trade_details.fiat_amount} {dispute.trade_details.fiat_symbol}
                                                </p>
                                            </div>
                                            <div className="text-right">
                                                <p className="text-sm text-text-dim">With</p>
                                                <p className="text-white font-semibold">
                                                    {dispute.trade_details.counterpart_username}
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                )}

                                {/* Dispute Details */}
                                <div className="mb-4">
                                    <div className="flex items-center gap-2 mb-2">
                                        <FileText className="w-4 h-4 text-primary" />
                                        <span className="text-sm font-semibold text-white">Reason</span>
                                    </div>
                                    <p className="text-white font-medium mb-2">{dispute.reason}</p>
                                    <p className="text-text-dim text-sm">{dispute.description}</p>
                                </div>

                                {/* Evidence Files */}
                                {dispute.evidence_files.length > 0 && (
                                    <div className="mb-4">
                                        <div className="flex items-center gap-2 mb-2">
                                            <FileText className="w-4 h-4 text-primary" />
                                            <span className="text-sm font-semibold text-white">Evidence Files</span>
                                        </div>
                                        <div className="flex flex-wrap gap-2">
                                            {dispute.evidence_files.map((file, index) => (
                                                <button
                                                    key={index}
                                                    className="flex items-center gap-2 px-3 py-2 bg-white/5 border border-white/10 rounded-lg text-sm text-text-dim hover:text-white hover:bg-white/10 transition-colors"
                                                >
                                                    <Download className="w-3 h-3" />
                                                    {file.split("/").pop()}
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Resolution */}
                                {dispute.resolution && (
                                    <div className="bg-green-500/10 border border-green-500/20 rounded-xl p-4 mb-4">
                                        <div className="flex items-center gap-2 mb-2">
                                            <CheckCircle2 className="w-4 h-4 text-green-500" />
                                            <span className="text-sm font-semibold text-green-500">Resolution</span>
                                        </div>
                                        <p className="text-white font-medium mb-1">{dispute.resolution}</p>
                                        {dispute.resolution_notes && (
                                            <p className="text-text-dim text-sm">{dispute.resolution_notes}</p>
                                        )}
                                        {dispute.resolved_at && (
                                            <p className="text-text-dim text-xs mt-2">
                                                Resolved on {formatDate(dispute.resolved_at)}
                                            </p>
                                        )}
                                    </div>
                                )}

                                {/* Actions */}
                                <div className="flex items-center justify-between pt-4 border-t border-white/10">
                                    <div className="flex items-center gap-4 text-xs text-text-dim">
                                        <span>Created {formatDate(dispute.created_at)}</span>
                                        {dispute.updated_at !== dispute.created_at && (
                                            <span>Updated {formatDate(dispute.updated_at)}</span>
                                        )}
                                    </div>
                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => router.push(`/dashboard/trades/${dispute.trade_id}`)}
                                            className="px-3 py-2 bg-white/5 border border-white/10 rounded-lg text-sm text-text-dim hover:text-white hover:bg-white/10 transition-colors flex items-center gap-2"
                                        >
                                            <Eye className="w-3 h-3" />
                                            View Trade
                                        </button>
                                        <button
                                            onClick={() => router.push(`/dashboard/disputes/${dispute.dispute_id}`)}
                                            className="px-3 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 transition-colors flex items-center gap-2"
                                        >
                                            <MessageSquare className="w-3 h-3" />
                                            View Details
                                        </button>
                                    </div>
                                </div>
                            </motion.div>
                        ))}
                    </div>
                )}
            </div>
        </DashboardLayout>
    );
};

export default UserDisputes;
