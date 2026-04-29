package referral

import "errors"

func ValidateReferral(referral *Referral) error {
	if referral == nil {
		return errors.New("referral is required")
	}
	return nil
}
