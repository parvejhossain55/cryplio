package trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// SendMessage posts a text message to the trade chat on behalf of senderID.
func (s *tradeService) SendMessage(ctx context.Context, tradeID, senderID uuid.UUID, content string) (*TradeMessage, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != senderID && trade.SellerID != senderID {
		return nil, errors.New("unauthorized")
	}

	msg := &TradeMessage{
		ID:       uuid.New(),
		TradeID:  tradeID,
		SenderID: senderID,
		Message:  content,
		IsSystem: false,
	}

	if err = s.tradeRepo.CreateTradeMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// SendFileMessage posts a file (or image) message to the trade chat on behalf of senderID.
// The message type is set to TradeMessageTypeImage when mimeType has an "image/" prefix,
// otherwise TradeMessageTypeFile.
func (s *tradeService) SendFileMessage(ctx context.Context, tradeID, senderID uuid.UUID, fileURL, mimeType string, fileSize int) (*TradeMessage, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != senderID && trade.SellerID != senderID {
		return nil, errors.New("unauthorized")
	}

	msg := &TradeMessage{
		ID:       uuid.New(),
		TradeID:  tradeID,
		SenderID: senderID,
		Message:  fmt.Sprintf("[File: %s]", fileURL),
		IsSystem: false,
	}

	if err = s.tradeRepo.CreateTradeMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// GetChatHistory returns the full message history for a trade, gated on userID being a participant.
func (s *tradeService) GetChatHistory(ctx context.Context, tradeID, userID uuid.UUID) ([]TradeMessage, error) {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return nil, err
	}
	if trade == nil {
		return nil, errors.New("trade not found")
	}
	if trade.BuyerID != userID && trade.SellerID != userID {
		return nil, errors.New("unauthorized")
	}

	return s.tradeRepo.ListTradeMessages(ctx, tradeID)
}

// LeaveFeedback records a rating and optional comment for a completed trade.
// Each trade may only receive one feedback entry.
func (s *tradeService) LeaveFeedback(ctx context.Context, tradeID, userID uuid.UUID, rating FeedbackRating, comment string) error {
	trade, err := s.tradeRepo.GetTradeByID(ctx, tradeID)
	if err != nil {
		return err
	}
	if trade == nil {
		return errors.New("trade not found")
	}
	if trade.Status != TradeStatusCompleted {
		return errors.New("can only leave feedback on completed trades")
	}
	if trade.BuyerID != userID && trade.SellerID != userID {
		return errors.New("unauthorized")
	}

	// Check if feedback already exists.
	existing, err := s.tradeRepo.GetFeedbackByTrade(ctx, tradeID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("feedback already exists for this trade")
	}

	// Determine recipient: the other party.
	var recipientID uuid.UUID
	if trade.BuyerID == userID {
		recipientID = trade.SellerID
	} else {
		recipientID = trade.BuyerID
	}

	feedback := &TradeFeedback{
		FeedbackID: uuid.New(),
		TradeID:    tradeID,
		FromUserID: userID,
		ToUserID:   recipientID,
		Rating:     rating,
		Comment:    &comment,
	}

	if err := s.tradeRepo.CreateFeedback(ctx, feedback); err != nil {
		return fmt.Errorf("create feedback: %w", err)
	}

	return nil
}
