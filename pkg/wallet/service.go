package wallet

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/a1ishm/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
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
	/* accounts export */

	if len(s.accounts) == 0 {
		return nil
	}

	file, err := os.Create((dir + "/accounts.dump"))
	if err != nil {
		return err
	}

	var data string

	for i, account := range s.accounts {
		id := strconv.FormatInt(account.ID, 10)
		phone := string(account.Phone)
		balance := strconv.FormatInt(int64(account.Balance), 10)

		var acc string

		if i == (len(s.accounts) - 1) {
			acc = id + ";" + phone + ";" + balance
		} else {
			acc = id + ";" + phone + ";" + balance + "\n"
		}

		data += acc
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	data = ""

	/* payments export */

	if len(s.payments) != 0 {
		file, err = os.Create((dir + "/payments.dump"))
		if err != nil {
			return err
		}

		for i, payment := range s.payments {
			id := payment.ID
			accountID := strconv.FormatInt(payment.AccountID, 10)
			amount := strconv.FormatInt(int64(payment.AccountID), 10)
			category := string(payment.Category)
			status := string(payment.Status)

			var paym string

			if i == (len(s.payments) - 1) {
				paym = id + ";" + accountID + ";" + amount + ";" + category + ";" + status
			} else {
				paym = id + ";" + accountID + ";" + amount + ";" + category + ";" + status + "\n"
			}

			data += paym
		}

		_, err = file.Write([]byte(data))
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
		data = ""
	}

	/* favorites export */

	if len(s.favorites) == 0 {
		return nil
	}

	file, err = os.Create((dir + "/favorites.dump"))
	if err != nil {
		return err
	}

	for i, favorite := range s.favorites {
		id := favorite.ID
		accountID := strconv.FormatInt(favorite.AccountID, 10)
		name := favorite.Name
		amount := strconv.FormatInt(int64(favorite.AccountID), 10)
		category := string(favorite.Category)

		var fav string

		if i == (len(s.favorites) - 1) {
			fav = id + ";" + accountID + ";" + name + ";" + amount + ";" + category
		} else {
			fav = id + ";" + accountID + ";" + name + ";" + amount + ";" + category + "\n"
		}

		data += fav
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func importAccount(dir string, s *Service) error {
	file, err := os.Open((dir + "/accounts.dump"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			var account types.Account
			props := strings.Split(line, ";")

			id, err := strconv.ParseInt(props[0], 10, 64)
			if err != nil {
				return err
			}

			phone := types.Phone(props[1])

			balance, err := strconv.ParseInt(strings.Trim(props[2], "\n"), 10, 64)
			if err != nil {
				return err
			}

			account.ID = id
			account.Phone = phone
			account.Balance = types.Money(balance)

			flag := false
			for i, acc := range s.accounts {
				if acc.ID == id {
					s.accounts[i] = &account
					flag = true
					break
				}
			}
			if !flag {
				s.accounts = append(s.accounts, &account)
				s.nextAccountID = id
				break
			}
		}
		if err != nil {
			return err
		}

		var account types.Account
		props := strings.Split(line, ";")

		id, err := strconv.ParseInt(props[0], 10, 64)
		if err != nil {
			return err
		}

		phone := types.Phone(props[1])

		balance, err := strconv.ParseInt(strings.Trim(props[2], "\n"), 10, 64)
		if err != nil {
			return err
		}

		account.ID = id
		account.Phone = phone
		account.Balance = types.Money(balance)
		flag := false
		for i, acc := range s.accounts {
			if acc.ID == id {
				s.accounts[i] = &account
				flag = true
				break
			}
		}
		if !flag {
			s.accounts = append(s.accounts, &account)
			s.nextAccountID = id
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func importPayment(dir string, s *Service) error {
	file, err := os.Open((dir + "/payments.dump"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			var payment types.Payment
			props := strings.Split(line, ";")

			id := props[0]

			accountID, err := strconv.ParseInt(props[1], 10, 64)
			if err != nil {
				return err
			}

			amount, err := strconv.ParseInt(props[2], 10, 64)
			if err != nil {
				return err
			}

			category := props[3]
			status := props[4]

			payment.ID = id
			payment.AccountID = accountID
			payment.Amount = types.Money(amount)
			payment.Category = types.PaymentCategory(category)
			payment.Status = types.PaymentStatus(status)

			idFound := false
			for i, paym := range s.payments {
				if paym.ID == id {
					s.payments[i] = &payment
					idFound = true
					break
				}
			}
			if !idFound {
				s.payments = append(s.payments, &payment)
				break
			}
		}
		if err != nil {
			return err
		}

		var payment types.Payment
		props := strings.Split(line, ";")

		id := props[0]

		accountID, err := strconv.ParseInt(props[1], 10, 64)
		if err != nil {
			return err
		}

		amount, err := strconv.ParseInt(props[2], 10, 64)
		if err != nil {
			return err
		}

		category := props[3]
		status := props[4]

		payment.ID = id
		payment.AccountID = accountID
		payment.Amount = types.Money(amount)
		payment.Category = types.PaymentCategory(category)
		payment.Status = types.PaymentStatus(status)
		idFound := false
		for i, paym := range s.payments {
			if paym.ID == id {
				s.payments[i] = &payment
				idFound = true
				break
			}
		}
		if !idFound {
			s.payments = append(s.payments, &payment)
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func importFavorite(dir string, s *Service) error {
	file, err := os.Open((dir + "/favorites.dump"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			var favorite types.Favorite
			props := strings.Split(line, ";")

			id := props[0]

			accountID, err := strconv.ParseInt(props[1], 10, 64)
			if err != nil {
				return err
			}

			name := props[2]

			amount, err := strconv.ParseInt(props[3], 10, 64)
			if err != nil {
				return err
			}

			category := props[4]

			favorite.ID = id
			favorite.AccountID = accountID
			favorite.Name = name
			favorite.Amount = types.Money(amount)
			favorite.Category = types.PaymentCategory(category)
			idFound := false
			for i, fav := range s.favorites {
				if fav.ID == id {
					s.favorites[i] = &favorite
					idFound = true
					break
				}
			}
			if !idFound {
				s.favorites = append(s.favorites, &favorite)
				break
			}
		}
		if err != nil {
			return err
		}

		var favorite types.Favorite
		props := strings.Split(line, ";")

		id := props[0]

		accountID, err := strconv.ParseInt(props[1], 10, 64)
		if err != nil {
			return err
		}

		name := props[2]

		amount, err := strconv.ParseInt(props[3], 10, 64)
		if err != nil {
			return err
		}

		category := props[4]

		favorite.ID = id
		favorite.AccountID = accountID
		favorite.Name = name
		favorite.Amount = types.Money(amount)
		favorite.Category = types.PaymentCategory(category)
		idFound := false
		for i, fav := range s.favorites {
			if fav.ID == id {
				s.favorites[i] = &favorite
				idFound = true
				break
			}
		}
		if !idFound {
			s.favorites = append(s.favorites, &favorite)
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Import(dir string) error {
	/* accounts import */
	err := importAccount(dir, s)
	if err != nil {
		return err
	}

	/* payments import */
	err = importPayment(dir, s)
	if err != nil {
		return err
	}

	/* favorites import */
	err = importFavorite(dir, s)
	if err != nil {
		return err
	}

	return nil
}
