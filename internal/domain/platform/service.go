package platform

import (
	"context"
	"fmt"
)

// CreateCryptoAssetInput holds the parameters for creating a cryptocurrency asset.
type CreateCryptoAssetInput struct {
	Symbol          string
	Name            string
	Blockchain      string
	ContractAddress *string
	Decimals        int
	MinConfirmation int
}

// UpdateCryptoAssetInput holds the parameters for updating a cryptocurrency asset.
type UpdateCryptoAssetInput struct {
	Symbol          string
	Name            string
	Blockchain      string
	ContractAddress *string
	Decimals        int
	MinConfirmation int
	IsActive        bool
}

// CreateFiatCurrencyInput holds the parameters for creating a fiat currency.
type CreateFiatCurrencyInput struct {
	Code   string
	Name   string
	Symbol string
}

// UpdateFiatCurrencyInput holds the parameters for updating a fiat currency.
type UpdateFiatCurrencyInput struct {
	Code     string
	Name     string
	Symbol   string
	IsActive bool
}

// CreatePaymentMethodInput holds the parameters for creating a payment method.
type CreatePaymentMethodInput struct {
	Code        string
	Name        string
	Category    PaymentCategory
	IconURL     *string
	Description *string
	SortOrder   int
}

// UpdatePaymentMethodInput holds the parameters for updating a payment method.
type UpdatePaymentMethodInput struct {
	Code        string
	Name        string
	Category    PaymentCategory
	IconURL     *string
	Description *string
	IsActive    bool
	SortOrder   int
}

