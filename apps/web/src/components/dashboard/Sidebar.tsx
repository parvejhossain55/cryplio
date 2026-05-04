"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import {
    Wallet,
    LayoutDashboard,
    ArrowLeftRight,
    History,
    Settings,
    Shield,
    LogOut,
    ChevronRight,
    TrendingUp,
    CreditCard,
    Users,
    Store,
    BarChart3,
    UserCheck,
    X,
    Coins,
    DollarSign
} from "lucide-react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

const DISMISSED_STATUS_CARDS_KEY = "cryplio_dismissed_status_cards";

interface SidebarItem {
    name: string;
    href: string;
    icon: React.ElementType;
}

interface SidebarProps {
    role: "user" | "merchant" | "admin";
    isMobile?: boolean;
}

const Sidebar = ({ role, isMobile }: SidebarProps) => {
    const pathname = usePathname();
    const { user, logout, isLoading } = useAuth();
    const [dismissedStatusCards, setDismissedStatusCards] = useState<string[]>(() => {
        if (typeof window === "undefined") {
            return [];
        }

        try {
            return JSON.parse(localStorage.getItem(DISMISSED_STATUS_CARDS_KEY) || "[]");
        } catch {
            return [];
        }
    });

    const showStatusCard = Boolean(user && !dismissedStatusCards.includes(user.id) && !user.emailVerified);

    const dismissStatusCard = () => {
        if (!user) {
            return;
        }

        const nextDismissedCards = Array.from(new Set([...dismissedStatusCards, user.id]));
        setDismissedStatusCards(nextDismissedCards);
        localStorage.setItem(DISMISSED_STATUS_CARDS_KEY, JSON.stringify(nextDismissedCards));
    };

    const navigation: Record<string, SidebarItem[]> = {
        user: [
            { name: "Overview", href: "/user/dashboard", icon: LayoutDashboard },
            { name: "Marketplace", href: "/marketplace", icon: Store },
            { name: "My Trades", href: "/user/dashboard/trades", icon: History },
            { name: "Wallet", href: "/user/dashboard/wallet", icon: CreditCard },
            { name: "Settings", href: "/user/dashboard/settings", icon: Settings },
        ],
        merchant: [
            { name: "Dashboard", href: "/merchant/dashboard", icon: BarChart3 },
            { name: "My Ads", href: "/merchant/dashboard/ads", icon: Store },
            { name: "Client Orders", href: "/merchant/dashboard/orders", icon: ArrowLeftRight },
            { name: "Earnings", href: "/merchant/dashboard/earnings", icon: TrendingUp },
            { name: "Settings", href: "/merchant/dashboard/settings", icon: Settings },
        ],
        admin: [
            { name: "Admin Panel", href: "/admin/dashboard", icon: Shield },
            { name: "User Management", href: "/admin/dashboard/users", icon: Users },
            { name: "Merchants", href: "/admin/dashboard/merchants", icon: Store },
            { name: "KYC Reviews", href: "/admin/dashboard/kyc", icon: UserCheck },
            { name: "System Stats", href: "/admin/dashboard/stats", icon: BarChart3 },
            { name: "Payment Methods", href: "/admin/dashboard/payment-methods", icon: CreditCard },
            { name: "Crypto Assets", href: "/admin/dashboard/crypto-assets", icon: Coins },
            { name: "Fiat Currencies", href: "/admin/dashboard/fiat-currencies", icon: DollarSign },
        ],
    };

    const navItems = navigation[role] || [];

    return (
        <aside className={cn(
            "w-72 flex flex-col bg-surface overflow-y-auto scrollbar-hide",
            isMobile
                ? "h-full border-0"
                : "hidden md:flex border-r border-border h-screen fixed top-0 left-0 z-20"
        )}>
            <div className="p-8">
                <Link href="/" className="flex items-center space-x-3 group">
                    <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-xl shadow-primary/20 group-hover:rotate-12 transition-transform duration-500">
                        <Wallet className="text-background w-6 h-6" />
                    </div>
                    <span className="text-2xl font-black italic uppercase tracking-tighter text-white">
                        CRYP<span className="text-primary truncate">LIO</span>
                    </span>
                </Link>
            </div>

            <nav className="flex-1 px-4 space-y-2">
                <div className="px-4 py-2 mb-2">
                    <span className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">
                        Menu
                    </span>
                </div>
                {navItems.map((item) => {
                    const isActive = pathname === item.href;
                    return (
                        <Link
                            key={item.name}
                            href={item.href}
                            className={cn(
                                "flex items-center group px-4 py-3 rounded-xl transition-all duration-300",
                                isActive
                                    ? "bg-primary/10 text-primary border border-primary/20 shadow-lg shadow-primary/5"
                                    : "text-text-dim hover:text-white hover:bg-white/5 border border-transparent"
                            )}
                        >
                            <item.icon className={cn(
                                "w-5 h-5 mr-3 transition-colors",
                                isActive ? "text-primary" : "text-text-dim group-hover:text-white"
                            )} />
                            <span className="font-bold text-sm tracking-tight">{item.name}</span>
                            {isActive && <ChevronRight className="ml-auto w-4 h-4 text-primary" />}
                        </Link>
                    );
                })}
            </nav>

            <div className="p-4 mt-auto">
                {showStatusCard && user && (
                    <div className="relative p-4 rounded-2xl bg-surface-light border border-white/5 mb-4">
                        <button
                            type="button"
                            onClick={dismissStatusCard}
                            className="absolute right-3 top-3 p-1.5 rounded-lg text-text-dim hover:text-white hover:bg-white/5 transition-colors"
                            aria-label="Dismiss account status"
                        >
                            <X className="w-3.5 h-3.5" />
                        </button>
                        <div className="flex items-center space-x-3 mb-3">
                            <div className={cn(
                                "w-8 h-8 rounded-lg flex items-center justify-center",
                                user.emailVerified ? "bg-accent/20" : "bg-yellow-500/15"
                            )}>
                                <Shield className={cn(
                                    "w-4 h-4",
                                    user.emailVerified ? "text-accent" : "text-yellow-500"
                                )} />
                            </div>
                            <span className="text-xs font-bold text-white tracking-tight">
                                {user.emailVerified ? "Email Verified" : "Email Not Verified"}
                            </span>
                        </div>
                        <p className="text-[10px] text-text-dim font-medium leading-relaxed mb-3">
                            {user.emailVerified
                                ? "Your email is verified and account is secured."
                                : "Please verify your email to secure your account."}
                        </p>
                        <div className="w-full bg-white/5 h-1 rounded-full overflow-hidden">
                            <div className={cn(
                                "h-full rounded-full",
                                user.emailVerified ? "w-full bg-accent" : "w-1/3 bg-yellow-500"
                            )} />
                        </div>
                    </div>
                )}

                <button
                    type="button"
                    onClick={() => void logout()}
                    disabled={isLoading}
                    className="flex items-center w-full px-4 py-3 text-text-dim hover:text-white hover:bg-white/5 rounded-xl transition-all group border border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
                >
                    <LogOut className="w-5 h-5 mr-3 group-hover:text-white transition-colors" />
                    <span className="font-bold text-sm tracking-tight text-white/60 group-hover:text-white">
                        {isLoading ? "Logging out..." : "Logout"}
                    </span>
                </button>
            </div>
        </aside>
    );
};

export default Sidebar;
