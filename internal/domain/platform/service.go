package platform

import (
	"context"
	"fmt"
)

// PlatformService defines business logic for platform management
type PlatformService interface {
	// Crypto Assets
	CreateCryptoAsset(ctx context.Context, symbol, name, blockchain string, contractAddress *string, decimals, minConfirmation int) (*CryptoAsset, error)
	GetCryptoAsset(ctx context.Context, id int) (*CryptoAsset, error)
	GetCryptoAssets(ctx context.Context, activeOnly bool) ([]*CryptoAsset, error)
	UpdateCryptoAsset(ctx context.Context, id int, symbol, name, blockchain string, contractAddress *string, decimals, minConfirmation int, isActive bool) (*CryptoAsset, error)
	DeleteCryptoAsset(ctx context.Context, id int) error

	// Fiat Currencies
	CreateFiatCurrency(ctx context.Context, code, name, symbol string) (*FiatCurrency, error)
	GetFiatCurrency(ctx context.Context, id int) (*FiatCurrency, error)
	GetFiatCurrencies(ctx context.Context, activeOnly bool) ([]*FiatCurrency, error)
	UpdateFiatCurrency(ctx context.Context, id int, code, name, symbol string, isActive bool) (*FiatCurrency, error)
	DeleteFiatCurrency(ctx context.Context, id int) error

	// Payment Methods
	CreatePaymentMethod(ctx context.Context, code, name string, category PaymentCategory, iconURL, description *string, sortOrder int) (*PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id int) (*PaymentMethod, error)
	GetPaymentMethods(ctx context.Context, activeOnly bool) ([]*PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, id int, code, name string, category PaymentCategory, iconURL, description *string, isActive bool, sortOrder int) (*PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, id int) error
}

// platformService implements PlatformService
type platformService struct {
	repo PlatformRepository
}

// NewPlatformService creates a new platform service
func NewPlatformService(repo PlatformRepository) PlatformService {
	return &platformService{repo: repo}
}

// Crypto Assets
func (s *platformService) CreateCryptoAsset(ctx context.Context, symbol, name, blockchain string, contractAddress *string, decimals, minConfirmation int) (*CryptoAsset, error) {
	if symbol == "" || name == "" || blockchain == "" {
		return nil, fmt.Errorf("symbol, name, and blockchain are required")
	}

	if decimals < 0 || minConfirmation < 0 {
		return nil, fmt.Errorf("decimals and minConfirmation must be non-negative")
	}

	asset := &CryptoAsset{
		Symbol:          symbol,
		Name:            name,
		Blockchain:      blockchain,
		ContractAddress: contractAddress,
		Decimals:        decimals,
		MinConfirmation: minConfirmation,
		IsActive:        true,
	}

	err := s.repo.CreateCryptoAsset(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("create crypto asset: %w", err)
	}

	return asset, nil
}

func (s *platformService) GetCryptoAsset(ctx context.Context, id int) (*CryptoAsset, error) {
	return s.repo.GetCryptoAsset(ctx, id)
}

func (s *platformService) GetCryptoAssets(ctx context.Context, activeOnly bool) ([]*CryptoAsset, error) {
	return s.repo.GetCryptoAssets(ctx, activeOnly)
}

func (s *platformService) UpdateCryptoAsset(ctx context.Context, id int, symbol, name, blockchain string, contractAddress *string, decimals, minConfirmation int, isActive bool) (*CryptoAsset, error) {
	if symbol == "" || name == "" || blockchain == "" {
		return nil, fmt.Errorf("symbol, name, and blockchain are required")
	}

	if decimals < 0 || minConfirmation < 0 {
		return nil, fmt.Errorf("decimals and minConfirmation must be non-negative")
	}

	asset := &CryptoAsset{
		ID:              id,
		Symbol:          symbol,
		Name:            name,
		Blockchain:      blockchain,
		ContractAddress: contractAddress,
		Decimals:        decimals,
		MinConfirmation: minConfirmation,
		IsActive:        isActive,
	}

	err := s.repo.UpdateCryptoAsset(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("update crypto asset: %w", err)
	}

	return asset, nil
}

