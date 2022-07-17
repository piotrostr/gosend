# gosend

Simple cobra-based cli tool to send ethereum over rinkeby/mainnet (or ganache
at localhost:8545).

Under the hood uses the `github.com/ethereum/go-ethereum` module.

## Usage

Dead simple:

```sh
gosend \
  --qty 0.042069 \
  --to $some_addr \
  --chain rinkeby
```

sends 0.042069 ethereum over rinkeby to `$some_addr`.

## License

[MIT](https://github.com/piotrostr/gosend/blob/main/LICENSE)
