package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fundCmd)
	rootCmd.AddCommand(stockCmd)
	rootCmd.AddCommand(cryptoCmd)
	rootCmd.AddCommand(fiatCmd)
}

var rootCmd = &cobra.Command{
  Use:   "ticker_fetcher",
  Short: "Fetch financial information.",
  Long: "Fetch financial information like stock tickers, funds, cryptos, etc.",
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
