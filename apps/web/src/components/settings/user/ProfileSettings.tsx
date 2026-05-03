"use client";

import React, { useState, useEffect } from "react";
import { User, Mail, Edit2, Save, X, CheckCircle, AlertCircle, Shield, Loader2 } from "lucide-react";
import { motion } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import type { User as UserType } from "@/context/AuthContext";
import { authService } from "@/services/authService";
import ConfirmModal from "@/components/ui/ConfirmModal";

interface ProfileEditProps {
    user: UserType;
}

const ProfileEdit = ({ user }: ProfileEditProps) => {
    const { refreshUser } = useAuth();
    const [isEditing, setIsEditing] = useState(false);
    const [username, setUsername] = useState(user.username);
    const [bio, setBio] = useState(user.bio || "");
    const [isSaving, setIsSaving] = useState(false);
    const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);
    const [showCancelModal, setShowCancelModal] = useState(false);
    const [hasChanges, setHasChanges] = useState(false);

    useEffect(() => {
        setUsername(user.username);
        setBio(user.bio || "");
        setIsEditing(false);
        setMessage(null);
    }, [user]);

    useEffect(() => {
        setHasChanges(username !== user.username || bio !== (user.bio || ""));
    }, [username, bio, user]);

    const handleSave = async () => {
        if (!username.trim()) {
            setMessage({ type: "error", text: "Username is required" });
            return;
        }
        setIsSaving(true);
        setMessage(null);
        try {
            await authService.updateCurrentUser({ username: username.trim(), bio: bio.trim() });
            await refreshUser();
            setMessage({ type: "success", text: "Profile updated successfully" });
            setIsEditing(false);
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to update profile" });
        } finally {
            setIsSaving(false);
        }
    };

    const handleCancel = () => {
        if (hasChanges) {
            setShowCancelModal(true);
        } else {
            setIsEditing(false);
            setUsername(user.username);
            setBio(user.bio || "");
        }
    };

    const confirmCancel = () => {
        setShowCancelModal(false);
        setIsEditing(false);
        setUsername(user.username);
        setBio(user.bio || "");
        setMessage(null);
    };

    const handleResendVerification = async () => {
        try {
            await authService.requestEmailVerification(user.id);
            setMessage({ type: "success", text: "Verification email sent" });
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to send verification" });
        }
    };

    return (
        <div className="space-y-8">
            {/* Profile Header Card */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <div className="flex flex-col md:flex-row items-center md:items-start gap-8">
                    {/* Avatar */}
                    <div className="relative group">
                        <div className="w-28 h-28 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center overflow-hidden transition-transform group-hover:scale-105">
                            <User className="w-12 h-12 text-white/50" />
                        </div>
                        {isEditing && (
                            <div className="absolute inset-0 mx-auto w-28 h-28 rounded-full bg-black/60 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                <span className="text-[10px] font-black uppercase tracking-widest text-white">Change</span>
                            </div>
                        )}
                    </div>

                    {/* User Info */}
                    <div className="flex-1 text-center md:text-left space-y-4">
                        <div>
                            {isEditing ? (
                                <input
                                    type="text"
                                    value={username}
                                    onChange={(e) => setUsername(e.target.value)}
                                    className="w-full max-w-xs bg-white/5 border border-white/10 py-3 px-5 rounded-2xl text-xl font-black text-white outline-none focus:border-primary transition-all"
                                    placeholder="Username"
                                />
                            ) : (
                                <h2 className="text-2xl font-black text-white tracking-tight">
                                    {user.username}
                                </h2>
                            )}
                            <p className="text-sm text-text-dim mt-1 flex items-center justify-center md:justify-start gap-2">
                                <Mail className="w-4 h-4" />
                                {user.email}
                            </p>
                        </div>

                        <div className="flex flex-wrap items-center justify-center md:justify-start gap-3">
                            {/* Email Verification Badge */}
                            {user.emailVerified ? (
                                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-accent/10 text-accent border border-accent/20 text-[10px] font-black uppercase tracking-widest">
                                    <CheckCircle className="w-3.5 h-3.5 mr-1.5" />
                                    Verified
                                </span>
                            ) : (
                                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-yellow-500/10 text-yellow-500 border border-yellow-500/20 text-[10px] font-black uppercase tracking-widest">
                                    <AlertCircle className="w-3.5 h-3.5 mr-1.5" />
                                    Unverified
                                </span>
                            )}

                            {/* 2FA Badge */}
                            {user.twoFAEnabled ? (
                                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-accent/10 text-accent border border-accent/20 text-[10px] font-black uppercase tracking-widest">
                                    <Shield className="w-3.5 h-3.5 mr-1.5" />
                                    2FA On
                                </span>
                            ) : (
                                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-yellow-500/10 text-yellow-500 border border-yellow-500/20 text-[10px] font-black uppercase tracking-widest">
                                    <Shield className="w-3.5 h-3.5 mr-1.5" />
                                    2FA Off
                                </span>
                            )}

                            {/* KYC Level */}
                            <span className={`inline-flex items-center px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border ${
                                user.kycLevel >= 2 ? "bg-accent/10 text-accent border-accent/20" :
                                user.kycLevel === 1 ? "bg-primary/10 text-primary border-primary/20" :
                                "bg-yellow-500/10 text-yellow-500 border-yellow-500/20"
                            }`}>
                                Level {user.kycLevel} Trader
                            </span>
                        </div>
                    </div>

                    {/* Edit/Save Buttons */}
                    <div className="flex flex-col gap-3">
                        {!isEditing ? (
                            <button
                                onClick={() => setIsEditing(true)}
                                className="flex items-center justify-center px-6 py-3 bg-white/5 border border-white/10 text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/10 transition-all"
                            >
                                <Edit2 className="w-4 h-4 mr-2" />
                                Edit Profile
                            </button>
                        ) : (
                            <>
                                <button
                                    onClick={handleSave}
                                    disabled={isSaving || !hasChanges}
                                    className="flex items-center justify-center px-6 py-3 bg-primary border border-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    {isSaving ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Save className="w-4 h-4 mr-2" />}
                                    Save
                                </button>
                                <button
                                    onClick={handleCancel}
                                    disabled={isSaving}
                                    className="flex items-center justify-center px-6 py-3 bg-white/5 border border-white/10 text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/10 transition-all"
                                >
                                    <X className="w-4 h-4 mr-2" />
                                    Cancel
                                </button>
                            </>
                        )}
                    </div>
                </div>
            </div>

            {/* Bio Section */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <h3 className="text-xl font-black text-white mb-6 uppercase tracking-tight flex items-center">
                    <User className="w-5 h-5 mr-3 text-primary" />
                    About Me
                </h3>

                {isEditing ? (
                    <div className="space-y-4">
                        <textarea
                            value={bio}
                            onChange={(e) => setBio(e.target.value)}
                            rows={4}
                            maxLength={300}
                            placeholder="Tell others about yourself (max 300 characters)"
                            className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all resize-none"
                        />
                        <p className="text-[10px] text-text-dim text-right">
                            {bio.length}/300 characters
                        </p>
                    </div>
                ) : (
                    <div>
                        {user.bio ? (
                            <p className="text-sm text-text-dim leading-relaxed whitespace-pre-wrap">
                                {user.bio}
                            </p>
                        ) : (
                            <p className="text-sm text-text-dim italic">
                                No bio added yet. Click Edit to add a description.
                            </p>
                        )}
                    </div>
                )}
            </div>

            {/* Email Verification Call-to-Action (if not verified) */}
            {!user.emailVerified && (
                <div className="bg-surface border border-yellow-500/20 rounded-[2.5rem] p-8 md:p-10 border-l-4 border-l-yellow-500">
                    <div className="flex flex-col md:flex-row items-center md:items-center justify-between gap-4">
                        <div className="flex items-start gap-4">
                            <div className="p-3 rounded-full bg-yellow-500/10">
                                <AlertCircle className="w-6 h-6 text-yellow-500" />
                            </div>
                            <div>
                                <h4 className="font-black text-white text-sm uppercase tracking-wider mb-1">
                                    Verify your email
                                </h4>
                                <p className="text-xs text-text-dim max-w-lg">
                                    Email verification is required to secure your account and unlock all features.
                                    Check your inbox or request a new link.
                                </p>
                            </div>
                        </div>
                        <button
                            onClick={handleResendVerification}
                            className="px-6 py-3 bg-yellow-500/10 text-yellow-500 border border-yellow-500/30 rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-yellow-500/20 transition-all whitespace-nowrap"
                        >
                            Resend Verification
                        </button>
                    </div>
                </div>
            )}

            {/* Feedback Message */}
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

            {/* Cancel Confirmation Modal */}
            <ConfirmModal
                open={showCancelModal}
                title="Discard changes?"
                description="You have unsaved changes. Are you sure you want to discard them?"
                onConfirm={confirmCancel}
                onClose={() => setShowCancelModal(false)}
                confirmText="Discard"
                cancelText="Keep Editing"
            />
        </div>
    );
};

export default ProfileEdit;
