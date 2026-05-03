"use client";

import React, { useState, useEffect } from "react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { motion } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { authService } from "@/services/authService";
import ConfirmModal from "@/components/ui/ConfirmModal";
import ProfileEdit from "@/components/settings/user/ProfileSettings";
import AlertsSettings from "@/components/settings/user/AlertsSettings";
import PaymentSettings from "@/components/settings/user/PaymentSettings";
import {
    User,
    Shield,
    Bell,
    CreditCard,
    UserCheck,
    ShieldCheck,
    Smartphone,
    LogOut,
    RefreshCw,
    Trash2,
    Loader2,
    AlertCircle,
    CheckCircle,
    Lock,
    Plus,
    ArrowUpRight,
    ArrowDownLeft,
    Building2,
    Globe,
    Volume2,
    Wallet
} from "lucide-react";

// QR Code component using qrcode library
const QRCodeSVG = ({ value, size = 200 }: { value: string; size?: number }) => {
    const [svgUrl, setSvgUrl] = useState<string>('');
    
    useEffect(() => {
        let mounted = true;
        
        const generateQR = async () => {
            try {
                // Dynamic import to avoid SSR issues
                const QRCode = (await import('qrcode')).default;
                if (!mounted) return;
                
                const svgString = await QRCode.toString(value, {
                    type: 'svg',
                    margin: 2,
                    width: size,
                    color: {
                        dark: '#000000',
                        light: '#ffffff'
                    }
                });
                
                if (mounted) {
                    const url = 'data:image/svg+xml;base64,' + btoa(svgString);
                    setSvgUrl(url);
                }
            } catch (err) {
                console.error('QR code generation failed:', err);
            }
        };

        generateQR();

        return () => {
            mounted = false;
        };
    }, [value, size]);

    if (!svgUrl) {
        return (
            <div
                style={{
                    width: size,
                    height: size,
                    backgroundColor: '#ffffff',
                    borderRadius: '8px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center'
                }}
            >
                <Loader2 className="w-6 h-6 animate-spin text-primary" />
            </div>
        );
    }

    return (
        <img
            src={svgUrl}
            alt="QR Code"
            width={size}
            height={size}
            style={{ borderRadius: '8px' }}
        />
    );
};

