package dispute

import "errors"

func ValidateRaise(dispute *Dispute) error {
	if dispute == nil {
		return errors.New("dispute is required")
	}
	return nil
}
