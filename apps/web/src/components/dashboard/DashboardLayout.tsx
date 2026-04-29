"use client";

import React, { useState } from "react";
import Sidebar from "./Sidebar";
import DashboardHeader from "./DashboardHeader";
import { motion, AnimatePresence } from "framer-motion";
import { X, Wallet } from "lucide-react";
import Link from "next/link";

interface DashboardLayoutProps {
    children: React.ReactNode;
    title: string;
    role: "user" | "merchant" | "admin";
}

const DashboardLayout = ({ children, title, role }: DashboardLayoutProps) => {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);

    return (
        <div className="min-h-screen bg-background flex text-white overflow-hidden">
            {/* Mobile Sidebar Overlay */}
            <AnimatePresence>
                {isSidebarOpen && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={() => setIsSidebarOpen(false)}
                        className="fixed inset-0 bg-background/60 backdrop-blur-sm z-40 md:hidden"
                    />
                )}
            </AnimatePresence>

            {/* Mobile Sidebar */}
            <motion.div
                className={`fixed top-0 left-0 bottom-0 w-72 bg-surface border-r border-border z-50 md:hidden`}
                initial={{ x: "-100%" }}
                animate={{ x: isSidebarOpen ? 0 : "-100%" }}
                transition={{ type: "spring", damping: 25, stiffness: 200 }}
            >
                <div className="p-6 flex items-center justify-between">
                    <Link href="/" className="flex items-center space-x-2">
                        <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                            <Wallet className="text-white w-5 h-5" />
                        </div>
                        <span className="text-xl font-black">Cryplio</span>
                    </Link>
                    <button onClick={() => setIsSidebarOpen(false)} className="p-2 hover:bg-surface-light rounded-lg">
                        <X className="w-5 h-5 text-white" />
                    </button>
                </div>
                <div className="mt-4 px-2">
                    {/* We can pass the Sidebar component here or its content */}
                    <Sidebar role={role} isMobile />
                </div>
            </motion.div>

            {/* Desktop Sidebar */}
            <Sidebar role={role} />

            {/* Main Content Area */}
            <div className="flex-1 flex flex-col min-h-screen md:pl-72 overflow-y-auto overflow-x-hidden relative">
                <DashboardHeader title={title} onMenuClick={() => setIsSidebarOpen(true)} />

                <main className="flex-1 p-6 md:p-10 relative">
                    {/* Subtle background glow */}
                    <div className="absolute top-0 right-0 w-[40%] h-[40%] bg-primary/5 rounded-full blur-[120px] -z-10 pointer-events-none" />
                    <div className="absolute bottom-0 left-0 w-[40%] h-[40%] bg-secondary/5 rounded-full blur-[120px] -z-10 pointer-events-none" />

                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.4 }}
                    >
                        {children}
                    </motion.div>
                </main>

                <footer className="px-10 py-6 border-t border-white/5 text-center md:text-left">
                    <p className="text-[10px] font-bold text-text-dim uppercase tracking-[0.2em]">
                        © 2026 Cryplio Infrastructure • Secure Gateway • V 2.4.0
                    </p>
                </footer>
            </div>
        </div>
    );
};

export default DashboardLayout;
