package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ctx = context.TODO()

func getEnv(variable string) string {
	val, exists := os.LookupEnv(variable)
	if !exists {
		log.Fatal("env var not found")
	}

	return val
}

type Eth struct {
	chainId *big.Int
	client  *ethclient.Client
	prv     *ecdsa.PrivateKey
	address *common.Address
	balance *big.Int
}

func (eth *Eth) Init() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	eth.client = client

	chainId, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	eth.chainId = chainId

	prv, err := crypto.HexToECDSA(getEnv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	eth.prv = prv

	address := crypto.PubkeyToAddress(prv.PublicKey)
	eth.address = &address

	eth.UpdateBalance()

	fmt.Println("ChainID:", eth.chainId)
	fmt.Println("Address:", eth.address.Hex())
	fmt.Println("Balance:", eth.balance)
}

func (eth *Eth) UpdateBalance() {
	balance, err := eth.client.BalanceAt(ctx, *eth.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	eth.balance = balance
}

func (eth *Eth) Send() {
	nonce, err := eth.client.PendingNonceAt(ctx, *eth.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := eth.client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(100000000000000)
	txdata := &types.LegacyTx{
		Nonce:    nonce,
		Gas:      21000,
		GasPrice: gasPrice,
		To:       eth.address,
		Value:    value,
		Data:     []byte{},
	}
	tx, _ := types.SignNewTx(
		eth.prv,
		types.LatestSignerForChainID(eth.chainId),
		txdata,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = eth.client.SendTransaction(ctx, tx)
	if err != nil {
		log.Println(err)
	}

	eth.UpdateBalance()
	fmt.Println("PostTxBalance:", eth.balance)
}
