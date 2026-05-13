"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff, X } from "lucide-react";
import { toast } from "sonner";

interface FiatCurrency {
    id: number;
    code: string;
    name: string;
    symbol: string;
    is_active: boolean;
    created_at: string;
}

const AdminFiatCurrencies = () => {
    const [fiatCurrencies, setFiatCurrencies] = useState<FiatCurrency[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [showInactive, setShowInactive] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingCurrency, setEditingCurrency] = useState<FiatCurrency | null>(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0,
        pages: 0,
    });

    useEffect(() => {
        const timer = setTimeout(() => {
            fetchFiatCurrencies(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchQuery, showInactive]);

    useEffect(() => {
        fetchFiatCurrencies(pagination.page);
    }, [pagination.page]);

    const fetchFiatCurrencies = async (page = 1) => {
        setLoading(true);
        try {
            const params = new URLSearchParams({
                active_only: showInactive ? 'false' : 'true',
                page: page.toString(),
                limit: pagination.limit.toString(),
            });

            if (searchQuery) {
                params.append('search', searchQuery);
            }

            const response = await fetch(`/api/v1/admin/fiat-currencies?${params}`);
            if (response.ok) {
                const data = await response.json();
                setFiatCurrencies(data.fiat_currencies || []);
                setPagination(data.pagination || { page: 1, limit: 50, total: 0, pages: 0 });
            }
        } catch (error) {
            console.error("Failed to fetch fiat currencies:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (currencyData: Partial<FiatCurrency>) => {
        try {
            const url = editingCurrency
                ? `/api/v1/admin/fiat-currencies/${editingCurrency.id}`
                : "/api/v1/admin/fiat-currencies";
            const method = editingCurrency ? "PUT" : "POST";

            const response = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(currencyData),
            });

            if (response.ok) {
                await fetchFiatCurrencies(pagination.page);
                setIsModalOpen(false);
                setEditingCurrency(null);
                toast.success(editingCurrency ? "Fiat currency updated" : "Fiat currency created");
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to save fiat currency");
            }
        } catch (error) {
            console.error("Failed to save fiat currency:", error);
            toast.error("Failed to save fiat currency");
        }
    };

    const handleDelete = async (id: number) => {
        const confirmDelete = window.confirm("Are you sure you want to delete this fiat currency? This may affect existing trades.");
        if (!confirmDelete) return;

        try {
            const response = await fetch(`/api/v1/admin/fiat-currencies/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (response.ok) {
                toast.success("Fiat currency deleted");
                await fetchFiatCurrencies(pagination.page);
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to delete fiat currency");
            }
        } catch (error) {
            console.error("Failed to delete fiat currency:", error);
            toast.error("Failed to delete fiat currency");
        }
    };

    const filteredCurrencies = fiatCurrencies.filter(currency => {
        const matchesSearch = currency.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            currency.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
            currency.symbol.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesActive = showInactive || currency.is_active;
        return matchesSearch && matchesActive;
    });

    if (loading) {
        return (
            <DashboardLayout title="Fiat Currencies" role="admin">
                <div className="flex items-center justify-center h-64">
                    <div className="text-white">Loading...</div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Fiat Currencies" role="admin">
            <div className="space-y-8">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-black text-white">Fiat Currencies Management</h1>
                        <p className="text-text-dim mt-1">Manage supported fiat currencies for trading</p>
                    </div>
                    <button
                        onClick={() => setIsModalOpen(true)}
                        className="flex items-center px-6 py-3 bg-primary text-white rounded-xl font-bold hover:scale-105 active:scale-95 transition-all"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Add Fiat Currency
                    </button>
                </div>

                {/* Filters */}
                <div className="flex items-center space-x-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search fiat currencies..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="w-full pl-10 pr-4 py-3 bg-surface border border-white/10 rounded-xl text-white placeholder-text-dim focus:border-primary/50 transition-all"
                        />
                    </div>
                    <button
                        onClick={() => setShowInactive(!showInactive)}
                        className={`flex items-center px-4 py-3 rounded-xl font-medium transition-all ${showInactive ? 'bg-primary/20 text-primary' : 'bg-white/5 text-text-dim'
                            }`}
                    >
                        {showInactive ? <Eye className="w-4 h-4 mr-2" /> : <EyeOff className="w-4 h-4 mr-2" />}
                        Show Inactive
                    </button>
                </div>

                {/* Table */}
                <div className="bg-surface border border-white/10 rounded-[2.5rem] overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-white/5">
                                <tr>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Code</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Name</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Symbol</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Status</th>
                                    <th className="px-6 py-4 text-right text-xs font-black text-text-dim uppercase tracking-widest">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {fiatCurrencies.map((currency) => (
                                    <tr key={currency.id} className="border-t border-white/5 hover:bg-white/2 transition-all">
                                        <td className="px-6 py-4">
                                            <code className="text-sm font-bold text-primary bg-primary/10 px-2 py-1 rounded">
                                                {currency.code}
                                            </code>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="font-bold text-white">{currency.name}</div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className="text-lg font-bold text-white">{currency.symbol}</span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${currency.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
                                                }`}>
                                                {currency.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-right space-x-2">
                                            <button
                                                onClick={() => {
                                                    setEditingCurrency(currency);
                                                    setIsModalOpen(true);
                                                }}
                                                className="inline-flex items-center px-3 py-2 bg-blue-500/20 text-blue-500 rounded-lg hover:bg-blue-500/30 transition-all"
                                            >
                                                <Edit className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleDelete(currency.id)}
                                                className="inline-flex items-center px-3 py-2 bg-red-500/20 text-red-500 rounded-lg hover:bg-red-500/30 transition-all"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {fiatCurrencies.length === 0 && (
                    <div className="text-center py-12 text-text-dim">
                        No fiat currencies found matching your criteria.
                    </div>
                )}

                {/* Pagination */}
                {pagination.pages > 1 && (
                    <div className="flex items-center justify-between pt-6 border-t border-white/5">
                        <div className="text-sm text-text-dim">
                            Showing {((pagination.page - 1) * pagination.limit) + 1} to {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} entries
                        </div>
                        <div className="flex items-center space-x-2">
                            <button
                                onClick={() => setPagination(prev => ({ ...prev, page: prev.page - 1 }))}
                                disabled={pagination.page <= 1}
                                className="px-3 py-2 bg-white/5 border border-white/5 rounded-lg text-sm font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-all"
                            >
                                Previous
                            </button>

                            {Array.from({ length: Math.min(5, pagination.pages) }, (_, i) => {
                                const pageNum = i + 1;
                                return (
                                    <button
                                        key={pageNum}
                                        onClick={() => setPagination(prev => ({ ...prev, page: pageNum }))}
                                        className={`px-3 py-2 rounded-lg text-sm font-bold transition-all ${pagination.page === pageNum
                                            ? 'bg-primary text-white'
                                            : 'bg-white/5 border border-white/5 text-text-dim hover:bg-white/10'
                                            }`}
                                    >
                                        {pageNum}
                                    </button>
                                );
                            })}

                            <button
                                onClick={() => setPagination(prev => ({ ...prev, page: prev.page + 1 }))}
                                disabled={pagination.page >= pagination.pages}
                                className="px-3 py-2 bg-white/5 border border-white/5 rounded-lg text-sm font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-all"
                            >
                                Next
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Fiat Currency Modal */}
            {isModalOpen && (
                <FiatCurrencyModal
                    currency={editingCurrency}
                    onClose={() => {
                        setIsModalOpen(false);
                        setEditingCurrency(null);
                    }}
                    onSave={handleSave}
                />
            )}
        </DashboardLayout>
    );
};

// Fiat Currency Modal Component
interface FiatCurrencyModalProps {
    currency: FiatCurrency | null;
    onClose: () => void;
    onSave: (data: Partial<FiatCurrency>) => void;
}

const FiatCurrencyModal: React.FC<FiatCurrencyModalProps> = ({ currency, onClose, onSave }) => {
    const [formData, setFormData] = useState({
        code: currency?.code || '',
        name: currency?.name || '',
        symbol: currency?.symbol || '',
        is_active: currency?.is_active ?? true,
    });
    const [saving, setSaving] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);

        try {
            await onSave(formData);
        } catch (error) {
            console.error('Failed to save fiat currency:', error);
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="p-8">
                    <div className="flex items-center justify-between mb-6">
                        <h2 className="text-xl font-black text-white">
                            {currency ? 'Edit Fiat Currency' : 'Add Fiat Currency'}
                        </h2>
                        <button
                            onClick={onClose}
                            className="p-2 hover:bg-white/5 rounded-xl transition-all"
                        >
                            <X className="w-5 h-5 text-text-dim" />
                        </button>
                    </div>

                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Code *
                            </label>
                            <input
                                type="text"
                                value={formData.code}
                                onChange={(e) => setFormData(prev => ({ ...prev, code: e.target.value.toUpperCase() }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="e.g., USD, EUR"
                                maxLength={3}
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Name *
                            </label>
                            <input
                                type="text"
                                value={formData.name}
                                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="e.g., US Dollar, Euro"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Symbol *
                            </label>
                            <input
                                type="text"
                                value={formData.symbol}
                                onChange={(e) => setFormData(prev => ({ ...prev, symbol: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="e.g., $, €"
                                maxLength={5}
                                required
                            />
                        </div>

                        <div className="flex items-center space-x-3">
                            <input
                                type="checkbox"
                                id="is_active"
                                checked={formData.is_active}
                                onChange={(e) => setFormData(prev => ({ ...prev, is_active: e.target.checked }))}
                                className="w-4 h-4 text-primary bg-white/5 border-white/5 rounded focus:ring-primary focus:ring-2"
                            />
                            <label htmlFor="is_active" className="text-sm font-bold text-white">
                                Active
                            </label>
                        </div>

                        <div className="flex space-x-4 pt-6">
                            <button
                                type="button"
                                onClick={onClose}
                                className="flex-1 py-3 px-6 bg-white/5 border border-white/5 rounded-xl text-sm font-bold text-text-dim hover:bg-white/10 transition-all"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={saving}
                                className="flex-1 py-3 px-6 bg-primary text-white rounded-xl text-sm font-bold hover:scale-105 active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {saving ? 'Saving...' : (currency ? 'Update' : 'Create')}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AdminFiatCurrencies;