"use client";

import React, { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import {
    User,
    Shield,
    Calendar,
    Star,
    TrendingUp,
    CheckCircle2,
    MessageSquare,
    ShieldCheck,
    Ban,
    Loader2,
    ArrowLeft,
    Clock,
    Zap,
    ThumbsUp,
    ThumbsDown,
    Activity
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { toast } from "sonner";
import { userService } from "@/services/userService";
import { BackendUser, UserStats } from "@/types/api";
import { useAuth } from "@/context/AuthContext";
import Navbar from "@/components/layout/Navbar";
import Footer from "@/components/layout/Footer";
import ConfirmModal from "@/components/ui/ConfirmModal";

const PublicProfilePage = () => {
    const { username } = useParams();
    const router = useRouter();
    const { user: currentUser } = useAuth();

    const [profile, setProfile] = useState<{ user: BackendUser; stats: UserStats } | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);


    useEffect(() => {
        const fetchProfile = async () => {
            try {
                setIsLoading(true);
                const data = await userService.getUserByUsername(username as string);
                setProfile(data);
            } catch (err: any) {
                setError(err.message || "Failed to load profile");
            } finally {
                setIsLoading(false);
            }
        };

        if (username) {
            fetchProfile();
        }
    }, [username]);



    if (isLoading) {
        return (
            <div className="min-h-screen bg-background flex flex-col items-center justify-center">
                <Loader2 className="w-12 h-12 text-primary animate-spin mb-4" />
                <p className="text-text-dim font-black uppercase tracking-widest text-xs animate-pulse">Loading Profile...</p>
            </div>
        );
    }

    if (error || !profile) {
        return (
            <div className="min-h-screen bg-background flex flex-col items-center justify-center p-6 text-center">
                <div className="w-20 h-20 rounded-3xl bg-red-500/10 flex items-center justify-center mb-6 border border-red-500/20">
                    <Ban className="w-10 h-10 text-red-500" />
                </div>
                <h1 className="text-3xl font-black text-white mb-2 uppercase tracking-tight">Profile Not Found</h1>
                <p className="text-text-dim max-w-md mb-8">{error || "The user you are looking for does not exist or has been deactivated."}</p>
                <button
                    onClick={() => router.push("/marketplace")}
                    className="px-8 py-4 bg-primary text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-105 transition-all shadow-xl shadow-primary/20"
                >
                    Back to Marketplace
                </button>
            </div>
        );
    }

    const { user, stats } = profile;
    const isOwnProfile = currentUser?.id === user.id;
    const joinDate = new Date(user.id ? (user as any).created_at || Date.now() : Date.now()).toLocaleDateString('en-US', {
        month: 'long',
        year: 'numeric'
    });

    const isOnline = user.is_online;

    return (
        <main className="min-h-screen bg-background">
            <Navbar />

            <div className="pt-28 pb-20">
                <div className="container mx-auto px-4 md:px-6">

                    {/* Back Button */}
                    <motion.button
                        initial={{ opacity: 0, x: -10 }}
                        animate={{ opacity: 1, x: 0 }}
                        onClick={() => router.back()}
                        className="flex items-center text-text-dim hover:text-white transition-colors mb-8 group"
                    >
                        <ArrowLeft className="w-5 h-5 mr-2 group-hover:-translate-x-1 transition-transform" />
                        <span className="text-xs font-black uppercase tracking-widest">Back</span>
                    </motion.button>

                    <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">

                        {/* ── LEFT COLUMN: INFO & BADGES ── */}
                        <div className="lg:col-span-4 space-y-6">

                            {/* Profile Card */}
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="glass rounded-[2.5rem] border-border p-8 md:p-10 relative overflow-hidden"
                            >
                                <div className="absolute top-0 right-0 p-6">
                                    <div className={`w-3 h-3 rounded-full ${isOnline ? "bg-accent animate-pulse shadow-[0_0_12px_rgba(46,213,115,0.6)]" : "bg-text-dim/30"}`} />
                                </div>

                                <div className="flex flex-col items-center text-center">
                                    <div className="relative mb-6">
                                        <div className="w-32 h-32 rounded-[2.5rem] bg-gradient-to-br from-primary/20 to-secondary/20 border-2 border-white/10 flex items-center justify-center overflow-hidden shadow-2xl">
                                            {user.avatar_url ? (
                                                <img src={user.avatar_url} alt={user.username} className="w-full h-full object-cover" />
                                            ) : (
                                                <User className="w-16 h-16 text-white/20" />
                                            )}
                                        </div>
                                        {user.is_merchant && (
                                            <div className="absolute -bottom-2 -right-2 bg-accent text-background p-2 rounded-xl shadow-xl border border-accent/20">
                                                <ShieldCheck className="w-5 h-5" />
                                            </div>
                                        )}
                                    </div>

                                    <h1 className="text-3xl font-black text-white tracking-tight mb-2 flex items-center gap-2">
                                        {user.username}
                                    </h1>

                                    <p className="inline-flex items-center px-4 py-1.5 rounded-full bg-white/5 border border-white/10 text-[10px] font-black uppercase tracking-widest text-text-dim mb-6">
                                        <Calendar className="w-3.5 h-3.5 mr-2" />
                                        Joined {joinDate}
                                    </p>

                                    {user.bio ? (
                                        <p className="text-sm text-text-dim leading-relaxed mb-8 italic">"{user.bio}"</p>
                                    ) : (
                                        <p className="text-sm text-white/10 italic mb-8 uppercase tracking-widest font-black text-[10px]">Strategic Trader</p>
                                    )}

                                    <div className="w-full space-y-3">
                                        <>
                                            <button className="w-full py-4 bg-primary text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-[1.02] active:scale-[0.98] transition-all shadow-xl shadow-primary/20 flex items-center justify-center gap-2">
                                                <MessageSquare className="w-4 h-4" /> Message
                                            </button>
                                        </>
                                        {isOwnProfile && (
                                            <button
                                                onClick={() => router.push("/user/dashboard/settings")}
                                                className="w-full py-4 bg-white/5 text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:bg-white/10 border border-white/10 transition-all"
                                            >
                                                Edit My Profile
                                            </button>
                                        )}
                                    </div>
                                </div>
                            </motion.div>

                            {/* Trust Metrics */}
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.1 }}
                                className="glass rounded-[2rem] border-border p-6 space-y-4"
                            >
                                <h3 className="text-[10px] font-black text-white/30 uppercase tracking-[0.2em] mb-2 px-2">Verification Status</h3>
                                <div className="space-y-3">
                                    <div className={`flex items-center justify-between p-4 rounded-2xl border ${user.email_verified ? "bg-accent/5 border-accent/20" : "bg-white/5 border-white/5"}`}>
                                        <div className="flex items-center gap-3">
                                            <Shield className={`w-5 h-5 ${user.email_verified ? "text-accent" : "text-white/20"}`} />
                                            <span className={`text-xs font-bold ${user.email_verified ? "text-white" : "text-white/30"}`}>Email Verified</span>
                                        </div>
                                        {user.email_verified && <CheckCircle2 className="w-4 h-4 text-accent" />}
                                    </div>
                                </div>
                            </motion.div>
                        </div>

                        {/* ── RIGHT COLUMN: STATS & ACTIVITY ── */}
                        <div className="lg:col-span-8 space-y-8">

                            {/* Summary Totals */}
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                {[
                                    { label: "Total Trades", value: stats.total_trades || 0, icon: Activity, color: "text-primary" },
                                    { label: "Completion", value: `${(stats.success_rate || 0).toFixed(1)}%`, icon: Zap, color: "text-accent" },
                                    { label: "Pos. Feedback", value: `${(((stats.positive_feedback_count || 0) / (stats.total_trades || 1)) * 100).toFixed(0)}%`, icon: ThumbsUp, color: "text-green-500" },
                                    { label: "Volume (USD)", value: `$${((stats.total_volume_usd || 0) / 1000).toFixed(1)}k`, icon: TrendingUp, color: "text-blue-500" },
                                ].map((stat, i) => (
                                    <motion.div
                                        key={i}
                                        initial={{ opacity: 0, scale: 0.9 }}
                                        animate={{ opacity: 1, scale: 1 }}
                                        transition={{ delay: 0.2 + i * 0.05 }}
                                        className="glass rounded-3xl border-border p-6 flex flex-col items-center justify-center text-center group hover:border-primary/30 transition-all cursor-default"
                                    >
                                        <stat.icon className={`w-6 h-6 ${stat.color} mb-3 group-hover:scale-110 transition-transform`} />
                                        <p className="text-2xl font-black text-white mb-1">{stat.value}</p>
                                        <p className="text-[9px] font-black text-text-dim uppercase tracking-widest">{stat.label}</p>
                                    </motion.div>
                                ))}
                            </div>

                            {/* Detailed Stats & Feedback Card */}
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.4 }}
                                className="glass rounded-[2.5rem] border-border overflow-hidden shadow-2xl"
                            >
                                <div className="p-8 border-b border-border bg-white/5 flex items-center justify-between">
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Trade Performance</h3>
                                    <div className="flex items-center gap-4">
                                        <div className="text-right">
                                            <p className="text-[9px] font-black text-text-dim uppercase tracking-widest mb-0.5">Rating</p>
                                            <div className="flex items-center gap-1.5">
                                                <Star className="w-4 h-4 text-amber-500 fill-amber-500" />
                                                <span className="text-lg font-black text-white">{stats.avg_rating?.toFixed(2) || "5.00"}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-8 grid grid-cols-1 md:grid-cols-2 gap-12">
                                    {/* Speed Metrics */}
                                    <div className="space-y-6">
                                        <h4 className="text-xs font-black text-white/40 uppercase tracking-widest flex items-center gap-2">
                                            <Clock className="w-4 h-4" /> Average speed
                                        </h4>
                                        <div className="space-y-6">
                                            <div className="space-y-2">
                                                <div className="flex justify-between text-xs font-bold mb-1">
                                                    <span className="text-text-dim">Payment (Avg)</span>
                                                    <span className="text-white">~4.5 min</span>
                                                </div>
                                                <div className="h-1.5 w-full bg-white/5 rounded-full overflow-hidden">
                                                    <div className="h-full bg-accent w-[85%] rounded-full shadow-[0_0_8px_rgba(46,213,115,0.4)]" />
                                                </div>
                                            </div>
                                            <div className="space-y-2">
                                                <div className="flex justify-between text-xs font-bold mb-1">
                                                    <span className="text-text-dim">Release (Avg)</span>
                                                    <span className="text-white">~2.1 min</span>
                                                </div>
                                                <div className="h-1.5 w-full bg-white/5 rounded-full overflow-hidden">
                                                    <div className="h-full bg-primary w-[92%] rounded-full shadow-[0_0_8px_rgba(255,255,255,0.2)]" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    {/* Feedback Breakdown */}
                                    <div className="space-y-6">
                                        <h4 className="text-xs font-black text-white/40 uppercase tracking-widest flex items-center gap-2">
                                            <ThumbsUp className="w-4 h-4" /> Feedback History
                                        </h4>
                                        <div className="space-y-4">
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-3">
                                                    <div className="w-8 h-8 rounded-lg bg-green-500/10 flex items-center justify-center border border-green-500/20">
                                                        <ThumbsUp className="w-4 h-4 text-green-500" />
                                                    </div>
                                                    <span className="text-sm font-bold text-white">Positive</span>
                                                </div>
                                                <span className="text-sm font-black text-white">{stats.positive_feedback_count}</span>
                                            </div>
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-3">
                                                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                                                        <Star className="w-4 h-4 text-text-dim" />
                                                    </div>
                                                    <span className="text-sm font-bold text-text-dim">Neutral</span>
                                                </div>
                                                <span className="text-sm font-black text-text-dim">{stats.neutral_feedback_count}</span>
                                            </div>
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-3">
                                                    <div className="w-8 h-8 rounded-lg bg-red-500/10 flex items-center justify-center border border-red-500/20">
                                                        <ThumbsDown className="w-4 h-4 text-red-500" />
                                                    </div>
                                                    <span className="text-sm font-bold text-text-dim">Negative</span>
                                                </div>
                                                <span className="text-sm font-black text-text-dim">{stats.negative_feedback_count}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </motion.div>

                            {/* Recent Activity / Ads placeholder */}
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.5 }}
                                className="glass rounded-[2.5rem] border-border p-8"
                            >
                                <div className="flex items-center justify-between mb-8">
                                    <h3 className="text-xl font-black text-white uppercase tracking-tight">Active Trade Offers</h3>
                                    <span className="px-3 py-1 bg-white/5 border border-white/10 rounded-full text-[10px] font-black text-text-dim uppercase tracking-widest">Live Now</span>
                                </div>
                                <div className="text-center py-16 border-2 border-dashed border-white/5 rounded-[2rem]">
                                    <TrendingUp className="w-12 h-12 text-white/5 mx-auto mb-4" />
                                    <p className="text-text-dim font-bold">No active public offers at this moment.</p>
                                    <p className="text-[10px] uppercase font-black text-white/20 tracking-widest mt-1">This merchant might be trading privately</p>
                                </div>
                            </motion.div>
                        </div>
                    </div>
                </div>
            </div>

            <Footer />
        </main>
    );
};

export default PublicProfilePage;
