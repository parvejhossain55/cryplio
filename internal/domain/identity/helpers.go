package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"unicode"

	"cryplio/pkg/apperrors"
)

// hashToken returns the SHA-256 hex digest of a plain-text token.
// Tokens are never stored in plain text — always store the hash.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// validatePasswordComplexity enforces the password policy:
//   - at least 8 characters
//   - at least one uppercase letter
//   - at least one digit
//   - at least one special character
func validatePasswordComplexity(password string) error {
	if len(password) < 8 {
		return apperrors.Validation("password must be at least 8 characters", nil)
	}

	var hasUpper, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasDigit || !hasSpecial {
		return apperrors.Validation(
			"password must include at least one uppercase letter, one number, and one special character",
			nil,
		)
	}
	return nil
}
