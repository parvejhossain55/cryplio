"use client";

import React from "react";
import { Store, ArrowLeftRight, CreditCard, TrendingUp } from "lucide-react";

const QuickActions = () => {
    const actions = [
        { name: "P2P Marketplace", icon: Store, color: "text-primary", bg: "bg-primary/10" },
        { name: "Instant Swap", icon: ArrowLeftRight, color: "text-secondary", bg: "bg-secondary/10" },
        { name: "Card Services", icon: CreditCard, color: "text-accent", bg: "bg-accent/10" },
        { name: "Refer & Earn", icon: TrendingUp, color: "text-primary", bg: "bg-primary/10" },
    ];

    return (
        <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
            <h3 className="text-lg font-black text-white mb-6">Quick Actions</h3>
            <div className="grid grid-cols-2 gap-4">
                {actions.map((action, i) => (
                    <button key={i} className="flex flex-col items-center justify-center p-6 rounded-3xl bg-white/5 border border-white/5 hover:border-white/10 hover:bg-white/[0.08] transition-all group">
                        <div className={`w-12 h-12 rounded-2xl ${action.bg} ${action.color} flex items-center justify-center mb-3 group-hover:scale-110 transition-transform`}>
                            <action.icon className="w-6 h-6" />
                        </div>
                        <span className="text-[10px] font-black text-white text-center uppercase tracking-widest leading-tight">{action.name}</span>
                    </button>
                ))}
            </div>
        </div>
    );
};

export default QuickActions;
