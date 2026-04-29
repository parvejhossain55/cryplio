"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { Menu, X, ChevronDown, Globe, Wallet } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const Navbar = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [scrolled, setScrolled] = useState(false);

    useEffect(() => {
        const handleScroll = () => {
            setScrolled(window.scrollY > 20);
        };
        window.addEventListener("scroll", handleScroll);
        return () => window.removeEventListener("scroll", handleScroll);
    }, []);

    const navLinks = [
        { name: "Marketplace", href: "/marketplace" },
        { name: "Features", href: "/#features" },
        { name: "Swap", href: "/swap" },
        { name: "Support", href: "/support" },
    ];

    return (
        <nav
            className={`fixed top-0 w-full z-50 transition-all duration-300 ${scrolled ? "glass py-3" : "bg-transparent py-5"
                }`}
        >
            <div className="container mx-auto px-4 md:px-6">
                <div className="flex items-center justify-between">
                    {/* Logo */}
                    <Link href="/" className="flex items-center space-x-2">
                        <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20">
                            <Wallet className="text-white w-6 h-6" />
                        </div>
                        <span className="text-2xl font-bold tracking-tight">
                            Cryp<span className="gradient-text">lio</span>
                        </span>
                    </Link>

                    {/* Desktop Nav */}
                    <div className="hidden md:flex items-center space-x-8">
                        {navLinks.map((link) => (
                            <Link
                                key={link.name}
                                href={link.href}
                                className="text-sm font-medium text-text-dim hover:text-white transition-colors"
                            >
                                {link.name}
                            </Link>
                        ))}
                    </div>

                    {/* Desktop Actions */}
                    <div className="hidden md:flex items-center space-x-4">
                        <button className="flex items-center space-x-1 text-sm font-medium text-text-dim hover:text-white transition-colors px-3 py-2">
                            <Globe className="w-4 h-4" />
                            <span>EN</span>
                            <ChevronDown className="w-4 h-4" />
                        </button>
                        <Link
                            href="/login"
                            className="text-sm font-medium hover:text-primary transition-colors"
                        >
                            Log In
                        </Link>
                        <Link
                            href="/register"
                            className="bg-primary hover:bg-primary-dark text-white px-6 py-2.5 rounded-full text-sm font-semibold transition-all shadow-lg shadow-primary/20 hover:scale-105 active:scale-95"
                        >
                            Sign Up
                        </Link>
                    </div>

                    {/* Mobile Menu Toggle */}
                    <button
                        className="md:hidden p-2 text-text hover:bg-surface-light rounded-lg transition-colors"
                        onClick={() => setIsOpen(!isOpen)}
                    >
                        {isOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                    </button>
                </div>
            </div>

            {/* Mobile Nav */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: "auto" }}
                        exit={{ opacity: 0, height: 0 }}
                        className="md:hidden glass border-t border-border mt-3"
                    >
                        <div className="container mx-auto px-4 py-6 flex flex-col space-y-4">
                            {navLinks.map((link) => (
                                <Link
                                    key={link.name}
                                    href={link.href}
                                    className="text-lg font-medium text-text-dim hover:text-white transition-colors py-2"
                                    onClick={() => setIsOpen(false)}
                                >
                                    {link.name}
                                </Link>
                            ))}
                            <hr className="border-border" />
                            <div className="flex flex-col space-y-4 pt-2">
                                <Link
                                    href="/login"
                                    className="w-full text-center py-3 text-lg font-medium border border-border rounded-xl hover:bg-surface-light transition-colors"
                                    onClick={() => setIsOpen(false)}
                                >
                                    Log In
                                </Link>
                                <Link
                                    href="/register"
                                    className="w-full text-center py-3 text-lg font-medium bg-primary text-white rounded-xl shadow-lg shadow-primary/20"
                                    onClick={() => setIsOpen(false)}
                                >
                                    Sign Up
                                </Link>
                            </div>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </nav>
    );
};

export default Navbar;
