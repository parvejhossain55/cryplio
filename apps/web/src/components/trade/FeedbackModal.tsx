"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Star, X, Loader2, MessageSquare } from "lucide-react";
import { toast } from "sonner";

interface FeedbackModalProps {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: (rating: string, comment: string) => Promise<void>;
}

const FeedbackModal = ({ isOpen, onClose, onConfirm }: FeedbackModalProps) => {
    const [rating, setRating] = useState("positive");
    const [comment, setComment] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);

    const ratingOptions = [
        { value: "positive", label: "Positive", color: "text-green-500", bg: "bg-green-500/10 border-green-500/20" },
        { value: "neutral", label: "Neutral", color: "text-yellow-500", bg: "bg-yellow-500/10 border-yellow-500/20" },
        { value: "negative", label: "Negative", color: "text-red-500", bg: "bg-red-500/10 border-red-500/20" },
    ];

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);
        try {
            await onConfirm(rating, comment);
            toast.success("Feedback submitted successfully");
            onClose();
        } catch (err: any) {
            toast.error(err.message || "Failed to submit feedback");
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
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        exit={{ opacity: 0, scale: 0.95 }}
                        className="relative w-full max-w-md glass rounded-[2rem] border border-white/5 p-8"
                    >
                        <button
                            onClick={onClose}
                            className="absolute top-6 right-6 p-2 rounded-full bg-white/5 text-text-dim hover:text-white transition-all"
                        >
                            <X className="w-5 h-5" />
                        </button>

                        <div className="space-y-6">
                            <div className="text-center space-y-2">
                                <div className="w-16 h-16 mx-auto rounded-full bg-primary/10 flex items-center justify-center mb-4">
                                    <Star className="w-8 h-8 text-primary" />
                                </div>
                                <h3 className="text-2xl font-black italic uppercase tracking-tight">Leave Feedback</h3>
                                <p className="text-sm text-text-dim font-bold">
                                    Rate your trading experience with this partner
                                </p>
                            </div>

                            <form onSubmit={handleSubmit} className="space-y-6">
                                <div>
                                    <label className="block text-sm font-black uppercase tracking-widest text-text-dim mb-3">
                                        Rating
                                    </label>
                                    <div className="grid grid-cols-3 gap-3">
                                        {ratingOptions.map((option) => (
                                            <button
                                                key={option.value}
                                                type="button"
                                                onClick={() => setRating(option.value)}
                                                className={`p-4 rounded-2xl border-2 font-black uppercase tracking-widest text-xs transition-all ${
                                                    rating === option.value
                                                        ? `${option.bg} border-current ${option.color}`
                                                        : "bg-white/5 border-white/10 text-text-dim hover:bg-white/10"
                                                }`}
                                            >
                                                {option.label}
                                            </button>
                                        ))}
                                    </div>
                                </div>

                                <div>
                                    <label className="block text-sm font-black uppercase tracking-widest text-text-dim mb-3">
                                        Comment (Optional)
                                    </label>
                                    <textarea
                                        value={comment}
                                        onChange={(e) => setComment(e.target.value)}
                                        placeholder="Share your experience..."
                                        rows={3}
                                        className="w-full bg-background/50 border border-white/5 rounded-2xl p-4 text-sm font-bold outline-none focus:border-primary transition-all resize-none"
                                    />
                                </div>

                                <button
                                    type="submit"
                                    disabled={isSubmitting}
                                    className="w-full py-4 bg-primary text-white rounded-2xl font-black uppercase tracking-widest text-xs hover:scale-[1.02] active:scale-95 transition-all shadow-xl shadow-primary/20 flex items-center justify-center gap-2 disabled:opacity-50 disabled:grayscale"
                                >
                                    {isSubmitting ? (
                                        <Loader2 className="w-4 h-4 animate-spin" />
                                    ) : (
                                        <>
                                            <MessageSquare className="w-4 h-4" />
                                            SUBMIT FEEDBACK
                                        </>
                                    )}
                                </button>
                            </form>
                        </div>
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};

export default FeedbackModal;