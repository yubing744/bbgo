package okexapi

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
	client *RestClient

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

	orderType OrderType `param:"ordType"`

	quantity string `param:"sz"`

	// price
	price *string `param:"px"`

	// Take-profit trigger price
	tpTriggerPx *string `param:"tpTriggerPx"`
	// Take-profit order price
	tpOrdPx *string `param:"tpOrdPx"`
	// Take-profit trigger price type
	tpTriggerPxType *string `param:"tpOrdPx"`

	// Stop-loss trigger price
	slTriggerPx *string `param:"slTriggerPx"`
	// Stop-loss order price
	slOrdPx *string `param:"slOrdPx"`
	// Stop-loss trigger price type
	slTriggerPxType *string `param:"slTriggerPxType"`
}

func (c *RestClient) NewPlaceOrderRequest() *PlaceOrderRequest {
	return &PlaceOrderRequest{
		client:    c,
		tradeMode: TradeModeCash,
	}
}
