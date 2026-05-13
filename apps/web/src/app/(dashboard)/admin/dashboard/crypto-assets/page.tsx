"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff, X } from "lucide-react";
import { toast } from "sonner";

interface CryptoAsset {
    id: number;
    symbol: string;
    name: string;
    blockchain: string;
    contract_address?: string;
    decimals: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

const AdminCryptoAssets = () => {
    const [cryptoAssets, setCryptoAssets] = useState<CryptoAsset[]>([]);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [showInactive, setShowInactive] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingAsset, setEditingAsset] = useState<CryptoAsset | null>(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0,
        pages: 0,
    });

    useEffect(() => {
        const timer = setTimeout(() => {
            fetchCryptoAssets(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchQuery, showInactive]);

    useEffect(() => {
        fetchCryptoAssets(pagination.page);
    }, [pagination.page]);

    const fetchCryptoAssets = async (page = 1) => {
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

            const response = await fetch(`/api/v1/admin/crypto-assets?${params}`);
            if (response.ok) {
                const data = await response.json();
                setCryptoAssets(data.crypto_assets || []);
                setPagination(data.pagination || { page: 1, limit: 50, total: 0, pages: 0 });
            }
        } catch (error) {
            console.error("Failed to fetch crypto assets:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (assetData: Partial<CryptoAsset>) => {
        try {
            const url = editingAsset
                ? `/api/v1/admin/crypto-assets/${editingAsset.id}`
                : "/api/v1/admin/crypto-assets";
            const method = editingAsset ? "PUT" : "POST";

            const response = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(assetData),
            });

            if (response.ok) {
                await fetchCryptoAssets(pagination.page);
                setIsModalOpen(false);
                setEditingAsset(null);
                toast.success(editingAsset ? "Crypto asset updated" : "Crypto asset created");
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to save crypto asset");
            }
        } catch (error) {
            console.error("Failed to save crypto asset:", error);
            toast.error("Failed to save crypto asset");
        }
    };

    const handleDelete = async (id: number) => {
        const confirmDelete = window.confirm("Are you sure you want to delete this crypto asset? This may affect existing trades.");
        if (!confirmDelete) return;

        try {
            const response = await fetch(`/api/v1/admin/crypto-assets/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (response.ok) {
                toast.success("Crypto asset deleted");
                await fetchCryptoAssets(pagination.page);
            } else {
                const errorData = await response.json();
                toast.error(errorData.error || "Failed to delete crypto asset");
            }
        } catch (error) {
            console.error("Failed to delete crypto asset:", error);
            toast.error("Failed to delete crypto asset");
        }
    };



    const filteredAssets = cryptoAssets.filter(asset => {
        const matchesSearch = asset.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            asset.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
            asset.blockchain.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesActive = showInactive || asset.is_active;
        return matchesSearch && matchesActive;
    });

    if (loading) {
        return (
            <DashboardLayout title="Crypto Assets" role="admin">
                <div className="flex items-center justify-center h-64">
                    <div className="text-white">Loading...</div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Crypto Assets" role="admin">
            <div className="space-y-8">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-black text-white">Crypto Assets Management</h1>
                        <p className="text-text-dim mt-1">Manage supported cryptocurrencies for trading</p>
                    </div>
                    <button
                        onClick={() => setIsModalOpen(true)}
                        className="flex items-center px-6 py-3 bg-primary text-white rounded-xl font-bold hover:scale-105 active:scale-95 transition-all"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Add Crypto Asset
                    </button>
                </div>

                {/* Filters */}
                <div className="flex items-center space-x-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search crypto assets..."
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
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Symbol</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Name</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Blockchain</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Decimals</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Status</th>
                                    <th className="px-6 py-4 text-right text-xs font-black text-text-dim uppercase tracking-widest">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {cryptoAssets.map((asset) => (
                                    <tr key={asset.id} className="border-t border-white/5 hover:bg-white/2 transition-all">
                                        <td className="px-6 py-4">
                                            <code className="text-sm font-bold text-primary bg-primary/10 px-2 py-1 rounded">
                                                {asset.symbol}
                                            </code>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className="font-bold text-white">{asset.name}</div>
                                            {asset.contract_address && (
                                                <div className="text-xs text-text-dim mt-1 font-mono">
                                                    {asset.contract_address.substring(0, 10)}...
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-bold bg-blue-500/20 text-blue-500">
                                                {asset.blockchain}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-white font-bold">{asset.decimals}</td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${asset.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
                                                }`}>
                                                {asset.is_active ? 'Active' : 'Inactive'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-right space-x-2">
                                            <button
                                                onClick={() => {
                                                    setEditingAsset(asset);
                                                    setIsModalOpen(true);
                                                }}
                                                className="inline-flex items-center px-3 py-2 bg-blue-500/20 text-blue-500 rounded-lg hover:bg-blue-500/30 transition-all"
                                            >
                                                <Edit className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleDelete(asset.id)}
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

                {cryptoAssets.length === 0 && (
                    <div className="text-center py-12 text-text-dim">
                        No crypto assets found matching your criteria.
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

            {/* Crypto Asset Modal */}
            {isModalOpen && (
                <CryptoAssetModal
                    asset={editingAsset}
                    onClose={() => {
                        setIsModalOpen(false);
                        setEditingAsset(null);
                    }}
                    onSave={handleSave}
                />
            )}
        </DashboardLayout>
    );
};

// Crypto Asset Modal Component
interface CryptoAssetModalProps {
    asset: CryptoAsset | null;
    onClose: () => void;
    onSave: (data: Partial<CryptoAsset>) => void;
}

const CryptoAssetModal: React.FC<CryptoAssetModalProps> = ({ asset, onClose, onSave }) => {
    const [formData, setFormData] = useState({
        symbol: asset?.symbol || '',
        name: asset?.name || '',
        blockchain: asset?.blockchain || '',
        contract_address: asset?.contract_address || '',
        decimals: asset?.decimals || 18,
        is_active: asset?.is_active ?? true,
    });
    const [saving, setSaving] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);

        try {
            await onSave(formData);
        } catch (error) {
            console.error('Failed to save crypto asset:', error);
        } finally {
            setSaving(false);
        }
    };

    const blockchains = [
        'ETH', 'BSC', 'TRC20', 'ERC20', 'BTC', 'SOL', 'AVAX', 'MATIC', 'ARB', 'OP'
    ];

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="p-8">
                    <div className="flex items-center justify-between mb-6">
                        <h2 className="text-xl font-black text-white">
                            {asset ? 'Edit Crypto Asset' : 'Add Crypto Asset'}
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
                                Symbol *
                            </label>
                            <input
                                type="text"
                                value={formData.symbol}
                                onChange={(e) => setFormData(prev => ({ ...prev, symbol: e.target.value.toUpperCase() }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all"
                                placeholder="e.g., BTC, ETH"
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
                                placeholder="e.g., Bitcoin, Ethereum"
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Blockchain *
                            </label>
                            <select
                                value={formData.blockchain}
                                onChange={(e) => setFormData(prev => ({ ...prev, blockchain: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white focus:border-primary/50 transition-all"
                                required
                            >
                                <option value="" className="bg-surface">Select blockchain</option>
                                {blockchains.map(bc => (
                                    <option key={bc} value={bc} className="bg-surface">
                                        {bc}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                Contract Address
                            </label>
                            <input
                                type="text"
                                value={formData.contract_address}
                                onChange={(e) => setFormData(prev => ({ ...prev, contract_address: e.target.value }))}
                                className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white placeholder-text-dim focus:border-primary/50 transition-all font-mono"
                                placeholder="0x..."
                            />
                        </div>

                        <div className="grid grid-cols-1 gap-4">
                            <div>
                                <label className="block text-xs font-black text-text-dim uppercase tracking-widest mb-2">
                                    Decimals *
                                </label>
                                <input
                                    type="number"
                                    value={formData.decimals}
                                    onChange={(e) => setFormData(prev => ({ ...prev, decimals: parseInt(e.target.value) || 18 }))}
                                    className="w-full bg-white/5 border border-white/5 py-3 px-4 rounded-xl text-sm font-bold text-white focus:border-primary/50 transition-all"
                                    min="0"
                                    max="18"
                                    required
                                />
                            </div>
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
                                {saving ? 'Saving...' : (asset ? 'Update' : 'Create')}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AdminCryptoAssets;