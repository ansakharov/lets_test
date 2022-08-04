package fake_order

import (
	"context"

	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	orders map[uint64]*order.Order
	currID uint64
}

// New instance of repository.
func New() *Repository {
	return &Repository{
		orders: make(map[uint64]*order.Order),
		currID: 1,
	}
}

// Save new order to DB.
func (r *Repository) Save(ctx context.Context, log logrus.FieldLogger, order *order.Order) error {
	order.ID = r.currID
	for idx, item := range order.Items {
		item.OrderID = r.currID
		order.Items[idx] = item
	}
	r.orders[r.currID] = order
	r.currID++

	return nil
}

// Get returns map of orders.
func (r *Repository) Get(ctx context.Context, log logrus.FieldLogger, IDs []uint64) (map[uint64]order.Order, error) {
	result := make(map[uint64]order.Order)

	for _, ID := range IDs {
		order, ok := r.orders[ID]
		if ok {
			result[ID] = *order
		}
	}

	return result, nil
}
