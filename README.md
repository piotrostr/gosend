# gosend

Simple cobra-based cli tool to send ethereum over rinkeby/mainnet (or ganache
at localhost:8545).

Under the hood uses the `github.com/ethereum/go-ethereum` module.

## Usage

Having `INFURA_KEY` and `PRIVATE_KEY` environment variables set, a simple:

```sh
gosend \
  --qty 0.042069 \
  --to $SOME_ADDRESS \
  --chain rinkeby
```

sends 0.042069 ethereum over rinkeby to `$SOME_ADDRESS`.

## License

[MIT](https://github.com/piotrostr/gosend/blob/main/LICENSE)
