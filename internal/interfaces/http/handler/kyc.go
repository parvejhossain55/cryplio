package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"cryplio/internal/application/user"
	"cryplio/internal/domain/identity"
	kycdomain "cryplio/internal/domain/kyc"
	kycinfra "cryplio/internal/infrastructure/kyc"
	"cryplio/internal/interfaces/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// KYCHandler handles KYC related HTTP requests
type KYCHandler struct {
	submitUC      *user.SubmitKYCUseCase
	verifyUC      *user.VerifyKYCUseCase
	personaClient kycinfra.PersonaClient
	webhookSecret string
}

// NewKYCHandler creates a new KYC handler
func NewKYCHandler(
	submitUC *user.SubmitKYCUseCase,
	verifyUC *user.VerifyKYCUseCase,
	personaClient kycinfra.PersonaClient,
	webhookSecret string,
) *KYCHandler {
	return &KYCHandler{
		submitUC:      submitUC,
		verifyUC:      verifyUC,
		personaClient: personaClient,
		webhookSecret: webhookSecret,
	}
}

// SubmitKYCHandler handles KYC document submission via Persona
func (h *KYCHandler) SubmitKYCHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.KYCSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := user.SubmitKYCInput{
		UserID:           userID,
		Level:            identity.KYCLevel(req.Level),
		DocumentType:     req.DocumentType,
		DocumentFrontURL: req.DocumentFrontURL,
		DocumentBackURL:  req.DocumentBackURL,
		SelfieURL:        req.SelfieURL,
		Provider:         "persona",
	}

	record, err := h.submitUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.mapToResponse(record))
}

// VerifyKYCHandler handles admin manual KYC verification
func (h *KYCHandler) VerifyKYCHandler(c *gin.Context) {
	adminIDStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin id"})
		return
	}

	kycIDStr := c.Param("id")
	kycID, err := uuid.Parse(kycIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kyc id"})
		return
	}

	var req dto.KYCVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := user.VerifyKYCInput{
		KYCID:           kycID,
		AdminID:         adminID,
		Approved:        req.Approved,
		RejectionReason: req.RejectionReason,
	}

	record, err := h.verifyUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.mapToResponse(record))
}

// PersonaWebhookHandler processes incoming events from Persona.
// Persona sends events such as inquiry.completed, inquiry.failed, etc.
// The handler validates the HMAC-SHA256 signature before processing.
func (h *KYCHandler) PersonaWebhookHandler(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Verify Persona webhook signature if secret is configured
	if h.webhookSecret != "" {
		sig := c.GetHeader("Persona-Signature")
		if !verifyPersonaSignature(payload, sig, h.webhookSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook signature"})
			return
		}
	}

	// Parse the event envelope
	var event struct {
		Data struct {
			Attributes struct {
				Name    string `json:"name"`
				Payload struct {
					Data struct {
						ID         string `json:"id"`
						Attributes struct {
							Status      string `json:"status"`
							ReferenceID string `json:"reference-id"` // our userID
						} `json:"attributes"`
					} `json:"data"`
				} `json:"payload"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// TODO: Route to a WebhookUseCase that looks up the record by provider_reference
	// and calls verifyUC to approve/reject based on Persona's outcome status.
	// Persona status values: created, pending, completed, failed, expired, needs_review

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// verifyPersonaSignature validates the HMAC-SHA256 signature on a Persona webhook.
func verifyPersonaSignature(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (h *KYCHandler) mapToResponse(record *kycdomain.KYCRecord) dto.KYCRecordResponse {
	return dto.KYCRecordResponse{
		KYCID:            record.KYCID.String(),
		UserID:           record.UserID.String(),
		Level:            string(record.Level),
		Status:           string(record.Status),
		DocumentType:     record.DocumentType,
		DocumentFrontURL: record.DocumentFrontURL,
		DocumentBackURL:  record.DocumentBackURL,
		SelfieURL:        record.SelfieURL,
		Provider:         record.Provider,
		RejectionReason:  record.RejectionReason,
		AMLScreened:      record.AMLScreened,
		SubmittedAt:      record.SubmittedAt,
		ReviewedAt:       record.ReviewedAt,
	}
}
