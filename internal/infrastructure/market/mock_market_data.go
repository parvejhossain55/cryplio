package market

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"cryplio/internal/interfaces/websocket"
)

// MockMarketDataProvider provides mock real-time market data
type MockMarketDataProvider struct {
	rates       map[string]*MarketRate
	mutex       sync.RWMutex
	wsService   websocket.Service
	subscribers map[string]bool
	lastUpdate  time.Time
}

// MarketRate represents a market rate
type MarketRate struct {
	CryptoSymbol string    `json:"crypto_symbol"`
	FiatSymbol   string    `json:"fiat_symbol"`
	Price        float64   `json:"price"`
	Source       string    `json:"source"`
	AsOf         time.Time `json:"as_of"`
	Change24h    float64   `json:"change_24h"`
	Volume24h    float64   `json:"volume_24h"`
	High24h      float64   `json:"high_24h"`
	Low24h       float64   `json:"low_24h"`
}

// NewMockMarketDataProvider creates a new mock market data provider
func NewMockMarketDataProvider(wsService websocket.Service) *MockMarketDataProvider {
	provider := &MockMarketDataProvider{
		rates:       make(map[string]*MarketRate),
		wsService:   wsService,
		subscribers: make(map[string]bool),
	}

	// Initialize with base rates
	provider.initializeRates()

	return provider
}

// initializeRates sets up initial market rates
func (m *MockMarketDataProvider) initializeRates() {
	baseRates := map[string]map[string]float64{
		"BTC": {
			"USD": 45000.0,
			"EUR": 41000.0,
			"GBP": 36000.0,
			"JPY": 6500000.0,
			"BDT": 3800000.0,
		},
		"ETH": {
			"USD": 3000.0,
			"EUR": 2750.0,
			"GBP": 2400.0,
			"JPY": 435000.0,
			"BDT": 255000.0,
		},
		"USDT": {
			"USD": 1.0,
			"EUR": 0.92,
			"GBP": 0.80,
			"JPY": 145.0,
			"BDT": 85.0,
		},
		"USDC": {
			"USD": 1.0,
			"EUR": 0.92,
			"GBP": 0.80,
			"JPY": 145.0,
			"BDT": 85.0,
		},
		"BNB": {
			"USD": 320.0,
			"EUR": 295.0,
			"GBP": 255.0,
			"JPY": 46500.0,
			"BDT": 27200.0,
		},
		"SOL": {
			"USD": 105.0,
			"EUR": 97.0,
			"GBP": 84.0,
			"JPY": 15250.0,
			"BDT": 8925.0,
		},
		"ADA": {
			"USD": 0.38,
			"EUR": 0.35,
			"GBP": 0.30,
			"JPY": 55.0,
			"BDT": 32.0,
		},
		"DOT": {
			"USD": 7.5,
			"EUR": 6.9,
			"GBP": 6.0,
			"JPY": 1085.0,
			"BDT": 635.0,
		},
	}

	cryptos := []string{"BTC", "ETH", "USDT", "USDC", "BNB", "SOL", "ADA", "DOT"}
	fiats := []string{"USD", "EUR", "GBP", "JPY", "BDT"}

	for _, crypto := range cryptos {
		for _, fiat := range fiats {
			if basePrice, exists := baseRates[crypto][fiat]; exists {
				key := fmt.Sprintf("%s-%s", crypto, fiat)
				m.rates[key] = &MarketRate{
					CryptoSymbol: crypto,
					FiatSymbol:   fiat,
					Price:        basePrice + (rand.Float64()-0.5)*basePrice*0.02, // ±1% variation
					Source:       "mock_exchange",
					AsOf:         time.Now(),
					Change24h:    (rand.Float64() - 0.5) * 10, // ±5% change
					Volume24h:    rand.Float64() * 1000000000,
					High24h:      basePrice * (1 + rand.Float64()*0.05),
					Low24h:       basePrice * (1 - rand.Float64()*0.05),
				}
			}
		}
	}

	m.lastUpdate = time.Now()
}

// StartRealTimeUpdates starts real-time market data updates
func (m *MockMarketDataProvider) StartRealTimeUpdates(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Update every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.updateRates()
		}
	}
}

// updateRates updates market rates with realistic variations
func (m *MockMarketDataProvider) updateRates() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, rate := range m.rates {
		// Simulate price movement (±0.5% per update)
		changePercent := (rand.Float64() - 0.5) * 0.01
		newPrice := rate.Price * (1 + changePercent)

		// Update rate
		rate.Price = newPrice
		rate.Change24h += changePercent * 100
		rate.AsOf = time.Now()
		rate.Volume24h *= (1 + (rand.Float64()-0.5)*0.1) // ±5% volume change

		// Update high/low if necessary
		if newPrice > rate.High24h {
			rate.High24h = newPrice
		}
		if newPrice < rate.Low24h {
			rate.Low24h = newPrice
		}

		// Broadcast update via WebSocket
		if m.wsService != nil {
			marketUpdate := websocket.MarketUpdate{
				CryptoSymbol: rate.CryptoSymbol,
				FiatSymbol:   rate.FiatSymbol,
				Price:        rate.Price,
				Change24h:    rate.Change24h,
				Timestamp:    rate.AsOf.Format(time.RFC3339),
			}

			m.wsService.BroadcastMessage("market_update", marketUpdate, "")
		}
	}

	m.lastUpdate = time.Now()
	log.Printf("Updated %d market rates", len(m.rates))
}

