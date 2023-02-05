package okexapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *ClosePositionRequest) InstrumentID(instrumentID string) *ClosePositionRequest {
	c.instrumentID = instrumentID
	return c
}

func (c *ClosePositionRequest) MarginMode(mgnMode string) *ClosePositionRequest {
	c.marginMode = mgnMode
	return c
}

func (c *ClosePositionRequest) CCY(ccy string) *ClosePositionRequest {
	c.ccy = &ccy
	return c
}

func (c *ClosePositionRequest) Tag(tag string) *ClosePositionRequest {
	c.tag = &tag
	return c
}

func (c *ClosePositionRequest) GetParameters() (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	// check instrumentID field -> json key instId
	instrumentID := c.instrumentID

	// assign parameter of instrumentID
	params["instId"] = instrumentID
	params["mgnMode"] = c.marginMode

	if c.posSide != nil {
		posSide := *c.posSide

		// assign parameter of orderID
		params["posSide"] = posSide
	}

	if c.ccy != nil {
		ccy := *c.ccy

		// assign parameter of clientOrderID
		params["ccy"] = ccy
	}

	if c.autoCxl != nil {
		autoCxl := *c.autoCxl

		// assign parameter of clientOrderID
		params["autoCxl"] = autoCxl
	}

	if c.clOrdId != nil {
		clOrdId := *c.clOrdId

		// assign parameter of clientOrderID
		params["clOrdId"] = clOrdId
	}

	if c.tag != nil {
		tag := *c.tag

		// assign parameter of clientOrderID
		params["tag"] = tag
	}

	return params, nil
}

func (c *ClosePositionRequest) GetParametersQuery() (url.Values, error) {
	query := url.Values{}

	params, err := c.GetParameters()
	if err != nil {
		return query, err
	}

	for k, v := range params {
		query.Add(k, fmt.Sprintf("%v", v))
	}

	return query, nil
}

func (c *ClosePositionRequest) GetParametersJSON() ([]byte, error) {
	params, err := c.GetParameters()
	if err != nil {
		return nil, err
	}

	return json.Marshal(params)
}
