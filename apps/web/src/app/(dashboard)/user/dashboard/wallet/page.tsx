"use client";

import React, { useEffect, useMemo, useState } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ArrowUpRight,
    ArrowDownLeft,
    History,
    Wallet as WalletIcon,
    RefreshCw
} from "lucide-react";
import { authService, WalletBalance, WalletTransaction } from "@/services/authService";

const UserWallet = () => {
    const [wallets, setWallets] = useState<WalletBalance[]>([]);
    const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const loadWalletData = async () => {
        setLoading(true);
        setError(null);
        try {
            const [walletData, txData] = await Promise.all([
                authService.getWalletBalances(),
                authService.getWalletTransactions({ limit: 10, offset: 0 }),
            ]);
            setWallets(walletData);
            setTransactions(txData.transactions);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load wallet data");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void loadWalletData();
    }, []);

    const totalBalance = useMemo(
        () => wallets.reduce((sum, wallet) => sum + Number(wallet.balance || 0), 0),
        [wallets]
    );

    return (
        <DashboardLayout title="Universal Wallet" role="user">
            <div className="space-y-8">
                <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                            <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center">
                                <WalletIcon className="text-white w-6 h-6" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Total Balance</p>
                                <h2 className="text-3xl font-black text-white">{totalBalance.toFixed(8)}</h2>
                            </div>
                        </div>
                        <button
                            onClick={() => void loadWalletData()}
                            className="px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-xs font-bold flex items-center gap-2"
                        >
                            <RefreshCw className="w-4 h-4" /> Refresh
                        </button>
                    </div>
                </div>

                {error && (
                    <div className="bg-primary/10 border border-primary/30 rounded-2xl p-4 text-sm text-primary">{error}</div>
                )}

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    <div className="bg-surface border border-white/10 rounded-[2rem] p-6">
                        <h3 className="text-lg font-black text-white mb-4 uppercase">Wallet Balances</h3>
                        <div className="space-y-3">
                            {loading && <p className="text-text-dim text-sm">Loading balances...</p>}
                            {!loading && wallets.length === 0 && <p className="text-text-dim text-sm">No wallets found.</p>}
                            {wallets.map((wallet) => (
                                <div key={wallet.wallet_id} className="p-4 rounded-xl bg-white/5 border border-white/5">
                                    <div className="flex items-center justify-between">
                                        <span className="text-xs font-bold text-white">Asset #{wallet.crypto_id}</span>
                                        <span className="text-xs text-text-dim">{wallet.is_active ? "Active" : "Inactive"}</span>
                                    </div>
                                    <p className="text-lg font-black text-white mt-2">{Number(wallet.balance).toFixed(8)}</p>
                                    <p className="text-[10px] text-text-dim">Locked: {Number(wallet.locked_balance).toFixed(8)}</p>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-[2rem] p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-black text-white uppercase">Recent Activity</h3>
                            <History className="w-4 h-4 text-text-dim" />
                        </div>
                        <div className="space-y-3">
                            {loading && <p className="text-text-dim text-sm">Loading transactions...</p>}
                            {!loading && transactions.length === 0 && <p className="text-text-dim text-sm">No transactions yet.</p>}
                            {transactions.map((tx) => (
                                <div key={tx.tx_id} className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/5">
                                    <div className="flex items-center gap-3">
                                        {tx.type === "withdrawal" ? (
                                            <ArrowUpRight className="w-4 h-4 text-primary" />
                                        ) : (
                                            <ArrowDownLeft className="w-4 h-4 text-accent" />
                                        )}
                                        <div>
                                            <p className="text-xs font-bold text-white uppercase">{tx.type}</p>
                                            <p className="text-[10px] text-text-dim">{new Date(tx.created_at).toLocaleString()}</p>
                                        </div>
                                    </div>
                                    <p className="text-xs font-black text-white">{Number(tx.amount).toFixed(8)}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default UserWallet;
