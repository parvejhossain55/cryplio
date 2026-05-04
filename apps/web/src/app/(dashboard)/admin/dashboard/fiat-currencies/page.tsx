"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff } from "lucide-react";

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

    useEffect(() => {
        fetchFiatCurrencies();
    }, []);

    const fetchFiatCurrencies = async () => {
        try {
            const response = await fetch(`/api/v1/admin/fiat-currencies?active_only=false`);
            if (response.ok) {
                const data = await response.json();
                setFiatCurrencies(data.fiat_currencies || []);
            }
        } catch (error) {
            console.error("Failed to fetch fiat currencies:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async (currencyData: Partial<FiatCurrency>) => {
        try {
            const response = await fetch("/api/v1/admin/fiat-currencies", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(currencyData),
            });

            if (response.ok) {
                await fetchFiatCurrencies();
                setIsModalOpen(false);
                setEditingCurrency(null);
            }
        } catch (error) {
            console.error("Failed to create fiat currency:", error);
        }
    };

    const handleUpdate = async (id: number, currencyData: Partial<FiatCurrency>) => {
        try {
            const response = await fetch(`/api/v1/admin/fiat-currencies/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(currencyData),
            });

            if (response.ok) {
                await fetchFiatCurrencies();
                setIsModalOpen(false);
                setEditingCurrency(null);
            }
        } catch (error) {
            console.error("Failed to update fiat currency:", error);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this fiat currency? This may affect existing trades.")) return;

        try {
            const response = await fetch(`/api/v1/admin/fiat-currencies/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (response.ok) {
                await fetchFiatCurrencies();
            }
        } catch (error) {
            console.error("Failed to delete fiat currency:", error);
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
                        className={`flex items-center px-4 py-3 rounded-xl font-medium transition-all ${
                            showInactive ? 'bg-primary/20 text-primary' : 'bg-white/5 text-text-dim'
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
                                {filteredCurrencies.map((currency) => (
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
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${
                                                currency.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
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

                {filteredCurrencies.length === 0 && (
                    <div className="text-center py-12 text-text-dim">
                        No fiat currencies found matching your criteria.
                    </div>
                )}
            </div>

            {/* Modal will be added next */}
        </DashboardLayout>
    );
};

export default AdminFiatCurrencies;