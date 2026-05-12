"use client";

import React, { useEffect, useState } from "react";
import { ChevronRight, ArrowDownLeft, ArrowUpRight, ArrowLeftRight, ExternalLink, Loader2 } from "lucide-react";
import { walletService } from "@/services/walletService";
import { WalletTransaction } from "@/types/api";

const RecentActivity = () => {
    const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchTransactions = async () => {
            try {
                const data = await walletService.getTransactions({ limit: 5 });
                setTransactions(data.transactions);
            } catch (error) {
                console.error("Failed to fetch transactions:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchTransactions();
    }, []);

    return (
        <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
            <div className="flex items-center justify-between mb-8">
                <h3 className="text-xl font-black text-white">Recent Activity</h3>
                <button className="text-xs font-bold text-primary hover:text-white transition-colors flex items-center">
                    View All <ChevronRight className="w-4 h-4 ml-1" />
                </button>
            </div>

            <div className="space-y-4">
                {loading ? (
                    <div className="flex justify-center py-8">
                        <Loader2 className="w-8 h-8 animate-spin text-primary" />
                    </div>
                ) : transactions.length === 0 ? (
                    <div className="text-center py-8 text-text-dim text-sm">
                        No recent activity found.
                    </div>
                ) : (
                    transactions.map((tx) => (
                        <div key={tx.tx_id} className="flex items-center justify-between p-4 rounded-2xl hover:bg-white/5 transition-all group">
                            <div className="flex items-center space-x-4">
                                <div className={`w-12 h-12 rounded-2xl flex items-center justify-center border border-white/5 ${tx.type === 'deposit' ? 'bg-accent/10 text-accent' :
                                    tx.type === 'withdrawal' ? 'bg-primary/10 text-primary' :
                                        'bg-secondary/10 text-secondary'
                                    }`}>
                                    {tx.type === 'deposit' ? <ArrowDownLeft className="w-6 h-6" /> :
                                        tx.type === 'withdrawal' ? <ArrowUpRight className="w-6 h-6" /> :
                                            <ArrowLeftRight className="w-5 h-5" />}
                                </div>
                                <div>
                                    <h4 className="font-bold text-white tracking-tight flex items-center uppercase">
                                        {tx.type} {tx.crypto_symbol || 'ETH'}
                                        <span className={`ml-2 text-[8px] px-1.5 py-0.5 rounded font-black uppercase tracking-widest ${tx.status === 'completed' || tx.status === 'confirmed' ? 'bg-accent/10 text-accent' :
                                            tx.status === 'pending' ? 'bg-primary/10 text-primary' :
                                                'bg-red-500/10 text-red-500'
                                            }`}>
                                            {tx.status}
                                        </span>
                                    </h4>
                                    <p className="text-[10px] font-medium text-text-dim mt-1">{new Date(tx.created_at).toLocaleString()} • ID: {tx.tx_id.slice(0, 8)}</p>
                                </div>
                            </div>
                            <div className="text-right">
                                <p className="font-black text-white">{Number(tx.amount).toFixed(4)} {tx.crypto_symbol || 'ETH'}</p>
                                <ExternalLink className="w-3 h-3 ml-auto mt-1 text-text-dim opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer hover:text-white" />
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default RecentActivity;
