package okexapi

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"

	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
)

func (c *RestClient) NewBatchPlaceOrderRequest() *BatchPlaceOrderRequest {
	return &BatchPlaceOrderRequest{
		client: c,
	}
}

func (c *RestClient) NewBatchCancelOrderRequest() *BatchCancelOrderRequest {
	return &BatchCancelOrderRequest{
		client: c,
	}
}

func (c *RestClient) NewGetOrderDetailsRequest() *GetOrderDetailsRequest {
	return &GetOrderDetailsRequest{
		client: c,
	}
}

func (c *RestClient) NewGetTransactionDetailsRequest() *GetTransactionDetailsRequest {
	return &GetTransactionDetailsRequest{
		client: c,
	}
}

func (c *RestClient) NewClosePositionOrderRequest() *ClosePositionRequest {
	return &ClosePositionRequest{
		client: c,
	}
}

type BatchCancelOrderRequest struct {
	client *RestClient

	reqs []*CancelOrderRequest
}

func (r *BatchCancelOrderRequest) Add(reqs ...*CancelOrderRequest) *BatchCancelOrderRequest {
	r.reqs = append(r.reqs, reqs...)
	return r
}

func (r *BatchCancelOrderRequest) Do(ctx context.Context) ([]OrderResponse, error) {
	var parameterList []map[string]interface{}

	for _, req := range r.reqs {
		params, err := req.GetParameters()
		if err != nil {
			return nil, err
		}
		parameterList = append(parameterList, params)
	}

	req, err := r.client.NewAuthenticatedRequest(ctx, "POST", "/api/v5/trade/cancel-batch-orders", nil, parameterList)
	if err != nil {
		return nil, err
	}

	response, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	var data []OrderResponse
	if err := json.Unmarshal(apiResponse.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type BatchPlaceOrderRequest struct {
	client *RestClient

	reqs []*PlaceOrderRequest
}

func (r *BatchPlaceOrderRequest) Add(reqs ...*PlaceOrderRequest) *BatchPlaceOrderRequest {
	r.reqs = append(r.reqs, reqs...)
	return r
}

func (r *BatchPlaceOrderRequest) Do(ctx context.Context) ([]OrderResponse, error) {
	var parameterList []map[string]interface{}

	for _, req := range r.reqs {
		params, _ := req.GetParameters()
		parameterList = append(parameterList, params)
	}

	req, err := r.client.NewAuthenticatedRequest(ctx, "POST", "/api/v5/trade/batch-orders", nil, parameterList)
	if err != nil {
		return nil, err
	}

	response, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	var data []OrderResponse
	if err := json.Unmarshal(apiResponse.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type OrderDetails struct {
	InstrumentType InstrumentType   `json:"instType"`
	InstrumentID   string           `json:"instId"`
	Tag            string           `json:"tag"`
	Price          fixedpoint.Value `json:"px"`
	Quantity       fixedpoint.Value `json:"sz"`

	OrderID       string    `json:"ordId"`
	ClientOrderID string    `json:"clOrdId"`
	OrderType     OrderType `json:"ordType"`
	Side          SideType  `json:"side"`
	PostionSide   string    `json:"posSide"`
	TradeMode     string    `json:"tdMode"`

	// Accumulated fill quantity
	FilledQuantity fixedpoint.Value `json:"accFillSz"`

	FeeCurrency string           `json:"feeCcy"`
	Fee         fixedpoint.Value `json:"fee"`

	// trade related fields
	LastTradeID           string                     `json:"tradeId,omitempty"`
	LastFilledPrice       fixedpoint.Value           `json:"fillPx"`
	LastFilledQuantity    fixedpoint.Value           `json:"fillSz"`
	LastFilledTime        types.MillisecondTimestamp `json:"fillTime"`
	LastFilledFee         fixedpoint.Value           `json:"fillFee"`
	LastFilledFeeCurrency string                     `json:"fillFeeCcy"`
	LastFilledPnl         fixedpoint.Value           `json:"fillPnl"`
	BillID                types.StrInt64             `json:"billId"`

	// ExecutionType = liquidity (M = maker or T = taker)
	ExecutionType string `json:"execType"`

	// Average filled price. If none is filled, it will return 0.
	AveragePrice fixedpoint.Value `json:"avgPx"`

	// Currency = Margin currency
	// Only applicable to cross MARGIN orders in Single-currency margin.
	Currency string `json:"ccy"`

	// Leverage = from 0.01 to 125.
	// Only applicable to MARGIN/FUTURES/SWAP
	Leverage fixedpoint.Value `json:"lever"`

	RebateCurrency string           `json:"rebateCcy"`
	Rebate         fixedpoint.Value `json:"rebate"`

	PnL fixedpoint.Value `json:"pnl"`

	UpdateTime   types.MillisecondTimestamp `json:"uTime"`
	CreationTime types.MillisecondTimestamp `json:"cTime"`

	State OrderState `json:"state"`
}

type GetOrderDetailsRequest struct {
	client *RestClient

	instId  string
	ordId   *string
	clOrdId *string
}

func (r *GetOrderDetailsRequest) InstrumentID(instId string) *GetOrderDetailsRequest {
	r.instId = instId
	return r
}

func (r *GetOrderDetailsRequest) OrderID(orderID string) *GetOrderDetailsRequest {
	r.ordId = &orderID
	return r
}

func (r *GetOrderDetailsRequest) ClientOrderID(clientOrderID string) *GetOrderDetailsRequest {
	r.clOrdId = &clientOrderID
	return r
}

func (r *GetOrderDetailsRequest) QueryParameters() url.Values {
	var values = url.Values{}

	values.Add("instId", r.instId)

	if r.ordId != nil {
		values.Add("ordId", *r.ordId)
	} else if r.clOrdId != nil {
		values.Add("clOrdId", *r.clOrdId)
	}

	return values
}

func (r *GetOrderDetailsRequest) Do(ctx context.Context) (*OrderDetails, error) {
	params := r.QueryParameters()
	req, err := r.client.NewAuthenticatedRequest(ctx, "GET", "/api/v5/trade/order", params, nil)
	if err != nil {
		return nil, err
	}

	response, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	var data []OrderDetails
	if err := json.Unmarshal(apiResponse.Data, &data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("get order details error")
	}

	return &data[0], nil
}

type GetTransactionDetailsRequest struct {
	client *RestClient

	instType *InstrumentType

	instId *string

	ordId *string
}

func (r *GetTransactionDetailsRequest) InstrumentType(instType InstrumentType) *GetTransactionDetailsRequest {
	r.instType = &instType
	return r
}

func (r *GetTransactionDetailsRequest) InstrumentID(instId string) *GetTransactionDetailsRequest {
	r.instId = &instId
	return r
}

func (r *GetTransactionDetailsRequest) OrderID(orderID string) *GetTransactionDetailsRequest {
	r.ordId = &orderID
	return r
}

func (r *GetTransactionDetailsRequest) Parameters() map[string]interface{} {
	var payload = map[string]interface{}{}

	if r.instType != nil {
		payload["instType"] = r.instType
	}

	if r.instId != nil {
		payload["instId"] = r.instId
	}

	if r.ordId != nil {
		payload["ordId"] = r.ordId
	}

	return payload
}

func (r *GetTransactionDetailsRequest) Do(ctx context.Context) ([]OrderDetails, error) {
	payload := r.Parameters()
	req, err := r.client.NewAuthenticatedRequest(ctx, "GET", "/api/v5/trade/fills", nil, payload)
	if err != nil {
		return nil, err
	}

	response, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	if err := response.DecodeJSON(&apiResponse); err != nil {
		return nil, err
	}
	var data []OrderDetails
	if err := json.Unmarshal(apiResponse.Data, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type PositionDetails struct {
	InstrumentType string `json:"instType"`
	InstrumentID   string `json:"instId"`

	MarginMode string           `json:"mgnMode"`
	PosId      string           `json:"posId"`
	PosSide    string           `json:"posSide"`
	Pos        fixedpoint.Value `json:"pos"`
	BaseBal    fixedpoint.Value `json:"baseBal"`
	QuoteBal   fixedpoint.Value `json:"quoteBal"`
	PosCCY     string           `json:"posCcy"`
	AvailPos   fixedpoint.Value `json:"availPos"`
	AvgPx      fixedpoint.Value `json:"avgPx"`
	Upl        fixedpoint.Value `json:"upl"`
	UplRatio   fixedpoint.Value `json:"uplRatio"`

	Lever       fixedpoint.Value `json:"lever"`
	MarkPx      fixedpoint.Value `json:"markPx"`
	IMR         fixedpoint.Value `json:"imr"`
	Margin      fixedpoint.Value `json:"margin"`
	MarginRatio fixedpoint.Value `json:"mgnRatio"`
	TradeID     string           `json:"tradeId"`
	CCY         string           `json:"ccy"`
	Last        fixedpoint.Value `json:"last"`

	LiabCcy string           `json:"liabCcy"`
	Liab    fixedpoint.Value `json:"liab"`

	PushTime     types.MillisecondTimestamp `json:"pTime"`
	UpdateTime   types.MillisecondTimestamp `json:"uTime"`
	CreationTime types.MillisecondTimestamp `json:"cTime"`
}

type ClosePositionRequest struct {
	client *RestClient

	instrumentID string  `param:"instId"`
	posSide      *string `param:"posSide"`
	marginMode   string  `param:"mgnMode"`
	ccy          *string `param:"ccy"`
	autoCxl      *string `param:"autoCxl"`
	clOrdId      *string `param:"clOrdId"`
	tag          *string `param:"tag"`
}

type ClosePositionResponse struct {
	InstrumentID string `param:"instId"`
	PosSide      string `param:"posSide"`
	ClOrdId      string `param:"clOrdId"`
	Tag          string `param:"tag"`
}

func (r *ClosePositionRequest) Parameters() map[string]interface{} {
	payload, _ := r.GetParameters()
	return payload
}

func (r *ClosePositionRequest) Do(ctx context.Context) ([]ClosePositionResponse, error) {
	payload, err := r.GetParameters()
	if err != nil {
		return nil, err
	}

	log.WithField("payload", payload).
		Debug("ClosePositionRequest payload")

	req, err := r.client.NewAuthenticatedRequest(ctx, "POST", "/api/v5/trade/close-position", nil, payload)
	if err != nil {
		return nil, err
	}

	response, err := r.client.SendRequest(req)
	if err != nil {
		return nil, err
	}

	log.WithField("response", response).
		Debug("ClosePositionRequest response")

	var orderResponse struct {
		Code    string                  `json:"code"`
		Message string                  `json:"msg"`
		Data    []ClosePositionResponse `json:"data"`
	}
	if err := response.DecodeJSON(&orderResponse); err != nil {
		return nil, err
	}

	return orderResponse.Data, nil
}
