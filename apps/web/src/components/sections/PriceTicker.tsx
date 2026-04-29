"use client";

import React from "react";
import { motion } from "framer-motion";
import { TrendingUp, TrendingDown } from "lucide-react";

const initialCoins = [
    { name: "BTC", price: 96420.50, change: "+2.4%", up: true },
    { name: "ETH", price: 2540.12, change: "-0.8%", up: false },
    { name: "USDT", price: 1.00, change: "+0.01%", up: true },
    { name: "BNB", price: 598.45, change: "+1.2%", up: true },
    { name: "SOL", price: 145.20, change: "+5.4%", up: true },
    { name: "XRP", price: 0.62, change: "-1.5%", up: false },
];

const PriceTicker = () => {
    const [coins, setCoins] = React.useState(initialCoins);

    React.useEffect(() => {
        const interval = setInterval(() => {
            setCoins(currentCoins =>
                currentCoins.map(coin => ({
                    ...coin,
                    price: coin.price + (Math.random() - 0.5) * (coin.price * 0.0001)
                }))
            );
        }, 3000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="bg-surface border-b border-border py-2 overflow-hidden">
            <div className="flex whitespace-nowrap animate-scroll items-center gap-12 px-6">
                {[...coins, ...coins].map((coin, i) => (
                    <div key={i} className="flex items-center space-x-3 group cursor-pointer">
                        <span className="font-bold text-sm text-white">{coin.name}</span>
                        <span className="text-sm font-medium text-text-dim">${coin.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                        <div className={`flex items-center text-xs font-bold ${coin.up ? "text-accent" : "text-rose-500"}`}>
                            {coin.up ? <TrendingUp className="w-3 h-3 mr-1" /> : <TrendingDown className="w-3 h-3 mr-1" />}
                            {coin.change}
                        </div>
                    </div>
                ))}
            </div>

            <style jsx>{`
        @keyframes scroll {
          0% { transform: translateX(0); }
          100% { transform: translateX(-50%); }
        }
        .animate-scroll {
          animation: scroll 30s linear infinite;
        }
        .animate-scroll:hover {
          animation-play-state: paused;
        }
      `}</style>
        </div>
    );
};

export default PriceTicker;
