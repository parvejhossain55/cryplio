"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Plus,
    Search,
    Filter,
    MoreHorizontal,
    Pause,
    Play,
    Edit,
    ChevronRight,
    TrendingUp,
    Users,
    Activity,
    AlertCircle,
    Loader2,
    LayoutGrid,
    List as ListIcon
} from "lucide-react";
import Link from "next/link";
import { authService } from "@/services/authService";
import DashboardLayout from "@/components/dashboard/DashboardLayout";

const MerchantAdsPage = () => {
    const [ads, setAds] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [viewMode, setViewMode] = useState<"grid" | "list">("list");

    useEffect(() => {
        fetchMyAds();
    }, []);

    const fetchMyAds = async () => {
        setIsLoading(true);
        try {
            const data = await authService.getMyAds();
            setAds(data.ads || []);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleToggleStatus = async (adId: string) => {
        try {
            await authService.toggleAdStatus(adId);
            setAds(ads.map(ad =>
                ad.ad_id === adId ? { ...ad, status: ad.status === "active" ? "paused" : "active" } : ad
            ));
        } catch (err: any) {
            alert(err.message);
        }
    };

    return (
        <DashboardLayout title="My Advertisements" role="merchant">
            <div className="space-y-8 pb-10">
                {/* Header section with Action */}
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
                    <div>
                        <h1 className="text-3xl font-black text-white italic uppercase tracking-tight">Active <span className="gradient-text">Inventory</span></h1>
                        <p className="text-text-dim text-xs mt-1 font-medium italic">Monitor and control your marketplace visibility</p>
                    </div>
                    <Link
                        href="/marketplace/create"
                        className="flex items-center justify-center px-6 py-3 bg-white text-background rounded-2xl font-black uppercase tracking-widest text-[10px] hover:scale-105 active:scale-95 transition-all shadow-xl shadow-white/5"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Create New Ad
                    </Link>
                </div>

                {/* Quick Stats */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    {[
                        { label: "Active", value: ads.filter(a => a.status !== 'paused').length, icon: Activity, color: "text-accent" },
                        { label: "Paused", value: ads.filter(a => a.status === 'paused').length, icon: Pause, color: "text-yellow-500" },
                        { label: "Views", value: "1.2K", icon: Users, color: "text-primary" },
                        { label: "Matches", value: "94%", icon: TrendingUp, color: "text-emerald-500" },
                    ].map((stat, i) => (
                        <div key={i} className="bg-surface border border-white/5 p-5 rounded-[1.5rem] space-y-2 group hover:border-white/10 transition-colors">
                            <div className={`p-2 rounded-xl bg-white/5 w-fit ${stat.color} group-hover:scale-110 transition-transform`}>
                                <stat.icon className="w-4 h-4" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim/50">{stat.label}</p>
                                <p className="text-xl font-black text-white">{stat.value}</p>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Filter Bar */}
                <div className="flex flex-wrap items-center justify-between gap-4 bg-surface/50 border border-white/5 p-3 rounded-[1.2rem]">
                    <div className="flex items-center gap-3">
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-dim" />
                            <input
                                type="text"
                                placeholder="Search by ID or Asset..."
                                className="bg-background/50 border border-white/5 rounded-xl py-2 pl-10 pr-4 text-[10px] font-bold text-white outline-none focus:border-primary/30 transition-all w-48 md:w-64"
                            />
                        </div>
                        <button className="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/5 text-text-dim hover:text-white transition-all text-[10px] font-black uppercase tracking-widest border border-transparent">
                            <Filter className="w-3.5 h-3.5" />
                            Filters
                        </button>
                    </div>
                </div>

                {/* Ads Content */}
                <AnimatePresence mode="wait">
                    {isLoading ? (
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            className="flex flex-col items-center justify-center py-24 space-y-4"
                        >
                            <Loader2 className="w-10 h-10 text-primary animate-spin" />
                            <p className="text-text-dim font-black uppercase tracking-widest text-[10px]">Accessing Vault...</p>
                        </motion.div>
                    ) : ads.length === 0 ? (
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            className="bg-surface border border-white/5 rounded-[2.5rem] p-16 text-center space-y-6"
                        >
                            <div className="w-16 h-16 bg-white/5 rounded-[1.5rem] flex items-center justify-center mx-auto">
                                <AlertCircle className="w-8 h-8 text-text-dim/30" />
                            </div>
                            <div>
                                <h3 className="text-xl font-black text-white italic">EMPTY INVENTORY</h3>
                                <p className="text-text-dim text-xs max-w-xs mx-auto mt-2">Scale your trading operations by listing your first peer-to-peer advertisement.</p>
                            </div>
                            <Link
                                href="/marketplace/create"
                                className="inline-flex items-center px-6 py-3 bg-white text-background rounded-xl font-black uppercase tracking-widest text-[10px] hover:scale-105 active:scale-95 transition-all"
                            >
                                Get Started
                            </Link>
                        </motion.div>
                    ) : (
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            className={viewMode === "grid" ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5" : "space-y-3"}
                        >
                            {ads.map((ad, i) => (
                                <motion.div
                                    key={ad.ad_id}
                                    initial={{ opacity: 0, y: 10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    transition={{ delay: i * 0.05 }}
                                    className={`bg-surface border border-white/5 rounded-[1.5rem] overflow-hidden group hover:border-white/10 transition-all ${ad.status === 'paused' ? 'opacity-60 saturate-50' : ''}`}
                                >
                                    <div className="p-6">
                                        <div className="flex items-center justify-between mb-5">
                                            <div className="flex items-center gap-2">
                                                <div className={`px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest ${ad.type === 'buy' ? 'bg-accent/10 text-accent' : 'bg-primary/10 text-primary'}`}>
                                                    {ad.type}
                                                </div>
                                                <div className="text-[8px] font-black text-text-dim uppercase tracking-widest">
                                                    #{ad.ad_id.slice(0, 8)}
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-1.5">
                                                <button
                                                    onClick={() => handleToggleStatus(ad.ad_id)}
                                                    className="p-1.5 rounded-lg bg-white/5 text-text-dim hover:text-white transition-all"
                                                >
                                                    {ad.status === 'paused' ? <Play className="w-3.5 h-3.5" /> : <Pause className="w-3.5 h-3.5" />}
                                                </button>
                                                <button className="p-1.5 rounded-lg bg-white/5 text-text-dim hover:text-white transition-all">
                                                    <Edit className="w-3.5 h-3.5" />
                                                </button>
                                            </div>
                                        </div>

                                        <div className="space-y-5">
                                            <div className="flex items-end justify-between">
                                                <div>
                                                    <p className="text-[8px] font-black text-text-dim uppercase tracking-widest mb-1">Exchange Rate</p>
                                                    <p className="text-xl font-black text-white italic">
                                                        {ad.price.toLocaleString()} <span className="text-[10px] font-medium not-italic text-text-dim/40">USD</span>
                                                    </p>
                                                </div>
                                                <div className="text-right">
                                                    <p className="text-[8px] font-black text-text-dim uppercase tracking-widest mb-1">Limit Range</p>
                                                    <p className="text-[10px] font-bold text-white">${ad.min_amount} - ${ad.max_amount}</p>
                                                </div>
                                            </div>

                                            <div className="border-t border-white/5 pt-4 flex items-center justify-between">
                                                <div className="flex items-center gap-2">
                                                    <div className="w-6 h-6 rounded bg-primary/10 flex items-center justify-center text-[8px] font-black text-primary">U</div>
                                                    <span className="text-[8px] font-black text-text-dim uppercase tracking-widest">USDT / USD</span>
                                                </div>
                                                <Link
                                                    href={`/marketplace`}
                                                    className="flex items-center gap-1 text-[8px] font-black text-primary uppercase tracking-widest hover:translate-x-1 transition-all"
                                                >
                                                    Live View <ChevronRight className="w-2.5 h-2.5" />
                                                </Link>
                                            </div>
                                        </div>
                                    </div>
                                </motion.div>
                            ))}
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>
        </DashboardLayout>
    );
};

export default MerchantAdsPage;
