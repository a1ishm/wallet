package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/a1ishm/wallet/pkg/types"
)



func main() {
	// accounts := []*types.Account{
	// 	{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
	// 	{ID: 2, Phone: "+992000000002", Balance: 20_000_00},
	// 	{ID: 3, Phone: "+992000000003", Balance: 30_000_00},
	// 	{ID: 4, Phone: "+992000000004", Balance: 40_000_00},
    // 	{ID: 5, Phone: "+992000000005", Balance: 50_000_00},
	// }

	arr := []*types.Account{}
	path := "files/accs.txt"

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Print(cerr)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			return 
		}

		content = append(content, buf[:read]...)
	}

	data := string(content)
	accProps := strings.Split(data, "|")

	for _, acc := range accProps {
		var account types.Account
		props := strings.Split(acc, ";")

		id, err := strconv.ParseInt(props[0], 10, 64)
		if err != nil {
			return 
		}

		phone := types.Phone(props[1])

		balance, err := strconv.ParseInt(props[0], 10, 64)
		if err != nil {
			return 
		}

		account.ID = id
		account.Phone = phone
		account.Balance = types.Money(balance)

		arr = append(arr, &account)
	}

	for _, acc := range arr {
		log.Print(*acc)
	}
}