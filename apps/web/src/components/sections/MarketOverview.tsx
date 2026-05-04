"use client";

import React, { useState, useMemo, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Search, Filter, ArrowUpRight, ChevronDown, CheckCircle2, ShieldCheck, Clock, Star, TrendingUp } from "lucide-react";
import Link from "next/link";
import Pagination from "@/components/ui/Pagination";

const p2pAds = [
    // USD - BUY
    { user: "FlashExchange", trades: "850+", rating: "99%", price: "1.02", currency: "USD", coin: "USDT", limits: "10 - 2,000", methods: ["PayPal", "Zelle"], type: "buy", time: "5 min", available: "12,000.00" },
    { user: "QuickCrypto_US", trades: "1,200+", rating: "98.5%", price: "1.01", currency: "USD", coin: "USDT", limits: "100 - 5,000", methods: ["Bank Transfer", "Cash App"], type: "buy", time: "10 min", available: "25,000.00" },
    { user: "EliteTrader_NY", trades: "5,400+", rating: "99.9%", price: "1.03", currency: "USD", coin: "BTC", limits: "500 - 10,000", methods: ["Bank Transfer", "Zelle"], type: "buy", time: "8 min", available: "0.50" },
    { user: "EasyBuy_Crypto", trades: "320+", rating: "97%", price: "1.05", currency: "USD", coin: "ETH", limits: "20 - 1,000", methods: ["Venmo", "PayPal"], type: "buy", time: "15 min", available: "5.00" },
    { user: "CryptoWhale_X", trades: "15k+", rating: "99.1%", price: "1.00", currency: "USD", coin: "USDT", limits: "1,000 - 50,000", methods: ["WebMoney", "Advcash"], type: "buy", time: "12 min", available: "100,000.00" },
    { user: "ZenTrader", trades: "2,100+", rating: "99.5%", price: "1.02", currency: "USD", coin: "USDT", limits: "50 - 3,000", methods: ["Zelle", "Bank Transfer"], type: "buy", time: "7 min", available: "15,000.00" },
    { user: "TitanLiquidity", trades: "8,900+", rating: "100%", price: "1.02", currency: "USD", coin: "USDT", limits: "200 - 15,000", methods: ["Bank Transfer"], type: "buy", time: "4 min", available: "80,000.00" },
    { user: "AlphaAssets", trades: "1,500+", rating: "99%", price: "1.01", currency: "USD", coin: "USDT", limits: "100 - 10,000", methods: ["Bank Transfer"], type: "buy", time: "3 min", available: "45,000.00" },
    { user: "BetaExchange", trades: "750+", rating: "98%", price: "1.03", currency: "USD", coin: "USDT", limits: "10 - 500", methods: ["PayPal"], type: "buy", time: "15 min", available: "2,500.00" },
    { user: "GammaTrade", trades: "4,200+", rating: "99.8%", price: "1.02", currency: "USD", coin: "USDT", limits: "500 - 20,000", methods: ["Zelle"], type: "buy", time: "2 min", available: "60,000.00" },

    // USD - SELL
    { user: "FastCash_US", trades: "1,100+", rating: "98%", price: "1.00", currency: "USD", coin: "USDT", limits: "50 - 2,000", methods: ["Zelle", "PayPal"], type: "sell", time: "10 min", available: "8,000.00" },
    { user: "SecureOut_X", trades: "2,300+", rating: "99.2%", price: "0.99", currency: "USD", coin: "USDT", limits: "100 - 5,000", methods: ["Bank Transfer"], type: "sell", time: "15 min", available: "20,000.00" },

    // NGN - BUY
    { user: "CryptoKing_99", trades: "1,200+", rating: "98%", price: "85,420,000", currency: "NGN", coin: "BTC", limits: "50,000 - 500,000", methods: ["Bank Transfer", "Kuda"], type: "buy", time: "15 min", available: "4,500.00" },
    { user: "NaijaTrader_Pro", trades: "3,500+", rating: "99.5%", price: "1,550", currency: "NGN", coin: "USDT", limits: "100,000 - 5M", methods: ["Bank Transfer", "Opay"], type: "buy", time: "10 min", available: "15,000.00" },
    { user: "KudaExpert", trades: "900+", rating: "97.8%", price: "1,560", currency: "NGN", coin: "USDT", limits: "10,000 - 200,000", methods: ["Kuda", "PalmPay"], type: "buy", time: "5 min", available: "2,000.00" },
    { user: "LagosWhale", trades: "12k+", rating: "99.9%", price: "1,540", currency: "NGN", coin: "USDT", limits: "500,000 - 20M", methods: ["Bank Transfer"], type: "buy", time: "20 min", available: "50,000.00" },
    { user: "SwiftNaira", trades: "1,500+", rating: "98.2%", price: "1,555", currency: "NGN", coin: "USDT", limits: "20,000 - 1M", methods: ["Opay", "Kuda"], type: "buy", time: "8 min", available: "6,500.00" },
    { user: "OpayMaster", trades: "2,800+", rating: "99%", price: "1,552", currency: "NGN", coin: "USDT", limits: "50,000 - 2M", methods: ["Opay", "Bank Transfer"], type: "buy", time: "12 min", available: "10,000.00" },
    { user: "EasySwap_NG", trades: "450+", rating: "96.5%", price: "1,565", currency: "NGN", coin: "USDT", limits: "5,000 - 50,000", methods: ["PalmPay"], type: "buy", time: "15 min", available: "1,500.00" },
    { user: "NairaLiquidity", trades: "6,700+", rating: "99.7%", price: "1,548", currency: "NGN", coin: "USDT", limits: "200,000 - 10M", methods: ["Bank Transfer"], type: "buy", time: "10 min", available: "30,000.00" },

    // NGN - SELL
    { user: "SecureTrader", trades: "3,500+", rating: "100%", price: "1,520", currency: "NGN", coin: "USDT", limits: "100,000 - 5M", methods: ["Bank Transfer"], type: "sell", time: "20 min", available: "45,000" },
    { user: "NairaOut_Fast", trades: "1,200+", rating: "98.5%", price: "1,515", currency: "NGN", coin: "USDT", limits: "50,000 - 2M", methods: ["Bank Transfer", "Kuda"], type: "sell", time: "15 min", available: "12,000" },

    // PKR - BUY
    { user: "GlobalNode", trades: "12k+", rating: "99.8%", price: "285", currency: "PKR", coin: "USDT", limits: "5,000 - 100,000", methods: ["Nayapay", "Easypaisa"], type: "buy", time: "10 min", available: "890.00" },
    { user: "PakCrypto_Hub", trades: "2,100+", rating: "99%", price: "286", currency: "PKR", coin: "USDT", limits: "10,000 - 500,000", methods: ["Bank Transfer", "Sadapay"], type: "buy", time: "12 min", available: "2,500.00" },
    { user: "EasyTrade_PK", trades: "4,500+", rating: "98.2%", price: "287", currency: "PKR", coin: "USDT", limits: "2,000 - 50,000", methods: ["Easypaisa", "JazzCash"], type: "buy", time: "5 min", available: "1,200.00" },
    { user: "NayaTrader", trades: "1,800+", rating: "99.5%", price: "285.5", currency: "PKR", coin: "USDT", limits: "5,000 - 200,000", methods: ["Nayapay", "Bank Transfer"], type: "buy", time: "8 min", available: "3,000.00" },
    { user: "DesiExchange", trades: "950+", rating: "97.5%", price: "288", currency: "PKR", coin: "USDT", limits: "1,000 - 20,000", methods: ["JazzCash", "Easypaisa"], type: "buy", time: "10 min", available: "500.00" },
    { user: "IndusLiquidity", trades: "7,200+", rating: "99.9%", price: "284.5", currency: "PKR", coin: "USDT", limits: "50,000 - 1M", methods: ["Bank Transfer"], type: "buy", time: "15 min", available: "10,000.00" },

    // EUR - BUY
    { user: "EuroCrypto_X", trades: "3,400+", rating: "99.2%", price: "0.94", currency: "EUR", coin: "USDT", limits: "100 - 5,000", methods: ["SEPA", "Revolut"], type: "buy", time: "12 min", available: "15,000.00" },
    { user: "BerlinTrader", trades: "1,200+", rating: "98.5%", price: "0.95", currency: "EUR", coin: "USDT", limits: "50 - 2,000", methods: ["Revolut", "Wise"], type: "buy", time: "5 min", available: "5,000.00" },
    { user: "ParisLiquidity", trades: "8,500+", rating: "100%", price: "0.93", currency: "EUR", coin: "USDT", limits: "500 - 20,000", methods: ["SEPA Instant"], type: "buy", time: "4 min", available: "40,000.00" },
];

