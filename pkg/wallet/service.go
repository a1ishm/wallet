package wallet

import (
	"errors"
	"github.com/a1ishm/wallet/pkg/types"
)

// ErrPhoneRegistered v
var ErrPhoneRegistered = errors.New("phone already registered")

// ErrAccountNotFound v
var ErrAccountNotFound = errors.New("account not found")

// Service t
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

// RegisterAccount f
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// FindAccountById f
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var acc *types.Account
	for _, account := range s.accounts {
		if account.ID == accountID {
			acc = account
			return acc, nil
		}
	}
	return nil, ErrAccountNotFound
}