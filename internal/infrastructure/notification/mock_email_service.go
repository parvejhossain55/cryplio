package notification

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MockEmailService provides a mock email delivery service
type MockEmailService struct {
	emails      map[string]*MockEmail
	mutex       sync.RWMutex
	enabled     bool
	fromEmail   string
	fromName    string
	deliveryRate float64 // 0.0 to 1.0, success rate
}

// MockEmail represents a mock email
type MockEmail struct {
	ID          string            `json:"id"`
	To          string            `json:"to"`
	From        string            `json:"from"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTMLBody    string            `json:"html_body"`
	Status      string            `json:"status"` // "pending", "sent", "delivered", "failed"
	CreatedAt   time.Time         `json:"created_at"`
	SentAt      *time.Time        `json:"sent_at,omitempty"`
	DeliveredAt *time.Time        `json:"delivered_at,omitempty"`
	FailedAt    *time.Time        `json:"failed_at,omitempty"`
	ErrorReason string            `json:"error_reason,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    string `json:"html"`
}

// NewMockEmailService creates a new mock email service
func NewMockEmailService(fromEmail, fromName string) *MockEmailService {
	return &MockEmailService{
		emails:      make(map[string]*MockEmail),
		enabled:     true,
		fromEmail:   fromEmail,
		fromName:    fromName,
		deliveryRate: 0.95, // 95% success rate
	}
}

// SetEnabled enables or disables email sending
func (m *MockEmailService) SetEnabled(enabled bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.enabled = enabled
}

// SetDeliveryRate sets the email delivery success rate (0.0 to 1.0)
func (m *MockEmailService) SetDeliveryRate(rate float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	m.deliveryRate = rate
}

// SendEmail sends an email
func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, textBody, htmlBody string, metadata map[string]string) (string, error) {
	if !m.enabled {
		return "", fmt.Errorf("email service is disabled")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	emailID := fmt.Sprintf("email_%d_%s", time.Now().UnixNano(), uuid.New().String()[:8])
	
	email := &MockEmail{
		ID:       emailID,
		To:       to,
		From:     fmt.Sprintf("%s <%s>", m.fromName, m.fromEmail),
		Subject:  subject,
		Body:     textBody,
		HTMLBody: htmlBody,
		Status:   "pending",
		CreatedAt: time.Now(),
		Metadata: metadata,
	}

	m.emails[emailID] = email

	// Simulate async email delivery
	go m.simulateEmailDelivery(emailID)

	log.Printf("Email queued for delivery: %s to %s", emailID, to)
	
	return emailID, nil
}

// SendTemplateEmail sends an email using a template
func (m *MockEmailService) SendTemplateEmail(ctx context.Context, to string, template *EmailTemplate, data map[string]interface{}) (string, error) {
	// Simple template substitution
	subject := m.replaceTemplateVariables(template.Subject, data)
	textBody := m.replaceTemplateVariables(template.Body, data)
	htmlBody := m.replaceTemplateVariables(template.HTML, data)

	return m.SendEmail(ctx, to, subject, textBody, htmlBody, map[string]string{"template": template.Name})
}

// GetEmailStatus returns the status of an email
func (m *MockEmailService) GetEmailStatus(ctx context.Context, emailID string) (*MockEmail, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	email, exists := m.emails[emailID]
	if !exists {
		return nil, fmt.Errorf("email not found: %s", emailID)
	}

	// Return a copy to avoid concurrent modification
	emailCopy := *email
	return &emailCopy, nil
}

// simulateEmailDelivery simulates the email delivery process
func (m *MockEmailService) simulateEmailDelivery(emailID string) {
	// Simulate sending delay (1-5 seconds)
	sendDelay := time.Duration(1+rand.Intn(4)) * time.Second
	time.Sleep(sendDelay)

	m.mutex.Lock()
	email, exists := m.emails[emailID]
	if !exists {
		m.mutex.Unlock()
		return
	}

	// Update to sent status
	email.Status = "sent"
	sentAt := time.Now()
	email.SentAt = &sentAt
	m.mutex.Unlock()

	log.Printf("Email sent: %s", emailID)

	// Simulate delivery delay (1-3 seconds)
	deliveryDelay := time.Duration(1+rand.Intn(2)) * time.Second
	time.Sleep(deliveryDelay)

	m.mutex.Lock()
	email, exists = m.emails[emailID]
	if !exists {
		m.mutex.Unlock()
		return
	}

	// Determine if delivery succeeds based on delivery rate
	if rand.Float64() < m.deliveryRate {
		// Success
		email.Status = "delivered"
		deliveredAt := time.Now()
		email.DeliveredAt = &deliveredAt
		log.Printf("Email delivered: %s", emailID)
	} else {
		// Failure
		email.Status = "failed"
		failedAt := time.Now()
		email.FailedAt = &failedAt
		reasons := []string{
			"Recipient mailbox full",
			"Invalid email address",
			"SMTP server timeout",
			"Recipient server not responding",
			"Email marked as spam",
		}
		email.ErrorReason = reasons[rand.Intn(len(reasons))]
		log.Printf("Email delivery failed: %s - %s", emailID, email.ErrorReason)
	}
	m.mutex.Unlock()
}

// replaceTemplateVariables replaces template variables with actual data
func (m *MockEmailService) replaceTemplateVariables(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		if strValue, ok := value.(string); ok {
			result = fmt.Sprintf("%s", result)
			// Simple string replacement
			for i := 0; i < len(result); i++ {
				if i+len(placeholder) <= len(result) && result[i:i+len(placeholder)] == placeholder {
					result = result[:i] + strValue + result[i+len(placeholder):]
					break
				}
			}
		}
	}
	return result
}