interface MarketOverviewProps {
    hideViewAll?: boolean;
}

const MarketOverview = ({ hideViewAll = false }: MarketOverviewProps) => {
    const [activeTab, setActiveTab] = useState("buy");
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedCoin, setSelectedCoin] = useState("USDT");
    const [selectedFiat, setSelectedFiat] = useState("USD");
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 6;

    const fiats = [
        { code: "USD", name: "US Dollar" },
        { code: "NGN", name: "Nigerian Naira" },
        { code: "PKR", name: "Pakistani Rupee" },
        { code: "EUR", name: "Euro" },
        { code: "GBP", name: "British Pound" },
    ];

    // Filter logic
    const filteredAds = useMemo(() => {
        return p2pAds.filter(ad =>
            (ad.user.toLowerCase().includes(searchQuery.toLowerCase()) ||
                ad.methods.some(m => m.toLowerCase().includes(searchQuery.toLowerCase()))) &&
            ad.type === activeTab &&
            ad.currency === selectedFiat &&
            ad.coin === selectedCoin
        );
    }, [searchQuery, activeTab, selectedFiat, selectedCoin]);

    // Reset page when filters change
    useEffect(() => {
        setCurrentPage(1);
    }, [searchQuery, activeTab, selectedFiat, selectedCoin]);

    // Pagination logic
    const totalPages = Math.ceil(filteredAds.length / itemsPerPage);
    const currentAds = filteredAds.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage
    );

    const handlePageChange = (page: number) => {
        setCurrentPage(page);
        // Smooth scroll to top of table/grid
        const element = document.getElementById("marketplace");
        if (element) {
            element.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    };

    return (
        <section className="py-24 bg-background relative overflow-hidden" id="marketplace">
            <div className="container mx-auto px-4 md:px-6">
                <div className="flex flex-col lg:flex-row lg:items-end justify-between mb-12 gap-8">
                    <div className="space-y-4">
                        <motion.div
                            initial={{ opacity: 0, x: -20 }}
                            whileInView={{ opacity: 1, x: 0 }}
                            className="inline-flex items-center space-x-2 bg-primary/10 border border-primary/20 px-3 py-1 rounded-full text-[10px] font-black text-primary uppercase tracking-[0.2em]"
                        >
                            <TrendingUp className="w-3 h-3" />
                            <span>Live Marketplace</span>
                        </motion.div>
                        <h2 className="text-4xl md:text-5xl font-black">P2P Marketplace</h2>
                        <p className="text-text-dim max-w-xl text-lg">
                            The most liquid decentralized exchange. Trade directly with verified merchants with zero hidden fees.
                        </p>
                    </div>

                    <div className="flex items-center space-x-2 bg-surface p-1.5 rounded-2xl border border-border self-start lg:self-auto shadow-2xl">
                        <button
                            onClick={() => setActiveTab("buy")}
                            className={`px-10 py-3 rounded-xl text-sm font-black transition-all ${activeTab === "buy" ? "bg-accent text-background shadow-lg shadow-accent/20" : "text-text-dim hover:text-white"
                                }`}
                        >
                            Buy
                        </button>
                        <button
                            onClick={() => setActiveTab("sell")}
                            className={`px-10 py-3 rounded-xl text-sm font-black transition-all ${activeTab === "sell" ? "bg-primary text-white shadow-lg shadow-primary/20" : "text-text-dim hover:text-white"
                                }`}
                        >
                            Sell
                        </button>
                    </div>
                </div>

                {/* Filters and Controls */}
                <div className="glass rounded-[32px] border-border mb-8 p-4 md:p-6 flex flex-col xl:flex-row items-center gap-4 shadow-2xl">
                    <div className="flex-1 w-full relative">
                        <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5" />
                        <input
                            type="text"
                            placeholder="Find merchant or payment method..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="w-full bg-surface/50 border border-border rounded-2xl py-4 pl-12 pr-4 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/50"
                        />
                    </div>
                    <div className="flex flex-wrap items-center gap-4 w-full xl:w-auto">
                        <div className="relative group flex-1 md:flex-none md:min-w-[140px]">
                            <select
                                value={selectedCoin}
                                onChange={(e) => setSelectedCoin(e.target.value)}
                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-5 text-base font-bold outline-none appearance-none cursor-pointer hover:border-primary transition-colors pr-10"
                            >
                                <option>USDT</option>
                                <option>BTC</option>
                                <option>ETH</option>
                            </select>
                            <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-dim pointer-events-none group-hover:text-primary transition-colors" />
                        </div>

                        <div className="relative group flex-1 md:flex-none md:min-w-[240px]">
                            <select
                                value={selectedFiat}
                                onChange={(e) => setSelectedFiat(e.target.value)}
                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-5 text-base font-bold outline-none appearance-none cursor-pointer hover:border-primary transition-colors pr-10"
                            >
                                {fiats.map(f => (
                                    <option key={f.code} value={f.code}>{f.name} ({f.code})</option>
                                ))}
                            </select>
                            <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-dim pointer-events-none group-hover:text-primary transition-colors" />
                        </div>

                        <button className="bg-surface/50 border border-border p-4 rounded-2xl hover:bg-primary/10 hover:border-primary transition-all group">
                            <Filter className="w-5 h-5 text-text-dim group-hover:text-primary transition-colors" />
                        </button>
                    </div>
                </div>

                {/* Table for Desktop */}
                <div className="hidden lg:block glass rounded-[40px] border-border overflow-hidden shadow-2xl relative">
                    <div className="overflow-x-auto">
                        <table className="w-full text-left border-collapse">
                            <thead>
                                <tr className="bg-surface/30 border-b border-border">
                                    <th className="px-8 py-6 text-xs font-black text-text-dim uppercase tracking-[0.15em]">Advertiser</th>
                                    <th className="px-8 py-6 text-xs font-black text-text-dim uppercase tracking-[0.15em]">Trade Price</th>
                                    <th className="px-8 py-6 text-xs font-black text-text-dim uppercase tracking-[0.15em]">Limits & Available</th>
                                    <th className="px-8 py-6 text-xs font-black text-text-dim uppercase tracking-[0.15em]">Payment Method</th>
                                    <th className="px-8 py-6 text-xs font-black text-text-dim uppercase tracking-[0.15em] text-right">Start Trade</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                <AnimatePresence mode="wait">
                                    {currentAds.length > 0 ? (
                                        <motion.tr
                                            key={`${activeTab}-${searchQuery}-${selectedFiat}-${selectedCoin}-${currentPage}`}
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                        >
                                            <td colSpan={5} className="p-0">
                                                <table className="w-full border-collapse">
                                                    <tbody className="divide-y divide-white/5">
                                                        {currentAds.map((ad, i) => (
                                                            <motion.tr
                                                                key={`${ad.user}-${i}`}
                                                                initial={{ opacity: 0, y: 10 }}
                                                                animate={{ opacity: 1, y: 0 }}
                                                                exit={{ opacity: 0, y: -10 }}
                                                                transition={{ duration: 0.2, delay: i * 0.05 }}
                                                                className="hover:bg-white/5 transition-all group cursor-pointer"
                                                            >
                                                                <td className="px-8 py-10">
                                                                    <div className="flex items-center space-x-4">
                                                                        <div className="w-14 h-14 rounded-2xl bg-surface border border-border flex items-center justify-center font-black text-xl text-primary shadow-inner">
                                                                            {ad.user[0]}
                                                                        </div>
                                                                        <div>
                                                                            <Link
                                                                                href={`/u/${ad.user}`}
                                                                                className="font-black text-lg text-white flex items-center gap-2 mb-1 hover:text-primary transition-colors cursor-pointer"
                                                                            >
                                                                                {ad.user} <CheckCircle2 className="w-4 h-4 text-accent fill-accent/10" />
                                                                            </Link>
                                                                            <div className="flex items-center space-x-3 text-xs font-bold text-text-dim">
                                                                                <span className="flex items-center gap-1"><Star className="w-3 h-3 text-amber-500 fill-amber-500" /> {ad.rating}</span>
                                                                                <span>{ad.trades} trades</span>
                                                                            </div>
                                                                        </div>
                                                                    </div>
                                                                </td>
                                                                <td className="px-8 py-10">
                                                                    <p className="text-3xl font-black text-white group-hover:text-primary transition-colors">{ad.price} <span className="text-xs font-bold text-text-dim uppercase tracking-widest">{ad.currency}</span></p>
                                                                    <div className="flex items-center gap-2 mt-2">
                                                                        <Clock className="w-3 h-3 text-text-dim" />
                                                                        <span className="text-[10px] font-bold text-text-dim uppercase">Avg speed: {ad.time}</span>
                                                                    </div>
                                                                </td>
                                                                <td className="px-8 py-10">
                                                                    <p className="text-sm font-bold text-white mb-2">Available: {ad.available} {ad.coin}</p>
                                                                    <div className="flex items-center space-x-2">
                                                                        <ShieldCheck className="w-3.5 h-3.5 text-accent" />
                                                                        <span className="text-xs font-medium text-text-dim">L: {ad.limits} {ad.currency}</span>
                                                                    </div>
                                                                </td>
                                                                <td className="px-8 py-10">
                                                                    <div className="flex flex-wrap gap-2">
                                                                        {ad.methods.map((m) => (
                                                                            <span key={m} className="px-4 py-1.5 rounded-full bg-surface-light text-[10px] font-black text-white uppercase border border-white/5 tracking-wider">
                                                                                {m}
                                                                            </span>
                                                                        ))}
                                                                    </div>
                                                                </td>
                                                                <td className="px-8 py-10 text-right">
                                                                    <button
                                                                        onClick={() => window.location.href = "/register"}
                                                                        className={`px-10 py-5 rounded-2xl font-black text-sm uppercase tracking-widest transition-all shadow-xl hover:scale-105 active:scale-95 ${activeTab === "buy" ? "bg-accent text-background shadow-accent/20 hover:shadow-accent/40" : "bg-primary text-white shadow-primary/20 hover:shadow-primary/40"
                                                                            }`}>
                                                                        {activeTab === "buy" ? `Buy ${ad.coin}` : `Sell ${ad.coin}`}
                                                                    </button>
                                                                </td>
                                                            </motion.tr>
                                                        ))}
                                                    </tbody>
                                                </table>
                                            </td>
                                        </motion.tr>
                                    ) : (
                                        <motion.tr
                                            key="empty-state"
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                        >
                                            <td colSpan={5} className="px-8 py-24 text-center">
                                                <div className="max-w-xs mx-auto">
                                                    <p className="text-xl font-bold text-white mb-2">No offers found</p>
                                                    <p className="text-text-dim">Try adjusting your filters or search terms to find more advertisements.</p>
                                                </div>
                                            </td>
                                        </motion.tr>
                                    )}
                                </AnimatePresence>
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Card Layout for Mobile/Tablet */}
                <div className="lg:hidden space-y-6">
                    <AnimatePresence mode="popLayout">
                        {currentAds.length > 0 ? (
                            currentAds.map((ad, i) => (
                                <motion.div
                                    key={`${ad.user}-${i}`}
                                    initial={{ opacity: 0, y: 20 }}
                                    whileInView={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, scale: 0.95 }}
                                    viewport={{ once: true }}
                                    className="glass rounded-[32px] border-border p-6 space-y-6"
                                >
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center space-x-3">
                                            <div className="w-12 h-12 rounded-xl bg-surface border border-border flex items-center justify-center font-black text-primary">
                                                {ad.user[0]}
                                            </div>
                                            <div>
                                                <Link
                                                    href={`/u/${ad.user}`}
                                                    className="font-black text-white flex items-center gap-1 hover:text-primary transition-colors"
                                                >
                                                    {ad.user} <CheckCircle2 className="w-4 h-4 text-accent" />
                                                </Link>
                                                <p className="text-[10px] font-bold text-text-dim uppercase tracking-wider">{ad.trades} TRADES | {ad.rating} RATING</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-xs font-black text-text-dim uppercase tracking-widest mb-1">Time</p>
                                            <p className="text-sm font-bold text-white">{ad.time}</p>
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-2 gap-4 pt-4 border-t border-white/5">
                                        <div>
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1">Price</p>
                                            <p className="text-xl font-black text-white">{ad.price} {ad.currency}</p>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-1">Limits</p>
                                            <p className="text-sm font-bold text-white">{ad.limits} {ad.currency}</p>
                                        </div>
                                    </div>

                                    <div className="pb-2">
                                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-3">Accepted Payments</p>
                                        <div className="flex flex-wrap gap-2">
                                            {ad.methods.map((m) => (
                                                <span key={m} className="px-3 py-1 rounded-lg bg-surface-light text-[10px] font-bold text-white uppercase border border-border">
                                                    {m}
                                                </span>
                                            ))}
                                        </div>
                                    </div>

                                    <button
                                        onClick={() => window.location.href = "/register"}
                                        className={`w-full py-5 rounded-2xl font-black text-sm uppercase tracking-widest transition-all ${activeTab === "buy" ? "bg-accent text-background" : "bg-primary text-white"}`}>
                                        {activeTab === "buy" ? `Buy ${ad.coin}` : `Sell ${ad.coin}`}
                                    </button>
                                </motion.div>
                            ))
                        ) : (
                            <div className="glass rounded-[32px] border-border p-12 text-center">
                                <p className="text-xl font-bold text-white mb-2">No offers found</p>
                                <p className="text-text-dim text-sm">Try adjusting your filters.</p>
                            </div>
                        )}
                    </AnimatePresence>
                </div>

                {/* Pagination component */}
                <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    onPageChange={handlePageChange}
                />

                {!hideViewAll && (
                    <div className="mt-12 text-center">
                        <a
                            href="/marketplace"
                            className="text-primary font-black text-sm uppercase tracking-[0.2em] inline-flex items-center space-x-3 mx-auto group hover:tracking-[0.3em] transition-all"
                        >
                            <span>View all market offers</span>
                            <ArrowUpRight className="w-5 h-5 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
                        </a>
                    </div>
                )}
            </div>
        </section>
    );
};

export default MarketOverview;