// PlatformService defines business logic for platform management
type PlatformService interface {
	// Crypto Assets
	CreateCryptoAsset(ctx context.Context, input CreateCryptoAssetInput) (*CryptoAsset, error)
	GetCryptoAsset(ctx context.Context, id int) (*CryptoAsset, error)
	GetCryptoAssets(ctx context.Context, activeOnly bool, page, limit int) ([]*CryptoAsset, int, error)
	UpdateCryptoAsset(ctx context.Context, id int, input UpdateCryptoAssetInput) (*CryptoAsset, error)
	DeleteCryptoAsset(ctx context.Context, id int) error

	// Fiat Currencies
	CreateFiatCurrency(ctx context.Context, input CreateFiatCurrencyInput) (*FiatCurrency, error)
	GetFiatCurrency(ctx context.Context, id int) (*FiatCurrency, error)
	GetFiatCurrencies(ctx context.Context, activeOnly bool, page, limit int) ([]*FiatCurrency, int, error)
	UpdateFiatCurrency(ctx context.Context, id int, input UpdateFiatCurrencyInput) (*FiatCurrency, error)
	DeleteFiatCurrency(ctx context.Context, id int) error

	// Payment Methods
	CreatePaymentMethod(ctx context.Context, input CreatePaymentMethodInput) (*PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, id int) (*PaymentMethod, error)
	GetPaymentMethods(ctx context.Context, activeOnly bool, page, limit int) ([]*PaymentMethod, int, error)
	UpdatePaymentMethod(ctx context.Context, id int, input UpdatePaymentMethodInput) (*PaymentMethod, error)
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

func (s *platformService) CreateCryptoAsset(ctx context.Context, input CreateCryptoAssetInput) (*CryptoAsset, error) {
	if input.Symbol == "" || input.Name == "" || input.Blockchain == "" {
		return nil, fmt.Errorf("symbol, name, and blockchain are required")
	}
	if input.Decimals < 0 || input.MinConfirmation < 0 {
		return nil, fmt.Errorf("decimals and minConfirmation must be non-negative")
	}

	asset := &CryptoAsset{
		Symbol:          input.Symbol,
		Name:            input.Name,
		Blockchain:      input.Blockchain,
		ContractAddress: input.ContractAddress,
		Decimals:        input.Decimals,
		MinConfirmation: input.MinConfirmation,
		IsActive:        true,
	}

	if err := s.repo.CreateCryptoAsset(ctx, asset); err != nil {
		return nil, fmt.Errorf("create crypto asset: %w", err)
	}
	return asset, nil
}

func (s *platformService) GetCryptoAsset(ctx context.Context, id int) (*CryptoAsset, error) {
	return s.repo.GetCryptoAsset(ctx, id)
}

func (s *platformService) GetCryptoAssets(ctx context.Context, activeOnly bool, page, limit int) ([]*CryptoAsset, int, error) {
	offset := (page - 1) * limit
	return s.repo.GetCryptoAssets(ctx, activeOnly, limit, offset)
}

func (s *platformService) UpdateCryptoAsset(ctx context.Context, id int, input UpdateCryptoAssetInput) (*CryptoAsset, error) {
	if input.Symbol == "" || input.Name == "" || input.Blockchain == "" {
		return nil, fmt.Errorf("symbol, name, and blockchain are required")
	}
	if input.Decimals < 0 || input.MinConfirmation < 0 {
		return nil, fmt.Errorf("decimals and minConfirmation must be non-negative")
	}

	asset := &CryptoAsset{
		ID:              id,
		Symbol:          input.Symbol,
		Name:            input.Name,
		Blockchain:      input.Blockchain,
		ContractAddress: input.ContractAddress,
		Decimals:        input.Decimals,
		MinConfirmation: input.MinConfirmation,
		IsActive:        input.IsActive,
	}

	if err := s.repo.UpdateCryptoAsset(ctx, asset); err != nil {
		return nil, fmt.Errorf("update crypto asset: %w", err)
	}
	return asset, nil
}

func (s *platformService) DeleteCryptoAsset(ctx context.Context, id int) error {
	return s.repo.DeleteCryptoAsset(ctx, id)
}

// Fiat Currencies

func (s *platformService) CreateFiatCurrency(ctx context.Context, input CreateFiatCurrencyInput) (*FiatCurrency, error) {
	if input.Code == "" || input.Name == "" || input.Symbol == "" {
		return nil, fmt.Errorf("code, name, and symbol are required")
	}

	currency := &FiatCurrency{
		Code:     input.Code,
		Name:     input.Name,
		Symbol:   input.Symbol,
		IsActive: true,
	}

	if err := s.repo.CreateFiatCurrency(ctx, currency); err != nil {
		return nil, fmt.Errorf("create fiat currency: %w", err)
	}
	return currency, nil
}

func (s *platformService) GetFiatCurrency(ctx context.Context, id int) (*FiatCurrency, error) {
	return s.repo.GetFiatCurrency(ctx, id)
}

func (s *platformService) GetFiatCurrencies(ctx context.Context, activeOnly bool, page, limit int) ([]*FiatCurrency, int, error) {
	offset := (page - 1) * limit
	return s.repo.GetFiatCurrencies(ctx, activeOnly, limit, offset)
}

func (s *platformService) UpdateFiatCurrency(ctx context.Context, id int, input UpdateFiatCurrencyInput) (*FiatCurrency, error) {
	if input.Code == "" || input.Name == "" || input.Symbol == "" {
		return nil, fmt.Errorf("code, name, and symbol are required")
	}

	currency := &FiatCurrency{
		ID:       id,
		Code:     input.Code,
		Name:     input.Name,
		Symbol:   input.Symbol,
		IsActive: input.IsActive,
	}

	if err := s.repo.UpdateFiatCurrency(ctx, currency); err != nil {
		return nil, fmt.Errorf("update fiat currency: %w", err)
	}
	return currency, nil
}

func (s *platformService) DeleteFiatCurrency(ctx context.Context, id int) error {
	return s.repo.DeleteFiatCurrency(ctx, id)
}

// Payment Methods

func (s *platformService) CreatePaymentMethod(ctx context.Context, input CreatePaymentMethodInput) (*PaymentMethod, error) {
	if input.Code == "" || input.Name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	method := &PaymentMethod{
		Code:        input.Code,
		Name:        input.Name,
		Category:    input.Category,
		IconURL:     input.IconURL,
		Description: input.Description,
		IsActive:    true,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.CreatePaymentMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("create payment method: %w", err)
	}
	return method, nil
}

func (s *platformService) GetPaymentMethod(ctx context.Context, id int) (*PaymentMethod, error) {
	return s.repo.GetPaymentMethod(ctx, id)
}

func (s *platformService) GetPaymentMethods(ctx context.Context, activeOnly bool, page, limit int) ([]*PaymentMethod, int, error) {
	offset := (page - 1) * limit
	return s.repo.GetPaymentMethods(ctx, activeOnly, limit, offset)
}

func (s *platformService) UpdatePaymentMethod(ctx context.Context, id int, input UpdatePaymentMethodInput) (*PaymentMethod, error) {
	if input.Code == "" || input.Name == "" {
		return nil, fmt.Errorf("code and name are required")
	}

	method := &PaymentMethod{
		ID:          id,
		Code:        input.Code,
		Name:        input.Name,
		Category:    input.Category,
		IconURL:     input.IconURL,
		Description: input.Description,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.UpdatePaymentMethod(ctx, method); err != nil {
		return nil, fmt.Errorf("update payment method: %w", err)
	}
	return method, nil
}

func (s *platformService) DeletePaymentMethod(ctx context.Context, id int) error {
	return s.repo.DeletePaymentMethod(ctx, id)
}
