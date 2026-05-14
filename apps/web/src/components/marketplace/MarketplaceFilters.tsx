"use client";

import React from "react";
import { Search, Filter, ChevronDown } from "lucide-react";

interface MarketplaceFiltersProps {
    searchTerm: string;
    setSearchTerm: (val: string) => void;
    filters: {
        type: string;
        fiat_currency: string;
        payment_method: string;
        sort_by: string;
    };
    setFilters: (filters: any) => void;
    showFilters: boolean;
    setShowFilters: (show: boolean) => void;
}

const MarketplaceFilters: React.FC<MarketplaceFiltersProps> = ({
    searchTerm,
    setSearchTerm,
    filters,
    setFilters,
    showFilters,
    setShowFilters
}) => {
    return (
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
                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50 transition-colors"
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
                        <label className="block text-text-dim text-sm mb-2 font-medium">Trade Type</label>
                        <select
                            value={filters.type}
                            onChange={(e) => setFilters({ ...filters, type: e.target.value })}
                            className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50 transition-colors"
                        >
                            <option value="all">All Types</option>
                            <option value="buy">Buy USDT</option>
                            <option value="sell">Sell USDT</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-text-dim text-sm mb-2 font-medium">Currency</label>
                        <select
                            value={filters.fiat_currency}
                            onChange={(e) => setFilters({ ...filters, fiat_currency: e.target.value })}
                            className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50 transition-colors"
                        >
                            <option value="all">All Currencies</option>
                            <option value="USD">USD</option>
                            <option value="BDT">BDT</option>
                            <option value="PKR">PKR</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-text-dim text-sm mb-2 font-medium">Payment Method</label>
                        <select
                            value={filters.payment_method}
                            onChange={(e) => setFilters({ ...filters, payment_method: e.target.value })}
                            className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50 transition-colors"
                        >
                            <option value="all">All Methods</option>
                            <option value="bkash">Bkash</option>
                            <option value="nagad">Nagad</option>
                            <option value="bank">Bank Transfer</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-text-dim text-sm mb-2 font-medium">Sort By</label>
                        <select
                            value={filters.sort_by}
                            onChange={(e) => setFilters({ ...filters, sort_by: e.target.value })}
                            className="w-full px-4 py-2 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50 transition-colors"
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
    );
};

export default MarketplaceFilters;
