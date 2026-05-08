"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import { motion } from "framer-motion";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { 
    Search, 
    Filter, 
    ArrowRight, 
    TrendingUp, 
    TrendingDown, 
    User, 
    Clock,
    Shield,
    Star,
    ChevronDown,
    Loader2,
    Plus
} from "lucide-react";
import { authService } from "@/services/authService";
import { toast } from "sonner";

interface TradeAd {
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
    payment_window_minutes: number;
    trade_terms?: string;
    status: "active" | "paused" | "closed";
    created_at: string;
    user_trades?: number;
    user_rating?: number;
}

const ADS_PER_PAGE = 12;

const MarketplacePage = () => {
    const router = useRouter();
    const [ads, setAds] = useState<TradeAd[]>([]);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [offset, setOffset] = useState(0);
    const [total, setTotal] = useState(0);
    const [searchTerm, setSearchTerm] = useState("");
    const [filters, setFilters] = useState({
        type: "all",
        fiat_currency: "all",
        payment_method: "all",
        sort_by: "best_price"
    });
    const [showFilters, setShowFilters] = useState(false);
    const observerRef = useRef<IntersectionObserver | null>(null);
    const loadMoreRef = useRef<HTMLDivElement>(null);

    // Reset and fetch when filters change
    useEffect(() => {
        setAds([]);
        setOffset(0);
        setHasMore(true);
        fetchAds(0, true);
    }, [filters]);

    const fetchAds = async (currentOffset: number, isInitial: boolean = false) => {
        if (isInitial) {
            setLoading(true);
        } else {
            setLoadingMore(true);
        }
        
        try {
            const params = new URLSearchParams();
            params.append("limit", ADS_PER_PAGE.toString());
            params.append("offset", currentOffset.toString());
            if (filters.type !== "all") params.append("type", filters.type);
            if (filters.fiat_currency !== "all") params.append("fiat_currency", filters.fiat_currency);
            if (filters.payment_method !== "all") params.append("payment_method", filters.payment_method);
            
            const response = await fetch(`/api/v1/marketplace/ads?${params.toString()}`);
            const data = await response.json();
            
            if (response.ok) {
                const newAds = data.ads || [];
                setTotal(data.total || 0);
                
                if (isInitial) {
                    setAds(newAds);
                } else {
                    setAds(prev => [...prev, ...newAds]);
                }
                
                // Check if there are more ads to load
                const loadedCount = isInitial ? newAds.length : currentOffset + newAds.length;
                setHasMore(loadedCount < (data.total || 0));
                setOffset(currentOffset + newAds.length);
            } else {
                throw new Error(data.message || "Failed to fetch ads");
            }
        } catch (error: any) {
            console.error("Error fetching ads:", error);
            toast.error("Failed to load marketplace ads");
        } finally {
            setLoading(false);
            setLoadingMore(false);
        }
    };

    // Intersection Observer for infinite scroll
    const handleObserver = useCallback((entries: IntersectionObserverEntry[]) => {
        const target = entries[0];
        if (target.isIntersecting && hasMore && !loadingMore && !loading) {
            fetchAds(offset);
        }
    }, [offset, hasMore, loadingMore, loading]);

    useEffect(() => {
        const option = {
            root: null,
            rootMargin: "100px",
            threshold: 0
        };
        
        observerRef.current = new IntersectionObserver(handleObserver, option);
        
        if (loadMoreRef.current) {
            observerRef.current.observe(loadMoreRef.current);
        }
        
        return () => {
            if (observerRef.current) {
                observerRef.current.disconnect();
            }
        };
    }, [handleObserver]);

    const filteredAds = ads.filter(ad => 
        ad.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ad.trade_terms?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const handleInitiateTrade = async (adId: string) => {
        try {
            const response = await fetch(`/api/v1/marketplace/ads/${adId}/trades`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            });
            
            if (response.ok) {
                const data = await response.json();
                toast.success("Trade initiated successfully");
                router.push(`/dashboard/trades/${data.trade_id}`);
            } else {
                const error = await response.json();
                throw new Error(error.message || "Failed to initiate trade");
            }
        } catch (error: any) {
            console.error("Error initiating trade:", error);
            toast.error(error.message || "Failed to initiate trade");
        }
    };

    const formatPrice = (ad: TradeAd) => {
        return `${ad.price.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ${ad.fiat_symbol}`;
    };

    const formatAmount = (amount: number) => {
        return amount.toLocaleString("en-US", { minimumFractionDigits: 0, maximumFractionDigits: 2 });
    };

    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            <div className="container mx-auto px-4 md:px-6 pt-24 pb-8">
                {/* Header */}
                <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6 mb-8">
                    <div>
                        <h1 className="text-4xl font-black text-white mb-2">P2P Marketplace</h1>
                        <p className="text-text-dim">Trade USDT with trusted peers</p>
                    </div>
                    <Link
                        href="/marketplace/create"
                        className="inline-flex items-center px-6 py-3 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-colors"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Create Ad
                    </Link>
                </div>

                {/* Search and Filters */}
                <div className="bg-surface border border-white/10 rounded-2xl p-6 mb-8">
                    <div className="flex flex-col lg:flex-row gap-4">
                        {/* Search */}
                        <div className="flex-1 relative">
                            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-dim" />
                            <input
                                type="text"
                                placeholder="Search by username or trade terms..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                            />
                        </div>

                        {/* Filter Toggle */}
                        <button
                            onClick={() => setShowFilters(!showFilters)}
                            className="px-6 py-3 bg-white/5 border border-white/10 rounded-xl text-white hover:bg-white/10 transition-colors flex items-center gap-2"
                        >
                            <Filter className="w-4 h-4" />
                            Filters
                            <ChevronDown className={`w-4 h-4 transition-transform ${showFilters ? "rotate-180" : ""}`} />
                        </button>
                    </div>

                    {/* Advanced Filters */}
                    {showFilters && (
                        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mt-6 pt-6 border-t border-white/10">
                            <div>
                                <label className="block text-text-dim text-sm mb-2">Trade Type</label>
                                <select
                                    value={filters.type}
                                    onChange={(e) => setFilters({...filters, type: e.target.value})}
                                    className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                                >
                                    <option value="all">All Types</option>
                                    <option value="buy">Buy USDT</option>
                                    <option value="sell">Sell USDT</option>
                                </select>
                            </div>
                            
                            <div>
                                <label className="block text-text-dim text-sm mb-2">Currency</label>
                                <select
                                    value={filters.fiat_currency}
                                    onChange={(e) => setFilters({...filters, fiat_currency: e.target.value})}
                                    className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                                >
                                    <option value="all">All Currencies</option>
                                    <option value="USD">USD</option>
                                    <option value="BDT">BDT</option>
                                    <option value="PKR">PKR</option>
                                </select>
                            </div>
                            
                            <div>
                                <label className="block text-text-dim text-sm mb-2">Payment Method</label>
                                <select
                                    value={filters.payment_method}
                                    onChange={(e) => setFilters({...filters, payment_method: e.target.value})}
                                    className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                                >
                                    <option value="all">All Methods</option>
                                    <option value="bkash">Bkash</option>
                                    <option value="nagad">Nagad</option>
                                    <option value="bank">Bank Transfer</option>
                                </select>
                            </div>
                            
                            <div>
                                <label className="block text-text-dim text-sm mb-2">Sort By</label>
                                <select
                                    value={filters.sort_by}
                                    onChange={(e) => setFilters({...filters, sort_by: e.target.value})}
                                    className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                                >
                                    <option value="best_price">Best Price</option>
                                    <option value="completion_rate">Completion Rate</option>
                                    <option value="newest">Newest</option>
                                    <option value="trade_count">Most Trades</option>
                                </select>
                            </div>
                        </div>
                    )}
                </div>

                {/* Ads List */}
                {loading ? (
                    <div className="flex items-center justify-center h-64">
                        <Loader2 className="w-8 h-8 animate-spin text-primary" />
                    </div>
                ) : filteredAds.length === 0 ? (
                    <div className="text-center py-12">
                        <p className="text-text-dim">No ads found matching your criteria</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {filteredAds.map((ad) => (
                            <motion.div
                                key={ad.ad_id}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="bg-surface border border-white/10 rounded-2xl p-6 hover:border-primary/50 transition-all"
                            >
                                {/* Header */}
                                <div className="flex items-center justify-between mb-4">
                                    <div className="flex items-center gap-2">
                                        <div className={`px-3 py-1 rounded-full text-xs font-black uppercase ${
                                            ad.type === "buy" 
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
                                <div className="flex items-center justify-between mb-4">
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
                                    onClick={() => handleInitiateTrade(ad.ad_id)}
                                    className="w-full py-3 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-colors flex items-center justify-center gap-2"
                                >
                                    Trade Now
                                    <ArrowRight className="w-4 h-4" />
                                </button>
                            </motion.div>
                        ))}
                    </div>
                )}

                {/* Infinite Scroll Sentinel */}
                <div ref={loadMoreRef} className="h-20 flex items-center justify-center">
                    {loadingMore && (
                        <div className="flex items-center gap-2 text-text-dim">
                            <Loader2 className="w-5 h-5 animate-spin" />
                            <span className="text-sm">Loading more...</span>
                        </div>
                    )}
                    {!hasMore && ads.length > 0 && (
                        <p className="text-text-dim text-sm">
                            Showing {filteredAds.length} of {total} ads
                        </p>
                    )}
                </div>
            </div>

            <Footer />
        </main>
    );
};

export default MarketplacePage;
