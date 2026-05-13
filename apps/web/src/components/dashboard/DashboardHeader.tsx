"use client";

import React from "react";
import {
    Search,
    Bell,
    Menu,
    ChevronDown,
    User as UserIcon,
    Activity
} from "lucide-react";
import { motion } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { useNotifications } from "@/context/NotificationContext";
import Link from "next/link";

interface DashboardHeaderProps {
    title: string;
    onMenuClick?: () => void;
}

const DashboardHeader = ({ title, onMenuClick }: DashboardHeaderProps) => {
    const { user } = useAuth();
    const { unreadCount, isConnected } = useNotifications();

    return (
        <header className="sticky top-0 z-30 flex items-center justify-between px-6 md:px-10 py-4 bg-background/80 backdrop-blur-xl border-b border-white/5">
            <div className="flex items-center space-x-4">
                <button
                    onClick={onMenuClick}
                    className="p-2 md:hidden hover:bg-surface-light rounded-lg transition-colors border border-white/5"
                >
                    <Menu className="w-5 h-5 text-white" />
                </button>
                <div className="hidden md:block">
                    <h1 className="text-xl font-black text-white tracking-tight">{title}</h1>
                </div>
            </div>

            <div className="flex items-center space-x-3 md:space-x-6">
                {/* Search Bar - Desktop */}
                <div className="hidden lg:flex items-center relative group">
                    <Search className="absolute left-4 w-4 h-4 text-text-dim group-focus-within:text-primary transition-colors" />
                    <input
                        type="text"
                        placeholder="Search transactions, assets..."
                        className="bg-surface border border-border py-2.5 pl-12 pr-6 rounded-xl text-sm w-80 outline-none focus:border-primary/50 transition-all font-medium placeholder:text-text-dim/40"
                    />
                </div>

                {/* Global Stats */}
                <div className="hidden sm:flex items-center space-x-4 px-4 py-2 bg-surface-light rounded-xl border border-white/5">
                    <div className="flex items-center space-x-2">
                        <Activity className={`w-3.5 h-3.5 ${isConnected ? 'text-accent' : 'text-primary animate-pulse'}`} />
                        <span className="text-[10px] font-black text-white tracking-widest uppercase">
                            {isConnected ? 'Escrow Live' : 'Connecting...'}
                        </span>
                    </div>
                    <div className="w-[1px] h-3 bg-white/10" />
                    <div className="flex items-center space-x-1">
                        <span className="text-xs font-bold text-white">$1.2M</span>
                        <span className="text-[8px] text-accent font-black">+2.4%</span>
                    </div>
                </div>

                {/* Profile Dropdown */}
                <div className="flex items-center space-x-2 pl-2">
                    <div className="relative cursor-pointer group">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary/20 to-secondary/20 border border-white/10 flex items-center justify-center overflow-hidden">
                            {user?.avatarUrl ? (
                                <img src={user.avatarUrl} alt={user.name} className="w-full h-full object-cover" />
                            ) : (
                                <UserIcon className="w-5 h-5 text-text-dim group-hover:text-white transition-colors" />
                            )}
                        </div>
                        <div className={`absolute -bottom-1 -right-1 w-4 h-4 ${user?.isOnline ? 'bg-accent' : 'bg-text-dim'} border-4 border-background rounded-full`} />
                    </div>
                    <div className="hidden lg:block">
                        <div className="flex items-center space-x-1 cursor-pointer">
                            <span className="text-sm font-bold text-white">{user?.name || "Guest"}</span>
                        </div>
                    </div>
                </div>
            </div>
        </header>
    );
};

export default DashboardHeader;
