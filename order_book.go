package blocktrade

import (
	"encoding/json"
	"fmt"
)

const ORDER_BOOK_ENDPOINT = "/order_book"

type OrderBookEntry struct {
	Amount string `json:"amount"`
	Price  string `json:"price"`
	Value  string `json:"value"`
}

type OrderBookResponse struct {
	Asks []OrderBookEntry `json:"asks"`
	Bids []OrderBookEntry `json:"bids"`
}

func (a *APIClient) GetOrderBook(tradingPairId int64) (*OrderBookResponse, error) {
	url := fmt.Sprintf("%v/%d", ORDER_BOOK_ENDPOINT, tradingPairId)
	b, err := a.requestPublicGET(url)
	if err != nil {
		return nil, err
	}

	resp := new(OrderBookResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
