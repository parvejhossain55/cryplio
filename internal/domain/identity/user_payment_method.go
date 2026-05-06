package identity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserPaymentMethod struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	PaymentMethodCode string          `json:"payment_method_code"`
	DisplayName       string          `json:"display_name"`
	AccountName       string          `json:"account_name,omitempty"`
	AccountNumber     string          `json:"account_number,omitempty"`
	BankName          string          `json:"bank_name,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	IsActive          bool            `json:"is_active"`
	IsDefault         bool            `json:"is_default"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
