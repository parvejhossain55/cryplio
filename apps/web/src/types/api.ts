export interface BackendUser {
    id: string;
    email: string;
    username: string;
    bio?: string;
    email_verified: boolean;
    is_merchant: boolean;
    two_fa_enabled: boolean;
    avatar_url?: string | null;
    is_online: boolean;
}

export interface BackendSession {
    id?: string;
    token_id: string;
    device_type?: string;
    ip_address?: string;
    last_used_at: string;
}

export interface UserStats {
    total_trades: number;
    successful_trades: number;
    dispute_rate: number;
    avg_rating?: number;
    positive_feedback_count: number;
    neutral_feedback_count: number;
    negative_feedback_count: number;
    total_volume_usd: number;
    last_trade_at?: string;
    success_rate: number;
}

export interface UserPaymentMethod {
    id: string;
    user_id: string;
    payment_method_code: string;
    display_name: string;
    account_name?: string;
    account_number?: string;
    bank_name?: string;
    is_active: boolean;
    is_default: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreatePaymentMethodRequest {
    payment_method_code: string;
    display_name: string;
    account_name?: string;
    account_number?: string;
    bank_name?: string;
}

export interface AdResponse {
    ad_id: string;
    user_id: string;
    username: string;
    user_avatar?: string;
    user_rating: number;
    user_trades: number;
    type: "buy" | "sell";
    crypto_symbol: string;
    fiat_symbol: string;
    price_type: string;
    price: number;
    min_amount: number;
    max_amount: number;
    payment_methods: string[];
    payment_window_minutes: number;
    is_online: boolean;
    trade_terms?: string;
}

export interface WalletBalance {
    wallet_id: string;
    user_id: string;
    crypto_id: number | null;
    crypto_symbol: string | null;
    address: string;
    balance: number;
    locked_balance: number;
    is_active: boolean;
}

export interface WalletTransaction {
    tx_id: string;
    wallet_id: string;
    type: string;
    status: string;
    amount: number;
    fee: number;
    created_at: string;
}

export interface AuthResponse {
    user: BackendUser;
    token?: string;
    refresh_token?: string;
    requires_2fa?: boolean;
    temp_token?: string;
    error?: string;
}

export interface HeaderProfileResponse {
    username: string;
    avatar_url?: string | null;
    is_online: boolean;
    trader_badge: string;
    unread_notification_count: number;
}
