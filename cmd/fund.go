// Fetches information for index funds.
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

const URL_API_FUND = "https://query2.finance.yahoo.com/v7/finance/options"

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Fetch performances for ETFs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		funds := ParseCommaString(args[0])

		c := make(chan fundInfo)
		var wg sync.WaitGroup

		for _, fund := range funds {
			wg.Add(1)
			go fetchFundInfo(fund, c, &wg)
		}

		// Necessary to avoid deadlock
		// https://stackoverflow.com/a/70877210/1585507
		go func() {
			wg.Wait()
			close(c)
		}()

		for data := range c {
			fmt.Printf(
				"${alignc}%s: %g %s (%.2f %%)\n",
				data.Symbol,
				data.Price,
				data.Currency,
				data.diff(),
			)
		}
	},
}

// Fetch the market information for one fund.
// Returns the info + pushes the info to a channel.
func fetchFundInfo(fund string, c chan<- fundInfo, wg *sync.WaitGroup) fundInfo {
	defer wg.Done()

	responseBody := queryApi(fmt.Sprintf("%s/%s", URL_API_FUND, fund))

	value := gjson.GetBytes(responseBody, "optionChain.result.0.quote")

	var target fundInfo
	if err := json.Unmarshal([]byte(value.Raw), &target); err != nil {
		log.Fatal(err)
	}

	c <- target

	return target
}
