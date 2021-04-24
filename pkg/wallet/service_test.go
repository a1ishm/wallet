package wallet

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	// "strconv"
	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
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

	got := favorite.Amount
	exp := payments[0].Amount

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_PayFromFavorite(t *testing.T) {
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

	payment, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	got := payment.Amount
	exp := payments[0].Amount

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_PayFromFavorite_notFound(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FavoritePayment(payments[0].ID, "fav")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.PayFromFavorite(uuid.New().String())
	if err == nil {
		t.Errorf("PayFromFavorite(): must return error, returned nil")
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound, returned %v", err)
		return
	}
}

func TestExport_all(t *testing.T) {
	s := newTestService()
	as := []*types.Account{
		{ID: 1, Phone: "+992100000001", Balance: 11_111_10},
		{ID: 2, Phone: "+992100000011", Balance: 11_111_00},
		{ID: 3, Phone: "+992100000111", Balance: 11_110_00},
		{ID: 4, Phone: "+992100001111", Balance: 11_100_00},
		{ID: 5, Phone: "+992100011111", Balance: 11_000_00},
	}

	ps := []*types.Payment{
		{ID: "aaa", AccountID: 1, Amount: 22_000_00, Category: "auto", Status: types.PaymentStatusOk},
		{ID: "bbb", AccountID: 1, Amount: 22_200_00, Category: "food", Status: types.PaymentStatusOk},
		{ID: "ccc", AccountID: 1, Amount: 22_220_00, Category: "food", Status: types.PaymentStatusOk},
		{ID: "ddd", AccountID: 4, Amount: 22_222_00, Category: "auto", Status: types.PaymentStatusOk},
		{ID: "eee", AccountID: 5, Amount: 22_222_20, Category: "auto", Status: types.PaymentStatusOk},
	}

	fvs := []*types.Favorite{
		{ID: "fff", AccountID: 1, Name: "Fav0", Amount: 30_000_00, Category: "auto"},
		{ID: "ggg", AccountID: 1, Name: "Fav1", Amount: 33_000_00, Category: "food"},
		{ID: "hhh", AccountID: 1, Name: "Fav2", Amount: 33_300_00, Category: "food"},
	}

	s.accounts = append(s.accounts, as...)
	s.payments = append(s.payments, ps...)
	s.favorites = append(s.favorites, fvs...)

	err := s.Export("../../files")
	if err != nil {
		t.Error(err)
	}
}

// func TestImportExport(t *testing.T) {
// 	s := newTestService()

// 	as := []*types.Account{
// 		{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
// 		{ID: 4, Phone: "+992000000002", Balance: 40_000_00},
// 		{ID: 5, Phone: "+992000000003", Balance: 50_000_00},
// 		{ID: 6, Phone: "+992000000004", Balance: 60_000_00},
// 		{ID: 7, Phone: "+992000000005", Balance: 70_000_00},
// 	}

// 	ps := []*types.Payment{
// 		{ID: "aaa", AccountID: 1, Amount: 11_000_00, Category: "auto", Status: types.PaymentStatusOk},
// 		{ID: "bbb", AccountID: 1, Amount: 22_000_00, Category: "food", Status: types.PaymentStatusOk},
// 		{ID: "ccc", AccountID: 1, Amount: 33_000_00, Category: "food", Status: types.PaymentStatusOk},
// 		{ID: "oooooooo", AccountID: 4, Amount: 44_000_00, Category: "auto", Status: types.PaymentStatusOk},
// 		{ID: "eee", AccountID: 5, Amount: 55_000_00, Category: "auto", Status: types.PaymentStatusOk},
// 	}

// 	s.accounts = append(s.accounts, as...)
// 	s.payments = append(s.payments, ps...)

// 	err := s.Import("../../files/")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = s.Export("../../txts/")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestImportExport(t *testing.T) {
	s := newTestService()

	// as := []*types.Account{
	// 	{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
	// 	{ID: 4, Phone: "+992000000002", Balance: 40_000_00},
	// 	{ID: 5, Phone: "+992000000003", Balance: 50_000_00},
	// 	{ID: 6, Phone: "+992000000004", Balance: 60_000_00},
	// 	{ID: 7, Phone: "+992000000005", Balance: 70_000_00},
	// }

	// ps := []*types.Payment{
	// 	{ID: "aaa", AccountID: 1, Amount: 11_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	{ID: "bbb", AccountID: 1, Amount: 22_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 	{ID: "ccc", AccountID: 1, Amount: 33_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 	{ID: "oooooooo", AccountID: 4, Amount: 44_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	{ID: "eee", AccountID: 5, Amount: 55_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// }

	// s.accounts = append(s.accounts, as...)
	// s.payments = append(s.payments, ps...)

	err := s.Import("../../files/")
	if err != nil {
		t.Error(err)
	}

	err = s.Export("../../txts/")
	if err != nil {
		t.Error(err)
	}
}