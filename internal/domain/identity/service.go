package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cryplio/internal/domain/events"
	"cryplio/pkg/apperrors"
	sharedcrypto "cryplio/pkg/crypto"
	sharedjwt "cryplio/pkg/jwt"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

// Token type constants
const (
	TokenTypeAccess     = "access"
	TokenTypeRefresh    = "refresh"
	TokenType2FAPending = "2fa_pending"
)

// TwoFactorRequiredError indicates that login requires 2FA verification
type TwoFactorRequiredError struct {
	User      *User
	TempToken string
}

func (e TwoFactorRequiredError) Error() string { return "two factor authentication required" }

// UserRegistrar handles user registration.
type UserRegistrar interface {
	Register(ctx context.Context, email, username, password string) (*User, error)
}

// Authenticator handles login/logout and token management.
type Authenticator interface {
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, user *User, err error)
	Logout(ctx context.Context, tokenID string) error
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, refreshTokenNew string, user *User, err error)
	Complete2FALogin(ctx context.Context, tempToken, code string) (accessToken, refreshToken string, user *User, err error)
}

// OAuthProvider handles third-party OAuth flows.
type OAuthProvider interface {
	GoogleOAuthURL() string
	GoogleCallback(ctx context.Context, code string) (string, *User, error)
}

// ProfileManager handles user profile queries and updates.
type ProfileManager interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, username, bio *string) (*User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*User, error)
}

// EmailVerifier handles email verification flows.
type EmailVerifier interface {
	RequestEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) (*User, error)
}

// PasswordResetter handles password reset flows.
type PasswordResetter interface {
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) (*User, error)
}

// EmailMailer sends account emails to a user.
type EmailMailer interface {
	SendPasswordReset(ctx context.Context, email, token string) error
	SendVerificationEmail(ctx context.Context, email, token string) error
}

// TwoFactorManager handles 2FA setup and verification.
type TwoFactorManager interface {
	Setup2FA(ctx context.Context, userID uuid.UUID) (secret, provisioningURI string, err error)
	Verify2FA(ctx context.Context, userID uuid.UUID, code string) error
	Disable2FA(ctx context.Context, userID uuid.UUID, password string) error
	Is2FAEnabled(user *User) bool
}

// SessionManager handles user session management.
type SessionManager interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tokenID string, deviceFingerprint, ipAddress, userAgent, deviceType, location *string, isRemembered bool, expiresAt time.Time) (*UserSession, error)
	GetSession(ctx context.Context, tokenID string) (*UserSession, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error)
	DeleteSession(ctx context.Context, tokenID string) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
}

// AuthService combines all auth operations (legacy composite).
type AuthService interface {
	UserRegistrar
	Authenticator
	OAuthProvider
	ProfileManager
	EmailVerifier
	PasswordResetter
	TwoFactorManager
	SessionManager
}

type authService struct {
	userRepo           UserRepository
	jwtSecret          string
	jwtExpiry          time.Duration
	refreshTokenExpiry time.Duration
	cookieName         string
	cookieSecure       bool
	cookieSameSite     string
	googleClientID     string
	googleClientSecret string
	oauthRedirectURL   string
	eventDispatcher    events.Dispatcher
	emailMailer        EmailMailer
	issuerName         string // for TOTP
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo UserRepository,
	jwtSecret string,
	jwtExpiry time.Duration,
	refreshTokenExpiry time.Duration,
	cookieName string,
	cookieSecure bool,
	cookieSameSite string,
	issuerName string,
) *authService {
	return &authService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		jwtExpiry:          jwtExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		cookieName:         cookieName,
		cookieSecure:       cookieSecure,
		cookieSameSite:     cookieSameSite,
		issuerName:         issuerName,
	}
}

// WithGoogleOAuth sets Google OAuth credentials
func (s *authService) WithGoogleOAuth(clientID, clientSecret, redirectURL string) *authService {
	s.googleClientID = clientID
	s.googleClientSecret = clientSecret
	s.oauthRedirectURL = redirectURL
	return s
}

// WithEventDispatcher sets the event dispatcher.
func (s *authService) WithEventDispatcher(dispatcher events.Dispatcher) *authService {
	s.eventDispatcher = dispatcher
	return s
}

