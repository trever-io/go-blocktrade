package blocktrade

import "encoding/json"

const CUSTOMER_ORDERS_ENDPOINT = "/customer_orders"

type Direction string

const Direction_BUY Direction = "BUY"
const Direction_SELL Direction = "SELL"

type Type string

const Type_LIMIT Type = "LIMIT"
const Type_MARKET Type = "MARKET"

type TimeInForce string

const TimeInForce_GTC TimeInForce = "GTC"

type CustomerOrderRequest struct {
	CustomerOrderId string      `json:"customer_order_id"`
	PortfolioId     int64       `json:"portfolio_id"`
	Direction       Direction   `json:"direction"`
	Type            Type        `json:"type"`
	TradingPairId   int64       `json:"trading_pair_id"`
	Amount          string      `json:"amount"`
	Price           string      `json:"price,omitempty"`
	TimeInForce     TimeInForce `json:"time_in_force,omitempty"`
	StopPrice       string      `json:"stop_price,omitempty"`
}

type OrderResponse struct {
	Id              int64  `json:"id"`
	CustomerOrderId string `json:"customer_order_id"`
}

func (a *APIClient) CreateCustomerOrder(request *CustomerOrderRequest) (*OrderResponse, error) {
	b, err := a.requestPOST(CUSTOMER_ORDERS_ENDPOINT, request)
	if err != nil {
		return nil, err
	}

	resp := new(OrderResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
