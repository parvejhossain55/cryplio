package trading

import (
	"context"
	"fmt"

	"cryplio/internal/domain/trading"

	"github.com/google/uuid"
)

// ─── Trade Chat Messages ──────────────────────────────────────────────────────

func (r *tradeRepository) CreateTradeMessage(ctx context.Context, msg *trading.TradeMessage) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO trade_messages (
			message_id, trade_id, sender_id, message_type, content, file_url, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING created_at`,
		msg.MessageID, msg.TradeID, msg.SenderID, msg.MessageType, msg.Content, msg.FileURL,
	).Scan(&msg.CreatedAt)
}

func (r *tradeRepository) ListTradeMessages(ctx context.Context, tradeID uuid.UUID) ([]trading.TradeMessage, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT message_id, trade_id, sender_id, message_type, content,
		       file_url, is_read, read_at, created_at
		FROM trade_messages
		WHERE trade_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC`, tradeID)
	if err != nil {
		return nil, fmt.Errorf("list trade messages: %w", err)
	}
	defer rows.Close()

	var msgs []trading.TradeMessage
	for rows.Next() {
		var m trading.TradeMessage
		if err := rows.Scan(
			&m.MessageID, &m.TradeID, &m.SenderID, &m.MessageType,
			&m.Content, &m.FileURL, &m.IsRead, &m.ReadAt, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan trade message: %w", err)
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
