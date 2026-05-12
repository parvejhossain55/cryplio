"use client";

import React from "react";
import { ChevronRight, ArrowDownLeft, ArrowUpRight, ArrowLeftRight, ExternalLink } from "lucide-react";

interface Activity {
    type: string;
    asset: string;
    amount: string;
    status: string;
    date: string;
    price: string;
}

const RecentActivity = () => {
    const activities: Activity[] = [
        { type: "buy", asset: "USDT", amount: "500", status: "completed", date: "2 mins ago", price: "$1.00" },
        { type: "sell", asset: "BTC", amount: "0.02", status: "pending", date: "15 mins ago", price: "$64,200" },
        { type: "swap", asset: "ETH to USDT", amount: "1.2", status: "completed", date: "2 hours ago", price: "$3,200" },
        { type: "deposit", asset: "USD", amount: "1,200", status: "failed", date: "1 day ago", price: "-" },
    ];

    return (
        <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
            <div className="flex items-center justify-between mb-8">
                <h3 className="text-xl font-black text-white">Recent Activity</h3>
                <button className="text-xs font-bold text-primary hover:text-white transition-colors flex items-center">
                    View All <ChevronRight className="w-4 h-4 ml-1" />
                </button>
            </div>

            <div className="space-y-4">
                {activities.map((tx, i) => (
                    <div key={i} className="flex items-center justify-between p-4 rounded-2xl hover:bg-white/5 transition-all group">
                        <div className="flex items-center space-x-4">
                            <div className={`w-12 h-12 rounded-2xl flex items-center justify-center border border-white/5 ${tx.type === 'buy' ? 'bg-accent/10 text-accent' :
                                tx.type === 'sell' ? 'bg-primary/10 text-primary' :
                                    'bg-secondary/10 text-secondary'
                                }`}>
                                {tx.type === 'buy' ? <ArrowDownLeft className="w-6 h-6" /> :
                                    tx.type === 'sell' ? <ArrowUpRight className="w-6 h-6" /> :
                                        <ArrowLeftRight className="w-5 h-5" />}
                            </div>
                            <div>
                                <h4 className="font-bold text-white tracking-tight flex items-center uppercase">
                                    {tx.type} {tx.asset}
                                    <span className={`ml-2 text-[8px] px-1.5 py-0.5 rounded font-black uppercase tracking-widest ${tx.status === 'completed' ? 'bg-accent/10 text-accent' :
                                        tx.status === 'pending' ? 'bg-primary/10 text-primary' :
                                            'bg-red-500/10 text-red-500'
                                        }`}>
                                        {tx.status}
                                    </span>
                                </h4>
                                <p className="text-[10px] font-medium text-text-dim mt-1">{tx.date} • Price: {tx.price}</p>
                            </div>
                        </div>
                        <div className="text-right">
                            <p className="font-black text-white">{tx.amount} {tx.asset.split(' ')[0]}</p>
                            <ExternalLink className="w-3 h-3 ml-auto mt-1 text-text-dim opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer hover:text-white" />
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default RecentActivity;
