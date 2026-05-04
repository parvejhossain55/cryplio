"use client";

import React, { useEffect } from "react";
import { X } from "lucide-react";

interface ConfirmModalProps {
    open: boolean;
    title: string;
    description?: string;
    onConfirm: () => void;
    onClose: () => void;
    confirmText?: string;
    cancelText?: string;
    isDestructive?: boolean;
}

const ConfirmModal = ({
    open,
    title,
    description,
    onConfirm,
    onClose,
    confirmText = "Confirm",
    cancelText = "Cancel",
    isDestructive = false
}: ConfirmModalProps) => {
    useEffect(() => {
        const handleEscape = (e: KeyboardEvent) => {
            if (e.key === "Escape") onClose();
        };

        if (open) {
            document.addEventListener("keydown", handleEscape);
            document.body.style.overflow = "hidden";
        }

        return () => {
            document.removeEventListener("keydown", handleEscape);
            document.body.style.overflow = "unset";
        };
    }, [open, onClose]);

    if (!open) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div
                className="absolute inset-0 bg-black/60 backdrop-blur-sm"
                onClick={onClose}
            />
            <div className="relative bg-surface border border-white/10 rounded-3xl p-8 max-w-md w-full mx-4 shadow-2xl">
                <button
                    onClick={onClose}
                    className="absolute top-4 right-4 p-2 text-text-dim hover:text-white transition-colors rounded-xl hover:bg-white/5"
                >
                    <X className="w-5 h-5" />
                </button>

                <h3 className="text-xl font-black text-white mb-2">
                    {title}
                </h3>

                {description && (
                    <p className="text-sm text-text-dim mb-6">
                        {description}
                    </p>
                )}

                <div className="flex items-center gap-3">
                    <button
                        onClick={onClose}
                        className="flex-1 px-6 py-3 bg-white/5 border border-white/10 text-white rounded-2xl text-sm font-black hover:bg-white/10 transition-all"
                    >
                        {cancelText}
                    </button>
                    <button
                        onClick={() => {
                            onConfirm();
                            onClose();
                        }}
                        className={`flex-1 px-6 py-3 border rounded-2xl text-sm font-black hover:scale-105 active:scale-95 transition-all shadow-lg ${isDestructive
                                ? "bg-red-500 border-red-500 text-white shadow-red-500/20"
                                : "bg-primary border-primary text-white shadow-primary/20"
                            }`}
                    >
                        {confirmText}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ConfirmModal;
