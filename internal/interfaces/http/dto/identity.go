package dto

// UserResponse is the delivery-safe projection of a user.
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	KYCLevel int    `json:"kyc_level"`
}

// AuthResponse is returned from login and registration flows.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
