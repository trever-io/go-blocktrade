package blocktrade

import "encoding/json"

const FEES_ENDPOINT = "/fees"

type FeeResponse struct {
	Trading                map[string]interface{} `json:"TRADING"`
	TransferInCreditCard   map[string]interface{} `json:"TRANSFER_IN_CREDIT_CARD"`
	TransferInSepa         map[string]interface{} `json:"TRANSFER_IN_SEPA"`
	TransferInOutsideSepa  map[string]interface{} `json:"TRANSFER_IN_OUTSIDE_SEPA"`
	TransferOutSepa        map[string]interface{} `json:"TRANSFER_OUT_SEPA"`
	TransferOutOutsideSepa map[string]interface{} `json:"TRANSFER_OUT_OUTSIDE_SEPA"`
	MinerFee               map[string]interface{} `json:"MINER_FEE"`
	TransferOut            map[string]interface{} `json:"TRANSFER_OUT"`
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
