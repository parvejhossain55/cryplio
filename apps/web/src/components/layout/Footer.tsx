import React from "react";
import Link from "next/link";
import { Wallet, Send, Code2, Briefcase, MessageCircle } from "lucide-react";

const Footer = () => {
    const footerLinks = [
        {
            title: "Product",
            links: [
                { name: "P2P Trading", href: "/marketplace" },
                { name: "Invest", href: "#" },
                { name: "Swap", href: "/swap" },
                { name: "Fees", href: "#" },
            ],
        },
        {
            title: "Account",
            links: [
                { name: "Login", href: "/login" },
                { name: "Register", href: "/register" },
                { name: "Forgot Password", href: "/forgot-password" },
                { name: "Security", href: "/security" },
            ],
        },
        {
            title: "Support",
            links: [
                { name: "Help Center", href: "/support" },
                { name: "API Docs", href: "/api" },
                { name: "Contact Us", href: "/support" },
                { name: "Community", href: "/community" },
            ],
        },
    ];

    return (
        <footer className="bg-surface border-t border-border pt-16 pb-8">
            <div className="container mx-auto px-4 md:px-6">
                <div className="grid grid-cols-2 md:grid-cols-5 gap-10 mb-16">
                    {/* Brand Column */}
                    <div className="col-span-2">
                        <Link href="/" className="flex items-center space-x-2 mb-6">
                            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                                <Wallet className="text-white w-5 h-5" />
                            </div>
                            <span className="text-xl font-bold">
                                Cryp<span className="gradient-text">lio</span>
                            </span>
                        </Link>
                        <p className="text-text-dim text-sm max-w-xs mb-8 leading-relaxed">
                            The most trusted P2P cryptocurrency exchange platform. Buy, sell, and trade over 50+ cryptocurrencies with local payment methods.
                        </p>
                        <div className="flex items-center space-x-5">
                            <Link href="#" className="text-text-dim hover:text-primary transition-colors">
                                <Send className="w-5 h-5" />
                            </Link>
                            <Link href="#" className="text-text-dim hover:text-primary transition-colors">
                                <Code2 className="w-5 h-5" />
                            </Link>
                            <Link href="#" className="text-text-dim hover:text-primary transition-colors">
                                <Briefcase className="w-5 h-5" />
                            </Link>
                            <Link href="#" className="text-text-dim hover:text-primary transition-colors">
                                <MessageCircle className="w-5 h-5" />
                            </Link>
                        </div>
                    </div>

                    {/* Links Columns */}
                    {footerLinks.map((group) => (
                        <div key={group.title} className="col-span-1">
                            <h3 className="text-white font-semibold mb-6">{group.title}</h3>
                            <ul className="space-y-4 text-sm">
                                {group.links.map((link) => (
                                    <li key={link.name}>
                                        <Link
                                            href={link.href}
                                            className="text-text-dim hover:text-white transition-colors"
                                        >
                                            {link.name}
                                        </Link>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>

                <hr className="border-border mb-8" />

                <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
                    <p className="text-text-dim text-xs">
                        © {new Date().getFullYear()} Cryplio Inc. All rights reserved.
                    </p>
                    <div className="flex items-center space-x-6 text-xs text-text-dim">
                        <Link href="#" className="hover:text-white transition-colors">
                            Privacy Policy
                        </Link>
                        <Link href="#" className="hover:text-white transition-colors">
                            Terms of Service
                        </Link>
                        <Link href="#" className="hover:text-white transition-colors">
                            Cookie Policy
                        </Link>
                    </div>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
