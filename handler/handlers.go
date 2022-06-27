package handler

import (
	"context"
	"fmt"

	"github.com/ansakharov/lets_test/cmd/config"
	create_order_handler "github.com/ansakharov/lets_test/handler/create_order"
	echo_handler "github.com/ansakharov/lets_test/handler/echo"
	get_orders_handler "github.com/ansakharov/lets_test/handler/get_orders"
	orderUCase "github.com/ansakharov/lets_test/internal/app/usecase/order"
	orderRepo "github.com/ansakharov/lets_test/internal/pkg/repository/order"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	echoRoute   = "/echo"
	orderRoute  = "/order"
	ordersRoute = "/orders"
)

// Router register necessary routes and returns an instance of a router.
func Router(ctx context.Context, log logrus.FieldLogger, config *config.Config) (*mux.Router, error) {
	r := mux.NewRouter()

	// echo
	r.HandleFunc(echoRoute, echo_handler.Handler("Your message: ").ServeHTTP).Methods("GET")

	pool, err := pgxpool.Connect(context.Background(), config.DbConnString)
	if err != nil {
		return nil, fmt.Errorf("can't create pg pool: %s", err.Error())
	}
	repo := orderRepo.New(pool)
	orderUCase := orderUCase.New(repo)

	createOrderHandleFunc := create_order_handler.New(orderUCase, log).Create(ctx).ServeHTTP
	// create order
	r.HandleFunc(orderRoute, createOrderHandleFunc).Methods("POST")

	getOrderHandlerFunc := get_orders_handler.New(orderUCase, log).Get(ctx).ServeHTTP
	// get orders
	r.HandleFunc(ordersRoute, getOrderHandlerFunc).Methods("GET")

	return r, nil
}
