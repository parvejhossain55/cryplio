"use client";

import React, { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    ArrowDown,
    Settings,
    RefreshCw,
    ChevronDown,
    Info,
    Search,
    X,
    Zap,
    ShieldCheck,
    CheckCircle2,
    AlertCircle,
    ArrowUpRight
} from "lucide-react";

interface Token {
    symbol: string;
    name: string;
    balance: number;
    price: number;
    color: string;
    icon: string;
}

const tokens: Token[] = [
    { symbol: "USDT", name: "Tether USD", balance: 1245.50, price: 1.00, color: "text-emerald-500", icon: "https://cryptologos.cc/logos/tether-usdt-logo.svg?v=040" },
    { symbol: "BTC", name: "Bitcoin", balance: 0.045, price: 65420.00, color: "text-amber-500", icon: "https://cryptologos.cc/logos/bitcoin-btc-logo.svg?v=040" },
    { symbol: "ETH", name: "Ethereum", balance: 1.20, price: 3450.00, color: "text-indigo-500", icon: "https://cryptologos.cc/logos/ethereum-eth-logo.svg?v=040" },
    { symbol: "BNB", name: "BNB", balance: 15.4, price: 580.00, color: "text-yellow-500", icon: "https://cryptologos.cc/logos/bnb-bnb-logo.svg?v=040" },
    { symbol: "SOL", name: "Solana", balance: 45.0, price: 145.00, color: "text-purple-500", icon: "https://cryptologos.cc/logos/solana-sol-logo.svg?v=040" },
];

