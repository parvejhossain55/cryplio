-- ============================================
-- Migration 014: Views
-- Common query views for easy reuse
-- ============================================

BEGIN;

-- Active trade ads with maker reputation info (SRS 3.2.2)
CREATE OR REPLACE VIEW active_trade_ads AS
SELECT 
    a.ad_id,
    a.type,
    ca.symbol AS crypto,
    fc.code AS fiat,
    a.price,
    a.min_amount,
    a.max_amount,
    a.payment_window_minutes,
    u.username AS maker_username,
    u.kyc_level,
    u.is_merchant,
    COALESCE(us.total_trades, 0) AS maker_total_trades,
    COALESCE(us.successful_trades, 0) AS maker_successful_trades,
    COALESCE(us.avg_rating, 0) AS maker_avg_rating,
    a.published_at,
    a.auto_repost
FROM trade_ads a
JOIN users u ON a.user_id = u.user_id
JOIN crypto_assets ca ON a.crypto_id = ca.id
JOIN fiat_currencies fc ON a.fiat_id = fc.id
LEFT JOIN user_stats us ON u.user_id = us.user_id
WHERE a.status = 'active'
    AND a.is_public = true
    AND a.deleted_at IS NULL
    AND (a.visibility_start_at IS NULL OR a.visibility_start_at <= NOW())
    AND (a.visibility_end_at IS NULL OR a.visibility_end_at >= NOW());

-- Completed trades with feedback status
CREATE OR REPLACE VIEW completed_trades AS
SELECT 
    t.trade_id,
    t.buyer_id,
    t.seller_id,
    t.crypto_amount,
    t.fiat_amount,
    t.status,
    t.completed_at,
    t.created_at,
    bf.feedback_id AS buyer_feedback_id,
    bf.rating AS buyer_rating,
    sf.feedback_id AS seller_feedback_id,
    sf.rating AS seller_rating
FROM trades t
LEFT JOIN trade_feedback bf ON t.trade_id = bf.trade_id AND bf.from_user_id = t.buyer_id
LEFT JOIN trade_feedback sf ON t.trade_id = sf.trade_id AND sf.from_user_id = t.seller_id
WHERE t.status IN ('completed', 'released');

-- User public profiles (for public display)
CREATE OR REPLACE VIEW public_user_profiles AS
SELECT 
    u.user_id,
    u.username,
    u.avatar_url,
    u.bio,
    u.kyc_level,
    u.is_merchant,
    u.created_at,
    COALESCE(us.total_trades, 0) AS total_trades,
    COALESCE(us.successful_trades, 0) AS successful_trades,
    COALESCE(us.avg_rating, 0) AS avg_rating,
    COALESCE(us.total_volume_usd, 0) AS total_volume_usd
FROM users u
LEFT JOIN user_stats us ON u.user_id = us.user_id
WHERE u.status = 'active' AND u.deleted_at IS NULL;

-- Open disputes requiring admin attention
CREATE OR REPLACE VIEW open_disputes AS
SELECT 
    d.dispute_id,
    d.trade_id,
    d.raised_by,
    d.reason_code,
    d.status,
    d.assigned_admin,
    t.buyer_id,
    t.seller_id,
    t.fiat_amount,
    t.crypto_amount,
    u.username AS raised_by_username,
    buyer.username AS buyer_username,
    seller.username AS seller_username
FROM disputes d
JOIN trades t ON d.trade_id = t.trade_id
JOIN users u ON d.raised_by = u.user_id
JOIN users buyer ON t.buyer_id = buyer.user_id
JOIN users seller ON t.seller_id = seller.user_id
WHERE d.status IN ('pending', 'assigned', 'under_review')
ORDER BY d.created_at ASC;

COMMIT;
