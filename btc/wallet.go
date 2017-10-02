package btc

import (
	"io"

	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/goany/slf4go"
)

// NetType btc net type
type NetType string

// btc net config name
const (
	NetTypeTestNet3 NetType = "testnet3"
	NetTypeRegTest  NetType = "testnet"
	NetTypeMainNet  NetType = "mainnet"
)

// Wallet BTC wallet
type Wallet struct {
	slf4go.Logger
	privateKey *btcec.PrivateKey
	publicKey  *btcec.PublicKey
	Address    btcutil.Address
	net        *chaincfg.Params
	compressed bool
}

// NewWallet create wallet from private key
func NewWallet(privateKeyString string, chainname NetType) (*Wallet, error) {

	logger := slf4go.Get("BTCWallet")

	compressed := false

	switch privateKeyString[0] {
	case 'L', 'K', 'c':
		compressed = true
	}

	bytes, _, _ := base58.CheckDecode(privateKeyString)
	priv, pub := btcec.PrivKeyFromBytes(btcec.S256(), bytes)
	// bytes, err := hex.DecodeString(privateKeyString)

	// priv, pub := btcec.PrivKeyFromBytes(btcec.S256(), bytes)

	wallet := &Wallet{
		Logger:     logger,
		privateKey: priv,
		publicKey:  pub,
		compressed: compressed,
	}

	switch chainname {
	case NetTypeTestNet3:
		wallet.net = &chaincfg.TestNet3Params
	case NetTypeRegTest:
		wallet.net = &chaincfg.RegressionNetParams
	case NetTypeMainNet:
		wallet.net = &chaincfg.MainNetParams
	default:
		return nil, fmt.Errorf("unknown btc net :%s", chainname)
	}

	address, err := btcutil.NewAddressPubKeyHash(
		btcutil.Hash160(
			wallet.publicKey.SerializeUncompressed(),
		),
		wallet.net,
	)

	if err != nil {
		return nil, err
	}

	wallet.Address = address

	return wallet, nil
}

// Pay pay btc to address
func (wallet *Wallet) Pay(
	inputs []UTXO,
	to string,
	amount btcutil.Amount,
	feeRate btcutil.Amount,
	writer io.Writer) error {

	addr, err := btcutil.DecodeAddress(to, wallet.net)

	wallet.Debug("???????", addr.EncodeAddress())

	if err != nil {
		return err
	}

	tx, err := txFrom(inputs).
		to(addr, amount).
		change(wallet.Address).
		feeRate(feeRate).
		sign(wallet.privateKey, wallet.compressed).
		done()

	if err != nil {
		return err
	}

	return tx.Serialize(writer)
}
