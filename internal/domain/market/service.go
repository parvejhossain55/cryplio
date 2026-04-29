package market

import "errors"

func ValidateRate(rate *Rate) error {
	if rate == nil {
		return errors.New("rate is required")
	}
	if rate.Price <= 0 {
		return errors.New("invalid rate price")
	}
	return nil
}
