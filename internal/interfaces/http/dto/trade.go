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
