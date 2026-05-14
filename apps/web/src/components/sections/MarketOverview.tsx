"use client";

import React, { useState, useMemo, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Search, Filter, ArrowUpRight, ChevronDown, Clock, Star, TrendingUp } from "lucide-react";
import { toast } from "sonner";
import Link from "next/link";
import Pagination from "@/components/ui/Pagination";
import { MarketplaceService, AdResponse } from "@/services/marketplaceService";
import { TradeService } from "@/services/tradeService";
import { useAuth } from "@/context/AuthContext";

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
    const [activeTab, setActiveTab] = useState<"buy" | "sell">("buy");
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedCoin, setSelectedCoin] = useState("USDT");
    const [selectedFiat, setSelectedFiat] = useState("USD");
    const [isCryptoOpen, setIsCryptoOpen] = useState(false);
    const [isFiatOpen, setIsFiatOpen] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);

    const closeDropdowns = () => {
        setIsCryptoOpen(false);
        setIsFiatOpen(false);
    };

    const [ads, setAds] = useState<AdResponse[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const itemsPerPage = 6;
    const [selectedAd, setSelectedAd] = useState<AdResponse | null>(null);
    const [tradeAmount, setTradeAmount] = useState<string>("");
    const [isInitiating, setIsInitiating] = useState(false);
    const [initiationError, setInitiationError] = useState<string | null>(null);

    useEffect(() => {
        const fetchAds = async () => {
            setIsLoading(true);
            try {
                const data = await MarketplaceService.getAds({
                    type: activeTab,
                } as any);

                if (!data || !data.ads || data.ads.length === 0) {
                    const mockMapped: AdResponse[] = p2pAds
                        .filter(ad => ad.type === activeTab && ad.currency === selectedFiat && ad.coin === selectedCoin)
                        .map((ad, idx) => ({
                            ad_id: `mock-${idx}`,
                            user_id: `user-${ad.user}`,
                            username: ad.user,
                            user_trades: parseInt(ad.trades.replace(/\D/g, "")),
                            user_rating: parseFloat(ad.rating.replace("%", "")),
                            type: ad.type as "buy" | "sell",
                            crypto_symbol: ad.coin,
                            fiat_symbol: ad.currency,
                            price: parseFloat(ad.price.replace(/,/g, "")),
                            min_amount: parseFloat(ad.limits.split("-")[0].replace(/,/g, "")),
                            max_amount: parseFloat(ad.limits.split("-")[1].replace(/,/g, "")),
                            payment_methods: ad.methods,
                            payment_method_ids: ad.methods.map((_, i) => i + 1), // Dummy IDs for mock
                            payment_window_minutes: parseInt(ad.time.split(" ")[0]),
                            is_online: true,
                            price_type: "fixed",
                            status: "active",
                            created_at: new Date().toISOString()
                        }));
                    setAds(mockMapped);
                } else {
                    setAds(data.ads);
                }
            } catch (error) {
                // Fallen back to mock data if API is temporarily unavailable
                const mockMapped: AdResponse[] = p2pAds
                    .filter(ad => ad.type === activeTab && ad.currency === selectedFiat && ad.coin === selectedCoin)
                    .map((ad, idx) => ({
                        ad_id: `mock-${idx}`,
                        user_id: `user-${ad.user}`,
                        username: ad.user,
                        user_trades: parseInt(ad.trades.replace(/\D/g, "")),
                        user_rating: parseFloat(ad.rating.replace("%", "")),
                        type: ad.type as "buy" | "sell",
                        crypto_symbol: ad.coin,
                        fiat_symbol: ad.currency,
                        price: parseFloat(ad.price.replace(/,/g, "")),
                        min_amount: parseFloat(ad.limits.split("-")[0].replace(/,/g, "")),
                        max_amount: parseFloat(ad.limits.split("-")[1].replace(/,/g, "")),
                        payment_methods: ad.methods,
                        payment_method_ids: ad.methods.map((_, i) => i + 1), // Dummy IDs for mock
                        payment_window_minutes: parseInt(ad.time.split(" ")[0]),
                        is_online: true,
                        price_type: "fixed",
                        status: "active",
                        created_at: new Date().toISOString()
                    }));
                setAds(mockMapped);
            } finally {
                setIsLoading(false);
            }
        };

        fetchAds();
    }, [activeTab, selectedCoin, selectedFiat]);

    const fiats = [
        { code: "USD", name: "US Dollar" },
        { code: "NGN", name: "Nigerian Naira" },
        { code: "PKR", name: "Pakistani Rupee" },
        { code: "EUR", name: "Euro" },
        { code: "GBP", name: "British Pound" },
    ];

    const filteredAds = useMemo(() => {
        return ads.filter(ad =>
        (ad.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (ad.payment_methods && ad.payment_methods.some((m: string) => m.toLowerCase().includes(searchQuery.toLowerCase()))))
        );
    }, [searchQuery, ads]);

    useEffect(() => {
        setCurrentPage(1);
    }, [searchQuery, activeTab, selectedFiat, selectedCoin]);

    const totalPages = Math.ceil(filteredAds.length / itemsPerPage);
    const currentAds = filteredAds.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage
    );

    const handlePageChange = (page: number) => {
        setCurrentPage(page);
        const element = document.getElementById("marketplace");
        if (element) {
            element.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    };

    const { user: authUser, isLoading: authLoading } = useAuth();

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
                        <h2 className="text-5xl md:text-7xl font-black italic uppercase tracking-tighter">P2P <span className="text-primary">MARKET</span></h2>
                        <p className="text-text-dim max-w-xl text-lg font-medium leading-tight">
                            The clearing layer for decentralized trade. Execute directly with verified liquidity providers with absolute transparency.
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

                <div className="glass rounded-[32px] border-border mb-8 p-4 md:p-6 flex flex-col xl:flex-row items-center gap-4 shadow-2xl">
                    <div className="flex-1 w-full relative">
                        <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5" />
                        <input
                            type="text"
                            placeholder="Find offer or payment method..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="w-full bg-surface/50 border border-border rounded-2xl py-4 pl-12 pr-4 text-base outline-none focus:border-primary transition-all font-medium placeholder:text-text-dim/50"
                        />
                    </div>
                    <div className="flex flex-wrap items-center gap-4 w-full xl:w-auto relative">
                        {/* Global click catcher */}
                        {(isCryptoOpen || isFiatOpen) && (
                            <div className="fixed inset-0 z-30" onClick={closeDropdowns} />
                        )}

                        <div className="relative group flex-1 md:flex-none md:min-w-[140px] z-50">
                            <button
                                type="button"
                                onClick={() => { closeDropdowns(); setIsCryptoOpen(!isCryptoOpen); }}
                                className="w-full h-[56px] px-5 flex items-center justify-between bg-surface/50 border border-border hover:border-primary/50 rounded-2xl text-base font-bold text-white outline-none focus:border-primary transition-all pr-4"
                            >
                                <span>{selectedCoin}</span>
                                <ChevronDown className={`w-5 h-5 text-text-dim transition-transform ${isCryptoOpen ? "rotate-180" : ""}`} />
                            </button>

                            {isCryptoOpen && (
                                <div className="absolute w-full mt-2 bg-surface bg-opacity-95 backdrop-blur-xl border border-border/50 rounded-2xl shadow-2xl overflow-hidden py-1 z-50">
                                    {["USDT", "BTC", "ETH"].map(coin => (
                                        <button
                                            key={coin}
                                            type="button"
                                            onClick={() => { setSelectedCoin(coin); setIsCryptoOpen(false); }}
                                            className={`w-full flex items-center justify-between px-5 py-3.5 text-sm transition-all text-left ${selectedCoin === coin ? "bg-primary/20 text-white font-black" : "text-text-dim font-bold hover:text-white hover:bg-white/5"}`}
                                        >
                                            {coin}
                                        </button>
                                    ))}
                                </div>
                            )}
                        </div>

                        <div className="relative group flex-1 md:flex-none md:min-w-[240px] z-40">
                            <button
                                type="button"
                                onClick={() => { closeDropdowns(); setIsFiatOpen(!isFiatOpen); }}
                                className="w-full h-[56px] px-5 flex items-center justify-between bg-surface/50 border border-border hover:border-primary/50 rounded-2xl text-base font-bold text-white outline-none focus:border-primary transition-all pr-4"
                            >
                                <div className="flex items-baseline gap-2 truncate pr-2">
                                    <span>{selectedFiat}</span>
                                    <span className="text-xs text-text-dim font-medium truncate">
                                        {fiats.find(f => f.code === selectedFiat)?.name}
                                    </span>
                                </div>
                                <ChevronDown className={`shrink-0 w-5 h-5 text-text-dim transition-transform ${isFiatOpen ? "rotate-180" : ""}`} />
                            </button>

                            {isFiatOpen && (
                                <div className="absolute w-full mt-2 bg-surface bg-opacity-95 backdrop-blur-xl border border-border/50 rounded-2xl shadow-2xl overflow-hidden py-1 z-50">
                                    {fiats.map(f => (
                                        <button
                                            key={f.code}
                                            type="button"
                                            onClick={() => { setSelectedFiat(f.code); setIsFiatOpen(false); }}
                                            className={`w-full flex items-center justify-between px-5 py-3.5 text-sm transition-all text-left group-item ${selectedFiat === f.code ? "bg-primary/20 text-white font-black" : "text-text-dim font-bold hover:text-white hover:bg-white/5"}`}
                                        >
                                            <span>{f.code}</span>
                                            <span className={`text-[10px] font-medium transition-colors ${selectedFiat === f.code ? "text-primary/80" : "text-white/30"}`}>{f.name}</span>
                                        </button>
                                    ))}
                                </div>
                            )}
                        </div>

                        <button
                            onClick={() => toast.info("Advanced filtering coming soon")}
                            className="bg-surface/50 border border-border p-4 rounded-2xl hover:bg-primary/10 hover:border-primary transition-all group"
                        >
                            <Filter className="w-5 h-5 text-text-dim group-hover:text-primary transition-colors" />
                        </button>
                    </div>
                </div>

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
                                <AnimatePresence>
                                    {isLoading ? (
                                        <tr><td colSpan={5} className="p-20 text-center text-text-dim">Loading offers...</td></tr>
                                    ) : currentAds.length > 0 ? (
                                        currentAds.map((ad, i) => (
                                            <motion.tr
                                                key={ad.ad_id}
                                                initial={{ opacity: 0, y: 10 }}
                                                animate={{ opacity: 1, y: 0 }}
                                                exit={{ opacity: 0, y: -10 }}
                                                transition={{ duration: 0.2, delay: i * 0.05 }}
                                                className="hover:bg-white/5 transition-all group cursor-pointer"
                                            >
                                                <td className="px-8 py-10">
                                                    <div className="flex items-center space-x-4">
                                                        <div className="w-14 h-14 rounded-2xl bg-surface border border-border flex items-center justify-center font-black text-xl text-primary shadow-inner">
                                                            {ad.username[0]}
                                                        </div>
                                                        <div>
                                                            <Link href={`/u/${ad.username}`} className="font-black text-lg text-white flex items-center gap-2 mb-1 hover:text-primary transition-colors cursor-pointer">
                                                                {ad.username}
                                                            </Link>
                                                            <div className="flex items-center space-x-3 text-xs font-bold text-text-dim">
                                                                <span className="flex items-center gap-1"><Star className="w-3 h-3 text-amber-500 fill-amber-500" /> {ad.user_rating ?? 0}{ad.user_rating && ad.user_rating > 10 ? "%" : ""}</span>
                                                                <span>{ad.user_trades ?? 0} trades</span>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-8 py-10">
                                                    <p className="text-3xl font-black text-white group-hover:text-primary transition-colors">{ad.price.toLocaleString()} <span className="text-xs font-bold text-text-dim uppercase tracking-widest">{ad.fiat_symbol}</span></p>
                                                    <div className="flex items-center gap-2 mt-2">
                                                        <Clock className="w-3 h-3 text-text-dim" />
                                                        <span className="text-[10px] font-bold text-text-dim uppercase">Avg speed: {ad.payment_window_minutes} min</span>
                                                    </div>
                                                </td>
                                                <td className="px-8 py-10">
                                                    <p className="text-sm font-bold text-white mb-2">Limits: {ad.min_amount.toLocaleString()} - {ad.max_amount.toLocaleString()} {ad.fiat_symbol}</p>
                                                    <div className="flex items-center space-x-2">
                                                        <span className="text-xs font-medium text-text-dim">Offer Available</span>
                                                    </div>
                                                </td>
                                                <td className="px-8 py-10">
                                                    <div className="flex flex-wrap gap-2">
                                                        {ad.payment_methods?.map((m: string, idx: number) => (
                                                            <span key={idx} className="px-4 py-1.5 rounded-full bg-surface-light text-[10px] font-black text-white uppercase border border-white/5 tracking-wider">
                                                                {m}
                                                            </span>
                                                        ))}
                                                    </div>
                                                </td>
                                                <td className="px-8 py-10 text-right">
                                                    <button
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            setSelectedAd(ad);
                                                            setTradeAmount(ad.min_amount.toString());
                                                            setInitiationError(null);
                                                        }}
                                                        className={`px-10 py-5 rounded-2xl font-black text-sm uppercase tracking-widest transition-all shadow-xl hover:scale-105 active:scale-95 ${activeTab === "buy" ? "bg-accent text-background shadow-accent/20 hover:shadow-accent/40" : "bg-primary text-white shadow-primary/20 hover:shadow-primary/40"}`}
                                                    >
                                                        {activeTab === "buy" ? `Buy ${ad.crypto_symbol}` : `Sell ${ad.crypto_symbol}`}
                                                    </button>
                                                </td>
                                            </motion.tr>
                                        ))
                                    ) : (
                                        <motion.tr key="empty-state" initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
                                            <td colSpan={5} className="px-8 py-24 text-center">
                                                <p className="text-xl font-bold text-white mb-2">No offers found</p>
                                                <p className="text-text-dim">Try adjusting your filters.</p>
                                            </td>
                                        </motion.tr>
                                    )}
                                </AnimatePresence>
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="lg:hidden space-y-6">
                    <AnimatePresence>
                        {isLoading ? (
                            <div className="text-center py-10 text-text-dim">Loading...</div>
                        ) : currentAds.length > 0 ? (
                            currentAds.map((ad, i) => (
                                <motion.div key={ad.ad_id} initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} exit={{ opacity: 0, scale: 0.95 }} viewport={{ once: true }} className="glass rounded-[1.5rem] border-border p-5 relative">
                                    <div className="flex items-center justify-between mb-4">
                                        <div className="flex items-center gap-3">
                                            <div className="w-12 h-12 rounded-xl bg-surface border border-border flex items-center justify-center font-black text-primary">{ad.username[0]}</div>
                                            <div>
                                                <Link href={`/u/${ad.username}`} className="font-black text-white flex items-center gap-1 hover:text-primary transition-colors">
                                                    {ad.username}
                                                </Link>
                                                <p className="text-[10px] font-bold text-text-dim uppercase tracking-wider">{ad.user_trades ?? 0} TRADES | {ad.user_rating ?? 0}{ad.user_rating && ad.user_rating > 10 ? "%" : ""} RATING</p>
                                            </div>
                                        </div>
                                        <div className="bg-accent/5 px-3 py-1 rounded-lg border border-accent/10">
                                            <span className="text-[10px] font-black text-accent">{ad.payment_window_minutes} MIN</span>
                                        </div>
                                    </div>
                                    <div className="mb-4 bg-white/5 p-4 rounded-2xl">
                                        <div className="flex justify-between items-baseline mb-1">
                                            <span className="text-[10px] font-black text-text-dim uppercase tracking-widest">Price</span>
                                            <span className="text-xl font-black text-white">{ad.price.toLocaleString()} {ad.fiat_symbol}</span>
                                        </div>
                                        <div className="flex justify-between text-[10px] font-bold">
                                            <span className="text-text-dim">Limits</span>
                                            <span className="text-white">{ad.min_amount.toLocaleString()} - {ad.max_amount.toLocaleString()} {ad.fiat_symbol}</span>
                                        </div>
                                    </div>
                                    <div className="flex gap-2 mb-5 overflow-x-auto pb-1 no-scrollbar">
                                        {ad.payment_methods?.map((method: string, idx: number) => (
                                            <span key={idx} className="px-2 py-1 bg-white/5 border border-white/10 rounded-md text-[8px] font-bold text-text-dim uppercase whitespace-nowrap">{method}</span>
                                        ))}
                                    </div>
                                    <button
                                        onClick={() => {
                                            setSelectedAd(ad);
                                            setTradeAmount(ad.min_amount.toString());
                                            setInitiationError(null);
                                        }}
                                        className="w-full py-4 bg-primary text-white rounded-xl font-black uppercase tracking-widest text-[10px] shadow-lg shadow-primary/20"
                                    >
                                        {activeTab === "buy" ? "Buy" : "Sell"} {ad.crypto_symbol}
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

                {/* Trade Initiation Modal */}
                <AnimatePresence>
                    {selectedAd && (
                        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
                            <motion.div
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                onClick={() => setSelectedAd(null)}
                                className="absolute inset-0 bg-background/80 backdrop-blur-md"
                            />
                            <motion.div
                                initial={{ scale: 0.9, opacity: 0, y: 20 }}
                                animate={{ scale: 1, opacity: 1, y: 0 }}
                                exit={{ scale: 0.9, opacity: 0, y: 20 }}
                                className="glass border-border p-8 md:p-10 rounded-[3rem] max-w-lg w-full relative overflow-hidden"
                            >
                                <div className="flex justify-between items-start mb-8">
                                    <div>
                                        <h3 className="text-2xl font-black text-white mb-2 uppercase tracking-tight">Initiate Trade</h3>
                                        <p className="text-text-dim text-sm font-bold uppercase tracking-widest">Trading with <span className="text-primary">{selectedAd.username}</span></p>
                                    </div>
                                    <button onClick={() => setSelectedAd(null)} className="w-10 h-10 rounded-full border border-white/10 flex items-center justify-center text-white hover:bg-white/10 transition-colors">×</button>
                                </div>

                                <div className="space-y-6 mb-10">
                                    <div>
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-3 block">Amount to {activeTab === "buy" ? "Pay" : "Receive"} ({selectedAd.fiat_symbol})</label>
                                        <div className="relative">
                                            <input
                                                type="text"
                                                value={tradeAmount}
                                                onChange={(e) => setTradeAmount(e.target.value)}
                                                placeholder={`Min: ${selectedAd.min_amount}`}
                                                className="w-full bg-white/5 border-2 border-white/10 rounded-2xl py-5 px-6 text-xl font-black text-white focus:border-primary outline-none transition-all pr-16"
                                            />
                                            <div className="absolute right-6 top-1/2 -translate-y-1/2 font-black text-text-dim">{selectedAd.fiat_symbol}</div>
                                        </div>
                                        <p className="mt-3 text-[10px] font-bold text-text-dim uppercase tracking-wider">
                                            Approx. {selectedAd && (parseFloat(tradeAmount || "0") / selectedAd.price).toFixed(8)} {selectedAd.crypto_symbol}
                                        </p>
                                    </div>

                                    {initiationError && (
                                        <motion.div
                                            initial={{ opacity: 0, height: 0 }}
                                            animate={{ opacity: 1, height: "auto" }}
                                            className="p-4 bg-red-500/10 border border-red-500/20 rounded-2xl"
                                        >
                                            <p className="text-xs font-black text-red-500 uppercase tracking-widest leading-loose">
                                                ⚠️ Error: {initiationError}
                                            </p>
                                        </motion.div>
                                    )}
                                </div>

                                <div className="flex gap-4">
                                    <button
                                        disabled={isInitiating}
                                        onClick={async () => {
                                            if (authLoading) return;
                                            if (!authUser) {
                                                toast.error("Please login to trade");
                                                window.location.href = `/login?redirect=${encodeURIComponent(window.location.pathname)}`;
                                                return;
                                            }
                                            setIsInitiating(true);
                                            setInitiationError(null);
                                            try {
                                                const amount = parseFloat(tradeAmount);
                                                const pmId = selectedAd.payment_method_ids ? selectedAd.payment_method_ids[0] : 0;
                                                const result = await TradeService.initiateTrade(selectedAd.ad_id, amount, pmId);
                                                window.location.href = `/user/dashboard/trades/${result.trade_id}`;
                                            } catch (err: any) {
                                                setInitiationError(err.message);
                                            } finally {
                                                setIsInitiating(false);
                                            }
                                        }}
                                        className="flex-1 py-6 bg-primary text-white rounded-[2rem] font-black uppercase tracking-widest text-xs hover:scale-105 active:scale-95 transition-all shadow-xl shadow-primary/20 disabled:opacity-50"
                                    >
                                        {isInitiating ? "Initiating..." : `Start ${activeTab === "buy" ? "Buy" : "Sell"} Trade`}
                                    </button>
                                </div>
                            </motion.div>
                        </div>
                    )}
                </AnimatePresence>

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
