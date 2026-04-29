package identity

import (
	"context"
	"database/sql"
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
)

// UserRegistrar handles user registration.
type UserRegistrar interface {
	Register(ctx context.Context, email, username, password string) (*User, error)
}

// Authenticator handles login/logout.
type Authenticator interface {
	Login(ctx context.Context, email, password string) (string, *User, error)
	Logout(ctx context.Context) error
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
}

// AuthService combines all auth operations (legacy composite).
type AuthService interface {
	UserRegistrar
	Authenticator
	OAuthProvider
	ProfileManager
}

type authService struct {
	userRepo           UserRepository
	jwtSecret          string
	jwtExpiry          time.Duration
	cookieName         string
	cookieSecure       bool
	cookieSameSite     string
	googleClientID     string
	googleClientSecret string
	oauthRedirectURL   string
	eventDispatcher    events.Dispatcher
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo UserRepository, jwtSecret string, jwtExpiry time.Duration, cookieName string, cookieSecure bool, cookieSameSite string) *authService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiry:      jwtExpiry,
		cookieName:     cookieName,
		cookieSecure:   cookieSecure,
		cookieSameSite: cookieSameSite,
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

// Login logs in a user and returns JWT token
func (s *authService) Login(ctx context.Context, email, password string) (string, *User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, apperrors.Unauthorized("invalid credentials", nil)
		}
		return "", nil, apperrors.Internal("database error", err)
	}
	if user == nil {
		return "", nil, apperrors.Unauthorized("invalid credentials", nil)
	}

	// Check if account is locked
	if user.IsLocked() {
		return "", nil, apperrors.RateLimited("account temporarily locked", nil)
	}

	// Verify password
	if err := sharedcrypto.CheckPassword(user.PasswordHash, password); err != nil {
		// increment failed attempts (ignore error)
		s.userRepo.IncrementFailedAttempts(ctx, user.UserID)
		return "", nil, apperrors.Unauthorized("invalid credentials", nil)
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

// Logout is a no-op for stateless JWT; client should discard token
func (s *authService) Logout(ctx context.Context) error {
	// Could add token to blacklist if needed
	// Publish event
	if s.eventDispatcher != nil {
		// Cannot get userID from context easily without token; could pass userID from context if middleware sets it.
		// For now, skip or fetch from context if available.
		// We'll leave as no-op event publishing for logout due to stateless nature.
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

// generateJWT creates a JWT token for user
func (s *authService) generateJWT(user *User) string {
	token, err := sharedjwt.Issue(s.jwtSecret, s.jwtExpiry, sharedjwt.Claims{
		sharedjwt.ClaimUserID:    user.UserID.String(),
		"email":                  user.Email,
		"username":               user.Username,
		"kyc_level":              string(user.KYCLevel),
		sharedjwt.ClaimTokenType: "access",
	})
	if err != nil {
		return ""
	}
	return token
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