// GetRate returns the current rate for a crypto-fiat pair
func (m *MockMarketDataProvider) GetRate(cryptoSymbol, fiatSymbol string) (*MarketRate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := fmt.Sprintf("%s-%s", cryptoSymbol, fiatSymbol)
	rate, exists := m.rates[key]
	if !exists {
		return nil, fmt.Errorf("rate not found for %s-%s", cryptoSymbol, fiatSymbol)
	}

	// Return a copy to avoid concurrent modification
	rateCopy := *rate
	return &rateCopy, nil
}

// GetRates returns all current rates
func (m *MockMarketDataProvider) GetRates() ([]*MarketRate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	rates := make([]*MarketRate, 0, len(m.rates))
	for _, rate := range m.rates {
		rateCopy := *rate
		rates = append(rates, &rateCopy)
	}

	return rates, nil
}

// GetRatesByFiat returns rates for a specific fiat currency
func (m *MockMarketDataProvider) GetRatesByFiat(fiatSymbol string) ([]*MarketRate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var rates []*MarketRate
	for _, rate := range m.rates {
		if rate.FiatSymbol == fiatSymbol {
			rateCopy := *rate
			rates = append(rates, &rateCopy)
		}
	}

	return rates, nil
}

// GetRatesByCrypto returns rates for a specific cryptocurrency
func (m *MockMarketDataProvider) GetRatesByCrypto(cryptoSymbol string) ([]*MarketRate, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var rates []*MarketRate
	for _, rate := range m.rates {
		if rate.CryptoSymbol == cryptoSymbol {
			rateCopy := *rate
			rates = append(rates, &rateCopy)
		}
	}

	return rates, nil
}

// SubscribeToMarketUpdates subscribes a client to market updates
func (m *MockMarketDataProvider) SubscribeToMarketUpdates(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.subscribers[clientID] = true
	log.Printf("Client %s subscribed to market updates", clientID)
}

// UnsubscribeFromMarketUpdates unsubscribes a client from market updates
func (m *MockMarketDataProvider) UnsubscribeFromMarketUpdates(clientID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.subscribers, clientID)
	log.Printf("Client %s unsubscribed from market updates", clientID)
}

// GetStats returns market data statistics
func (m *MockMarketDataProvider) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	totalVolume := 0.0
	gainers := 0
	losers := 0

	for _, rate := range m.rates {
		totalVolume += rate.Volume24h
		if rate.Change24h > 0 {
			gainers++
		} else if rate.Change24h < 0 {
			losers++
		}
	}

	return map[string]interface{}{
		"total_rates":  len(m.rates),
		"total_volume": totalVolume,
		"gainers":      gainers,
		"losers":       losers,
		"last_update":  m.lastUpdate.Format(time.RFC3339),
		"subscribers":  len(m.subscribers),
	}
}

// SimulateMarketEvent simulates a market event (e.g., sudden price movement)
func (m *MockMarketDataProvider) SimulateMarketEvent(cryptoSymbol, fiatSymbol string, priceChangePercent float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := fmt.Sprintf("%s-%s", cryptoSymbol, fiatSymbol)
	rate, exists := m.rates[key]
	if !exists {
		return
	}

	// Apply sudden price change
	rate.Price *= (1 + priceChangePercent/100)
	rate.Change24h += priceChangePercent
	rate.AsOf = time.Now()

	// Update high/low
	if rate.Price > rate.High24h {
		rate.High24h = rate.Price
	}
	if rate.Price < rate.Low24h {
		rate.Low24h = rate.Price
	}

	// Broadcast the significant update
	if m.wsService != nil {
		marketUpdate := websocket.MarketUpdate{
			CryptoSymbol: rate.CryptoSymbol,
			FiatSymbol:   rate.FiatSymbol,
			Price:        rate.Price,
			Change24h:    rate.Change24h,
			Timestamp:    rate.AsOf.Format(time.RFC3339),
		}

		m.wsService.BroadcastMessage("market_event", marketUpdate, "")
	}

	log.Printf("Market event simulated: %s-%s changed by %.2f%%", cryptoSymbol, fiatSymbol, priceChangePercent)
}

// ExportRates exports rates to JSON format
func (m *MockMarketDataProvider) ExportRates() ([]byte, error) {
	rates, err := m.GetRates()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(rates, "", "  ")
}
