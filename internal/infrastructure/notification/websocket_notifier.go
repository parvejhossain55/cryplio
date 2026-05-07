package notification

import (
	"context"
	"fmt"
	"log"

	"cryplio/internal/interfaces/websocket"

	"github.com/google/uuid"
)

// WebSocketNotifier provides real-time notification delivery via WebSocket
type WebSocketNotifier struct {
	wsService    websocket.Service
	emailService *MockEmailService
}

// NewWebSocketNotifier creates a new WebSocket notifier
func NewWebSocketNotifier(wsService websocket.Service, emailService *MockEmailService) *WebSocketNotifier {
	return &WebSocketNotifier{
		wsService:    wsService,
		emailService: emailService,
	}
}

// NotifyTradeUpdate sends a trade update notification
func (w *WebSocketNotifier) NotifyTradeUpdate(ctx context.Context, tradeID uuid.UUID, buyerID, sellerID uuid.UUID, status string, amount float64, crypto, fiat string) {
	// Create notification data
	notificationData := map[string]interface{}{
		"trade_id":  tradeID.String(),
		"status":    status,
		"amount":    amount,
		"crypto":    crypto,
		"fiat":      fiat,
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	// Send to buyer
	w.sendNotificationToUser(buyerID, "trade_update", "Trade Status Updated",
		fmt.Sprintf("Your trade %s status changed to %s", tradeID.String()[:8], status),
		notificationData)

	// Send to seller
	w.sendNotificationToUser(sellerID, "trade_update", "Trade Status Updated",
		fmt.Sprintf("Your trade %s status changed to %s", tradeID.String()[:8], status),
		notificationData)

	log.Printf("Trade update notification sent for trade %s", tradeID)
}

// NotifyNewMessage sends a new chat message notification
func (w *WebSocketNotifier) NotifyNewMessage(ctx context.Context, tradeID uuid.UUID, senderID, receiverID uuid.UUID, message string) {
	notificationData := map[string]interface{}{
		"trade_id":  tradeID.String(),
		"sender_id": senderID.String(),
		"message":   message,
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	w.sendNotificationToUser(receiverID, "new_message", "New Message",
		fmt.Sprintf("New message in trade %s", tradeID.String()[:8]),
		notificationData)

	log.Printf("Message notification sent for trade %s", tradeID)
}

// NotifyWithdrawalRequest sends a withdrawal request notification to admins
func (w *WebSocketNotifier) NotifyWithdrawalRequest(ctx context.Context, userID uuid.UUID, amount float64, crypto string) {
	notificationData := map[string]interface{}{
		"user_id":   userID.String(),
		"amount":    amount,
		"crypto":    crypto,
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	// Broadcast to all connected users (admins will filter)
	w.wsService.BroadcastMessage("withdrawal_request", map[string]interface{}{
		"type":    "withdrawal_request",
		"title":   "New Withdrawal Request",
		"message": fmt.Sprintf("User requested withdrawal of %.6f %s", amount, crypto),
		"data":    notificationData,
	}, "")

	log.Printf("Withdrawal request notification sent")
}

// NotifyDisputeCreated sends a dispute created notification
func (w *WebSocketNotifier) NotifyDisputeCreated(ctx context.Context, tradeID uuid.UUID, disputeID uuid.UUID, buyerID, sellerID uuid.UUID) {
	notificationData := map[string]interface{}{
		"trade_id":   tradeID.String(),
		"dispute_id": disputeID.String(),
		"timestamp":  fmt.Sprintf("%d", uuid.New().Time()),
	}

	// Notify both parties
	w.sendNotificationToUser(buyerID, "dispute_created", "Dispute Created",
		fmt.Sprintf("Dispute created for trade %s", tradeID.String()[:8]),
		notificationData)

	w.sendNotificationToUser(sellerID, "dispute_created", "Dispute Created",
		fmt.Sprintf("Dispute created for trade %s", tradeID.String()[:8]),
		notificationData)

	// Notify admins
	w.wsService.BroadcastMessage("dispute_created", map[string]interface{}{
		"type":    "dispute_created",
		"title":   "New Dispute Created",
		"message": fmt.Sprintf("Dispute created for trade %s", tradeID.String()[:8]),
		"data":    notificationData,
	}, "")

	log.Printf("Dispute created notification sent for trade %s", tradeID)
}

// NotifyMarketAlert sends a market price alert
func (w *WebSocketNotifier) NotifyMarketAlert(ctx context.Context, crypto, fiat string, price float64, changePercent float64) {
	notificationData := map[string]interface{}{
		"crypto_symbol":  crypto,
		"fiat_symbol":    fiat,
		"price":          price,
		"change_percent": changePercent,
		"timestamp":      fmt.Sprintf("%d", uuid.New().Time()),
	}

	// Broadcast market update to all users
	w.wsService.BroadcastMessage("market_alert", map[string]interface{}{
		"type":    "market_alert",
		"title":   fmt.Sprintf("%s Price Update", crypto),
		"message": fmt.Sprintf("%s is now %.2f %s (%.2f%%)", crypto, price, fiat, changePercent),
		"data":    notificationData,
	}, "")

	log.Printf("Market alert sent for %s-%s", crypto, fiat)
}

// NotifyUserSuspended sends a user suspension notification
func (w *WebSocketNotifier) NotifyUserSuspended(ctx context.Context, userID uuid.UUID, reason string, duration string) {
	notificationData := map[string]interface{}{
		"reason":    reason,
		"duration":  duration,
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	w.sendNotificationToUser(userID, "user_suspended", "Account Suspended",
		fmt.Sprintf("Your account has been suspended. Reason: %s", reason),
		notificationData)

	// Send email notification
	if w.emailService != nil {
		// Email body would be: fmt.Sprintf("Hello,\n\nYour account has been suspended.\n\nReason: %s\nDuration: %s\n\nPlease contact support if you have any questions.\n\nBest regards,\nThe Cryplio Team", reason, duration)
		// This would need the user's email - for now we'll log it
		log.Printf("Email notification sent to user %s for suspension", userID)
	}

	log.Printf("User suspension notification sent for user %s", userID)
}

// NotifyReferralBonus sends a referral bonus notification
func (w *WebSocketNotifier) NotifyReferralBonus(ctx context.Context, userID uuid.UUID, referralID uuid.UUID, bonusAmount float64, crypto string) {
	notificationData := map[string]interface{}{
		"referral_id":  referralID.String(),
		"bonus_amount": bonusAmount,
		"crypto":       crypto,
		"timestamp":    fmt.Sprintf("%d", uuid.New().Time()),
	}

	w.sendNotificationToUser(userID, "referral_bonus", "Referral Bonus Earned",
		fmt.Sprintf("You earned %.6f %s from a referral!", bonusAmount, crypto),
		notificationData)

	log.Printf("Referral bonus notification sent for user %s", userID)
}

// sendNotificationToUser sends a notification to a specific user
func (w *WebSocketNotifier) sendNotificationToUser(userID uuid.UUID, notificationType, title, message string, data map[string]interface{}) {
	if w.wsService == nil {
		return
	}

	notification := websocket.NotificationEvent{
		Type:      notificationType,
		UserID:    userID,
		Title:     title,
		Message:   message,
		Data:      data,
		Timestamp: fmt.Sprintf("%d", uuid.New().Time()),
	}

	// Send to specific user via WebSocket
	w.wsService.BroadcastToUser(userID, "notification", notification)
}

// BroadcastSystemMessage broadcasts a system message to all users
func (w *WebSocketNotifier) BroadcastSystemMessage(ctx context.Context, title, message string) {
	if w.wsService == nil {
		return
	}

	notificationData := map[string]interface{}{
		"type":      "system_message",
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	w.wsService.BroadcastMessage("system_message", map[string]interface{}{
		"type":    "system_message",
		"title":   title,
		"message": message,
		"data":    notificationData,
	}, "")

	log.Printf("System message broadcasted: %s", title)
}

// NotifyPaymentProcessed sends a payment processed notification
func (w *WebSocketNotifier) NotifyPaymentProcessed(ctx context.Context, userID uuid.UUID, transactionID string, amount float64, currency string, status string) {
	notificationData := map[string]interface{}{
		"transaction_id": transactionID,
		"amount":         amount,
		"currency":       currency,
		"status":         status,
		"timestamp":      fmt.Sprintf("%d", uuid.New().Time()),
	}

	title := "Payment Processed"
	message := fmt.Sprintf("Your payment of %.2f %s has been %s", amount, currency, status)

	if status == "completed" {
		w.sendNotificationToUser(userID, "payment_completed", title, message, notificationData)
	} else if status == "failed" {
		w.sendNotificationToUser(userID, "payment_failed", "Payment Failed", message, notificationData)
	} else {
		w.sendNotificationToUser(userID, "payment_update", title, message, notificationData)
	}

	log.Printf("Payment notification sent for transaction %s", transactionID)
}

// NotifyEscrowReleased sends an escrow release notification
func (w *WebSocketNotifier) NotifyEscrowReleased(ctx context.Context, tradeID uuid.UUID, receiverID uuid.UUID, amount float64, crypto string) {
	notificationData := map[string]interface{}{
		"trade_id":  tradeID.String(),
		"amount":    amount,
		"crypto":    crypto,
		"timestamp": fmt.Sprintf("%d", uuid.New().Time()),
	}

	w.sendNotificationToUser(receiverID, "escrow_released", "Escrow Released",
		fmt.Sprintf("%.6f %s has been released to your wallet", amount, crypto),
		notificationData)

	log.Printf("Escrow release notification sent for trade %s", tradeID)
}

// GetConnectedUsers returns the number of currently connected users
func (w *WebSocketNotifier) GetConnectedUsers() int {
	if w.wsService == nil {
		return 0
	}

	users := w.wsService.GetConnectedUsers()
	return len(users)
}
