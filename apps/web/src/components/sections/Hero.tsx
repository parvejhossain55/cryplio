"use client";

import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, ChevronRight, CheckCircle2, ShieldCheck, Zap, ChevronDown } from "lucide-react";

const Hero = () => {
    const [tradeType, setTradeType] = React.useState<"buy" | "sell">("buy");
    const [amount, setAmount] = React.useState("1000");
    const [liveCounter, setLiveCounter] = React.useState(2450);
    const [fiat, setFiat] = React.useState("USD");
    const [asset, setAsset] = React.useState("BTC");
    const [showFiatMenu, setShowFiatMenu] = React.useState(false);
    const [showAssetMenu, setShowAssetMenu] = React.useState(false);
    const fiatRef = React.useRef<HTMLDivElement>(null);
    const assetRef = React.useRef<HTMLDivElement>(null);

    const assets = [
        { name: "BTC", icon: "₿", color: "bg-[#f7931a]", rate: 0.000014 },
        { name: "ETH", icon: "Ξ", color: "bg-[#627eea]", rate: 0.00045 },
        { name: "USDT", icon: "₮", color: "bg-[#26a17b]", rate: 1.0 },
    ];

    const fiats = ["USD", "EUR", "GBP", "NGN", "PKR"];

    const currentAsset = assets.find(a => a.name === asset) || assets[0];
    const baseRate = currentAsset.rate;
    const finalRate = tradeType === "buy" ? baseRate : baseRate * 1.05;
    const receive = (parseFloat(amount.replace(/,/g, "")) * finalRate).toFixed(asset === "USDT" ? 2 : 6);

    React.useEffect(() => {
        const interval = setInterval(() => {
            setLiveCounter(prev => prev + Math.floor(Math.random() * 3));
        }, 3000);

        const handleClickOutside = (event: MouseEvent) => {
            if (fiatRef.current && !fiatRef.current.contains(event.target as Node)) {
                setShowFiatMenu(false);
            }
            if (assetRef.current && !assetRef.current.contains(event.target as Node)) {
                setShowAssetMenu(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            clearInterval(interval);
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    return (
        <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden">
            {/* Background Glow */}
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-6xl h-full -z-10 bg-radial-[circle_at_50%_0%] from-primary/10 via-transparent to-transparent opacity-50" />

            <div className="container mx-auto px-4 md:px-6">
                <div className="flex flex-col lg:flex-row items-center gap-16">
                    {/* Text Content */}
                    <div className="flex-1 text-center lg:text-left">
                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.5 }}
                            className="inline-flex items-center space-x-3 bg-surface-light border border-border px-4 py-2 rounded-full text-xs font-medium mb-6 shadow-sm"
                        >
                            <span className="relative flex h-2 w-2">
                                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent opacity-75"></span>
                                <span className="relative inline-flex rounded-full h-2 w-2 bg-accent"></span>
                            </span>
                            <span className="text-text-dim"><span className="text-white font-bold">{liveCounter.toLocaleString()}</span> trades completed in last 24h</span>
                        </motion.div>

                        <motion.h1
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.5, delay: 0.1 }}
                            className="text-4xl md:text-6xl lg:text-7xl font-extrabold leading-[1.1] mb-6"
                        >
                            The safest way to <br />
                            <span className="gradient-text">{tradeType === "buy" ? "Buy" : "Sell"} {asset}</span>
                        </motion.h1>

                        <motion.p
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.5, delay: 0.2 }}
                            className="text-text-dim text-lg md:text-xl max-w-2xl mx-auto lg:mx-0 mb-8 leading-relaxed"
                        >
                            Cryplio is the premier P2P marketplace for traders worldwide. Secure escrow, 0% fees for makers, and instant settlement.
                        </motion.p>

                        <motion.div
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ duration: 0.5, delay: 0.3 }}
                            className="flex flex-col sm:flex-row items-center justify-center lg:justify-start space-y-4 sm:space-y-0 sm:space-x-4 mb-10"
                        >
                            <button
                                onClick={() => window.location.href = "/register"}
                                className="w-full sm:w-auto bg-primary hover:bg-primary-dark text-white px-8 py-4 rounded-2xl text-lg font-bold transition-all shadow-xl shadow-primary/30 flex items-center justify-center space-x-2 group"
                            >
                                <span>Start Trading</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </button>
                            <button
                                onClick={() => document.getElementById("marketplace")?.scrollIntoView({ behavior: "smooth" })}
                                className="w-full sm:w-auto glass hover:bg-surface-light px-8 py-4 rounded-2xl text-lg font-bold transition-all flex items-center justify-center border border-border"
                            >
                                Explore Market
                            </button>
                        </motion.div>

                        <motion.div
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            transition={{ duration: 1, delay: 0.5 }}
                            className="grid grid-cols-3 gap-4 md:gap-8 max-w-md mx-auto lg:mx-0 pt-4 border-t border-border"
                        >
                            {[
                                { icon: ShieldCheck, label: "Escrow Protec" },
                                { icon: CheckCircle2, label: "KYC Verified" },
                                { icon: Zap, label: "Instant Fill" },
                            ].map((item, i) => (
                                <div key={i} className="flex flex-col items-center lg:items-start space-y-2">
                                    <item.icon className="w-5 h-5 text-accent" />
                                    <span className="text-[10px] md:text-xs font-semibold uppercase tracking-wider text-text-dim">
                                        {item.label}
                                    </span>
                                </div>
                            ))}
                        </motion.div>
                    </div>

                    {/* Quick Trade Widget */}
                    <motion.div
                        initial={{ opacity: 0, x: 50 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ duration: 0.8, delay: 0.2 }}
                        className="flex-1 w-full max-w-xl"
                    >
                        <div className="relative">
                            {/* Decorative elements */}
                            <div className="absolute -top-10 -right-10 w-40 h-40 bg-primary/20 rounded-full blur-[80px]" />
                            <div className="absolute -bottom-10 -left-10 w-40 h-40 bg-secondary/20 rounded-full blur-[80px]" />

                            <div className="glass rounded-[32px] p-8 md:p-10 border-white/5 relative overflow-visible shadow-2xl">
                                <div className="flex items-center justify-between mb-8">
                                    <h3 className="text-xl font-bold">Quick {tradeType === "buy" ? "Buy" : "Sell"}</h3>
                                    <div className="flex bg-surface rounded-lg p-1">
                                        <button
                                            onClick={() => setTradeType("buy")}
                                            className={`px-4 py-1.5 rounded-md text-sm font-semibold transition-all ${tradeType === "buy" ? "bg-primary text-white shadow-sm" : "text-text-dim hover:text-white"}`}
                                        >
                                            Buy
                                        </button>
                                        <button
                                            onClick={() => setTradeType("sell")}
                                            className={`px-4 py-1.5 rounded-md text-sm font-semibold transition-all ${tradeType === "sell" ? "bg-primary text-white shadow-sm" : "text-text-dim hover:text-white"}`}
                                        >
                                            Sell
                                        </button>
                                    </div>
                                </div>

                                <div className="space-y-6">
                                    <div className="relative">
                                        <label className="text-xs font-bold text-text-dim uppercase tracking-widest block mb-2">{tradeType === "buy" ? "Spend" : "Sell"}</label>
                                        <div className="bg-surface rounded-2xl p-4 border border-border flex items-center justify-between group focus-within:border-primary transition-colors">
                                            <input
                                                type="text"
                                                value={amount}
                                                onChange={(e) => setAmount(e.target.value)}
                                                className="bg-transparent text-2xl font-bold outline-none w-full"
                                            />
                                            <div className="relative" ref={fiatRef}>
                                                <button
                                                    type="button"
                                                    onClick={() => {
                                                        setShowFiatMenu(!showFiatMenu);
                                                        setShowAssetMenu(false);
                                                    }}
                                                    className="flex items-center space-x-2 bg-surface-light px-3 py-1.5 rounded-xl border border-border cursor-pointer hover:bg-border transition-colors outline-none"
                                                >
                                                    <span className="font-bold">{fiat}</span>
                                                    <ChevronDown className={`w-4 h-4 transition-transform ${showFiatMenu ? "rotate-180" : ""}`} />
                                                </button>
                                                {showFiatMenu && (
                                                    <div className="absolute right-0 mt-2 w-24 bg-surface-light border border-border rounded-xl shadow-2xl z-50 overflow-hidden">
                                                        {fiats.map(f => (
                                                            <button
                                                                key={f}
                                                                type="button"
                                                                onClick={() => { setFiat(f); setShowFiatMenu(false); }}
                                                                className="w-full text-left px-4 py-2 text-sm hover:bg-primary hover:text-white transition-colors font-bold"
                                                            >
                                                                {f}
                                                            </button>
                                                        ))}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex justify-center -my-3 relative z-10">
                                        <div className="w-10 h-10 bg-primary rounded-full flex items-center justify-center shadow-lg shadow-primary/30 border-4 border-background">
                                            <ArrowRight className="w-5 h-5 text-white rotate-90" />
                                        </div>
                                    </div>

                                    <div className="relative">
                                        <label className="text-xs font-bold text-text-dim uppercase tracking-widest block mb-2">Receive</label>
                                        <div className="bg-surface rounded-2xl p-4 border border-border flex items-center justify-between group focus-within:border-primary transition-colors">
                                            <input
                                                type="text"
                                                value={isNaN(parseFloat(receive)) ? "0.00" : receive}
                                                readOnly
                                                className="bg-transparent text-2xl font-bold outline-none w-full text-accent"
                                            />
                                            <div className="relative" ref={assetRef}>
                                                <button
                                                    type="button"
                                                    onClick={() => {
                                                        setShowAssetMenu(!showAssetMenu);
                                                        setShowFiatMenu(false);
                                                    }}
                                                    className="flex items-center space-x-2 bg-surface-light px-3 py-1.5 rounded-xl border border-border cursor-pointer hover:bg-border transition-colors outline-none"
                                                >
                                                    <div className={`w-6 h-6 ${currentAsset.color} rounded-full flex items-center justify-center text-[10px] text-white`}>{currentAsset.icon}</div>
                                                    <span className="font-bold">{asset}</span>
                                                    <ChevronDown className={`w-4 h-4 transition-transform ${showAssetMenu ? "rotate-180" : ""}`} />
                                                </button>
                                                {showAssetMenu && (
                                                    <div className="absolute right-0 mt-2 w-32 bg-surface-light border border-border rounded-xl shadow-2xl z-50 overflow-hidden">
                                                        {assets.map(a => (
                                                            <button
                                                                key={a.name}
                                                                type="button"
                                                                onClick={() => { setAsset(a.name); setShowAssetMenu(false); }}
                                                                className="w-full flex items-center space-x-3 px-4 py-3 text-sm hover:bg-primary hover:text-white transition-colors font-bold"
                                                            >
                                                                <div className={`w-5 h-5 ${a.color} rounded-full flex items-center justify-center text-[8px] text-white`}>{a.icon}</div>
                                                                <span>{a.name}</span>
                                                            </button>
                                                        ))}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    </div>

                                    <div className="pt-4">
                                        <button
                                            onClick={() => window.location.href = "/register"}
                                            className="w-full bg-accent hover:bg-accent/90 text-background px-8 py-4 rounded-2xl text-lg font-bold transition-all shadow-xl shadow-accent/20 active:scale-[0.98]"
                                        >
                                            {tradeType === "buy" ? `Buy ${asset} Now` : `Sell ${asset} Now`}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};

export default Hero;
