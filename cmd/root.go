package cmd

import (
	"os"

	ethereum "github.com/piotrostr/gosend/eth"
	"github.com/spf13/cobra"
)

// quantity in eth as per 1 or 0.15
// display a confirmation with usd/pln equivalent
var rootCmd = &cobra.Command{
	Use:   "gosend",
	Short: "send ethereum from command-line",
	Run: func(cmd *cobra.Command, args []string) {
		qty := cmd.Flag("qty").Value.String()
		to := cmd.Flag("to").Value.String()
		// TODO validate here
		eth := ethereum.Eth{}
		eth.Init()
		eth.Send(qty, to)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("qty", "q", "", "quantity in eth as per 1 or 0.15")
	rootCmd.MarkFlagRequired("qty")

	rootCmd.Flags().StringP("to", "t", "", "address to send to")
	rootCmd.MarkFlagRequired("to")
}
