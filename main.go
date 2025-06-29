package main

import (
	"fmt"
	"github.com/restartfu/coinmarket/coinmarket"
	"log"
)

func main() {
	cmc, err := coinmarket.Subscribe(coinmarket.CurrencyLTC, coinmarket.CurrencyBTC)
	if err != nil {
		log.Fatal(err)
	}

	for {
		data := cmc.Poll(coinmarket.CurrencyBTC)
		fmt.Printf("$%v\n", data.Price)
	}
}
