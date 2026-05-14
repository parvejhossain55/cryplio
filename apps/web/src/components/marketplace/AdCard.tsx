"use client";

import React from "react";
import { motion } from "framer-motion";
import {
    Clock,
    TrendingUp,
    TrendingDown,
    User,
    Star,
    ArrowRight
} from "lucide-react";

export interface TradeAd {
    ad_id: string;
    user_id: string;
    username: string;
    user_avatar?: string;
    is_online?: boolean;
    type: "buy" | "sell";
    crypto_symbol: string;
    fiat_symbol: string;
    price_type: "fixed" | "floating";
    price: number;
    min_amount: number;
    max_amount: number;
    payment_methods: string[];
    payment_method_ids: number[];
    payment_window_minutes: number;
    trade_terms?: string;
    status: "active" | "paused" | "closed";
    created_at: string;
    user_trades?: number;
    user_rating?: number;
}

interface AdCardProps {
    ad: TradeAd;
    onTrade: (adId: string) => void;
}

const AdCard: React.FC<AdCardProps> = ({ ad, onTrade }) => {
    const formatPrice = (ad: TradeAd) => {
        return `${ad.price.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ${ad.fiat_symbol}`;
    };

    const formatAmount = (amount: number) => {
        return amount.toLocaleString("en-US", { minimumFractionDigits: 0, maximumFractionDigits: 2 });
    };

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-surface border border-white/10 rounded-2xl p-6 hover:border-primary/50 transition-all group"
        >
            {/* Header */}
            <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                    <div className={`px-3 py-1 rounded-full text-xs font-black uppercase ${ad.type === "buy"
                        ? "bg-green-500/20 text-green-500 border border-green-500/20"
                        : "bg-red-500/20 text-red-500 border border-red-500/20"
                        }`}>
                        {ad.type === "buy" ? "BUY" : "SELL"}
                    </div>
                    <span className="text-xs font-bold text-white/60 uppercase tracking-wider">
                        {ad.crypto_symbol}
                    </span>
                </div>
                <div className="flex items-center gap-1">
                    <Clock className="w-3 h-3 text-text-dim" />
                    <span className="text-xs text-text-dim">{ad.payment_window_minutes}m</span>
                </div>
            </div>

            {/* Price */}
            <div className="mb-4">
                <div className="flex items-center gap-2">
                    {ad.type === "buy" ? (
                        <TrendingDown className="w-4 h-4 text-green-500" />
                    ) : (
                        <TrendingUp className="w-4 h-4 text-red-500" />
                    )}
                    <span className="text-2xl font-black text-white">
                        {formatPrice(ad)}
                    </span>
                </div>
                <p className="text-text-dim text-sm">
                    {formatAmount(ad.min_amount)} - {formatAmount(ad.max_amount)} {ad.fiat_symbol}
                </p>
            </div>

            {/* User Info */}
            <div className="flex items-center justify-between mb-4 bg-white/5 p-3 rounded-xl">
                <div className="flex items-center gap-3">
                    <div className="relative">
                        {ad.user_avatar ? (
                            <img
                                src={ad.user_avatar}
                                alt={ad.username}
                                className="w-10 h-10 rounded-full object-cover border border-white/10"
                            />
                        ) : (
                            <div className="w-10 h-10 bg-primary/20 rounded-full flex items-center justify-center border border-white/10">
                                <User className="w-5 h-5 text-primary" />
                            </div>
                        )}
                        {ad.is_online && (
                            <span className="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 rounded-full border-2 border-surface"></span>
                        )}
                    </div>
                    <div>
                        <p className="text-white font-medium">{ad.username}</p>
                        <p className="text-xs text-text-dim">
                            {ad.user_trades || 0} trades
                        </p>
                    </div>
                </div>
                {ad.user_rating ? (
                    <div className="flex items-center gap-1">
                        <Star className="w-4 h-4 text-yellow-500 fill-current" />
                        <span className="text-sm text-white">{ad.user_rating.toFixed(1)}</span>
                    </div>
                ) : null}
            </div>

            {/* Payment Methods */}
            <div className="flex flex-wrap gap-2 mb-4">
                {(ad.payment_methods || []).slice(0, 2).map((method) => (
                    <span
                        key={method}
                        className="px-2 py-1 bg-white/5 border border-white/10 rounded-lg text-xs text-text-dim"
                    >
                        {method}
                    </span>
                ))}
                {(ad.payment_methods || []).length > 2 && (
                    <span className="px-2 py-1 bg-white/5 border border-white/10 rounded-lg text-xs text-text-dim">
                        +{(ad.payment_methods || []).length - 2}
                    </span>
                )}
            </div>

            {/* Trade Terms */}
            {ad.trade_terms && (
                <p className="text-text-dim text-sm mb-4 line-clamp-2">
                    {ad.trade_terms}
                </p>
            )}

            {/* Action Button */}
            <button
                onClick={() => onTrade(ad.ad_id)}
                className="w-full py-3 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-all flex items-center justify-center gap-2 group-hover:scale-[1.02]"
            >
                Trade Now
                <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </button>
        </motion.div>
    );
};

export default AdCard;
