package coinmarket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
)

const (
	WSCoinMarketURL           = "wss://push.coinmarketcap.com/"
	currencySubscribeEndpoint = "ws?device=web&client_source=coin_detail_page"
)

const (
	CurrencyBTC = iota + 1
	CurrencyLTC
)

func Subscribe(currencies ...int) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(WSCoinMarketURL+currencySubscribeEndpoint, map[string][]string{})
	if err != nil {
		return nil, err
	}
	newConn := &Conn{conn: conn, data: make(map[int]cryptoData)}

	var stringCurrencies []string
	for _, curr := range currencies {
		stringCurrencies = append(stringCurrencies, strconv.Itoa(curr))
	}

	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_15s@{}@detail\",\"%s\"]}", strings.Join(stringCurrencies, ","))))
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_5s@{}@normal\",\"%s\"]}", strings.Join(stringCurrencies, ","))))
	if err != nil {
		return nil, err
	}

	go newConn.startUpdatingPrice()
	return newConn, nil
}
