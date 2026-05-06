"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { AlertTriangle, X, Loader2, Info } from "lucide-react";
import { toast } from "sonner";

interface DisputeModalProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: (reasonCode: string, reasonText: string) => Promise<void>;
}

const REASON_CODES = [
    { code: "payment_not_received", label: "Payment Not Received" },
    { code: "buyer_claiming_unpaid", label: "Buyer Claiming Unpaid" },
    { code: "scam_attempt", label: "Suspected Scam Attempt" },
    { code: "other", label: "Other / Communication Issue" },
];

const DisputeModal = ({ isOpen, onClose, onConfirm }: DisputeModalProps) => {
    const [reasonCode, setReasonCode] = useState(REASON_CODES[0].code);
    const [reasonText, setReasonText] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        try {
            await onConfirm(reasonCode, reasonText);
            toast.success("Dispute raised successfully");
            onClose();
        } catch (err: any) {
            toast.error(err.message || "Failed to raise dispute");
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <AnimatePresence>
            {isOpen && (
                <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={onClose}
                        className="absolute inset-0 bg-background/80 backdrop-blur-sm"
                    />
                    <motion.div
                        initial={{ opacity: 0, scale: 0.9, y: 20 }}
                        animate={{ opacity: 1, scale: 1, y: 0 }}
                        exit={{ opacity: 0, scale: 0.9, y: 20 }}
                        className="relative w-full max-w-lg bg-surface border border-white/5 rounded-[2.5rem] overflow-hidden shadow-2xl"
                    >
                        {/* Header */}
                        <div className="p-8 border-b border-white/5 flex items-center justify-between bg-red-500/5">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-2xl bg-red-500/10 flex items-center justify-center">
                                    <AlertTriangle className="w-6 h-6 text-red-500" />
                                </div>
                                <div>
                                    <h3 className="text-xl font-black italic uppercase tracking-tight text-white">RAISE <span className="text-red-500">DISPUTE</span></h3>
                                    <p className="text-[10px] font-bold text-text-dim uppercase tracking-widest">Protocol Initiation</p>
                                </div>
                            </div>
                            <button onClick={onClose} className="p-2 hover:bg-white/5 rounded-xl transition-all">
                                <X className="w-5 h-5 text-text-dim" />
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="p-8 space-y-6">
                            <div className="p-4 bg-white/5 rounded-2xl border border-white/5 flex gap-3">
                                <Info className="w-5 h-5 text-primary shrink-0" />
                                <p className="text-[9px] font-bold text-text-dim leading-relaxed uppercase tracking-widest">
                                    Raising a dispute will lock the escrowed assets and notify an administrator. A moderator will join the chat to investigate.
                                </p>
                            </div>

                            <div className="space-y-4">
                                <div>
                                    <label className="text-[10px] font-black uppercase tracking-widest text-text-dim mb-2 block">Reason Category</label>
                                    <div className="grid grid-cols-1 gap-2">
                                        {REASON_CODES.map((rc) => (
                                            <button
                                                key={rc.code}
                                                type="button"
                                                onClick={() => setReasonCode(rc.code)}
                                                className={`p-4 rounded-xl text-left text-xs font-bold transition-all border ${reasonCode === rc.code
                                                    ? "bg-primary/10 border-primary text-white"
                                                    : "bg-background border-white/5 text-text-dim hover:border-white/10"
                                                    }`}
                                            >
                                                {rc.label}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                <div>
                                    <label className="text-[10px] font-black uppercase tracking-widest text-text-dim mb-2 block">Detailed Explanation</label>
                                    <textarea
                                        required
                                        value={reasonText}
                                        onChange={(e) => setReasonText(e.target.value)}
                                        placeholder="Describe the situation in detail..."
                                        rows={4}
                                        className="w-full bg-background border border-white/5 rounded-2xl p-4 text-sm font-bold outline-none focus:border-primary transition-all resize-none"
                                    />
                                </div>
                            </div>

                            <div className="flex gap-4 pt-4">
                                <button
                                    type="button"
                                    onClick={onClose}
                                    className="flex-1 py-4 rounded-2xl border border-white/5 font-black uppercase tracking-widest text-[10px] hover:bg-white/5 transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="flex-[2] py-4 bg-red-500 text-white rounded-2xl font-black uppercase tracking-widest text-[10px] hover:scale-[1.02] active:scale-95 transition-all shadow-xl shadow-red-500/20 disabled:grayscale"
                                >
                                    {isSubmitting ? <Loader2 className="w-4 h-4 animate-spin mx-auto" /> : "INITIATE DISPUTE"}
                                </button>
                            </div>
                        </form>
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};

export default DisputeModal;
