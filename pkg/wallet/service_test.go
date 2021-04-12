package wallet

import (
	"fmt"
	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"testing"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
}

func TestService_FindAccountByID(t *testing.T) {
	s := newTestService()
	exp, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	got, err := s.FindAccountByID(1)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindAccountByID(rand.Int63())
	if err == nil {
		t.Errorf("FindAccountByID(): must return error, returned nil")
		return
	}

	if err != ErrAccountNotFound {
		t.Errorf("FindAccountByID(): must return ErrAccountNotFound, returned %v", err)
	}

}

func TestService_Reject(t *testing.T) {
	s := newTestService()
	account, err := s.addAccountWithBalance("+992000000001", 1000)
	if err != nil {
		t.Error(err)
		return
	}

	payment, err := s.Pay(account.ID, 1000, "food")
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Reject(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	got := account.Balance
	var exp types.Money = 1000

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_Reject_notFound(t *testing.T) {
	s := newTestService()
	account, err := s.addAccountWithBalance("+992000000001", 1000)
	if err != nil {
		t.Error(err)
		return
	}

	s.Pay(account.ID, 1000, "food")

	err = s.Reject("999")
	if err == nil {
		t.Errorf("Reject(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("Reject(): must return ErrPaymentNotFound, returned %v", err)
		return
	}
}

func TestService_Repeat(t *testing.T) {
	s := newTestService()
	account, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.Repeat(payments[0].ID)
	if err != nil {
		t.Error(err)
		return
	}

	got := account.Balance
	exp := types.Money(8_000_00)

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_Repeat_notFound(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.Repeat(uuid.New().String())
	if err == nil {
		t.Errorf("Repeat(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("Repeat(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_FavoritePayment(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	favorite, err := s.FavoritePayment(payments[0].ID, "fav")
	if err != nil {
		t.Error(err)
		return
	}

	got := favorite.ID
	exp := payments[0].ID

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

