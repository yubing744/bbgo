package okexapi

import "github.com/c9s/requestgen"

type TradeMode string

const (
	TradeModeCash     TradeMode = "cash"
	TradeModeIsolated TradeMode = "isolated"
	TradeModeCross    TradeMode = "cross"
)

type TargetCurrency string

const (
	TargetCurrencyBase  TargetCurrency = "base_ccy"
	TargetCurrencyQuote TargetCurrency = "quote_ccy"
)

//go:generate -command GetRequest requestgen -method GET -responseType .APIResponse -responseDataField Data
//go:generate -command PostRequest requestgen -method POST -responseType .APIResponse -responseDataField Data

type OrderResponse struct {
	OrderID       string `json:"ordId"`
	ClientOrderID string `json:"clOrdId"`
	Tag           string `json:"tag"`
	Code          string `json:"sCode"`
	Message       string `json:"sMsg"`
}

//go:generate PostRequest -url "/api/v5/trade/order" -type PlaceOrderRequest -responseDataType []OrderResponse
type PlaceOrderRequest struct {
	client requestgen.AuthenticatedAPIClient

	instrumentID string `param:"instId"`

	// tdMode
	// margin mode: "cross", "isolated"
	// non-margin mode cash
	tradeMode TradeMode `param:"tdMode" validValues:"cross,isolated,cash"`

	// Margin currency
	// Only applicable to cross MARGIN orders in Single-currency margin
	marginCurrency *string `param:"ccy"`

	// A combination of case-sensitive alphanumerics, all numbers, or all letters of up to 32 characters.
	clientOrderID *string `param:"clOrdId"`

	// A combination of case-sensitive alphanumerics, all numbers, or all letters of up to 8 characters.
	tag *string `param:"tag"`

	// "buy" or "sell"
	side SideType `param:"side" validValues:"buy,sell"`

	// Position side
	// The default is net in the net mode
	// It is required in the long/short mode, and can only be long or short.
	// Only applicable to FUTURES/SWAP.
	positionSide *string `param:"posSide"`

	//Order type
	//  market: Market order
	//  limit: Limit order
	//  post_only: Post-only order
	//  fok: Fill-or-kill order
	//  ioc: Immediate-or-cancel order
	//  optimal_limit_ioc: Market order with immediate-or-cancel order (applicable only to Expiry Futures and Perpetual Futures).
	//  mmp: Market Maker Protection (only applicable to Option in Portfolio Margin mode)
	//  mmp_and_post_only: Market Maker Protection and Post-only order(only applicable to Option in Portfolio Margin mode)
	orderType OrderType `param:"ordType"`

	// Quantity to buy or sell
	size string `param:"sz"`

	// Order price. Only applicable to limit,post_only,fok,ioc,mmp,mmp_and_post_only order.
	// When placing an option order, one of px/pxUsd/pxVol must be filled in, and only one can be filled in
	price *string `param:"px"`

	// Take-profit trigger price
	takeProfitTriggerPx *string `param:"tpTriggerPx"`
	// Take-profit order price
	takeProfitOrdPx *string `param:"tpOrdPx"`
	// Take-profit trigger price type
	takeProfitTriggerPxType *string `param:"tpOrdPx"`

	// Stop-loss trigger price
	stopLossTriggerPx *string `param:"slTriggerPx"`
	// Stop-loss order price
	stopLossOrdPx *string `param:"slOrdPx"`
	// Stop-loss trigger price type
	stopLossTriggerPxType *string `param:"slTriggerPxType"`

	// Whether the target currency uses the quote or base currency.
	// base_ccy: Base currency ,quote_ccy: Quote currency
	// Only applicable to SPOT Market Orders
	// Default is quote_ccy for buy, base_ccy for sell
	targetCurrency *TargetCurrency `param:"tgtCcy" validValues:"quote_ccy,base_ccy"`
}

func (c *RestClient) NewPlaceOrderRequest() *PlaceOrderRequest {
	return &PlaceOrderRequest{
		client:    c,
		tradeMode: TradeModeCash,
	}
}
