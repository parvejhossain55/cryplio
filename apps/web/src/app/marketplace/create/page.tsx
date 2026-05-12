"use client";

import React, { useState } from "react";
import { motion } from "framer-motion";
import {
    ChevronLeft,
    Info,
    Zap,
    ShieldCheck,
    ArrowRight,
    CheckCircle2,
    AlertCircle,
    Loader2,
    DollarSign,
    Plus,
    X
} from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import { tradeService } from "@/services/tradeService";

const CreateAdPage = () => {
    const router = useRouter();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    // Form State
    const [adType, setAdType] = useState<"buy" | "sell">("buy");
    const [cryptoId, setCryptoId] = useState(1); // USDT ERC20
    const [fiatId, setFiatId] = useState(1); // USD
    const [priceType, setPriceType] = useState<"fixed" | "floating">("fixed");
    const [price, setPrice] = useState("");
    const [minAmount, setMinAmount] = useState("");
    const [maxAmount, setMaxAmount] = useState("");
    const [paymentWindow, setPaymentWindow] = useState(15);
    const [tradeTerms, setTradeTerms] = useState("");
    const [selectedMethods, setSelectedMethods] = useState<number[]>([]);

    // Dropdown states
    const [isCryptoOpen, setIsCryptoOpen] = useState(false);
    const [isFiatOpen, setIsFiatOpen] = useState(false);
    const [isWindowOpen, setIsWindowOpen] = useState(false);

    // Close all dropdowns when clicking outside
    const closeDropdowns = () => {
        setIsCryptoOpen(false);
        setIsFiatOpen(false);
        setIsWindowOpen(false);
    };

    const assets = [
        { id: 1, symbol: "USDT", name: "Tether (ERC20)" },
        { id: 2, symbol: "USDT", name: "Tether (TRC20)" },
        { id: 4, symbol: "ETH", name: "Ethereum" },
        { id: 5, symbol: "BTC", name: "Bitcoin" },
    ];

    const fiats = [
        { id: 1, symbol: "USD", name: "US Dollar" },
        { id: 2, symbol: "BDT", name: "Bangladeshi Taka" },
        { id: 5, symbol: "NGN", name: "Nigerian Naira" },
        { id: 3, symbol: "EUR", name: "Euro" },
    ];

    const paymentMethods = [
        { id: 1, name: "bKash" },
        { id: 2, name: "Nagad" },
        { id: 3, name: "Bank Transfer" },
        { id: 4, name: "Wise" },
        { id: 5, name: "PayPal" },
    ];

    const toggleMethod = (id: number) => {
        if (selectedMethods.includes(id)) {
            setSelectedMethods(selectedMethods.filter(m => m !== id));
        } else {
            setSelectedMethods([...selectedMethods, id]);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        try {
            const adData = {
                type: adType,
                crypto_id: cryptoId,
                fiat_id: fiatId,
                price_type: priceType,
                price: parseFloat(price),
                min_amount: parseFloat(minAmount),
                max_amount: parseFloat(maxAmount),
                payment_methods: selectedMethods,
                trade_terms: tradeTerms,
                payment_window_minutes: paymentWindow,
            };

            await tradeService.createAd(adData);
            setSuccess(true);
            setTimeout(() => router.push("/marketplace"), 2000);
        } catch (err: any) {
            setError(err.message || "Failed to create advertisement");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (success) {
        return (
            <main className="min-h-screen bg-background flex flex-col justify-center items-center p-6">
                <motion.div
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    className="glass p-12 rounded-[3.5rem] border-accent/20 text-center space-y-6 max-w-md"
                >
                    <div className="w-24 h-24 bg-accent/20 rounded-full flex items-center justify-center mx-auto">
                        <CheckCircle2 className="w-12 h-12 text-accent" />
                    </div>
                    <h1 className="text-3xl font-black text-white italic">POSTED SUCCESSFULLY</h1>
                    <p className="text-text-dim font-medium">Your advertisement is now live on the marketplace. Redirecting you...</p>
                </motion.div>
            </main>
        );
    }

    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            <section className="pt-32 pb-24">
                <div className="container mx-auto px-4 max-w-4xl">
                    <Link href="/marketplace" className="inline-flex items-center text-text-dim hover:text-white transition-colors mb-8 group">
                        <ChevronLeft className="w-5 h-5 mr-1 group-hover:-translate-x-1 transition-transform" />
                        <span className="text-xs font-black uppercase tracking-widest">Back to Marketplace</span>
                    </Link>

                    <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
                        <div className="space-y-4">
                            <h1 className="text-5xl md:text-6xl font-black tracking-tight leading-none italic">
                                CREATE <br />
                                <span className="gradient-text uppercase">Advertisement</span>
                            </h1>
                            <p className="text-text-dim text-lg font-medium max-w-md">
                                Set your own prices and trade crypto directly with users worldwide.
                            </p>
                        </div>
                        <div className="hidden lg:flex items-center gap-4 bg-surface border border-border p-6 rounded-[2rem]">
                            <div className="w-12 h-12 bg-primary/10 rounded-2xl flex items-center justify-center">
                                <ShieldCheck className="w-6 h-6 text-primary" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">Escrow Protected</p>
                                <p className="text-sm font-bold text-white">Institutional Grade Security</p>
                            </div>
                        </div>
                    </div>

                    <form onSubmit={handleSubmit} className="space-y-8 relative">
                        {/* Global dropdown overlay */}
                        {(isCryptoOpen || isFiatOpen || isWindowOpen) && (
                            <div className="fixed inset-0 z-40" onClick={closeDropdowns} />
                        )}

                        {/* Step 1: Ad Type & Asset */}
                        <div className="glass rounded-[3rem] p-8 md:p-12 border-border/50 relative">
                            {/* Watermark restricted to container bounds */}
                            <div className="absolute inset-0 overflow-hidden rounded-[3rem] pointer-events-none">
                                <div className="absolute top-0 right-0 p-8 text-primary/10 select-none">
                                    <span className="text-9xl font-black italic">01</span>
                                </div>
                            </div>

                            <h2 className="text-xl font-black text-white uppercase tracking-tight mb-8 flex items-center relative z-10">
                                <Plus className="w-5 h-5 mr-3 text-primary" />
                                Basic Information
                            </h2>

                            <div className="grid md:grid-cols-2 gap-8">
                                <div className="space-y-3">
                                    <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">I want to</label>
                                    <div className="flex p-1.5 bg-background/50 border border-border rounded-2xl">
                                        <button
                                            type="button"
                                            onClick={() => setAdType("buy")}
                                            className={`flex-1 py-3.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all ${adType === "buy" ? "bg-white text-background shadow-lg" : "text-text-dim hover:text-white"}`}
                                        >
                                            Buy Crypto
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => setAdType("sell")}
                                            className={`flex-1 py-3.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all ${adType === "sell" ? "bg-white text-background shadow-lg" : "text-text-dim hover:text-white"}`}
                                        >
                                            Sell Crypto
                                        </button>
                                    </div>
                                </div>

                                <div className="grid grid-cols-2 gap-4">
                                    <div className="space-y-3 relative z-50">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Asset</label>
                                        <button
                                            type="button"
                                            onClick={() => { closeDropdowns(); setIsCryptoOpen(!isCryptoOpen); }}
                                            className="w-full h-[52px] px-4 flex items-center justify-between bg-background/50 border border-border hover:border-white/20 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary transition-all"
                                        >
                                            <div className="flex items-baseline gap-1.5 truncate pr-2">
                                                <span>{assets.find(a => a.id === cryptoId)?.symbol}</span>
                                                <span className="text-[10px] text-text-dim font-medium truncate">{assets.find(a => a.id === cryptoId)?.name}</span>
                                            </div>
                                            <svg className={`shrink-0 w-4 h-4 text-text-dim transition-transform ${isCryptoOpen ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                        </button>

                                        {isCryptoOpen && (
                                            <div className="absolute w-full mt-2 bg-surface border border-border rounded-2xl shadow-2xl overflow-hidden shadow-black/50 py-1">
                                                {assets.map(a => (
                                                    <button
                                                        key={a.id}
                                                        type="button"
                                                        onClick={() => { setCryptoId(a.id); setIsCryptoOpen(false); }}
                                                        className={`w-full flex items-center justify-between px-4 py-3 text-sm transition-all text-left group ${cryptoId === a.id ? "bg-primary/10 text-white font-black" : "text-text-dim hover:text-white hover:bg-white/5 font-bold"}`}
                                                    >
                                                        <span>{a.symbol}</span>
                                                        <span className={`text-[10px] font-medium transition-colors ${cryptoId === a.id ? "text-primary/80" : "text-white/30 group-hover:text-white/50"}`}>{a.name}</span>
                                                    </button>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                    <div className="space-y-3 relative z-40">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Fiat</label>
                                        <button
                                            type="button"
                                            onClick={() => { closeDropdowns(); setIsFiatOpen(!isFiatOpen); }}
                                            className="w-full h-[52px] px-4 flex items-center justify-between bg-background/50 border border-border hover:border-white/20 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary transition-all"
                                        >
                                            <div className="flex items-baseline gap-1.5 truncate pr-2">
                                                <span>{fiats.find(f => f.id === fiatId)?.symbol}</span>
                                                <span className="text-[10px] text-text-dim font-medium truncate">{fiats.find(f => f.id === fiatId)?.name}</span>
                                            </div>
                                            <svg className={`shrink-0 w-4 h-4 text-text-dim transition-transform ${isFiatOpen ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                        </button>

                                        {isFiatOpen && (
                                            <div className="absolute w-full mt-2 bg-surface border border-border rounded-2xl shadow-2xl overflow-hidden shadow-black/50 py-1">
                                                {fiats.map(f => (
                                                    <button
                                                        key={f.id}
                                                        type="button"
                                                        onClick={() => { setFiatId(f.id); setIsFiatOpen(false); }}
                                                        className={`w-full flex items-center justify-between px-4 py-3 text-sm transition-all text-left group ${fiatId === f.id ? "bg-primary/10 text-white font-black" : "text-text-dim hover:text-white hover:bg-white/5 font-bold"}`}
                                                    >
                                                        <span>{f.symbol}</span>
                                                        <span className={`text-[10px] font-medium transition-colors ${fiatId === f.id ? "text-primary/80" : "text-white/30 group-hover:text-white/50"}`}>{f.name}</span>
                                                    </button>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Step 2: Pricing */}
                        <div className="glass rounded-[3rem] p-8 md:p-12 border-border/50 relative">
                            <div className="absolute inset-0 overflow-hidden rounded-[3rem] pointer-events-none">
                                <div className="absolute top-0 right-0 p-8 text-primary/10 select-none">
                                    <span className="text-9xl font-black italic">02</span>
                                </div>
                            </div>

                            <h2 className="text-xl font-black text-white uppercase tracking-tight mb-8 flex items-center relative z-10">
                                <Zap className="w-5 h-5 mr-3 text-primary" />
                                Pricing Structure
                            </h2>

                            <div className="space-y-8">
                                <div className="flex flex-wrap gap-4">
                                    <button
                                        type="button"
                                        onClick={() => setPriceType("fixed")}
                                        className={`px-8 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest border transition-all ${priceType === "fixed" ? "bg-primary border-primary text-white" : "bg-surface border-border text-text-dim hover:border-primary"}`}
                                    >
                                        Fixed Price
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setPriceType("floating")}
                                        className={`px-8 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest border transition-all ${priceType === "floating" ? "bg-primary border-primary text-white" : "bg-surface border-border text-text-dim hover:border-primary"}`}
                                    >
                                        Floating Margin
                                    </button>
                                </div>

                                <div className="grid md:grid-cols-3 gap-8">
                                    <div className="space-y-3">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">
                                            {priceType === "fixed" ? "My Price" : "Margin (%)"}
                                        </label>
                                        <div className="relative">
                                            <input
                                                type="number"
                                                value={price}
                                                onChange={(e) => setPrice(e.target.value)}
                                                required
                                                step="any"
                                                className="w-full bg-background/50 border border-border rounded-2xl p-4 pl-12 text-sm font-bold text-white outline-none focus:border-primary transition-all"
                                                placeholder={priceType === "fixed" ? "0.00" : "100.0"}
                                            />
                                            {priceType === "fixed" ? (
                                                <DollarSign className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-text-dim" />
                                            ) : (
                                                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-sm font-black text-text-dim">%</span>
                                            )}
                                        </div>
                                    </div>
                                    <div className="space-y-3">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Min Order</label>
                                        <input
                                            type="number"
                                            value={minAmount}
                                            onChange={(e) => setMinAmount(e.target.value)}
                                            required
                                            className="w-full bg-background/50 border border-border rounded-2xl p-4 text-sm font-bold text-white outline-none focus:border-primary transition-all"
                                            placeholder="10.00"
                                        />
                                    </div>
                                    <div className="space-y-3">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Max Order</label>
                                        <input
                                            type="number"
                                            value={maxAmount}
                                            onChange={(e) => setMaxAmount(e.target.value)}
                                            required
                                            className="w-full bg-background/50 border border-border rounded-2xl p-4 text-sm font-bold text-white outline-none focus:border-primary transition-all"
                                            placeholder="5000.00"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* Step 3: Methods & Terms */}
                        <div className="glass rounded-[3rem] p-8 md:p-12 border-border/50 relative">
                            <div className="absolute inset-0 overflow-hidden rounded-[3rem] pointer-events-none">
                                <div className="absolute top-0 right-0 p-8 text-primary/10 select-none">
                                    <span className="text-9xl font-black italic">03</span>
                                </div>
                            </div>

                            <h2 className="text-xl font-black text-white uppercase tracking-tight mb-8 flex items-center relative z-10">
                                <ShieldCheck className="w-5 h-5 mr-3 text-primary" />
                                Trade Requirements
                            </h2>

                            <div className="space-y-8">
                                <div className="space-y-4">
                                    <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Payment Methods</label>
                                    <div className="flex flex-wrap gap-3">
                                        {paymentMethods.map((method) => (
                                            <button
                                                key={method.id}
                                                type="button"
                                                onClick={() => toggleMethod(method.id)}
                                                className={`px-4 py-2.5 rounded-xl text-xs font-bold border transition-all ${selectedMethods.includes(method.id) ? "bg-white text-background border-white" : "bg-surface border-border text-text-dim hover:border-primary"}`}
                                            >
                                                {method.name}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                <div className="grid md:grid-cols-2 gap-8">
                                    <div className="space-y-3 relative z-50">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Order Expiry (Minutes)</label>
                                        <button
                                            type="button"
                                            onClick={() => { closeDropdowns(); setIsWindowOpen(!isWindowOpen); }}
                                            className="w-full flex items-center justify-between bg-background/50 border border-border rounded-2xl p-4 text-sm font-bold text-white outline-none focus:border-primary transition-all text-left"
                                        >
                                            <span>{paymentWindow} Minutes</span>
                                            <svg className={`w-4 h-4 text-text-dim transition-transform ${isWindowOpen ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                        </button>

                                        {isWindowOpen && (
                                            <div className="absolute w-full mt-2 bg-surface border border-border rounded-2xl shadow-2xl overflow-hidden shadow-black/50">
                                                {[15, 30, 45, 60].map(val => (
                                                    <button
                                                        key={val}
                                                        type="button"
                                                        onClick={() => { setPaymentWindow(val); setIsWindowOpen(false); }}
                                                        className={`w-full flex items-center justify-start p-4 text-sm font-bold transition-all text-left ${paymentWindow === val ? "bg-primary text-white" : "text-text-dim hover:text-white hover:bg-white/5"}`}
                                                    >
                                                        {val} Minutes
                                                    </button>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                    <div className="space-y-3">
                                        <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Security</label>
                                        <div className="bg-surface/50 border border-border p-4 rounded-2xl flex items-center gap-3">
                                            <LockIcon className="w-4 h-4 text-accent" />
                                            <span className="text-xs font-bold text-text-dim">Escrow protected</span>
                                        </div>
                                    </div>
                                </div>

                                <div className="space-y-3">
                                    <label className="block text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-1">Terms & Conditions</label>
                                    <textarea
                                        rows={4}
                                        value={tradeTerms}
                                        onChange={(e) => setTradeTerms(e.target.value)}
                                        className="w-full bg-background/50 border border-border rounded-2xl p-6 text-sm font-bold text-white outline-none focus:border-primary transition-all resize-none placeholder:text-text-dim/30"
                                        placeholder="Enter your trade instructions and terms here..."
                                    />
                                </div>
                            </div>
                        </div>

                        {error && (
                            <motion.div
                                initial={{ opacity: 0, x: -10 }}
                                animate={{ opacity: 1, x: 0 }}
                                className="bg-red-500/10 border border-red-500/20 p-4 rounded-2xl flex items-center gap-3 text-red-500"
                            >
                                <AlertCircle className="w-5 h-5 flex-shrink-0" />
                                <span className="text-sm font-bold">{error}</span>
                            </motion.div>
                        )}

                        <div className="flex flex-col md:flex-row items-center gap-6 pt-4">
                            <button
                                type="submit"
                                disabled={isSubmitting || selectedMethods.length === 0}
                                className="w-full md:w-auto px-12 py-5 bg-white text-background rounded-[2rem] font-black uppercase tracking-widest text-sm hover:scale-105 active:scale-95 transition-all shadow-xl shadow-white/10 disabled:opacity-50 disabled:hover:scale-100 flex items-center justify-center gap-2"
                            >
                                {isSubmitting ? <Loader2 className="w-5 h-5 animate-spin" /> : <>Post Advertisement <ArrowRight className="w-5 h-5" /></>}
                            </button>
                            <p className="text-[10px] text-text-dim font-medium max-w-xs text-center md:text-left leading-relaxed">
                                By posting this advertisement, you agree to Cryplio's P2P Trading Policies and Fee Schedule.
                            </p>
                        </div>
                    </form>
                </div>
            </section>

            <Footer />
        </main>
    );
};

const LockIcon = ({ className }: { className?: string }) => (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
    </svg>
);

export default CreateAdPage;
