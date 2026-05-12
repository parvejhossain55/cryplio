"use client";

import React from "react";
import { ShieldCheck } from "lucide-react";

const SecurityStatus = () => {
    return (
        <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
            <div className="flex items-center space-x-4 mb-6">
                <div className="w-12 h-12 rounded-2xl bg-accent/20 flex items-center justify-center">
                    <ShieldCheck className="w-6 h-6 text-accent" />
                </div>
                <div>
                    <h4 className="font-black text-white uppercase text-xs tracking-widest">Account Health</h4>
                    <p className="text-[10px] font-bold text-accent uppercase tracking-widest">Excellent</p>
                </div>
            </div>
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">Account Security</span>
                    <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Verified</span>
                </div>
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">2FA Status</span>
                    <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Enabled</span>
                </div>
                <div className="flex items-center justify-between">
                    <span className="text-xs text-text-dim font-medium">Login Notifications</span>
                    <span className="text-[10px] bg-accent/10 text-accent px-2 py-0.5 rounded font-black uppercase tracking-widest">Active</span>
                </div>
            </div>
            <button className="w-full mt-8 py-4 bg-white/5 border border-white/5 rounded-2xl text-[10px] font-black text-white uppercase tracking-[0.2em] hover:bg-white/10 transition-all">
                Security Hub
            </button>
        </div>
    );
};

export default SecurityStatus;
