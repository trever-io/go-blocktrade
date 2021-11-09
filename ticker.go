package blocktrade

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
