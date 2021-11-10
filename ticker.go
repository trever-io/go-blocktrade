package blocktrade

import (
	"encoding/json"
	"fmt"
)

const TICKER_ENDPOINT = "/ticker"

type TickerData struct {
	AskPrice  string `json:"ask_price"`
	BidPrice  string `json:"bid_price"`
	LastPrice string `json:"last_price"`
	Volume    string `json:"volume"`
	High      string `json:"high"`
	Low       string `json:"low"`
}

type TickerResponse struct {
	TradingPairId int64      `json:"trading_pair_id"`
	Data          TickerData `json:"data"`
}

func (a *APIClient) GetTicker(tradingPairId int64) (*TickerData, error) {
	url := fmt.Sprintf("%v/%d", TICKER_ENDPOINT, tradingPairId)
	b, err := a.requestPublicGET(url)
	if err != nil {
		return nil, err
	}

	resp := new(TickerData)
	err = json.Unmarshal(b, resp)
	return resp, err
}
