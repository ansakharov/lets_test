package order_handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	get_order_handler "github.com/ansakharov/lets_test/handler/get_orders"
	order_ucase "github.com/ansakharov/lets_test/internal/app/usecase/order"
	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	mock_order "github.com/ansakharov/lets_test/internal/pkg/repository/order/mocks"
	"github.com/ansakharov/lets_test/logger"
	"github.com/ansakharov/lets_test/metrics"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	// pool, err := pgxpool.Connect(ctx, "postgres://alesakharov@localhost:5432/postgres")
	// repo := order_repo.New(pool)
	// require.NoError(t, err)

	reqID := 1

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)

	exp := map[uint64]order.Order{
		1: {
			ID:               uint64(reqID),
			Status:           0,
			UserID:           1,
			PaymentType:      1,
			OriginalAmount:   0,
			DiscountedAmount: 0,
			Items: []order.Item{
				{
					OrderID: uint64(reqID),
					ID:      1, Amount: 100, DiscountedAmount: 0,
				},
			},
		},
	}
	repo.EXPECT().Get(ctx, log, []uint64{uint64(reqID)}).Return(exp, nil).Times(1)

	uCase := order_ucase.New(repo)
	h := get_order_handler.New(uCase, log)

	serverFunc := h.Get(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(fmt.Sprintf(`{"ids": [%d]}`, reqID)),
		)),
	)
	// if you need it
	req.Header.Set("Content-Type", "Application/Json")

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected :=
		`[{"ID":1,"Status":0,"UserID":1,"PaymentType":1,"OriginalAmount":100,"DiscountedAmount":0,"Items":[{"OrderID":1,"ID":1,"Amount":100,"DiscountedAmount":0}]}]` +
			"\n"

	require.Equal(t, expected, string(data))
}

func TestGetOrdersBadJSON(t *testing.T) {
	// {"ids": [999]
	// {"ids": asdfoij}
	// {: asdfoij}
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)
	uCase := order_ucase.New(repo)
	h := get_order_handler.New(uCase, log)

	serverFunc := h.Get(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(`{"ids": [999]`),
		)),
	)

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "bad json: unexpected EOF\n", string(data))
}

func TestGetOrdersBadReq(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)
	uCase := order_ucase.New(repo)
	h := get_order_handler.New(uCase, log)

	serverFunc := h.Get(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(`{"ids": []}`),
		)),
	)

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "bad request: no order ids passed\n", string(data))
}

func TestGetOrdersUcaseError(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	// pool, err := pgxpool.Connect(ctx, "postgres://alesakharov@localhost:5432/postgres")
	// repo := order_repo.New(pool)
	// require.NoError(t, err)

	reqID := 1

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)

	repoErr := errors.New("can't get orders DB is down")

	repo.EXPECT().Get(ctx, log, []uint64{uint64(reqID)}).Return(nil, repoErr).Times(1)

	uCase := order_ucase.New(repo)
	h := get_order_handler.New(uCase, log)

	serverFunc := h.Get(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(fmt.Sprintf(`{"ids": [%d]}`, reqID)),
		)),
	)
	// if you need it
	req.Header.Set("Content-Type", "Application/Json")

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected := "can't get orders: err from orders_repository: can't get orders DB is down\n"

	require.Equal(t, expected, string(data))
}
