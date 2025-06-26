package main

import (
	"fmt"
	"github.com/restartfu/coinmarket/coinmarket"
	"log"
	"time"
)

func main() {
	conn, err := coinmarket.Subscribe(coinmarket.CurrencyBTC)
	if err != nil {
		log.Fatal(err)
	}

	for {
		<-time.After(time.Second)
		fmt.Printf("$%.2f\n", conn.Price())
	}
}
