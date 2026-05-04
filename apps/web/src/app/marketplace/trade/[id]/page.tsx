"use client";

import React, { useState, useEffect, useRef } from "react";
import { useParams, useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import {
    Shield,
    Clock,
    Send,
    Lock,
    CheckCircle2,
    AlertTriangle,
    Info,
    ChevronLeft,
    Loader2,
    Check,
    X,
    FileText,
    ExternalLink
} from "lucide-react";
import Link from "next/link";
import { authService } from "@/services/authService";
import Navbar from "@/components/layout/Navbar";
import { useAuth } from "@/context/AuthContext";

const TradeDetailPage = () => {
    const { id } = useParams();
    const router = useRouter();
    const { user } = useAuth();
    const [trade, setTrade] = useState<any>(null);
    const [messages, setMessages] = useState<any[]>([]);
    const [newMessage, setNewMessage] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [isUpdating, setIsUpdating] = useState(false);
    const chatEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (id) {
            fetchTradeDetails();
            const interval = setInterval(fetchTradeDetails, 5000); // Poll for updates
            return () => clearInterval(interval);
        }
    }, [id]);

    useEffect(() => {
        chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);

    const fetchTradeDetails = async () => {
        try {
            const [tradeData, msgData] = await Promise.all([
                authService.getTradeDetails(id as string),
                authService.getTradeMessages(id as string)
            ]);
            setTrade(tradeData);
            setMessages(msgData || []);
        } catch (err: any) {
            console.error(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSendMessage = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newMessage.trim()) return;

        try {
            const msg = await authService.sendTradeMessage(id as string, newMessage);
            setMessages([...messages, msg]);
            setNewMessage("");
        } catch (err: any) {
            alert(err.message);
        }
    };

    const handleAction = async (action: string) => {
        setIsUpdating(true);
        try {
            await authService.updateTradeStatus(id as string, action);
            fetchTradeDetails();
        } catch (err: any) {
            alert(err.message);
        } finally {
            setIsUpdating(false);
        }
    };

    if (isLoading || !trade) {
        return (
            <div className="min-h-screen bg-background flex flex-col items-center justify-center space-y-4">
                <Loader2 className="w-12 h-12 text-primary animate-spin" />
                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">Establishing Secure Uplink...</p>
            </div>
        );
    }

    const isBuyer = user?.id === trade.buyer_id;
    const isSeller = user?.id === trade.seller_id;

    return (
        <main className="min-h-screen bg-background text-white">
            <Navbar />

            <div className="container mx-auto px-4 md:px-6 pt-32 pb-20">
                {/* Back Link */}
                <button
                    onClick={() => router.back()}
                    className="flex items-center gap-2 text-[10px] font-black text-text-dim uppercase tracking-widest hover:text-white transition-all mb-8"
                >
                    <ChevronLeft className="w-4 h-4" /> Go Back
                </button>

                <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                    {/* Left Column: Trade Info */}
                    <div className="lg:col-span-4 space-y-6">
                        <section className="glass rounded-[2.5rem] border-white/5 p-8 relative overflow-hidden">
                            <div className="absolute top-0 right-0 w-32 h-32 bg-primary/10 blur-3xl -z-10" />

                            <div className="flex items-center justify-between mb-8">
                                <div className={`px-4 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest border ${trade.status === 'released' ? 'bg-accent/10 text-accent border-accent/20' :
                                        trade.status === 'paid' ? 'bg-blue-500/10 text-blue-500 border-blue-500/20' :
                                            'bg-primary/10 text-primary border-primary/20'
                                    }`}>
                                    {trade.status}
                                </div>
                                <div className="text-right">
                                    <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">Order ID</p>
                                    <p className="text-xs font-bold text-white">#{trade.trade_id.slice(0, 12)}</p>
                                </div>
                            </div>

                            <div className="space-y-6">
                                <div className="flex justify-between items-end">
                                    <div className="space-y-1">
                                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">Receiving Asset</p>
                                        <h2 className="text-3xl font-black italic">{trade.crypto_amount.toFixed(2)} USDT</h2>
                                    </div>
                                    <div className="text-right space-y-1">
                                        <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">Total Price</p>
                                        <h2 className="text-2xl font-black italic text-accent">{trade.fiat_amount.toLocaleString()} USD</h2>
                                    </div>
                                </div>

                                <div className="p-4 bg-white/5 rounded-2xl border border-white/5 space-y-3">
                                    <div className="flex justify-between text-xs">
                                        <span className="text-text-dim uppercase tracking-widest font-black text-[10px]">Exchange Rate</span>
                                        <span className="font-bold">1 USDT = {trade.exchange_rate} USD</span>
                                    </div>
                                    <div className="flex justify-between text-xs">
                                        <span className="text-text-dim uppercase tracking-widest font-black text-[10px]">Payment Method</span>
                                        <span className="font-bold">Bank Transfer</span>
                                    </div>
                                    <div className="flex justify-between text-xs">
                                        <span className="text-text-dim uppercase tracking-widest font-black text-[10px]">Trade Window</span>
                                        <span className="font-bold flex items-center gap-1">
                                            <Clock className="w-3.5 h-3.5" /> 15 Minutes
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </section>

                        <section className="glass rounded-[2.5rem] border-white/5 p-8 space-y-6 sticky top-24">
                            <div className="flex items-center gap-3">
                                <Shield className="w-6 h-6 text-primary" />
                                <h3 className="text-lg font-black italic uppercase tracking-tight">ESCROW PROTECTION</h3>
                            </div>
                            <p className="text-[10px] font-bold text-text-dim leading-relaxed uppercase tracking-widest">
                                Assets are securely held in our institutional-grade vault. Do not release until you have confirmed receipt of funds in your account.
                            </p>

                            <div className="space-y-3">
                                {isBuyer && trade.status === 'pending' && (
                                    <button
                                        disabled={isUpdating}
                                        onClick={() => handleAction('pay')}
                                        className="w-full py-4 bg-white text-background rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-[1.02] active:scale-95 transition-all shadow-xl shadow-white/5"
                                    >
                                        {isUpdating ? <Loader2 className="w-4 h-4 animate-spin mx-auto" /> : "I HAVE PAID"}
                                    </button>
                                )}

                                {isSeller && trade.status === 'paid' && (
                                    <button
                                        disabled={isUpdating}
                                        onClick={() => handleAction('release')}
                                        className="w-full py-4 bg-accent text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-[1.02] active:scale-95 transition-all shadow-xl shadow-accent/20"
                                    >
                                        {isUpdating ? <Loader2 className="w-4 h-4 animate-spin mx-auto" /> : "RELEASE ASSETS"}
                                    </button>
                                )}

                                {trade.status !== 'completed' && trade.status !== 'cancelled' && (
                                    <button
                                        disabled={isUpdating}
                                        onClick={() => handleAction('cancel')}
                                        className="w-full py-4 bg-red-500/10 text-red-500 border border-red-500/20 rounded-2xl font-black uppercase tracking-widest text-xs hover:bg-red-500/20 transition-all"
                                    >
                                        CANCEL ORDER
                                    </button>
                                )}
                            </div>

                            <div className="pt-4 border-t border-white/5 flex items-center gap-2 text-[8px] font-black text-text-dim uppercase tracking-widest justify-center">
                                <Lock className="w-3 h-3" /> End-to-End Encrypted Communication
                            </div>
                        </section>
                    </div>

                    {/* Right Column: Chat */}
                    <div className="lg:col-span-8 flex flex-col h-[700px] lg:h-auto">
                        <div className="glass rounded-[3rem] border-white/5 flex flex-col flex-1 overflow-hidden">
                            {/* Chat Header */}
                            <div className="p-6 border-b border-white/5 flex items-center justify-between bg-white/5">
                                <div className="flex items-center gap-4">
                                    <div className="w-12 h-12 rounded-2xl bg-surface border border-white/10 flex items-center justify-center font-black text-primary">
                                        {isBuyer ? 'S' : 'B'}
                                    </div>
                                    <div>
                                        <p className="font-black italic uppercase tracking-tight">TRADING PARTNER</p>
                                        <div className="flex items-center gap-1.5">
                                            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                                            <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">ACTIVE NOW</p>
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center gap-3">
                                    <button className="p-3 rounded-2xl bg-white/5 text-text-dim hover:text-white transition-all">
                                        <Info className="w-5 h-5" />
                                    </button>
                                </div>
                            </div>

                            {/* Chat Messages */}
                            <div className="flex-1 overflow-y-auto p-8 space-y-6 scrollbar-hide">
                                {messages.map((msg, i) => {
                                    const isMe = msg.sender_id === user?.id;
                                    return (
                                        <motion.div
                                            key={msg.message_id || i}
                                            initial={{ opacity: 0, y: 10, scale: 0.95 }}
                                            animate={{ opacity: 1, y: 0, scale: 1 }}
                                            className={`flex ${isMe ? 'justify-end' : 'justify-start'}`}
                                        >
                                            <div className={`max-w-[80%] p-4 rounded-2xl text-sm font-bold leading-relaxed ${isMe ? 'bg-primary text-white ml-12 rounded-tr-none' : 'bg-white/5 border border-white/5 mr-12 rounded-tl-none'
                                                }`}>
                                                {msg.content}
                                                <div className={`mt-2 text-[8px] font-black tracking-widest uppercase flex items-center gap-2 ${isMe ? 'text-white/50 justify-end' : 'text-text-dim/50'
                                                    }`}>
                                                    {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                                    {isMe && <Check className="w-3 h-3" />}
                                                </div>
                                            </div>
                                        </motion.div>
                                    );
                                })}
                                <div ref={chatEndRef} />
                            </div>

                            {/* Chat Input */}
                            <form onSubmit={handleSendMessage} className="p-6 bg-surface-light/50 border-t border-white/5">
                                <div className="flex items-center gap-4">
                                    <div className="flex-1 relative">
                                        <input
                                            type="text"
                                            value={newMessage}
                                            onChange={(e) => setNewMessage(e.target.value)}
                                            placeholder="Write a message to your partner..."
                                            className="w-full bg-background/50 border border-white/5 rounded-2xl py-4 pl-6 pr-12 text-sm font-bold outline-none focus:border-primary transition-all shadow-inner"
                                        />
                                        <button
                                            type="button"
                                            className="absolute right-4 top-1/2 -translate-y-1/2 p-2 text-text-dim hover:text-white transition-all"
                                        >
                                            <AlertTriangle className="w-5 h-5" />
                                        </button>
                                    </div>
                                    <button
                                        type="submit"
                                        disabled={!newMessage.trim()}
                                        className="p-4 bg-primary text-white rounded-2xl shadow-xl shadow-primary/20 hover:scale-105 active:scale-95 transition-all text-xs font-black disabled:opacity-50 disabled:grayscale"
                                    >
                                        <Send className="w-5 h-5" />
                                    </button>
                                </div>
                                <p className="text-[8px] text-center text-text-dim font-black uppercase tracking-widest mt-4">
                                    Institutional Escrow Active • Do not share sensitive payment details outside this chat
                                </p>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    );
};

export default TradeDetailPage;
