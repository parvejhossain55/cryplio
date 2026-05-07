package market

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// RateService provides market exchange rates
type RateService interface {
	GetRate(ctx context.Context, cryptoSymbol, fiatCode string) (*Rate, error)
	GetRates(ctx context.Context, cryptoSymbol string) ([]Rate, error)
	GetAllRates(ctx context.Context) ([]Rate, error)
}

type rateService struct {
	// In production, this would connect to a price feed API
	baseURL string
}

// NewRateService creates a new rate service
func NewRateService() RateService {
	return &rateService{}
}

func (s *rateService) GetRate(ctx context.Context, cryptoSymbol, fiatCode string) (*Rate, error) {
	if cryptoSymbol == "" || fiatCode == "" {
		return nil, errors.New("crypto symbol and fiat code are required")
	}

	// For MVP, return mock rates based on fiat currency
	// In production, this would fetch from a live API like CoinGecko/CoinMarketCap
	var price float64
	switch fiatCode {
	case "BDT":
		price = 118.50 + rand.Float64()*2.0 // Mock USDT/BDT rate
	case "PKR":
		price = 278.30 + rand.Float64()*3.0 // Mock USDT/PKR rate
	case "USD":
		price = 1.00 + rand.Float64()*0.01 // Mock USDT/USD rate
	case "NGN":
		price = 1550.00 + rand.Float64()*20.0 // Mock USDT/NGN rate
	case "EGP":
		price = 49.50 + rand.Float64()*0.5 // Mock USDT/EGP rate
	default:
		price = 1.00 + rand.Float64()*0.01
	}

	return &Rate{
		CryptoSymbol: cryptoSymbol,
		FiatSymbol:   fiatCode,
		Price:        price,
		Source:       "mock",
		AsOf:         time.Now(),
	}, nil
}

func (s *rateService) GetRates(ctx context.Context, cryptoSymbol string) ([]Rate, error) {
	if cryptoSymbol == "" {
		return nil, errors.New("crypto symbol is required")
	}

	fiats := []string{"BDT", "PKR", "USD", "NGN", "EGP"}
	rates := make([]Rate, 0, len(fiats))
	for _, fiat := range fiats {
		rate, err := s.GetRate(ctx, cryptoSymbol, fiat)
		if err != nil {
			return nil, fmt.Errorf("get rate for %s: %w", fiat, err)
		}
		rates = append(rates, *rate)
	}
	return rates, nil
}

func (s *rateService) GetAllRates(ctx context.Context) ([]Rate, error) {
	cryptos := []string{"USDT", "USDC", "BTC", "ETH"}
	fiats := []string{"BDT", "PKR", "USD", "NGN", "EGP"}

	allRates := make([]Rate, 0, len(cryptos)*len(fiats))
	for _, crypto := range cryptos {
		for _, fiat := range fiats {
			rate, err := s.GetRate(ctx, crypto, fiat)
			if err != nil {
				continue // Skip invalid rates
			}
			allRates = append(allRates, *rate)
		}
	}
	return allRates, nil
}

func ValidateRate(rate *Rate) error {
	if rate == nil {
		return errors.New("rate is required")
	}
	if rate.Price <= 0 {
		return errors.New("invalid rate price")
	}
	return nil
}