// WithPasswordResetMailer sets the password reset mailer.
func (s *authService) WithPasswordResetMailer(mailer EmailMailer) *authService {
	s.emailMailer = mailer
	return s
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, email, username, password string) (*User, error) {
	// Check if user already exists by email
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, apperrors.Internal("database error", err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("user already exists", nil)
	}

	// Hash password
	hashedPassword, err := sharedcrypto.HashPassword(password)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	// Create user
	user := NewUser(email, username, hashedPassword)
	if err := s.userRepo.Create(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, apperrors.Conflict("user already exists", nil)
		}
		return nil, apperrors.Internal("failed to create user", err)
	}

	// Publish event
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.UserRegisteredEvent{
			UserID:    user.UserID.String(),
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		})
	}

	return user, nil
}

// Login logs in a user and returns access and refresh tokens.
// If 2FA is enabled, returns a TwoFactorRequiredError.
func (s *authService) Login(ctx context.Context, email, password string) (string, string, *User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
		}
		return "", "", nil, apperrors.Internal("database error", err)
	}
	if user == nil {
		return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
	}

	// Check if account is locked
	if user.IsLocked() {
		return "", "", nil, apperrors.RateLimited("account temporarily locked", nil)
	}

	// Verify password
	if err := sharedcrypto.CheckPassword(user.PasswordHash, password); err != nil {
		// increment failed attempts (ignore error)
		s.userRepo.IncrementFailedAttempts(ctx, user.UserID)
		return "", "", nil, apperrors.Unauthorized("invalid credentials", nil)
	}

	// Check for 2FA
	if s.Is2FAEnabled(user) {
		tempToken := s.generate2FATempToken(user)
		// Return error to indicate 2FA required
		return "", "", user, TwoFactorRequiredError{User: user, TempToken: tempToken}
	}

	// Complete login without 2FA
	access, refresh, err := s.completeLogin(ctx, user)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

// Logout logs out a user by invalidating the session token
func (s *authService) Logout(ctx context.Context, tokenID string) error {
	// Delete session
	if tokenID != "" {
		_ = s.userRepo.DeleteSession(ctx, tokenID)
	}
	// Publish event
	if s.eventDispatcher != nil {
		if userID, ok := ctx.Value("user_id").(string); ok {
			_ = s.eventDispatcher.Dispatch(ctx, events.UserLoggedOutEvent{
				UserID:      userID,
				LoggedOutAt: time.Now(),
			})
		}
	}
	return nil
}

// GetUserByID fetches a user
func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}
	return user, nil
}

// UpdateProfile updates user profile fields
func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, username, bio *string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}

	if username != nil && *username != "" {
		user.Username = *username
	}
	if bio != nil {
		user.Bio = bio
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	return user, nil
}

// UpdateAvatar updates user's avatar URL
func (s *authService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return nil, apperrors.NotFound("user not found", nil)
	}

	user.AvatarURL = &avatarURL
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	return user, nil
}

// GoogleOAuthURL returns Google OAuth URL
func (s *authService) GoogleOAuthURL() string {
	if s.googleClientID == "" {
		return ""
	}
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + url.QueryEscape(s.googleClientID) +
		"&redirect_uri=" + url.QueryEscape(s.oauthRedirectURL) +
		"&response_type=code" +
		"&scope=openid%20email%20profile" +
		"&access_type=offline" +
		"&prompt=consent"
	return authURL
}

