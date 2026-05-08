"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Users,
    Search,
    Filter,
    Loader2,
    Shield,
    ShieldOff,
    UserCheck,
    Calendar,
    Mail,
    MoreVertical,
    AlertTriangle,
    CheckCircle2,
    XCircle
} from "lucide-react";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import { toast } from "sonner";

interface User {
    user_id: string;
    username: string;
    email: string;
    status: string;
    is_suspended: boolean;
    is_merchant: boolean;
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

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            const response = await fetch("/api/v1/admin/users");
            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch users");
            }
            
            setUsers(data.users || []);
        } catch (error) {
            console.error("Error fetching users:", error);
            toast.error("Failed to load users");
        } finally {
            setIsLoading(false);
        }
    };

    
    const handleSuspendUser = async (userId: string) => {
        const reason = prompt("Enter suspension reason:");
        const duration = prompt("Enter suspension duration in hours (optional):");
        
        try {
            const response = await fetch(`/api/admin/users/${userId}/suspend`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ 
                    reason,
                    duration: duration ? parseInt(duration) : undefined
                }),
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || "Failed to suspend user");
            }

            toast.success("User suspended successfully");
            fetchUsers();
        } catch (error) {
            console.error("Error suspending user:", error);
            toast.error("Failed to suspend user");
        }
    };

    const handleUnsuspendUser = async (userId: string) => {
        try {
            const response = await fetch(`/api/admin/users/${userId}/unsuspend`, {
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
            (filter === "suspended" && user.is_suspended) ||
            (filter === "merchants" && user.is_merchant);
        
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
                                <p className="text-text-dim text-sm font-medium">Merchants</p>
                                <p className="text-2xl font-bold text-white mt-1">
                                    {users.filter(u => u.is_merchant).length}
                                </p>
                            </div>
                            <ShieldOff className="w-8 h-8 text-purple-500 opacity-50" />
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
                        <option value="merchants">Merchants</option>
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
                                    {filteredUsers.map((user) => (
                                        <motion.tr
                                            key={user.user_id}
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
                                                        {user.is_merchant && (
                                                            <span className="px-2 py-1 bg-purple-500/20 text-purple-400 text-xs rounded-full">
                                                                Merchant
                                                            </span>
                                                        )}
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
                                                    {new Date(user.created_at).toLocaleDateString()}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4">
                                                <div className="flex items-center space-x-2">
                                                    {user.is_suspended ? (
                                                        <button
                                                            onClick={() => handleUnsuspendUser(user.user_id)}
                                                            className="p-2 bg-green-500/20 text-green-400 rounded-lg hover:bg-green-500/30 transition-colors"
                                                            title="Unsuspend User"
                                                        >
                                                            <CheckCircle2 className="w-4 h-4" />
                                                        </button>
                                                    ) : (
                                                        <button
                                                            onClick={() => handleSuspendUser(user.user_id)}
                                                            className="p-2 bg-orange-500/20 text-orange-400 rounded-lg hover:bg-orange-500/30 transition-colors"
                                                            title="Suspend User"
                                                        >
                                                            <Shield className="w-4 h-4" />
                                                        </button>
                                                    )}
                                                    
                                                    <button
                                                        className="p-2 bg-white/10 text-white rounded-lg hover:bg-white/20 transition-colors"
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
                        
                        {filteredUsers.length === 0 && (
                            <div className="text-center py-12">
                                <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                <p className="text-text-dim">No users found</p>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default AdminUsersPage;
