"use client";

import React, { useState, useEffect, useRef, useCallback } from "react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Plus, Loader2, ArrowUpRight } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";

import AdCard, { TradeAd } from "@/components/marketplace/AdCard";
import MarketplaceFilters from "@/components/marketplace/MarketplaceFilters";
import { MarketplaceService, AdResponse } from "@/services/marketplaceService";
import { TradeService } from "@/services/tradeService";

const ADS_PER_PAGE = 12;

const MarketplacePage = () => {
    const router = useRouter();
    const [ads, setAds] = useState<AdResponse[]>([]);
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
    const [selectedAd, setSelectedAd] = useState<AdResponse | null>(null);
    const [tradeAmount, setTradeAmount] = useState("");
    const [isInitiating, setIsInitiating] = useState(false);
    const [initiationError, setInitiationError] = useState<string | null>(null);
    const observerRef = useRef<IntersectionObserver | null>(null);
    const loadMoreRef = useRef<HTMLDivElement>(null);

    const fetchAds = useCallback(async (currentOffset: number, isInitial: boolean = false) => {
        if (isInitial) {
            setLoading(true);
        } else {
            setLoadingMore(true);
        }

        try {
            const params: Record<string, string> = {
                limit: ADS_PER_PAGE.toString(),
                offset: currentOffset.toString(),
                sort_by: filters.sort_by
            };
            if (filters.type !== "all") params.type = filters.type;
            if (filters.fiat_currency !== "all") params.fiat_currency = filters.fiat_currency;
            if (filters.payment_method !== "all") params.payment_method = filters.payment_method;

            const data = await MarketplaceService.getAds(params);

            const newAds = data.ads || [];
            setTotal(data.total || 0);

            if (isInitial) {
                setAds(newAds);
            } else {
                setAds(prev => [...prev, ...newAds]);
            }

            const loadedCount = isInitial ? newAds.length : currentOffset + newAds.length;
            setHasMore(loadedCount < (data.total || 0));
            setOffset(loadedCount);
        } catch (error: any) {
            console.error("Error fetching ads:", error);
            toast.error("Failed to load marketplace ads");
        } finally {
            setLoading(false);
            setLoadingMore(false);
        }
    }, [filters]);

    // Reset and fetch when filters change
    useEffect(() => {
        setOffset(0);
        fetchAds(0, true);
    }, [fetchAds]);

    // Intersection Observer for infinite scroll
    const handleObserver = useCallback((entries: IntersectionObserverEntry[]) => {
        const target = entries[0];
        if (target.isIntersecting && hasMore && !loadingMore && !loading) {
            fetchAds(offset);
        }
    }, [offset, hasMore, loadingMore, loading, fetchAds]);

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

    const { user: authUser, isLoading: authLoading } = useAuth();

    const handleInitiateTrade = async (adId: string) => {
        if (authLoading) return;

        if (!authUser) {
            toast.error("Please login to trade");
            router.push(`/login?redirect=${encodeURIComponent(window.location.pathname)}`);
            return;
        }

        const ad = ads.find(a => a.ad_id === adId);
        if (ad) {
            setSelectedAd(ad);
            setTradeAmount(ad.min_amount.toString());
            setInitiationError(null);
        }
    };

    const filteredAds = ads.filter(ad =>
        ad.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ad.trade_terms?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            <div className="container mx-auto px-4 md:px-6 pt-24 pb-8">
                {/* Header */}
                <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6 mb-8">
                    <div>
                        <h1 className="text-4xl md:text-5xl font-black text-white mb-2 italic">MARKETPLACE</h1>
                        <p className="text-text-dim uppercase tracking-[0.2em] text-sm">P2P Peer-to-Peer Trading</p>
                    </div>
                    <Link
                        href="/marketplace/create"
                        className="inline-flex items-center px-8 py-4 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:translate-y-[-2px] hover:shadow-[0_8px_20px_-5px_rgba(var(--primary-rgb),0.5)] transition-all active:translate-y-[0px]"
                    >
                        <Plus className="w-5 h-5 mr-2" />
                        Create Advertisement
                    </Link>
                </div>

                {/* Filters */}
                <MarketplaceFilters
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                    filters={filters}
                    setFilters={setFilters}
                    showFilters={showFilters}
                    setShowFilters={setShowFilters}
                />

                {/* Ads Grid */}
                {loading && ads.length === 0 ? (
                    <div className="flex flex-col items-center justify-center h-80 gap-4">
                        <Loader2 className="w-10 h-10 animate-spin text-primary" />
                        <span className="text-white font-bold uppercase tracking-widest text-sm animate-pulse">Loading Assets...</span>
                    </div>
                ) : filteredAds.length === 0 ? (
                    <div className="text-center py-20 bg-surface/50 border border-white/5 rounded-3xl">
                        <p className="text-text-dim font-bold uppercase tracking-widest">No matching advertisements found</p>
                        <button
                            onClick={() => { setSearchTerm(""); setFilters({ type: "all", fiat_currency: "all", payment_method: "all", sort_by: "best_price" }) }}
                            className="mt-4 text-primary hover:underline font-black uppercase text-xs"
                        >Clear all filters</button>
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {filteredAds.map((ad) => (
                            <AdCard
                                key={ad.ad_id}
                                ad={ad}
                                onTrade={handleInitiateTrade}
                            />
                        ))}
                    </div>
                )}

                {/* Pagination Status / Sentinel */}
                <div ref={loadMoreRef} className="mt-12 py-8 flex flex-col items-center justify-center border-t border-white/5">
                    {loadingMore ? (
                        <div className="flex items-center gap-3 text-text-dim">
                            <Loader2 className="w-5 h-5 animate-spin" />
                            <span className="text-sm font-bold uppercase tracking-widest">Scanning more ads...</span>
                        </div>
                    ) : hasMore ? (
                        <div className="h-4 w-4 rounded-full bg-primary/20 animate-ping"></div>
                    ) : ads.length > 0 && (
                        <div className="text-center">
                            <p className="text-white/40 text-xs font-black uppercase tracking-[0.3em]">
                                End of List — {total} Ads available
                            </p>
                        </div>
                    )}
                </div>

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
                                        <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] mb-3 block">Amount to {selectedAd.type === "sell" ? "Pay" : "Receive"} ({selectedAd.fiat_symbol})</label>
                                        <div className="relative">
                                            <input
                                                type="number"
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
                                            setIsInitiating(true);
                                            setInitiationError(null);
                                            try {
                                                const amount = parseFloat(tradeAmount);
                                                if (isNaN(amount) || amount < selectedAd.min_amount) {
                                                    throw new Error(`Minimum amount is ${selectedAd.min_amount} ${selectedAd.fiat_symbol}`);
                                                }
                                                if (amount > selectedAd.max_amount) {
                                                    throw new Error(`Maximum amount is ${selectedAd.max_amount} ${selectedAd.fiat_symbol}`);
                                                }
                                                const pmId = selectedAd.payment_method_ids ? selectedAd.payment_method_ids[0] : 0;
                                                const result = await TradeService.initiateTrade(selectedAd.ad_id, amount, pmId);
                                                router.push(`/user/dashboard/trades/${result.trade_id}`);
                                            } catch (err: any) {
                                                setInitiationError(err.message);
                                            } finally {
                                                setIsInitiating(false);
                                            }
                                        }}
                                        className="flex-1 py-6 bg-primary text-white rounded-[2rem] font-black uppercase tracking-widest text-xs hover:scale-105 active:scale-95 transition-all shadow-xl shadow-primary/20 disabled:opacity-50"
                                    >
                                        {isInitiating ? "Initiating..." : `Start ${selectedAd.type === "sell" ? "Buy" : "Sell"} Trade`}
                                    </button>
                                </div>
                            </motion.div>
                        </div>
                    )}
                </AnimatePresence>
            </div>

            <Footer />
        </main>
    );
};

export default MarketplacePage;
