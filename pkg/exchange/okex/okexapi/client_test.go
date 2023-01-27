package okexapi

import (
	"context"
	"testing"

	"github.com/c9s/bbgo/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
}

func TestQueryCandlesticks(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)

	req := client.MarketDataService.NewCandlesticksRequest("ETH-USDT")
	req.Bar(string(types.Interval1m))

	candles, err := req.Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, candles)
}

func TestQueryHistoryCandlesticks(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)

	req := client.MarketDataService.NewHistoryCandlesticksRequest("ETH-USDT")
	req.Bar(string(types.Interval1m))

	req.Before(1669852798) // 2022-12-01 07:59:58
	req.After(1701475202)  // 2023-12-02 08:00:02

	candles, err := req.Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, candles)
}
