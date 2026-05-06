"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff, X } from "lucide-react";
import { toast } from "sonner";

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
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 50,
        total: 0,
        pages: 0,
    });

    useEffect(() => {
        fetchPaymentMethods(pagination.page);
    }, [pagination.page, showInactive]);

    const fetchPaymentMethods = async (page = 1) => {
        setLoading(true);
        try {
            const params = new URLSearchParams({
                active_only: showInactive ? 'false' : 'true',
                page: page.toString(),
                limit: pagination.limit.toString(),
            });

            const response = await fetch(`/api/v1/admin/payment-methods?${params}`);
            if (response.ok) {
                const data = await response.json();
                setPaymentMethods(data.payment_methods || []);
                setPagination(data.pagination || { page: 1, limit: 50, total: 0, pages: 0 });
            }
        } catch (error) {
            console.error("Failed to fetch payment methods:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (methodData: Partial<PaymentMethod>) => {
        try {
            const url = editingMethod
                ? `/api/v1/admin/payment-methods/${editingMethod.id}`
                : "/api/v1/admin/payment-methods";
            const method = editingMethod ? "PUT" : "POST";

            const response = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(methodData),
            });

            if (response.ok) {
                await fetchPaymentMethods(pagination.page);
                setIsModalOpen(false);
                setEditingMethod(null);
                toast.success(editingMethod ? "Payment method updated" : "Payment method created");
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to save payment method");
            }
        } catch (error) {
            console.error("Failed to save payment method:", error);
            toast.error("Failed to save payment method");
        }
    };

    const handleDelete = async (id: number) => {
        toast("Delete Payment Method", {
            description: "Are you sure you want to delete this payment method?",
            action: {
                label: 'Delete',
                onClick: async () => {
                    try {
                        const response = await fetch(`/api/v1/admin/payment-methods/${id}`, {
                            method: "DELETE",
                            credentials: "include",
                        });

                        if (response.ok) {
                            toast.success("Payment method deleted");
                            await fetchPaymentMethods(pagination.page);
                        } else {
                            toast.error("Failed to delete payment method");
                        }
                    } catch (error) {
                        console.error("Failed to delete payment method:", error);
                        toast.error("Failed to delete payment method");
                    }
                }
            },
            cancel: {
                label: 'Cancel',
                onClick: () => { }
            }
        });
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
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${method.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
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

            {/* Payment Method Modal */}
            {isModalOpen && (
                <PaymentMethodModal
                    method={editingMethod}
                    onClose={() => {
                        setIsModalOpen(false);
                        setEditingMethod(null);
                    }}
                    onSave={handleSave}
                />
            )}
        </DashboardLayout>
    );
};

// Payment Method Modal Component
interface PaymentMethodModalProps {
    method: PaymentMethod | null;
    onClose: () => void;
    onSave: (data: Partial<PaymentMethod>) => void;
}

const PaymentMethodModal: React.FC<PaymentMethodModalProps> = ({ method, onClose, onSave }) => {
    const [formData, setFormData] = useState({
        code: method?.code || '',
        name: method?.name || '',
        category: method?.category || 'bank_transfer',
        icon_url: method?.icon_url || '',
        description: method?.description || '',
        is_active: method?.is_active ?? true,
        sort_order: method?.sort_order || 0,
    });
    const [saving, setSaving] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);

        try {
            const data = {
                ...formData,
                ...(method?.id && { id: method.id }),
            };
            await onSave(data);
        } catch (error) {
            console.error('Failed to save payment method:', error);
        } finally {
            setSaving(false);
        }
    };

    const categories = [
        { value: 'mobile_money', label: 'Mobile Money' },
        { value: 'bank_transfer', label: 'Bank Transfer' },
        { value: 'online_wallet', label: 'Online Wallet' },
        { value: 'crypto', label: 'Crypto' },
        { value: 'cash', label: 'Cash' },
    ];

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="p-8">
                    <div className="flex items-center justify-between mb-6">
                        <h2 className="text-xl font-black text-white">
                            {method ? 'Edit Payment Method' : 'Add Payment Method'}
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
                                onChange={(e) => setFormData(prev => ({ ...prev, code: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="e.g., bkash, paypal"
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
                                placeholder="e.g., bKash, PayPal"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Category *
                            </label>
                            <select
                                value={formData.category}
                                onChange={(e) => setFormData(prev => ({ ...prev, category: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white focus:border-primary/50 transition-all"
                                required
                            >
                                {categories.map(cat => (
                                    <option key={cat.value} value={cat.value} className="bg-surface">
                                        {cat.label}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Icon URL
                            </label>
                            <input
                                type="url"
                                value={formData.icon_url}
                                onChange={(e) => setFormData(prev => ({ ...prev, icon_url: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="https://example.com/icon.png"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Description
                            </label>
                            <textarea
                                value={formData.description}
                                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all resize-none"
                                rows={3}
                                placeholder="Optional description..."
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Sort Order
                            </label>
                            <input
                                type="number"
                                value={formData.sort_order}
                                onChange={(e) => setFormData(prev => ({ ...prev, sort_order: parseInt(e.target.value) || 0 }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white focus:border-primary/50 transition-all"
                                min="0"
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
                                {saving ? 'Saving...' : (method ? 'Update' : 'Create')}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AdminPaymentMethods;