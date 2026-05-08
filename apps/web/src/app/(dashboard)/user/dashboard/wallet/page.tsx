"use client";

import React, { useEffect, useMemo, useState } from "react";
import { motion } from "framer-motion";
import DashboardLayout from "@/components/dashboard/DashboardLayout";
import {
    ArrowUpRight,
    ArrowDownLeft,
    History,
    Wallet as WalletIcon,
    RefreshCw,
    Copy,
    ExternalLink,
    AlertTriangle,
    Shield,
    Loader2,
    X,
    CheckCircle2,
    QrCode,
    Plus
} from "lucide-react";
import { authService, WalletBalance, WalletTransaction } from "@/services/authService";
import { toast } from "sonner";

const UserWallet = () => {
    const [wallets, setWallets] = useState<WalletBalance[]>([]);
    const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    
    // Modal states
    const [showDepositModal, setShowDepositModal] = useState(false);
    const [showWithdrawModal, setShowWithdrawModal] = useState(false);
    const [showAddWalletModal, setShowAddWalletModal] = useState(false);
    const [selectedWallet, setSelectedWallet] = useState<WalletBalance | null>(null);
    const [depositAddress, setDepositAddress] = useState<string>("");
    const [withdrawForm, setWithdrawForm] = useState({
        amount: "",
        address: "",
        twoFACode: ""
    });
    const [selectedCrypto, setSelectedCrypto] = useState("USDT");
    const [isSubmitting, setIsSubmitting] = useState(false);

    const loadWalletData = async () => {
        setLoading(true);
        setError(null);
        try {
            const [walletData, txData] = await Promise.all([
                authService.getWalletBalances(),
                authService.getWalletTransactions({ limit: 10, offset: 0 }),
            ]);
            setWallets(walletData);
            setTransactions(txData.transactions);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load wallet data");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void loadWalletData();
    }, []);

    const totalBalance = useMemo(
        () => wallets.reduce((sum, wallet) => sum + Number(wallet.balance || 0), 0),
        [wallets]
    );

    const handleDeposit = async (wallet: WalletBalance) => {
        setSelectedWallet(wallet);
        setShowDepositModal(true);
        
        try {
            const addressData = await authService.getDepositAddress(wallet.crypto_symbol);
            setDepositAddress(addressData.address);
        } catch (error: any) {
            toast.error("Failed to get deposit address");
            setShowDepositModal(false);
        }
    };

    const handleWithdraw = async (wallet: WalletBalance) => {
        if (Number(wallet.balance) <= 0) {
            toast.error("Insufficient balance");
            return;
        }
        
        setSelectedWallet(wallet);
        setShowWithdrawModal(true);
    };

    const handleWithdrawSubmit = async () => {
        if (!withdrawForm.amount || !withdrawForm.address || !withdrawForm.twoFACode) {
            toast.error("Please fill all fields");
            return;
        }

        if (Number(withdrawForm.amount) > Number(selectedWallet?.balance)) {
            toast.error("Insufficient balance");
            return;
        }

        setIsSubmitting(true);
        try {
            await authService.withdrawFunds({
                crypto_symbol: selectedWallet?.crypto_symbol || "USDT",
                amount: Number(withdrawForm.amount),
                address: withdrawForm.address,
                two_fa_code: withdrawForm.twoFACode
            });

            toast.success("Withdrawal request submitted successfully");
            setShowWithdrawModal(false);
            setWithdrawForm({ amount: "", address: "", twoFACode: "" });
            loadWalletData();
        } catch (error: any) {
            toast.error(error.message || "Failed to submit withdrawal");
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleAddWallet = async () => {
        if (!selectedCrypto) {
            toast.error("Please select a cryptocurrency");
            return;
        }

        setIsSubmitting(true);
        try {
            await authService.createWallet(selectedCrypto);
            toast.success(`${selectedCrypto} wallet created successfully`);
            setShowAddWalletModal(false);
            loadWalletData();
        } catch (error: any) {
            toast.error(error.message || "Failed to create wallet");
        } finally {
            setIsSubmitting(false);
        }
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        toast.success("Copied to clipboard");
    };

    return (
        <DashboardLayout title="Universal Wallet" role="user">
            <div className="space-y-8">
                <div className="bg-surface rounded-[2.5rem] border border-white/10 p-8">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                            <div className="w-12 h-12 bg-primary rounded-2xl flex items-center justify-center">
                                <WalletIcon className="text-white w-6 h-6" />
                            </div>
                            <div>
                                <p className="text-[10px] font-black text-text-dim uppercase tracking-[0.2em]">Total Balance</p>
                                <h2 className="text-3xl font-black text-white">{totalBalance.toFixed(8)}</h2>
                            </div>
                        </div>
                        <button
                            onClick={() => void loadWalletData()}
                            className="px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-xs font-bold flex items-center gap-2"
                        >
                            <RefreshCw className="w-4 h-4" /> Refresh
                        </button>
                    </div>
                </div>

                {error && (
                    <div className="bg-primary/10 border border-primary/30 rounded-2xl p-4 text-sm text-primary">{error}</div>
                )}

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    <div className="bg-surface border border-white/10 rounded-[2rem] p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-black text-white uppercase">Wallet Balances</h3>
                            <div className="flex items-center gap-2">
                                <button
                                    onClick={() => setShowAddWalletModal(true)}
                                    className="flex items-center gap-2 px-3 py-2 bg-gradient-to-r from-accent/20 to-accent/10 hover:from-accent/30 hover:to-accent/20 text-accent rounded-xl transition-all duration-200 border border-accent/20 hover:border-accent/40"
                                    title="Add Wallet"
                                >
                                    <Plus className="w-4 h-4" />
                                    <span className="text-xs font-bold">Add</span>
                                </button>
                                <WalletIcon className="w-4 h-4 text-text-dim" />
                            </div>
                        </div>
                        <div className="space-y-4">
                            {loading && <p className="text-text-dim text-sm">Loading balances...</p>}
                            {!loading && wallets.length === 0 && <p className="text-text-dim text-sm">No wallets found.</p>}
                            {wallets.map((wallet) => (
                                <div key={wallet.wallet_id} className="p-4 rounded-xl bg-white/5 border border-white/5">
                                    <div className="flex items-center justify-between mb-3">
                                        <div>
                                            <span className="text-sm font-bold text-white">{wallet.crypto_symbol || "USDT"}</span>
                                            <span className={`ml-2 text-xs px-2 py-0.5 rounded-full ${
                                                wallet.is_active ? "bg-green-500/20 text-green-500" : "bg-red-500/20 text-red-500"
                                            }`}>
                                                {wallet.is_active ? "Active" : "Inactive"}
                                            </span>
                                        </div>
                                    </div>
                                    <div className="flex items-center justify-between mb-3">
                                        <div>
                                            <p className="text-xl font-black text-white">{Number(wallet.balance).toFixed(4)}</p>
                                            <p className="text-xs text-text-dim">Available</p>
                                        </div>
                                        {Number(wallet.locked_balance) > 0 && (
                                            <div className="text-right">
                                                <p className="text-sm font-bold text-orange-500">{Number(wallet.locked_balance).toFixed(4)}</p>
                                                <p className="text-xs text-text-dim">In Escrow</p>
                                            </div>
                                        )}
                                    </div>
                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => handleDeposit(wallet)}
                                            className="flex-1 py-2.5 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-xl text-xs font-black uppercase tracking-wide hover:from-emerald-400 hover:to-emerald-500 transition-all duration-200 flex items-center justify-center gap-1.5 shadow-lg shadow-emerald-500/20 hover:shadow-emerald-500/30"
                                        >
                                            <ArrowDownLeft className="w-3.5 h-3.5" />
                                            Deposit
                                        </button>
                                        <button
                                            onClick={() => handleWithdraw(wallet)}
                                            disabled={Number(wallet.balance) <= 0}
                                            className="flex-1 py-2.5 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-xl text-xs font-black uppercase tracking-wide hover:from-orange-400 hover:to-orange-500 transition-all duration-200 disabled:opacity-40 disabled:cursor-not-allowed disabled:shadow-none flex items-center justify-center gap-1.5 shadow-lg shadow-orange-500/20 hover:shadow-orange-500/30"
                                        >
                                            <ArrowUpRight className="w-3.5 h-3.5" />
                                            Withdraw
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="bg-surface border border-white/10 rounded-[2rem] p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="text-lg font-black text-white uppercase">Recent Activity</h3>
                            <History className="w-4 h-4 text-text-dim" />
                        </div>
                        <div className="space-y-3">
                            {loading && <p className="text-text-dim text-sm">Loading transactions...</p>}
                            {!loading && transactions.length === 0 && <p className="text-text-dim text-sm">No transactions yet.</p>}
                            {transactions.map((tx) => (
                                <div key={tx.tx_id} className="flex items-center justify-between p-3 rounded-xl bg-white/5 border border-white/5">
                                    <div className="flex items-center gap-3">
                                        {tx.type === "withdrawal" ? (
                                            <ArrowUpRight className="w-4 h-4 text-primary" />
                                        ) : (
                                            <ArrowDownLeft className="w-4 h-4 text-accent" />
                                        )}
                                        <div>
                                            <p className="text-xs font-bold text-white uppercase">{tx.type}</p>
                                            <p className="text-[10px] text-text-dim">{new Date(tx.created_at).toLocaleString()}</p>
                                        </div>
                                    </div>
                                    <p className="text-xs font-black text-white">{Number(tx.amount).toFixed(8)}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>

            {/* Deposit Modal */}
            {showDepositModal && selectedWallet && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
                    onClick={() => setShowDepositModal(false)}
                >
                    <motion.div
                        initial={{ scale: 0.95, opacity: 0 }}
                        animate={{ scale: 1, opacity: 1 }}
                        className="bg-surface border border-white/10 rounded-2xl p-6 max-w-md w-full"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-xl font-black text-white">Deposit {selectedWallet.crypto_symbol || "USDT"}</h3>
                            <button
                                onClick={() => setShowDepositModal(false)}
                                className="text-text-dim hover:text-white transition-colors"
                            >
                                <X className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="space-y-4">
                            <div className="bg-white/5 border border-white/10 rounded-xl p-4">
                                <div className="flex items-center justify-between mb-2">
                                    <span className="text-sm font-bold text-white">Your Deposit Address</span>
                                    <button
                                        onClick={() => copyToClipboard(depositAddress)}
                                        className="p-2 text-primary hover:bg-primary/20 rounded-lg transition-colors"
                                    >
                                        <Copy className="w-4 h-4" />
                                    </button>
                                </div>
                                <p className="text-xs text-text-dim font-mono break-all">{depositAddress}</p>
                            </div>

                            <div className="bg-yellow-500/10 border border-yellow-500/20 rounded-xl p-4">
                                <div className="flex items-start gap-3">
                                    <AlertTriangle className="w-5 h-5 text-yellow-500 flex-shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-sm font-bold text-yellow-500 mb-1">Important Notice</p>
                                        <ul className="text-xs text-yellow-500/80 space-y-1">
                                            <li>• Only send {selectedWallet.crypto_symbol || "USDT"} to this address</li>
                                            <li>• Minimum deposit: 10 USDT</li>
                                            <li>• Deposits require 12 network confirmations</li>
                                            <li>• Sending other tokens may result in permanent loss</li>
                                        </ul>
                                    </div>
                                </div>
                            </div>

                            <div className="flex items-center gap-2 text-xs text-text-dim">
                                <Shield className="w-4 h-4 text-accent" />
                                <span>All deposits are protected by multi-signature escrow</span>
                            </div>
                        </div>
                    </motion.div>
                </motion.div>
            )}

            {/* Withdrawal Modal */}
            {showWithdrawModal && selectedWallet && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
                    onClick={() => setShowWithdrawModal(false)}
                >
                    <motion.div
                        initial={{ scale: 0.95, opacity: 0 }}
                        animate={{ scale: 1, opacity: 1 }}
                        className="bg-surface border border-white/10 rounded-2xl p-6 max-w-md w-full"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-xl font-black text-white">Withdraw {selectedWallet.crypto_symbol || "USDT"}</h3>
                            <button
                                onClick={() => setShowWithdrawModal(false)}
                                className="text-text-dim hover:text-white transition-colors"
                            >
                                <X className="w-5 h-5" />
                            </button>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-bold text-white mb-2">Amount</label>
                                <div className="relative">
                                    <input
                                        type="number"
                                        value={withdrawForm.amount}
                                        onChange={(e) => setWithdrawForm({...withdrawForm, amount: e.target.value})}
                                        placeholder="0.0000"
                                        step="0.0001"
                                        min="0.0001"
                                        max={Number(selectedWallet.balance)}
                                        className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50"
                                    />
                                    <div className="absolute right-4 top-1/2 -translate-y-1/2 text-xs text-text-dim">
                                        Available: {Number(selectedWallet.balance).toFixed(4)}
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-bold text-white mb-2">Recipient Address</label>
                                <input
                                    type="text"
                                    value={withdrawForm.address}
                                    onChange={(e) => setWithdrawForm({...withdrawForm, address: e.target.value})}
                                    placeholder="0x..."
                                    className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50 font-mono text-sm"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-bold text-white mb-2">2FA Code</label>
                                <input
                                    type="text"
                                    value={withdrawForm.twoFACode}
                                    onChange={(e) => setWithdrawForm({...withdrawForm, twoFACode: e.target.value})}
                                    placeholder="000000"
                                    maxLength={6}
                                    className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-white placeholder-text-dim focus:outline-none focus:border-primary/50 font-mono text-center text-lg"
                                />
                            </div>

                            <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-4">
                                <div className="flex items-start gap-3">
                                    <AlertTriangle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
                                    <div>
                                        <p className="text-sm font-bold text-red-500 mb-1">Withdrawal Warning</p>
                                        <ul className="text-xs text-red-500/80 space-y-1">
                                            <li>• Withdrawals are irreversible</li>
                                            <li>• Daily limit: $500 USD equivalent</li>
                                            <li>• Network fees will be deducted</li>
                                            <li>• Double-check the recipient address</li>
                                        </ul>
                                    </div>
                                </div>
                            </div>

                            <button
                                onClick={handleWithdrawSubmit}
                                disabled={isSubmitting || !withdrawForm.amount || !withdrawForm.address || !withdrawForm.twoFACode}
                                className="w-full py-3 bg-primary text-white rounded-xl font-black uppercase tracking-wider text-sm hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                            >
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="w-4 h-4 animate-spin" />
                                        Processing...
                                    </>
                                ) : (
                                    <>
                                        <ArrowUpRight className="w-4 h-4" />
                                        Withdraw Funds
                                    </>
                                )}
                            </button>
                        </div>
                    </motion.div>
                </motion.div>
            )}

            {/* Add Wallet Modal */}
            {showAddWalletModal && (
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex items-center justify-center p-4"
                    onClick={() => setShowAddWalletModal(false)}
                >
                    <motion.div
                        initial={{ scale: 0.9, opacity: 0, y: 20 }}
                        animate={{ scale: 1, opacity: 1, y: 0 }}
                        transition={{ type: "spring", damping: 25, stiffness: 300 }}
                        className="bg-gradient-to-br from-surface to-surface/95 border border-white/10 rounded-3xl p-8 max-w-md w-full shadow-2xl shadow-black/50"
                        onClick={(e) => e.stopPropagation()}
                    >
                        {/* Header with Icon */}
                        <div className="flex items-center justify-between mb-8">
                            <div className="flex items-center gap-3">
                                <div className="w-12 h-12 bg-gradient-to-br from-accent/20 to-accent/5 rounded-2xl flex items-center justify-center border border-accent/20">
                                    <WalletIcon className="w-6 h-6 text-accent" />
                                </div>
                                <div>
                                    <h3 className="text-xl font-black text-white">New Wallet</h3>
                                    <p className="text-xs text-text-dim">Create a crypto wallet</p>
                                </div>
                            </div>
                            <button
                                onClick={() => setShowAddWalletModal(false)}
                                className="p-2 hover:bg-white/5 rounded-xl transition-colors"
                            >
                                <X className="w-5 h-5 text-text-dim hover:text-white" />
                            </button>
                        </div>

                        <div className="space-y-6">
                            {/* Crypto Selection Cards */}
                            <div>
                                <label className="block text-xs font-bold text-text-dim uppercase tracking-wider mb-3">Select Asset</label>
                                <div className="grid grid-cols-1 gap-2">
                                    {[
                                        { symbol: "USDT", name: "Tether", color: "#26A17B", desc: "Stablecoin" },
                                        { symbol: "BTC", name: "Bitcoin", color: "#F7931A", desc: "Digital gold" },
                                        { symbol: "ETH", name: "Ethereum", color: "#627EEA", desc: "Smart contracts" },
                                    ].map((crypto) => (
                                        <button
                                            key={crypto.symbol}
                                            onClick={() => setSelectedCrypto(crypto.symbol)}
                                            className={`flex items-center gap-4 p-4 rounded-2xl border transition-all duration-200 ${
                                                selectedCrypto === crypto.symbol
                                                    ? "bg-white/10 border-accent/50 shadow-lg shadow-accent/10"
                                                    : "bg-white/5 border-white/5 hover:bg-white/10 hover:border-white/20"
                                            }`}
                                        >
                                            <div
                                                className="w-12 h-12 rounded-xl flex items-center justify-center text-white font-black text-lg"
                                                style={{ backgroundColor: `${crypto.color}20`, color: crypto.color }}
                                            >
                                                {crypto.symbol[0]}
                                            </div>
                                            <div className="flex-1 text-left">
                                                <div className="flex items-center gap-2">
                                                    <span className="font-bold text-white">{crypto.symbol}</span>
                                                    {selectedCrypto === crypto.symbol && (
                                                        <div className="w-2 h-2 bg-accent rounded-full" />
                                                    )}
                                                </div>
                                                <p className="text-xs text-text-dim">{crypto.name} • {crypto.desc}</p>
                                            </div>
                                        </button>
                                    ))}
                                </div>
                            </div>

                            {/* Info Card */}
                            <div className="bg-gradient-to-r from-yellow-500/10 to-orange-500/10 border border-yellow-500/20 rounded-2xl p-4">
                                <div className="flex items-start gap-3">
                                    <div className="p-2 bg-yellow-500/20 rounded-xl">
                                        <AlertTriangle className="w-4 h-4 text-yellow-500" />
                                    </div>
                                    <div>
                                        <p className="text-sm font-bold text-yellow-500 mb-1">One per Asset</p>
                                        <p className="text-xs text-yellow-500/70 leading-relaxed">
                                            You can create only one wallet per cryptocurrency. The wallet will be ready for deposits and withdrawals immediately.
                                        </p>
                                    </div>
                                </div>
                            </div>

                            {/* Action Button */}
                            <button
                                onClick={handleAddWallet}
                                disabled={isSubmitting}
                                className="w-full py-4 bg-gradient-to-r from-accent to-accent/80 text-white rounded-2xl font-black uppercase tracking-wider text-sm hover:from-accent/90 hover:to-accent/70 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 shadow-lg shadow-accent/20 hover:shadow-accent/30"
                            >
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="w-5 h-5 animate-spin" />
                                        <span>Creating...</span>
                                    </>
                                ) : (
                                    <>
                                        <Plus className="w-5 h-5" />
                                        <span>Create {selectedCrypto} Wallet</span>
                                    </>
                                )}
                            </button>
                        </div>
                    </motion.div>
                </motion.div>
            )}
        </DashboardLayout>
    );
};

export default UserWallet;
