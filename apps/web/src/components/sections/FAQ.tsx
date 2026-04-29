"use client";

import React, { useState } from "react";
import { ChevronDown, ChevronUp } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const FAQS = [
    {
        q: "Is Cryplio safe to use?",
        a: "Yes, Cryplio uses a secure escrow system. When a trade starts, the seller's crypto is held in nuestra locked account until payment is confirmed by both parties."
    },
    {
        q: "How long does a trade take?",
        a: "Average trades are completed within 5-15 minutes. It mostly depends on how quickly the buyer sends the payment and the seller confirms receipt."
    },
    {
        q: "What payment methods are supported?",
        a: "We support over 100+ payment methods including Bank Transfers, PayPal, Mobile Wallets like EaseyPaisa, bKash, and local bank methods in over 50 countries."
    },
    {
        q: "Are there any fees for trading?",
        a: "Trading is free for buyers. Sellers (makers) pay a small service fee of 0.5% per completed trade. There are no hidden charges."
    }
];

const FAQ = () => {
    const [openIndex, setOpenIndex] = useState<number | null>(0);

    return (
        <section className="py-24 px-4">
            <div className="container mx-auto max-w-3xl">
                <div className="text-center mb-16">
                    <h2 className="text-3xl md:text-5xl font-black mb-4">Common Questions</h2>
                    <p className="text-text-dim px-4">Everything you need to know about trading on Cryplio.</p>
                </div>

                <div className="space-y-4">
                    {FAQS.map((faq, i) => (
                        <div key={i} className="glass rounded-3xl border-white/5 overflow-hidden">
                            <button
                                onClick={() => setOpenIndex(openIndex === i ? null : i)}
                                className="w-full flex items-center justify-between p-6 text-left hover:bg-white/5 transition-colors"
                                aria-expanded={openIndex === i}
                            >
                                <span className="font-bold md:text-lg">{faq.q}</span>
                                {openIndex === i ? <ChevronUp className="w-5 h-5 text-primary" /> : <ChevronDown className="w-5 h-5 text-text-dim" />}
                            </button>
                            <AnimatePresence>
                                {openIndex === i && (
                                    <motion.div
                                        initial={{ height: 0, opacity: 0 }}
                                        animate={{ height: "auto", opacity: 1 }}
                                        exit={{ height: 0, opacity: 0 }}
                                        transition={{ duration: 0.3 }}
                                    >
                                        <div className="px-6 pb-6 text-text-dim text-sm leading-relaxed border-t border-white/5 pt-4">
                                            {faq.a}
                                        </div>
                                    </motion.div>
                                )}
                            </AnimatePresence>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
};

export default FAQ;
