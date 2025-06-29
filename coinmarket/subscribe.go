package coinmarket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"sync"
)

const cmcURL = "wss://push.coinmarketcap.com/ws?device=web&client_source=coin_detail_page"

func Subscribe(currencies ...int) (*CMC, error) {
	conn, _, err := websocket.DefaultDialer.Dial(cmcURL, nil)
	if err != nil {
		return nil, err
	}
	newConn := &CMC{
		conn:         conn,
		queues:       make(map[int][]CurrencyUpdateData),
		latestUpdate: make(map[int]CurrencyUpdateData),
	}
	newConn.cond = sync.NewCond(&newConn.mu)

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
