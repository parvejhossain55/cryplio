"use client";

import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
    Bell, 
    BellRing, 
    Check, 
    CheckCheck, 
    X, 
    ExternalLink,
    Wifi,
    WifiOff
} from 'lucide-react';
import { useNotifications } from '@/context/NotificationContext';
import { formatRelativeTime } from '@/lib/utils';

export const NotificationDropdown: React.FC = () => {
    const { 
        notifications, 
        unreadCount, 
        markAsRead, 
        markAllAsRead, 
        clearNotifications,
        isConnected 
    } = useNotifications();
    
    const [isOpen, setIsOpen] = useState(false);
    const dropdownRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleNotificationClick = (notification: any) => {
        markAsRead(notification.id);
        setIsOpen(false);
        
        if (notification.actionUrl) {
            window.location.href = notification.actionUrl;
        }
    };

    const getNotificationIcon = (type: string) => {
        switch (type) {
            case 'success': return Check;
            case 'warning': return BellRing;
            case 'error': return X;
            default: return Bell;
        }
    };

    const getNotificationColor = (type: string) => {
        switch (type) {
            case 'success': return 'text-green-400 bg-green-400/10';
            case 'warning': return 'text-yellow-400 bg-yellow-400/10';
            case 'error': return 'text-red-400 bg-red-400/10';
            default: return 'text-blue-400 bg-blue-400/10';
        }
    };

    return (
        <div className="relative" ref={dropdownRef}>
            {/* Notification Bell */}
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="relative p-2 text-text-dim hover:text-white transition-colors"
            >
                {isConnected ? (
                    <Wifi className="w-5 h-5" />
                ) : (
                    <WifiOff className="w-5 h-5 text-red-400" />
                )}
                
                {unreadCount > 0 && (
                    <span className="absolute top-1 right-1 h-2 w-2 bg-danger-500 rounded-full animate-pulse" />
                )}
                
                {unreadCount > 0 && (
                    <span className="absolute -top-1 -right-1 h-5 w-5 bg-danger-500 text-white text-xs rounded-full flex items-center justify-center font-bold">
                        {unreadCount > 9 ? '9+' : unreadCount}
                    </span>
                )}
            </button>

            {/* Dropdown */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -10, scale: 0.95 }}
                        animate={{ opacity: 1, y: 0, scale: 1 }}
                        exit={{ opacity: 0, y: -10, scale: 0.95 }}
                        className="absolute right-0 mt-2 w-96 bg-surface border border-white/10 rounded-2xl shadow-2xl z-50 max-h-96 overflow-hidden"
                    >
                        {/* Header */}
                        <div className="p-4 border-b border-white/10">
                            <div className="flex items-center justify-between">
                                <h3 className="text-sm font-bold text-white">Notifications</h3>
                                <div className="flex items-center space-x-2">
                                    {unreadCount > 0 && (
                                        <button
                                            onClick={markAllAsRead}
                                            className="text-xs text-text-dim hover:text-white transition-colors"
                                        >
                                            Mark all read
                                        </button>
                                    )}
                                    <button
                                        onClick={clearNotifications}
                                        className="text-xs text-text-dim hover:text-white transition-colors"
                                    >
                                        Clear
                                    </button>
                                </div>
                            </div>
                            
                            <div className="flex items-center space-x-2 mt-2">
                                <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-400' : 'bg-red-400'}`} />
                                <span className="text-xs text-text-dim">
                                    {isConnected ? 'Connected' : 'Disconnected'}
                                </span>
                            </div>
                        </div>

                        {/* Notifications List */}
                        <div className="max-h-80 overflow-y-auto">
                            {notifications.length === 0 ? (
                                <div className="p-8 text-center">
                                    <Bell className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                    <p className="text-sm text-text-dim">No notifications</p>
                                </div>
                            ) : (
                                <div className="divide-y divide-white/5">
                                    {notifications.map((notification) => {
                                        const Icon = getNotificationIcon(notification.type);
                                        return (
                                            <motion.div
                                                key={notification.id}
                                                initial={{ opacity: 0, x: -20 }}
                                                animate={{ opacity: 1, x: 0 }}
                                                className={`p-4 hover:bg-white/5 cursor-pointer transition-colors ${
                                                    !notification.read ? 'bg-white/5' : ''
                                                }`}
                                                onClick={() => handleNotificationClick(notification)}
                                            >
                                                <div className="flex items-start space-x-3">
                                                    <div className={`p-2 rounded-lg ${getNotificationColor(notification.type)}`}>
                                                        <Icon className="w-4 h-4" />
                                                    </div>
                                                    
                                                    <div className="flex-1 min-w-0">
                                                        <div className="flex items-start justify-between">
                                                            <div className="flex-1">
                                                                <p className="text-sm font-medium text-white">
                                                                    {notification.title}
                                                                </p>
                                                                <p className="text-xs text-text-dim mt-1">
                                                                    {notification.message}
                                                                </p>
                                                                <p className="text-xs text-text-dim mt-2">
                                                                    {formatRelativeTime(notification.timestamp)}
                                                                </p>
                                                            </div>
                                                            
                                                            {!notification.read && (
                                                                <div className="w-2 h-2 bg-primary-500 rounded-full mt-1" />
                                                            )}
                                                        </div>
                                                        
                                                        {notification.actionUrl && (
                                                            <div className="flex items-center space-x-1 mt-2 text-primary-400">
                                                                <ExternalLink className="w-3 h-3" />
                                                                <span className="text-xs">View</span>
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            </motion.div>
                                        );
                                    })}
                                </div>
                            )}
                        </div>

                        {/* Footer */}
                        {notifications.length > 0 && (
                            <div className="p-3 border-t border-white/10 text-center">
                                <button
                                    onClick={() => setIsOpen(false)}
                                    className="text-xs text-text-dim hover:text-white transition-colors"
                                >
                                    Close
                                </button>
                            </div>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
};
