package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const apiUrl = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"

type CryptoData struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
}

func fetchCryptoData() ([]CryptoData, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP запрос вернул код ошибки: %d", resp.StatusCode)
	}

	var data []CryptoData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func updateCryptoPrices(prices *map[string]float64, mu *sync.Mutex) {
	for {
		data, err := fetchCryptoData()
		if err != nil {
			fmt.Printf("Error fetching data: %v\n", err)
		}
		p, err := preparePriceMap(data)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		mu.Lock()
		*prices = p
		mu.Unlock()
		time.Sleep(10 * time.Minute)
	}
}
func preparePriceMap(data []CryptoData) (map[string]float64, error) {
	prices := make(map[string]float64)

	for _, crypto := range data {
		prices[crypto.Name] = crypto.CurrentPrice
	}
	return prices, nil
}

func main() {
	prices := make(map[string]float64)
	mu := &sync.Mutex{}
	go updateCryptoPrices(&prices, mu)
	for {
		fmt.Print("Enter a crypto: ")
		var inputString string
		_, err := fmt.Scanln(&inputString)
		if err != nil {
			fmt.Println("Error reading input:", err)
		}
		mu.Lock()
		if prices[inputString] <= 0 {
			fmt.Println("Input error or cryptocurrency not found ", err)
		} else {
			fmt.Printf("%s price: $%.2f\n", inputString, prices[inputString])
		}
		mu.Unlock()

		var newline string
		fmt.Scanln(&newline)
	}
}
