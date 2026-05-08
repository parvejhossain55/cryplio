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
	crypto_amount, fiat_amount, exchange_rate, payment_method,
	price_type, agreed_price, status, dispute_id, chat_room_id,
	started_at, payment_marked_at, released_at, cancelled_at,
	completed_at, expired_at, payment_window_minutes,
	is_auto_dispute_triggered, cancel_reason, escrow_txn_hash,
	escrow_contract_address, created_at, updated_at, deleted_at
`

// scanTrade scans the 29-column trade projection into t.
func scanTrade(row scanner, t *trading.Trade) error {
	return row.Scan(
		&t.TradeID, &t.AdID, &t.BuyerID, &t.SellerID, &t.CryptoID, &t.FiatID,
		&t.CryptoAmount, &t.FiatAmount, &t.ExchangeRate, &t.PaymentMethod,
		&t.PriceType, &t.AgreedPrice, &t.Status, &t.DisputeID, &t.ChatRoomID,
		&t.StartedAt, &t.PaymentMarkedAt, &t.ReleasedAt, &t.CancelledAt,
		&t.CompletedAt, &t.ExpiredAt, &t.PaymentWindowMinutes,
		&t.IsAutoDisputeTriggered, &t.CancelReason, &t.EscrowTxnHash,
		&t.EscrowContractAddress, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
	)
}
