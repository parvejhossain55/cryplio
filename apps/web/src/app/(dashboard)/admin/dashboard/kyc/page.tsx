"use client";
import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { UserCheck, ShieldCheck, Search, Filter } from "lucide-react";

const AdminKYC = () => {
    return (
        <DashboardLayout title="Compliance Queue" role="admin">
            <div className="space-y-8">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="bg-surface border border-white/10 p-6 rounded-3xl">
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-2">Pending Reviews</p>
                        <h3 className="text-2xl font-black text-white">42</h3>
                    </div>
                    <div className="bg-surface border border-white/10 p-6 rounded-3xl">
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-2">Approved Today</p>
                        <h3 className="text-2xl font-black text-accent">128</h3>
                    </div>
                    <div className="bg-surface border border-white/10 p-6 rounded-3xl">
                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest mb-2">Flagged</p>
                        <h3 className="text-2xl font-black text-primary">3</h3>
                    </div>
                </div>

                <div className="bg-surface border border-white/10 rounded-[2.5rem] p-20 flex flex-col items-center justify-center text-center">
                    <div className="w-20 h-20 bg-accent/10 rounded-[2rem] flex items-center justify-center mb-6">
                        <UserCheck className="w-10 h-10 text-accent" />
                    </div>
                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Queue Clear</h3>
                    <p className="text-xs text-text-dim max-w-xs mt-2">All high-priority KYC applications have been processed for this region.</p>
                </div>
            </div>
        </DashboardLayout>
    );
};
export default AdminKYC;
