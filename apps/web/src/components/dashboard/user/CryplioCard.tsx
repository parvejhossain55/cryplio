"use client";

import React from "react";
import { Wallet } from "lucide-react";

const CryplioCard = () => {
    return (
        <div className="bg-gradient-to-br from-primary to-secondary p-8 rounded-[2.5rem] relative overflow-hidden group shadow-2xl shadow-primary/20 h-64">
            <div className="absolute top-0 right-0 p-8">
                <Wallet className="w-10 h-10 text-white/50 group-hover:rotate-12 transition-transform duration-500" />
            </div>
            <div className="absolute bottom-0 left-0 p-8 w-full">
                <div className="flex justify-between items-end">
                    <div>
                        <p className="text-[10px] font-black text-white/60 uppercase tracking-[0.2em] mb-1">Virtual Card</p>
                        <p className="text-2xl font-black text-white tracking-widest">•••• 4820</p>
                    </div>
                    <div className="text-right">
                        <p className="text-[10px] font-black text-white/60 uppercase tracking-[0.2em] mb-1">Exp</p>
                        <p className="text-lg font-black text-white italic">08/28</p>
                    </div>
                </div>
            </div>
            {/* Card visual elements */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[150%] h-[150%] border border-white/20 rounded-full opacity-20 group-hover:scale-110 transition-transform duration-700" />
        </div>
    );
};

export default CryplioCard;
