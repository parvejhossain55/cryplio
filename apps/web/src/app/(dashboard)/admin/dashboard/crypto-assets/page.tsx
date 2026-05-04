"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Plus, Edit, Trash2, Search, Eye, EyeOff } from "lucide-react";

interface CryptoAsset {
    id: number;
    symbol: string;
    name: string;
    blockchain: string;
    contract_address?: string;
    decimals: number;
    min_confirmation: number;
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

    useEffect(() => {
        fetchCryptoAssets();
    }, []);

    const fetchCryptoAssets = async () => {
        try {
            const response = await fetch(`/api/v1/admin/crypto-assets?active_only=false`);
            if (response.ok) {
                const data = await response.json();
                setCryptoAssets(data.crypto_assets || []);
            }
        } catch (error) {
            console.error("Failed to fetch crypto assets:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async (assetData: Partial<CryptoAsset>) => {
        try {
            const response = await fetch("/api/v1/admin/crypto-assets", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(assetData),
            });

            if (response.ok) {
                await fetchCryptoAssets();
                setIsModalOpen(false);
                setEditingAsset(null);
            }
        } catch (error) {
            console.error("Failed to create crypto asset:", error);
        }
    };

    const handleUpdate = async (id: number, assetData: Partial<CryptoAsset>) => {
        try {
            const response = await fetch(`/api/v1/admin/crypto-assets/${id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify(assetData),
            });

            if (response.ok) {
                await fetchCryptoAssets();
                setIsModalOpen(false);
                setEditingAsset(null);
            }
        } catch (error) {
            console.error("Failed to update crypto asset:", error);
        }
    };

    const handleDelete = async (id: number) => {
        if (!confirm("Are you sure you want to delete this crypto asset? This may affect existing trades.")) return;

        try {
            const response = await fetch(`/api/v1/admin/crypto-assets/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (response.ok) {
                await fetchCryptoAssets();
            }
        } catch (error) {
            console.error("Failed to delete crypto asset:", error);
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
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Symbol</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Name</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Blockchain</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Decimals</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Min Confirm</th>
                                    <th className="px-6 py-4 text-left text-xs font-black text-text-dim uppercase tracking-widest">Status</th>
                                    <th className="px-6 py-4 text-right text-xs font-black text-text-dim uppercase tracking-widest">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredAssets.map((asset) => (
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
                                        <td className="px-6 py-4 text-white font-bold">{asset.min_confirmation}</td>
                                        <td className="px-6 py-4">
                                            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-bold ${
                                                asset.is_active ? 'bg-accent/20 text-accent' : 'bg-red-500/20 text-red-500'
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

                {filteredAssets.length === 0 && (
                    <div className="text-center py-12 text-text-dim">
                        No crypto assets found matching your criteria.
                    </div>
                )}
            </div>

            {/* Modal will be added next */}
        </DashboardLayout>
    );
};

export default AdminCryptoAssets;