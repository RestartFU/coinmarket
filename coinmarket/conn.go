package coinmarket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn

	price float64
}

func (c *Conn) Price() float64 {
	return c.price
}

func (c *Conn) startUpdatingPrice() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var cryptoData CryptoData
		err = json.Unmarshal(message, &cryptoData)
		if err != nil {
			return
		}

		c.price = cryptoData.D.P
	}
}
