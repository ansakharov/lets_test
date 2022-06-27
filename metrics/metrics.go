package metrics

import (
	"fmt"

	metrics "github.com/rcrowley/go-metrics"
)

const (
	GetOrdersSuccess = "get_orders.ok"
	GetOrdersError   = "get_orders.error"
	GetOrdersCount   = "get_orders.count"

	SaveOrderSuccess = "save_order.ok"
	SaveOrderError   = "save_order.error"
	SaveOrderCount   = "save_order.count"
)

func Init() {
	metrics.Unregister(GetOrdersError)
	metrics.MustRegister(GetOrdersError, metrics.NewCounter())
	metrics.Unregister(GetOrdersSuccess)
	metrics.MustRegister(GetOrdersSuccess, metrics.NewCounter())
	metrics.Unregister(GetOrdersCount)
	metrics.MustRegister(GetOrdersCount, metrics.NewCounter())

	metrics.Unregister(SaveOrderCount)
	metrics.MustRegister(SaveOrderCount, metrics.NewCounter())
	metrics.Unregister(SaveOrderError)
	metrics.MustRegister(SaveOrderError, metrics.NewCounter())
	metrics.Unregister(SaveOrderSuccess)
	metrics.MustRegister(SaveOrderSuccess, metrics.NewCounter())
}

func IncCounter(name string) {
	counter := metrics.Get(name).(metrics.Counter)
	counter.Inc(1)

	fmt.Printf("counter: %s, count: %d\n", name, counter.Count())

}