// GoogleCallback handles Google OAuth callback, returns JWT token and user
func (s *authService) GoogleCallback(ctx context.Context, code string) (string, *User, error) {
	if s.googleClientID == "" || s.googleClientSecret == "" || s.oauthRedirectURL == "" {
		return "", nil, apperrors.Internal("OAuth not configured", nil)
	}

	// Exchange code for tokens
	tokenURL := "https://oauth2.googleapis.com/token"
	formData := url.Values{
		"code":          {code},
		"client_id":     {s.googleClientID},
		"client_secret": {s.googleClientSecret},
		"redirect_uri":  {s.oauthRedirectURL},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", nil, apperrors.Internal("failed to create token request", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, apperrors.Internal("failed to exchange code for token", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", nil, apperrors.Internal("failed to parse token response", err)
	}

	// Get user info from Google
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokenResp.AccessToken)
	if err != nil {
		return "", nil, apperrors.Internal("failed to get user info", err)
	}
	defer userInfoResp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&googleUser); err != nil {
		return "", nil, apperrors.Internal("failed to parse user info", err)
	}

	// Check if OAuth already linked
	existingOAuth, err := s.userRepo.GetOAuthByProviderID(ctx, "google", googleUser.ID)
	if err != nil {
		return "", nil, apperrors.Internal("database error", err)
	}

	var user *User
	if existingOAuth != nil {
		user, err = s.userRepo.GetByID(ctx, existingOAuth.UserID)
		if err != nil {
			return "", nil, apperrors.NotFound("user not found", err)
		}
		existingOAuth.AccessToken = &tokenResp.AccessToken
		existingOAuth.RefreshToken = &tokenResp.RefreshToken
		expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		existingOAuth.TokenExpiry = &expiry
		if err := s.userRepo.UpdateOAuth(ctx, existingOAuth); err != nil {
			return "", nil, apperrors.Internal("failed to update oauth", err)
		}
	} else {
		user, err = s.userRepo.GetByEmail(ctx, googleUser.Email)
		if err != nil && err != sql.ErrNoRows {
			return "", nil, apperrors.Internal("database error", err)
		}

		if user == nil {
			// Create new user
			baseUsername := googleUser.GivenName
			if baseUsername == "" {
				baseUsername = "user"
			}
			uniqueUsername := s.generateUniqueUsername(ctx, baseUsername)
			user = NewUser(googleUser.Email, uniqueUsername, "") // password empty for OAuth
			user.EmailVerified = googleUser.EmailVerified
			user.AvatarURL = &googleUser.Picture
			if err := s.userRepo.Create(ctx, user); err != nil {
				return "", nil, apperrors.Internal("failed to create user", err)
			}
		}

		// Link OAuth
		oauth := &UserOAuth{
			ID:               uuid.New(),
			UserID:           user.UserID,
			Provider:         "google",
			ProviderUserID:   googleUser.ID,
			ProviderEmail:    &googleUser.Email,
			ProviderUsername: &googleUser.Name,
			AccessToken:      &tokenResp.AccessToken,
			RefreshToken:     &tokenResp.RefreshToken,
			TokenExpiry:      nil,
		}
		if err := s.userRepo.CreateOAuth(ctx, oauth); err != nil {
			return "", nil, apperrors.Internal("failed to link OAuth", err)
		}
	}

	// Update login stats
	_ = s.userRepo.IncrementLogin(ctx, user.UserID)

	// Generate JWT
	token := s.generateJWT(user)

	// Publish event
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.UserLoggedInEvent{
			UserID:     user.UserID.String(),
			Email:      user.Email,
			LoggedInAt: time.Now(),
		})
	}

	return token, user, nil
}

// generateJWT creates an access JWT token for user
func (s *authService) generateJWT(user *User) string {
	token, err := sharedjwt.Issue(s.jwtSecret, s.jwtExpiry, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		"email":                  user.Email,
		"username":               user.Username,
		"kyc_level":              string(user.KYCLevel),
		sharedjwt.ClaimTokenType: TokenTypeAccess,
	})
	if err != nil {
		return ""
	}
	return token
}

// generate2FATempToken creates a short-lived JWT for 2FA verification step
func (s *authService) generate2FATempToken(user *User) string {
	token, err := sharedjwt.Issue(s.jwtSecret, 5*time.Minute, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		sharedjwt.ClaimTokenType: TokenType2FAPending,
	})
	if err != nil {
		return ""
	}
	return token
}

// ============================================================
// Helper Functions
// ============================================================

