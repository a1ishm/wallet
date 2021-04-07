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

	got, err := svc.FindAccountByID(999)
	if err != nil {
		fmt.Println(err)
		return
	}

	var exp *types.Account = nil

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", exp, got)
	}
}