func (s *platformService) DeleteCryptoAsset(ctx context.Context, id int) error {
	return s.repo.DeleteCryptoAsset(ctx, id)
}

// Fiat Currencies
func (s *platformService) CreateFiatCurrency(ctx context.Context, code, name, symbol string) (*FiatCurrency, error) {
	if code == "" || name == "" || symbol == "" {
		return nil, fmt.Errorf("code, name, and symbol are required")
	}

	currency := &FiatCurrency{
		Code:     code,
		Name:     name,
		Symbol:   symbol,
		IsActive: true,
	}

	err := s.repo.CreateFiatCurrency(ctx, currency)
	if err != nil {
		return nil, fmt.Errorf("create fiat currency: %w", err)
	}

	return currency, nil
}

func (s *platformService) GetFiatCurrency(ctx context.Context, id int) (*FiatCurrency, error) {
	return s.repo.GetFiatCurrency(ctx, id)
}

func (s *platformService) GetFiatCurrencies(ctx context.Context, activeOnly bool) ([]*FiatCurrency, error) {
	return s.repo.GetFiatCurrencies(ctx, activeOnly)
}

func (s *platformService) UpdateFiatCurrency(ctx context.Context, id int, code, name, symbol string, isActive bool) (*FiatCurrency, error) {
	if code == "" || name == "" || symbol == "" {
		return nil, fmt.Errorf("code, name, and symbol are required")
	}

	currency := &FiatCurrency{
		ID:       id,
		Code:     code,
		Name:     name,
		Symbol:   symbol,
		IsActive: isActive,
	}

	err := s.repo.UpdateFiatCurrency(ctx, currency)
	if err != nil {
		return nil, fmt.Errorf("update fiat currency: %w", err)
	}

	return currency, nil
}

func (s *platformService) DeleteFiatCurrency(ctx context.Context, id int) error {
	return s.repo.DeleteFiatCurrency(ctx, id)
}

// Payment Methods
func (s *platformService) CreatePaymentMethod(ctx context.Context, code, name string, category PaymentCategory, iconURL, description *string, sortOrder int) (*PaymentMethod, error) {
	if code == "" || name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	method := &PaymentMethod{
		Code:        code,
		Name:        name,
		Category:    category,
		IconURL:     iconURL,
		Description: description,
		IsActive:    true,
		SortOrder:   sortOrder,
	}

	err := s.repo.CreatePaymentMethod(ctx, method)
	if err != nil {
		return nil, fmt.Errorf("create payment method: %w", err)
	}

	return method, nil
}

func (s *platformService) GetPaymentMethod(ctx context.Context, id int) (*PaymentMethod, error) {
	return s.repo.GetPaymentMethod(ctx, id)
}

func (s *platformService) GetPaymentMethods(ctx context.Context, activeOnly bool) ([]*PaymentMethod, error) {
	return s.repo.GetPaymentMethods(ctx, activeOnly)
}

func (s *platformService) UpdatePaymentMethod(ctx context.Context, id int, code, name string, category PaymentCategory, iconURL, description *string, isActive bool, sortOrder int) (*PaymentMethod, error) {
	if code == "" || name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	method := &PaymentMethod{
		ID:          id,
		Code:        code,
		Name:        name,
		Category:    category,
		IconURL:     iconURL,
		Description: description,
		IsActive:    isActive,
		SortOrder:   sortOrder,
	}

	err := s.repo.UpdatePaymentMethod(ctx, method)
	if err != nil {
		return nil, fmt.Errorf("update payment method: %w", err)
	}

	return method, nil
}

func (s *platformService) DeletePaymentMethod(ctx context.Context, id int) error {
	return s.repo.DeletePaymentMethod(ctx, id)
}
