package trading

import (
	"cryplio/internal/domain/trading"
)

// rowScanner is satisfied by both *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

// tradeColumns is the canonical SELECT column list matching scanTrade's Scan order.
const tradeColumns = `
	trade_id, ad_id, buyer_id, seller_id, crypto_id, fiat_id,
	crypto_amount, fiat_amount, exchange_rate, payment_method,
	price_type, agreed_price, status, payment_details,
	payment_window_minutes, started_at, payment_marked_at,
	released_at, cancelled_at, completed_at, expired_at,
	cancel_reason, created_at, updated_at`

// scanTrade scans a single row into a Trade value.
func scanTrade(row rowScanner, t *trading.Trade) error {
	var pd []byte
	err := row.Scan(
		&t.TradeID, &t.AdID, &t.BuyerID, &t.SellerID, &t.CryptoID, &t.FiatID,
		&t.CryptoAmount, &t.FiatAmount, &t.ExchangeRate, &t.PaymentMethod,
		&t.PriceType, &t.AgreedPrice, &t.Status, &pd,
		&t.PaymentWindowMinutes, &t.StartedAt, &t.PaymentMarkedAt,
		&t.ReleasedAt, &t.CancelledAt, &t.CompletedAt, &t.ExpiredAt,
		&t.CancelReason, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if len(pd) > 0 {
		t.PaymentDetails = pd
	}
	return nil
}
