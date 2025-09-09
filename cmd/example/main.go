package main

import (
	"fmt"
	"log"
	"time"

	coinmarket "github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/restartfu/coinmarketcap/coinmarketcap/fiat"
)

func main() {
	rateConverter := fiat.NewDefaultRateConverter([]string{"socks5://localhost:9050"})
	go rateConverter.Start(time.Hour)

	subscriber := coinmarket.NewSubscriber(fiat.CAD, rateConverter)
	err := subscriber.Subscribe(coinmarket.CurrencyXMR)
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, err := subscriber.Poll(coinmarket.CurrencyXMR)
		if err != nil {
			subscriber.Revive()
			log.Println(err)
			continue
		}
		fmt.Printf("$%v\n", data.Price)
	}

}
