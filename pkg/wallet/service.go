package wallet

import (
	"bufio"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment(s) not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Error string

func (e Error) Error() string {
	return string(e)
}

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, paym := range s.payments {
		if paym.ID == paymentID {
			payment = paym
			break
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	repeated, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return repeated, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}

	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}

	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			log.Print(cerr)
		}
	}()

	var data string

	for i, account := range s.accounts {
		id := strconv.FormatInt(account.ID, 10)
		phone := string(account.Phone)
		balance := strconv.FormatInt(int64(account.Balance), 10)

		var acc string

		if i == (len(s.accounts) - 1) {
			acc = id + ";" + phone + ";" + balance
		} else {
			acc = id + ";" + phone + ";" + balance + "|"
		}

		data += acc
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
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
			return err
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
			return err
		}

		phone := types.Phone(props[1])

		balance, err := strconv.ParseInt(props[2], 10, 64)
		if err != nil {
			return err
		}

		account.ID = id
		account.Phone = phone
		account.Balance = types.Money(balance)

		s.accounts = append(s.accounts, &account)
	}

	return nil
}

func (s *Service) Export(dir string) error {
	var err error
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	aPath := abs + "/accounts.dump" // если не примет, то убрать слэш в названии файла
	pPath := abs + "/payments.dump"
	fPath := abs + "/favorites.dump"

	var accounts *os.File
	var payments *os.File
	var favorites *os.File

	aExist := true
	pExist := true
	fExist := true

	if len(s.accounts) == 0 {
		aExist = false
	}
	if len(s.payments) == 0 {
		pExist = false
	}
	if len(s.favorites) == 0 {
		fExist = false
	}

	if aExist {
		accounts, err = os.Create(aPath)
		if err != nil {
			return err
		}
	}
	if pExist {
		payments, err = os.Create(pPath)
		if err != nil {
			return err
		}
	}
	if fExist {
		favorites, err = os.Create(fPath)
		if err != nil {
			return err
		}
	}

	aData := ""
	pData := ""
	fData := ""

	for i, account := range s.accounts {
		if !aExist {
			break
		}

		id := strconv.Itoa(int(account.ID))
		phone := string(account.Phone)
		balance := strconv.Itoa(int(account.Balance))

		line := id + ";" + phone + ";" + balance + "\n"
		if i == len(s.accounts)-1 {
			line = id + ";" + phone + ";" + balance
		}

		aData += line
	}

	for i, payment := range s.payments {
		if !pExist {
			break
		}

		id := payment.ID
		accountID := strconv.Itoa(int(payment.AccountID))
		amount := strconv.Itoa(int(payment.Amount))
		category := string(payment.Category)
		status := string(payment.Status)

		line := id + ";" + accountID + ";" + amount + ";" + category + ";" + status + "\n"
		if i == len(s.payments)-1 {
			line = id + ";" + accountID + ";" + amount + ";" + category + ";" + status
		}

		pData += line
	}

	for i, favorite := range s.favorites {
		if !fExist {
			break
		}

		id := favorite.ID
		accountID := strconv.Itoa(int(favorite.AccountID))
		name := favorite.Name
		amount := strconv.Itoa(int(favorite.Amount))
		category := string(favorite.Category)

		line := id + ";" + accountID + ";" + name + ";" + amount + ";" + category + "\n"
		if i == len(s.favorites)-1 {
			line = id + ";" + accountID + ";" + name + ";" + amount + ";" + category
		}

		fData += line
	}

	if aExist {
		_, err = accounts.Write([]byte(aData))
		if err != nil {
			return err
		}
	}
	if pExist {
		_, err = payments.Write([]byte(pData))
		if err != nil {
			return err
		}
	}
	if fExist {
		_, err = favorites.Write([]byte(fData))
		if err != nil {
			return err
		}
	}

	if aExist {
		err = accounts.Close()
		if err != nil {
			return err
		}
	}
	if pExist {
		err = payments.Close()
		if err != nil {
			return err
		}
	}
	if fExist {
		err = favorites.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Import(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	aPath := abs + "/accounts.dump" // если не примет, то убрать слэш в названии файла
	pPath := abs + "/payments.dump"
	fPath := abs + "/favorites.dump"

	aExist := true
	pExist := true
	fExist := true
	eof := false
	accountFound := false
	var reader *bufio.Reader
	var nextAccountID int64

	accounts, err := os.Open(aPath)
	if err != nil {
		if os.IsNotExist(err) {
			aExist = false
		} else {
			return err
		}
	}
	payments, err := os.Open(pPath)
	if err != nil {
		if os.IsNotExist(err) {
			pExist = false
		} else {
			return err
		}
	}
	favorites, err := os.Open(fPath)
	if err != nil {
		if os.IsNotExist(err) {
			fExist = false
		} else {
			return err
		}
	}

	if aExist {
		reader = bufio.NewReader(accounts)
	}
	for {
		if !aExist {
			break
		}

		line, err := reader.ReadString('\n') // есть очень тупая идея на запас*
		if err == io.EOF {
			eof = true
		}
		var account types.Account
		aProps := strings.Split(line, ";")

		id, err := strconv.Atoi(strings.Trim(aProps[0], "\n"))
		if err != nil {
			return err
		}
		phone := strings.Trim(aProps[1], "\n")
		balance, err := strconv.Atoi(strings.Trim(aProps[2], "\n"))
		if err != nil {
			return err
		}

		account.ID = int64(id)
		account.Phone = types.Phone(phone)
		account.Balance = types.Money(balance)

		for i, acc := range s.accounts {
			if acc.ID == account.ID {
				s.accounts[i] = &account
				accountFound = true
				break
			}
		}
		if !accountFound {
			s.accounts = append(s.accounts, &account)
		}
		accountFound = false

		if eof {
			break
		}
	}
	eof = false

	if pExist {
		reader = bufio.NewReader(payments)
	}
	for {
		if !pExist {
			break
		}

		line, err := reader.ReadString('\n') // есть очень тупая идея на запас*
		if err == io.EOF {
			eof = true
		}
		var payment types.Payment
		pProps := strings.Split(line, ";")

		id := strings.Trim(pProps[0], "\n")
		accountID, err := strconv.Atoi(strings.Trim(pProps[1], "\n"))
		if err != nil {
			return err
		}
		amount, err := strconv.Atoi(strings.Trim(pProps[2], "\n"))
		if err != nil {
			return err
		}
		category := strings.Trim(pProps[3], "\n")
		status := strings.Trim(pProps[4], "\n")

		payment.ID = id
		payment.AccountID = int64(accountID)
		payment.Amount = types.Money(amount)
		payment.Category = types.PaymentCategory(category)
		payment.Status = types.PaymentStatus(status)

		for i, paym := range s.payments {
			if paym.ID == payment.ID {
				s.payments[i] = &payment
				accountFound = true
				break
			}
		}
		if !accountFound {
			s.payments = append(s.payments, &payment)
		}
		accountFound = false

		if eof {
			break
		}
	}
	eof = false

	if fExist {
		reader = bufio.NewReader(favorites)
	}
	for {
		if !fExist {
			break
		}

		line, err := reader.ReadString('\n') // есть очень тупая идея на запас*
		if err == io.EOF {
			eof = true
		}
		var favorite types.Favorite
		fProps := strings.Split(line, ";")

		id := strings.Trim(fProps[0], "\n")
		accountID, err := strconv.Atoi(strings.Trim(fProps[1], "\n"))
		if err != nil {
			return err
		}
		name := strings.Trim(fProps[2], "\n")
		amount, err := strconv.Atoi(strings.Trim(fProps[3], "\n"))
		if err != nil {
			return err
		}
		category := strings.Trim(fProps[4], "\n")

		favorite.ID = id
		favorite.AccountID = int64(accountID)
		favorite.Name = name
		favorite.Amount = types.Money(amount)
		favorite.Category = types.PaymentCategory(category)

		for i, fav := range s.favorites {
			if fav.ID == favorite.ID {
				s.favorites[i] = &favorite
				accountFound = true
				break
			}
		}
		if !accountFound {
			s.favorites = append(s.favorites, &favorite)
		}
		accountFound = false

		if eof {
			break
		}
	}

	for _, acc := range s.accounts {
		if acc.ID > nextAccountID {
			nextAccountID = acc.ID
		}
	}
	s.nextAccountID = nextAccountID

	if aExist {
		err = accounts.Close()
		if err != nil {
			return err
		}
	}
	if pExist {
		err = payments.Close()
		if err != nil {
			return err
		}
	}
	if fExist {
		err = favorites.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	var payments []types.Payment
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, *payment)
		}
	}

	return payments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if payments == nil {
		return nil
	}

	if records < 1 {
		return Error("there must be at least 1 record")
	}

	if records > len(payments) {
		records = len(payments)
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	num := (len(payments) / records)
	if len(payments)%records != 0 && num >= 1 {
		num++
	}

	var iter int
	for i := 0; i < num; i++ {
		path := abs + "/payments" + strconv.Itoa(i+1) + ".dump"
		if num <= 1 {
			path = abs + "/payments.dump"
		}

		record, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			cerr := record.Close()
			if cerr != nil {
				log.Print(cerr)
			}
		}()

		if iter == records {
			if len(payments)%records != 0 {
				records = len(payments) % records
			}
		}

		var data string
		for x := 0; x < records; x++ {
			payment := payments[iter]
			iter++

			id := payment.ID
			accountID := strconv.Itoa(int(payment.AccountID))
			amount := strconv.Itoa(int(payment.Amount))
			category := string(payment.Category)
			status := string(payment.Status)

			line := id + ";" + accountID + ";" + amount + ";" + category + ";" + status + "\n"
			if x == records-1 {
				line = id + ";" + accountID + ";" + amount + ";" + category + ";" + status
			}

			data += line
		}

		_, err = record.Write([]byte(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SumPayments(goroutines int) types.Money {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	sum := int64(0)
	start := 0
	end := 0

	if goroutines > len(s.payments) {
		goroutines = len(s.payments)
	}
	if goroutines == 0 {
		goroutines = 1
	}

	ratio := []int{}
	add := float64(len(s.payments)) / float64(goroutines)
	ratioSum := 0

	for i := 0; i < goroutines; i++ {
		if i == goroutines-1 {
			ratio = append(ratio, len(s.payments)-ratioSum)
			break
		}

		ratio = append(ratio, int(math.Ceil(add)))
		ratioSum += ratio[i]
	}

	for i := 0; i < goroutines; i++ {
		if goroutines == 1 {
			payments := s.payments

			val := int64(0)
			for _, payment := range payments {
				val += int64(payment.Amount)
			}
			sum += val
			break
		}

		wg.Add(1)
		go func(iter int) {
			defer wg.Done()

			end += ratio[iter]
			payments := append([]*types.Payment{}, s.payments[start:end]...)
			start += ratio[iter]

			val := int64(0)
			for _, payment := range payments {
				val += int64(payment.Amount)
			}

			mu.Lock()
			defer mu.Unlock()
			sum += val
		}(i)
		wg.Wait()
	}

	return types.Money(sum)
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	filteredPayments := []types.Payment{}
	start := 0
	end := 0
	accountFound := false

	for _, account := range s.accounts {
		if account.ID == accountID {
			accountFound = true
			break
		}
	}
	if !accountFound {
		return nil, ErrAccountNotFound
	}

	if goroutines > len(s.payments) {
		goroutines = len(s.payments)
	}
	if goroutines == 0 {
		goroutines = 1
	}

	ratio := []int{}
	add := float64(len(s.payments)) / float64(goroutines)
	ratioSum := 0

	for i := 0; i < goroutines; i++ {
		if i == goroutines-1 {
			ratio = append(ratio, len(s.payments)-ratioSum)
			break
		}

		ratio = append(ratio, int(math.Ceil(add)))
		ratioSum += ratio[i]
	}

	for i := 0; i < goroutines; i++ {
		if goroutines == 1 {
			payments := s.payments

			var filtered []*types.Payment
			for _, payment := range payments {
				if payment.AccountID == accountID {
					filtered = append(filtered, payment)
				}
			}

			for _, payment := range filtered {
				filteredPayments = append(filteredPayments, *payment)
			}
			break
		}

		wg.Add(1)
		go func(iter int) {
			defer wg.Done()

			end += ratio[iter]
			payments := append([]*types.Payment{}, s.payments[start:end]...)
			start += ratio[iter]

			var filtered []*types.Payment
			for _, payment := range payments {
				if payment.AccountID == accountID {
					filtered = append(filtered, payment)
				}
			}
			
			mu.Lock()
			defer mu.Unlock()
			for _, payment := range filtered {
				filteredPayments = append(filteredPayments, *payment)
			}
		}(i)
		wg.Wait()
	}

	return filteredPayments, nil
}

