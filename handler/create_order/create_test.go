package order_handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	h := Handler{}
	in := &OrderIn{
		UserID:      1,
		PaymentType: "card",
		Items: []Item{
			{ID: 1, Amount: 10},
		},
	}
	err := h.validateReq(in)
	require.NoError(t, err)
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name   string
		in     *OrderIn
		expErr error
	}{
		{
			name:   "bad_user_id",
			in:     &OrderIn{UserID: 0},
			expErr: ErrInvalidUserID,
		},
		{
			name:   "bad_payment_type",
			in:     &OrderIn{UserID: 1, PaymentType: "bad"},
			expErr: ErrInvalidPaymentType,
		},
		{
			name:   "no_items",
			in:     &OrderIn{UserID: 1, PaymentType: "card"},
			expErr: ErrEmptyItems,
		},
		{
			name: "bad_item_id",
			in: &OrderIn{
				UserID:      1,
				PaymentType: "card",
				Items: []Item{
					{ID: 0},
				},
			},
			expErr: ErrInvalidItemID,
		},
		{
			name: "bad_item_amount",
			in: &OrderIn{
				UserID:      1,
				PaymentType: "card",
				Items: []Item{
					{ID: 1, Amount: 0},
				},
			},
			expErr: ErrInvalidAmount,
		},
	}
	h := Handler{}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := h.validateReq(tCase.in)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
}
