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

    const handleApprove = async (txnId: string) => {
        const txHash = prompt("Enter transaction hash:");
        if (!txHash) return;

        try {
            const response = await fetch(`/api/admin/withdrawals/${txnId}/approve`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ tx_hash: txHash }),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to approve withdrawal");
            }

            toast.success("Withdrawal approved successfully");
            fetchWithdrawals();
        } catch (error) {
            console.error("Error approving withdrawal:", error);
            toast.error("Failed to approve withdrawal");
        }
    };

    const handleReject = async (txnId: string) => {
        const reason = prompt("Enter rejection reason:");
        if (!reason) return;

        try {
            const response = await fetch(`/api/admin/withdrawals/${txnId}/reject`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ reason }),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to reject withdrawal");
            }

            toast.success("Withdrawal rejected successfully");
            fetchWithdrawals();
        } catch (error) {
            console.error("Error rejecting withdrawal:", error);
            toast.error("Failed to reject withdrawal");
        }
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
                                                <div className="flex items-center space-x-2">
                                                    {withdrawal.status === 'pending' && (
                                                        <>
                                                            <button
                                                                onClick={() => handleApprove(withdrawal.txn_id)}
                                                                className="p-2 bg-green-500/20 text-green-400 rounded-lg hover:bg-green-500/30 transition-colors"
                                                                title="Approve"
                                                            >
                                                                <CheckCircle2 className="w-4 h-4" />
                                                            </button>
                                                            <button
                                                                onClick={() => handleReject(withdrawal.txn_id)}
                                                                className="p-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
                                                                title="Reject"
                                                            >
                                                                <XCircle className="w-4 h-4" />
                                                            </button>
                                                        </>
                                                    )}
                                                    <button
                                                        className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors"
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
        </DashboardLayout>
    );
};

export default AdminWithdrawalsPage;
