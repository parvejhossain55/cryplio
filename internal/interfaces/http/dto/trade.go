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
	PaymentMethodIDs     []int    `json:"payment_method_ids"`
	PaymentWindowMinutes int      `json:"payment_window_minutes"`
	IsOnline             bool     `json:"is_online"`
	TradeTerms           string   `json:"trade_terms,omitempty"`
	Status               string   `json:"status"`
	CreatedAt            string   `json:"created_at"`
}

type ListAdsResponse struct {
	Ads   []AdResponse `json:"ads"`
	Total int          `json:"total"`
}

type CreateAdRequest struct {
	Type                 string  `json:"type" binding:"required,oneof=buy sell"`
	CryptoID             int     `json:"crypto_id" binding:"required"`
	FiatID               int     `json:"fiat_id" binding:"required"`
	PriceType            string  `json:"price_type" binding:"required,oneof=fixed floating"`
	Price                float64 `json:"price" binding:"required,gt=0"`
	MinAmount            float64 `json:"min_amount" binding:"required,gt=0"`
	MaxAmount            float64 `json:"max_amount" binding:"required,gt=0"`
	PaymentMethodIDs     []int   `json:"payment_method_ids" binding:"required,min=1"`
	TradeTerms           string  `json:"trade_terms"`
	PaymentWindowMinutes int     `json:"payment_window_minutes" binding:"required"`
}

type UpdateAdRequest struct {
	Type                 string  `json:"type" binding:"omitempty,oneof=buy sell"`
	CryptoID             int     `json:"crypto_id"`
	FiatID               int     `json:"fiat_id"`
	PriceType            string  `json:"price_type" binding:"omitempty,oneof=fixed floating"`
	Price                float64 `json:"price" binding:"omitempty,gt=0"`
	MinAmount            float64 `json:"min_amount" binding:"omitempty,gt=0"`
	MaxAmount            float64 `json:"max_amount" binding:"omitempty,gt=0"`
	PaymentMethodIDs     []int   `json:"payment_method_ids"`
	TradeTerms           string  `json:"trade_terms"`
	PaymentWindowMinutes int     `json:"payment_window_minutes"`
}

type InitiateTradeRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethodID int     `json:"payment_method_id" binding:"required"`
}

type LeaveFeedbackRequest struct {
	Rating  string `json:"rating" binding:"required,oneof=positive neutral negative"`
	Comment string `json:"comment"`
}
