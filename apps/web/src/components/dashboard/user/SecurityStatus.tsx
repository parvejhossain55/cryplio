"use client";

import React, { useEffect, useState } from "react";
import { ShieldCheck } from "lucide-react";
import { fetchHeaderProfile } from "@/services/apiClient";
import { HeaderProfileResponse } from "@/types/api";

const SecurityStatus = () => {
    const [data, setData] = useState<HeaderProfileResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetchHeaderProfile();
                setData(response);
            } catch (error) {
                console.error("Failed to fetch header profile:", error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const getHealthColor = (health: string) => {
        switch (health) {
            case "EXCELLENT":
                return "text-accent";
            case "GOOD":
                return "text-green-400";
            case "FAIR":
                return "text-yellow-400";
            case "POOR":
                return "text-red-400";
            default:
                return "text-text-dim";
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "VERIFIED":
            case "ENABLED":
            case "ACTIVE":
                return "bg-accent/10 text-accent";
            case "UNVERIFIED":
            case "DISABLED":
            case "INACTIVE":
                return "bg-red-500/10 text-red-400";
            default:
                return "bg-white/5 text-text-dim";
        }
    };

    if (isLoading) {
        return (
            <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                <div className="animate-pulse">
                    <div className="flex items-center space-x-4 mb-6">
                        <div className="w-12 h-12 rounded-2xl bg-white/5" />
                        <div className="space-y-2">
                            <div className="h-3 w-24 bg-white/5 rounded" />
                            <div className="h-2 w-16 bg-white/5 rounded" />
                        </div>
                    </div>
                    <div className="space-y-4">
                        <div className="h-8 bg-white/5 rounded" />
                        <div className="h-8 bg-white/5 rounded" />
                        <div className="h-8 bg-white/5 rounded" />
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
            <div className="flex items-center space-x-4 mb-6">
                <div className="w-12 h-12 rounded-2xl bg-accent/20 flex items-center justify-center">
                    <ShieldCheck className="w-6 h-6 text-accent" />
                </div>
                <div>
                    <h4 className="font-black text-white uppercase text-xs tracking-widest">Account Health</h4>
                    <p className={`text-[10px] font-bold uppercase tracking-widest ${getHealthColor(data?.account_health || "")}`}>
                        {data?.account_health || "Loading..."}
                    </p>
                </div>
            </div>
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">Account Security</span>
                    <span className={`text-[10px] px-2 py-0.5 rounded font-black uppercase tracking-widest ${getStatusColor(data?.account_security || "")}`}>
                        {data?.account_security || "Loading..."}
                    </span>
                </div>
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">2FA Status</span>
                    <span className={`text-[10px] px-2 py-0.5 rounded font-black uppercase tracking-widest ${getStatusColor(data?.two_factor_status || "")}`}>
                        {data?.two_factor_status || "Loading..."}
                    </span>
                </div>
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">Login Notifications</span>
                    <span className={`text-[10px] px-2 py-0.5 rounded font-black uppercase tracking-widest ${getStatusColor(data?.login_notifications || "")}`}>
                        {data?.login_notifications || "Loading..."}
                    </span>
                </div>
            </div>
            <button className="w-full mt-8 py-4 bg-white/5 border border-white/5 rounded-2xl text-[10px] font-black text-white uppercase tracking-[0.2em] hover:bg-white/10 transition-all">
                Security Hub
            </button>
        </div>
    );
};

export default SecurityStatus;
