package coinmarketcap

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/restartfu/coinmarketcap/coinmarketcap/fiat"
)

const cmcURL = "wss://push.coinmarketcap.com/ws?device=web&client_source=coin_detail_page"

type Subscriber struct {
	conn *websocket.Conn

	mu     sync.Mutex
	cond   *sync.Cond
	queues map[int][]CurrencyUpdateData

	latestUpdate  map[int]CurrencyUpdateData
	currencies    []int
	fiatCurrency  fiat.Currency
	rateConverter fiat.RateConverter

	pollingErr error
}

func NewSubscriber(fiatCurrency fiat.Currency, rateConverter fiat.RateConverter) *Subscriber {
	sub := &Subscriber{
		fiatCurrency:  fiatCurrency,
		rateConverter: rateConverter,
		queues:        make(map[int][]CurrencyUpdateData),
		latestUpdate:  make(map[int]CurrencyUpdateData),
	}

	sub.cond = sync.NewCond(&sub.mu)
	return sub
}

func (c *Subscriber) Subscribe(currencies ...int) error {
	c.currencies = currencies

	conn, _, err := websocket.DefaultDialer.Dial(cmcURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	var stringIDs []string
	for _, curr := range currencies {
		stringIDs = append(stringIDs, strconv.Itoa(curr))
	}
	idString := strings.Join(stringIDs, ",")

	topics := []string{
		"main-site@crypto_price_15s@{}@detail",
		"main-site@crypto_price_5s@{}@normal",
	}

	for _, topic := range topics {
		msg := fmt.Sprintf(`{"method":"RSUBSCRIPTION","params":["%s","%s"]}`, topic, idString)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return err
		}
	}

	go c.startUpdatingPrice()
	return nil
}

func (c *Subscriber) LatestCurrencyUpdate(currencyID int) (CurrencyUpdateData, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	d, ok := c.latestUpdate[currencyID]
	return d, ok
}

func (c *Subscriber) Revive() error {
	if c.pollingErr == nil {
		return errors.New("subscriber is already active")
	}
	c.Subscribe(c.currencies...)
	c.pollingErr = nil
	return nil
}

func (c *Subscriber) Poll(currencyID int) (CurrencyUpdateData, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pollingErr != nil {
		return CurrencyUpdateData{}, c.pollingErr
	}

	for len(c.queues[currencyID]) == 0 {
		c.cond.Wait()
	}

	update := c.queues[currencyID][0]
	rate, err := c.rateConverter.ConvertRate(fiat.USD, c.fiatCurrency)
	if err != nil {
		return CurrencyUpdateData{}, err
	}

	update.Price = update.Price * rate
	update.Price1H = update.Price1H * rate
	update.Price24H = update.Price24H * rate
	update.Price7D = update.Price7D * rate
	update.Price30D = update.Price30D * rate
	update.Price3M = update.Price3M * rate
	update.Price1Y = update.Price1Y * rate
	update.PriceAllTime = update.PriceAllTime * rate
	update.PriceYearToDate = update.PriceYearToDate * rate

	c.queues[currencyID] = c.queues[currencyID][1:]
	return update, nil
}

func (c *Subscriber) startUpdatingPrice() {
	var err error

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var data CurrencyUpdate
		err = json.Unmarshal(message, &data)
		if err != nil {
			break
		}
		c.mu.Lock()
		c.latestUpdate[data.Data.ID] = data.Data
		c.queues[data.Data.ID] = append(c.queues[data.Data.ID], data.Data)
		c.mu.Unlock()
		c.cond.Broadcast()
	}

	c.pollingErr = err
	c.cond.Broadcast()
}
