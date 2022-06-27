package order

// Order represents clients order.
type Order struct {
	ID               uint64
	Status           Status
	UserID           uint64
	PaymentType      PaymentType
	OriginalAmount   uint64
	DiscountedAmount uint64
	Items            []Item
}

// Order status.
type Status uint8

const (
	UnknownStatus Status = iota
	CreatedStatus
	ProcessedStatus
	CanceledStatus
)

// Way of payment
type PaymentType uint8

const (
	UnknownType PaymentType = iota
	Card
	Wallet
)

type Item struct {
	OrderID          uint64 `db:"order_id"`
	ID               uint64 `db:"item_id"`
	Amount           uint64 `db:"amount"`
	DiscountedAmount uint64 `db:"discounted_amount"`
}
