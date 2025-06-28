package coinmarket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
)

const (
	WSCoinMarketURL           = "wss://push.coinmarketcap.com/"
	currencySubscribeEndpoint = "ws?device=web&client_source=coin_detail_page"
)

const (
	CurrencyBTC = iota + 1
	CurrencyLTC
)

func Subscribe(currency int) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(WSCoinMarketURL+currencySubscribeEndpoint, nil)
	if err != nil {
		return nil, err
	}
	newConn := &Conn{conn: conn}

	stringCurrency := strconv.Itoa(currency)

	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_15s@{}@detail\",\"%s\"]}", stringCurrency)))
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_5s@{}@normal\",\"%s\"]}", stringCurrency)))
	if err != nil {
		return nil, err
	}

	go newConn.startUpdatingPrice()
	return newConn, nil
}
