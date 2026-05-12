"use client";

import React, { useEffect, useState } from "react";
import { TrendingUp, ArrowUpRight, ArrowDownLeft, Loader2 } from "lucide-react";
import { walletService } from "@/services/walletService";
import { WalletBalance } from "@/types/api";

const WalletCard = () => {
    const [wallets, setWallets] = useState<WalletBalance[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchBalances = async () => {
            try {
                const balances = await walletService.getBalances();
                setWallets(balances);
            } catch (error) {
                console.error("Failed to fetch balances:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchBalances();
    }, []);

    const totalBalance = wallets.reduce((acc, wallet) => acc + Number(wallet.balance), 0);

    return (
        <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-8 group">
            <div className="absolute top-0 right-0 w-[50%] h-[150%] bg-primary/10 -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

            <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                <div>
                    <p className="text-xs font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Balance</p>
                    <div className="flex items-baseline space-x-3">
                        {loading ? (
                            <Loader2 className="w-8 h-8 animate-spin text-primary" />
                        ) : (
                            <h2 className="text-5xl font-black text-white">{totalBalance.toFixed(8)}</h2>
                        )}
                        <span className="text-accent font-bold bg-accent/10 px-2 py-0.5 rounded-lg text-sm flex items-center">
                            <TrendingUp className="w-3 h-3 mr-1" /> +0.0%
                        </span>
                    </div>
                </div>

                <div className="flex items-center space-x-3">
                    <button className="flex-1 md:flex-none px-6 py-4 bg-white text-background rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:scale-105 active:scale-95">
                        <ArrowUpRight className="w-4 h-4 mr-2" /> Send
                    </button>
                    <button className="flex-1 md:flex-none px-6 py-4 bg-surface-light border border-white/5 text-white rounded-2xl font-black text-sm flex items-center justify-center transition-all hover:bg-white/5 active:scale-95">
                        <ArrowDownLeft className="w-4 h-4 mr-2" /> Receive
                    </button>
                </div>
            </div>

            <div className="relative z-10 mt-10 grid grid-cols-2 sm:grid-cols-4 gap-4">
                {loading ? (
                    <div className="col-span-full flex justify-center py-4">
                        <Loader2 className="w-6 h-6 animate-spin text-primary" />
                    </div>
                ) : wallets.length === 0 ? (
                    <div className="col-span-full text-center py-4 text-text-dim text-xs">
                        No active wallets found.
                    </div>
                ) : (
                    wallets.map((wallet) => (
                        <div key={wallet.wallet_id} className="p-4 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 transition-all cursor-pointer group/card">
                            <div className="w-10 h-10 rounded-2xl bg-surface flex items-center justify-center text-lg font-black mb-3 group-hover/card:scale-110 transition-transform">
                                {wallet.crypto_symbol === 'BTC' ? '₿' : wallet.crypto_symbol === 'ETH' ? 'Ξ' : wallet.crypto_symbol === 'USDT' ? '₮' : 'S'}
                            </div>
                            <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{wallet.crypto_symbol || 'ETH'}</p>
                            <p className="text-sm font-bold text-white mt-1">{Number(wallet.balance).toFixed(4)}</p>
                            <p className="text-[10px] font-medium text-text-dim mt-0.5">Primary Wallet</p>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default WalletCard;
