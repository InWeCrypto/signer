package unichain

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/goany/eth"
)

// ETHWallet eth wallet facade
type ETHWallet struct {
	impl *eth.Wallet
}

// ETHWalletFromMnemonic create eth wallet from mnemonic
func ETHWalletFromMnemonic(mnemonic string) (*ETHWallet, error) {

	wallet, err := eth.WalletFromMnemonic(mnemonic)

	if err != nil {
		return nil, err
	}

	return &ETHWallet{
		impl: wallet,
	}, nil
}

// OpenETHWallet create eth wallet from json metadata
func OpenETHWallet(json []byte, password string) (*ETHWallet, error) {

	wallet, err := eth.OpenWallet(json, password)

	return &ETHWallet{
		impl: wallet,
	}, err
}

// ETHWalletFromPrivateKey create eth wallet from private key
func ETHWalletFromPrivateKey(privatekey string) (*ETHWallet, error) {

	pk, err := hex.DecodeString(privatekey)

	if err != nil {
		return nil, err
	}

	wallet, err := eth.WalletFromPrivateKey(pk)

	return &ETHWallet{
		impl: wallet,
	}, err
}

// NewETHWallet create new wallet with new pk pair
func NewETHWallet() (*ETHWallet, error) {
	wallet, err := eth.NewWallet()

	if err != nil {
		return nil, err
	}

	return &ETHWallet{
		impl: wallet,
	}, nil
}

// Encrypt package wallet as json format
func (wallet *ETHWallet) Encrypt(password string) (data []byte, err error) {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("c")
		if e := recover(); err != nil {
			err = fmt.Errorf("%s", e)
		}
	}()

	return wallet.impl.Encrypt(password)
}

// Mnemonic generate wallet's mnemonic words
func (wallet *ETHWallet) Mnemonic() (string, error) {
	return wallet.impl.Mnemonic()
}

// TransferCurrency transfer currency to special account address
func (wallet *ETHWallet) TransferCurrency(
	nonceString string,
	gasPriceString string,
	gasLimitString string,
	to string,
	amountString string) ([]byte, error) {

	var nonce hexutil.Big

	err := nonce.UnmarshalText([]byte(nonceString))

	if err != nil {
		return nil, err
	}

	var gasPrice hexutil.Big

	err = gasPrice.UnmarshalText([]byte(gasPriceString))

	if err != nil {
		return nil, err
	}

	var gasLimit hexutil.Big

	err = gasLimit.UnmarshalText([]byte(gasLimitString))

	if err != nil {
		return nil, err
	}

	return wallet.impl.TransferCurrency(nonce.ToInt().Uint64(), gasPrice.ToInt(), gasLimit.ToInt(), to, amountString)
}

// TransferToken transfer token to special account address
func (wallet *ETHWallet) TransferToken(
	nonceString string,
	gasPriceString string,
	gasLimitString string,
	contract string,
	data []byte) ([]byte, error) {

	var nonce hexutil.Big

	err := nonce.UnmarshalText([]byte(nonceString))

	if err != nil {
		return nil, err
	}

	var gasPrice hexutil.Big

	err = gasPrice.UnmarshalText([]byte(gasPriceString))

	if err != nil {
		return nil, err
	}

	var gasLimit hexutil.Big

	err = gasLimit.UnmarshalText([]byte(gasLimitString))

	if err != nil {
		return nil, err
	}

	return wallet.impl.TransferToken(nonce.ToInt().Uint64(), gasPrice.ToInt(), gasLimit.ToInt(), contract, data)
}

// Address get wallet address
func (wallet *ETHWallet) Address() string {
	return wallet.impl.Address()
}
