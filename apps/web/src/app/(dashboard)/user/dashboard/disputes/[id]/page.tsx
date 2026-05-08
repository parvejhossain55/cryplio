"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    AlertTriangle,
    Clock,
    CheckCircle2,
    XCircle,
    ArrowLeft,
    Upload,
    X,
    FileText,
    Download,
    MessageSquare,
    Send,
    Loader2
} from "lucide-react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

interface DisputeMessage {
    id: string;
    sender_id: string;
    sender_username: string;
    content: string;
    is_admin: boolean;
    created_at: string;
}

interface Dispute {
    dispute_id: string;
    trade_id: string;
    raiser_id: string;
    respondent_id: string;
    reason: string;
    description: string;
    status: "open" | "investigating" | "resolved" | "closed";
    resolution?: string;
    resolution_notes?: string;
    resolved_by?: string;
    resolved_at?: string;
    evidence_files: string[];
    created_at: string;
    updated_at: string;
    messages: DisputeMessage[];
    trade_details?: {
        crypto_amount: number;
        fiat_amount: number;
        crypto_symbol: string;
        fiat_symbol: string;
        counterpart_username: string;
        payment_method: string;
    };
}

const DisputeDetail = ({ params }: { params: { id: string } }) => {
    const router = useRouter();
    const { id } = params;
    const [dispute, setDispute] = useState<Dispute | null>(null);
    const [loading, setLoading] = useState(true);
    const [newMessage, setNewMessage] = useState("");
    const [isSending, setIsSending] = useState(false);
    const [isUploading, setIsUploading] = useState(false);

    useEffect(() => {
        if (id) {
            fetchDisputeDetail();
        }
    }, [id]);

    const fetchDisputeDetail = async () => {
        setLoading(true);
        try {
            const response = await fetch(`/api/disputes/${id}`, {
                credentials: "include",
            });
            
            if (!response.ok) {
                throw new Error("Failed to fetch dispute details");
            }
            
            const data = await response.json();
            setDispute(data);
        } catch (error) {
            console.error("Error fetching dispute:", error);
            toast.error("Failed to load dispute details");
        } finally {
            setLoading(false);
        }
    };

    const sendMessage = async () => {
        if (!newMessage.trim()) return;

        setIsSending(true);
        try {
            const response = await fetch(`/api/disputes/${id}/messages`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ content: newMessage.trim() }),
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to send message");
            }

            setNewMessage("");
            fetchDisputeDetail(); // Refresh to show new message
            toast.success("Message sent successfully");
        } catch (error) {
            console.error("Error sending message:", error);
            toast.error("Failed to send message");
        } finally {
            setIsSending(false);
        }
    };

    const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        if (file.size > 10 * 1024 * 1024) {
            toast.error("File size must be less than 10MB");
            return;
        }

        setIsUploading(true);
        try {
            const formData = new FormData();
            formData.append("file", file);

            const response = await fetch(`/api/disputes/${id}/evidence`, {
                method: "POST",
                body: formData,
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Failed to upload file");
            }

            fetchDisputeDetail(); // Refresh to show new evidence
            toast.success("Evidence uploaded successfully");
        } catch (error) {
            console.error("Error uploading file:", error);
            toast.error("Failed to upload file");
        } finally {
            setIsUploading(false);
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "open":
                return <AlertTriangle className="w-5 h-5 text-yellow-500" />;
            case "investigating":
                return <Clock className="w-5 h-5 text-blue-500" />;
            case "resolved":
                return <CheckCircle2 className="w-5 h-5 text-green-500" />;
            case "closed":
                return <XCircle className="w-5 h-5 text-gray-500" />;
            default:
                return <AlertTriangle className="w-5 h-5 text-gray-500" />;
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case "open":
                return "bg-yellow-500/20 text-yellow-500 border-yellow-500/20";
            case "investigating":
                return "bg-blue-500/20 text-blue-500 border-blue-500/20";
            case "resolved":
                return "bg-green-500/20 text-green-500 border-green-500/20";
            case "closed":
                return "bg-gray-500/20 text-gray-500 border-gray-500/20";
            default:
                return "bg-gray-500/20 text-gray-500 border-gray-500/20";
        }
    };

    if (loading) {
        return (
            <DashboardLayout title="Dispute Details" role="user">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            </DashboardLayout>
        );
    }

    if (!dispute) {
        return (
            <DashboardLayout title="Dispute Details" role="user">
                <div className="text-center py-12">
                    <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                    <h3 className="text-lg font-semibold text-white mb-2">Dispute not found</h3>
                    <p className="text-text-dim">The dispute you're looking for doesn't exist.</p>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Dispute Details" role="user">
            <div className="space-y-6">
                {/* Header */}
                <div className="flex items-center gap-4">
                    <button
                        onClick={() => router.back()}
                        className="p-2 text-text-dim hover:text-white transition-colors"
                    >
                        <ArrowLeft className="w-5 h-5" />
                    </button>
                    <div className="flex-1">
                        <h1 className="text-2xl font-bold text-white">Dispute #{dispute.dispute_id.slice(0, 8)}</h1>
                        <p className="text-text-dim">Trade #{dispute.trade_id.slice(0, 8)}</p>
                    </div>
                    <div className={`px-4 py-2 rounded-full text-sm font-black uppercase border flex items-center gap-2 ${getStatusColor(dispute.status)}`}>
                        {getStatusIcon(dispute.status)}
                        {dispute.status}
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Main Content */}
                    <div className="lg:col-span-2 space-y-6">
                        {/* Trade Details */}
                        {dispute.trade_details && (
                            <div className="bg-surface border border-white/10 rounded-2xl p-6">
                                <h3 className="text-lg font-semibold text-white mb-4">Trade Details</h3>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <p className="text-sm text-text-dim">Amount</p>
                                        <p className="text-white font-semibold">
                                            {dispute.trade_details.crypto_amount} {dispute.trade_details.crypto_symbol}
                                        </p>
                                        <p className="text-text-dim">
                                            {dispute.trade_details.fiat_amount} {dispute.trade_details.fiat_symbol}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-sm text-text-dim">Payment Method</p>
                                        <p className="text-white font-semibold">{dispute.trade_details.payment_method}</p>
                                        <p className="text-text-dim">With {dispute.trade_details.counterpart_username}</p>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Dispute Details */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h3 className="text-lg font-semibold text-white mb-4">Dispute Information</h3>
                            <div className="space-y-4">
                                <div>
                                    <p className="text-sm text-text-dim mb-1">Reason</p>
                                    <p className="text-white font-semibold">{dispute.reason}</p>
                                </div>
                                <div>
                                    <p className="text-sm text-text-dim mb-1">Description</p>
                                    <p className="text-text-dim">{dispute.description}</p>
                                </div>
                                <div className="flex items-center gap-4 text-sm text-text-dim">
                                    <span>Created {new Date(dispute.created_at).toLocaleDateString()}</span>
                                    <span>Updated {new Date(dispute.updated_at).toLocaleDateString()}</span>
                                </div>
                            </div>
                        </div>

                        {/* Resolution */}
                        {dispute.resolution && (
                            <div className="bg-green-500/10 border border-green-500/20 rounded-2xl p-6">
                                <h3 className="text-lg font-semibold text-green-500 mb-4">Resolution</h3>
                                <div className="space-y-2">
                                    <p className="text-white font-semibold">{dispute.resolution}</p>
                                    {dispute.resolution_notes && (
                                        <p className="text-text-dim">{dispute.resolution_notes}</p>
                                    )}
                                    {dispute.resolved_at && (
                                        <p className="text-text-dim text-sm">
                                            Resolved on {new Date(dispute.resolved_at).toLocaleDateString()}
                                        </p>
                                    )}
                                </div>
                            </div>
                        )}

                        {/* Messages */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h3 className="text-lg font-semibold text-white mb-4">Communication</h3>
                            <div className="space-y-4 max-h-96 overflow-y-auto">
                                {dispute.messages.map((message) => (
                                    <motion.div
                                        key={message.id}
                                        initial={{ opacity: 0, y: 20 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        className={`flex ${message.is_admin ? "justify-start" : "justify-end"}`}
                                    >
                                        <div className={`max-w-md ${
                                            message.is_admin 
                                                ? "bg-white/10 text-white" 
                                                : "bg-primary/20 text-white"
                                        } rounded-2xl p-4`}>
                                            <div className="flex items-center gap-2 mb-2">
                                                <span className="text-sm font-semibold">
                                                    {message.sender_username}
                                                </span>
                                                {message.is_admin && (
                                                    <span className="px-2 py-1 bg-accent/20 text-accent text-xs rounded-full">
                                                        Admin
                                                    </span>
                                                )}
                                            </div>
                                            <p className="text-sm">{message.content}</p>
                                            <p className="text-xs text-text-dim mt-2">
                                                {new Date(message.created_at).toLocaleString()}
                                            </p>
                                        </div>
                                    </motion.div>
                                ))}
                            </div>

                            {/* Message Input */}
                            {dispute.status === "open" && (
                                <div className="mt-4 pt-4 border-t border-white/10">
                                    <div className="flex gap-2">
                                        <input
                                            type="text"
                                            value={newMessage}
                                            onChange={(e) => setNewMessage(e.target.value)}
                                            onKeyPress={(e) => e.key === "Enter" && sendMessage()}
                                            placeholder="Type your message..."
                                            className="flex-1 px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                                            disabled={isSending}
                                        />
                                        <button
                                            onClick={sendMessage}
                                            disabled={!newMessage.trim() || isSending}
                                            className="px-4 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                                        >
                                            {isSending ? (
                                                <Loader2 className="w-4 h-4 animate-spin" />
                                            ) : (
                                                <Send className="w-4 h-4" />
                                            )}
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Sidebar */}
                    <div className="space-y-6">
                        {/* Evidence Files */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h3 className="text-lg font-semibold text-white mb-4">Evidence Files</h3>
                            
                            {dispute.status === "open" && (
                                <div className="mb-4">
                                    <label className="flex items-center justify-center w-full p-4 border-2 border-dashed border-white/20 rounded-xl hover:border-primary/50 transition-colors cursor-pointer">
                                        <input
                                            type="file"
                                            onChange={handleFileUpload}
                                            className="hidden"
                                            accept="image/*,.pdf,.doc,.docx,.txt"
                                        />
                                        <div className="text-center">
                                            {isUploading ? (
                                                <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2" />
                                            ) : (
                                                <Upload className="w-6 h-6 mx-auto mb-2" />
                                            )}
                                            <span className="text-sm text-text-dim">
                                                {isUploading ? "Uploading..." : "Upload Evidence"}
                                            </span>
                                        </div>
                                    </label>
                                </div>
                            )}

                            <div className="space-y-2">
                                {dispute.evidence_files.length === 0 ? (
                                    <p className="text-text-dim text-sm text-center py-4">
                                        No evidence files uploaded
                                    </p>
                                ) : (
                                    dispute.evidence_files.map((file, index) => (
                                        <div
                                            key={index}
                                            className="flex items-center justify-between p-3 bg-white/5 rounded-lg"
                                        >
                                            <div className="flex items-center gap-2">
                                                <FileText className="w-4 h-4 text-primary" />
                                                <span className="text-sm text-white truncate">
                                                    {file.split("/").pop()}
                                                </span>
                                            </div>
                                            <button
                                                onClick={() => window.open(file, "_blank")}
                                                className="p-1 text-text-dim hover:text-white transition-colors"
                                            >
                                                <Download className="w-4 h-4" />
                                            </button>
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>

                        {/* Quick Actions */}
                        <div className="bg-surface border border-white/10 rounded-2xl p-6">
                            <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
                            <div className="space-y-2">
                                <button
                                    onClick={() => router.push(`/dashboard/trades/${dispute.trade_id}`)}
                                    className="w-full px-4 py-2 bg-white/5 border border-white/10 rounded-xl text-sm text-text-dim hover:text-white hover:bg-white/10 transition-colors flex items-center justify-center gap-2"
                                >
                                    <FileText className="w-4 h-4" />
                                    View Trade Details
                                </button>
                                {dispute.status === "open" && (
                                    <button
                                        onClick={() => router.push("/support")}
                                        className="w-full px-4 py-2 bg-accent text-white rounded-xl text-sm font-medium hover:bg-accent/90 transition-colors flex items-center justify-center gap-2"
                                    >
                                        <MessageSquare className="w-4 h-4" />
                                        Contact Support
                                    </button>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </DashboardLayout>
    );
};

export default DisputeDetail;
