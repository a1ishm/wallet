package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/a1ishm/wallet/pkg/types"
)

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	acc, err := svc.FindAccountByID(1)

	if err != nil {
		fmt.Println(err)
		return
	}

	acc.Balance += 5_000_00

	payment, errr := svc.Pay(1, 1_000_00, "A")

	if errr != nil {
		fmt.Println(errr)
		return
	}

	paymentID := payment.ID 

	errrr := svc.Reject(paymentID)

	if errrr != nil {
		fmt.Println(errrr)
		return
	}

	result := acc.Balance

	expected := types.Money(5_000_00)

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}

func TestService_Reject_notFound(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	acc, err := svc.FindAccountByID(1)

	if err != nil {
		fmt.Println(err)
		return
	}

	acc.Balance += 5_000_00

	svc.Pay(1, 1_000_00, "A")

	paymentID := "AXAX" 

	errrr := svc.Reject(paymentID)

	if errrr != nil {
		fmt.Println(errrr)
		return
	}

	result := acc.Balance

	expected := ErrPaymentNotFound

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, result)
	}
}
