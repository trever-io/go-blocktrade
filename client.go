package blocktrade

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var Debug = false

const TOO_MANY_REQUEST_MSG = "Too Many Requests"

type APIClient struct {
	apiKey    string
	apiSecret string

	assetCache map[int64]*TradingAsset
	pairCache  map[int64]*TradingPair
	portfolio  *Portfolio

	nonceMtx     sync.Mutex
	wsConn       *websocket.Conn
	wsHandlers   map[MessageType]interface{}
	wsHandlerMtx sync.Mutex
}

type APIError struct {
	Message         string      `json:"-"`
	MessageInternal interface{} `json:"message"`
	Code            int         `json:"-"`
}

func newAPIError(code int) *APIError {
	return &APIError{
		Code: code,
	}
}

func (e *APIError) Error() string {
	message := fmt.Sprintf("API Error: Code(%d) %v", e.Code, e.Message)
	return message
}

func NewClient(apiKey, apiSecret string) *APIClient {
	return &APIClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,

		assetCache: make(map[int64]*TradingAsset),
		pairCache:  make(map[int64]*TradingPair),

		wsHandlers: make(map[MessageType]interface{}),
	}
}

func (a *APIClient) Close() {
	if a.wsConn != nil {
		a.wsConn.Close()
		a.wsConn = nil
	}
}

const API_URL = "https://trade.blocktrade.com/api/v1"
const API_KEY_HEADER = "X-Api-Key"
const NONCE_HEADER = "X-Nonce"
const SIGNATURE_HEADER = "X-Signature"
const CONTENT_TYPE_HEADER = "Content-Type"
const CONTENT_TYPE = "application/json"

func (a *APIClient) nonceAndSignature(request interface{}) (int64, string, error) {
	nonce := time.Now().UTC().UnixNano() / 1e3
	message := fmt.Sprintf("%v.%v", a.apiKey, nonce)
	if request != nil {
		requestb, err := json.Marshal(request)
		if err != nil {
			return 0, "", err
		}

		message += fmt.Sprintf(".%v", string(requestb))
	}

	h := hmac.New(sha256.New, []byte(a.apiSecret))
	h.Write([]byte(message))
	sha := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return nonce, sha, nil
}

func (a *APIClient) requestPublicGET(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%v%v", API_URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return a.doRequest(req)
}

func (a *APIClient) requestPOST(endpoint string, request interface{}) ([]byte, error) {
	if a.apiKey == "" || a.apiSecret == "" {
		return nil, errors.New("missing credentials")
	}

	a.nonceMtx.Lock()
	defer a.nonceMtx.Unlock()
	nonce, sig, err := a.nonceAndSignature(request)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)
	url := fmt.Sprintf("%v%v", API_URL, endpoint)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(API_KEY_HEADER, a.apiKey)
	req.Header.Add(NONCE_HEADER, fmt.Sprint(nonce))
	req.Header.Add(SIGNATURE_HEADER, sig)
	req.Header.Add(CONTENT_TYPE_HEADER, CONTENT_TYPE)

	return a.doRequest(req)
}

func (a *APIClient) requestGET(endpoint string) ([]byte, error) {
	return a.requestNoBody(endpoint, "GET")
}

func (a *APIClient) requestNoBody(endpoint string, method string) ([]byte, error) {
	if a.apiKey == "" || a.apiSecret == "" {
		return nil, errors.New("missing credentials")
	}

	a.nonceMtx.Lock()
	defer a.nonceMtx.Unlock()
	nonce, sig, err := a.nonceAndSignature(nil)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v%v", API_URL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(API_KEY_HEADER, a.apiKey)
	req.Header.Add(NONCE_HEADER, fmt.Sprint(nonce))
	req.Header.Add(SIGNATURE_HEADER, sig)

	return a.doRequest(req)
}

func (a *APIClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if Debug {
		log.Println(string(b))
	}

	if resp.StatusCode >= 300 {
		apiErr := newAPIError(resp.StatusCode)

		if resp.StatusCode == http.StatusTooManyRequests {
			apiErr.Message = TOO_MANY_REQUEST_MSG
		}

		err := json.Unmarshal(b, &apiErr)
		if err != nil {
			return nil, err
		}

		if vList, ok := apiErr.MessageInternal.([]interface{}); ok {
			for _, v := range vList {
				if s, ok := v.(string); ok {
					apiErr.Message += s + ", "
				}
			}
		}

		if s, ok := apiErr.MessageInternal.(string); ok {
			apiErr.Message = s
		}

		return nil, apiErr
	}

	return b, nil
}
