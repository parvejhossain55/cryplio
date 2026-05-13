"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { wsService } from '@/services/websocketService';
import { toast } from 'sonner';
import { NotificationService } from "@/services/notificationService";

export interface Notification {
    id: string;
    title: string;
    message: string;
    type: "trade_update" | "system" | "payment_received" | "dispute_raised" | "trade_message";
    is_read: boolean;
    created_at: string;
    data?: any;
}

interface NotificationContextType {
    notifications: Notification[];
    unreadCount: number;
    addNotification: (notification: Partial<Notification>) => void;
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
                type: 'trade_update',
                title: 'Trade Update',
                message: `Trade ${data.trade_id} status changed to ${data.status}`,
                data: { trade_id: data.trade_id }
            });
        });

        wsService.on('dispute_created', (data) => {
            addNotification({
                type: 'dispute_raised',
                title: 'New Dispute',
                message: `Dispute created for trade ${data.trade_id}`,
                data: { trade_id: data.trade_id }
            });
        });

        wsService.on('message', (data) => {
            addNotification({
                type: 'system',
                title: data.title || 'Notification',
                message: data.message || 'You have a new notification'
            });
        });

        // Fetch existing notifications
        const fetchNotifications = async () => {
            try {
                const data = await NotificationService.getNotifications();
                setNotifications(data);
            } catch (error) {
                console.error("Failed to fetch notifications:", error);
            }
        };

        fetchNotifications();

        // Cleanup on unmount
        return () => {
            wsService.disconnect();
        };
    }, []);

    const addNotification = (notification: Partial<Notification>) => {
        const newNotification: Notification = {
            id: Date.now().toString() + Math.random().toString(36).substr(2, 9),
            title: notification.title || 'Notification',
            message: notification.message || '',
            type: notification.type || 'system',
            is_read: false,
            created_at: new Date().toISOString(),
            data: notification.data
        };

        setNotifications(prev => [newNotification, ...prev]);

        // Show toast notification
        const toastType = notification.type === 'trade_update' ? 'info' : 
                          notification.type === 'dispute_raised' ? 'error' : 'info';
        
        (toast as any)[toastType](newNotification.title, {
            description: newNotification.message,
        });
    };

    const markAsRead = async (id: string) => {
        try {
            await NotificationService.markAsRead(id);
            setNotifications(prev => 
                prev.map(notification => 
                    notification.id === id ? { ...notification, is_read: true } : notification
                )
            );
        } catch (error) {
            console.error("Failed to mark notification as read:", error);
        }
    };

    const markAllAsRead = () => {
        // Implementation for marking all as read would go here
        setNotifications(prev => 
            prev.map(notification => ({ ...notification, is_read: true }))
        );
    };

    const clearNotifications = () => {
        setNotifications([]);
    };

    const unreadCount = notifications.filter(n => !n.is_read).length;

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
