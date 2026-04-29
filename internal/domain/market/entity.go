package market

import "time"

type Rate struct {
	CryptoSymbol string    `json:"crypto_symbol"`
	FiatSymbol   string    `json:"fiat_symbol"`
	Price        float64   `json:"price"`
	Source       string    `json:"source"`
	AsOf         time.Time `json:"as_of"`
}
