package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// quantity in eth as per 1 or 0.15
// display a confirmation with usd/pln equivalent
var rootCmd = &cobra.Command{
	Use:   "gosend",
	Short: "send ethereum from command-line",
	Long:  "TODO: Long description",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
