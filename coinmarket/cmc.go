package coinmarket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type CMC struct {
	conn *websocket.Conn

	mu     sync.Mutex
	cond   *sync.Cond
	queues map[int][]CurrencyUpdateData

	latestUpdate map[int]CurrencyUpdateData
}

func (c *CMC) LatestCurrencyUpdate(currencyID int) (CurrencyUpdateData, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	d, ok := c.latestUpdate[currencyID]
	return d, ok
}
func (c *CMC) Poll(currencyID int) CurrencyUpdateData {
	c.mu.Lock()
	defer c.mu.Unlock()

	for len(c.queues[currencyID]) == 0 {
		c.cond.Wait()
	}

	update := c.queues[currencyID][0]
	c.queues[currencyID] = c.queues[currencyID][1:] // dequeue
	return update
}

func (c *CMC) startUpdatingPrice() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var data CurrencyUpdate
		err = json.Unmarshal(message, &data)
		if err != nil {
			return
		}
		c.mu.Lock()
		c.latestUpdate[data.Data.ID] = data.Data
		c.queues[data.Data.ID] = append(c.queues[data.Data.ID], data.Data)
		c.mu.Unlock()
		c.cond.Broadcast() // notify all waiting Poll() calls

	}
}
