"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    Users,
    TrendingUp,
    DollarSign,
    AlertTriangle,
    Eye,
    Shield,
    Activity,
    Clock,
    CheckCircle2,
    XCircle,
    ArrowUpRight,
    ArrowDownLeft
} from "lucide-react";
import { useRouter } from "next/navigation";

interface AdminStats {
    total_users: number;
    active_users: number;
    total_trades: number;
    active_trades: number;
    total_volume: number;
    pending_disputes: number;
    pending_withdrawals: number;
    system_health: "healthy" | "warning" | "error";
}

interface RecentActivity {
    id: string;
    type: "trade" | "dispute" | "withdrawal" | "user";
    title: string;
    description: string;
    status: "pending" | "completed" | "failed";
    created_at: string;
    user_id?: string;
    username?: string;
}

interface SystemAlert {
    id: string;
    type: "error" | "warning" | "info";
    title: string;
    message: string;
    created_at: string;
    resolved: boolean;
}

const AdminDashboard = () => {
    const router = useRouter();
    const [stats, setStats] = useState<AdminStats | null>(null);
    const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
    const [alerts, setAlerts] = useState<SystemAlert[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fetchDashboardData = async () => {
        setLoading(true);
        try {
            const [statsRes, activityRes, alertsRes] = await Promise.all([
                fetch("/api/v1/admin/stats", { credentials: "include" }),
                fetch("/api/v1/admin/activity", { credentials: "include" }),
                fetch("/api/v1/admin/alerts", { credentials: "include" })
            ]);

            if (statsRes.ok) {
                const statsData = await statsRes.json();
                setStats(statsData);
            }

            if (activityRes.ok) {
                const activityData = await activityRes.json();
                setRecentActivity(activityData.activities || []);
            }

            if (alertsRes.ok) {
                const alertsData = await alertsRes.json();
                setAlerts(alertsData.alerts || []);
            }
        } catch (error) {
            console.error("Error fetching dashboard data:", error);
        } finally {
            setLoading(false);
        }
    };

    const getHealthColor = (health: string) => {
        switch (health) {
            case "healthy":
                return "text-green-500";
            case "warning":
                return "text-yellow-500";
            case "error":
                return "text-red-500";
            default:
                return "text-gray-500";
        }
    };

    const getActivityIcon = (type: string) => {
        switch (type) {
            case "trade":
                return <DollarSign className="w-4 h-4" />;
            case "dispute":
                return <AlertTriangle className="w-4 h-4" />;
            case "withdrawal":
                return <ArrowUpRight className="w-4 h-4" />;
            case "user":
                return <Users className="w-4 h-4" />;
            default:
                return <Activity className="w-4 h-4" />;
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "pending":
                return "text-yellow-500";
            case "completed":
                return "text-green-500";
            case "failed":
                return "text-red-500";
            default:
                return "text-gray-500";
        }
    };

    const getAlertIcon = (type: string) => {
        switch (type) {
            case "error":
                return <XCircle className="w-4 h-4 text-red-500" />;
            case "warning":
                return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
            case "info":
                return <Shield className="w-4 h-4 text-blue-500" />;
            default:
                return <Activity className="w-4 h-4 text-gray-500" />;
        }
    };

    if (loading) {
        return (
            <DashboardLayout title="Admin Dashboard" role="admin">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Admin Dashboard" role="admin">
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-white">Admin Dashboard</h1>
                        <p className="text-text-dim">Platform overview and management</p>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className={`flex items-center gap-2 px-3 py-2 rounded-lg border ${
                            stats?.system_health === "healthy" 
                                ? "bg-green-500/10 border-green-500/20"
                                : stats?.system_health === "warning"
                                ? "bg-yellow-500/10 border-yellow-500/20"
                                : "bg-red-500/10 border-red-500/20"
                        }`}>
                            <Shield className={`w-4 h-4 ${getHealthColor(stats?.system_health || "")}`} />
                            <span className={`text-sm font-medium ${getHealthColor(stats?.system_health || "")}`}>
                                System {stats?.system_health}
                            </span>
                        </div>
                    </div>
                </div>

                {/* Stats Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <div className="p-2 bg-blue-500/20 rounded-lg">
                                <Users className="w-5 h-5 text-blue-500" />
                            </div>
                            <span className="text-xs text-text-dim">Total</span>
                        </div>
                        <h3 className="text-2xl font-bold text-white">{stats?.total_users || 0}</h3>
                        <p className="text-sm text-text-dim">Users</p>
                        <div className="mt-2 flex items-center gap-1 text-xs">
                            <span className="text-green-500">{stats?.active_users || 0}</span>
                            <span className="text-text-dim">active</span>
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <div className="p-2 bg-green-500/20 rounded-lg">
                                <TrendingUp className="w-5 h-5 text-green-500" />
                            </div>
                            <span className="text-xs text-text-dim">Total</span>
                        </div>
                        <h3 className="text-2xl font-bold text-white">{stats?.total_trades || 0}</h3>
                        <p className="text-sm text-text-dim">Trades</p>
                        <div className="mt-2 flex items-center gap-1 text-xs">
                            <span className="text-yellow-500">{stats?.active_trades || 0}</span>
                            <span className="text-text-dim">active</span>
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <div className="p-2 bg-purple-500/20 rounded-lg">
                                <DollarSign className="w-5 h-5 text-purple-500" />
                            </div>
                            <span className="text-xs text-text-dim">Total</span>
                        </div>
                        <h3 className="text-2xl font-bold text-white">
                            ${(stats?.total_volume || 0).toLocaleString()}
                        </h3>
                        <p className="text-sm text-text-dim">Volume</p>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <div className="p-2 bg-orange-500/20 rounded-lg">
                                <AlertTriangle className="w-5 h-5 text-orange-500" />
                            </div>
                            <span className="text-xs text-text-dim">Pending</span>
                        </div>
                        <h3 className="text-2xl font-bold text-white">
                            {(stats?.pending_disputes || 0) + (stats?.pending_withdrawals || 0)}
                        </h3>
                        <p className="text-sm text-text-dim">Actions Required</p>
                        <div className="mt-2 flex items-center gap-1 text-xs">
                            <span className="text-yellow-500">{stats?.pending_disputes || 0}</span>
                            <span className="text-text-dim">disputes,</span>
                            <span className="text-yellow-500">{stats?.pending_withdrawals || 0}</span>
                            <span className="text-text-dim">withdrawals</span>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Recent Activity */}
                    <div className="lg:col-span-2 bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-semibold text-white">Recent Activity</h3>
                            <button
                                onClick={() => router.push("/admin/activity")}
                                className="text-primary hover:text-primary/80 text-sm font-medium"
                            >
                                View All
                            </button>
                        </div>
                        <div className="space-y-3">
                            {recentActivity.length === 0 ? (
                                <p className="text-text-dim text-center py-8">No recent activity</p>
                            ) : (
                                recentActivity.slice(0, 5).map((activity) => (
                                    <motion.div
                                        key={activity.id}
                                        initial={{ opacity: 0, x: -20 }}
                                        animate={{ opacity: 1, x: 0 }}
                                        className="flex items-center gap-3 p-3 bg-white/5 rounded-lg"
                                    >
                                        <div className="p-2 bg-white/10 rounded-lg">
                                            {getActivityIcon(activity.type)}
                                        </div>
                                        <div className="flex-1">
                                            <p className="text-white font-medium">{activity.title}</p>
                                            <p className="text-text-dim text-sm">{activity.description}</p>
                                            {activity.username && (
                                                <p className="text-text-dim text-xs">by {activity.username}</p>
                                            )}
                                        </div>
                                        <div className="text-right">
                                            <div className={`flex items-center gap-1 ${getStatusColor(activity.status)}`}>
                                                {activity.status === "completed" && <CheckCircle2 className="w-3 h-3" />}
                                                {activity.status === "pending" && <Clock className="w-3 h-3" />}
                                                {activity.status === "failed" && <XCircle className="w-3 h-3" />}
                                                <span className="text-xs capitalize">{activity.status}</span>
                                            </div>
                                            <p className="text-text-dim text-xs mt-1">
                                                {new Date(activity.created_at).toLocaleDateString()}
                                            </p>
                                        </div>
                                    </motion.div>
                                ))
                            )}
                        </div>
                    </div>

                    {/* System Alerts */}
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-semibold text-white">System Alerts</h3>
                            <span className="px-2 py-1 bg-red-500/20 text-red-500 text-xs rounded-full">
                                {alerts.filter(a => !a.resolved).length}
                            </span>
                        </div>
                        <div className="space-y-3">
                            {alerts.length === 0 ? (
                                <p className="text-text-dim text-center py-8">No alerts</p>
                            ) : (
                                alerts.slice(0, 5).map((alert) => (
                                    <motion.div
                                        key={alert.id}
                                        initial={{ opacity: 0, y: -20 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        className={`p-3 rounded-lg border ${
                                            alert.type === "error" 
                                                ? "bg-red-500/10 border-red-500/20"
                                                : alert.type === "warning"
                                                ? "bg-yellow-500/10 border-yellow-500/20"
                                                : "bg-blue-500/10 border-blue-500/20"
                                        }`}
                                    >
                                        <div className="flex items-start gap-2">
                                            {getAlertIcon(alert.type)}
                                            <div className="flex-1">
                                                <p className="text-white font-medium text-sm">{alert.title}</p>
                                                <p className="text-text-dim text-xs mt-1">{alert.message}</p>
                                                <p className="text-text-dim text-xs mt-2">
                                                    {new Date(alert.created_at).toLocaleDateString()}
                                                </p>
                                            </div>
                                        </div>
                                    </motion.div>
                                ))
                            )}
                        </div>
                        {alerts.length > 5 && (
                            <button
                                onClick={() => router.push("/admin/alerts")}
                                className="w-full mt-3 py-2 bg-white/5 border border-white/10 rounded-lg text-sm text-text-dim hover:text-white hover:bg-white/10 transition-colors"
                            >
                                View All Alerts
                            </button>
                        )}
                    </div>
                </div>

                {/* Quick Actions */}
                <div className="bg-surface border border-white/10 rounded-2xl p-6">
                    <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <button
                            onClick={() => router.push("/admin/users")}
                            className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-center"
                        >
                            <Users className="w-6 h-6 mx-auto mb-2 text-primary" />
                            <p className="text-white text-sm font-medium">Manage Users</p>
                        </button>
                        <button
                            onClick={() => router.push("/admin/disputes")}
                            className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-center"
                        >
                            <AlertTriangle className="w-6 h-6 mx-auto mb-2 text-yellow-500" />
                            <p className="text-white text-sm font-medium">Disputes</p>
                        </button>
                        <button
                            onClick={() => router.push("/admin/withdrawals")}
                            className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-center"
                        >
                            <ArrowUpRight className="w-6 h-6 mx-auto mb-2 text-orange-500" />
                            <p className="text-white text-sm font-medium">Withdrawals</p>
                        </button>
                        <button
                            onClick={() => router.push("/admin/settings")}
                            className="p-4 bg-white/5 border border-white/10 rounded-xl hover:bg-white/10 transition-colors text-center"
                        >
                            <Shield className="w-6 h-6 mx-auto mb-2 text-accent" />
                            <p className="text-white text-sm font-medium">Settings</p>
                        </button>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminDashboard;
