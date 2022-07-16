package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
		// parse flags
		qty := cmd.Flag("qty").Value.String()
		to := cmd.Flag("to").Value.String()
		chain := cmd.Flag("chain").Value.String()

		// initialize eth client
		eth := ethereum.Eth{}
		eth.Init(chain)

		// parse and check address
		addr := common.HexToAddress(to)
		bytecode, err := eth.Client.CodeAt(context.Background(), addr, nil)
		if err != nil {
			log.Fatal(err)
		}
		isContract := len(bytecode) > 0

		// print the tx data
		rawMsg := "Sending %s to %s (is contract: %v) on %s\n"
		fmt.Printf(rawMsg, qty, to, isContract, chain)

		// double check with user
		fmt.Println("Go for it? [Y/n]")
		if !ask() {
			fmt.Println("Aborted")
			return
		}

		// send the tx
		eth.Send(&addr, ethereum.EthStringToWei(qty))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("qty", "", "quantity in eth as per 1 or 0.15")
	rootCmd.MarkFlagRequired("qty")

	rootCmd.Flags().String("to", "", "address to send to")
	rootCmd.MarkFlagRequired("to")

	rootCmd.Flags().String(
		"chain",
		"localhost",
		"chain to send on (mainnet or rinkeby, defaults to localhost)",
	)
}