export default function SettingsPage() {
    const { user, refreshUser } = useAuth();
    const [activeTab, setActiveTab] = useState("profile");
    const [is2FAEnabled, setIs2FAEnabled] = useState(false);
    const [isSettingUp2FA, setIsSettingUp2FA] = useState(false);
    const [showQR, setShowQR] = useState(false);
    const [qrCode, setQrCode] = useState<string | null>(null);
    const [secret, setSecret] = useState<string | null>(null);
    const [setupCode, setSetupCode] = useState("");
    const [sessions, setSessions] = useState<any[]>([]);
    const [isLoadingSessions, setIsLoadingSessions] = useState(false);
    const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);
    const [modalOpen, setModalOpen] = useState(false);
    const [selectedTokenId, setSelectedTokenId] = useState<string | null>(null);

    const tabs = [
        { id: "profile", label: "Profile", icon: User },
        { id: "security", label: "Security", icon: Shield },
        { id: "notifications", label: "Alerts", icon: Bell },
        { id: "billing", label: "Payments", icon: CreditCard },
    ];

    useEffect(() => {
        if (user) {
            setIs2FAEnabled(user.twoFAEnabled);
            loadSessions();
        }
    }, [user]);

    const loadSessions = async () => {
        setIsLoadingSessions(true);
        try {
            const sessionList = await authService.getSessions();
            setSessions(sessionList);
        } catch (error) {
            console.error("Failed to load sessions:", error);
        } finally {
            setIsLoadingSessions(false);
        }
    };

    const handleEnable2FA = async () => {
        setIsSettingUp2FA(true);
        setMessage(null);
        try {
            const result = await authService.setup2FA();
            setSecret(result.secret);
            setQrCode(result.provisioning_uri);
            setShowQR(true);
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to setup 2FA" });
        } finally {
            setIsSettingUp2FA(false);
        }
    };

    const handleVerify2FA = async () => {
        if (!setupCode || setupCode.length !== 6) {
            setMessage({ type: "error", text: "Enter a valid 6-digit code" });
            return;
        }
        setIsSettingUp2FA(true);
        try {
            await authService.verify2FA(setupCode);
            setIs2FAEnabled(true);
            setShowQR(false);
            setSecret(null);
            setQrCode(null);
            setSetupCode("");
            setMessage({ type: "success", text: "Two-factor authentication enabled successfully!" });
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Verification failed" });
        } finally {
            setIsSettingUp2FA(false);
            refreshUser()
        }
    };

    const handleRevokeSession = async (tokenId: string) => {
        setSelectedTokenId(tokenId);
        setModalOpen(true);
    };

    const confirmRevokeSession = async () => {
        if (!selectedTokenId) return;
        try {
            await authService.revokeSession(selectedTokenId);
            setSessions(sessions.filter(s => s.token_id !== selectedTokenId));
            setMessage({ type: "success", text: "Session revoked successfully" });
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to revoke session" });
        } finally {
            setModalOpen(false);
            setSelectedTokenId(null);
        }
    };

    const handleDisable2FA = async () => {
        setSelectedTokenId("disable-2fa");
        setModalOpen(true);
    };

    const confirmDisable2FA = async () => {
        const password = prompt("Enter your password to disable 2FA:");
        if (!password) {
            setModalOpen(false);
            setSelectedTokenId(null);
            return;
        }
        setIsSettingUp2FA(true);
        try {
            await authService.disable2FA(password);
            setIs2FAEnabled(false);
            setMessage({ type: "success", text: "Two-factor authentication disabled" });
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to disable 2FA" });
        } finally {
            setIsSettingUp2FA(false);
            setModalOpen(false);
            setSelectedTokenId(null);
        }
    };

    const handleLogoutAll = async () => {
        setSelectedTokenId("logout-all");
        setModalOpen(true);
    };

    const confirmLogoutAll = async () => {
        try {
            for (const session of sessions) {
                await authService.revokeSession(session.token_id);
            }
            setSessions([]);
            setMessage({ type: "success", text: "Logged out from all devices successfully" });
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to logout from all devices" });
        } finally {
            setModalOpen(false);
            setSelectedTokenId(null);
        }
    };

    return (
        <DashboardLayout title="Settings" role="user">
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
                <div className="lg:col-span-1 space-y-2">
                    {/* <div className="p-6 bg-surface border border-white/10 rounded-[2.5rem] mb-6">
                        <div className="flex flex-col items-center">
                            <div className="relative group cursor-pointer">
                                <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center overflow-hidden transition-transform group-hover:scale-105">
                                    <User className="w-10 h-10 text-white/50" />
                                </div>
                                <div className="absolute bottom-0 right-0 w-8 h-8 bg-primary rounded-full border-4 border-surface flex items-center justify-center">
                                    <UserCheck className="w-3.5 h-3.5 text-white" />
                                </div>
                            </div>
                            <h3 className="mt-4 text-xl font-black text-white">
                                {user?.username || "User"}
                            </h3>
                            <p className="text-[10px] font-medium text-text-dim uppercase tracking-widest mt-1">
                                {user?.kycLevel ? `Level ${user.kycLevel} Trader` : "Unverified"}
                            </p>
                        </div>
                    </div> */}

                    <div className="bg-surface border border-white/10 rounded-[2.5rem] p-4 flex flex-col space-y-1">
                        {tabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`flex items-center px-6 py-4 rounded-2xl transition-all group ${
                                    activeTab === tab.id
                                        ? "bg-primary text-white shadow-xl shadow-primary/20"
                                        : "text-text-dim hover:bg-white/5 hover:text-white"
                                }`}
                            >
                                <tab.icon className={`w-5 h-5 mr-4 transition-colors ${activeTab === tab.id ? "text-white" : "text-text-dim group-hover:text-white"}`} />
                                <span className="text-sm font-black uppercase tracking-widest">{tab.label}</span>
                                {activeTab === tab.id && <ShieldCheck className="ml-auto w-4 h-4" />}
                            </button>
                        ))}
                    </div>
                </div>

                <div className="lg:col-span-3 space-y-8">
                    {activeTab === "security" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} className="space-y-8">
                            {message && (
                                <motion.div
                                    initial={{ opacity: 0, y: -10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    className={`p-4 rounded-2xl flex items-center space-x-3 ${
                                        message.type === "success"
                                            ? "bg-green-500/10 border border-green-500/30 text-green-400"
                                            : "bg-red-500/10 border border-red-500/30 text-red-400"
                                    }`}
                                >
                                    {message.type === "success" ? (
                                        <CheckCircle className="w-5 h-5 flex-shrink-0" />
                                    ) : (
                                        <AlertCircle className="w-5 h-5 flex-shrink-0" />
                                    )}
                                    <span className="text-sm font-medium">{message.text}</span>
                                </motion.div>
                            )}

                            <ConfirmModal
                                open={modalOpen}
                                title={selectedTokenId === "disable-2fa" ? "Disable 2FA?" : selectedTokenId === "logout-all" ? "Logout all devices?" : "Revoke this session?"}
                                description={selectedTokenId === "disable-2fa" ? "Two-factor authentication will be turned off. You'll need to set it up again to secure your account." : selectedTokenId === "logout-all" ? "This will log you out from all devices including this one." : "Are you sure you want to revoke this session? The device will be logged out immediately."}
                                onConfirm={() => {
                                    if (selectedTokenId === "disable-2fa") {
                                        confirmDisable2FA();
                                    } else if (selectedTokenId === "logout-all") {
                                        confirmLogoutAll();
                                    } else {
                                        confirmRevokeSession();
                                    }
                                }}
                                onClose={() => setModalOpen(false)}
                                confirmText="OK"
                                cancelText="Cancel"
                            />

                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <h3 className="text-xl font-black text-white mb-6 uppercase tracking-tight flex items-center">
                                    <Shield className="w-6 h-6 mr-3 text-primary" />
                                    Two-Factor Authentication
                                </h3>

                                {!is2FAEnabled && !showQR ? (
                                    <div className="flex items-center justify-between p-6 rounded-3xl bg-white/5 border border-white/5 hover:border-white/20 transition-all">
                                        <div className="flex items-start space-x-6">
                                            <div className="p-4 rounded-2xl bg-surface-light border border-white/10">
                                                <Smartphone className="w-6 h-6 text-white" />
                                            </div>
                                            <div>
                                                <h4 className="font-black text-white uppercase text-sm tracking-widest">Authenticator App</h4>
                                                <p className="text-xs text-text-dim font-medium mt-1">
                                                    Use Google Authenticator, Authy, or similar app for 2FA
                                                </p>
                                            </div>
                                        </div>
                                        <button
                                            onClick={handleEnable2FA}
                                            disabled={isSettingUp2FA}
                                            className="px-6 py-3 bg-primary text-white rounded-2xl text-xs font-black hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20"
                                        >
                                            {isSettingUp2FA ? "Preparing..." : "Enable 2FA"}
                                        </button>
                                    </div>
                                ) : showQR && secret ? (
                                    <div className="p-8 rounded-3xl bg-white/5 border border-white/10 space-y-6">
                                        <div className="text-center">
                                            <ShieldCheck className="w-12 h-12 text-primary mx-auto mb-4" />
                                            <h4 className="text-lg font-black text-white">Scan QR Code</h4>
                                            <p className="text-sm text-text-dim mt-2">
                                                Use your authenticator app to scan this QR code or enter the secret manually.
                                            </p>
                                        </div>

                                        <div className="flex justify-center p-6 bg-white rounded-xl">
                                            {qrCode && (
                                                <QRCodeSVG
                                                    value={qrCode}
                                                    size={200}
                                                />
                                            )}
                                        </div>

                                        <div className="space-y-4">
                                            <div>
                                                <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block mb-2">
                                                    Or enter this secret manually:
                                                </label>
                                                <code className="block p-4 bg-surface rounded-2xl text-sm font-mono break-all select-all">
                                                    {secret}
                                                </code>
                                            </div>

                                            <div className="space-y-2">
                                                <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block">
                                                    Enter 6-digit code from app
                                                </label>
                                                <div className="flex space-x-3">
                                                    <input
                                                        type="text"
                                                        inputMode="numeric"
                                                        maxLength={6}
                                                        placeholder="000000"
                                                        value={setupCode}
                                                        onChange={(e) => setSetupCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                                                        className="flex-1 bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-center text-lg font-mono outline-none focus:border-primary transition-all"
                                                    />
                                                    <button
                                                        onClick={handleVerify2FA}
                                                        disabled={isSettingUp2FA || setupCode.length !== 6}
                                                        className="px-6 bg-primary text-white rounded-2xl font-black disabled:opacity-50"
                                                    >
                                                        {isSettingUp2FA ? <Loader2 className="w-5 h-5 animate-spin" /> : "Verify"}
                                                    </button>
                                                </div>
                                            </div>
                                        </div>

                                        <button
                                            onClick={() => {
                                                setShowQR(false);
                                                setSecret(null);
                                                setQrCode(null);
                                            }}
                                            className="w-full text-center text-xs text-text-dim hover:text-white mt-4"
                                        >
                                            Cancel setup
                                        </button>
                                    </div>
                                ) : is2FAEnabled ? (
                                    <div className="flex items-center justify-between p-6 rounded-3xl bg-green-500/10 border border-green-500/30">
                                        <div className="flex items-center space-x-4">
                                            <div className="p-4 rounded-2xl bg-green-500/20">
                                                <ShieldCheck className="w-6 h-6 text-green-500" />
                                            </div>
                                            <div>
                                                <h4 className="font-bold text-white">Two-Factor Authentication Enabled</h4>
                                                <p className="text-xs text-text-dim mt-1">Your account is protected with 2FA</p>
                                            </div>
                                        </div>
                                        <button
                                            onClick={handleDisable2FA}
                                            className="px-6 py-3 bg-red-500/10 text-red-400 border border-red-500/30 rounded-2xl text-xs font-black hover:bg-red-500/20 transition-all"
                                        >
                                            Disable
                                        </button>
                                    </div>
                                ) : null}
                            </div>

                            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-10">
                                <div className="flex items-center justify-between mb-8">
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight flex items-center">
                                        <RefreshCw className="w-6 h-6 mr-3 text-primary" />
                                        Active Sessions
                                    </h3>
                                    <button
                                        onClick={handleLogoutAll}
                                        className="text-[8px] font-black text-red-400 uppercase tracking-wider hover:underline"
                                    >
                                        Logout All Devices
                                    </button>
                                </div>

                                {isLoadingSessions ? (
                                    <div className="flex justify-center py-10">
                                        <Loader2 className="w-8 h-8 animate-spin text-primary" />
                                    </div>
                                ) : sessions.length === 0 ? (
                                    <p className="text-text-dim text-center py-10">No active sessions found.</p>
                                ) : (
                                    <div className="space-y-4">
                                        {sessions.map((session) => (
                                            <div
                                                key={session.id || session.token_id}
                                                className="flex items-center justify-between p-6 rounded-3xl bg-white/5 border border-white/5 group hover:border-white/20 transition-all"
                                            >
                                                <div className="flex items-center space-x-4">
                                                    <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                                                        <Smartphone className="w-5 h-5 text-primary" />
                                                    </div>
                                                    <div>
                                                        <p className="font-bold text-white text-sm">
                                                            {session.device_type || "Unknown device"}
                                                        </p>
                                                        <p className="text-xs text-text-dim">
                                                            IP: {session.ip_address || "N/A"} • Last used: {new Date(session.last_used_at).toLocaleString()}
                                                        </p>
                                                    </div>
                                                </div>
                                                <button
                                                    onClick={() => handleRevokeSession(session.token_id)}
                                                    className="p-2 text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                                                    title="Revoke session"
                                                >
                                                    <Trash2 className="w-5 h-5" />
                                                </button>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>
                        </motion.div>
                    )}

                    {activeTab === "profile" && user && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
                            <ProfileEdit user={user} />
                        </motion.div>
                    )}
                    
                    {activeTab === "notifications" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
                            <AlertsSettings />
                        </motion.div>
                    )}

                    {activeTab === "billing" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
                            <PaymentSettings />
                        </motion.div>
                    )}
                    
                    {activeTab === "preferences" && (
                        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}>
                        </motion.div>
                    )}
                </div>
            </div>
        </DashboardLayout>
    );
}

