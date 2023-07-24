// Fetches information for index funds.
// Uses the Yahoo API.

package cmd

import (
	"fmt"

	"local/ticker_fetcher/yahoo"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Fetch indexes performances",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		display := func(data yahoo.YahooInfo) {
			if data.Currency == "USD" {
				fmt.Printf(
					"${alignc}%s = %.2f (%.2f %%)\n\n",
					data.ExchangeName,
					data.Price,
					data.Diff(),
				)
			} else {
				fmt.Printf(
					"${alignc}%s = %.2f %s (%.2f %%)\n\n",
					data.ExchangeName,
					data.Price,
					data.Currency,
					data.Diff(),
				)

			}
		}

		yahoo.ProcessFromYahoo(args[0], display)
	},
}
