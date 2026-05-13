"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    DollarSign,
    CheckCircle2,
    XCircle,
    Clock,
    Search,
    Filter,
    Loader2,
    ExternalLink,
    AlertTriangle,
    TrendingUp,
    Eye,
    EyeOff
} from "lucide-react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { toast } from "sonner";

interface Withdrawal {
    txn_id: string;
    user_id: string;
    username: string;
    email: string;
    amount: number;
    fee: number;
    net_amount: number;
    crypto_symbol: string;
    destination_address: string;
    status: string;
    requires_approval: boolean;
    created_at: string;
    updated_at: string;
}

const AdminWithdrawalsPage = () => {
    const [withdrawals, setWithdrawals] = useState<Withdrawal[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState("pending");
    const [searchTerm, setSearchTerm] = useState("");
    const [showAddresses, setShowAddresses] = useState<Record<string, boolean>>({});

    useEffect(() => {
        fetchWithdrawals();
    }, []);

    const fetchWithdrawals = async () => {
        try {
            const response = await fetch("/api/v1/admin/withdrawals/pending");
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch withdrawals");
            }
            
            setWithdrawals(data.withdrawals || []);
        } catch (error) {
            console.error("Error fetching withdrawals:", error);
            toast.error("Failed to load withdrawals");
        } finally {
            setIsLoading(false);
        }
    };

    const [isApproveModalOpen, setIsApproveModalOpen] = useState(false);
    const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);
    const [selectedTxnId, setSelectedTxnId] = useState<string | null>(null);
    const [txHash, setTxHash] = useState("");
    const [rejectReason, setRejectReason] = useState("");

    const handleApprove = async () => {
        if (!selectedTxnId || !txHash) return;

        try {
            const response = await fetch(`/api/v1/admin/withdrawals/${selectedTxnId}/approve`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ tx_hash: txHash }),
            });

            if (response.ok) {
                toast.success("Withdrawal approved");
                setIsApproveModalOpen(false);
                setTxHash("");
                fetchWithdrawals();
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to approve withdrawal");
            }
        } catch (error) {
            console.error("Failed to approve withdrawal:", error);
            toast.error("Failed to approve withdrawal");
        }
    };

    const handleReject = async () => {
        if (!selectedTxnId || !rejectReason) return;

        try {
            const response = await fetch(`/api/v1/admin/withdrawals/${selectedTxnId}/reject`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ reason: rejectReason }),
            });

            if (response.ok) {
                toast.success("Withdrawal rejected");
                setIsRejectModalOpen(false);
                setRejectReason("");
                fetchWithdrawals();
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to reject withdrawal");
            }
        } catch (error) {
            console.error("Failed to reject withdrawal:", error);
            toast.error("Failed to reject withdrawal");
        }
    };

    const openApproveModal = (txnId: string) => {
        setSelectedTxnId(txnId);
        setIsApproveModalOpen(true);
    };

    const openRejectModal = (txnId: string) => {
        setSelectedTxnId(txnId);
        setIsRejectModalOpen(true);
    };

    const toggleAddressVisibility = (txnId: string) => {
        setShowAddresses(prev => ({
            ...prev,
            [txnId]: !prev[txnId]
        }));
    };

    const filteredWithdrawals = withdrawals.filter(withdrawal => {
        const matchesFilter = filter === "all" || withdrawal.status === filter;
        const matchesSearch = 
            withdrawal.username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            withdrawal.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            withdrawal.crypto_symbol?.toLowerCase().includes(searchTerm.toLowerCase());
        
        return matchesFilter && matchesSearch;
    });

    const formatAddress = (address: string) => {
        if (!address) return "N/A";
        return `${address.slice(0, 6)}...${address.slice(-4)}`;
    };

    const formatCurrency = (amount: number, symbol: string) => {
        return `${amount.toFixed(8)} ${symbol}`;
    };

    if (isLoading) {
        return (
            <DashboardLayout title="Withdrawal Management" role="admin">
                <div className="flex items-center justify-center h-64">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Withdrawal Management" role="admin">
            <div className="space-y-6">
                {/* Stats Overview */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Pending</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {withdrawals.filter(w => w.status === 'pending').length}
                                </p>
                            </div>
                            <Clock className="w-8 h-8 text-yellow-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Total Amount</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {withdrawals.reduce((sum, w) => sum + w.amount, 0).toFixed(2)}
                                </p>
                            </div>
                            <DollarSign className="w-8 h-8 text-green-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Approved Today</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {withdrawals.filter(w => w.status === 'approved').length}
                                </p>
                            </div>
                            <CheckCircle2 className="w-8 h-8 text-green-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Rejected Today</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {withdrawals.filter(w => w.status === 'rejected').length}
                                </p>
                            </div>
                            <XCircle className="w-8 h-8 text-red-500 opacity-50" />
                        </div>
                    </div>
                </div>

                {/* Filters and Search */}
                <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search by username, email, or crypto..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="w-full pl-10 pr-4 py-3 bg-surface border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                        />
                    </div>
                    
                    <select
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        className="px-4 py-3 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                    >
                        <option value="pending">Pending</option>
                        <option value="approved">Approved</option>
                        <option value="rejected">Rejected</option>
                        <option value="all">All</option>
                    </select>
                </div>

                {/* Withdrawals Table */}
                <div className="bg-surface border border-white/10 rounded-2xl overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-white/5">
                                <tr>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        User
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Amount
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Destination
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Status
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Date
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                <AnimatePresence>
                                    {filteredWithdrawals.map((withdrawal) => (
                                        <motion.tr
                                            key={withdrawal.txn_id}
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                            className="hover:bg-white/5 transition-colors"
                                        >
                                            <td className="px-6 py-4">
                                                <div>
                                                    <div className="text-sm font-medium text-white">
                                                        {withdrawal.username || 'Unknown'}
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        {withdrawal.email || 'No email'}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div>
                                                    <div className="text-sm font-medium text-white">
                                                        {formatCurrency(withdrawal.amount, withdrawal.crypto_symbol)}
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        Fee: {formatCurrency(withdrawal.fee, withdrawal.crypto_symbol)}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="flex items-center space-x-2">
                                                    <span className="text-sm text-white">
                                                        {showAddresses[withdrawal.txn_id] 
                                                            ? withdrawal.destination_address 
                                                            : formatAddress(withdrawal.destination_address)}
                                                    </span>
                                                    <button
                                                        onClick={() => toggleAddressVisibility(withdrawal.txn_id)}
                                                        className="text-text-dim hover:text-white transition-colors"
                                                    >
                                                        {showAddresses[withdrawal.txn_id] ? (
                                                            <EyeOff className="w-4 h-4" />
                                                        ) : (
                                                            <Eye className="w-4 h-4" />
                                                        )}
                                                    </button>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                                                    withdrawal.status === 'pending' ? 'bg-yellow-500/20 text-yellow-400' :
                                                    withdrawal.status === 'approved' ? 'bg-green-500/20 text-green-400' :
                                                    withdrawal.status === 'rejected' ? 'bg-red-500/20 text-red-400' :
                                                    'bg-gray-500/20 text-gray-400'
                                                }`}>
                                                    {withdrawal.status}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="text-sm text-text-dim">
                                                    {new Date(withdrawal.created_at).toLocaleDateString()}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="flex items-center space-x-2 relative z-10">
                                                    {withdrawal.status === 'pending' && (
                                                        <>
                                                            <button
                                                                onClick={(e) => {
                                                                    e.stopPropagation();
                                                                    openApproveModal(withdrawal.txn_id);
                                                                }}
                                                                className="p-2 bg-green-500/20 text-green-400 rounded-lg hover:bg-green-500/30 transition-colors pointer-events-auto"
                                                                title="Approve Withdrawal"
                                                            >
                                                                <CheckCircle2 className="w-4 h-4" />
                                                            </button>
                                                            <button
                                                                onClick={(e) => {
                                                                    e.stopPropagation();
                                                                    openRejectModal(withdrawal.txn_id);
                                                                }}
                                                                className="p-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors pointer-events-auto"
                                                                title="Reject Withdrawal"
                                                            >
                                                                <XCircle className="w-4 h-4" />
                                                            </button>
                                                        </>
                                                    )}
                                                    <button
                                                        onClick={(e) => e.stopPropagation()}
                                                        className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors pointer-events-auto"
                                                        title="View Details"
                                                    >
                                                        <ExternalLink className="w-4 h-4" />
                                                    </button>
                                                </div>
                                            </td>
                                        </motion.tr>
                                    ))}
                                </AnimatePresence>
                            </tbody>
                        </table>
                        
                        {filteredWithdrawals.length === 0 && (
                            <div className="text-center py-12">
                                <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                <p className="text-text-dim">No withdrawals found</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* Approve Withdrawal Modal */}
            <AnimatePresence>
                {isApproveModalOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-surface border border-white/10 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl"
                        >
                            <div className="p-6 border-b border-white/5 flex items-center justify-between">
                                <h3 className="text-xl font-bold text-white">Approve Withdrawal</h3>
                                <button
                                    onClick={() => setIsApproveModalOpen(false)}
                                    className="p-2 hover:bg-white/5 rounded-lg transition-colors"
                                >
                                    <XCircle className="w-5 h-5 text-text-dim" />
                                </button>
                            </div>

                            <div className="p-6 space-y-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-text-dim">Transaction Hash (TXID)</label>
                                    <input
                                        type="text"
                                        value={txHash}
                                        onChange={(e) => setTxHash(e.target.value)}
                                        placeholder="Enter the blockchain transaction hash..."
                                        className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-text-dim/50 focus:outline-none focus:border-primary/50 font-mono text-sm"
                                    />
                                </div>

                                <div className="p-4 bg-green-500/10 border border-green-500/20 rounded-2xl flex gap-3">
                                    <CheckCircle2 className="w-5 h-5 text-green-500 shrink-0 mt-0.5" />
                                    <p className="text-xs text-green-200/80 leading-relaxed">
                                        Please ensure the transaction hash is correct. This will mark the withdrawal as completed and notify the user.
                                    </p>
                                </div>
                            </div>

                            <div className="p-6 bg-white/5 flex gap-3">
                                <button
                                    onClick={() => setIsApproveModalOpen(false)}
                                    className="flex-1 px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white font-bold hover:bg-white/10 transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleApprove}
                                    disabled={!txHash}
                                    className="flex-1 px-4 py-3 bg-green-500 text-white rounded-xl font-bold hover:bg-green-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    Complete Approval
                                </button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>

            {/* Reject Withdrawal Modal */}
            <AnimatePresence>
                {isRejectModalOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-surface border border-white/10 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl"
                        >
                            <div className="p-6 border-b border-white/5 flex items-center justify-between">
                                <h3 className="text-xl font-bold text-white">Reject Withdrawal</h3>
                                <button
                                    onClick={() => setIsRejectModalOpen(false)}
                                    className="p-2 hover:bg-white/5 rounded-lg transition-colors"
                                >
                                    <XCircle className="w-5 h-5 text-text-dim" />
                                </button>
                            </div>

                            <div className="p-6 space-y-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-text-dim">Reason for Rejection</label>
                                    <textarea
                                        value={rejectReason}
                                        onChange={(e) => setRejectReason(e.target.value)}
                                        placeholder="e.g. Invalid destination address, Suspicious activity..."
                                        className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-text-dim/50 focus:outline-none focus:border-primary/50 min-h-[100px] resize-none"
                                    />
                                </div>

                                <div className="p-4 bg-red-500/10 border border-red-500/20 rounded-2xl flex gap-3">
                                    <AlertTriangle className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
                                    <p className="text-xs text-red-200/80 leading-relaxed">
                                        Rejecting this withdrawal will return the funds (minus fees) to the user's wallet.
                                    </p>
                                </div>
                            </div>

                            <div className="p-6 bg-white/5 flex gap-3">
                                <button
                                    onClick={() => setIsRejectModalOpen(false)}
                                    className="flex-1 px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white font-bold hover:bg-white/10 transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleReject}
                                    disabled={!rejectReason}
                                    className="flex-1 px-4 py-3 bg-red-500 text-white rounded-xl font-bold hover:bg-red-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    Confirm Rejection
                                </button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </DashboardLayout>
    );
};

export default AdminWithdrawalsPage;
