package wallet

import (
	"errors"
	"fmt"
	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
)

// ErrPhoneRegistered v
var ErrPhoneRegistered = errors.New("phone already registered")

// ErrAccountNotFound v
var ErrAccountNotFound = errors.New("account not found")

// ErrPaymentNotFound v
var ErrPaymentNotFound = errors.New("payment not found")

// ErrAmountMustBePositive v
var ErrAmountMustBePositive = errors.New("amount must be positive")

// ErrNotEnoughBalance v
var ErrNotEnoughBalance = errors.New("not enough balance")

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
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// FindAccountByID f
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

// Pay f
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindPaymentByID f
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var paym *types.Payment
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			paym = payment
			return paym, nil
		}
	}
	return nil, ErrPaymentNotFound
}

// Reject f
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		fmt.Println(err)
		return err
	}

	payment.Status = types.PaymentStatusFail
	
	account, errr := s.FindAccountByID(payment.AccountID)

	if errr != nil {
		fmt.Println(errr)
		return errr
	}

	account.Balance += payment.Amount
	return nil
}

// Deposit f
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

// Repeat f
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	paym, errr := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	
	if errr != nil {
		return nil, errr
	}

	return paym, errr
}
