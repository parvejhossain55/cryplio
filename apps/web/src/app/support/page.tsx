"use client";

import React, { useState, useMemo } from "react";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import {
    Search,
    Rocket,
    ShieldCheck,
    Users,
    Zap,
    Wallet,
    Code,
    ChevronDown,
    MessageCircle,
    Mail,
    Send,
    CheckCircle2,
    X,
    LifeBuoy,
    ExternalLink,
    RefreshCw
} from "lucide-react";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";

interface FAQ {
    question: string;
    answer: string;
    category: string;
}

const faqs: FAQ[] = [
    {
        category: "Verification",
        question: "How do I verify my account?",
        answer: "VERIFY YOUR IDENTITY BY UPLOADING A VALID GOVERNMENT-ISSUED ID AND A SELFIE IN YOUR PROFILE SETTINGS. VERIFICATION TYPICALLY TAKES 2-24 HOURS."
    },
    {
        category: "Security",
        question: "Is my money safe with Cryplio?",
        answer: "YES. CRYPLIO USES INSTITUTIONAL-GRADE ESCROW PROTECTION. FUNDS ARE HELD SECURELY UNTIL BOTH PARTIES COMPLY WITH THE TRADE TERMS."
    },
    {
        category: "Trading",
        question: "What happens if a seller doesn't release crypto?",
        answer: "IF A SELLER FAILS TO RELEASE CRYPTO AFTER PAYMENT, YOU CAN OPEN AN APPEAL. OUR 24/7 MODERATORS WILL REVIEW EVIDENCE AND RELEASE FUNDS TO THE RIGHTFUL OWNER."
    },
    {
        category: "Swap",
        question: "How are swap fees calculated?",
        answer: "WE PROVIDE ZERO-FEE SWAPS BY AGGREGATING LIQUIDITY FROM MULTIPLE SOURCES. THE ONLY COST IS THE NETWORK TRANSACTION FEE AND A MINIMAL SPREAD."
    },
    {
        category: "Wallet",
        question: "How do I withdraw my funds?",
        answer: "GO TO YOUR WALLET, SELECT THE ASSET, AND CLICK WITHDRAW. ENTER THE DESTINATION ADDRESS OR SELECT A LINKED BANK ACCOUNT FOR FIAT WITHDRAWALS."
    },
    {
        category: "Account",
        question: "Can I use multiple bank accounts?",
        answer: "ABSOLUTELY. YOU CAN ADD MULTIPLE PAYMENT METHODS IN YOUR ACCOUNT SETTINGS TO MAKE TRADING MORE CONVENIENT."
    }
];

