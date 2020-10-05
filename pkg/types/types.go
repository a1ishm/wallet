package types

// Money t
type Money int64

// PaymentCategory t
type PaymentCategory string

// PaymentStatus t
type PaymentStatus string

const (
	// PaymentStatusOk v
	PaymentStatusOk PaymentStatus = "OK"

	// PaymentStatusFail v
	PaymentStatusFail PaymentStatus = "FAIL"

	// PaymentStatusInProgress v
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment t
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Phone t
type Phone string

// Account t
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
