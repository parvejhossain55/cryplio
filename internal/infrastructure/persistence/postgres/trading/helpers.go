package trading

import (
	"cryplio/internal/domain/trading"
)

// scanner is satisfied by *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

// tradeColumns is the canonical SELECT column list for the trades table.
const tradeColumns = `
	trade_id, ad_id, buyer_id, seller_id, crypto_id, fiat_id,
	crypto_amount, fiat_amount, rate, status, payment_method_code,
	escrow_wallet_id, tx_hash, payment_window_minutes,
	expires_at, paid_at, released_at, completed_at, cancelled_at,
	disputed_at, created_at, updated_at
`

// scanTrade scans the trade projection into t.
func scanTrade(row scanner, t *trading.Trade) error {
	return row.Scan(
		&t.TradeID, &t.AdID, &t.BuyerID, &t.SellerID, &t.CryptoID, &t.FiatID,
		&t.CryptoAmount, &t.FiatAmount, &t.Rate, &t.Status, &t.PaymentMethodCode,
		&t.EscrowWalletID, &t.TxHash, &t.PaymentWindowMinutes,
		&t.ExpiresAt, &t.PaidAt, &t.ReleasedAt, &t.CompletedAt, &t.CancelledAt,
		&t.DisputedAt, &t.CreatedAt, &t.UpdatedAt,
	)
}
