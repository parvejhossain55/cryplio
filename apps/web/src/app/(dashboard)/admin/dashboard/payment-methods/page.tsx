"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff } from "lucide-react";

interface PaymentMethod {
    id: number;
    code: string;
    name: string;
    category: string;
    icon_url?: string;
    description?: string;
    is_active: boolean;
    sort_order: number;
    created_at: string;
}

const AdminPaymentMethods = () => {
    const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [showInactive, setShowInactive] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingMethod, setEditingMethod] = useState<PaymentMethod | null>(null);

    useEffect(() => {
        fetchPaymentMethods();
    }, []);

    const fetchPaymentMethods = async () => {
        try {
            const response = await fetch(`/api/v1/admin/payment-methods?active_only=false`);
            if (response.ok) {
                const data = await response.json();
                setPaymentMethods(data.payment_methods || []);
            }
        } catch (error) {
            console.error("Failed to fetch payment methods:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async (methodData: Partial<PaymentMethod>) => {
        try {
            const response = await fetch("/api/v1/admin/payment-methods", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(methodData),
            });

            if (response.ok) {
                await fetchPaymentMethods();
                setIsModalOpen(false);
                setEditingMethod(null);
            }
        } catch (error) {
            console.error("Failed to create payment method:", error);
        }
    };

    const handleUpdate = async (id: number, methodData: Partial<PaymentMethod>) => {
        try {
            const response = await fetch(`/api/v1/admin/payment-methods/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(methodData),
            });

            if (response.ok) {
                await fetchPaymentMethods();
                setIsModalOpen(false);
                setEditingMethod(null);
            }
        } catch (error) {
            console.error("Failed to update payment method:", error);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this payment method?")) return;

        try {
            const response = await fetch(`/api/v1/admin/payment-methods/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (response.ok) {
                await fetchPaymentMethods();
            }
        } catch (error) {
            console.error("Failed to delete payment method:", error);
        }
    };

    const filteredMethods = paymentMethods.filter(method => {
        const matchesSearch = method.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                            method.code.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesActive = showInactive || method.is_active;
        return matchesSearch && matchesActive;
    });

    const getCategoryColor = (category: string) => {
        switch (category) {
            case 'mobile_money': return 'bg-blue-500/20 text-blue-500';
            case 'bank_transfer': return 'bg-green-500/20 text-green-500';
            case 'online_wallet': return 'bg-purple-500/20 text-purple-500';
            case 'crypto': return 'bg-orange-500/20 text-orange-500';
            case 'cash': return 'bg-gray-500/20 text-gray-500';
            default: return 'bg-gray-500/20 text-gray-500';
        }
    };

    if (loading) {
        return (
            <DashboardLayout title="Payment Methods" role="admin">
                <div className="flex items-center justify-center h-64">
                    <div className="text-white">Loading...</div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Payment Methods" role="admin">
            <div className="space-y-8">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-black text-white">Payment Methods Management</h1>
                        <p className="text-text-dim mt-1">Manage supported payment methods for the platform</p>
                    </div>
                    <button
                        onClick={() => setIsModalOpen(true)}
                        className="flex items-center px-6 py-3 bg-primary text-white rounded-xl font-bold hover:scale-105 active:scale-95 transition-all"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Add Payment Method
                    </button>
                </div>

                {/* Filters */}
                <div className="flex items-center space-x-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search payment methods..."
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
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Category</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Status</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Order</th>
                                    <th className="px-6 py-4 text-right text-xs font-black text-text-dim uppercase tracking-widest">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredMethods.map((method) => (
                                    <tr key={method.id} className="border-t border-white/5 hover:bg-white/2 transition-all">
                                        <td className="px-6 py-4">
                                            <code className="text-xs font-bold text-primary bg-primary/10 px-2 py-1 rounded">
                                                {method.code}
                                            </code>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="font-bold text-white">{method.name}</div>
                                            {method.description && (
                                                <div className="text-xs text-text-dim mt-1">{method.description}</div>
                                            )}
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${getCategoryColor(method.category)}`}>
                                                {method.category.replace('_', ' ')}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${
                                                method.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
                                            }`}>
                                                {method.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-white font-bold">{method.sort_order}</td>
                                        <td className="px-6 py-4 text-right space-x-2">
                                            <button
                                                onClick={() => {
                                                    setEditingMethod(method);
                                                    setIsModalOpen(true);
                                                }}
                                                className="inline-flex items-center px-3 py-2 bg-blue-500/20 text-blue-500 rounded-lg hover:bg-blue-500/30 transition-all"
                                            >
                                                <Edit className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleDelete(method.id)}
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

                {filteredMethods.length === 0 && (
                    <div className="text-center py-12 text-text-dim">
                        No payment methods found matching your criteria.
                    </div>
                )}
            </div>

            {/* Modal will be added next */}
        </DashboardLayout>
    );
};

export default AdminPaymentMethods;