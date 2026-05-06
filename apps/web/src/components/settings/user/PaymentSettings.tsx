"use client";

import React, { useState, useEffect } from "react";
import {
    CreditCard, Wallet, Plus, Trash2, CheckCircle, AlertCircle,
    Loader2, ArrowUpRight, ArrowDownLeft, Building2, Edit2, X, Star
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { toast } from "sonner";
import ConfirmModal from "@/components/ui/ConfirmModal";
import { authService, UserPaymentMethod } from "@/services/authService";

// Supported payment method types (matches DB seed data)
const PAYMENT_METHOD_OPTIONS = [
    { code: "bkash", label: "bKash", category: "Mobile Money" },
    { code: "nagad", label: "Nagad", category: "Mobile Money" },
    { code: "bank_transfer", label: "Bank Transfer", category: "Bank" },
    { code: "wise", label: "Wise", category: "Online Wallet" },
    { code: "paypal", label: "PayPal", category: "Online Wallet" },
    { code: "sepa", label: "SEPA Transfer", category: "Bank" },
    { code: "upi", label: "UPI", category: "Mobile Money" },
    { code: "gcash", label: "GCash", category: "Mobile Money" },
    { code: "mpesa", label: "M-Pesa", category: "Mobile Money" },
] as const;

type MethodCode = typeof PAYMENT_METHOD_OPTIONS[number]["code"];

const getMethodIcon = (code: string) => {
    if (code === "bank_transfer" || code === "sepa") return Building2;
    if (code === "paypal" || code === "wise") return Wallet;
    return CreditCard;
};

const getMethodLabel = (code: string) =>
    PAYMENT_METHOD_OPTIONS.find(m => m.code === code)?.label ?? code;

// ─── Add / Edit Modal ─────────────────────────────────────────────────────
interface MethodFormProps {
    initial?: UserPaymentMethod;
    onClose: () => void;
    onSaved: () => void;
}

const MethodModal = ({ initial, onClose, onSaved }: MethodFormProps) => {
    const [code, setCode] = useState<MethodCode>((initial?.payment_method_code as MethodCode) ?? "bank_transfer");
    const [displayName, setDisplayName] = useState(initial?.display_name ?? "");
    const [accountName, setAccountName] = useState(initial?.account_name ?? "");
    const [accountNumber, setAccountNumber] = useState(initial?.account_number ?? "");
    const [bankName, setBankName] = useState(initial?.bank_name ?? "");
    const [isSaving, setIsSaving] = useState(false);
    const [dropdownOpen, setDropdownOpen] = useState(false);

    const isEdit = !!initial?.id;

    const needsBank = code === "bank_transfer" || code === "sepa";
    const needsAccount = code !== "paypal" && code !== "wise";

    const handleSubmit = async () => {
        if (!displayName.trim()) {
            toast.error("Display name is required");
            return;
        }
        if (!code) {
            toast.error("Payment method type is required");
            return;
        }

        setIsSaving(true);
        try {
            const payload = {
                payment_method_code: code,
                display_name: displayName.trim(),
                account_name: accountName.trim() || undefined,
                account_number: accountNumber.trim() || undefined,
                bank_name: bankName.trim() || undefined,
            };

            if (isEdit && initial?.id) {
                await authService.updatePaymentMethod(initial.id, payload);
                toast.success("Payment method updated");
            } else {
                await authService.createPaymentMethod(payload);
                toast.success("Payment method added");
            }
            onSaved();
            onClose();
        } catch (err: any) {
            toast.error(err.message || "Failed to save payment method");
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div
                className="absolute inset-0 bg-black/60 backdrop-blur-sm"
                onClick={() => { onClose(); setDropdownOpen(false); }}
            />
            <motion.div
                initial={{ opacity: 0, scale: 0.95, y: 20 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95, y: 20 }}
                className="relative bg-surface border border-white/10 rounded-3xl p-8 max-w-md w-full shadow-2xl z-10 max-h-[90vh] overflow-y-auto"
            >
                {/* Header */}
                <div className="flex items-center justify-between mb-6">
                    <h3 className="text-xl font-black text-white uppercase tracking-tight">
                        {isEdit ? "Edit Payment Method" : "Add Payment Method"}
                    </h3>
                    <button onClick={onClose} className="p-2 text-white/40 hover:text-white rounded-xl hover:bg-white/5 transition-all">
                        <X className="w-5 h-5" />
                    </button>
                </div>

                <div className="space-y-4">
                    {/* Method Type — custom dropdown to avoid white OS option background */}
                    {!isEdit && (
                        <div className="space-y-2">
                            <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">Type</label>
                            <div className="relative">
                                <button
                                    type="button"
                                    onClick={() => setDropdownOpen(o => !o)}
                                    className="w-full flex items-center justify-between bg-white/5 border border-white/10 py-3.5 px-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary transition-all text-left"
                                >
                                    <span>
                                        {PAYMENT_METHOD_OPTIONS.find(o => o.code === code)?.label ?? code}
                                        <span className="ml-2 text-[10px] text-white/30 font-medium">
                                            {PAYMENT_METHOD_OPTIONS.find(o => o.code === code)?.category}
                                        </span>
                                    </span>
                                    <svg className={`w-4 h-4 text-white/40 transition-transform ${dropdownOpen ? "rotate-180" : ""}`} fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                </button>

                                {dropdownOpen && (
                                    <div className="absolute z-20 w-full mt-2 bg-[#1a1a2e] border border-white/10 rounded-2xl shadow-2xl overflow-hidden">
                                        {PAYMENT_METHOD_OPTIONS.map(opt => (
                                            <button
                                                key={opt.code}
                                                type="button"
                                                onClick={() => { setCode(opt.code); setDropdownOpen(false); }}
                                                className={`w-full flex items-center justify-between px-5 py-3 text-sm font-bold transition-all text-left ${code === opt.code
                                                    ? "bg-primary/20 text-white"
                                                    : "text-white/70 hover:bg-white/5 hover:text-white"
                                                    }`}
                                            >
                                                <span>{opt.label}</span>
                                                <span className="text-[10px] text-white/30 font-medium">{opt.category}</span>
                                            </button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        </div>
                    )}


                    {/* Display Name */}
                    <div className="space-y-2">
                        <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">Display Name</label>
                        <input
                            type="text"
                            value={displayName}
                            onChange={e => setDisplayName(e.target.value)}
                            placeholder={`e.g. My ${getMethodLabel(code)} Account`}
                            className="w-full bg-white/5 border border-white/10 py-3.5 px-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all placeholder:text-white/20"
                        />
                    </div>

                    {/* Account Name */}
                    <div className="space-y-2">
                        <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">Account Holder Name</label>
                        <input
                            type="text"
                            value={accountName}
                            onChange={e => setAccountName(e.target.value)}
                            placeholder="Full name as on account"
                            className="w-full bg-white/5 border border-white/10 py-3.5 px-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all placeholder:text-white/20"
                        />
                    </div>

                    {/* Account Number (context-aware) */}
                    {needsAccount && (
                        <div className="space-y-2">
                            <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">
                                {code === "upi" ? "UPI ID" : code === "bkash" || code === "nagad" || code === "gcash" || code === "mpesa" ? "Phone / Account Number" : "Account Number / IBAN"}
                            </label>
                            <input
                                type="text"
                                value={accountNumber}
                                onChange={e => setAccountNumber(e.target.value)}
                                placeholder={code === "upi" ? "name@bank" : "e.g. 015XXXXXXXX"}
                                className="w-full bg-white/5 border border-white/10 py-3.5 px-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all placeholder:text-white/20"
                            />
                        </div>
                    )}

                    {/* Bank Name */}
                    {needsBank && (
                        <div className="space-y-2">
                            <label className="block text-[10px] font-black text-white/40 uppercase tracking-widest">Bank Name</label>
                            <input
                                type="text"
                                value={bankName}
                                onChange={e => setBankName(e.target.value)}
                                placeholder="e.g. BRAC Bank"
                                className="w-full bg-white/5 border border-white/10 py-3.5 px-5 rounded-2xl text-sm font-bold text-white outline-none focus:border-primary focus:bg-primary/5 transition-all placeholder:text-white/20"
                            />
                        </div>
                    )}
                </div>

                {/* Actions */}
                <div className="flex gap-3 mt-8">
                    <button
                        onClick={onClose}
                        className="flex-1 px-6 py-3.5 bg-white/5 border border-white/10 text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/10 transition-all"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={isSaving || !displayName.trim()}
                        className="flex-1 flex items-center justify-center gap-2 px-6 py-3.5 bg-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
                    >
                        {isSaving ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
                        {isEdit ? "Save Changes" : "Add Method"}
                    </button>
                </div>
            </motion.div>
        </div>
    );
};

// ─── Main Component ────────────────────────────────────────────────────────
const PaymentSettings = () => {
    const [methods, setMethods] = useState<UserPaymentMethod[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [editTarget, setEditTarget] = useState<UserPaymentMethod | null>(null);
    const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const fetchMethods = async () => {
        setIsLoading(true);
        try {
            const data = await authService.getPaymentMethods();
            setMethods(data);
        } catch (err: any) {
            toast.error(err.message || "Failed to load payment methods");
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => { fetchMethods(); }, []);

    const handleDelete = async () => {
        if (!deleteTarget) return;
        setIsDeleting(true);
        try {
            await authService.deletePaymentMethod(deleteTarget);
            setMethods(prev => prev.filter(m => m.id !== deleteTarget));
            toast.success("Payment method removed");
        } catch (err: any) {
            toast.error(err.message || "Failed to remove payment method");
        } finally {
            setIsDeleting(false);
            setDeleteTarget(null);
        }
    };

    const handleSetDefault = async (id: string) => {
        try {
            await authService.setDefaultPaymentMethod(id);
            setMethods(prev => prev.map(m => ({ ...m, is_default: m.id === id })));
            toast.success("Default payment method updated");
        } catch (err: any) {
            toast.error(err.message || "Failed to set default");
        }
    };

    const MethodCard = ({ method }: { method: UserPaymentMethod }) => {
        const Icon = getMethodIcon(method.payment_method_code);
        return (
            <motion.div
                layout
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                className={`p-5 rounded-2xl border transition-all ${method.is_default
                    ? "bg-accent/5 border-accent/30"
                    : "bg-white/[0.03] border-white/5 hover:border-white/15"
                    }`}
            >
                <div className="flex items-start justify-between gap-4">
                    {/* Icon + info */}
                    <div className="flex items-center gap-4">
                        <div className={`p-3 rounded-xl flex-shrink-0 ${method.is_default ? "bg-accent/20 text-accent" : "bg-white/5 text-text-dim"}`}>
                            <Icon className="w-5 h-5" />
                        </div>
                        <div className="min-w-0">
                            <div className="flex items-center gap-2 flex-wrap">
                                <h4 className="font-bold text-white text-sm truncate">{method.display_name}</h4>
                                {method.is_default && (
                                    <span className="px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest bg-accent/20 text-accent border border-accent/30 flex-shrink-0">
                                        Default
                                    </span>
                                )}
                            </div>
                            <p className="text-[11px] text-text-dim mt-0.5">{getMethodLabel(method.payment_method_code)}</p>
                            {method.account_number && (
                                <p className="text-[10px] text-white/25 mt-0.5 font-mono">{method.account_number}</p>
                            )}
                            {method.bank_name && (
                                <p className="text-[10px] text-white/25 mt-0.5">{method.bank_name}</p>
                            )}
                        </div>
                    </div>

                    {/* Actions */}
                    <div className="flex items-center gap-1 flex-shrink-0">
                        {!method.is_default && (
                            <button
                                onClick={() => handleSetDefault(method.id)}
                                className="p-2 text-text-dim hover:text-accent hover:bg-accent/5 rounded-xl transition-all"
                                title="Set as default"
                            >
                                <Star className="w-4 h-4" />
                            </button>
                        )}
                        <button
                            onClick={() => { setEditTarget(method); setShowModal(true); }}
                            className="p-2 text-text-dim hover:text-white hover:bg-white/5 rounded-xl transition-all"
                            title="Edit"
                        >
                            <Edit2 className="w-4 h-4" />
                        </button>
                        <button
                            onClick={() => setDeleteTarget(method.id)}
                            className="p-2 text-text-dim hover:text-red-400 hover:bg-red-500/5 rounded-xl transition-all"
                            title="Remove"
                        >
                            <Trash2 className="w-4 h-4" />
                        </button>
                    </div>
                </div>
            </motion.div>
        );
    };

    return (
        <div className="space-y-8">
            {/* Header card */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <h3 className="text-xl font-black text-white mb-2 uppercase tracking-tight flex items-center">
                            <CreditCard className="w-6 h-6 mr-3 text-primary" />
                            Payment Methods
                        </h3>
                        <p className="text-xs text-text-dim">Your stored accounts for P2P trading</p>
                    </div>
                    <button
                        onClick={() => { setEditTarget(null); setShowModal(true); }}
                        className="flex items-center gap-2 px-6 py-3 bg-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20"
                    >
                        <Plus className="w-4 h-4" />
                        Add Method
                    </button>
                </div>

                {isLoading ? (
                    <div className="flex items-center justify-center py-16 gap-3">
                        <Loader2 className="w-6 h-6 animate-spin text-primary" />
                        <p className="text-xs text-text-dim font-black uppercase tracking-widest">Loading methods…</p>
                    </div>
                ) : methods.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-16 border-2 border-dashed border-white/5 rounded-2xl space-y-3">
                        <div className="w-14 h-14 bg-white/5 rounded-full flex items-center justify-center">
                            <CreditCard className="w-7 h-7 text-white/20" />
                        </div>
                        <p className="text-sm text-text-dim">No payment methods yet</p>
                        <button
                            onClick={() => { setEditTarget(null); setShowModal(true); }}
                            className="text-xs text-primary font-black uppercase tracking-widest hover:underline"
                        >
                            + Add your first method
                        </button>
                    </div>
                ) : (
                    <div className="space-y-3">
                        {methods.map(m => <MethodCard key={m.id} method={m} />)}
                    </div>
                )}
            </div>

            {/* Modals */}
            <AnimatePresence>
                {showModal && (
                    <MethodModal
                        initial={editTarget ?? undefined}
                        onClose={() => { setShowModal(false); setEditTarget(null); }}
                        onSaved={fetchMethods}
                    />
                )}
            </AnimatePresence>

            <ConfirmModal
                open={!!deleteTarget}
                title="Remove payment method?"
                description="This payment method will be permanently removed from your account."
                onConfirm={handleDelete}
                onClose={() => setDeleteTarget(null)}
                confirmText={isDeleting ? "Removing…" : "Remove"}
                cancelText="Keep"
            />
        </div>
    );
};

export default PaymentSettings;
