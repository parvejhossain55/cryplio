"use client";

import React, { useState, useEffect } from "react";
import { CreditCard, Wallet, Plus, Trash2, CheckCircle, AlertCircle, Loader2, Save, ArrowUpRight, ArrowDownLeft, Building2 } from "lucide-react";
import { motion } from "framer-motion";
import ConfirmModal from "@/components/ui/ConfirmModal";

interface PaymentMethod {
    id: string;
    type: "card" | "bank" | "crypto";
    name: string;
    last4?: string;
    bankName?: string;
    address?: string;
    isDefault: boolean;
}

const PaymentSettings = () => {
    const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([
        { id: "1", type: "card", name: "Visa ending in 4242", last4: "4242", isDefault: true },
        { id: "2", type: "bank", name: "Chase Bank", bankName: "Chase", isDefault: false },
    ]);
    const [showAddModal, setShowAddModal] = useState(false);
    const [newMethodType, setNewMethodType] = useState<"card" | "bank" | "crypto">("card");
    const [newMethodDetails, setNewMethodDetails] = useState({ name: "", last4: "", bankName: "", address: "" });
    const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

    useEffect(() => {
        const saved = localStorage.getItem("cryplio_payment_methods");
        if (saved) {
            try {
                setPaymentMethods(JSON.parse(saved));
            } catch (e) { console.error(e); }
        }
    }, []);

    const saveToStorage = (methods: PaymentMethod[]) => {
        localStorage.setItem("cryplio_payment_methods", JSON.stringify(methods));
    };

    const handleAdd = () => {
        if (!newMethodDetails.name.trim()) return;
        const newMethod: PaymentMethod = {
            id: Date.now().toString(),
            type: newMethodType,
            name: newMethodDetails.name.trim(),
            last4: newMethodDetails.last4.trim() || undefined,
            bankName: newMethodDetails.bankName.trim() || undefined,
            address: newMethodDetails.address.trim() || undefined,
            isDefault: paymentMethods.length === 0,
        };
        const updated = [...paymentMethods, newMethod];
        setPaymentMethods(updated);
        saveToStorage(updated);
        setShowAddModal(false);
        setNewMethodDetails({ name: "", last4: "", bankName: "", address: "" });
        setMessage({ type: "success", text: "Payment method added successfully" });
    };

    const handleDelete = (id: string) => {
        const updated = paymentMethods.filter(m => m.id !== id);
        setPaymentMethods(updated);
        saveToStorage(updated);
        setDeleteTarget(null);
        setMessage({ type: "success", text: "Payment method removed" });
    };

    const handleSetDefault = (id: string) => {
        const updated = paymentMethods.map(m => ({ ...m, isDefault: m.id === id }));
        setPaymentMethods(updated);
        saveToStorage(updated);
        setMessage({ type: "success", text: "Default payment method updated" });
    };

    const PaymentCard = ({ method }: { method: PaymentMethod }) => (
        <div className={`p-5 rounded-2xl border transition-all group ${
            method.isDefault 
                ? "bg-accent/10 border-accent/30" 
                : "bg-white/5 border-white/5 hover:border-white/20"
        }`}>
            <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                    <div className={`p-3 rounded-xl ${method.isDefault ? "bg-accent/20 text-accent" : "bg-surface-light text-text-dim"}`}>
                        {method.type === "card" ? <CreditCard className="w-6 h-6" /> : 
                         method.type === "bank" ? <Building2 className="w-6 h-6" /> : <Wallet className="w-6 h-6" />}
                    </div>
                    <div>
                        <div className="flex items-center gap-2">
                            <h4 className="font-bold text-white text-sm">{method.name}</h4>
                            {method.isDefault && (
                                <span className="px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest bg-accent/20 text-accent border border-accent/30">
                                    Default
                                </span>
                            )}
                        </div>
                        <p className="text-[10px] text-text-dim mt-1">
                            {method.type === "card" ? `Card ending in ${method.last4}` : 
                             method.type === "bank" ? method.bankName : "Crypto Wallet"}
                        </p>
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    {!method.isDefault && (
                        <button
                            onClick={() => handleSetDefault(method.id)}
                            className="p-2 text-text-dim hover:text-white hover:bg-white/5 rounded-lg transition-colors"
                            title="Set as default"
                        >
                            <CheckCircle className="w-4 h-4" />
                        </button>
                    )}
                    <button
                        onClick={() => setDeleteTarget(method.id)}
                        className="p-2 text-text-dim hover:text-red-400 hover:bg-red-500/5 rounded-lg transition-colors"
                        title="Remove payment method"
                    >
                        <Trash2 className="w-4 h-4" />
                    </button>
                </div>
            </div>
        </div>
    );

    return (
        <div className="space-y-8">
            {/* Payment Methods */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <div className="flex items-center justify-between mb-8">
                    <div>
                        <h3 className="text-xl font-black text-white mb-2 uppercase tracking-tight flex items-center">
                            <CreditCard className="w-6 h-6 mr-3 text-primary" />
                            Payment Methods
                        </h3>
                        <p className="text-xs text-text-dim">Manage your cards, bank accounts, and wallets</p>
                    </div>
                    <button
                        onClick={() => setShowAddModal(true)}
                        className="flex items-center px-6 py-3 bg-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20"
                    >
                        <Plus className="w-4 h-4 mr-2" />
                        Add Method
                    </button>
                </div>

                <div className="space-y-4">
                    {paymentMethods.length === 0 ? (
                        <div className="text-center py-12 border-2 border-dashed border-white/5 rounded-2xl">
                            <CreditCard className="w-12 h-12 text-text-dim mx-auto mb-4" />
                            <p className="text-sm text-text-dim">No payment methods added yet</p>
                        </div>
                    ) : (
                        paymentMethods.map(method => (
                            <PaymentCard key={method.id} method={method} />
                        ))
                    )}
                </div>
            </div>

            {/* Withdrawal Settings */}
            <div className="bg-surface border border-white/10 rounded-[2.5rem] p-8 md:p-10">
                <h3 className="text-xl font-black text-white mb-6 uppercase tracking-tight flex items-center">
                    <ArrowUpRight className="w-5 h-5 mr-3 text-primary" />
                    Withdrawal Settings
                </h3>

                <div className="space-y-6">
                    <div className="flex items-center justify-between p-5 rounded-2xl bg-white/5 border border-white/5">
                        <div className="flex items-center gap-4">
                            <div className="p-3 rounded-lg bg-surface-light border border-white/10">
                                <ArrowDownLeft className="w-5 h-5 text-primary" />
                            </div>
                            <div>
                                <h4 className="font-bold text-white text-sm">Auto-Withdrawal</h4>
                                <p className="text-[10px] text-text-dim mt-1">Automatically transfer funds to your default payment method</p>
                            </div>
                        </div>
                        <button className="relative w-11 h-6 rounded-full bg-white/10">
                            <span className="absolute left-1 top-1 w-4 h-4 rounded-full bg-white/40 transition-transform" />
                        </button>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Withdrawal Fee Cover</label>
                            <select className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all appearance-none">
                                <option value="sender">Sender pays fee</option>
                                <option value="recipient">Recipient pays fee</option>
                                <option value="split">Split equally</option>
                            </select>
                        </div>
                        <div className="space-y-2">
                            <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Min Withdrawal Amount</label>
                            <div className="relative">
                                <input
                                    type="text"
                                    value="0.001 BTC"
                                    readOnly
                                    className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold text-white cursor-not-allowed"
                                />
                                <span className="absolute right-6 top-1/2 -translate-y-1/2 text-xs text-text-dim">Contact support to change</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Add Payment Method Modal */}
            {showAddModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center">
                    <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={() => setShowAddModal(false)} />
                    <div className="relative bg-surface border border-white/10 rounded-3xl p-8 max-w-md w-full mx-4 shadow-2xl">
                        <h3 className="text-xl font-black text-white mb-6">Add Payment Method</h3>

                        <div className="space-y-4">
                            <div>
                                <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block mb-2">Type</label>
                                <div className="flex gap-3">
                                    {["card", "bank", "crypto"].map(type => (
                                        <button
                                            key={type}
                                            onClick={() => setNewMethodType(type as any)}
                                            className={`flex-1 py-3 rounded-2xl text-xs font-black uppercase tracking-widest border transition-all ${
                                                newMethodType === type
                                                    ? "bg-primary text-white border-primary"
                                                    : "bg-white/5 text-text-dim border-white/10 hover:border-white/20"
                                            }`}
                                        >
                                            {type === "card" ? "Card" : type === "bank" ? "Bank" : "Crypto"}
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <div className="space-y-2">
                                <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Display Name</label>
                                <input
                                    type="text"
                                    value={newMethodDetails.name}
                                    onChange={(e) => setNewMethodDetails({...newMethodDetails, name: e.target.value})}
                                    placeholder="e.g. Visa ending in 4242"
                                    className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all"
                                />
                            </div>

                            {newMethodType === "card" && (
                                <div className="space-y-2">
                                    <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Last 4 Digits</label>
                                    <input
                                        type="text"
                                        value={newMethodDetails.last4}
                                        onChange={(e) => setNewMethodDetails({...newMethodDetails, last4: e.target.value.slice(0,4)})}
                                        placeholder="4242"
                                        maxLength={4}
                                        className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all"
                                    />
                                </div>
                            )}

                            {newMethodType === "bank" && (
                                <div className="space-y-2">
                                    <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Bank Name</label>
                                    <input
                                        type="text"
                                        value={newMethodDetails.bankName}
                                        onChange={(e) => setNewMethodDetails({...newMethodDetails, bankName: e.target.value})}
                                        placeholder="e.g. Chase Bank"
                                        className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all"
                                    />
                                </div>
                            )}

                            {newMethodType === "crypto" && (
                                <div className="space-y-2">
                                    <label className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em] block px-1">Wallet Address</label>
                                    <textarea
                                        value={newMethodDetails.address}
                                        onChange={(e) => setNewMethodDetails({...newMethodDetails, address: e.target.value})}
                                        placeholder="0x..."
                                        rows={3}
                                        className="w-full bg-white/5 border border-white/10 py-4 px-6 rounded-2xl text-sm font-bold outline-none focus:border-primary transition-all resize-none font-mono"
                                    />
                                </div>
                            )}
                        </div>

                        <div className="flex gap-3 mt-8">
                            <button
                                onClick={() => setShowAddModal(false)}
                                className="flex-1 px-6 py-3 bg-white/5 border border-white/10 text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-white/10 transition-all"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleAdd}
                                disabled={!newMethodDetails.name.trim()}
                                className="flex-1 px-6 py-3 bg-primary text-white rounded-2xl text-xs font-black uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                Add Method
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Delete Confirmation Modal */}
            <ConfirmModal
                open={!!deleteTarget}
                title="Remove payment method?"
                description="This payment method will be removed from your account. This action cannot be undone."
                onConfirm={() => deleteTarget && handleDelete(deleteTarget)}
                onClose={() => setDeleteTarget(null)}
                confirmText="Remove"
                cancelText="Keep"
            />

            {/* Success/Error Message */}
            {message && (
                <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className={`p-4 rounded-2xl flex items-center space-x-3 ${
                        message.type === "success"
                            ? "bg-green-500/10 border border-green-500/30 text-green-400"
                            : "bg-red-500/10 border border-red-500/30 text-red-400"
                    }`}
                >
                    {message.type === "success" ? <CheckCircle className="w-5 h-5 flex-shrink-0" /> : <AlertCircle className="w-5 h-5 flex-shrink-0" />}
                    <span className="text-sm font-medium">{message.text}</span>
                </motion.div>
            )}
        </div>
    );
};

export default PaymentSettings;
