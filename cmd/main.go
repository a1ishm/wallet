package main

import (
	"github.com/a1ishm/wallet/pkg/types"
	"log"
	"os"
	"strconv"
)



func main() {
	accounts := []*types.Account{
		{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
		{ID: 2, Phone: "+992000000002", Balance: 20_000_00},
		{ID: 3, Phone: "+992000000003", Balance: 30_000_00},
	}

	path := "files/accs.txt"

	file, _ := os.Create(path)
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Print(cerr)
		}
	}()
	
	var content string

	for i, account := range accounts {
		id := strconv.FormatInt(account.ID, 10)
		phone := string(account.Phone)
		balance := strconv.FormatInt(int64(account.Balance), 10)

		var acc string

		if i == (len(accounts) - 1) {
			acc = id + ";" + phone + ";" + balance 
		} else {
			acc = id + ";" + phone + ";" + balance + "|"
		}
		
		content += acc
	}
	
	_, err := file.Write([]byte(content))

	if err != nil {
		log.Print(err)
		return
	}
}