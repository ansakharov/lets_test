package create_order

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ansakharov/lets_test/internal/pkg/entity/order"
	repoMock "github.com/ansakharov/lets_test/internal/pkg/repository/order/mocks"
	log "github.com/ansakharov/lets_test/logger"
	"github.com/ansakharov/lets_test/metrics"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	metrics.Init()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockOrderRepo(ctl)

	ctx := context.Background()
	log := log.New()
	in := []uint64{1, 2, 3}

	mockResp := map[uint64]order.Order{
		1: {
			ID:          1,
			PaymentType: 1,
			Items: []order.Item{
				{ID: 2, Amount: 100, DiscountedAmount: 10},
				{ID: 3, Amount: 1000, DiscountedAmount: 20},
			},
		},
		2: {
			ID:          2,
			PaymentType: 1,
			Items: []order.Item{
				{ID: 2, Amount: 100, DiscountedAmount: 10},
			},
		},
	}

	expected := []order.Order{
		{
			ID:               1,
			PaymentType:      1,
			OriginalAmount:   1100,
			DiscountedAmount: 30,
			Items: []order.Item{
				{ID: 2, Amount: 100, DiscountedAmount: 10},
				{ID: 3, Amount: 1000, DiscountedAmount: 20},
			},
		},
		{
			ID:               2,
			PaymentType:      1,
			OriginalAmount:   100,
			DiscountedAmount: 10,
			Items: []order.Item{
				{ID: 2, Amount: 100, DiscountedAmount: 10},
			},
		},
	}
	repo.EXPECT().Get(ctx, log, in).Return(mockResp, nil).Times(1)

	Usecase := New(repo)
	orders, err := Usecase.Get(ctx, log, in)
	require.NoError(t, err)
	require.ElementsMatch(t, expected, orders)
}

func TestGetError(t *testing.T) {
	metrics.Init()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockOrderRepo(ctl)

	repoErr := errors.New("db is down")
	ctx := context.Background()
	log := log.New()
	in := []uint64{1, 2, 3}
	repo.EXPECT().Get(ctx, log, in).Return(nil, repoErr).Times(1)

	Usecase := New(repo)
	orders, err := Usecase.Get(ctx, log, in)
	require.Error(t, err)
	require.EqualError(t,
		fmt.Errorf("err from orders_repository: %s", repoErr.Error()),
		err.Error(),
	)
	require.Nil(t, orders)
}

func TestSaveError(t *testing.T) {
	metrics.Init()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockOrderRepo(ctl)

	repoErr := errors.New("db is down")
	ctx := context.Background()
	log := log.New()
	in := &order.Order{}
	repo.EXPECT().Save(ctx, log, in).Return(repoErr).Times(1)

	Usecase := New(repo)
	err := Usecase.Save(ctx, log, in)
	require.Error(t, err)
}

func TestSave(t *testing.T) {
	metrics.Init()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := repoMock.NewMockOrderRepo(ctl)

	ctx := context.Background()
	log := log.New()
	in := &order.Order{}
	repo.EXPECT().Save(ctx, log, in).Return(nil).Times(1)

	Usecase := New(repo)
	err := Usecase.Save(ctx, log, in)
	require.NoError(t, err)
}
