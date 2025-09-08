package fiat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type RateConverter interface {
	ConvertRate(from, to Currency) (float64, error)
}

var (
	ErrRateNotFound   = errors.New("rate not found")
	fiatConversionURL = "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/quote/latest?id=2781,2782,2783,2784,2785,2786,2787,2788,2789,2790,2791,2792,2793,2794,2795,2796,2797,2798,2799,2800,2801,2802,2803,2804,2805,2806,2807,2808,2809,2810,2811,2812,2823,3554,3544,2821,2817,2824,2819,2813,2820,3538,3566,3530,3540,2814,3573,1,1027,2010,1839,6636,52,1975,2,512,1831,7083,74,9023,9022,5824,6783&convertId=2781"
)

type DefaultRateConverter struct {
	usdRates map[string]float64
	close    chan struct{}

	proxies []string
}

func NewDefaultRateConverter(proxies []string) *DefaultRateConverter {
	return &DefaultRateConverter{usdRates: make(map[string]float64), close: make(chan struct{}, 1), proxies: proxies}
}

func (c *DefaultRateConverter) Start(interval time.Duration) {
	c.updateUSDConversionRates()
	for {
		select {
		case <-time.After(interval):
			c.updateUSDConversionRates()
		case <-c.close:
			return
		}
	}
}

func (c *DefaultRateConverter) Stop() {
	close(c.close)
	c.close = make(chan struct{}, 1)
}

type data struct {
	Data []struct {
		Symbol string `json:"symbol"`
		Quotes []struct {
			Price float64 `json:"price"`
		} `json:"quotes"`
	} `json:"data"`
}

func (c *DefaultRateConverter) updateUSDConversionRates() {
	req, err := http.NewRequest("GET", fiatConversionURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}

	if len(c.proxies) >= 1 {
		randomProxyIndex := rand.Intn(len(c.proxies))
		fixedURL, err := url.Parse(c.proxies[randomProxyIndex])
		if err != nil {
			fmt.Println("Error parsing proxy URL:", err)
			return
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(fixedURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var currencies data

	err = json.Unmarshal(body, &currencies)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	for _, currency := range currencies.Data {
		c.usdRates[currency.Symbol] = currency.Quotes[0].Price
	}
}

func (c *DefaultRateConverter) ConvertRate(from, to Currency) (float64, error) {
	if from == to {
		return 1, nil
	}

	usdRate1, ok := c.usdRates[from.symbol]
	if !ok && from != USD {
		return 0, fmt.Errorf("%s: %w", from.String(), ErrRateNotFound)
	}
	usdRate2, ok := c.usdRates[to.symbol]
	if !ok && from != USD {
		return 0, fmt.Errorf("%s: %w", to.String(), ErrRateNotFound)
	}

	rate := usdRate1 / usdRate2
	return rate, nil
}
