"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { wsService } from '@/services/websocketService';
import { toast } from 'sonner';

interface Notification {
    id: string;
    type: 'info' | 'success' | 'warning' | 'error';
    title: string;
    message: string;
    timestamp: Date;
    read: boolean;
    actionUrl?: string;
}

interface NotificationContextType {
    notifications: Notification[];
    unreadCount: number;
    addNotification: (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => void;
    markAsRead: (id: string) => void;
    markAllAsRead: () => void;
    clearNotifications: () => void;
    isConnected: boolean;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export const useNotifications = () => {
    const context = useContext(NotificationContext);
    if (!context) {
        throw new Error('useNotifications must be used within NotificationProvider');
    }
    return context;
};

interface NotificationProviderProps {
    children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        // Connect to WebSocket
        wsService.connect();

        // Listen for connection events
        wsService.on('connected', () => {
            setIsConnected(true);
        });

        wsService.on('disconnected', () => {
            setIsConnected(false);
        });

        // Listen for different notification types
        wsService.on('trade_update', (data) => {
            addNotification({
                type: 'info',
                title: 'Trade Update',
                message: `Trade ${data.trade_id} status changed to ${data.status}`,
                actionUrl: `/trades/${data.trade_id}`
            });
        });

        wsService.on('withdrawal_request', (data) => {
            addNotification({
                type: 'warning',
                title: 'New Withdrawal Request',
                message: `User ${data.username} requested withdrawal of ${data.amount} ${data.crypto_symbol}`,
                actionUrl: '/admin/withdrawals'
            });
        });

        wsService.on('dispute_created', (data) => {
            addNotification({
                type: 'error',
                title: 'New Dispute',
                message: `Dispute created for trade ${data.trade_id}`,
                actionUrl: `/admin/disputes`
            });
        });

        wsService.on('user_blocked', (data) => {
            addNotification({
                type: 'warning',
                title: 'User Blocked',
                message: `User ${data.username} has been blocked`,
                actionUrl: '/admin/users'
            });
        });

        wsService.on('market_update', (data) => {
            addNotification({
                type: 'info',
                title: 'Market Update',
                message: `${data.crypto_symbol} price: ${data.price} ${data.fiat_symbol}`,
                actionUrl: '/marketplace'
            });
        });

        wsService.on('message', (data) => {
            addNotification({
                type: 'info',
                title: data.title || 'Notification',
                message: data.message || 'You have a new notification'
            });
        });

        // Cleanup on unmount
        return () => {
            wsService.disconnect();
        };
    }, []);

    const addNotification = (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => {
        const newNotification: Notification = {
            ...notification,
            id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
            timestamp: new Date(),
            read: false
        };

        setNotifications(prev => [newNotification, ...prev]);

        // Show toast notification
        toast[notification.type](notification.title, {
            description: notification.message,
            action: notification.actionUrl && {
                label: 'View',
                onClick: () => {
                    window.location.href = notification.actionUrl!;
                }
            }
        });
    };

    const markAsRead = (id: string) => {
        setNotifications(prev => 
            prev.map(notification => 
                notification.id === id ? { ...notification, read: true } : notification
            )
        );
    };

    const markAllAsRead = () => {
        setNotifications(prev => 
            prev.map(notification => ({ ...notification, read: true }))
        );
    };

    const clearNotifications = () => {
        setNotifications([]);
    };

    const unreadCount = notifications.filter(n => !n.read).length;

    return (
        <NotificationContext.Provider
            value={{
                notifications,
                unreadCount,
                addNotification,
                markAsRead,
                markAllAsRead,
                clearNotifications,
                isConnected
            }}
        >
            {children}
        </NotificationContext.Provider>
    );
};
