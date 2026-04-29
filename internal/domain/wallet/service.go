package wallet

import "errors"

func DepositFunds(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid deposit")
	}
	wallet.Credit(amount)
	return nil
}

func WithdrawFunds(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid withdrawal")
	}
	return wallet.Debit(amount)
}

func LockBalance(wallet *Wallet, amount float64) error {
	if wallet == nil || amount <= 0 {
		return errors.New("invalid lock amount")
	}
	return wallet.Lock(amount)
}
