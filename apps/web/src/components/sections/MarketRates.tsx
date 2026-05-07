"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import {
    TrendingUp,
    TrendingDown,
    RefreshCw,
    DollarSign,
    Euro,
    PoundSterling,
    JapaneseYen
} from "lucide-react";
import { toast } from "sonner";

interface MarketRate {
    crypto_symbol: string;
    fiat_symbol: string;
    price: number;
    source: string;
    as_of: string;
    change_24h?: number;
}

const MarketRates = () => {
    const [rates, setRates] = useState<MarketRate[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [selectedFiat, setSelectedFiat] = useState("USD");

    const fiatCurrencies = [
        { symbol: "USD", icon: DollarSign, name: "US Dollar" },
        { symbol: "EUR", icon: Euro, name: "Euro" },
        { symbol: "GBP", icon: PoundSterling, name: "British Pound" },
        { symbol: "JPY", icon: JapaneseYen, name: "Japanese Yen" }
    ];

    const cryptoCurrencies = ["BTC", "ETH", "USDT", "USDC", "BNB", "SOL", "ADA", "DOT"];

    useEffect(() => {
        fetchRates();
        const interval = setInterval(fetchRates, 30000); // Refresh every 30 seconds
        return () => clearInterval(interval);
    }, [selectedFiat]);

    const fetchRates = async () => {
        try {
            const response = await fetch(`/api/market/rates?fiat=${selectedFiat}`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch market rates");
            }
            
            setRates(data.rates || []);
        } catch (error) {
            console.error("Error fetching market rates:", error);
            toast.error("Failed to load market rates");
        } finally {
            setIsLoading(false);
        }
    };

    const formatPrice = (price: number, fiat: string) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: fiat,
            minimumFractionDigits: 2,
            maximumFractionDigits: price < 1 ? 8 : 2
        }).format(price);
    };

    const getChangeColor = (change?: number) => {
        if (!change) return "text-text-dim";
        return change > 0 ? "text-green-400" : "text-red-400";
    };

    const getChangeIcon = (change?: number) => {
        if (!change) return null;
        return change > 0 ? TrendingUp : TrendingDown;
    };

    const groupedRates = rates.reduce((acc, rate) => {
        if (!acc[rate.crypto_symbol]) {
            acc[rate.crypto_symbol] = [];
        }
        acc[rate.crypto_symbol].push(rate);
        return acc;
    }, {} as Record<string, MarketRate[]>);

    if (isLoading) {
        return (
            <div className="bg-surface border border-white/10 rounded-2xl p-8">
                <div className="flex items-center justify-center h-64">
                    <RefreshCw className="w-8 h-8 animate-spin text-primary" />
                </div>
            </div>
        );
    }

    return (
        <div className="bg-surface border border-white/10 rounded-2xl p-8">
            <div className="flex items-center justify-between mb-8">
                <h3 className="text-xl font-black text-white uppercase tracking-tight">
                    Market Rates
                </h3>
                
                <div className="flex items-center space-x-4">
                    <select
                        value={selectedFiat}
                        onChange={(e) => setSelectedFiat(e.target.value)}
                        className="px-4 py-2 bg-white/10 border border-white/20 rounded-xl text-white text-sm focus:outline-none focus:border-primary/50"
                    >
                        {fiatCurrencies.map(fiat => (
                            <option key={fiat.symbol} value={fiat.symbol}>
                                {fiat.name}
                            </option>
                        ))}
                    </select>
                    
                    <button
                        onClick={fetchRates}
                        className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors"
                        title="Refresh Rates"
                    >
                        <RefreshCw className="w-4 h-4" />
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {Object.entries(groupedRates).map(([crypto, cryptoRates]) => {
                    const primaryRate = cryptoRates.find(r => r.fiat_symbol === selectedFiat) || cryptoRates[0];
                    const ChangeIcon = getChangeIcon(primaryRate?.change_24h);
                    
                    return (
                        <motion.div
                            key={crypto}
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            className="bg-white/5 border border-white/10 rounded-xl p-6 hover:bg-white/10 transition-colors"
                        >
                            <div className="flex items-center justify-between mb-4">
                                <div className="flex items-center space-x-2">
                                    <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                                        <span className="text-xs font-bold text-primary">
                                            {crypto.slice(0, 2)}
                                        </span>
                                    </div>
                                    <span className="text-sm font-bold text-white">
                                        {crypto}
                                    </span>
                                </div>
                                
                                {ChangeIcon && (
                                    <div className={`flex items-center space-x-1 ${getChangeColor(primaryRate?.change_24h)}`}>
                                        <ChangeIcon className="w-4 h-4" />
                                        <span className="text-xs font-medium">
                                            {Math.abs(primaryRate.change_24h || 0).toFixed(2)}%
                                        </span>
                                    </div>
                                )}
                            </div>
                            
                            <div className="space-y-2">
                                <div className="text-2xl font-bold text-white">
                                    {formatPrice(primaryRate.price, primaryRate.fiat_symbol)}
                                </div>
                                
                                <div className="text-xs text-text-dim">
                                    per {primaryRate.fiat_symbol}
                                </div>
                                
                                {cryptoRates.length > 1 && (
                                    <div className="pt-2 border-t border-white/5">
                                        <div className="text-xs text-text-dim mb-2">Other pairs:</div>
                                        <div className="space-y-1">
                                            {cryptoRates
                                                .filter(r => r.fiat_symbol !== selectedFiat)
                                                .slice(0, 2)
                                                .map(rate => (
                                                    <div key={rate.fiat_symbol} className="flex justify-between text-xs">
                                                        <span className="text-text-dim">{rate.fiat_symbol}</span>
                                                        <span className="text-white">
                                                            {formatPrice(rate.price, rate.fiat_symbol)}
                                                        </span>
                                                    </div>
                                                ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        </motion.div>
                    );
                })}
            </div>

            {rates.length === 0 && (
                <div className="text-center py-12">
                    <TrendingDown className="w-12 h-12 text-text-dim mx-auto mb-4" />
                    <p className="text-text-dim">No market rates available</p>
                </div>
            )}
        </div>
    );
};

export default MarketRates;
