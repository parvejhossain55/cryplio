package identity

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cryplio/pkg/apperrors"
	"cryplio/pkg/logger"

	"github.com/google/uuid"
)

// ─── Google OAuth ─────────────────────────────────────────────────────────────

// GoogleOAuthURL builds the Google authorization URL. Returns an empty string
// if OAuth is not configured.
func (s *authService) GoogleOAuthURL() string {
	if s.googleClientID == "" {
		return ""
	}
	return "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + url.QueryEscape(s.googleClientID) +
		"&redirect_uri=" + url.QueryEscape(s.oauthRedirectURL) +
		"&response_type=code" +
		"&scope=openid%20email%20profile" +
		"&access_type=offline" +
		"&prompt=consent"
}

// GoogleCallback exchanges the authorization code for Google tokens, finds or
// creates the local user, and completes the login flow.
func (s *authService) GoogleCallback(ctx context.Context, code string) (string, string, *User, error) {
	if s.googleClientID == "" || s.googleClientSecret == "" || s.oauthRedirectURL == "" {
		return "", "", nil, apperrors.Internal("OAuth not configured", nil)
	}

	// ── Exchange code for Google tokens ───────────────────────────────────────
	formData := url.Values{
		"code":          {code},
		"client_id":     {s.googleClientID},
		"client_secret": {s.googleClientSecret},
		"redirect_uri":  {s.oauthRedirectURL},
		"grant_type":    {"authorization_code"},
	}
	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to create token request", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}
	tokenResp, err := httpClient.Do(tokenReq)
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to exchange code for token", err)
	}
	defer tokenResp.Body.Close()

	var googleTokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&googleTokens); err != nil {
		return "", "", nil, apperrors.Internal("failed to parse token response", err)
	}

	// ── Fetch user info from Google ───────────────────────────────────────────
	userInfoReq, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+googleTokens.AccessToken,
		nil,
	)
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to create user info request", err)
	}
	userInfoResp, err := httpClient.Do(userInfoReq)
	if err != nil {
		return "", "", nil, apperrors.Internal("failed to get user info", err)
	}
	defer userInfoResp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&googleUser); err != nil {
		return "", "", nil, apperrors.Internal("failed to parse user info", err)
	}

	// ── Find or create user ───────────────────────────────────────────────────
	user, err := s.findOrCreateOAuthUser(ctx, googleUser.ID, googleUser.Email,
		googleUser.EmailVerified, googleUser.GivenName, googleUser.Picture,
		googleUser.Name, googleTokens.AccessToken, googleTokens.RefreshToken,
		time.Duration(googleTokens.ExpiresIn)*time.Second,
	)
	if err != nil {
		return "", "", nil, err
	}

	// ── Complete login ────────────────────────────────────────────────────────
	accessToken, refreshToken, err := s.completeLogin(ctx, user)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

// findOrCreateOAuthUser resolves the local User for a given Google account.
// It either updates an existing OAuth link or creates a brand-new user.
func (s *authService) findOrCreateOAuthUser(
	ctx context.Context,
	providerID, email string,
	emailVerified bool,
	givenName, picture, displayName string,
	accessToken, refreshToken string,
	tokenTTL time.Duration,
) (*User, error) {
	existingOAuth, err := s.userRepo.GetOAuthByProviderID(ctx, "google", providerID)
	if err != nil {
		return nil, apperrors.Internal("database error", err)
	}

	if existingOAuth != nil {
		// Known OAuth link — refresh the stored tokens.
		user, err := s.userRepo.GetByID(ctx, existingOAuth.UserID)
		if err != nil {
			return nil, apperrors.NotFound("user not found", err)
		}
		existingOAuth.AccessToken = &accessToken
		existingOAuth.RefreshToken = &refreshToken
		expiry := time.Now().Add(tokenTTL)
		existingOAuth.TokenExpiry = &expiry
		if err := s.userRepo.UpdateOAuth(ctx, existingOAuth); err != nil {
			return nil, apperrors.Internal("failed to update oauth tokens", err)
		}
		return user, nil
	}

	// No OAuth link yet — look up by email or create a new user.
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, apperrors.Internal("database error", err)
	}

	if user == nil {
		base := givenName
		if base == "" {
			base = "user"
		}
		username := s.generateUniqueUsername(ctx, base)
		user = NewUser(email, username, "") // no password for OAuth accounts
		user.EmailVerified = emailVerified
		user.AvatarURL = &picture
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, apperrors.Internal("failed to create user", err)
		}

		// Auto-create a default wallet for the new user
		if s.walletService != nil {
			if _, err := s.walletService.CreateDefaultWallet(ctx, user.UserID); err != nil {
				logger.Error("failed to auto-create wallet during OAuth registration", logger.Fields{
					"user_id": user.UserID,
					"error":   err.Error(),
				})
			}
		}
	}

	oauth := &UserOAuth{
		ID:               uuid.New(),
		UserID:           user.UserID,
		Provider:         "google",
		ProviderUserID:   providerID,
		ProviderEmail:    &email,
		ProviderUsername: &displayName,
		AccessToken:      &accessToken,
		RefreshToken:     &refreshToken,
	}
	if err := s.userRepo.CreateOAuth(ctx, oauth); err != nil {
		return nil, apperrors.Internal("failed to link OAuth account", err)
	}
	return user, nil
}
