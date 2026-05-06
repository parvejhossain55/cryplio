"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    AlertTriangle,
    ShieldAlert,
    Gavel,
    User,
    Clock,
    CheckCircle2,
    XCircle,
    ChevronRight,
    Loader2,
    Search,
    Filter
} from "lucide-react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { authService } from "@/services/authService";
import { toast } from "sonner";
import Link from "next/link";

const AdminDisputesPage = () => {
    const [disputes, setDisputes] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState("all");

    useEffect(() => {
        fetchDisputes();
    }, []);

    const fetchDisputes = async () => {
        setIsLoading(true);
        try {
            const data = await authService.getAdminDisputes();
            setDisputes(data || []);
        } catch (err: any) {
            console.error(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleAssign = async (disputeId: string) => {
        try {
            await authService.assignDispute(disputeId);
            toast.success("Dispute assigned to you");
            fetchDisputes();
        } catch (err: any) {
            toast.error(err.message || "Failed to assign dispute");
        }
    };

    const handleResolve = async (disputeId: string, resolution: string, winnerId: 'buyer' | 'seller') => {
        toast(`Resolve in favor of ${winnerId === 'buyer' ? 'Buyer' : 'Seller'}?`, {
            description: "This action will release or refund the escrowed assets. It cannot be undone.",
            action: {
                label: 'Confirm Resolve',
                onClick: async () => {
                    try {
                        const dispute = disputes.find(d => d.dispute_id === disputeId);
                        await authService.resolveDispute(disputeId, resolution, dispute.raised_by);
                        toast.success("Dispute resolved successfully");
                        fetchDisputes();
                    } catch (err: any) {
                        toast.error(err.message || "Failed to resolve dispute");
                    }
                }
            },
            cancel: {
                label: 'Cancel',
                onClick: () => { }
            }
        });
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'open': return 'text-primary bg-primary/10 border-primary/20';
            case 'in_review': return 'text-blue-500 bg-blue-500/10 border-blue-500/20';
            case 'resolved': return 'text-accent bg-accent/10 border-accent/20';
            default: return 'text-text-dim bg-white/5 border-white/10';
        }
    };

    return (
        <DashboardLayout title="Dispute Management" role="admin">
            <div className="space-y-8 pb-10">
                {/* Header */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                    <div>
                        <h1 className="text-3xl font-black text-white italic uppercase tracking-tight">Active <span className="gradient-text">Conflicts</span></h1>
                        <p className="text-[10px] text-text-dim mt-1 font-medium italic">Protocol-level arbitration for peer-to-peer operations.</p>
                    </div>
                    <div className="flex items-center gap-4">
                        <div className="flex items-center gap-2 bg-surface/50 p-1 rounded-xl border border-white/5">
                            {["all", "open", "in_review", "resolved"].map((f) => (
                                <button
                                    key={f}
                                    onClick={() => setFilter(f)}
                                    className={`px-4 py-2 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all ${filter === f ? "bg-white text-background" : "text-text-dim hover:text-white"}`}
                                >
                                    {f}
                                </button>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Disputes List */}
                <div className="grid grid-cols-1 gap-4">
                    {isLoading ? (
                        <div className="flex flex-col items-center justify-center py-24 space-y-4">
                            <Loader2 className="w-10 h-10 text-primary animate-spin" />
                            <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">Retrieving Conflict Data...</p>
                        </div>
                    ) : disputes.length === 0 ? (
                        <div className="bg-surface border border-white/5 rounded-[2.5rem] p-20 text-center space-y-4">
                            <div className="w-16 h-16 bg-white/5 rounded-full flex items-center justify-center mx-auto">
                                <ShieldAlert className="w-8 h-8 text-text-dim/20" />
                            </div>
                            <h3 className="text-xl font-black text-white italic uppercase tracking-tight">System Integrity Confirmed</h3>
                            <p className="text-xs text-text-dim max-w-xs mx-auto">No active disputes requiring administrative intervention at this time.</p>
                        </div>
                    ) : (
                        disputes.filter(d => filter === 'all' || d.status === filter).map((dispute, i) => (
                            <motion.div
                                key={dispute.dispute_id}
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                transition={{ delay: i * 0.05 }}
                                className="bg-surface border border-white/5 rounded-[1.5rem] p-6 hover:border-white/10 transition-all group"
                            >
                                <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                                    <div className="flex items-center gap-6">
                                        <div className={`p-4 rounded-2xl bg-white/5 ${dispute.status === 'open' ? "text-primary" : "text-accent"}`}>
                                            <AlertTriangle className="w-6 h-6" />
                                        </div>
                                        <div className="space-y-1">
                                            <div className="flex items-center gap-2">
                                                <span className="text-[10px] font-black uppercase tracking-widest text-text-dim/50">ORDER #{dispute.trade_id.substring(0, 8)}</span>
                                                <span className={`px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border ${getStatusColor(dispute.status)}`}>
                                                    {dispute.status}
                                                </span>
                                            </div>
                                            <h3 className="text-lg font-black text-white">
                                                {dispute.reason_text || dispute.reason_code}
                                            </h3>
                                            <div className="flex items-center gap-3">
                                                <p className="text-[9px] text-text-dim font-bold uppercase tracking-widest flex items-center gap-1.5">
                                                    <User className="w-3 h-3" />
                                                    INITIATED BY: {dispute.raised_by.substring(0, 8)}
                                                </p>
                                                <div className="w-1 h-1 rounded-full bg-white/10" />
                                                <p className="text-[9px] text-text-dim font-bold uppercase tracking-widest flex items-center gap-1.5">
                                                    <Clock className="w-3 h-3" />
                                                    {new Date(dispute.created_at).toLocaleString()}
                                                </p>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-3">
                                        {dispute.status === 'pending' && (
                                            <button
                                                onClick={() => handleAssign(dispute.dispute_id)}
                                                className="px-6 py-3 bg-primary text-white rounded-xl text-[10px] font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/10"
                                            >
                                                Assign to Me
                                            </button>
                                        )}
                                        {dispute.status === 'assigned' && (
                                            <div className="flex items-center gap-2">
                                                <button
                                                    onClick={() => handleResolve(dispute.dispute_id, 'release_to_buyer', 'buyer')}
                                                    className="px-4 py-3 bg-accent text-white rounded-xl text-[8px] font-black uppercase tracking-widest hover:bg-accent/80 transition-all"
                                                >
                                                    Refund Buyer
                                                </button>
                                                <button
                                                    onClick={() => handleResolve(dispute.dispute_id, 'return_to_seller', 'seller')}
                                                    className="px-4 py-3 bg-white text-background rounded-xl text-[8px] font-black uppercase tracking-widest hover:bg-white/80 transition-all"
                                                >
                                                    Release Seller
                                                </button>
                                            </div>
                                        )}
                                        <Link
                                            href={`/marketplace/trade/${dispute.trade_id}`}
                                            className="p-3 rounded-xl bg-white/5 text-text-dim hover:text-white border border-white/5 transition-all"
                                        >
                                            <ChevronRight className="w-5 h-5" />
                                        </Link>
                                    </div>
                                </div>
                            </motion.div>
                        ))
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminDisputesPage;
