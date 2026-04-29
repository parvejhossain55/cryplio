package dto

// RegisterRequest is the auth registration payload.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is the auth login payload.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateProfileRequest is the authenticated profile update payload.
type UpdateProfileRequest struct {
	Username *string `json:"username,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}
