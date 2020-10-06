package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}

// newTestService f
func newTestService() *testService {
	return &testService{Service: &Service{}}
}

// testAccount t
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

// addAccountWithBalance f
func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	// выполняем платежи
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)

	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	repeated, err := s.Repeat(payment.ID)

	if err != nil {
		t.Error(err)
		return
	}

	got := repeated.Amount

	expected := types.Money(1_000_00)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, got)
	}
}

func TestService_Repeat_notFound(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(testAccount{
		phone:   "+992000000001",
		balance: 10_000_00,
		payments: []struct {
			amount   types.Money
			category types.PaymentCategory
		}{},
	})

	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.Repeat(uuid.New().String())

	if err == nil {
		t.Error(err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("AAA ERROR: %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)

	if err != nil {
		t.Errorf("AXAXXAX %v", err)
		return
	}

	acc, err := s.FindAccountByID(payment.AccountID)

	if err != nil {
		t.Error(err)
		return
	}

	res := acc.Balance
	expected := types.Money(10_000_00)

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, res)
	}
}

func TestService_Reject_notFound(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(testAccount{
		phone:   "+992000000001",
		balance: 10_000_00,
		payments: []struct {
			amount   types.Money
			category types.PaymentCategory
		}{},
	})

	if err != nil {
		t.Error(err)
		return
	}

	if len(payments) > 0 {
		t.Error(err)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService();
	_, payments, err := s.addAccount(defaultTestAccount)

	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "WISH")
	
	if err != nil {
		t.Error(err)
		return
	}

	got := favorite
	exp := s.favorites[0]

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_FavoritePayment_notFound(t *testing.T) {
	s := newTestService();
	_, _, err := s.addAccount(defaultTestAccount)

	paymID := "AXAX"
	_, err = s.FavoritePayment(paymID, "WISH")
	
	if err == nil {
		t.Error(err)
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("AAA ERROR: %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService();
	_, payments, err := s.addAccount(defaultTestAccount)

	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "WISH")
	
	if err != nil {
		t.Error(err)
		return
	}

	favPay, err := s.PayFromFavorite(favorite.ID)

	if err != nil {
		t.Error(err)
		return
	}

	lastPayIndex := len(payments) - 1
	lastPay := payments[lastPayIndex]

	got := favPay.Amount
	exp := lastPay.Amount

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("invalid result, expected: %v, actual: %v", got, exp)
	}
}

func TestService_PayFromFavorite_notFound(t *testing.T) {
	s := newTestService();
	_, _, err := s.addAccount(defaultTestAccount)

	paymID := "AXAX"
	_, err = s.FavoritePayment(paymID, "WISH")
	
	if err == nil {
		t.Error(err)
		return
	}

	_, err = s.PayFromFavorite(paymID)

	if err == nil {
		t.Error(err)
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("AAA ERROR: %v", err)
		return
	}
}
