package blocktrade

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

const WS_URL = "wss://trade.blocktrade.com/ws/v1/notification"
const MESSAGE_BUFFER_SIZE = 10

type MessageType string

const MessageType_UserOrders = "user_orders"

type UserOrderHandlerFunc func(orderResponse *OrderResponse, err error)

type blocktradeWebsocketMessage struct {
	MessageType MessageType            `json:"message_type"`
	Payload     map[string]interface{} `json:"payload"`
}

type blockTradeUserOrdersWsResponse struct {
	Data []*OrderResponse `json:"data"`
}

type websocketMessage struct {
	Message []byte
	Error   error
}

func (a *APIClient) Websocket() (<-chan error, error) {
	wsChan := make(chan websocketMessage, MESSAGE_BUFFER_SIZE)
	wsCloseChan := make(chan error)

	conn, _, err := websocket.DefaultDialer.Dial(WS_URL, nil)
	if err != nil {
		return nil, err
	}
	a.wsConn = conn
	go a.receiveWsMessages(a.wsConn, wsChan)
	go a.handleWsMessages(wsChan, wsCloseChan)

	return wsCloseChan, nil
}

func (a *APIClient) handleWsMessages(wsChan chan websocketMessage, wsCloseChan chan error) {
	for {
		msg := <-wsChan
		if msg.Error != nil {
			wsCloseChan <- msg.Error
			close(wsCloseChan)
			return
		}

		if Debug {
			log.Printf("WS EVENT: %v\n", string(msg.Message))
		}

		wsMsg := new(blocktradeWebsocketMessage)
		err := json.Unmarshal(msg.Message, &wsMsg)
		if err != nil {
			log.Printf("WS ERROR: %v\n", err)
			continue
		}

		a.wsHandlerMtx.Lock()
		switch wsMsg.MessageType {
		case MessageType_UserOrders:
			var f UserOrderHandlerFunc
			if val, ok := a.wsHandlers[MessageType_UserOrders]; ok {
				f = val.(UserOrderHandlerFunc)
			} else {
				break
			}

			b, err := json.Marshal(wsMsg.Payload)
			if err != nil {
				f(nil, err)
				break
			}

			orderResponse := new(blockTradeUserOrdersWsResponse)
			err = json.Unmarshal(b, &orderResponse)
			if err != nil {
				f(nil, err)
				break
			}

			for _, order := range orderResponse.Data {
				f(order, nil)
			}
		default:
			log.Printf("Unhandled message_type: %v\n", wsMsg.MessageType)
		}
		a.wsHandlerMtx.Unlock()
	}
}

func (a *APIClient) receiveWsMessages(conn *websocket.Conn, wsChan chan websocketMessage) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			wsChan <- websocketMessage{Error: err}
			close(wsChan)
			return
		}

		wsChan <- websocketMessage{Message: msg}
	}
}

func (a *APIClient) SubscribeUserOrders(f UserOrderHandlerFunc) error {
	if a.wsConn == nil {
		return errors.New("websocket not initialized")
	}

	a.wsHandlerMtx.Lock()
	a.wsHandlers[MessageType_UserOrders] = f
	a.wsHandlerMtx.Unlock()

	userResp, err := a.User()
	if err != nil {
		return err
	}

	subscribeMessage := map[string]interface{}{
		"subscribe_user_orders": map[string]interface{}{
			"auth_token": userResp.WebsocketAuthToken,
		},
	}

	err = a.wsConn.WriteJSON(subscribeMessage)
	return err
}

func (a *APIClient) UnsubscribeUserOrders() error {
	if a.wsConn == nil {
		return errors.New("websocket not initialized")
	}

	unsubcribeMessage := map[string]interface{}{
		"unsubscribe_user_orders": map[string]interface{}{},
	}

	err := a.wsConn.WriteJSON(unsubcribeMessage)
	if err != nil {
		return err
	}

	a.wsHandlerMtx.Lock()
	delete(a.wsHandlers, MessageType_UserOrders)
	a.wsHandlerMtx.Unlock()

	return nil
}
