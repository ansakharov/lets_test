package order_handler_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	create_order_handler "github.com/ansakharov/lets_test/handler/create_order"
	get_orders_handler "github.com/ansakharov/lets_test/handler/get_orders"
	order_ucase "github.com/ansakharov/lets_test/internal/app/usecase/order"
	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	fake_order "github.com/ansakharov/lets_test/internal/pkg/repository/order/fake_order_repo"
	mock_order "github.com/ansakharov/lets_test/internal/pkg/repository/order/mocks"
	"github.com/ansakharov/lets_test/logger"
	"github.com/ansakharov/lets_test/metrics"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateOrders(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)

	toSave := order.Order{
		Status:      1,
		UserID:      1,
		PaymentType: 1,
		Items: []order.Item{
			{ID: 2, Amount: 10000, DiscountedAmount: 100},
			{ID: 2, Amount: 2, DiscountedAmount: 3},
		},
	}
	repo.EXPECT().Save(ctx, log, &toSave).Return(nil).Times(1)

	uCase := order_ucase.New(repo)
	h := create_order_handler.New(uCase, log)

	serverFunc := h.Create(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewBuffer([]byte(
			[]byte(`
			{
				"user_id": 1,
				"payment_type": "card",
				"items": [
					{
						"id": 2,
						"amount": 10000,
						"discount": 100
					},
					{
						"id": 2,
						"amount": 2,
						"discount": 3
					}
				]
			}
			`),
		)),
	)
	// if you need it
	req.Header.Set("Content-Type", "Application/Json")

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected := "{\"success\":\"ok\"}\n"

	require.Equal(t, expected, string(data))
}

func TestCreateOrderBadJSON(t *testing.T) {
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)
	uCase := order_ucase.New(repo)
	h := create_order_handler.New(uCase, log)

	serverFunc := h.Create(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewBuffer([]byte(
			[]byte(`
			{
				"payment_type": "card",
				"items": [
					{
						"id": 2,
						"amount": 10000,
						"discount": 100
					},
					{
						"id": 2,
						"amount": 2,
						"discount": 3
					}
				]
			`),
		)),
	)

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "bad json: unexpected EOF\n", string(data))
}

func TestCreateOrderBadReq(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)
	uCase := order_ucase.New(repo)
	h := create_order_handler.New(uCase, log)

	serverFunc := h.Create(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewBuffer([]byte(
			[]byte(`
			{
				"payment_type": "card",
				"items": [
					{
						"id": 2,
						"amount": 10000,
						"discount": 100
					},
					{
						"id": 2,
						"amount": 2,
						"discount": 3
					}
				]
			}
			`),
		)),
	)

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "bad request: invalid user ID\n", string(data))
}

func TestCreateOrderUcaseError(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := mock_order.NewMockOrderRepo(ctl)

	repoErr := errors.New("can't get orders DB is down")
	toSave := order.Order{
		Status:      1,
		UserID:      1,
		PaymentType: 1,
		Items: []order.Item{
			{ID: 2, Amount: 10000, DiscountedAmount: 100},
			{ID: 2, Amount: 2, DiscountedAmount: 3},
		},
	}
	repo.EXPECT().Save(ctx, log, &toSave).Return(repoErr).Times(1)

	uCase := order_ucase.New(repo)
	h := create_order_handler.New(uCase, log)

	serverFunc := h.Create(ctx).ServeHTTP

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewBuffer([]byte(
			[]byte(`
			{
				"user_id": 1,
				"payment_type": "card",
				"items": [
					{
						"id": 2,
						"amount": 10000,
						"discount": 100
					},
					{
						"id": 2,
						"amount": 2,
						"discount": 3
					}
				]
			}
			`),
		)),
	)
	// if you need it
	req.Header.Set("Content-Type", "Application/Json")

	serverFunc(rec, req)

	res := rec.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected := "can't create order: can't get orders DB is down\n"

	require.Equal(t, expected, string(data))
}

func TestCreateAndGetOrder(t *testing.T) {
	metrics.Init()
	log := logger.New()
	ctx := context.Background()

	repo := fake_order.New()

	uCase := order_ucase.New(repo)
	hSave := create_order_handler.New(uCase, log)
	hGet := get_orders_handler.New(uCase, log)

	getFunc := hGet.Get(ctx).ServeHTTP

	// no orders
	recGet := httptest.NewRecorder()
	reqGet := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(`{"ids": [1]}`),
		)),
	)
	getFunc(recGet, reqGet)
	res := recGet.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected := "[]\n"
	require.Equal(t, expected, string(data))

	// save one
	saveFunc := hSave.Create(ctx).ServeHTTP
	saveRec := httptest.NewRecorder()
	saveReq := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewBuffer([]byte(
			[]byte(`
			{
				"user_id": 1,
				"payment_type": "card",
				"items": [
					{
						"id": 2,
						"amount": 10000,
						"discount": 100
					},
					{
						"id": 2,
						"amount": 2,
						"discount": 3
					}
				]
			}
			`),
		)),
	)
	saveFunc(saveRec, saveReq)
	saveRes := saveRec.Result()
	defer saveRes.Body.Close()

	saveData, err := ioutil.ReadAll(saveRes.Body)
	require.NoError(t, err)

	expected = "{\"success\":\"ok\"}\n"
	require.Equal(t, expected, string(saveData))

	// now there is order in db
	recGet = httptest.NewRecorder()
	reqGet = httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBuffer([]byte(
			[]byte(`{"ids": [1]}`),
		)),
	)
	getFunc(recGet, reqGet)
	res = recGet.Result()
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	expected = `[{"ID":1,"Status":1,"UserID":1,"PaymentType":1,"OriginalAmount":10002,"DiscountedAmount":103,"Items":[{"OrderID":1,"ID":2,"Amount":10000,"DiscountedAmount":100},{"OrderID":1,"ID":2,"Amount":2,"DiscountedAmount":3}]}]
`
	require.Equal(t, expected, string(data))
}
