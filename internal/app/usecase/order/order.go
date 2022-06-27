package create_order

import (
	"context"
	"fmt"

	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	orderRepo "github.com/ansakharov/lets_test/internal/pkg/repository/order"
	"github.com/ansakharov/lets_test/metrics"
	"github.com/sirupsen/logrus"
)

// Usecase responsible for saving request.
type Usecase struct {
	repo orderRepo.OrderRepo
}

// New gives Usecase.
func New(orderRepo orderRepo.OrderRepo) *Usecase {
	return &Usecase{repo: orderRepo}
}

// Save single order.
func (uc *Usecase) Save(ctx context.Context, log logrus.FieldLogger, order *order.Order) error {
	if err := uc.repo.Save(ctx, log, order); err != nil {
		metrics.IncCounter(metrics.SaveOrderError)
		metrics.IncCounter(metrics.SaveOrderCount)

		return err
	}

	metrics.IncCounter(metrics.SaveOrderSuccess)
	metrics.IncCounter(metrics.SaveOrderCount)

	return nil
}

// Get orders by ids.
func (uc *Usecase) Get(ctx context.Context, log logrus.FieldLogger, IDs []uint64) ([]order.Order, error) {
	ordersMap, err := uc.repo.Get(ctx, log, IDs)
	if err != nil {
		metrics.IncCounter(metrics.GetOrdersError)
		metrics.IncCounter(metrics.GetOrdersCount)
		return nil, fmt.Errorf("err from orders_repository: %s", err.Error())
	}

	// count amount and discount for all orders.
	for idx, singleOrder := range ordersMap {
		for _, singleService := range singleOrder.Items {
			// TODO fix here bug hehehe
			singleOrder.OriginalAmount += singleService.Amount
			singleOrder.DiscountedAmount += singleService.DiscountedAmount
			ordersMap[idx] = singleOrder
		}
	}
	result := make([]order.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		result = append(result, order)
	}

	metrics.IncCounter(metrics.GetOrdersSuccess)
	metrics.IncCounter(metrics.GetOrdersCount)
	return result, nil
}