const SwapSection = () => {
    const [fromToken, setFromToken] = useState(tokens[0]);
    const [toToken, setToToken] = useState(tokens[2]);
    const [fromAmount, setFromAmount] = useState("");
    const [toAmount, setToAmount] = useState("");
    const [lastInput, setLastInput] = useState<"from" | "to">("from");

    const [isLoading, setIsLoading] = useState(false);
    const [isSuccess, setIsSuccess] = useState(false);
    const [isTokenModalOpen, setIsTokenModalOpen] = useState(false);
    const [isSettingsOpen, setIsSettingsOpen] = useState(false);
    const [modalType, setModalType] = useState<"from" | "to">("from");
    const [searchQuery, setSearchQuery] = useState("");
    const [slippage, setSlippage] = useState("0.5");
    const [error, setError] = useState<string | null>(null);

    // Calculate exchange
    useEffect(() => {
        if (lastInput === "from") {
            if (fromAmount && !isNaN(parseFloat(fromAmount))) {
                const amount = parseFloat(fromAmount);
                const result = (amount * fromToken.price) / toToken.price;
                setToAmount(result.toFixed(6));

                if (amount > fromToken.balance) {
                    setError(`Insufficient ${fromToken.symbol} balance`);
                } else {
                    setError(null);
                }
            } else {
                setToAmount("");
                setError(null);
            }
        }
    }, [fromAmount, fromToken, toToken, lastInput]);

    useEffect(() => {
        if (lastInput === "to") {
            if (toAmount && !isNaN(parseFloat(toAmount))) {
                const amount = parseFloat(toAmount);
                const result = (amount * toToken.price) / fromToken.price;
                setFromAmount(result.toFixed(6));

                if (parseFloat(result.toFixed(6)) > fromToken.balance) {
                    setError(`Insufficient ${fromToken.symbol} balance`);
                } else {
                    setError(null);
                }
            } else {
                setFromAmount("");
                setError(null);
            }
        }
    }, [toAmount, fromToken, toToken, lastInput]);

    const handleSwap = () => {
        if (error) return;
        setIsLoading(true);
        setTimeout(() => {
            setIsLoading(false);
            setIsSuccess(true);
        }, 2000);
    };

    const reverseTokens = () => {
        const temp = fromToken;
        setFromToken(toToken);
        setToToken(temp);
        setFromAmount(toAmount);
        setLastInput("from");
    };

    const setMaxAmount = () => {
        setFromAmount(fromToken.balance.toString());
        setLastInput("from");
    };

    const openTokenModal = (type: "from" | "to") => {
        setModalType(type);
        setIsTokenModalOpen(true);
    };

    const selectToken = (token: Token) => {
        if (modalType === "from") {
            setFromToken(token);
        } else {
            setToToken(token);
        }
        setIsTokenModalOpen(false);
        setSearchQuery("");
    };

    const filteredTokens = tokens.filter(t =>
        t.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
        t.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <section className="py-24 relative overflow-hidden">
            {/* Background Decorations */}
            <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/10 blur-[120px] rounded-full -z-10" />
            <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-secondary/10 blur-[120px] rounded-full -z-10" />

            <div className="container mx-auto px-4">
                <div className="max-w-[480px] mx-auto">
                    {/* Header */}
                    <div className="flex items-center justify-between mb-8">
                        <div>
                            <h1 className="text-3xl font-black tracking-tight">Instant <span className="gradient-text">Swap</span></h1>
                            <p className="text-text-dim text-sm font-medium">Flash-speed asset exchange</p>
                        </div>
                        <div className="flex items-center space-x-2">
                            <button className="p-3 bg-surface border border-border rounded-2xl hover:border-primary transition-all text-text-dim hover:text-white">
                                <RefreshCw className="w-5 h-5" />
                            </button>
                            <button
                                onClick={() => setIsSettingsOpen(true)}
                                className="p-3 bg-surface border border-border rounded-2xl hover:border-primary transition-all text-text-dim hover:text-white"
                            >
                                <Settings className="w-5 h-5" />
                            </button>
                        </div>
                    </div>

                    {/* Swap Card */}
                    <div className="glass rounded-[40px] border-border p-6 shadow-2xl relative">
                        <div className="space-y-2">
                            {/* From Section */}
                            <div className={`bg-surface/50 border rounded-3xl p-6 transition-all focus-within:border-primary/50 group ${error ? "border-red-500/50" : "border-border"}`}>
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-xs font-black text-text-dim uppercase tracking-widest">You Sell</span>
                                    <div className="flex items-center space-x-2">
                                        <span className="text-xs font-bold text-text-dim">Balance: {fromToken.balance}</span>
                                        <button
                                            onClick={setMaxAmount}
                                            className="text-[10px] font-black uppercase text-primary hover:text-white bg-primary/10 px-2 py-0.5 rounded-lg transition-colors"
                                        >
                                            Max
                                        </button>
                                    </div>
                                </div>
                                <div className="flex items-center space-x-4">
                                    <input
                                        type="number"
                                        placeholder="0.0"
                                        value={fromAmount}
                                        onChange={(e) => {
                                            setFromAmount(e.target.value);
                                            setLastInput("from");
                                        }}
                                        className="bg-transparent text-3xl font-black outline-none w-full placeholder:text-text-dim/20"
                                    />
                                    <button
                                        onClick={() => openTokenModal("from")}
                                        className="flex items-center space-x-2 bg-surface border border-border px-4 py-2.5 rounded-2xl hover:border-primary transition-all whitespace-nowrap"
                                    >
                                        <img src={fromToken.icon} alt={fromToken.symbol} className="w-6 h-6 rounded-full" />
                                        <span className="font-black">{fromToken.symbol}</span>
                                        <ChevronDown className="w-4 h-4 text-text-dim" />
                                    </button>
                                </div>
                            </div>

                            {/* Splitter Button */}
                            <div className="relative h-2 flex justify-center items-center z-10">
                                <button
                                    onClick={reverseTokens}
                                    className="absolute bg-surface border-4 border-background w-12 h-12 rounded-2xl flex items-center justify-center hover:scale-110 transition-all hover:border-primary/20 group shadow-xl"
                                >
                                    <ArrowDown className="w-5 h-5 text-primary group-hover:scale-110 transition-transform" />
                                </button>
                            </div>

                            {/* To Section */}
                            <div className="bg-surface/50 border border-border rounded-3xl p-6 transition-all focus-within:border-primary/50 group">
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-xs font-black text-text-dim uppercase tracking-widest">You Get</span>
                                    <span className="text-xs font-bold text-text-dim">Balance: {toToken.balance}</span>
                                </div>
                                <div className="flex items-center space-x-4">
                                    <input
                                        type="number"
                                        placeholder="0.0"
                                        value={toAmount}
                                        onChange={(e) => {
                                            setToAmount(e.target.value);
                                            setLastInput("to");
                                        }}
                                        className="bg-transparent text-3xl font-black outline-none w-full placeholder:text-text-dim/20"
                                    />
                                    <button
                                        onClick={() => openTokenModal("to")}
                                        className="flex items-center space-x-2 bg-surface border border-border px-4 py-2.5 rounded-2xl hover:border-primary transition-all whitespace-nowrap"
                                    >
                                        <img src={toToken.icon} alt={toToken.symbol} className="w-6 h-6 rounded-full" />
                                        <span className="font-black">{toToken.symbol}</span>
                                        <ChevronDown className="w-4 h-4 text-text-dim" />
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Error Message */}
                        {error && (
                            <motion.div
                                initial={{ opacity: 0, y: -10 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="mt-4 flex items-center space-x-2 text-red-400 bg-red-400/10 p-3 rounded-2xl border border-red-400/20"
                            >
                                <AlertCircle className="w-4 h-4 shrink-0" />
                                <span className="text-xs font-bold">{error}</span>
                            </motion.div>
                        )}

                        {/* Order Details */}
                        <div className="mt-8 space-y-4 px-2">
                            <div className="flex justify-between items-center text-sm">
                                <div className="flex items-center space-x-1 text-text-dim font-medium">
                                    <span>Exchange Rate</span>
                                    <Info className="w-3.5 h-3.5" />
                                </div>
                                <span className="font-bold">1 {fromToken.symbol} = {(fromToken.price / toToken.price).toFixed(6)} {toToken.symbol}</span>
                            </div>
                            <div className="flex justify-between items-center text-sm">
                                <div className="flex items-center space-x-1 text-text-dim font-medium">
                                    <span>Slippage Tolerance</span>
                                </div>
                                <span className="text-accent font-bold">{slippage}%</span>
                            </div>
                            <div className="flex justify-between items-center text-sm">
                                <div className="flex items-center space-x-1 text-text-dim font-medium">
                                    <span>Network Fee</span>
                                </div>
                                <span className="text-text-dim/50 font-bold">~${(parseFloat(fromAmount || "0") * 0.001 + 2).toFixed(2)}</span>
                            </div>
                        </div>

                        {/* Action Button */}
                        <button
                            disabled={!fromAmount || isLoading || !!error}
                            onClick={handleSwap}
                            className={`w-full mt-8 py-5 rounded-[24px] font-black text-lg uppercase tracking-widest transition-all shadow-xl flex items-center justify-center space-x-3
                                ${!fromAmount || isLoading || !!error
                                    ? "bg-surface text-text-dim border border-border cursor-not-allowed"
                                    : "bg-primary text-white hover:scale-[1.02] active:scale-95 shadow-primary/20 hover:shadow-primary/40"}
                            `}
                        >
                            {isLoading ? (
                                <>
                                    <RefreshCw className="w-6 h-6 animate-spin" />
                                    <span>Processing...</span>
                                </>
                            ) : (
                                <span>{error || "Swap Assets"}</span>
                            )}
                        </button>
                    </div>

                    {/* Additional Info Cards */}
                    <div className="grid grid-cols-2 gap-4 mt-6">
                        <div className="glass p-4 rounded-3xl border-border flex items-center space-x-3">
                            <div className="w-10 h-10 bg-accent/10 rounded-xl flex items-center justify-center">
                                <Zap className="w-5 h-5 text-accent" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">Speed</p>
                                <p className="text-sm font-bold">Instant</p>
                            </div>
                        </div>
                        <div className="glass p-4 rounded-3xl border-border flex items-center space-x-3">
                            <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center">
                                <ShieldCheck className="w-5 h-5 text-primary" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-widest">Security</p>
                                <p className="text-sm font-bold">Verified</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Token Selection Modal */}
            <AnimatePresence>
                {isTokenModalOpen && (
                    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-6">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setIsTokenModalOpen(false)}
                            className="absolute inset-0 bg-background/80 backdrop-blur-sm"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            className="relative w-full max-w-md glass rounded-[40px] border-border shadow-2xl overflow-hidden"
                        >
                            <div className="p-6 pb-0">
                                <div className="flex items-center justify-between mb-6">
                                    <h2 className="text-xl font-black">Select Asset</h2>
                                    <button
                                        onClick={() => setIsTokenModalOpen(false)}
                                        className="p-2 hover:bg-surface rounded-xl transition-colors"
                                    >
                                        <X className="w-6 h-6" />
                                    </button>
                                </div>
                                <div className="relative mb-6">
                                    <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5" />
                                    <input
                                        type="text"
                                        placeholder="Search by name or symbol"
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        className="w-full bg-surface/50 border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-primary transition-all font-medium"
                                    />
                                </div>
                            </div>
                            <div className="max-h-[400px] overflow-y-auto px-2 pb-6 custom-scrollbar">
                                {filteredTokens.map((token) => (
                                    <button
                                        key={token.symbol}
                                        onClick={() => selectToken(token)}
                                        className="w-full flex items-center justify-between p-4 rounded-2xl hover:bg-surface transition-all group"
                                    >
                                        <div className="flex items-center space-x-4">
                                            <div className="w-12 h-12 rounded-2xl bg-surface-light border border-border flex items-center justify-center">
                                                <img src={token.icon} alt={token.symbol} className="w-8 h-8 rounded-full" />
                                            </div>
                                            <div className="text-left">
                                                <p className="font-black group-hover:text-primary transition-colors">{token.symbol}</p>
                                                <p className="text-xs text-text-dim font-medium">{token.name}</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="font-bold">{token.balance}</p>
                                            <p className="text-[10px] text-text-dim font-bold uppercase tracking-widest">${token.price}</p>
                                        </div>
                                    </button>
                                ))}
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>

            {/* Success Modal */}
            <AnimatePresence>
                {isSuccess && (
                    <div className="fixed inset-0 z-[101] flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setIsSuccess(false)}
                            className="absolute inset-0 bg-background/80 backdrop-blur-md"
                        />
                        <motion.div
                            initial={{ opacity: 0, scale: 0.9, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.9, y: 20 }}
                            className="relative w-full max-w-sm glass rounded-[40px] border-primary/30 p-8 text-center shadow-[0_0_50px_rgba(99,102,241,0.2)]"
                        >
                            <div className="w-20 h-20 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-6">
                                <CheckCircle2 className="w-10 h-10 text-primary animate-[bounce_1s_infinite]" />
                            </div>
                            <h2 className="text-2xl font-black mb-2">Swap Successful!</h2>
                            <p className="text-text-dim text-sm mb-8 font-medium">
                                You have successfully exchanged {fromAmount} {fromToken.symbol} for {toAmount} {toToken.symbol}.
                            </p>

                            <div className="bg-surface/50 border border-border rounded-3xl p-4 mb-8 text-left space-y-3">
                                <div className="flex justify-between items-center text-xs">
                                    <span className="text-text-dim font-bold uppercase tracking-widest">Transaction ID</span>
                                    <span className="font-mono font-bold text-primary">0x8a...4b2e</span>
                                </div>
                                <div className="flex justify-between items-center text-xs">
                                    <span className="text-text-dim font-bold uppercase tracking-widest">Status</span>
                                    <span className="text-accent font-bold">Confirmed</span>
                                </div>
                            </div>

                            <button
                                onClick={() => {
                                    setIsSuccess(false);
                                    setFromAmount("");
                                    setToAmount("");
                                }}
                                className="w-full py-4 bg-primary text-white rounded-2xl font-black uppercase tracking-widest hover:scale-[1.02] active:scale-95 transition-all shadow-lg shadow-primary/20"
                            >
                                Done
                            </button>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>

            {/* Settings Sidebar */}
            <AnimatePresence>
                {isSettingsOpen && (
                    <div className="fixed inset-0 z-[101] flex justify-end">
                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            onClick={() => setIsSettingsOpen(false)}
                            className="absolute inset-0 bg-background/60 backdrop-blur-sm"
                        />
                        <motion.div
                            initial={{ x: "100%" }}
                            animate={{ x: 0 }}
                            exit={{ x: "100%" }}
                            transition={{ type: "spring", damping: 25, stiffness: 200 }}
                            className="relative w-full max-w-sm glass border-l border-border h-full p-8 shadow-2xl"
                        >
                            <div className="flex items-center justify-between mb-12">
                                <h2 className="text-2xl font-black">Swap Settings</h2>
                                <button
                                    onClick={() => setIsSettingsOpen(false)}
                                    className="p-2 hover:bg-surface rounded-xl transition-colors"
                                >
                                    <X className="w-7 h-7" />
                                </button>
                            </div>

                            <div className="space-y-12">
                                {/* Slippage Toggle */}
                                <div>
                                    <div className="flex items-center space-x-2 mb-4">
                                        <span className="text-xs font-black uppercase tracking-widest text-text-dim">Slippage Tolerance</span>
                                        <Info className="w-3.5 h-3.5 text-text-dim" />
                                    </div>
                                    <div className="grid grid-cols-4 gap-2">
                                        {["0.1", "0.5", "1.0"].map((s) => (
                                            <button
                                                key={s}
                                                onClick={() => setSlippage(s)}
                                                className={`py-3 rounded-xl text-sm font-black border transition-all ${slippage === s ? "bg-primary border-primary text-white" : "bg-surface border-border text-text-dim hover:border-primary/50"}`}
                                            >
                                                {s}%
                                            </button>
                                        ))}
                                        <div className="relative">
                                            <input
                                                type="text"
                                                placeholder="Custom"
                                                className="w-full h-full py-3 bg-surface border border-border rounded-xl text-center text-sm font-black outline-none focus:border-primary placeholder:text-[10px]"
                                            />
                                        </div>
                                    </div>
                                </div>

                                {/* Transaction Speed */}
                                <div>
                                    <div className="flex items-center space-x-2 mb-4">
                                        <span className="text-xs font-black uppercase tracking-widest text-text-dim">Transaction Speed</span>
                                    </div>
                                    <div className="space-y-2">
                                        {[
                                            { name: "Default", speed: "Instant", fee: "~$2.45" },
                                            { name: "Turbo", speed: "< 50ms", fee: "~$5.12" },
                                        ].map((opt, i) => (
                                            <button
                                                key={i}
                                                className={`w-full flex items-center justify-between p-4 rounded-2xl border transition-all ${i === 0 ? "border-primary bg-primary/5" : "border-border bg-surface/50 hover:border-primary/50"}`}
                                            >
                                                <div className="text-left">
                                                    <p className="font-bold">{opt.name}</p>
                                                    <p className="text-[10px] uppercase font-black text-text-dim tracking-widest">{opt.speed}</p>
                                                </div>
                                                <span className="text-xs font-bold text-primary">{opt.fee}</span>
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            </div>

                            <div className="absolute bottom-8 left-8 right-8">
                                <button
                                    onClick={() => setIsSettingsOpen(false)}
                                    className="w-full py-4 bg-surface border border-border rounded-2xl font-black uppercase tracking-widest hover:border-primary transition-all"
                                >
                                    Save Changes
                                </button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </section>
    );
};

export default SwapSection;
