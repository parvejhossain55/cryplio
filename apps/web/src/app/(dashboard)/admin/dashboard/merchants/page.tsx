"use client";
import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Store, ShieldCheck, Search } from "lucide-react";

const AdminMerchants = () => {
    return (
        <DashboardLayout title="Merchant Registrar" role="admin">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                <div className="flex items-center justify-between mb-8">
                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Verified Liquidators</h3>
                    <ShieldCheck className="w-6 h-6 text-accent" />
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="p-6 rounded-3xl bg-white/5 border border-white/5">
                            <div className="w-12 h-12 bg-primary/10 rounded-2xl mb-4 flex items-center justify-center">
                                <Store className="w-6 h-6 text-primary" />
                            </div>
                            <h4 className="font-black text-white uppercase text-sm">Merchant #{i * 124}</h4>
                            <p className="text-[10px] text-text-dim mt-1">Status: Active • Volume: $12k+/mo</p>
                        </div>
                    ))}
                </div>
            </div>
        </DashboardLayout>
    );
};
export default AdminMerchants;
