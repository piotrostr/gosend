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

func getEnv(variable string) string {
	val, exists := os.LookupEnv(variable)
	if !exists {
		log.Fatalf("env var %s not found", variable)
	}

	return val
}

type Eth struct {
	chainId *big.Int
	Client  *ethclient.Client
	prv     *ecdsa.PrivateKey
	address *common.Address
	balance *big.Int
}

func (eth *Eth) Init(chainName string) {
	var Client *ethclient.Client
	var err error
	switch chainName {
	case "mainnet":
	case "rinkeby":
		rawUrl := "https://%s.infura.io/v3/%s"
		url := fmt.Sprintf(rawUrl, chainName, getEnv("INFURA_KEY"))
		Client, err = ethclient.Dial(url)
		if err != nil {
			log.Fatal(err)
		}
	case "localhost":
		Client, err = ethclient.Dial("http://localhost:8545")
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("chain %s not supported", chainName)
	}
	eth.Client = Client

	chainId, err := Client.NetworkID(ctx)
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
	nonce, err := eth.Client.PendingNonceAt(ctx, *eth.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := eth.Client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	txdata := &types.LegacyTx{
		Nonce:    nonce,
		Gas:      21000,
		GasPrice: gasPrice,
		To:       to,
		Value:    qty,
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

	err = eth.Client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}

	_, isPending, err := eth.Client.TransactionByHash(ctx, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

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

	rec, err := eth.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	pretty, _ := json.MarshalIndent(rec, "", "  ")
	fmt.Printf("Receipt: %+v\n", string(pretty))

	eth.UpdateBalance()
	fmt.Println("PostTxBalance:", eth.balance, "wei")
}
