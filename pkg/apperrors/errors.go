package apperrors

import (
	"errors"
	"fmt"
)

// AppError represents an application error with code and message
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithMessage overrides the message
func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

// Common error codes
const (
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodePermissionDenied  = "PERMISSION_DENIED"
	ErrCodeRateLimited       = "RATE_LIMITED"
	ErrCodeDisabledAccount   = "DISABLED_ACCOUNT"
	ErrCodeSuspendedAccount  = "SUSPENDED_ACCOUNT"
	ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrCodeEscrowLocked      = "ESCROW_LOCKED"
	ErrCodeTradeNotFound     = "TRADE_NOT_FOUND"
	ErrCodeAdNotFound        = "AD_NOT_FOUND"
	ErrCodeWalletNotFound    = "WALLET_NOT_FOUND"
	ErrCodeDisputeExists     = "DISPUTE_EXISTS"
	ErrCodeInvalidState      = "INVALID_STATE"
)

// Common error constructors
func NotFound(message string, err error) *AppError {
	return New(ErrCodeNotFound, message, err)
}

func InvalidInput(message string, err error) *AppError {
	return New(ErrCodeInvalidInput, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(ErrCodeUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return New(ErrCodeForbidden, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(ErrCodeConflict, message, err)
}

func Internal(message string, err error) *AppError {
	return New(ErrCodeInternal, message, err)
}

func Validation(message string, err error) *AppError {
	return New(ErrCodeValidation, message, err)
}

func RateLimited(message string, err error) *AppError {
	return New(ErrCodeRateLimited, message, err)
}

func InsufficientFunds(message string, err error) *AppError {
	return New(ErrCodeInsufficientFunds, message, err)
}

func InvalidState(message string, err error) *AppError {
	return New(ErrCodeInvalidState, message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from error
func GetAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
