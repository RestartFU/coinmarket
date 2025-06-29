package main

import (
	"fmt"
	"github.com/restartfu/coinmarket/coinmarket"
	"log"
	"time"
)

func main() {
	conn, err := coinmarket.Subscribe(coinmarket.CurrencyLTC, coinmarket.CurrencyBTC)
	if err != nil {
		log.Fatal(err)
	}

	for {
		<-time.After(time.Second)
		p, _ := conn.Price(coinmarket.CurrencyLTC)
		fmt.Printf("$%.2f\n", p)
		p, _ = conn.Price(coinmarket.CurrencyBTC)
		fmt.Printf("$%.2f\n", p)
	}
}
