"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Users,
    Search,
    Loader2,
    Shield,
    UserCheck,
    AlertTriangle,
    CheckCircle2,
    MoreVertical,
    X
} from "lucide-react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { toast } from "sonner";

interface User {
    id: string;
    username: string;
    email: string;
    status: string;
    is_suspended: boolean;
    created_at: string;
    last_login_at?: string;
    phone_verified: boolean;
    email_verified: boolean;
    login_count: number;
    failed_login_attempts: number;
}

const AdminUsersPage = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [filter, setFilter] = useState("all");
    const [searchTerm, setSearchTerm] = useState("");
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 20,
        total: 0,
        pages: 0,
    });

    const [isSuspendModalOpen, setIsSuspendModalOpen] = useState(false);
    const [suspendingUser, setSuspendingUser] = useState<string | null>(null);
    const [suspensionReason, setSuspensionReason] = useState("");
    const [suspensionDuration, setSuspensionDuration] = useState("");

    useEffect(() => {
        const timer = setTimeout(() => {
            fetchUsers(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchTerm, filter]);

    useEffect(() => {
        fetchUsers(pagination.page);
    }, [pagination.page]);

    const fetchUsers = async (page = 1) => {
        setIsLoading(true);
        try {
            const offset = (page - 1) * pagination.limit;
            const params = new URLSearchParams({
                limit: pagination.limit.toString(),
                offset: offset.toString(),
            });

            if (searchTerm) params.append("search", searchTerm);
            if (filter !== "all") params.append("status", filter);

            const response = await fetch(`/api/v1/admin/users?${params}`);
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch users");
            }
            
            setUsers(data.users || []);
            setPagination(prev => ({
                ...prev,
                page,
                total: data.total || 0,
                pages: Math.ceil((data.total || 0) / prev.limit)
            }));
        } catch (error) {
            console.error("Error fetching users:", error);
            toast.error("Failed to load users");
        } finally {
            setIsLoading(false);
        }
    };

    
    const handleSuspendUser = async () => {
        if (!suspendingUser || !suspensionReason) return;
        
        try {
            const response = await fetch(`/api/v1/admin/users/${suspendingUser}/suspend`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ 
                    reason: suspensionReason,
                    duration_minutes: suspensionDuration ? parseInt(suspensionDuration) : undefined
                }),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to suspend user");
            }

            toast.success("User suspended successfully");
            setIsSuspendModalOpen(false);
            setSuspendingUser(null);
            setSuspensionReason("");
            setSuspensionDuration("");
            fetchUsers();
        } catch (error) {
            console.error("Error suspending user:", error);
            toast.error("Failed to suspend user");
        }
    };

    const confirmSuspendUser = (userId: string) => {
        setSuspendingUser(userId);
        setIsSuspendModalOpen(true);
    };

    const handleUnsuspendUser = async (userId: string) => {
        try {
            const response = await fetch(`/api/v1/admin/users/${userId}/unsuspend`, {
                method: "POST",
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to unsuspend user");
            }

            toast.success("User unsuspended successfully");
            fetchUsers();
        } catch (error) {
            console.error("Error unsuspending user:", error);
            toast.error("Failed to unsuspend user");
        }
    };

    const filteredUsers = users.filter(user => {
        const matchesFilter = 
            filter === "all" || 
            (filter === "active" && !user.is_suspended && user.status === 'active') ||
            (filter === "suspended" && user.is_suspended);
        
        const matchesSearch = 
            user.username?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            user.email?.toLowerCase().includes(searchTerm.toLowerCase());
        
        return matchesFilter && matchesSearch;
    });

    const getStatusColor = (user: User) => {
        if (user.is_suspended) return "bg-orange-500/20 text-orange-400";
        if (user.status === 'active') return "bg-green-500/20 text-green-400";
        return "bg-gray-500/20 text-gray-400";
    };

    const getStatusText = (user: User) => {
        if (user.is_suspended) return "Suspended";
        return user.status;
    };

    if (isLoading) {
        return (
            <DashboardLayout title="User Management" role="admin">
                <div className="flex items-center justify-center h-64">
                    <Loader2 className="w-8 h-8 animate-spin text-primary" />
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="User Management" role="admin">
            <div className="space-y-6">
                {/* Stats Overview */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Total Users</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {users.length}
                                </p>
                            </div>
                            <Users className="w-8 h-8 text-blue-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Active</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {users.filter(u => !u.is_suspended && u.status === 'active').length}
                                </p>
                            </div>
                            <UserCheck className="w-8 h-8 text-green-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Suspended</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {users.filter(u => u.is_suspended).length}
                                </p>
                            </div>
                            <Shield className="w-8 h-8 text-orange-500 opacity-50" />
                        </div>
                    </div>
                    
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-text-dim text-sm font-medium">Verified Emails</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {users.filter(u => u.email_verified).length}
                                </p>
                            </div>
                            <CheckCircle2 className="w-8 h-8 text-green-500 opacity-50" />
                        </div>
                    </div>
                </div>

                {/* Filters and Search */}
                <div className="flex flex-col md:flex-row gap-4">
                    <div className="flex-1 relative">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-text-dim" />
                        <input
                            type="text"
                            placeholder="Search by username or email..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className="w-full pl-10 pr-4 py-3 bg-surface border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                        />
                    </div>
                    
                    <select
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        className="px-4 py-3 bg-surface border border-white/10 rounded-xl text-white focus:outline-none focus:border-primary/50"
                    >
                        <option value="all">All Users</option>
                        <option value="active">Active</option>
                        <option value="suspended">Suspended</option>
                    </select>
                </div>

                {/* Users Table */}
                <div className="bg-surface border border-white/10 rounded-2xl overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full">
                            <thead className="bg-white/5">
                                <tr>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        User
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Status
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Verification
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Activity
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Joined
                                    </th>
                                    <th className="px-6 py-4 text-left text-xs font-medium text-text-dim uppercase tracking-wider">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                <AnimatePresence>
                                    {users.map((user) => (
                                        <motion.tr
                                            key={user.id}
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                            className="hover:bg-white/5 transition-colors"
                                        >
                                            <td className="px-6 py-4">
                                                <div>
                                                    <div className="flex items-center space-x-2">
                                                        <span className="text-sm font-medium text-white">
                                                            {user.username}
                                                        </span>
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        {user.email}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(user)}`}>
                                                    {getStatusText(user)}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="flex items-center space-x-2">
                                                    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                                                        user.email_verified ? 'bg-green-500/20 text-green-400' : 'bg-gray-500/20 text-gray-400'
                                                    }`}>
                                                        {user.email_verified ? 'Email' : 'No Email'}
                                                    </span>
                                                    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                                                        user.phone_verified ? 'bg-green-500/20 text-green-400' : 'bg-gray-500/20 text-gray-400'
                                                    }`}>
                                                        {user.phone_verified ? 'Phone' : 'No Phone'}
                                                    </span>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div>
                                                    <div className="text-sm text-white">
                                                        {user.login_count} logins
                                                    </div>
                                                    <div className="text-xs text-text-dim">
                                                        {user.failed_login_attempts} failed attempts
                                                    </div>
                                                    {user.last_login_at && (
                                                        <div className="text-xs text-text-dim">
                                                            Last: {new Date(user.last_login_at).toLocaleDateString()}
                                                        </div>
                                                    )}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="text-sm text-text-dim">
                                                    {user.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="flex items-center space-x-2 relative z-10">
                                                    {user.is_suspended ? (
                                                        <button
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                handleUnsuspendUser(user.id);
                                                            }}
                                                            className="p-2 bg-green-500/20 text-green-400 rounded-lg hover:bg-green-500/30 transition-colors pointer-events-auto"
                                                            title="Unsuspend User"
                                                        >
                                                            <CheckCircle2 className="w-4 h-4" />
                                                        </button>
                                                    ) : (
                                                        <button
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                confirmSuspendUser(user.id);
                                                            }}
                                                            className="p-2 bg-orange-500/20 text-orange-400 rounded-lg hover:bg-orange-500/30 transition-colors pointer-events-auto"
                                                            title="Suspend User"
                                                        >
                                                            <Shield className="w-4 h-4" />
                                                        </button>
                                                    )}
                                                    
                                                    <button
                                                        onClick={(e) => e.stopPropagation()}
                                                        className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors pointer-events-auto"
                                                        title="More Options"
                                                    >
                                                        <MoreVertical className="w-4 h-4" />
                                                    </button>
                                                </div>
                                            </td>
                                        </motion.tr>
                                    ))}
                                </AnimatePresence>
                            </tbody>
                        </table>
                        
                        {users.length === 0 && (
                            <div className="text-center py-12">
                                <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                <p className="text-text-dim">No users found</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Pagination */}
                {pagination.pages > 1 && (
                    <div className="flex items-center justify-between pt-6 border-t border-white/5">
                        <div className="text-sm text-text-dim">
                            Showing {((pagination.page - 1) * pagination.limit) + 1} to {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} users
                        </div>
                        <div className="flex items-center space-x-2">
                            <button
                                onClick={() => setPagination(prev => ({ ...prev, page: prev.page - 1 }))}
                                disabled={pagination.page <= 1}
                                className="px-3 py-2 bg-white/5 border border-white/5 rounded-lg text-sm font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-all"
                            >
                                Previous
                            </button>

                            {Array.from({ length: Math.min(5, pagination.pages) }, (_, i) => {
                                const pageNum = i + 1;
                                return (
                                    <button
                                        key={pageNum}
                                        onClick={() => setPagination(prev => ({ ...prev, page: pageNum }))}
                                        className={`px-3 py-2 rounded-lg text-sm font-bold transition-all ${pagination.page === pageNum
                                            ? 'bg-primary text-white'
                                            : 'bg-white/5 border border-white/5 text-text-dim hover:bg-white/10'
                                            }`}
                                    >
                                        {pageNum}
                                    </button>
                                );
                            })}

                            <button
                                onClick={() => setPagination(prev => ({ ...prev, page: prev.page + 1 }))}
                                disabled={pagination.page >= pagination.pages}
                                className="px-3 py-2 bg-white/5 border border-white/5 rounded-lg text-sm font-bold text-white disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/10 transition-all"
                            >
                                Next
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Suspend User Modal */}
            <AnimatePresence>
                {isSuspendModalOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-surface border border-white/10 rounded-3xl w-full max-w-md overflow-hidden shadow-2xl"
                        >
                            <div className="p-6 border-b border-white/5 flex items-center justify-between">
                                <h3 className="text-xl font-bold text-white">Suspend User</h3>
                                <button
                                    onClick={() => setIsSuspendModalOpen(false)}
                                    className="p-2 hover:bg-white/5 rounded-lg transition-colors"
                                >
                                    <X className="w-5 h-5 text-text-dim" />
                                </button>
                            </div>

                            <div className="p-6 space-y-4">
                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-text-dim">Reason for Suspension</label>
                                    <textarea
                                        value={suspensionReason}
                                        onChange={(e) => setSuspensionReason(e.target.value)}
                                        placeholder="e.g. Violation of terms, Suspicious activity..."
                                        className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-text-dim/50 focus:outline-none focus:border-primary/50 min-h-[100px] resize-none"
                                    />
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-medium text-text-dim">Duration (minutes, optional)</label>
                                    <input
                                        type="number"
                                        value={suspensionDuration}
                                        onChange={(e) => setSuspensionDuration(e.target.value)}
                                        placeholder="Leave empty for permanent"
                                        className="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-text-dim/50 focus:outline-none focus:border-primary/50"
                                    />
                                </div>

                                <div className="p-4 bg-orange-500/10 border border-orange-500/20 rounded-2xl flex gap-3">
                                    <AlertTriangle className="w-5 h-5 text-orange-500 shrink-0 mt-0.5" />
                                    <p className="text-xs text-orange-200/80 leading-relaxed">
                                        Suspending this user will prevent them from trading and accessing their account features.
                                    </p>
                                </div>
                            </div>

                            <div className="p-6 bg-white/5 flex gap-3">
                                <button
                                    onClick={() => setIsSuspendModalOpen(false)}
                                    className="flex-1 px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white font-bold hover:bg-white/10 transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleSuspendUser}
                                    disabled={!suspensionReason}
                                    className="flex-1 px-4 py-3 bg-orange-500 text-white rounded-xl font-bold hover:bg-orange-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    Suspend User
                                </button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </DashboardLayout>
    );
};

export default AdminUsersPage;
