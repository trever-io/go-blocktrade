package blocktrade

type TradeResponse struct {
	Id            int64     `json:"id"`
	OrderId       int64     `json:"order_id"`
	TradingPairId int64     `json:"trading_pair_id"`
	Symbol        string    `json:"symbol"`
	Direction     Direction `json:"direction"`
	Amount        string    `json:"amount"`
	Price         string    `json:"price"`
	Date          int64     `json:"date"`
	FeeValue      string    `json:"fee_value"`
	TradeValue    string    `json:"trade_value"`
	Make          bool      `json:"maker"`
}
