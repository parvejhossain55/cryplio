package platform

import "context"

// PlatformRepository defines methods for platform data management
type PlatformRepository interface {
	// Crypto Assets
	CreateCryptoAsset(ctx context.Context, asset *CryptoAsset) error
	GetCryptoAsset(ctx context.Context, id int) (*CryptoAsset, error)
	GetCryptoAssets(ctx context.Context, activeOnly bool, limit, offset int) ([]*CryptoAsset, int, error)
	UpdateCryptoAsset(ctx context.Context, asset *CryptoAsset) error
	DeleteCryptoAsset(ctx context.Context, id int) error

	// Fiat Currencies
	CreateFiatCurrency(ctx context.Context, currency *FiatCurrency) error
	GetFiatCurrency(ctx context.Context, id int) (*FiatCurrency, error)
	GetFiatCurrencies(ctx context.Context, activeOnly bool, limit, offset int) ([]*FiatCurrency, int, error)
	UpdateFiatCurrency(ctx context.Context, currency *FiatCurrency) error
	DeleteFiatCurrency(ctx context.Context, id int) error

	// Payment Methods
	CreatePaymentMethod(ctx context.Context, method *PaymentMethod) error
	GetPaymentMethod(ctx context.Context, id int) (*PaymentMethod, error)
	GetPaymentMethods(ctx context.Context, activeOnly bool, limit, offset int) ([]*PaymentMethod, int, error)
	UpdatePaymentMethod(ctx context.Context, method *PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id int) error
}
