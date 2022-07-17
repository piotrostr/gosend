package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ctx = context.TODO()

// get an environment variable, throw err if does not exist
func getEnv(variable string) string {
	val, exists := os.LookupEnv(variable)
	if !exists {
		log.Fatalf("env var %s not found", variable)
	}

	return val
}

type Eth struct {
	Client  *ethclient.Client
	chainId *big.Int
	prv     *ecdsa.PrivateKey
	address *common.Address
	balance *big.Int
}

func (eth *Eth) Init(chainName string) {
	// get client depending on chainName param
	var client *ethclient.Client
	var err error
	switch chainName {
	case "mainnet", "rinkeby":
		rawUrl := "https://%s.infura.io/v3/%s"
		url := fmt.Sprintf(rawUrl, chainName, getEnv("INFURA_KEY"))
		client, err = ethclient.Dial(url)
		if err != nil {
			log.Fatal(err)
		}
	case "localhost":
		client, err = ethclient.Dial("http://localhost:8545")
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("chain %s not supported", chainName)
	}
	eth.Client = client

	// get chainId
	chainId, err := eth.Client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	eth.chainId = chainId

	// get private key env var and ecdsa-ify it
	prv, err := crypto.HexToECDSA(getEnv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	eth.prv = prv

	// get address from pub key
	address := crypto.PubkeyToAddress(prv.PublicKey)
	eth.address = &address

	eth.UpdateBalance()

	// print initialized eth struct vals
	fmt.Println("ChainID:", eth.chainId)
	fmt.Println("Address:", eth.address.Hex())
	fmt.Println("Balance:", eth.balance, "wei")
}

func (eth *Eth) UpdateBalance() {
	balance, err := eth.Client.BalanceAt(ctx, *eth.address, nil)
	if err != nil {
		log.Fatal(err)
	}
	eth.balance = balance
}

func (eth *Eth) Send(to *common.Address, qty *big.Int) {
	// get nonce
	nonce, err := eth.Client.PendingNonceAt(ctx, *eth.address)
	if err != nil {
		log.Fatal(err)
	}

	// get gasPrice
	gasPrice, err := eth.Client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// make the tx struct
	txdata := &types.LegacyTx{
		Nonce:    nonce,
		Gas:      21000,
		GasPrice: gasPrice,
		To:       to,
		Value:    qty,
		Data:     []byte{},
	}

	// sign the tx
	tx, _ := types.SignNewTx(
		eth.prv,
		types.LatestSignerForChainID(eth.chainId),
		txdata,
	)
	if err != nil {
		log.Fatal(err)
	}

	// send the tx
	err = eth.Client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}

	// cooldown (wait for the tx to be added to the block)
	time.Sleep(time.Second * 1)

	// check if tx pending
	_, isPending, err := eth.Client.TransactionByHash(ctx, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// if tx is pending, let user know and wait for it to be mined
	if isPending {
		fmt.Println("Transaction is pending")
	}
	for isPending {
		_, isPending, err = eth.Client.TransactionByHash(ctx, tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(".")
		time.Sleep(1 * time.Second) // add cooldown var
	}

	// get the tx receipt
	rec, err := eth.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// prettify and print to stdout
	pretty, _ := json.MarshalIndent(rec, "", "  ")
	fmt.Printf("Receipt: %+v\n", string(pretty))

	eth.UpdateBalance()
}
