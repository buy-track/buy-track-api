package ws_client

import (
	"encoding/json"
	"fmt"
	"log"
	"my-stocks/coins/app"
	"my-stocks/coins/infrastructure/binance"
	"my-stocks/domains"
	"strconv"
	"strings"
)

func RunUpdateCoinPrice(url string, service app.CoinService) {
	ws := binance.Connect(&binance.WsOptions{Url: url})

	ch := make(chan []byte)
	coins := service.GetAllActiveCoins("btc", "ftm", "eth", "sol", "bnb", "doge")
	ln := len(coins)
	if ln == 0 {
		log.Println("There is no coin to update their prices.")
		return
	}
	params := make([]string, ln)
	idHash := make(map[string]string, ln)
	for i, coin := range coins {
		params[i] = coin.Symbol + "usdt@miniTicker"
		idHash[strings.ToUpper(coin.Symbol+"usdt")] = coin.ID
	}
	jParams, _ := json.Marshal(params)

	err := ws.Subscribe(string(jParams), "1")
	if err != nil {
		fmt.Println(err)
	}

	go ws.Listen(ch)
	for {
		select {
		case message := <-ch:
			var tmp binance.MiniTicker
			err := json.Unmarshal(message, &tmp)
			if err == nil && tmp.EventType == "24hrMiniTicker" {
				if found, ok := idHash[tmp.Symbol]; ok {
					price, _ := strconv.ParseFloat(strings.TrimSpace(tmp.Close), 64)
					service.UpdatePrice(domains.CoinPrice{
						Timestamp: tmp.Timestamp,
						Price:     price,
						CoinId:    found,
					})
				}
			}
		}
	}
}

func chunkSlice(items []*domains.Coin, chunkSize int) [][]*domains.Coin {
	if chunkSize <= 0 {
		panic("chunkSize must be positive")
	}

	if len(items) == 0 {
		return [][]*domains.Coin{}
	}

	var chunks [][]*domains.Coin
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunk := items[i:end]
		chunks = append(chunks, chunk)
	}
	return chunks
}
