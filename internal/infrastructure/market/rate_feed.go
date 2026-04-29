package market

import "context"

// Rate represents a crypto-fiat quote from a market data provider.
type Rate struct {
	CryptoSymbol string
	FiatSymbol   string
	Price        float64
	Provider     string
}

// RateFeed loads quotes from an external market data provider.
type RateFeed interface {
	GetRate(ctx context.Context, cryptoSymbol, fiatSymbol string) (*Rate, error)
}

// NoopRateFeed is a placeholder until a provider is integrated.
type NoopRateFeed struct {
	provider string
}

func NewCoinGeckoRateFeed() *NoopRateFeed {
	return &NoopRateFeed{provider: "coingecko"}
}

func NewBinanceRateFeed() *NoopRateFeed {
	return &NoopRateFeed{provider: "binance"}
}

func (f *NoopRateFeed) GetRate(context.Context, string, string) (*Rate, error) {
	return nil, nil
}