// generateRandomToken creates a cryptographically secure random token
func (s *authService) generateRandomToken() (string, error) {
	b := make([]byte, 16) // 128-bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken creates a SHA256 hash of a token for secure storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Is2FAEnabled checks if user has 2FA enabled
func (s *authService) Is2FAEnabled(user *User) bool {
	return user.TwoFASecret != nil && *user.TwoFASecret != ""
}

// generateUniqueUsername creates a unique username
func (s *authService) generateUniqueUsername(ctx context.Context, base string) string {
	username := strings.ToLower(strings.ReplaceAll(base, " ", "."))
	for i := 1; i < 100; i++ {
		candidate := username
		if i > 1 {
			candidate = username + fmt.Sprintf("%d", i)
		}
		existing, _ := s.userRepo.GetByUsername(ctx, candidate)
		if existing == nil {
			return candidate
		}
	}
	// Fallback with random suffix
	return username + fmt.Sprintf(".%d", time.Now().UnixNano()%1000)
}

// completeLogin finalizes the login process after successful authentication
func (s *authService) completeLogin(ctx context.Context, user *User) (string, string, error) {
	// Increment login stats
	_ = s.userRepo.IncrementLogin(ctx, user.UserID)

	// Create session
	sessionID := uuid.New()
	session := &UserSession{
		ID:           sessionID,
		UserID:       user.UserID,
		TokenID:      sessionID.String(),
		IsRemembered: false,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return "", "", apperrors.Internal("failed to create session", err)
	}

	// Generate tokens
	accessToken := s.generateJWT(user)
	refreshToken := s.generateRefreshJWT(user, session.ID.String())

	// Publish event
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.UserLoggedInEvent{
			UserID:     user.UserID.String(),
			Email:      user.Email,
			LoggedInAt: time.Now(),
		})
	}

	return accessToken, refreshToken, nil
}

// generateRefreshJWT creates a refresh token JWT
func (s *authService) generateRefreshJWT(user *User, sessionID string) string {
	token, err := sharedjwt.Issue(s.jwtSecret, s.refreshTokenExpiry, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		sharedjwt.ClaimTokenType: TokenTypeRefresh,
		"jti":                    sessionID,
	})
	if err != nil {
		return ""
	}
	return token
}

// ============================================================
// Email Verification
// ============================================================

// RequestEmailVerification requests an email verification token
func (s *authService) RequestEmailVerification(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	if user.EmailVerified {
		return apperrors.Conflict("email already verified", nil)
	}

	plainToken, err := s.generateRandomToken()
	if err != nil {
		return apperrors.Internal("failed to generate token", err)
	}
	tokenHash := hashToken(plainToken)

	token := &EmailVerificationToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.userRepo.CreateEmailVerificationToken(ctx, token); err != nil {
		return apperrors.Internal("failed to create verification token", err)
	}

	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.EmailVerificationRequestedEvent{
			UserID:  userID.String(),
			Email:   user.Email,
			Token:   plainToken,
			Expires: token.ExpiresAt,
		})
	}
	if s.emailMailer != nil {
		if err := s.emailMailer.SendVerificationEmail(ctx, user.Email, plainToken); err != nil {
			return apperrors.Internal("failed to send verification email", err)
		}
	}

	return nil
}

// VerifyEmail verifies an email using the token
func (s *authService) VerifyEmail(ctx context.Context, tokenString string) (*User, error) {
	tokenHash := hashToken(tokenString)
	token, err := s.userRepo.GetEmailVerificationTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.Internal("database error", err)
	}
	if token == nil {
		return nil, apperrors.Unauthorized("invalid or expired token", nil)
	}
	if token.IsExpired() {
		return nil, apperrors.Unauthorized("token expired", nil)
	}
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	if token.IsVerified() || user.EmailVerified {
		return user, nil
	}
	if err := s.userRepo.MarkEmailVerificationTokenVerified(ctx, token.ID); err != nil {
		return nil, apperrors.Internal("failed to mark token verified", err)
	}
	user.EmailVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.EmailVerifiedEvent{
			UserID:     user.UserID.String(),
			Email:      user.Email,
			VerifiedAt: time.Now(),
		})
	}
	return user, nil
}

// ============================================================
// Password Reset
// ============================================================

// RequestPasswordReset requests a password reset token
func (s *authService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // don't reveal existence
		}
		return apperrors.Internal("database error", err)
	}
	plainToken, err := s.generateRandomToken()
	if err != nil {
		return apperrors.Internal("failed to generate token", err)
	}
	tokenHash := hashToken(plainToken)

	token := &PasswordResetToken{
		UserID:    user.UserID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.userRepo.CreatePasswordResetToken(ctx, token); err != nil {
		return apperrors.Internal("failed to create reset token", err)
	}
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.PasswordResetRequestedEvent{
			UserID: user.UserID.String(),
			Email:  user.Email,
			Token:  plainToken,
		})
	}
	if s.emailMailer != nil {
		if err := s.emailMailer.SendPasswordReset(ctx, user.Email, plainToken); err != nil {
			return apperrors.Internal("failed to send password reset email", err)
		}
	}
	return nil
}

