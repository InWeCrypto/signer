package eth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/goany/bip39"
	"github.com/goany/slf4go"
)

var logger = slf4go.Get("eth")

// Wallet .
type Wallet struct {
	slf4go.Logger
	key *keystore.Key
}

// WalletFromMnemonic create wallet from mnemonic words
func WalletFromMnemonic(mnemonic string) (*Wallet, error) {
	walletMnemonic := OpenWalletMnemonic(mnemonic)

	key, err := walletMnemonic.ToKey()

	if err != nil {
		return nil, err
	}

	return &Wallet{
		Logger: slf4go.Get("wallet"),
		key:    key,
	}, nil
}

// WalletFromPrivateKey create wallet direct from private key
func WalletFromPrivateKey(privatekey []byte) (*Wallet, error) {
	privateKeyECDSA, err := crypto.ToECDSA(privatekey)

	if err != nil {
		return nil, err
	}

	key := keystore.NewKeyForDirectICAP(rand.Reader)

	key.PrivateKey = privateKeyECDSA
	key.Address = crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

	return &Wallet{
		Logger: slf4go.Get("wallet"),
		key:    key,
	}, nil
}

// OpenWallet .
func OpenWallet(wallet []byte, password string) (*Wallet, error) {

	key, err := keystore.DecryptKey(wallet, password)

	runtime.GC()

	if err != nil {
		return nil, err
	}

	return &Wallet{
		Logger: slf4go.Get("wallet"),
		key:    key,
	}, nil
}

// NewWallet create new wallet
func NewWallet() (*Wallet, error) {

	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)

	if err != nil {
		return nil, err
	}

	key := keystore.NewKeyFromECDSA(privateKeyECDSA)

	return &Wallet{
		Logger: slf4go.Get("wallet"),
		key:    key,
	}, err
}

// Encrypt encrypt wallet into json format
func (wallet *Wallet) Encrypt(password string) ([]byte, error) {
	return keystore.EncryptKey(wallet.key, password, keystore.LightScryptN, keystore.LightScryptP)
}

// Mnemonic .
func (wallet *Wallet) Mnemonic() (string, error) {
	mnemonic, err := NewWalletMnemonic(wallet.key)

	if err != nil {
		return "", err
	}

	return mnemonic.String(), nil
}

// TransferCurrency transfer eth currency to special address
func (wallet *Wallet) TransferCurrency(
	nonce uint64,
	gasPrice *big.Int,
	gasLimit *big.Int,
	to string,
	amount string) ([]byte, error) {

	logger.DebugF("try get nonce from server ..")

	var count hexutil.Big

	err := count.UnmarshalText([]byte(amount))

	wallet.Debug(count.String())

	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(strings.Trim(to, " "))

	if addr.Hex() == "0x0000000000000000000000000000000000000000" {
		return nil, fmt.Errorf("bad address fmt:%s", to)
	}

	tx := types.NewTransaction(
		nonce,
		addr,
		count.ToInt(),
		gasLimit,
		gasPrice,
		nil,
	)

	signedTx, err := types.SignTx(tx, types.FrontierSigner{}, wallet.key.PrivateKey)

	if err != nil {
		return nil, nil
	}

	return rlp.EncodeToBytes(signedTx)
}

// Address get wallet address string .
func (wallet *Wallet) Address() string {
	return wallet.key.Address.Hex()
}

// TransferToken .
func (wallet *Wallet) TransferToken(
	nonce uint64,
	gasPrice *big.Int,
	gasLimit *big.Int,
	contract string,
	data []byte) ([]byte, error) {

	var bytes hexutil.Bytes

	bytes.UnmarshalText(data)

	addr := common.HexToAddress(strings.Trim(contract, " "))

	if addr.Hex() == "0x0000000000000000000000000000000000000000" {
		return nil, fmt.Errorf("bad address fmt:%s", contract)
	}

	tx := types.NewTransaction(
		nonce,
		addr,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		bytes,
	)

	signedTx, err := types.SignTx(tx, types.FrontierSigner{}, wallet.key.PrivateKey)

	if err != nil {
		return nil, nil
	}

	return rlp.EncodeToBytes(signedTx)
}

// WalletMnemonic .
type WalletMnemonic struct {
	mnemonic string
}

// OpenWalletMnemonic .
func OpenWalletMnemonic(mnemonic string) *WalletMnemonic {
	return &WalletMnemonic{
		mnemonic: mnemonic,
	}
}

// NewWalletMnemonic .
func NewWalletMnemonic(key *keystore.Key) (*WalletMnemonic, error) {

	privateKey := crypto.FromECDSA(key.PrivateKey)

	mnemonic, err := bip39.NewMnemonic(privateKey)

	return &WalletMnemonic{
		mnemonic: mnemonic,
	}, err
}

// ToKey convert mnemonic to keystore.Key
func (wm *WalletMnemonic) ToKey() (*keystore.Key, error) {

	privateKey, err := bip39.MnemonicToByteArray(wm.mnemonic)

	if err != nil {
		return nil, err
	}

	privateKey = privateKey[1 : len(privateKey)-1]

	privateKeyECDSA, err := crypto.ToECDSA(privateKey)

	if err != nil {
		return nil, err
	}

	key := keystore.NewKeyForDirectICAP(rand.Reader)

	key.PrivateKey = privateKeyECDSA
	key.Address = crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

	return key, err
}

func (wm *WalletMnemonic) String() string {
	return wm.mnemonic
}