// GetAllEmails returns all emails (for testing)
func (m *MockEmailService) GetAllEmails() []*MockEmail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	emails := make([]*MockEmail, 0, len(m.emails))
	for _, email := range m.emails {
		emailCopy := *email
		emails = append(emails, &emailCopy)
	}

	return emails
}

// GetStats returns email service statistics
func (m *MockEmailService) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]interface{}{
		"total":      0,
		"pending":    0,
		"sent":       0,
		"delivered":  0,
		"failed":     0,
		"enabled":    m.enabled,
		"success_rate": m.deliveryRate,
	}

	for _, email := range m.emails {
		stats["total"] = stats["total"].(int) + 1
		stats[email.Status] = stats[email.Status].(int) + 1
	}

	if total := stats["total"].(int); total > 0 {
		delivered := stats["delivered"].(int)
		stats["actual_delivery_rate"] = float64(delivered) / float64(total)
	} else {
		stats["actual_delivery_rate"] = 0.0
	}

	return stats
}

// CleanupOldEmails removes old emails (for testing/maintenance)
func (m *MockEmailService) CleanupOldEmails(olderThan time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	cutoff := time.Now().Add(-olderThan)
	removed := 0

	for id, email := range m.emails {
		if email.CreatedAt.Before(cutoff) {
			delete(m.emails, id)
			removed++
		}
	}

	log.Printf("Cleaned up %d old emails", removed)
}

// GetEmailsByRecipient returns emails for a specific recipient
func (m *MockEmailService) GetEmailsByRecipient(to string) []*MockEmail {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var emails []*MockEmail
	for _, email := range m.emails {
		if email.To == to {
			emailCopy := *email
			emails = append(emails, &emailCopy)
		}
	}

	return emails
}

// ResendFailedEmail attempts to resend a failed email
func (m *MockEmailService) ResendFailedEmail(ctx context.Context, emailID string) error {
	m.mutex.Lock()
	email, exists := m.emails[emailID]
	if !exists {
		m.mutex.Unlock()
		return fmt.Errorf("email not found: %s", emailID)
	}

	if email.Status != "failed" {
		m.mutex.Unlock()
		return fmt.Errorf("email is not in failed status: %s", email.Status)
	}

	// Reset email status
	email.Status = "pending"
	email.SentAt = nil
	email.DeliveredAt = nil
	email.FailedAt = nil
	email.ErrorReason = ""
	m.mutex.Unlock()

	// Simulate redelivery
	go m.simulateEmailDelivery(emailID)

	log.Printf("Email resend initiated: %s", emailID)
	return nil
}

// Predefined email templates
var (
	WelcomeEmailTemplate = &EmailTemplate{
		Name:    "welcome",
		Subject: "Welcome to Cryplio - Your P2P Crypto Journey Begins!",
		Body:    `Hello {{.Username}},

Welcome to Cryplio! Your account has been successfully created.

You can now:
- Browse and create trade ads
- Chat with trade partners
- Manage your wallet
- Track your transactions

Get started now: {{.DashboardURL}}

Best regards,
The Cryplio Team`,
		HTML: `<h2>Welcome to Cryplio, {{.Username}}!</h2>
<p>Your account has been successfully created. You can now start trading cryptocurrencies securely.</p>
<p><a href="{{.DashboardURL}}">Go to Dashboard</a></p>
<p>Best regards,<br>The Cryplio Team</p>`,
	}

	TradeInitiatedTemplate = &EmailTemplate{
		Name:    "trade_initiated",
		Subject: "Trade Initiated - {{.TradeID}}",
		Body:    `Hello {{.Username}},

Your trade has been initiated successfully.

Trade Details:
- Trade ID: {{.TradeID}}
- Amount: {{.Amount}} {{.Crypto}}
- Price: {{.Price}} {{.Fiat}}
- Total: {{.Total}} {{.Fiat}}

Please complete the payment within the specified time.

Best regards,
The Cryplio Team`,
		HTML: `<h2>Trade Initiated</h2>
<p>Your trade <strong>{{.TradeID}}</strong> has been initiated successfully.</p>
<ul>
<li>Amount: {{.Amount}} {{.Crypto}}</li>
<li>Price: {{.Price}} {{.Fiat}}</li>
<li>Total: {{.Total}} {{.Fiat}}</li>
</ul>
<p>Please complete the payment within the specified time.</p>`,
	}

	PaymentReceivedTemplate = &EmailTemplate{
		Name:    "payment_received",
		Subject: "Payment Received - Trade {{.TradeID}}",
		Body:    `Hello {{.Username}},

Payment has been received for your trade.

Trade ID: {{.TradeID}}
Amount: {{.Amount}} {{.Crypto}}

The seller has been notified and will release the crypto assets shortly.

Best regards,
The Cryplio Team`,
		HTML: `<h2>Payment Received</h2>
<p>Payment has been received for trade <strong>{{.TradeID}}</strong>.</p>
<p>Amount: {{.Amount}} {{.Crypto}}</p>
<p>The seller has been notified and will release the crypto assets shortly.</p>`,
	}
)
