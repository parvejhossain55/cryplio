"use client";

import React from "react";
import { TrendingUp, ArrowUpRight, ArrowDownLeft } from "lucide-react";

interface Coin {
    name: string;
    symbol: string;
    amount: string;
    price: string;
    icon: string;
}

const WalletCard = () => {
    const coins: Coin[] = [
        { name: "Bitcoin", symbol: "BTC", amount: "0.45", price: "$28,450", icon: "₿" },
        { name: "Ethereum", symbol: "ETH", amount: "12.5", price: "$12,400", icon: "Ξ" },
        { name: "USDT", symbol: "USDT", amount: "2,000", price: "$2,000", icon: "₮" },
        { name: "Solana", symbol: "SOL", amount: "450.0", price: "$125", icon: "S" },
    ];

    return (
        <div className="relative overflow-hidden bg-surface rounded-[2.5rem] border border-white/10 p-8 group">
            <div className="absolute top-0 right-0 w-[50%] h-[150%] bg-primary/10 -rotate-45 translate-x-[20%] -translate-y-[20%] blur-3xl pointer-events-none" />

            <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
                <div>
                    <p className="text-xs font-black text-text-dim uppercase tracking-[0.2em] mb-2">Total Balance</p>
                    <div className="flex items-baseline space-x-3">
                        <h2 className="text-5xl font-black text-white">$42,850.24</h2>
                        <span className="text-accent font-bold bg-accent/10 px-2 py-0.5 rounded-lg text-sm flex items-center">
                            <TrendingUp className="w-3 h-3 mr-1" /> +12.5%
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
                {coins.map((coin) => (
                    <div key={coin.symbol} className="p-4 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 transition-all cursor-pointer group/card">
                        <div className="w-10 h-10 rounded-2xl bg-surface flex items-center justify-center text-lg font-black mb-3 group-hover/card:scale-110 transition-transform">
                            {coin.icon}
                        </div>
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">{coin.symbol}</p>
                        <p className="text-sm font-bold text-white mt-1">{coin.amount}</p>
                        <p className="text-[10px] font-medium text-text-dim mt-0.5">{coin.price}</p>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default WalletCard;
