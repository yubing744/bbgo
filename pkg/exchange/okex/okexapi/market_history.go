package okexapi

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/c9s/bbgo/pkg/fixedpoint"
)

type HistoryCandlesticksRequest struct {
	client *RestClient

	instId string `param:"instId"`

	limit *int `param:"limit"`

	bar *string `param:"bar"`

	after *int64 `param:"after,seconds"`

	before *int64 `param:"before,seconds"`
}

func (r *HistoryCandlesticksRequest) After(after int64) *HistoryCandlesticksRequest {
	r.after = &after
	return r
}

func (r *HistoryCandlesticksRequest) Before(before int64) *HistoryCandlesticksRequest {
	r.before = &before
	return r
}

func (r *HistoryCandlesticksRequest) Bar(bar string) *HistoryCandlesticksRequest {
	r.bar = &bar
	return r
}

func (r *HistoryCandlesticksRequest) Limit(limit int) *HistoryCandlesticksRequest {
	r.limit = &limit
	return r
}

func (r *HistoryCandlesticksRequest) InstrumentID(instId string) *HistoryCandlesticksRequest {
	r.instId = instId
	return r
}

func (r *HistoryCandlesticksRequest) Do(ctx context.Context) ([]Candle, error) {
	// SPOT, SWAP, FUTURES, OPTION
	var params = url.Values{}
	params.Add("instId", r.instId)

	if r.bar != nil {
		params.Add("bar", *r.bar)
	}

	if r.before != nil {
		params.Add("before", strconv.FormatInt(*r.before*1000, 10))
	}

	if r.after != nil {
		params.Add("after", strconv.FormatInt(*r.after*1000, 10))
	}

	if r.limit != nil {
		params.Add("limit", strconv.Itoa(*r.limit))
	}

	req, err := r.client.NewRequest(ctx, "GET", "/api/v5/market/history-candles", params, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	type candleEntry [7]string
	var candlesResponse struct {
		Code    string        `json:"code"`
		Message string        `json:"msg"`
		Data    []candleEntry `json:"data"`
	}

	if err := resp.DecodeJSON(&candlesResponse); err != nil {
		return nil, err
	}

	var candles []Candle
	for _, entry := range candlesResponse.Data {
		timestamp, err := strconv.ParseInt(entry[0], 10, 64)
		if err != nil {
			return candles, err
		}

		open, err := fixedpoint.NewFromString(entry[1])
		if err != nil {
			return candles, err
		}

		high, err := fixedpoint.NewFromString(entry[2])
		if err != nil {
			return candles, err
		}

		low, err := fixedpoint.NewFromString(entry[3])
		if err != nil {
			return candles, err
		}

		cls, err := fixedpoint.NewFromString(entry[4])
		if err != nil {
			return candles, err
		}

		vol, err := fixedpoint.NewFromString(entry[5])
		if err != nil {
			return candles, err
		}

		volCcy, err := fixedpoint.NewFromString(entry[6])
		if err != nil {
			return candles, err
		}

		var interval = "1m"
		if r.bar != nil {
			interval = *r.bar
		}

		candles = append(candles, Candle{
			InstrumentID:     r.instId,
			Interval:         interval,
			Time:             time.Unix(0, timestamp*int64(time.Millisecond)),
			Open:             open,
			High:             high,
			Low:              low,
			Close:            cls,
			Volume:           vol,
			VolumeInCurrency: volCcy,
		})
	}

	return candles, nil
}

func (c *RestClient) NewHistoryCandlesticksRequest(instId string) *HistoryCandlesticksRequest {
	return &HistoryCandlesticksRequest{
		client: c,
		instId: instId,
	}
}