const SupportPage = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [openFaqIndex, setOpenFaqIndex] = useState<number | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isSubmitted, setIsSubmitted] = useState(false);

    // Ticket Form State
    const [ticket, setTicket] = useState({
        name: "",
        email: "",
        subject: "Verification Issue",
        message: ""
    });

    const filteredFaqs = useMemo(() => {
        return faqs.filter(f =>
            f.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
            f.answer.toLowerCase().includes(searchQuery.toLowerCase())
        );
    }, [searchQuery]);

    const handleTicketSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        setTimeout(() => {
            setIsSubmitting(false);
            setIsSubmitted(true);
        }, 2000);
    };

    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            {/* Hero Section */}
            <section className="pt-32 pb-24 relative overflow-hidden">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-primary/5 blur-[120px] rounded-full -z-10" />

                <div className="container mx-auto px-4 text-center">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="space-y-6"
                    >
                        <div className="inline-flex items-center space-x-2 bg-surface border border-border px-4 py-2 rounded-2xl mx-auto">
                            <LifeBuoy className="w-4 h-4 text-primary" />
                            <span className="text-xs font-black uppercase tracking-widest text-text-dim">Cryplio Help Center</span>
                        </div>
                        <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-[0.9]">
                            How can we <br />
                            <span className="gradient-text">help you today?</span>
                        </h1>

                        <div className="max-w-2xl mx-auto relative mt-12">
                            <Search className="absolute left-6 top-1/2 -translate-y-1/2 text-text-dim w-6 h-6" />
                            <input
                                type="text"
                                placeholder="Search for verification, trading, security..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full bg-surface border border-border rounded-[32px] py-6 pl-16 pr-8 text-lg outline-none focus:border-primary transition-all shadow-2xl font-medium"
                            />
                        </div>
                    </motion.div>
                </div>
            </section>

            {/* Categories */}
            <section className="pb-24">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                        {[
                            { icon: Rocket, title: "Start", color: "text-blue-500", href: "/register" },
                            { icon: ShieldCheck, title: "Security", color: "text-emerald-500", href: "/security" },
                            { icon: Users, title: "P2P", color: "text-amber-500", href: "/marketplace" },
                            { icon: Zap, title: "Swap", color: "text-indigo-500", href: "/swap" },
                            { icon: Wallet, title: "Wallet", color: "text-purple-500", href: "#" },
                            { icon: Code, title: "API", color: "text-pink-500", href: "/api" },
                        ].map((cat, i) => (
                            <Link
                                key={i}
                                href={cat.href}
                                className="block"
                            >
                                <motion.div
                                    initial={{ opacity: 0, scale: 0.9 }}
                                    whileInView={{ opacity: 1, scale: 1 }}
                                    transition={{ delay: i * 0.05 }}
                                    className="glass p-6 rounded-[32px] border-border hover:border-primary/30 transition-all group text-center h-full"
                                >
                                    <div className={`w-12 h-12 bg-surface rounded-2xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform`}>
                                        <cat.icon className={`w-6 h-6 ${cat.color}`} />
                                    </div>
                                    <span className="font-black text-sm uppercase tracking-widest">{cat.title}</span>
                                </motion.div>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* FAQ & Contact Grid */}
            <section className="pb-24">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">

                        {/* FAQ Column */}
                        <div className="lg:col-span-7 space-y-8">
                            <div>
                                <h2 className="text-3xl font-black mb-2">Frequently Asked <span className="text-primary">Questions</span></h2>
                                <p className="text-text-dim font-medium">Find instant answers to common queries</p>
                            </div>

                            <div className="space-y-4">
                                {filteredFaqs.length > 0 ? (
                                    filteredFaqs.map((faq, i) => (
                                        <div
                                            key={i}
                                            className="glass rounded-3xl border-border overflow-hidden transition-all"
                                        >
                                            <button
                                                onClick={() => setOpenFaqIndex(openFaqIndex === i ? null : i)}
                                                className="w-full p-6 text-left flex items-center justify-between"
                                            >
                                                <div className="flex items-center space-x-4">
                                                    <span className="text-[10px] font-black bg-primary/10 text-primary px-2 py-1 rounded-lg uppercase">{faq.category}</span>
                                                    <span className="font-bold text-lg">{faq.question}</span>
                                                </div>
                                                <ChevronDown className={`w-5 h-5 transition-transform ${openFaqIndex === i ? "rotate-180 text-primary" : "text-text-dim"}`} />
                                            </button>
                                            <AnimatePresence>
                                                {openFaqIndex === i && (
                                                    <motion.div
                                                        initial={{ height: 0, opacity: 0 }}
                                                        animate={{ height: "auto", opacity: 1 }}
                                                        exit={{ height: 0, opacity: 0 }}
                                                        className="px-6 pb-6"
                                                    >
                                                        <p className="text-text-dim text-sm leading-relaxed font-medium pt-2 border-t border-white/5">
                                                            {faq.answer}
                                                        </p>
                                                    </motion.div>
                                                )}
                                            </AnimatePresence>
                                        </div>
                                    ))
                                ) : (
                                    <div className="glass p-12 rounded-[32px] text-center">
                                        <p className="text-text-dim font-bold">No results found for "{searchQuery}"</p>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Contact Form Column */}
                        <div className="lg:col-span-5">
                            <div className="glass p-8 rounded-[40px] border-border sticky top-32">
                                <div className="space-y-6 mb-8">
                                    <h3 className="text-2xl font-black">Open a <span className="gradient-text">Support Ticket</span></h3>
                                    <p className="text-text-dim text-sm font-medium">Can't find what you need? Our team is here to help you 24/7.</p>
                                </div>

                                {isSubmitted ? (
                                    <motion.div
                                        initial={{ opacity: 0, scale: 0.9 }}
                                        animate={{ opacity: 1, scale: 1 }}
                                        className="text-center py-12"
                                    >
                                        <div className="w-20 h-20 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-6">
                                            <CheckCircle2 className="w-10 h-10 text-primary" />
                                        </div>
                                        <h4 className="text-xl font-black mb-2">Ticket Submitted!</h4>
                                        <p className="text-text-dim text-sm mb-8 font-medium">We've received your inquiry. Check your email for updates.</p>
                                        <button
                                            onClick={() => setIsSubmitted(false)}
                                            className="w-full py-4 bg-surface border border-border rounded-2xl font-black uppercase tracking-widest hover:border-primary transition-all"
                                        >
                                            Submit Another
                                        </button>
                                    </motion.div>
                                ) : (
                                    <form onSubmit={handleTicketSubmit} className="space-y-4">
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-2">Full Name</label>
                                            <input
                                                required
                                                type="text"
                                                placeholder="John Doe"
                                                value={ticket.name}
                                                onChange={(e) => setTicket({ ...ticket, name: e.target.value })}
                                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-6 outline-none focus:border-primary transition-all font-medium"
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-2">Email Address</label>
                                            <input
                                                required
                                                type="email"
                                                placeholder="john@example.com"
                                                value={ticket.email}
                                                onChange={(e) => setTicket({ ...ticket, email: e.target.value })}
                                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-6 outline-none focus:border-primary transition-all font-medium"
                                            />
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-2">Subject</label>
                                            <select
                                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-6 outline-none focus:border-primary transition-all font-bold appearance-none"
                                                value={ticket.subject}
                                                onChange={(e) => setTicket({ ...ticket, subject: e.target.value })}
                                            >
                                                <option>Verification Issue</option>
                                                <option>Trading Dispute</option>
                                                <option>Technical Bug</option>
                                                <option>Feature Request</option>
                                                <option>Other</option>
                                            </select>
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] ml-2">Message</label>
                                            <textarea
                                                required
                                                rows={4}
                                                placeholder="Please describe your issue in detail..."
                                                value={ticket.message}
                                                onChange={(e) => setTicket({ ...ticket, message: e.target.value })}
                                                className="w-full bg-surface/50 border border-border rounded-2xl py-4 px-6 outline-none focus:border-primary transition-all font-medium resize-none"
                                            ></textarea>
                                        </div>

                                        <button
                                            disabled={isSubmitting}
                                            className="w-full mt-4 py-5 bg-primary text-white rounded-[24px] font-black text-lg uppercase tracking-widest transition-all shadow-xl shadow-primary/20 hover:scale-[1.02] active:scale-95 flex items-center justify-center space-x-3"
                                        >
                                            {isSubmitting ? (
                                                <RefreshCw className="w-6 h-6 animate-spin" />
                                            ) : (
                                                <>
                                                    <span>Submit Ticket</span>
                                                    <Send className="w-5 h-5" />
                                                </>
                                            )}
                                        </button>
                                    </form>
                                )}

                                {/* <div className="mt-8 pt-8 border-t border-white/5 grid grid-cols-2 gap-4">
                                    <button className="flex items-center justify-center space-x-2 bg-surface/30 p-4 rounded-2xl border border-border hover:border-primary transition-all group">
                                        <MessageCircle className="w-5 h-5 text-primary" />
                                        <span className="text-xs font-black">Live Chat</span>
                                    </button>
                                    <button className="flex items-center justify-center space-x-2 bg-surface/30 p-4 rounded-2xl border border-border hover:border-primary transition-all group">
                                        <Mail className="w-5 h-5 text-secondary" />
                                        <span className="text-xs font-black">Email Us</span>
                                    </button>
                                </div> */}
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Quick Links */}
            <section className="py-24 bg-surface/30">
                <div className="container mx-auto px-4">
                    <div className="text-center mb-12">
                        <h2 className="text-3xl font-black mb-4">Quick <span className="text-primary">Resources</span></h2>
                        <p className="text-text-dim font-medium">Access helpful guides and community support</p>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        {[
                            { title: "User Guide", desc: "Step-by-step masterclass on Cryplio", link: "#" },
                            { title: "Community Forum", desc: "Join 10k+ traders in discussion", link: "#" },
                            { title: "Status Page", desc: "Real-time system performance", link: "#" },
                        ].map((res, i) => (
                            <div key={i} className="glass p-8 rounded-[32px] border-border flex items-center justify-between group cursor-pointer hover:border-primary/50 transition-all">
                                <div>
                                    <h4 className="text-xl font-black mb-1 group-hover:text-primary transition-colors">{res.title}</h4>
                                    <p className="text-text-dim text-sm font-medium">{res.desc}</p>
                                </div>
                                <ExternalLink className="w-6 h-6 text-text-dim group-hover:text-primary transition-colors" />
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            <Footer />
        </main>
    );
};

export default SupportPage;
