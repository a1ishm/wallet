package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/a1ishm/wallet/pkg/types"
)

func TestService_FindAccountByID(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	got, err := svc.FindAccountByID(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	exp := account

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	got, _ := svc.FindAccountByID(999)

	var exp *types.Account = nil

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_Reject(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = svc.Deposit(account.ID, 1000)
	payment, _ := svc.Pay(account.ID, 1000, "food")

	err = svc.Reject(payment.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	got := account.Balance
	var exp types.Money = 1000

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}

func TestService_Reject_notFound(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = svc.Deposit(account.ID, 1000)
	svc.Pay(account.ID, 1000, "food")

	_ = svc.Reject("999")
	
	got := account.Balance
	var exp types.Money = 0

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}
