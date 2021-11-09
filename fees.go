package blocktrade

import "encoding/json"

const FEES_ENDPOINT = "/fees"

type Fee struct {
	MinFee       string `json:"min_fee"`
	PercentValue string `json:"percent_value"`
}

type FeeResponse struct {
	Trading                map[string]Fee `json:"TRADING"`
	TransferInCreditCard   map[string]Fee `json:"TRANSFER_IN_CREDIT_CARD"`
	TransferInSepa         map[string]Fee `json:"TRANSFER_IN_SEPA"`
	TransferInOutsideSepa  map[string]Fee `json:"TRANSFER_IN_OUTSIDE_SEPA"`
	TransferOutSepa        map[string]Fee `json:"TRANSFER_OUT_SEPA"`
	TransferOutOutsideSepa map[string]Fee `json:"TRANSFER_OUT_OUTSIDE_SEPA"`
	MinerFee               map[string]Fee `json:"MINER_FEE"`
	TransferOut            map[string]Fee `json:"TRANSFER_OUT"`
}

func (a *APIClient) Fees() (*FeeResponse, error) {
	b, err := a.requestGET(FEES_ENDPOINT)
	if err != nil {
		return nil, err
	}

	resp := new(FeeResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
