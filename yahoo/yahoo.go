package yahoo

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"

	"local/ticker_fetcher/utils"

	"github.com/tidwall/gjson"
)

const URL_API = "https://query2.finance.yahoo.com/v7/finance/options"

// Fetch the market information for one fund.
// Returns the info + pushes the info to a channel.
func FetchInfoFromYahoo(fund string, c chan<- YahooInfo, wg *sync.WaitGroup) YahooInfo {
	defer wg.Done()

	responseBody := utils.QueryApi(fmt.Sprintf("%s/%s", URL_API, fund))

	value := gjson.GetBytes(responseBody, "optionChain.result.0.quote")

	var target YahooInfo
	if err := json.Unmarshal([]byte(value.Raw), &target); err != nil {
		log.Fatal(err)
	}

	c <- target

	return target
}

type YahooInfo struct {
	Symbol        string  `json:"symbol"`
	ExchangeName  string  `json:"fullExchangeName"`
	Price         float64 `json:"regularMarketPrice"`
	PreviousClose float64 `json:"regularMarketPreviousClose"`
	Currency      string  `json:"currency"`
	FromCurrency  string  `json:"fromCurrency"`
	ShortName     string  `json:"shortName"`
}

func (fi *YahooInfo) Diff() float64 {
	return (fi.Price - fi.PreviousClose) / fi.PreviousClose * 100
}

func ProcessFromYahoo(symbols string, display func(YahooInfo)) {

	tickers := utils.ParseCommaString(symbols)

	c := make(chan YahooInfo)
	var wg sync.WaitGroup

	for _, tick := range tickers {
		wg.Add(1)
		go FetchInfoFromYahoo(tick, c, &wg)
	}

	// Necessary to avoid deadlock
	// https://stackoverflow.com/a/70877210/1585507
	go func() {
		wg.Wait()
		close(c)
	}()

	var results []YahooInfo

	for data := range c {
		results = append(results, data)
	}

	// Ascending sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Symbol < results[j].Symbol
	})

	for _, data := range results {
		display(data)
	}

}
