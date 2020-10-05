package main

import (
	"fmt"
	"github.com/a1ishm/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account)
}