package okexapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaceOrderRequest_AttachAlgoOrdsReduceOnly(t *testing.T) {
	req := &PlaceOrderRequest{
		instrumentID: "SUI-USDT",
		tradeMode:    TradeModeCross,
		side:         SideTypeBuy,
		orderType:    OrderTypeMarket,
		size:         "10",
	}

	req.AttachAlgoOrds([]AttachAlgoOrder{{
		Sz:              "10",
		TpTriggerPx:     "1.1",
		TpOrdPx:         "-1",
		TpTriggerPxType: "last",
		SlTriggerPx:     "0.9",
		SlOrdPx:         "-1",
		SlTriggerPxType: "last",
		ReduceOnly:      true,
	}})

	params, err := req.GetParameters()
	require.NoError(t, err)

	attachRaw, ok := params["attachAlgoOrds"]
	require.True(t, ok, "attachAlgoOrds should be serialized into parameters")

	attachSlice, ok := attachRaw.([]AttachAlgoOrder)
	require.True(t, ok, "attachAlgoOrds should be a slice of AttachAlgoOrder")
	require.Len(t, attachSlice, 1)
	require.True(t, attachSlice[0].ReduceOnly, "reduceOnly should be true to avoid opposite positions")
	require.Equal(t, "10", attachSlice[0].Sz)
}
