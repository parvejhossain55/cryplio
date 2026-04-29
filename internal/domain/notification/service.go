package notification

import "errors"

func ValidateNotification(notification *Notification) error {
	if notification == nil {
		return errors.New("notification is required")
	}
	return nil
}
