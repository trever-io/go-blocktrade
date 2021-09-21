package blocktrade

import (
	"encoding/json"
	"fmt"
)

const ORDERS_ENDPOINT = "/orders"

func (a *APIClient) GetOrder(id int64) (*OrderResponse, error) {
	url := fmt.Sprintf("%v/%d", ORDERS_ENDPOINT, id)
	b, err := a.requestGET(url)
	if err != nil {
		return nil, err
	}

	resp := new(OrderResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
