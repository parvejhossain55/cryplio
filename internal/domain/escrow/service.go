package escrow

import (
	"errors"
	"time"
)

func LockFunds(lock *Lock) error {
	if lock == nil {
		return errors.New("escrow lock is required")
	}
	lock.Status = StatusLocked
	lock.LockedAt = time.Now()
	return nil
}

func ReleaseFunds(lock *Lock) error {
	if lock == nil || lock.Status != StatusLocked {
		return errors.New("escrow is not locked")
	}
	now := time.Now()
	lock.Status = StatusReleased
	lock.ReleasedAt = &now
	return nil
}

func ReturnFunds(lock *Lock) error {
	if lock == nil || lock.Status != StatusLocked {
		return errors.New("escrow is not locked")
	}
	now := time.Now()
	lock.Status = StatusReturned
	lock.ReturnedAt = &now
	return nil
}
