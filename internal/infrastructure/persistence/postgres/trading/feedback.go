package trading

import (
	"context"
	"database/sql"
	"fmt"

	"cryplio/internal/domain/trading"

	"github.com/google/uuid"
)

// ─── Trade Feedback ───────────────────────────────────────────────────────────

func (r *tradeRepository) CreateFeedback(ctx context.Context, fb *trading.TradeFeedback) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO trade_feedback (
			feedback_id, trade_id, from_user_id, to_user_id,
			rating, comment, is_public, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, true, NOW())
		RETURNING created_at`,
		fb.FeedbackID, fb.TradeID, fb.FromUserID, fb.ToUserID, fb.Rating, fb.Comment,
	).Scan(&fb.CreatedAt)
}

func (r *tradeRepository) GetFeedbackByTrade(ctx context.Context, tradeID uuid.UUID) (*trading.TradeFeedback, error) {
	var fb trading.TradeFeedback
	err := r.db.QueryRowContext(ctx, `
		SELECT feedback_id, trade_id, from_user_id, to_user_id, rating, comment, created_at
		FROM trade_feedback
		WHERE trade_id = $1`, tradeID,
	).Scan(&fb.FeedbackID, &fb.TradeID, &fb.FromUserID, &fb.ToUserID, &fb.Rating, &fb.Comment, &fb.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get feedback by trade: %w", err)
	}
	return &fb, nil
}
