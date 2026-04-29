"use client";
import React from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { Users, Search, Filter } from "lucide-react";

const AdminUsers = () => {
    return (
        <DashboardLayout title="User Logistics" role="admin">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                <div className="flex items-center justify-between mb-8">
                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Active Accounts</h3>
                    <div className="flex space-x-2">
                        <div className="relative">
                            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-4 h-4" />
                            <input type="text" placeholder="Search Identity..." className="bg-white/5 border border-white/10 py-3 pl-12 pr-6 rounded-xl text-xs outline-none" />
                        </div>
                    </div>
                </div>
                <div className="border-t border-white/5 pt-10 text-center">
                    <Users className="w-12 h-12 text-text-dim/20 mx-auto mb-4" />
                    <p className="text-sm font-bold text-text-dim">User database synchronizing with kernel...</p>
                </div>
            </div>
        </DashboardLayout>
    );
};
export default AdminUsers;
