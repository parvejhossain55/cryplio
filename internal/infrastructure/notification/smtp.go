package notification

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// EmailMessage describes an outbound email notification.
type EmailMessage struct {
	To      string
	Subject string
	HTML    string
	Text    string
}

// EmailSender sends transactional email notifications.
type EmailSender interface {
	Send(ctx context.Context, message EmailMessage) error
}

// PasswordResetMailer sends password reset instructions to a user.
type PasswordResetMailer interface {
	SendPasswordReset(ctx context.Context, email, token string) error
}

// SMTPClient sends transactional email notifications through SMTP (e.g., Gmail).
type SMTPClient struct {
	host        string
	port        int
	username    string
	password    string
	from        string
	frontendURL string
}

// NewSMTPClient creates a new SMTP email client.
func NewSMTPClient(host string, port int, username, password, from, frontendURL string) *SMTPClient {
	return &SMTPClient{
		host:        strings.TrimSpace(host),
		port:        port,
		username:    strings.TrimSpace(username),
		password:    strings.TrimSpace(password),
		from:        strings.TrimSpace(from),
		frontendURL: strings.TrimRight(strings.TrimSpace(frontendURL), "/"),
	}
}

// Send sends an email via SMTP.
func (c *SMTPClient) Send(ctx context.Context, message EmailMessage) error {
	if c.host == "" {
		return fmt.Errorf("smtp host is not configured")
	}
	if c.username == "" {
		return fmt.Errorf("smtp username is not configured")
	}
	if c.password == "" {
		return fmt.Errorf("smtp password is not configured")
	}
	if c.from == "" {
		return fmt.Errorf("email sender address is not configured")
	}
	if strings.TrimSpace(message.To) == "" {
		return fmt.Errorf("recipient address is required")
	}

	// Build email message
	var body bytes.Buffer

	// Headers
	body.WriteString("From: " + c.from + "\r\n")
	body.WriteString("To: " + message.To + "\r\n")
	body.WriteString("Subject: " + message.Subject + "\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: multipart/alternative; boundary=\"boundary\"\r\n\r\n")

	// Plain text part
	body.WriteString("--boundary\r\n")
	body.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString(message.Text + "\r\n\r\n")

	// HTML part
	body.WriteString("--boundary\r\n")
	body.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString(message.HTML + "\r\n\r\n")

	// End boundary
	body.WriteString("--boundary--\r\n")

	// Authenticate
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	address := net.JoinHostPort(c.host, fmt.Sprintf("%d", c.port))
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	var client *smtp.Client
	if c.port == 465 {
		tlsConn := tls.Client(conn, &tls.Config{ServerName: c.host, MinVersion: tls.VersionTLS12})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("tls handshake failed: %w", err)
		}
		client, err = smtp.NewClient(tlsConn, c.host)
	} else {
		client, err = smtp.NewClient(conn, c.host)
	}
	if err != nil {
		return fmt.Errorf("failed to create smtp client: %w", err)
	}
	defer client.Close()

	if c.port == 587 {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return fmt.Errorf("smtp server does not support STARTTLS")
		}
		if err := client.StartTLS(&tls.Config{ServerName: c.host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(c.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(message.To); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	_, err = w.Write(body.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Quit
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit smtp client: %w", err)
	}

	return nil
}

func (c *SMTPClient) SendEmail(ctx context.Context, to, subject, body string) error {
	return c.Send(ctx, EmailMessage{
		To:      to,
		Subject: subject,
		Text:    body,
		HTML:    "<html>" + body + "</html>",
	})
}

// SendPasswordReset sends a password reset email (implements PasswordResetMailer interface).
func (c *SMTPClient) SendPasswordReset(ctx context.Context, email, token string) error {
	resetLink := c.passwordResetURL(token)
	return c.Send(ctx, EmailMessage{
		To:      email,
		Subject: "Reset your Cryplio password",
		Text:    smtpPasswordResetText(resetLink),
		HTML:    smtpPasswordResetHTML(resetLink),
	})
}

// SendVerificationEmail sends an email verification email.
func (c *SMTPClient) SendVerificationEmail(ctx context.Context, email, token string) error {
	verifyLink := c.emailVerifyURL(token)
	return c.Send(ctx, EmailMessage{
		To:      email,
		Subject: "Verify your Cryplio email",
		Text:    smtpVerificationEmailText(verifyLink),
		HTML:    smtpVerificationEmailHTML(verifyLink),
	})
}

func (c *SMTPClient) passwordResetURL(token string) string {
	baseURL := c.frontendURL
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return baseURL + "/reset-password?token=" + token
}

func (c *SMTPClient) emailVerifyURL(token string) string {
	baseURL := c.frontendURL
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	return baseURL + "/verify-email?token=" + token
}

func smtpPasswordResetText(resetLink string) string {
	return "We received a request to reset your Cryplio password.\n\n" +
		"Open this link to choose a new password:\n" + resetLink + "\n\n" +
		"This link expires in 15 minutes. If you did not request a password reset, you can ignore this email."
}

func smtpPasswordResetHTML(resetLink string) string {
	return `<!doctype html>
<html>
  <body style="margin:0;padding:0;background:#f5f7fb;font-family:Arial,sans-serif;color:#111827;">
    <div style="max-width:560px;margin:0 auto;padding:32px 20px;">
      <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:8px;padding:28px;">
        <h1 style="margin:0 0 16px;font-size:22px;line-height:1.3;color:#111827;">Reset your Cryplio password</h1>
        <p style="margin:0 0 20px;font-size:15px;line-height:1.6;color:#374151;">We received a request to reset your password. Use the button below to choose a new one.</p>
        <p style="margin:0 0 24px;">
          <a href="` + resetLink + `" style="display:inline-block;background:#111827;color:#ffffff;text-decoration:none;border-radius:6px;padding:12px 18px;font-weight:700;">Reset password</a>
        </p>
        <p style="margin:0 0 12px;font-size:13px;line-height:1.5;color:#6b7280;">This link expires in 15 minutes. If you did not request a password reset, you can ignore this email.</p>
        <p style="margin:0;font-size:12px;line-height:1.5;color:#9ca3af;word-break:break-all;">` + resetLink + `</p>
      </div>
    </div>
  </body>
</html>`
}

func smtpVerificationEmailText(verifyLink string) string {
	return "Thanks for signing up for Cryplio!\n\n" +
		"Please click the link below to verify your email address:\n" + verifyLink + "\n\n" +
		"This link expires in 24 hours. If you did not create an account, you can ignore this email."
}

func smtpVerificationEmailHTML(verifyLink string) string {
	return `<!doctype html>
<html>
  <body style="margin:0;padding:0;background:#f5f7fb;font-family:Arial,sans-serif;color:#111827;">
    <div style="max-width:560px;margin:0 auto;padding:32px 20px;">
      <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:8px;padding:28px;">
        <h1 style="margin:0 0 16px;font-size:22px;line-height:1.3;color:#111827;">Verify your Cryplio email</h1>
        <p style="margin:0 0 20px;font-size:15px;line-height:1.6;color:#374151;">Thanks for signing up for Cryplio! Please click the button below to verify your email address.</p>
        <p style="margin:0 0 24px;">
          <a href="` + verifyLink + `" style="display:inline-block;background:#111827;color:#ffffff;text-decoration:none;border-radius:6px;padding:12px 18px;font-weight:700;">Verify email</a>
        </p>
        <p style="margin:0 0 12px;font-size:13px;line-height:1.5;color:#6b7280;">This link expires in 24 hours. If you did not create an account, you can ignore this email.</p>
        <p style="margin:0;font-size:12px;line-height:1.5;color:#9ca3af;word-break:break-all;">` + verifyLink + `</p>
      </div>
    </div>
  </body>
</html>`
}
