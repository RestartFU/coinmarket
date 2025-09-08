package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type data struct {
	Data []struct {
		Symbol string `json:"symbol"`
	} `json:"data"`
}

func main() {
	buf, err := os.ReadFile("./assets/latest_fiat.json")
	if err != nil {
		panic(err)
	}

	var currencies data
	err = json.Unmarshal(buf, &currencies)
	if err != nil {
		panic(err)
	}

	dat := `
	package fiat
	var (
	`
	for _, currency := range currencies.Data {
		dat += fmt.Sprintf(`	%s = newCurrency("%s")`, currency.Symbol, currency.Symbol) + "\n"
	}
	dat += `)`

	err = os.Remove("coinmarket/fiat/currencies_generated.go")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	err = os.WriteFile("coinmarketcap/fiat/currencies_generated.go", []byte(dat), 0777)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "fmt", "coinmarketcap/fiat/currencies_generated.go")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
