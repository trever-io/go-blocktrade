package blocktrade

import (
	"encoding/json"
	"fmt"
)

const CUSTOMER_ORDERS_ENDPOINT = "/customer_orders"

type Direction string

const Direction_BUY Direction = "BUY"
const Direction_SELL Direction = "SELL"

type Type string

const Type_LIMIT Type = "LIMIT"
const Type_MARKET Type = "MARKET"

type TimeInForce string

const TimeInForce_GTC TimeInForce = "GTC"

type Status string

const Status_NEW = "NEW"
const Status_PARTIALLY_FILLED = "PARTIALLY_FILLED"
const Status_FILLED = "FILLED"
const Status_CANCELLED = "CANCELLED"

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

type CreateOrderResponse struct {
	Id              int64  `json:"id"`
	CustomerOrderId string `json:"customer_order_id"`
}

type OrderResponse struct {
	Id              int64                 `json:"id"`
	CustomerOrderId string                `json:"customer_order_id"`
	PortfolioId     int64                 `json:"portfolio_id"`
	TradingPairId   int64                 `json:"trading_pair_id"`
	Direction       Direction             `json:"direction"`
	Type            Type                  `json:"type"`
	Amount          string                `json:"amount"`
	RemainingAmount string                `json:"remaining_amount"`
	Price           string                `json:"price"`
	TimeInForce     TimeInForce           `json:"time_in_force"`
	StopPrice       string                `json:"stop_price"`
	Date            int64                 `json:"date"`
	Status          Status                `json:"status"`
	Trades          []*OrderTradeResponse `json:"trades"`
}

type OrderTradeResponse struct {
	Id    int64  `json:"id"`
	Value string `json:"value"`
	Price string `json:"price"`
	Time  int64  `json:"time"`
}

func (a *APIClient) CreateCustomerOrder(request *CustomerOrderRequest) (*CreateOrderResponse, error) {
	b, err := a.requestPOST(CUSTOMER_ORDERS_ENDPOINT, request)
	if err != nil {
		return nil, err
	}

	resp := new(CreateOrderResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}

func (a *APIClient) GetCustomerOrder(customerOrderId string) (*OrderResponse, error) {
	url := fmt.Sprintf("%v/%v", CUSTOMER_ORDERS_ENDPOINT, customerOrderId)
	b, err := a.requestGET(url)
	if err != nil {
		return nil, err
	}

	resp := new(OrderResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}

func (a *APIClient) CancelCustomerOrder(customerOrderId string) error {
	url := fmt.Sprintf("%v/%v/cancel", CUSTOMER_ORDERS_ENDPOINT, customerOrderId)
	_, err := a.requestNoBody(url, "POST")
	return err
}
