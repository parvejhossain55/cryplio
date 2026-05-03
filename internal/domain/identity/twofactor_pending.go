package identity

import (
	"time"

	"github.com/google/uuid"
)

// TwoFactorPending holds a pending 2FA secret awaiting verification.
type TwoFactorPending struct {
	ID        uuid.UUID `db:"id" json:"-"`
	UserID    uuid.UUID `db:"user_id" json:"-"`
	Secret    string    `db:"secret" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"-"`
}
