package kycinfra

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PersonaClient handles inquiry creation and webhook verification with Persona.
type PersonaClient interface {
	CreateInquiry(ctx context.Context, userID, templateID string) (inquiryID string, err error)
	GetInquiry(ctx context.Context, inquiryID string) (status string, err error)
	VerifyWebhookSignature(payload []byte, signature string) error
}

type personaClient struct {
	apiKey        string
	webhookSecret string
	baseURL       string
	httpClient    *http.Client
}

// NewPersonaClient creates a new real Persona client.
func NewPersonaClient(apiKey, webhookSecret string) PersonaClient {
	return &personaClient{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
		baseURL:       "https://withpersona.com/api/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateInquiry creates a new inquiry in Persona.
func (c *personaClient) CreateInquiry(ctx context.Context, userID, templateID string) (string, error) {
	path := "/inquiries"

	bodyData := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "inquiry",
			"attributes": map[string]interface{}{
				"template-id":  templateID,
				"reference-id": userID,
			},
		},
	}
	bodyBytes, _ := json.Marshal(bodyData)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Persona-Version", "2023-01-05")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("persona error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

// GetInquiry retrieves the status of an inquiry.
func (c *personaClient) GetInquiry(ctx context.Context, inquiryID string) (string, error) {
	path := "/inquiries/" + inquiryID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Persona-Version", "2023-01-05")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Attributes struct {
				Status string `json:"status"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.Attributes.Status, nil
}

// VerifyWebhookSignature verifies the signature of a webhook request from Persona.
func (c *personaClient) VerifyWebhookSignature(payload []byte, signature string) error {
	if c.webhookSecret == "" {
		return nil // Skip if not configured
	}

	h := hmac.New(sha256.New, []byte(c.webhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return fmt.Errorf("invalid persona webhook signature")
	}
	return nil
}

// NoopPersonaClient is a placeholder for testing.
type NoopPersonaClient struct{}

func NewNoopPersonaClient() PersonaClient {
	return &NoopPersonaClient{}
}

func (c *NoopPersonaClient) CreateInquiry(ctx context.Context, userID, templateID string) (string, error) {
	return "inq_noop_" + userID, nil
}

func (c *NoopPersonaClient) GetInquiry(ctx context.Context, inquiryID string) (string, error) {
	return "completed", nil
}

func (c *NoopPersonaClient) VerifyWebhookSignature(payload []byte, signature string) error {
	return nil
}
