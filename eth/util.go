package eth

import (
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/params"
)

func EthStringToWei(qty string) *big.Int {
	qtyF, err := strconv.ParseFloat(qty, 64)
	if err != nil {
		log.Fatal(err)
	}
	ether := big.NewFloat(params.Ether)
	qtyWei, _ := new(big.Float).Mul(big.NewFloat(qtyF), ether).Int64()
	return big.NewInt(qtyWei)
}
