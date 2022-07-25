package order_handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	h := Handler{}
	in := &GetOrdersIn{
		IDs: []uint64{3, 2, 1},
	}
	err := h.validateReq(in)
	require.NoError(t, err)
}

func TestValidateError(t *testing.T) {
	h := Handler{}
	in := &GetOrdersIn{
		IDs: nil,
	}
	err := h.validateReq(in)
	require.Error(t, err)
}
