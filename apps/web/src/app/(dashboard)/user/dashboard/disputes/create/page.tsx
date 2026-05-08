"use client";

import React, { useState, useEffect } from "react";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ArrowLeft,
    Upload,
    X,
    FileText,
    AlertTriangle,
    CheckCircle2,
    Loader2
} from "lucide-react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

interface Trade {
    trade_id: string;
    ad_id: string;
    crypto_amount: number;
    fiat_amount: number;
    crypto_symbol: string;
    fiat_symbol: string;
    status: string;
    counterpart_username: string;
    payment_method: string;
    created_at: string;
}

const CreateDispute = () => {
    const router = useRouter();
    const [trades, setTrades] = useState<Trade[]>([]);
    const [selectedTrade, setSelectedTrade] = useState<Trade | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [uploadedFiles, setUploadedFiles] = useState<File[]>([]);
    const [formData, setFormData] = useState({
        reason: "",
        description: ""
    });

    const disputeReasons = [
        "PAYMENT_NOT_RECEIVED",
        "PAYMENT_NOT_SENT", 
        "FAKE_PAYMENT_PROOF",
        "WRONG_AMOUNT",
        "LATE_PAYMENT",
        "FRAUDULENT_ACTIVITY",
        "OTHER"
    ];

    useEffect(() => {
        fetchEligibleTrades();
    }, []);

    const fetchEligibleTrades = async () => {
        setLoading(true);
        try {
            const response = await fetch("/api/v1/trades/eligible-for-dispute", {
                credentials: "include",
            });
            
            if (!response.ok) {
                throw new Error("Failed to fetch eligible trades");
            }
            
            const data = await response.json();
            setTrades(data.trades || []);
        } catch (error) {
            console.error("Error fetching trades:", error);
            toast.error("Failed to load eligible trades");
        } finally {
            setLoading(false);
        }
    };

    const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(event.target.files || []);
        const validFiles = files.filter(file => {
            if (file.size > 10 * 1024 * 1024) {
                toast.error(`File ${file.name} is too large (max 10MB)`);
                return false;
            }
            return true;
        });
        
        setUploadedFiles(prev => [...prev, ...validFiles]);
    };

    const removeFile = (index: number) => {
        setUploadedFiles(prev => prev.filter((_, i) => i !== index));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (!selectedTrade) {
            toast.error("Please select a trade");
            return;
        }
        
        if (!formData.reason) {
            toast.error("Please select a dispute reason");
            return;
        }
        
        if (!formData.description.trim()) {
            toast.error("Please provide a description");
            return;
        }

        setSubmitting(true);
        try {
            // Upload evidence files first
            const evidenceUrls: string[] = [];
            
            for (const file of uploadedFiles) {
                const formData = new FormData();
                formData.append("file", file);
                
                const uploadResponse = await fetch("/api/upload/evidence", {
                    method: "POST",
                    body: formData,
                    credentials: "include",
                });
                
                if (!uploadResponse.ok) {
                    throw new Error(`Failed to upload ${file.name}`);
                }
                
                const uploadData = await uploadResponse.json();
                evidenceUrls.push(uploadData.url);
            }

            // Create dispute
            const disputeData = {
                trade_id: selectedTrade.trade_id,
                reason: formData.reason,
                description: formData.description.trim(),
                evidence_files: evidenceUrls
            };

            const response = await fetch("/api/disputes", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(disputeData),
                credentials: "include",
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || "Failed to create dispute");
            }

            const dispute = await response.json();
            toast.success("Dispute created successfully");
            router.push(`/dashboard/disputes/${dispute.dispute_id}`);
        } catch (error: any) {
            console.error("Error creating dispute:", error);
            toast.error(error.message || "Failed to create dispute");
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <DashboardLayout title="Create Dispute" role="user">
                <div className="flex items-center justify-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            </DashboardLayout>
        );
    }

    return (
        <DashboardLayout title="Create Dispute" role="user">
            <div className="max-w-4xl mx-auto space-y-6">
                {/* Header */}
                <div className="flex items-center gap-4">
                    <button
                        onClick={() => router.back()}
                        className="p-2 text-text-dim hover:text-white transition-colors"
                    >
                        <ArrowLeft className="w-5 h-5" />
                    </button>
                    <div>
                        <h1 className="text-2xl font-bold text-white">Create Dispute</h1>
                        <p className="text-text-dim">Report an issue with a trade</p>
                    </div>
                </div>

                {/* Warning Notice */}
                <div className="bg-yellow-500/10 border border-yellow-500/20 rounded-2xl p-4">
                    <div className="flex items-start gap-3">
                        <AlertTriangle className="w-5 h-5 text-yellow-500 flex-shrink-0 mt-0.5" />
                        <div>
                            <h3 className="text-sm font-semibold text-yellow-500 mb-1">Important Notice</h3>
                            <p className="text-yellow-500/80 text-sm">
                                Disputes should only be created for legitimate issues. False disputes may result in account penalties. 
                                Please ensure you have attempted to resolve the issue directly with your trading partner first.
                            </p>
                        </div>
                    </div>
                </div>

                <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Trade Selection */}
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <h3 className="text-lg font-semibold text-white mb-4">Select Trade</h3>
                        
                        {trades.length === 0 ? (
                            <div className="text-center py-8">
                                <AlertTriangle className="w-12 h-12 text-text-dim mx-auto mb-4" />
                                <h4 className="text-lg font-semibold text-white mb-2">No Eligible Trades</h4>
                                <p className="text-text-dim">
                                    You don't have any active trades that are eligible for dispute.
                                </p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {trades.map((trade) => (
                                    <label
                                        key={trade.trade_id}
                                        className={`block p-4 rounded-xl border-2 cursor-pointer transition-all ${
                                            selectedTrade?.trade_id === trade.trade_id
                                                ? "border-primary bg-primary/10"
                                                : "border-white/10 hover:border-white/20"
                                        }`}
                                    >
                                        <input
                                            type="radio"
                                            name="trade"
                                            value={trade.trade_id}
                                            checked={selectedTrade?.trade_id === trade.trade_id}
                                            onChange={() => setSelectedTrade(trade)}
                                            className="sr-only"
                                        />
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <div className="flex items-center gap-2 mb-2">
                                                    <span className="text-white font-semibold">
                                                        Trade #{trade.trade_id.slice(0, 8)}
                                                    </span>
                                                    <span className={`px-2 py-1 rounded-full text-xs ${
                                                        trade.status === "active" || trade.status === "paid"
                                                            ? "bg-yellow-500/20 text-yellow-500"
                                                            : "bg-gray-500/20 text-gray-500"
                                                    }`}>
                                                        {trade.status}
                                                    </span>
                                                </div>
                                                <p className="text-text-dim text-sm mb-1">
                                                    {trade.crypto_amount} {trade.crypto_symbol} / {trade.fiat_amount} {trade.fiat_symbol}
                                                </p>
                                                <p className="text-text-dim text-sm">
                                                    With {trade.counterpart_username} • {trade.payment_method}
                                                </p>
                                            </div>
                                            <div className="w-5 h-5 rounded-full border-2 border-white/30 flex items-center justify-center">
                                                {selectedTrade?.trade_id === trade.trade_id && (
                                                    <div className="w-3 h-3 rounded-full bg-primary"></div>
                                                )}
                                            </div>
                                        </div>
                                    </label>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Dispute Reason */}
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <h3 className="text-lg font-semibold text-white mb-4">Dispute Reason</h3>
                        <div className="grid grid-cols-2 gap-3">
                            {disputeReasons.map((reason) => (
                                <label
                                    key={reason}
                                    className={`block p-3 rounded-xl border-2 cursor-pointer transition-all ${
                                        formData.reason === reason
                                            ? "border-primary bg-primary/10"
                                            : "border-white/10 hover:border-white/20"
                                    }`}
                                >
                                    <input
                                        type="radio"
                                        name="reason"
                                        value={reason}
                                        checked={formData.reason === reason}
                                        onChange={(e) => setFormData({...formData, reason: e.target.value})}
                                        className="sr-only"
                                    />
                                    <div className="flex items-center gap-2">
                                        <div className="w-4 h-4 rounded-full border-2 border-white/30 flex items-center justify-center">
                                            {formData.reason === reason && (
                                                <div className="w-2 h-2 rounded-full bg-primary"></div>
                                            )}
                                        </div>
                                        <span className="text-sm text-white">
                                            {reason.replace(/_/g, " ").toLowerCase()}
                                        </span>
                                    </div>
                                </label>
                            ))}
                        </div>
                    </div>

                    {/* Description */}
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <h3 className="text-lg font-semibold text-white mb-4">Description</h3>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData({...formData, description: e.target.value})}
                            placeholder="Please provide a detailed description of the issue. Include relevant dates, amounts, and any communication with the trading partner."
                            rows={6}
                            className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50 resize-none"
                            required
                        />
                        <p className="text-text-dim text-xs mt-2">
                            Be as detailed as possible. This information will be used to resolve your dispute.
                        </p>
                    </div>

                    {/* Evidence Files */}
                    <div className="bg-surface border border-white/10 rounded-2xl p-6">
                        <h3 className="text-lg font-semibold text-white mb-4">Evidence Files</h3>
                        
                        <label className="flex items-center justify-center w-full p-8 border-2 border-dashed border-white/20 rounded-xl hover:border-primary/50 transition-colors cursor-pointer">
                            <input
                                type="file"
                                multiple
                                onChange={handleFileSelect}
                                className="hidden"
                                accept="image/*,.pdf,.doc,.docx,.txt,.jpg,.jpeg,.png"
                            />
                            <div className="text-center">
                                <Upload className="w-8 h-8 mx-auto mb-2 text-text-dim" />
                                <p className="text-white font-medium mb-1">Upload Evidence</p>
                                <p className="text-text-dim text-sm">
                                    Images, PDFs, documents (Max 10MB each)
                                </p>
                            </div>
                        </label>

                        {uploadedFiles.length > 0 && (
                            <div className="mt-4 space-y-2">
                                <p className="text-sm font-medium text-white">Selected Files:</p>
                                {uploadedFiles.map((file, index) => (
                                    <div
                                        key={index}
                                        className="flex items-center justify-between p-3 bg-white/5 rounded-lg"
                                    >
                                        <div className="flex items-center gap-2">
                                            <FileText className="w-4 h-4 text-primary" />
                                            <span className="text-sm text-white">{file.name}</span>
                                            <span className="text-text-dim text-xs">
                                                ({(file.size / 1024 / 1024).toFixed(2)} MB)
                                            </span>
                                        </div>
                                        <button
                                            type="button"
                                            onClick={() => removeFile(index)}
                                            className="p-1 text-text-dim hover:text-red-500 transition-colors"
                                        >
                                            <X className="w-4 h-4" />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Submit Button */}
                    <div className="flex justify-end">
                        <button
                            type="submit"
                            disabled={!selectedTrade || !formData.reason || !formData.description.trim() || submitting}
                            className="px-8 py-3 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                        >
                            {submitting ? (
                                <>
                                    <Loader2 className="w-4 h-4 animate-spin" />
                                    Creating Dispute...
                                </>
                            ) : (
                                <>
                                    <CheckCircle2 className="w-4 h-4" />
                                    Create Dispute
                                </>
                            )}
                        </button>
                    </div>
                </form>
            </div>
        </DashboardLayout>
    );
};

export default CreateDispute;
