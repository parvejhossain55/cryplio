"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Bell, CheckCircle2, AlertTriangle, Info, Clock, Check } from "lucide-react";
import { notificationService } from "@/services/notificationService";
import { useAuth } from "@/context/AuthContext";
import { wsService } from "@/services/websocketService";
import Link from "next/link";

const NotificationBell = () => {
    const { user } = useAuth();
    const [isOpen, setIsOpen] = useState(false);
    const [notifications, setNotifications] = useState<any[]>([]);
    const [unreadCount, setUnreadCount] = useState(0);

    // Connect to WebSocket and listen for real-time notifications
    useEffect(() => {
        if (user) {
            // Initial fetch
            fetchNotifications();
            
            // Connect WebSocket
            wsService.connect();
            
            // Listen for notification events
            const handleNotification = (data: any) => {
                if (data.type === 'notification') {
                    setNotifications(prev => [data.data, ...prev]);
                    setUnreadCount(prev => prev + 1);
                }
            };
            
            wsService.on('notification', handleNotification);
            
            return () => {
                wsService.off('notification', handleNotification);
            };
        }
    }, [user]);

    const fetchNotifications = async () => {
        if (!user) return;
        try {
            const data = await notificationService.getNotifications();
            setNotifications(data);
            setUnreadCount(data.filter((n: any) => !n.is_read).length);
        } catch (err) {
            // Silence errors for expected guest states or temporary service restarts
        }
    };

    const handleMarkRead = async (id: string) => {
        try {
            await notificationService.markRead(id);
            setNotifications(notifications.map(n => n.id === id ? { ...n, is_read: true } : n));
            setUnreadCount(prev => Math.max(0, prev - 1));
        } catch (err) {
            console.error(err);
        }
    };

    const getIcon = (type: string) => {
        switch (type) {
            case "trade_update": return <Clock className="w-4 h-4 text-primary" />;
            case "dispute_raised": return <AlertTriangle className="w-4 h-4 text-red-500" />;
            case "payment_received": return <CheckCircle2 className="w-4 h-4 text-accent" />;
            default: return <Info className="w-4 h-4 text-blue-500" />;
        }
    };

    return (
        <div className="relative">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="relative p-2.5 rounded-xl bg-white/5 border border-white/5 hover:border-white/10 transition-all group"
            >
                <Bell className={`w-5 h-5 transition-all ${unreadCount > 0 ? "text-primary animate-pulse" : "text-text-dim group-hover:text-white"}`} />
                {unreadCount > 0 && (
                    <span className="absolute top-2 right-2 w-2 h-2 bg-primary rounded-full border-2 border-background" />
                )}
            </button>

            <AnimatePresence>
                {isOpen && (
                    <>
                        <div className="fixed inset-0 z-40" onClick={() => setIsOpen(false)} />
                        <motion.div
                            initial={{ opacity: 0, y: 10, scale: 0.95 }}
                            animate={{ opacity: 1, y: 0, scale: 1 }}
                            exit={{ opacity: 0, y: 10, scale: 0.95 }}
                            className="absolute right-0 mt-4 w-80 z-50 glass border-white/10 rounded-[2rem] overflow-hidden shadow-2xl"
                        >
                            <div className="p-6 border-b border-white/5 flex items-center justify-between bg-white/5">
                                <h3 className="text-[10px] font-black uppercase tracking-widest text-white">Security Alerts</h3>
                                <span className="px-2 py-0.5 rounded bg-primary/10 text-primary text-[8px] font-black">{unreadCount} NEW</span>
                            </div>

                            <div className="max-h-[400px] overflow-y-auto scrollbar-hide">
                                {notifications.length === 0 ? (
                                    <div className="p-10 text-center space-y-2">
                                        <Bell className="w-8 h-8 text-text-dim/20 mx-auto" />
                                        <p className="text-[10px] font-black uppercase tracking-widest text-text-dim">No alerts found</p>
                                    </div>
                                ) : (
                                    notifications.map((n) => (
                                        <div
                                            key={n.id}
                                            className={`p-4 border-b border-white/5 last:border-0 hover:bg-white/5 transition-all cursor-pointer ${!n.is_read ? "bg-primary/5" : ""}`}
                                            onClick={() => handleMarkRead(n.id)}
                                        >
                                            <div className="flex gap-4">
                                                <div className="shrink-0 mt-1">
                                                    {getIcon(n.type)}
                                                </div>
                                                <div className="space-y-1">
                                                    <p className={`text-[11px] font-bold leading-tight ${!n.is_read ? "text-white" : "text-text-dim"}`}>
                                                        {n.message}
                                                    </p>
                                                    <p className="text-[8px] font-black uppercase tracking-widest text-text-dim/50">
                                                        {new Date(n.created_at).toLocaleTimeString()}
                                                    </p>
                                                </div>
                                                {!n.is_read && (
                                                    <div className="shrink-0">
                                                        <div className="w-1.5 h-1.5 bg-primary rounded-full" />
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>

                            {notifications.length > 0 && (
                                <Link
                                    href="/user/dashboard"
                                    className="block p-4 text-center text-[8px] font-black uppercase tracking-widest text-text-dim hover:text-white border-t border-white/5 bg-white/5 transition-all"
                                    onClick={() => setIsOpen(false)}
                                >
                                    View Command Center
                                </Link>
                            )}
                        </motion.div>
                    </>
                )}
            </AnimatePresence>
        </div>
    );
};

export default NotificationBell;
