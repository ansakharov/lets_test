package order_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	create_order "github.com/ansakharov/lets_test/internal/app/usecase/order"
	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	"github.com/sirupsen/logrus"
)

// Requst validation errors.
var ErrInvalidUserID = errors.New("invalid user ID")
var ErrInvalidAmount = errors.New("invalid price")
var ErrInvalidPaymentType = errors.New("invalid payment type")
var ErrEmptyItems = errors.New("items can't be empty")
var ErrInvalidItemID = errors.New("invalid service id")

// Handler creates orders
type Handler struct {
	uCase *create_order.Usecase
	log   logrus.FieldLogger
}

// New gives Handler.
func New(
	uCase *create_order.Usecase,
	log logrus.FieldLogger,
) *Handler {
	return &Handler{
		uCase: uCase,
		log:   log,
	}
}

// OrderIn is dto for http req.
type OrderIn struct {
	UserID      uint64 `json:"user_id"` // 0
	PaymentType string `json:"payment_type"`
	Items       []Item `json:"items"`
}

type Item struct {
	ID       uint64 `json:"id"`
	Amount   uint64 `json:"amount"`
	Discount uint64 `json:"discount"`
}

// OrderFromDTO creates Order for business layer.
func (in OrderIn) OrderFromDTO() order.Order {
	items := []order.Item{}
	for _, item := range in.Items {
		items = append(items, order.Item{
			ID:               item.ID,
			Amount:           item.Amount,
			DiscountedAmount: item.Discount,
		})
	}
	return order.Order{
		Status:      order.CreatedStatus,
		UserID:      in.UserID,
		PaymentType: order.PaymentType(paymentTypes[in.PaymentType]),
		Items:       items,
	}
}

var paymentTypes = map[string]PaymentType{
	"card":   Card,
	"wallet": Wallet,
}

type PaymentType uint8

const (
	UndefinedType PaymentType = iota
	Card
	Wallet
)

// validates request.
func (h Handler) validateReq(in *OrderIn) error {
	// user ID can't be 0
	if in.UserID == 0 {
		return ErrInvalidUserID
	}
	// payment type must be in paymentTypes
	if _, ok := paymentTypes[in.PaymentType]; !ok {
		return ErrInvalidPaymentType
	}
	// no services passed in request
	if len(in.Items) == 0 {
		return ErrEmptyItems
	}
	// service doesn't contain valid id
	for i := range in.Items {
		if in.Items[i].ID == 0 {
			return ErrInvalidItemID
		}
		if in.Items[i].Amount == 0 {
			return ErrInvalidAmount
		}
	}
	return nil
}

// Create responsible for saving new order.
func (h Handler) Create(ctx context.Context) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// prepare dto to parse request
		in := &OrderIn{}
		// parse req body to dto
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			h.log.Errorf("can't parse req: %s", err.Error())
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}

		// check that request valid
		err = h.validateReq(in)
		if err != nil {
			h.log.Errorf("bad req: %v: %s", in, err.Error())
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		order := in.OrderFromDTO()
		err = h.uCase.Save(ctx, h.log, &order)
		if err != nil {
			h.log.Errorf("can't create order: %v: %s", order, err.Error())
			http.Error(w, "can't create order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		m := make(map[string]interface{})
		m["success"] = "ok"
		json.NewEncoder(w).Encode(m)

	}
	return http.HandlerFunc(fn)
}
