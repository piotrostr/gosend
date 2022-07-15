package main

// import "github.com/piotrostr/gosend/cmd"

import (
	ethereum "github.com/piotrostr/gosend/eth"
)

func main() {
	eth := ethereum.Eth{}
	eth.Init()
	eth.Send()
}
