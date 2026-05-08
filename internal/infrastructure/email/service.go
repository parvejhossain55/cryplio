package email

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"cryplio/pkg/database"

	"github.com/google/uuid"
)

// EmailService handles email sending and template management
type EmailService struct {
	db         *database.DB
	fromEmail  string
	fromName   string
	workerSize int
}

// NewEmailService creates a new email service
func NewEmailService(db *database.DB, fromEmail, fromName string, workerSize int) *EmailService {
	return &EmailService{
		db:         db,
		fromEmail:  fromEmail,
		fromName:   fromName,
		workerSize: workerSize,
	}
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	TemplateID string   `json:"template_id"`
	Name       string   `json:"name"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	Variables  []string `json:"variables"`
	IsActive   bool     `json:"is_active"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// EmailQueueItem represents an email in the queue
type EmailQueueItem struct {
	EmailID       string                 `json:"email_id"`
	ToEmail       string                 `json:"to_email"`
	FromEmail     string                 `json:"from_email"`
	Subject       string                 `json:"subject"`
	Body          string                 `json:"body"`
	TemplateID    string                 `json:"template_id"`
	Variables     map[string]interface{} `json:"variables"`
	Status        string                 `json:"status"`
	Attempts      int                    `json:"attempts"`
	MaxAttempts   int                    `json:"max_attempts"`
	LastAttemptAt *time.Time             `json:"last_attempt_at"`
	SentAt        *time.Time             `json:"sent_at"`
	ErrorMessage  string                 `json:"error_message"`
	Priority      int                    `json:"priority"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
}

// SendEmailWithTemplate sends an email using a template
func (s *EmailService) SendEmailWithTemplate(ctx context.Context, toEmail, templateName string, variables map[string]interface{}, priority int) error {
	// Get template
	tmpl, err := s.GetTemplateByName(ctx, templateName)
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	if !tmpl.IsActive {
		return fmt.Errorf("template is not active")
	}

	// Validate required variables
	for _, requiredVar := range tmpl.Variables {
		if _, exists := variables[requiredVar]; !exists {
			return fmt.Errorf("missing required variable: %s", requiredVar)
		}
	}

	// Process template
	subject, err := s.processTemplate(tmpl.Subject, variables)
	if err != nil {
		return fmt.Errorf("process subject template: %w", err)
	}

	body, err := s.processTemplate(tmpl.Body, variables)
	if err != nil {
		return fmt.Errorf("process body template: %w", err)
	}

	// Add to queue
	return s.AddToQueue(ctx, toEmail, subject, body, tmpl.TemplateID, variables, priority)
}

// SendDirectEmail sends an email directly without template
func (s *EmailService) SendDirectEmail(ctx context.Context, toEmail, subject, body string, priority int) error {
	return s.AddToQueue(ctx, toEmail, subject, body, "", nil, priority)
}

// AddToQueue adds an email to the sending queue
func (s *EmailService) AddToQueue(ctx context.Context, toEmail, subject, body, templateID string, variables map[string]interface{}, priority int) error {
	if priority == 0 {
		priority = 5 // Default priority
	}

	emailID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	query := `
		INSERT INTO email_queue (
			email_id, to_email, from_email, subject, body, template_id, variables,
			status, attempts, max_attempts, priority, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := s.db.ExecContext(ctx, query,
		emailID, toEmail, s.fromEmail, subject, body, templateID, variables,
		"pending", 0, 3, priority, now, now,
	)

	if err != nil {
		return fmt.Errorf("add to queue: %w", err)
	}

	return nil
}

// GetTemplateByName retrieves an email template by name
func (s *EmailService) GetTemplateByName(ctx context.Context, name string) (*EmailTemplate, error) {
	query := `
		SELECT template_id, name, subject, body, variables, is_active, created_at, updated_at
		FROM email_templates
		WHERE name = $1
	`

	var tmpl EmailTemplate
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&tmpl.TemplateID, &tmpl.Name, &tmpl.Subject, &tmpl.Body,
		&tmpl.Variables, &tmpl.IsActive, &tmpl.CreatedAt, &tmpl.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("get template: %w", err)
	}

	return &tmpl, nil
}

// processTemplate processes a template string with variables
func (s *EmailService) processTemplate(tmplStr string, variables map[string]interface{}) (string, error) {
	// Create template
	tmpl, err := template.New("email").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	// Execute template
	var result strings.Builder
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return result.String(), nil
}

// ProcessQueue processes emails in the queue
func (s *EmailService) ProcessQueue(ctx context.Context) error {
	// Get pending emails ordered by priority and creation time
	query := `
		SELECT email_id, to_email, from_email, subject, body, template_id, variables,
			   attempts, max_attempts, last_attempt_at, sent_at, error_message, priority, created_at, updated_at
		FROM email_queue
		WHERE status = 'pending' AND attempts < max_attempts
		ORDER BY priority ASC, created_at ASC
		LIMIT 50
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("query email queue: %w", err)
	}
	defer rows.Close()

	var emails []EmailQueueItem
	for rows.Next() {
		var email EmailQueueItem
		err := rows.Scan(
			&email.EmailID, &email.ToEmail, &email.FromEmail, &email.Subject, &email.Body,
			&email.TemplateID, &email.Variables, &email.Attempts, &email.MaxAttempts,
			&email.LastAttemptAt, &email.SentAt, &email.ErrorMessage, &email.Priority,
			&email.CreatedAt, &email.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("scan email queue item: %w", err)
		}
		emails = append(emails, email)
	}

	// Process each email
	for _, email := range emails {
		err := s.sendEmail(ctx, &email)
		if err != nil {
			log.Printf("Failed to send email %s: %v", email.EmailID, err)
		}
	}

	return nil
}

// sendEmail sends an individual email
func (s *EmailService) sendEmail(ctx context.Context, email *EmailQueueItem) error {
	// Update attempt count
	now := time.Now().UTC()
	err := s.updateEmailAttempt(ctx, email.EmailID, now)
	if err != nil {
		return fmt.Errorf("update attempt: %w", err)
	}

	// TODO: Implement actual email sending
	// For now, we'll just log and mark as sent
	log.Printf("Sending email to %s: %s", email.ToEmail, email.Subject)

	// Mark as sent (placeholder implementation)
	return s.markEmailSent(ctx, email.EmailID, now)
}

// updateEmailAttempt updates the attempt count for an email
func (s *EmailService) updateEmailAttempt(ctx context.Context, emailID string, attemptTime time.Time) error {
	query := `
		UPDATE email_queue
		SET attempts = attempts + 1, last_attempt_at = $2, updated_at = $2
		WHERE email_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, emailID, attemptTime)
	if err != nil {
		return fmt.Errorf("update attempt: %w", err)
	}

	return nil
}

