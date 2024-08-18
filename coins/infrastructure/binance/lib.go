package binance

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

type MiniTicker struct {
	EventType string `json:"e"`
	Timestamp int64  `json:"E"`
	Symbol    string `json:"s"`
	Close     string `json:"c"`
	Open      string `json:"o"`
	High      string `json:"h"`
	Low       string `json:"l"`
	V         string `json:"v"`
	Q         string `json:"q"`
}

const pingPeriod = 3 * time.Minute

type WsOptions struct {
	Url string
}

type WsConnection struct {
	conn *websocket.Conn
}

func (ws *WsConnection) Close() {
	err := ws.conn.Close()
	if err != nil {
		log.Fatalf("failed to close connection: %v", err)
	}
}

func (ws *WsConnection) Ping() error {
	return ws.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
}

func (ws *WsConnection) Pong() error {
	return ws.conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
}

func (ws *WsConnection) Send(data string) error {
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (ws *WsConnection) Subscribe(params string, id string) error {
	return ws.Send(`{"method": "SUBSCRIBE", "params": ` + params + `, "id": ` + id + `}`)
}

func (ws *WsConnection) Listen(receiver chan<- []byte) {
	go ws.pingLoop()
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		receiver <- message
	}
}

func (ws *WsConnection) pingLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := ws.Pong()
			if err != nil {
				log.Println("Error sending pong:", err)
			}
		}
	}
}

var conn *WsConnection

func Connect(opt *WsOptions) *WsConnection {
	if conn != nil {
		return conn
	}

	c, _, err := websocket.DefaultDialer.Dial(opt.Url, nil)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	conn = &WsConnection{conn: c}

	return conn
}
