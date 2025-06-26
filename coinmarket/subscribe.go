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
)

type CryptoData struct {
	D struct {
		ID       int     `json:"id"`
		P        float64 `json:"p"`        // price
		V        float64 `json:"v"`        // 24h volume
		P1H      float64 `json:"p1h"`      // 1 hour %
		P24H     float64 `json:"p24h"`     // 24 hour %
		P7D      float64 `json:"p7d"`      // 7 day %
		P30D     float64 `json:"p30d"`     // 30 day %
		P3M      float64 `json:"p3m"`      // 3 month %
		P1Y      float64 `json:"p1y"`      // 1 year %
		PYTD     float64 `json:"pytd"`     // Year to date %
		PAll     float64 `json:"pall"`     // All time %
		TS       float64 `json:"ts"`       // total supply
		AS       float64 `json:"as"`       // available supply
		FMC      float64 `json:"fmc"`      // fully diluted market cap
		MC       float64 `json:"mc"`       // market cap
		MC24HPC  float64 `json:"mc24hpc"`  // market cap 24h %
		Vol24HPC float64 `json:"vol24hpc"` // volume 24h %
		FMC24HPC float64 `json:"fmc24hpc"` // full market cap 24h %
		D        float64 `json:"d"`        // dominance
		VD       float64 `json:"vd"`       // volume dominance
	} `json:"d"`
	T string `json:"t"` // timestamp, appears to be a string representation of an int
	C string `json:"c"` // context or source tag
}

func Subscribe(currencies ...int) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(WSCoinMarketURL+currencySubscribeEndpoint, map[string][]string{})
	if err != nil {
		return nil, err
	}
	newConn := &Conn{conn: conn}

	var stringCurrencies []string
	for _, curr := range currencies {
		stringCurrencies = append(stringCurrencies, strconv.Itoa(curr))
	}

	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_15s@{}@detail\",\"%s\"]}", strings.Join(stringCurrencies, "\",\""))))
	if err != nil {
		return nil, err
	}
	err = conn.WriteMessage(1, []byte(fmt.Sprintf("{\"method\":\"RSUBSCRIPTION\",\"params\":[\"main-site@crypto_price_5s@{}@normal\",\"%s\"]}", strings.Join(stringCurrencies, "\",\""))))
	if err != nil {
		return nil, err
	}

	go newConn.startUpdatingPrice()
	return newConn, nil
}
