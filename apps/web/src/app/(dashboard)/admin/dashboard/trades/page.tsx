"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    TrendingUp,
    Search,
    Filter,
    Loader2,
    AlertTriangle,
    CheckCircle2,
    XCircle,
    Clock,
    DollarSign,
    Users,
    Eye,
    MoreVertical,
    Calendar,
    ArrowUpRight,
    ArrowDownRight
} from "lucide-react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { toast } from "sonner";

interface Trade {
    trade_id: string;
    ad_id: string;
    buyer_id: string;
    seller_id: string;
    buyer_username: string;
    seller_username: string;
    amount: number;
    price: number;
    total_amount: number;
    status: 'pending' | 'paid' | 'completed' | 'cancelled' | 'disputed';
    crypto_symbol: string;
    fiat_symbol: string;
    payment_method: string;
    payment_details?: string;
    dispute_id?: string;
    created_at: string;
    updated_at: string;
}

const AdminTradesPage = () => {
    const [trades, setTrades] = useState<Trade[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState("all");
    const [searchTerm, setSearchTerm] = useState("");
    const [selectedTrade, setSelectedTrade] = useState<Trade | null>(null);

    useEffect(() => {
        fetchTrades();
    }, [filter]);

    const fetchTrades = async () => {
        try {
            const params = new URLSearchParams();
            if (filter !== "all") {
                params.append("status", filter);
            }
            
            const response = await fetch(`/api/admin/trades?${params.toString()}`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch trades");
            }
            
            setTrades(data.trades || []);
        } catch (error) {
            console.error("Error fetching trades:", error);
            toast.error("Failed to load trades");
        } finally {
            setIsLoading(false);
        }
    };

    const filteredTrades = trades.filter(trade => {
        const matchesSearch = 
            trade.buyer_username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            trade.seller_username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            trade.crypto_symbol?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            trade.fiat_symbol?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            trade.trade_id?.toLowerCase().includes(searchTerm.toLowerCase());
        
        return matchesSearch;
    });

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'pending': return 'bg-yellow-500/20 text-yellow-400';
            case 'paid': return 'bg-blue-500/20 text-blue-400';
            case 'completed': return 'bg-green-500/20 text-green-400';
            case 'cancelled': return 'bg-red-500/20 text-red-400';
            case 'disputed': return 'bg-orange-500/20 text-orange-400';
            default: return 'bg-gray-500/20 text-gray-400';
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'pending': return Clock;
            case 'paid': return DollarSign;
            case 'completed': return CheckCircle2;
            case 'cancelled': return XCircle;
            case 'disputed': return AlertTriangle;
            default: return Clock;
        }
    };

    const formatCurrency = (amount: number, symbol: string) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: symbol,
        }).format(amount);
    };

    const calculateVolume = () => {
        return trades.reduce((sum, trade) => sum + trade.total_amount, 0);
    };

    const getTradeStats = () => {
        const stats = trades.reduce((acc, trade) => {
            acc[trade.status] = (acc[trade.status] || 0) + 1;
            return acc;
        }, {} as Record<string, number>);

        return {
            total: trades.length,
            pending: stats.pending || 0,
            paid: stats.paid || 0,
            completed: stats.completed || 0,
            cancelled: stats.cancelled || 0,
            disputed: stats.disputed || 0,
        };
    };

    if (isLoading) {
        return (
            <DashboardLayout title="Trade Monitoring" role="admin">
                <div className="flex items-center justify-center h-64">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
            </DashboardLayout>
        );
    }

    const stats = getTradeStats();

    return (
        <DashboardLayout title="Trade Monitoring" role="admin">
            <div className="space-y-6">
                {/* Stats Overview */}
                <div className="grid grid-cols-1 md:grid-cols-6 gap-4">
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Total</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.total}</p>
                            </div>
                            <TrendingUp className="w-6 h-6 text-blue-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Pending</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.pending}</p>
                            </div>
                            <Clock className="w-6 h-6 text-yellow-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Paid</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.paid}</p>
                            </div>
                            <DollarSign className="w-6 h-6 text-blue-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Completed</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.completed}</p>
                            </div>
                            <CheckCircle2 className="w-6 h-6 text-green-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Cancelled</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.cancelled}</p>
                            </div>
                            <XCircle className="w-6 h-6 text-red-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-xs font-medium">Disputed</p>
                                <p className="text-xl font-bold text-white mt-1">{stats.disputed}</p>
                            </div>
                            <AlertTriangle className="w-6 h-6 text-orange-500 opacity-50" />
                        </div>
                    </div>
                </div>

                {/* Volume Overview */}
                <div className="bg-surface border border-white/10 rounded-2xl p-6">
                    <div className="flex items-center justify-between">
                        <div>
                            <h3 className="text-lg font-bold text-white">Total Volume</h3>
                            <p className="text-2xl font-bold text-primary mt-2">
                                {formatCurrency(calculateVolume(), 'USD')}
                            </p>
                        </div>
                        <div className="flex items-center space-x-2 text-green-400">
                            <ArrowUpRight className="w-5 h-5" />
                            <span className="text-sm font-medium">+12.5%</span>
                        </div>
                    </div>
                </div>

                {/* Filters and Search */}
                <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search by users, crypto, or trade ID..."
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
                        <option value="all">All Status</option>
                        <option value="pending">Pending</option>
                        <option value="paid">Paid</option>
                        <option value="completed">Completed</option>
                        <option value="cancelled">Cancelled</option>
                        <option value="disputed">Disputed</option>
                    </select>
                </div>

                {/* Trades Table */}
                <div className="bg-surface border border-white/10 rounded-2xl overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-white/5">
                                <tr>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Trade ID
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Participants
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Amount
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Price
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
                                    {filteredTrades.map((trade) => {
                                        const StatusIcon = getStatusIcon(trade.status);
                                        return (
                                            <motion.tr
                                                key={trade.trade_id}
                                                initial={{ opacity: 0 }}
                                                animate={{ opacity: 1 }}
                                                exit={{ opacity: 0 }}
                                                className="hover:bg-white/5 transition-colors"
                                            >
                                                <td className="px-6 py-4">
                                                    <div className="text-sm font-medium text-white">
                                                        {trade.trade_id.slice(0, 8)}...
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        {trade.crypto_symbol}/{trade.fiat_symbol}
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="space-y-1">
                                                        <div className="flex items-center space-x-2">
                                                            <span className="text-xs text-green-400">Buyer:</span>
                                                            <span className="text-xs text-white">{trade.buyer_username}</span>
                                                        </div>
                                                        <div className="flex items-center space-x-2">
                                                            <span className="text-xs text-red-400">Seller:</span>
                                                            <span className="text-xs text-white">{trade.seller_username}</span>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div>
                                                        <div className="text-sm font-medium text-white">
                                                            {trade.amount.toFixed(6)} {trade.crypto_symbol}
                                                        </div>
                                                        <div className="text-xs text-text-dim">
                                                            {formatCurrency(trade.total_amount, trade.fiat_symbol)}
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="text-sm font-medium text-white">
                                                        {formatCurrency(trade.price, trade.fiat_symbol)}
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        per {trade.crypto_symbol}
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="flex items-center space-x-2">
                                                        <StatusIcon className="w-4 h-4 text-text-dim" />
                                                        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(trade.status)}`}>
                                                            {trade.status}
                                                        </span>
                                                    </div>
                                                    {trade.dispute_id && (
                                                        <div className="text-xs text-orange-400 mt-1">
                                                            Dispute: {trade.dispute_id.slice(0, 8)}...
                                                        </div>
                                                    )}
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="text-sm text-text-dim">
                                                        {new Date(trade.created_at).toLocaleDateString()}
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        {new Date(trade.created_at).toLocaleTimeString()}
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="flex items-center space-x-2">
                                                        <button
                                                            onClick={() => setSelectedTrade(trade)}
                                                            className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors"
                                                            title="View Details"
                                                        >
                                                            <Eye className="w-4 h-4" />
                                                        </button>
                                                        <button
                                                            className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors"
                                                            title="More Options"
                                                        >
                                                            <MoreVertical className="w-4 h-4" />
                                                        </button>
                                                    </div>
                                                </td>
                                            </motion.tr>
                                        );
                                    })}
                                </AnimatePresence>
                            </tbody>
                        </table>
                        
                        {filteredTrades.length === 0 && (
                            <div className="text-center py-12">
                                <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                <p className="text-text-dim">No trades found</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Trade Details Modal */}
                {selectedTrade && (
                    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                        <div className="bg-surface border border-white/10 rounded-2xl p-6 max-w-2xl w-full mx-4">
                            <div className="flex items-center justify-between mb-4">
                                <h3 className="text-lg font-bold text-white">Trade Details</h3>
                                <button
                                    onClick={() => setSelectedTrade(null)}
                                    className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20"
                                >
                                    <XCircle className="w-4 h-4" />
                                </button>
                            </div>
                            
                            <div className="space-y-4">
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-xs text-text-dim">Trade ID</p>
                                        <p className="text-sm text-white">{selectedTrade.trade_id}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-text-dim">Status</p>
                                        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(selectedTrade.status)}`}>
                                            {selectedTrade.status}
                                        </span>
                                    </div>
                                </div>
                                
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-xs text-text-dim">Buyer</p>
                                        <p className="text-sm text-white">{selectedTrade.buyer_username}</p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-text-dim">Seller</p>
                                        <p className="text-sm text-white">{selectedTrade.seller_username}</p>
                                    </div>
                                </div>
                                
                                <div className="grid grid-cols-3 gap-4">
                                    <div>
                                        <p className="text-xs text-text-dim">Amount</p>
                                        <p className="text-sm text-white">
                                            {selectedTrade.amount.toFixed(6)} {selectedTrade.crypto_symbol}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-text-dim">Price</p>
                                        <p className="text-sm text-white">
                                            {formatCurrency(selectedTrade.price, selectedTrade.fiat_symbol)}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-text-dim">Total</p>
                                        <p className="text-sm text-white">
                                            {formatCurrency(selectedTrade.total_amount, selectedTrade.fiat_symbol)}
                                        </p>
                                    </div>
                                </div>
                                
                                <div>
                                    <p className="text-xs text-text-dim">Payment Method</p>
                                    <p className="text-sm text-white">{selectedTrade.payment_method}</p>
                                </div>
                                
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-xs text-text-dim">Created</p>
                                        <p className="text-sm text-white">
                                            {new Date(selectedTrade.created_at).toLocaleString()}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-xs text-text-dim">Updated</p>
                                        <p className="text-sm text-white">
                                            {new Date(selectedTrade.updated_at).toLocaleString()}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </DashboardLayout>
    );
};

export default AdminTradesPage;
