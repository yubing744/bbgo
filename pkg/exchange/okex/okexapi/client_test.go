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

	req.After(1672574086000)
	req.Before(1674647686000)

	candles, err := req.Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, candles)
}

func TestQueryHistoryCandlesticks(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)

	req := client.MarketDataService.NewHistoryCandlesticksRequest("ETH-USTD")
	req.Bar(string(types.Interval1m))

	//req.After(1674561167)
	//req.Before(1674647686)

	candles, err := req.Do(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, candles)
}
