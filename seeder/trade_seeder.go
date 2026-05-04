package seeder

import (
	"context"
	"math/rand"
	"time"

	domainidentity "cryplio/internal/domain/identity"
	domaintrading "cryplio/internal/domain/trading"

	"github.com/google/uuid"
)

func (s *Seeder) SeedTradeAds(ctx context.Context, users []*domainidentity.User, cryptoMap, fiatMap, pmMap map[string]int) ([]*domaintrading.TradeAd, error) {
	var ads []*domaintrading.TradeAd
	for _, user := range users {
		if !user.IsMerchant && user.Username != "AliceTrader" && user.Username != "DianaCoin" {
			continue
		}

		for i := 0; i < 4; i++ {
			adType := domaintrading.AdTypeSell
			if i%2 == 1 {
				adType = domaintrading.AdTypeBuy
			}

			cryptoSym := "USDT"
			if i == 2 {
				cryptoSym = "BTC"
			}
			if i == 3 {
				cryptoSym = "ETH"
			}

			fiatCode := "USD"
			if i%3 == 1 {
				fiatCode = "BDT"
			}
			if i%3 == 2 {
				fiatCode = "NGN"
			}

			price := 1.0
			if cryptoSym == "BTC" {
				price = 67000.0
			} else if cryptoSym == "ETH" {
				price = 3400.0
			}
			if fiatCode == "BDT" {
				price *= 118.0
			} else if fiatCode == "NGN" {
				price *= 1450.0
			}

			ad := &domaintrading.TradeAd{
				AdID: uuid.New(), UserID: user.UserID, Type: adType,
				CryptoID: cryptoMap[cryptoSym], FiatID: fiatMap[fiatCode],
				PriceType: domaintrading.PriceTypeFixed, Price: price,
				MinAmount: 10.0, MaxAmount: 2000.0,
				PaymentMethods:       []int{pmMap["bkash"], pmMap["bank_transfer"]},
				PaymentWindowMinutes: 15, RequiresKYCLevel: domainidentity.KYCLevel1,
				IsPublic: true, IsPaused: false, Timezone: "UTC", Status: domaintrading.TradeAdStatusActive,
				PublishedAt: time.Now(),
			}
			if fiatCode == "USD" {
				ad.PaymentMethods = []int{pmMap["wise"], pmMap["paypal"]}
			}

			if err := s.tradeRepo.CreateAd(ctx, ad); err != nil {
				return nil, err
			}
			ads = append(ads, ad)
		}
	}
	return ads, nil
}

func (s *Seeder) SeedTrades(ctx context.Context, users []*domainidentity.User, ads []*domaintrading.TradeAd) ([]*domaintrading.Trade, error) {
	var trades []*domaintrading.Trade
	for i := 0; i < 20; i++ {
		ad := ads[rand.Intn(len(ads))]
		buyer := users[rand.Intn(len(users))]
		if buyer.UserID == ad.UserID {
			buyer = users[(rand.Intn(len(users)-1)+1)%len(users)]
		}

		buyerID, sellerID := buyer.UserID, ad.UserID
		if ad.Type == domaintrading.AdTypeBuy {
			buyerID, sellerID = ad.UserID, buyer.UserID
		}

		status := domaintrading.TradeStatusCompleted
		if i == 0 {
			status = domaintrading.TradeStatusActive
		}
		if i == 1 {
			status = domaintrading.TradeStatusDisputed
		}

		trade := &domaintrading.Trade{
			TradeID: uuid.New(), AdID: ad.AdID, BuyerID: buyerID, SellerID: sellerID,
			CryptoID: ad.CryptoID, FiatID: ad.FiatID,
			CryptoAmount: 150.0 / ad.Price, FiatAmount: 150.0,
			ExchangeRate: ad.Price, PaymentMethod: ad.PaymentMethods[0],
			PriceType: ad.PriceType, AgreedPrice: ad.Price, Status: status,
			PaymentWindowMinutes: ad.PaymentWindowMinutes, CreatedAt: time.Now().Add(-time.Duration(i*2) * time.Hour),
		}
		if status == domaintrading.TradeStatusCompleted {
			compAt := trade.CreatedAt.Add(20 * time.Minute)
			trade.CompletedAt = &compAt
		}

		if err := s.tradeRepo.CreateTrade(ctx, trade); err != nil {
			continue
		}
		trades = append(trades, trade)

		if status == domaintrading.TradeStatusCompleted {
			comment := "Fast and reliable!"
			if i%2 == 0 {
				comment = "Highly recommended merchant."
			}
			_, _ = s.db.ExecContext(ctx, `
				INSERT INTO trade_feedback (feedback_id, trade_id, from_user_id, to_user_id, rating, comment, created_at)
				VALUES ($1, $2, $3, $4, 'positive', $5, NOW()) ON CONFLICT DO NOTHING`,
				uuid.New(), trade.TradeID, buyerID, sellerID, comment)
		}
	}
	return trades, nil
}

func (s *Seeder) SeedDisputes(ctx context.Context, trades []*domaintrading.Trade) error {
	for _, t := range trades {
		if t.Status == domaintrading.TradeStatusDisputed {
			_, err := s.db.ExecContext(ctx, `
				INSERT INTO disputes (dispute_id, trade_id, raised_by, reason_code, reason_text, status, created_at)
				VALUES ($1, $2, $3, $4, $5, 'pending', NOW())
				ON CONFLICT DO NOTHING`,
				uuid.New(), t.TradeID, t.BuyerID, "payment_not_received", "Buyer says they paid but seller hasn't released.",
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
