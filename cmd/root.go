package cmd

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"

	ethereum "github.com/piotrostr/gosend/eth"
	"github.com/spf13/cobra"
)

func ask() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if len(s) > 1 {
			fmt.Fprintln(os.Stderr, "Please enter Y or N")
			continue
		}
		if strings.Compare(s, "n") == 0 {
			return false
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
	return true
}

// quantity in eth as per 1 or 0.15
// display a confirmation with usd/pln equivalent
var rootCmd = &cobra.Command{
	Use:   "gosend",
	Short: "send ethereum from command-line",
	Run: func(cmd *cobra.Command, args []string) {
		qty := cmd.Flag("qty").Value.String()
		to := cmd.Flag("to").Value.String()
		// TODO validate here and parse eth to wei
		// verify address as wellas the chainId (add param too)
		// qtyWei := new(big.Int).Mul(big.NewFloat())
		eth := ethereum.Eth{}
		eth.Init()
		fmt.Println("Sending", qty, "to", to)
		fmt.Println("Go for it? [Y/n]")
		if !ask() {
			return
		}

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
