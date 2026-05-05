"use client";

import React, { useState, useEffect, useRef } from "react";
import { User, Mail, Edit2, Save, X, CheckCircle, AlertCircle, Shield, Loader2, Camera } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
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
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);
    const [showCancelModal, setShowCancelModal] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const avatarUrl = user.avatarUrl || (user as any).avatar_url || null;

    useEffect(() => {
        if (!isEditing) {
            setUsername(user.username);
            setBio(user.bio || "");
            setAvatarPreview(null);
            setAvatarFile(null);
            setMessage(null);
        }
    }, [user, isEditing]);

    const handleSave = async () => {
        if (!username.trim()) {
            setMessage({ type: "error", text: "Username is required" });
            return;
        }
        setIsSaving(true);
        setMessage(null);
        try {
            await authService.updateCurrentUser({
                username: username.trim(),
                bio: bio.trim()
            });

            if (avatarFile) {
                const formData = new FormData();
                formData.append("avatar", avatarFile);
                const response = await fetch("/api/users/me/avatar", {
                    method: "POST",
                    body: formData,
                    credentials: "include",
                });
                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || "Failed to upload avatar");
                }
            }

            await refreshUser();
            setMessage({ type: "success", text: "Profile updated successfully" });
            setIsEditing(false);
            setAvatarFile(null);
            setAvatarPreview(null);
        } catch (error: any) {
            setMessage({ type: "error", text: error.message || "Failed to update profile" });
        } finally {
            setIsSaving(false);
        }
    };

    const handleCancel = () => {
        if (isEditing && (username !== user.username || bio !== (user.bio || "") || avatarFile)) {
            setShowCancelModal(true);
        } else {
            setIsEditing(false);
        }
    };

    const confirmCancel = () => {
        setShowCancelModal(false);
        setIsEditing(false);
        setAvatarFile(null);
        setAvatarPreview(null);
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

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        if (file.size > 2 * 1024 * 1024) {
            setMessage({ type: "error", text: "File size must be less than 2MB" });
            e.target.value = '';
            return;
        }
        if (!file.type.match(/image\/(jpeg|png)/)) {
            setMessage({ type: "error", text: "Only JPEG and PNG files are allowed" });
            e.target.value = '';
            return;
        }
        setAvatarFile(file);
        const reader = new FileReader();
        reader.onloadend = () => setAvatarPreview(reader.result as string);
        reader.readAsDataURL(file);
    };

    const displayUrl = avatarPreview || avatarUrl;
    const hasChanges = !!(avatarFile || username !== user.username || bio !== (user.bio || ""));

    const StatusBadges = () => (
        <div className="flex flex-wrap gap-2">
            {user.emailVerified ? (
                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-accent/10 text-accent border border-accent/20 text-[10px] font-black uppercase tracking-widest">
                    <CheckCircle className="w-3.5 h-3.5 mr-1.5" /> Verified
                </span>
            ) : (
                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-yellow-500/10 text-yellow-500 border border-yellow-500/20 text-[10px] font-black uppercase tracking-widest">
                    <AlertCircle className="w-3.5 h-3.5 mr-1.5" /> Unverified
                </span>
            )}
            {user.twoFAEnabled ? (
                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-accent/10 text-accent border border-accent/20 text-[10px] font-black uppercase tracking-widest">
                    <Shield className="w-3.5 h-3.5 mr-1.5" /> 2FA On
                </span>
            ) : (
                <span className="inline-flex items-center px-3 py-1.5 rounded-full bg-yellow-500/10 text-yellow-500 border border-yellow-500/20 text-[10px] font-black uppercase tracking-widest">
                    <Shield className="w-3.5 h-3.5 mr-1.5" /> 2FA Off
                </span>
            )}
        </div>
    );

    return (
        <div className="space-y-6">

            {/* Profile Header Card */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] overflow-hidden">
                {/* Top gradient accent line */}
                <div className="h-px w-full bg-gradient-to-r from-transparent via-primary/60 to-transparent" />

                <div className="p-8 md:p-10">
                    <AnimatePresence mode="wait">
                        {!isEditing ? (
                            /* ── VIEW MODE ── */
                            <motion.div
                                key="view"
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 0.18 }}
                                className="flex flex-col md:flex-row items-center md:items-start gap-8"
                            >
                                {/* Avatar */}
                                <div className="relative flex-shrink-0">
                                    {displayUrl ? (
                                        <img src={displayUrl} alt="Avatar" className="w-28 h-28 rounded-full border-2 border-white/10 object-cover" />
                                    ) : (
                                        <div className="w-28 h-28 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center">
                                            <User className="w-12 h-12 text-white/50" />
                                        </div>
                                    )}
                                </div>

                                {/* Info */}
                                <div className="flex-1 text-center md:text-left space-y-4">
                                    <div>
                                        <h2 className="text-2xl font-black text-white tracking-tight">{user.username}</h2>
                                        <p className="text-sm text-text-dim mt-1.5 flex items-center justify-center md:justify-start gap-2">
                                            <Mail className="w-4 h-4 flex-shrink-0" />
                                            {user.email}
                                        </p>
                                    </div>
                                    <StatusBadges />
                                </div>

                                {/* Edit button */}
                                <div className="flex-shrink-0">
                                    <button
                                        onClick={() => setIsEditing(true)}
                                        className="flex items-center justify-center px-6 py-3 bg-white/5 border border-white/10 text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/10 transition-all"
                                    >
                                        <Edit2 className="w-4 h-4 mr-2" />
                                        Edit Profile
                                    </button>
                                </div>
                            </motion.div>
                        ) : (
                            /* ── EDIT MODE ── */
                            <motion.div
                                key="edit"
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -10 }}
                                transition={{ duration: 0.22 }}
                                className="space-y-8"
                            >
                                {/* Edit header row: title + action buttons */}
                                <div className="flex items-center justify-between flex-wrap gap-4">
                                    <div>
                                        <p className="text-[10px] font-black text-white/30 uppercase tracking-widest mb-0.5">Editing</p>
                                        <h3 className="text-lg font-black text-white uppercase tracking-tight">Your Profile</h3>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <button
                                            onClick={handleCancel}
                                            disabled={isSaving}
                                            className="flex items-center px-5 py-2.5 bg-white/5 border border-white/10 text-white/60 rounded-xl text-xs font-black uppercase tracking-widest hover:bg-white/10 hover:text-white transition-all disabled:opacity-40"
                                        >
                                            <X className="w-3.5 h-3.5 mr-1.5" />
                                            Cancel
                                        </button>
                                        <button
                                            onClick={handleSave}
                                            disabled={isSaving || !hasChanges}
                                            className="flex items-center px-6 py-2.5 bg-primary border border-primary/80 text-white rounded-xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:scale-100"
                                        >
                                            {isSaving
                                                ? <><Loader2 className="w-3.5 h-3.5 mr-1.5 animate-spin" />Saving…</>
                                                : <><Save className="w-3.5 h-3.5 mr-1.5" />Save Changes</>
                                            }
                                        </button>
                                    </div>
                                </div>

                                {/* Divider */}
                                <div className="border-t border-white/5" />

                                {/* Avatar + fields */}
                                <div className="flex flex-col md:flex-row gap-8 md:gap-12 md:items-start">

                                    {/* Avatar uploader column */}
                                    <div className="flex flex-col items-center gap-4 flex-shrink-0">
                                        <div
                                            className="relative group cursor-pointer"
                                            onClick={() => fileInputRef.current?.click()}
                                        >
                                            {displayUrl ? (
                                                <img
                                                    src={displayUrl}
                                                    alt="Avatar"
                                                    className="w-28 h-28 rounded-full border-2 border-white/10 object-cover transition-all duration-300 group-hover:brightness-60"
                                                />
                                            ) : (
                                                <div className="w-28 h-28 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center transition-all duration-300 group-hover:brightness-60">
                                                    <User className="w-12 h-12 text-white/50" />
                                                </div>
                                            )}
                                            {/* Hover overlay */}
                                            <div className="absolute inset-0 rounded-full flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 gap-1">
                                                <Camera className="w-6 h-6 text-white drop-shadow-lg" />
                                                <span className="text-[9px] font-black text-white uppercase tracking-wider drop-shadow-lg">Change</span>
                                            </div>
                                        </div>

                                        <input
                                            ref={fileInputRef}
                                            type="file"
                                            accept="image/jpeg,image/png"
                                            onChange={handleFileChange}
                                            className="hidden"
                                        />

                                        <button
                                            type="button"
                                            onClick={() => fileInputRef.current?.click()}
                                            className="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/5 border border-white/10 text-white/50 text-[11px] font-black uppercase tracking-widest hover:bg-white/10 hover:text-white transition-all"
                                        >
                                            <Camera className="w-3.5 h-3.5" />
                                            {avatarFile ? "Change photo" : "Upload photo"}
                                        </button>

                                        <p className="text-[10px] text-white/20 text-center leading-relaxed">
                                            JPG or PNG · Max 2MB
                                        </p>
                                    </div>

                                    {/* Form fields column */}
                                    <div className="flex-1 space-y-5">

                                        {/* Username field */}
                                        <div className="space-y-2">
                                            <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">
                                                Username
                                            </label>
                                            <div className="relative">
                                                <User className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-white/25 pointer-events-none" />
                                                <input
                                                    type="text"
                                                    value={username}
                                                    onChange={(e) => setUsername(e.target.value)}
                                                    className="w-full bg-white/5 border border-white/10 py-3.5 pl-11 pr-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all placeholder:text-white/20"
                                                    placeholder="Your username"
                                                />
                                            </div>
                                        </div>

                                        {/* Email (read-only) */}
                                        <div className="space-y-2">
                                            <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">
                                                Email &nbsp;
                                                <span className="text-white/20 normal-case font-medium tracking-normal">· cannot be changed</span>
                                            </label>
                                            <div className="relative">
                                                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-white/15 pointer-events-none" />
                                                <input
                                                    type="email"
                                                    value={user.email}
                                                    readOnly
                                                    className="w-full bg-white/[0.03] border border-white/5 py-3.5 pl-11 pr-5 rounded-2xl text-sm font-bold text-white/25 outline-none cursor-not-allowed"
                                                />
                                            </div>
                                        </div>

                                        {/* Status badges in edit mode */}
                                        <div className="pt-1">
                                            <p className="text-[10px] font-black text-white/25 uppercase tracking-widest mb-2.5">Account Status</p>
                                            <StatusBadges />
                                        </div>
                                    </div>
                                </div>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>

            {/* Bio Section */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <h3 className="text-xl font-black text-white mb-6 uppercase tracking-tight flex items-center">
                    <User className="w-5 h-5 mr-3 text-primary" />
                    About Me
                </h3>
                {isEditing ? (
                    <div className="space-y-3">
                        <textarea
                            value={bio}
                            onChange={(e) => setBio(e.target.value)}
                            rows={4}
                            maxLength={300}
                            placeholder="Tell others about yourself…"
                            className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all resize-none placeholder:text-white/20"
                        />
                        <div className="flex items-center justify-between">
                            <p className="text-[10px] text-white/20">Plain text only</p>
                            <p className={`text-[10px] font-black tabular-nums ${bio.length > 270 ? "text-yellow-500" : "text-white/25"}`}>
                                {bio.length} / 300
                            </p>
                        </div>
                    </div>
                ) : (
                    user.bio ? (
                        <p className="text-sm text-text-dim leading-relaxed whitespace-pre-wrap">{user.bio}</p>
                    ) : (
                        <p className="text-sm text-white/25 italic">No bio added yet. Click Edit Profile to add one.</p>
                    )
                )}
            </div>

            {/* Email Verification Banner */}
            {!user.emailVerified && (
                <div className="bg-surface border border-yellow-500/20 border-l-4 border-l-yellow-500 rounded-[2.5rem] p-8 md:p-10">
                    <div className="flex flex-col md:flex-row items-center justify-between gap-4">
                        <div className="flex items-start gap-4">
                            <div className="p-3 rounded-full bg-yellow-500/10 flex-shrink-0">
                                <AlertCircle className="w-6 h-6 text-yellow-500" />
                            </div>
                            <div>
                                <h4 className="font-black text-white text-sm uppercase tracking-wider mb-1">Verify your email</h4>
                                <p className="text-xs text-text-dim max-w-lg">
                                    Email verification is required to secure your account and unlock all features. Check your inbox or request a new link.
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
            <AnimatePresence>
                {message && (
                    <motion.div
                        initial={{ opacity: 0, y: -10 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -10 }}
                        className={`p-4 rounded-2xl flex items-center space-x-3 ${
                            message.type === "success"
                                ? "bg-green-500/10 border border-green-500/30 text-green-400"
                                : "bg-red-500/10 border border-red-500/30 text-red-400"
                        }`}
                    >
                        {message.type === "success"
                            ? <CheckCircle className="w-5 h-5 flex-shrink-0" />
                            : <AlertCircle className="w-5 h-5 flex-shrink-0" />
                        }
                        <span className="text-sm font-medium">{message.text}</span>
                    </motion.div>
                )}
            </AnimatePresence>

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