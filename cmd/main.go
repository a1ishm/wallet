package main

import (
	// "io"
	"log"
	"math"
	// "sync"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "github.com/a1ishm/wallet/pkg/types"
)

func main() {
	// goroutines := 3

	// ps := []*types.Payment{
	// 	{ID: "A", AccountID: 1, Amount: 10_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	{ID: "B", AccountID: 2, Amount: 20_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	{ID: "C", AccountID: 3, Amount: 30_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 	{ID: "D", AccountID: 4, Amount: 40_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 	{ID: "E", AccountID: 5, Amount: 50_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	{ID: "F", AccountID: 6, Amount: 60_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// }

	// wg := sync.WaitGroup{}

	// mu := sync.Mutex{}
	// sum := int64(0)

	// if goroutines > len(ps) {
	// 	goroutines = len(ps)
	// }

	// div := []int{}

	// add := len(ps) / goroutines

	// floatAdd := float64(len(ps)) / float64(goroutines)

	// for i := 0; i < goroutines; i++ {
	// 	if i == goroutines-1 {
	// 		if int(math.Round(floatAdd)) == add {
	// 			div = append(div, add+(len(ps)%goroutines))
	// 		} else {
	// 			div = append(div, add)
	// 		}
	// 		break
	// 	}

	// 	if int(math.Round(floatAdd)) == add {
	// 		div = append(div, add)
	// 	} else {
	// 		div = append(div, add+1)
	// 	}
	// }
	// log.Print(div)
	// // lastAdd := len(s.payments) % goroutines

	// start := 0
	// end := 0

	// for i := 0; i < goroutines; i++ {
	// 	wg.Add(1)
	// 	go func(iter int) {
	// 		defer wg.Done()

	// 		end += div[iter]
	// 		payments := append([]*types.Payment{}, ps[start:end]...)
	// 		log.Printf("start: %v; end: %v", start, end)
	// 		start += div[iter]

	// 		val := int64(0)
	// 		for _, payment := range payments {
	// 			val += int64(payment.Amount)
	// 		}
	// 		mu.Lock()
	// 		defer mu.Unlock()
	// 		sum += val
	// 	}(i)
	// 	wg.Wait()
	// }

	// log.Print(types.Money(sum))

	// --------------------------------------------------------------------------------------

	pss := []int{4, 7, 9, 13, 27, 33, 49}
	grts := 4

	for _, ps := range pss {
		divSum := 0
	
		arr := []int{}
	
		floatAdd := float64(ps) / float64(grts)
	
		for i := 0; i < grts; i++ {
			if i == grts-1 {
				arr = append(arr, ps-divSum)
				break
			}
	
			arr = append(arr, int(math.Ceil(floatAdd)))
			divSum += arr[i]
		}

		want := 0
		for _, a := range arr {
			want += a
		}
	
		log.Printf("%v + %v + %v + %v = %v; want: %v", arr[0], arr[1], arr[2], arr[3], ps, want)	
	}

	// 	payments := []*types.Payment{
	// 		{ID: "A", AccountID: 1, Amount: 10_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 		{ID: "B", AccountID: 2, Amount: 20_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 		{ID: "C", AccountID: 3, Amount: 30_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 		{ID: "D", AccountID: 4, Amount: 40_000_00, Category: "food", Status: types.PaymentStatusOk},
	// 		{ID: "E", AccountID: 5, Amount: 50_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 		{ID: "F", AccountID: 6, Amount: 60_000_00, Category: "auto", Status: types.PaymentStatusOk},
	// 	}

	// 	goroutines := 2

	// 	wg := sync.WaitGroup{}

	// 	mu := sync.Mutex{}
	// 	sum := ""

	// 	add := len(payments) / goroutines
	// 	lastAdd := len(payments) % goroutines

	// 	start := 0
	// 	end := add

	// 	for i := 0; i < goroutines; i++ {
	// 		wg.Add(1)

	// 		go func(iter int) {
	// 			defer wg.Done()

	// 			payments := append([]*types.Payment{}, payments[start:end]...)

	// 			if iter == goroutines-1 {
	// 				start += lastAdd
	// 				end += lastAdd
	// 			} else {
	// 				start += add
	// 				end += add
	// 			}

	// 			val := ""
	// 			for _, payment := range payments {
	// 				val += payment.ID
	// 			}
	// 			mu.Lock()
	// 			defer mu.Unlock()
	// 			sum += val
	// 		}(i)

	// 		wg.Wait()
	// 	}

	// 	log.Print(sum)

	// ----------------------------------------------------------------------------------------------------------------------------

	// a := []*types.Payment{
	// 	{ID: "A", AccountID: 1, Amount: 10_000_00, Category: "auto", Status: types.PaymentStatusInProgress},
	// 	{ID: "B", AccountID: 2, Amount: 20_000_00, Category: "auto", Status: types.PaymentStatusInProgress},
	// 	{ID: "C", AccountID: 3, Amount: 30_000_00, Category: "auto", Status: types.PaymentStatusInProgress},
	// 	{ID: "D", AccountID: 4, Amount: 40_000_00, Category: "auto", Status: types.PaymentStatusInProgress},
	// }

	// b := append([]*types.Payment{}, a[1:4]...)

	// log.Print(a)
	// log.Print(b)

	// accounts := []*types.Account{
	// 	{ID: 1, Phone: "+992000000001", Balance: 10_000_00},
	// 	{ID: 2, Phone: "+992000000002", Balance: 20_000_00},
	// 	{ID: 3, Phone: "+992000000003", Balance: 30_000_00},
	// 	{ID: 4, Phone: "+992000000004", Balance: 40_000_00},
	// 	{ID: 5, Phone: "+992000000005", Balance: 50_000_00},
	// }

	// path := "files"
	// abs, err := filepath.Abs(path)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Print(abs)

	// _, err := os.Open(path)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		log.Print("file doesn't exist")
	// 	}
	// } else {
	// 	log.Print("error, not expected")
	// }
}