// ResetPassword resets user password using token
func (s *authService) ResetPassword(ctx context.Context, tokenString, newPassword string) (*User, error) {
	tokenHash := hashToken(tokenString)
	token, err := s.userRepo.GetPasswordResetToken(ctx, tokenHash)
	if err != nil {
		return nil, apperrors.Internal("database error", err)
	}
	if token == nil {
		return nil, apperrors.Unauthorized("invalid or expired token", nil)
	}
	if token.IsExpired() {
		return nil, apperrors.Unauthorized("token expired", nil)
	}
	if token.IsUsed() {
		return nil, apperrors.Conflict("token already used", nil)
	}
	hashedPassword, err := sharedcrypto.HashPassword(newPassword)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, apperrors.NotFound("user not found", err)
	}
	user.PasswordHash = hashedPassword
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	if err := s.userRepo.MarkPasswordResetTokenUsed(ctx, token.ID); err != nil {
		return nil, apperrors.Internal("failed to mark token used", err)
	}
	// Invalidate all sessions (security)
	_ = s.userRepo.DeleteSessionsByUserID(ctx, user.UserID)
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.PasswordResetCompletedEvent{
			UserID:  user.UserID.String(),
			Email:   user.Email,
			ResetAt: time.Now(),
		})
	}
	return user, nil
}

// ============================================================
// Token Refresh
// ============================================================

// RefreshToken exchanges a refresh token for a new access token and refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, *User, error) {
	claims, err := sharedjwt.Parse(s.jwtSecret, refreshTokenString)
	if err != nil {
		return "", "", nil, apperrors.Unauthorized("invalid refresh token", err)
	}
	tokenType, ok := claims[sharedjwt.ClaimTokenType].(string)
	if !ok || tokenType != TokenTypeRefresh {
		return "", "", nil, apperrors.Unauthorized("invalid token type", nil)
	}
	userIDStr, ok := claims[sharedjwt.ClaimUserID].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token claims", nil)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", nil, apperrors.InvalidInput("invalid user id", err)
	}
	jti, ok := claims["jti"].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token", nil)
	}
	session, err := s.userRepo.GetSession(ctx, jti)
	if err != nil {
		return "", "", nil, apperrors.Internal("database error", err)
	}
	if session == nil {
		return "", "", nil, apperrors.Unauthorized("invalid session", nil)
	}
	if session.IsExpired() {
		_ = s.userRepo.DeleteSession(ctx, jti)
		return "", "", nil, apperrors.Unauthorized("session expired", nil)
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, apperrors.NotFound("user not found", err)
	}
	if user == nil || user.IsDeleted() {
		return "", "", nil, apperrors.NotFound("user not found", nil)
	}
	// Rotate refresh token
	_ = s.userRepo.DeleteSession(ctx, jti)
	newSessionID := uuid.New()
	newSession := &UserSession{
		ID:           newSessionID,
		UserID:       user.UserID,
		TokenID:      newSessionID.String(),
		IsRemembered: session.IsRemembered,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}
	if err := s.userRepo.CreateSession(ctx, newSession); err != nil {
		return "", "", nil, apperrors.Internal("failed to create session", err)
	}
	accessToken := s.generateJWT(user)
	refreshToken := s.generateRefreshJWT(user, newSessionID.String())
	return accessToken, refreshToken, user, nil
}

// ============================================================
// Two-Factor Authentication
// ============================================================

func (s *authService) Setup2FA(ctx context.Context, userID uuid.UUID) (secret, provisioningURI string, err error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", apperrors.NotFound("user not found", err)
	}
	if s.Is2FAEnabled(user) {
		return "", "", apperrors.Conflict("2FA already enabled", nil)
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuerName,
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", apperrors.Internal("failed to generate 2FA secret", err)
	}
	secret = key.Secret()
	provisioningURI = key.URL()

	// Store pending secret in database (expires in 10 minutes)
	pending := &TwoFactorPending{
		ID:        uuid.New(),
		UserID:    userID,
		Secret:    secret,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.userRepo.CreateTwoFactorPending(ctx, pending); err != nil {
		return "", "", apperrors.Internal("failed to store pending 2FA", err)
	}

	return secret, provisioningURI, nil
}

