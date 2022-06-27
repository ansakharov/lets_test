package order_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	order_ucase "github.com/ansakharov/lets_test/internal/app/usecase/order"
	"github.com/sirupsen/logrus"
)

// Requst validation errors.
var ErrEmptyOrderIDs = errors.New("no order ids passed")

// Handler creates orders
type Handler struct {
	uCase *order_ucase.Usecase
	log   logrus.FieldLogger
}

// New gives Handler.
func New(
	uCase *order_ucase.Usecase,
	log logrus.FieldLogger,
) *Handler {
	return &Handler{
		uCase: uCase,
		log:   log,
	}
}

// GetOrdersIn is dto for http req.
type GetOrdersIn struct {
	IDs []uint64 `json:"ids"`
}

// validates request.
func (h Handler) validateReq(in *GetOrdersIn) error {
	if len(in.IDs) == 0 {
		return ErrEmptyOrderIDs
	}

	return nil
}

// Create responsible for saving new order.
func (h Handler) Get(ctx context.Context) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// prepare dto to parse request
		in := &GetOrdersIn{}
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

		orders, err := h.uCase.Get(ctx, h.log, in.IDs)
		if err != nil {
			h.log.Errorf("can't get orders: %s", err.Error())
			http.Error(w, "can't get orders: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
	return http.HandlerFunc(fn)
}
