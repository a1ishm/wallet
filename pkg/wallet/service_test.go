package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/a1ishm/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")

	expected := &types.Account{
		ID:      3,
		Phone:   "+992000000003",
		Balance: 0,
	}

	result, err := svc.FindAccountByID(3)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
 
	result, err := svc.FindAccountByID(4)

	expected := ErrAccountNotFound

	if err != nil {
		fmt.Println(err)
		return
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}

}