func (s *authService) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	// Retrieve pending record from database
	pending, err := s.userRepo.GetTwoFactorPendingByUserID(ctx, userID)
	if err != nil {
		return apperrors.Internal("database error", err)
	}
	if pending == nil {
		return apperrors.Unauthorized("2FA setup not initiated or expired", nil)
	}
	if pending.ExpiresAt.Before(time.Now()) {
		// Clean up expired pending
		_ = s.userRepo.DeleteTwoFactorPending(ctx, userID)
		return apperrors.Unauthorized("2FA setup expired, please start again", nil)
	}
	if !totp.Validate(code, pending.Secret) {
		return apperrors.Unauthorized("invalid 2FA code", nil)
	}
	// Enable 2FA on user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	user.TwoFASecret = &pending.Secret
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to update user", err)
	}
	// Remove pending record
	if err := s.userRepo.DeleteTwoFactorPending(ctx, userID); err != nil {
		// log but not critical
	}
	if s.eventDispatcher != nil {
		_ = s.eventDispatcher.Dispatch(ctx, events.TwoFactorEnabledEvent{
			UserID:    user.UserID.String(),
			EnabledAt: time.Now(),
		})
	}
	return nil
}

func (s *authService) Disable2FA(ctx context.Context, userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NotFound("user not found", err)
	}
	if err := sharedcrypto.CheckPassword(user.PasswordHash, password); err != nil {
		return apperrors.Unauthorized("invalid password", nil)
	}
	user.TwoFASecret = nil
	// Clear any pending
	_ = s.userRepo.DeleteTwoFactorPending(ctx, userID)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return apperrors.Internal("failed to update user", err)
	}
	return nil
}

// ============================================================
// Session Management (SessionManager)
// ============================================================

func (s *authService) CreateSession(ctx context.Context, userID uuid.UUID, tokenID string, deviceFingerprint, ipAddress, userAgent, deviceType, location *string, isRemembered bool, expiresAt time.Time) (*UserSession, error) {
	session := &UserSession{
		ID:                uuid.New(),
		UserID:            userID,
		TokenID:           tokenID,
		DeviceFingerprint: deviceFingerprint,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		DeviceType:        deviceType,
		Location:          location,
		IsRemembered:      isRemembered,
		ExpiresAt:         expiresAt,
		LastUsedAt:        time.Now(),
		CreatedAt:         time.Now(),
	}
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *authService) GetSession(ctx context.Context, tokenID string) (*UserSession, error) {
	return s.userRepo.GetSession(ctx, tokenID)
}

func (s *authService) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]UserSession, error) {
	return s.userRepo.GetSessionsByUserID(ctx, userID)
}

func (s *authService) DeleteSession(ctx context.Context, tokenID string) error {
	return s.userRepo.DeleteSession(ctx, tokenID)
}

func (s *authService) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteSessionsByUserID(ctx, userID)
}

// ============================================================
// 2FA Login Completion
// ============================================================

// Complete2FALogin completes login after successful 2FA verification
func (s *authService) Complete2FALogin(ctx context.Context, tempToken, code string) (string, string, *User, error) {
	claims, err := sharedjwt.Parse(s.jwtSecret, tempToken)
	if err != nil {
		return "", "", nil, apperrors.Unauthorized("invalid verification token", err)
	}
	tokenType, ok := claims[sharedjwt.ClaimTokenType].(string)
	if !ok || tokenType != TokenType2FAPending {
		return "", "", nil, apperrors.Unauthorized("invalid token type", nil)
	}
	userIDStr, ok := claims[sharedjwt.ClaimUserID].(string)
	if !ok {
		return "", "", nil, apperrors.Unauthorized("invalid token claims", nil)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", nil, apperrors.InvalidInput("invalid user id", err)
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", nil, apperrors.NotFound("user not found", err)
	}
	if !s.Is2FAEnabled(user) {
		return "", "", nil, apperrors.Unauthorized("2FA not enabled", nil)
	}
	if code == "" {
		return "", "", nil, apperrors.Validation("verification code required", nil)
	}
	if user.TwoFASecret == nil || *user.TwoFASecret == "" {
		return "", "", nil, apperrors.Unauthorized("2FA secret not configured", nil)
	}
	if !totp.Validate(code, *user.TwoFASecret) {
		return "", "", nil, apperrors.Unauthorized("invalid verification code", nil)
	}
	access, refresh, err := s.completeLogin(ctx, user)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}
