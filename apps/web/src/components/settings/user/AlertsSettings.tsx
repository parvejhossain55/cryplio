"use client";

import React, { useState, useEffect } from "react";
import { Bell, BellOff, Mail, Smartphone, Globe, Volume2, VolumeX, Loader2, CheckCircle, AlertCircle, Save } from "lucide-react";
import { motion } from "framer-motion";
import { authService } from "@/services/authService";
import { toast } from "sonner";

interface NotificationChannel {
    id: string;
    name: string;
    description: string;
    icon: React.ElementType;
    enabled: boolean;
}

interface NotificationCategory {
    id: string;
    name: string;
    description: string;
    icon: React.ElementType;
    channels: NotificationChannel[];
}

const AlertsSettings = () => {
    const [isSaving, setIsSaving] = useState(false);

    // Initialize from localStorage or defaults
    const getInitialChannels = (): NotificationChannel[] => [
        { id: "email_security", name: "Email Security Alerts", description: "Login alerts, password changes, security notifications", icon: Mail, enabled: true },
        { id: "email_transactions", name: "Email Transaction Updates", description: "Trade confirmations, deposits, withdrawals", icon: Mail, enabled: true },
        { id: "email_marketing", name: "Email Marketing", description: "Newsletter, promotions, new features", icon: Globe, enabled: false },
        { id: "push_security", name: "Push Security Alerts", description: "Real-time security notifications on device", icon: Bell, enabled: true },
        { id: "push_trades", name: "Push Trade Alerts", description: "Price alerts, order fills, market updates", icon: Bell, enabled: false },
        { id: "sms_critical", name: "SMS Critical Alerts", description: "Urgent security issues only", icon: Smartphone, enabled: false },
    ];

    const [channels, setChannels] = useState<NotificationChannel[]>(getInitialChannels());

    useEffect(() => {
        const saved = localStorage.getItem("cryplio_notification_preferences");
        if (saved) {
            try {
                const parsed = JSON.parse(saved);
                setChannels(prev => prev.map(ch => ({ ...ch, enabled: parsed[ch.id] ?? ch.enabled })));
            } catch (e) {
                console.error("Failed to load notifications prefs", e);
            }
        }
    }, []);

    const toggleChannel = (id: string) => {
        setChannels(prev => prev.map(ch =>
            ch.id === id ? { ...ch, enabled: !ch.enabled } : ch
        ));
    };

    const handleSave = async () => {
        setIsSaving(true);
        try {
            const prefs = channels.reduce((acc, ch) => ({ ...acc, [ch.id]: ch.enabled }), {});
            localStorage.setItem("cryplio_notification_preferences", JSON.stringify(prefs));
            // In future: POST to /api/users/notifications/preferences
            await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API
            toast.success("Notification preferences saved successfully");
        } catch (error: any) {
            toast.error(error.message || "Failed to save preferences");
        } finally {
            setIsSaving(false);
        }
    };

    // Group by notification type
    const emailChannels = channels.filter(ch => ch.icon === Mail);
    const pushChannels = channels.filter(ch => ch.icon === Bell);
    const otherChannels = channels.filter(ch => ch.icon !== Mail && ch.icon !== Bell);

    const ChannelRow = ({ channel }: { channel: NotificationChannel }) => (
        <div className="flex items-center justify-between p-4 rounded-2xl bg-white/5 border border-white/5 hover:border-white/10 transition-all">
            <div className="flex items-center gap-4">
                <div className={`p-2.5 rounded-lg ${channel.enabled ? "bg-primary/20 text-primary" : "bg-white/10 text-text-dim"}`}>
                    {(() => { const Icon = channel.icon; return <Icon className="w-5 h-5" />; })()}
                </div>
                <div>
                    <h4 className="font-bold text-white text-sm">{channel.name}</h4>
                    <p className="text-[10px] text-text-dim mt-0.5 leading-relaxed">{channel.description}</p>
                </div>
            </div>
            <button
                onClick={() => toggleChannel(channel.id)}
                className={`relative w-11 h-6 rounded-full transition-colors duration-200 ${channel.enabled ? "bg-accent" : "bg-white/10"}`}
                aria-label={`Toggle ${channel.name}`}
            >
                <span className={`absolute top-1 w-4 h-4 rounded-full bg-white shadow transition-transform duration-200 ${channel.enabled ? "right-1" : "left-1"}`} />
            </button>
        </div>
    );

    return (
        <div className="space-y-8">
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <h3 className="text-xl font-black text-white mb-2 uppercase tracking-tight flex items-center">
                            <Bell className="w-6 h-6 mr-3 text-primary" />
                            Notification Channels
                        </h3>
                        <p className="text-xs text-text-dim">Choose how you receive notifications</p>
                    </div>
                </div>

                <div className="space-y-6">
                    {/* Email Section */}
                    <div>
                        <h4 className="text-sm font-black text-white uppercase tracking-wider mb-4 flex items-center">
                            <Mail className="w-4 h-4 mr-2 text-primary" />
                            Email Notifications
                        </h4>
                        <div className="space-y-3">
                            {emailChannels.map(channel => (
                                <ChannelRow key={channel.id} channel={channel} />
                            ))}
                        </div>
                    </div>

                    {/* Push Notifications */}
                    <div>
                        <h4 className="text-sm font-black text-white uppercase tracking-wider mb-4 flex items-center">
                            <Bell className="w-4 h-4 mr-2 text-primary" />
                            Push Notifications
                        </h4>
                        <div className="space-y-3">
                            {pushChannels.map(channel => (
                                <ChannelRow key={channel.id} channel={channel} />
                            ))}
                        </div>
                    </div>

                </div>
            </div>


            {/* Save Button */}
            <div className="flex justify-end">
                <button
                    onClick={handleSave}
                    disabled={isSaving}
                    className="flex items-center px-8 py-4 bg-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {isSaving ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Save className="w-4 h-4 mr-2" />}
                    Save Preferences
                </button>
            </div>
        </div>
    );
};

export default AlertsSettings;
