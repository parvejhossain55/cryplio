"use client";
import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Store, Plus, Search, Filter } from "lucide-react";

const MerchantAds = () => {
    return (
        <DashboardLayout title="Ad Management" role="merchant">
            <div className="space-y-8">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div className="flex items-center space-x-4">
                        <div className="relative">
                            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-4 h-4" />
                            <input type="text" placeholder="Search Ads..." className="bg-surface border border-white/5 py-3 pl-12 pr-6 rounded-2xl text-xs w-64 outline-none focus:border-primary/50" />
                        </div>
                    </div>
                    <button className="bg-primary text-white px-8 py-3 rounded-2xl font-black uppercase text-[10px] tracking-widest flex items-center shadow-xl shadow-primary/20">
                        <Plus className="w-4 h-4 mr-2" /> Create New Ad
                    </button>
                </div>

                <div className="bg-surface border border-white/10 rounded-[2.5rem] p-20 flex flex-col items-center justify-center text-center">
                    <div className="w-20 h-20 bg-white/5 rounded-[2rem] flex items-center justify-center mb-6">
                        <Store className="w-10 h-10 text-text-dim" />
                    </div>
                    <h3 className="text-xl font-black text-white uppercase tracking-tight">No Active Ads</h3>
                    <p className="text-xs text-text-dim max-w-xs mt-2">You haven't listed any buying or selling ads yet. Start by creating your first offer.</p>
                </div>
            </div>
        </DashboardLayout>
    );
};
export default MerchantAds;
