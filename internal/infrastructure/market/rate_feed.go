package market

import (
	"context"

	domainmarket "cryplio/internal/domain/market"
)

// RateFeed is a source of live exchange rate data.
// Implementations should connect to external price APIs (CoinGecko, Binance, etc.).
type RateFeed interface {
	GetRate(ctx context.Context, cryptoSymbol, fiatSymbol string) (*domainmarket.Rate, error)
}

// NoopRateFeed is a placeholder that returns nil for all rate queries.
// Replace with a real implementation (CoinGeckoRateFeed, BinanceRateFeed) for production.
type NoopRateFeed struct {
	provider string
}

func NewCoinGeckoRateFeed() *NoopRateFeed { return &NoopRateFeed{provider: "coingecko"} }
func NewBinanceRateFeed() *NoopRateFeed   { return &NoopRateFeed{provider: "binance"} }

func (f *NoopRateFeed) GetRate(context.Context, string, string) (*domainmarket.Rate, error) {
	return nil, nil
}
