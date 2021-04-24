package main

import (
	// "io"
	"log"
	// "os"
	"path/filepath"
	// "strconv"
	// "strings"
	// "github.com/a1ishm/wallet/pkg/types"
)

func main() {
	// accounts := []*types.Account{
	// 	{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
	// 	{ID: 2, Phone: "+992000000002", Balance: 20_000_00},
	// 	{ID: 3, Phone: "+992000000003", Balance: 30_000_00},
	// 	{ID: 4, Phone: "+992000000004", Balance: 40_000_00},
	// 	{ID: 5, Phone: "+992000000005", Balance: 50_000_00},
	// }

	// arr := []*types.Account{}
	path := "../../files"
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Print(err)
	}
	log.Print(abs)

	// _, err := os.Open(path)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		log.Print("file doesn't exist")
	// 	}
	// } else {
	// 	log.Print("error, not expected")
	// }
}