// markEmailSent marks an email as sent
func (s *EmailService) markEmailSent(ctx context.Context, emailID string, sentTime time.Time) error {
	query := `
		UPDATE email_queue
		SET status = 'sent', sent_at = $2, updated_at = $2
		WHERE email_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, emailID, sentTime)
	if err != nil {
		return fmt.Errorf("mark as sent: %w", err)
	}

	return nil
}

// markEmailFailed marks an email as failed
func (s *EmailService) markEmailFailed(ctx context.Context, emailID string, errorMsg string, failTime time.Time) error {
	query := `
		UPDATE email_queue
		SET status = 'failed', error_message = $2, updated_at = $3
		WHERE email_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, emailID, errorMsg, failTime)
	if err != nil {
		return fmt.Errorf("mark as failed: %w", err)
	}

	return nil
}

// GetQueueStats returns statistics about the email queue
func (s *EmailService) GetQueueStats(ctx context.Context) (map[string]int64, error) {
	query := `
		SELECT status, COUNT(*)
		FROM email_queue
		GROUP BY status
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query queue stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("scan queue stats: %w", err)
		}
		stats[status] = count
	}

	return stats, nil
}

// StartWorker starts the email processing worker
func (s *EmailService) StartWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			err := s.ProcessQueue(ctx)
			if err != nil {
				log.Printf("Failed to process email queue: %v", err)
			}
		}
	}
}

// SendTradeCreatedEmail sends an email when a trade is created
func (s *EmailService) SendTradeCreatedEmail(ctx context.Context, buyerEmail, sellerEmail, buyerName, sellerName, tradeID string, cryptoAmount, fiatAmount float64, cryptoSymbol, fiatSymbol, paymentMethod string, paymentWindow int) error {
	variables := map[string]interface{}{
		"user_name":            buyerName,
		"trade_id":             tradeID,
		"crypto_amount":        fmt.Sprintf("%.6f", cryptoAmount),
		"crypto_symbol":        cryptoSymbol,
		"fiat_amount":          fmt.Sprintf("%.2f", fiatAmount),
		"fiat_symbol":          fiatSymbol,
		"counterpart_username": sellerName,
		"payment_method":       paymentMethod,
		"payment_window":       paymentWindow,
	}

	// Send to buyer
	err := s.SendEmailWithTemplate(ctx, buyerEmail, "trade_created", variables, 5)
	if err != nil {
		return fmt.Errorf("send buyer email: %w", err)
	}

	// Send to seller with different variables
	variables["user_name"] = sellerName
	variables["counterpart_username"] = buyerName
	err = s.SendEmailWithTemplate(ctx, sellerEmail, "trade_created", variables, 5)
	if err != nil {
		return fmt.Errorf("send seller email: %w", err)
	}

	return nil
}

// SendTradeCompletedEmail sends an email when a trade is completed
func (s *EmailService) SendTradeCompletedEmail(ctx context.Context, buyerEmail, sellerEmail, buyerName, sellerName, tradeID string, cryptoAmount, cryptoSymbol string) error {
	variables := map[string]interface{}{
		"user_name":            buyerName,
		"trade_id":             tradeID,
		"crypto_amount":        fmt.Sprintf("%.6f", cryptoAmount),
		"crypto_symbol":        cryptoSymbol,
		"counterpart_username": sellerName,
	}

	// Send to buyer
	err := s.SendEmailWithTemplate(ctx, buyerEmail, "trade_completed", variables, 5)
	if err != nil {
		return fmt.Errorf("send buyer email: %w", err)
	}

	// Send to seller
	variables["user_name"] = sellerName
	variables["counterpart_username"] = buyerName
	err = s.SendEmailWithTemplate(ctx, sellerEmail, "trade_completed", variables, 5)
	if err != nil {
		return fmt.Errorf("send seller email: %w", err)
	}

	return nil
}

// SendDisputeCreatedEmail sends an email when a dispute is created
func (s *EmailService) SendDisputeCreatedEmail(ctx context.Context, userEmail, userName, tradeID, reason, description string) error {
	variables := map[string]interface{}{
		"user_name":           userName,
		"trade_id":            tradeID,
		"dispute_reason":      reason,
		"dispute_description": description,
	}

	return s.SendEmailWithTemplate(ctx, userEmail, "dispute_created", variables, 3)
}

// SendWithdrawalApprovedEmail sends an email when a withdrawal is approved
func (s *EmailService) SendWithdrawalApprovedEmail(ctx context.Context, userEmail, userName string, amount float64, toAddress, txHash string) error {
	variables := map[string]interface{}{
		"user_name":  userName,
		"amount":     fmt.Sprintf("%.6f", amount),
		"to_address": toAddress,
		"tx_hash":    txHash,
	}

	return s.SendEmailWithTemplate(ctx, userEmail, "withdrawal_approved", variables, 5)
}

// SendSecurityAlertEmail sends a security alert email
func (s *EmailService) SendSecurityAlertEmail(ctx context.Context, userEmail, userName, alertMessage string) error {
	variables := map[string]interface{}{
		"user_name":     userName,
		"alert_message": alertMessage,
	}

	return s.SendEmailWithTemplate(ctx, userEmail, "security_alert", variables, 1)
}
