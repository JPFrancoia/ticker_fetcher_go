// Fetches information for index funds.
// Uses the Yahoo API.

package cmd

import (
	"fmt"

	"local/ticker_fetcher/utils"
	"local/ticker_fetcher/yahoo"

	"github.com/spf13/cobra"
)

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Fetch performances for ETFs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		display := func(data yahoo.YahooInfo) {
			if data.Currency == "USD" {
				fmt.Printf(
					"${alignc}%s = %.2f (%.2f %%)\n\n",
					data.Symbol,
					data.Price,
					data.Diff(),
				)
			} else {
				price, currency := utils.ConvertDecimal(data.Price, data.Currency)
				fmt.Printf(
					"${alignc}%s = %.2f %s (%.2f %%)\n\n",
					data.Symbol,
					price,
					currency,
					data.Diff(),
				)
			}
		}

		yahoo.ProcessFromYahoo(args[0], display)
	},
}
