package dto

type AdResponse struct {
	AdID                 string   `json:"ad_id"`
	UserID               string   `json:"user_id"`
	Username             string   `json:"username"`
	UserAvatar           string   `json:"user_avatar,omitempty"`
	UserRating           float64  `json:"user_rating"`
	UserTrades           int      `json:"user_trades"`
	Type                 string   `json:"type"`
	CryptoSymbol         string   `json:"crypto_symbol"`
	FiatSymbol           string   `json:"fiat_symbol"`
	PriceType            string   `json:"price_type"`
	Price                float64  `json:"price"`
	MinAmount            float64  `json:"min_amount"`
	MaxAmount            float64  `json:"max_amount"`
	PaymentMethods       []string `json:"payment_methods"`
	PaymentWindowMinutes int      `json:"payment_window_minutes"`
	IsOnline             bool     `json:"is_online"`
}

type ListAdsResponse struct {
	Ads   []AdResponse `json:"ads"`
	Total int          `json:"total"`
}

type CreateAdRequest struct {
	Type                 string   `json:"type" binding:"required,oneof=buy sell"`
	CryptoID             int      `json:"crypto_id" binding:"required"`
	FiatID               int      `json:"fiat_id" binding:"required"`
	PriceType            string   `json:"price_type" binding:"required,oneof=fixed floating"`
	Price                float64  `json:"price" binding:"required,gt=0"`
	FloatingMarkup       *float64 `json:"floating_markup"`
	MinAmount            float64  `json:"min_amount" binding:"required,gt=0"`
	MaxAmount            float64  `json:"max_amount" binding:"required,gt=0"`
	PaymentMethods       []string `json:"payment_methods" binding:"required,min=1"`
	TradeTerms           string   `json:"trade_terms"`
	PaymentWindowMinutes int      `json:"payment_window_minutes" binding:"required,min=5,max=60"`
	RequiresKYCLevel     int      `json:"requires_kyc_level"`
}

type InitiateTradeRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
