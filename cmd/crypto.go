// Fetches information for stocks.
// Uses the Yahoo API.

package cmd

import (
	"fmt"
	"local/ticker_fetcher/yahoo"
	"strings"

	"github.com/spf13/cobra"
)

var cryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Fetch exchange rates for cryptos.",
	Long:  "Fetch exchange rates for cryptos. Use Yahoo symbols. E.g BTC to USD is BTC-USD.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		display := func(data yahoo.YahooInfo) {

			ticker := strings.Split(data.Symbol, "-")[0]

			if data.Currency == "USD" {
				fmt.Printf(
					"${alignc}1 %s = %.2f (%.2f %%)\n\n",
					ticker,
					data.Price,
					data.Diff(),
				)
			} else {
				fmt.Printf(
					"${alignc}1 %s = %.2f %s (%.2f %%)\n\n",
					ticker,
					data.Price,
					data.Currency,
					data.Diff(),
				)
			}
		}

		yahoo.ProcessFromYahoo(args[0], display)
	},
}
