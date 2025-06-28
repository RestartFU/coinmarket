package coinmarket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn

	data cryptoData
}

func (c *Conn) Price() float64 {
	return c.data.D.P
}

func (c *Conn) PercentChange() (hour, day, week, month, trimester, year, yearToDate, allTime float64) {
	return c.data.D.P1H, c.data.D.P24H, c.data.D.P7D,
		c.data.D.P30D, c.data.D.P3M, c.data.D.P1Y, c.data.D.PYTD, c.data.D.PAll
}

type cryptoData struct {
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

func (c *Conn) startUpdatingPrice() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		err = json.Unmarshal(message, &c.data)
		if err != nil {
			return
		}
	}
}
