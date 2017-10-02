package rpc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/goany/btc"
	"github.com/goany/eth"
)

// var client = NewClient("http://192.168.2.254:8888")

var client = NewClient("http://120.77.208.222:8888")

var wallet *btc.Wallet

func init() {
	var err error
	wallet, err = btc.NewWallet("93RdvTuHcJfMbihdZVqjSHvbEXpUzVEgFMxkwqK6jRqyDXBJf6w", btc.NetTypeRegTest)

	if err != nil {
		panic(err)
	}
}

func TestCreate(t *testing.T) {
	w, _ := eth.NewWallet()

	w.Debug(w.Address())
}

func TestBTC(t *testing.T) {

	utxos, err := client.UTXO(wallet.Address.EncodeAddress())

	if err != nil {
		t.Fatal(err)
	}

	// feeRate, err := client.BTCFeeRate(1)

	// if err != nil {
	// 	t.Fatal(err)
	// }

	var buff bytes.Buffer

	err = wallet.Pay(
		utxos, "n3aVWR1nkPJ1bq3DTHs8grsqc6un4MNVR1",
		5000000,
		200,
		&buff)

	if err != nil {
		t.Fatal(err)
	}

	rawtx := hex.EncodeToString(buff.Bytes())

	client.Debug("===========", rawtx)

	err = client.BTCSend(rawtx)

	if err != nil {
		t.Fatal(err)
	}
}

func TestETH(t *testing.T) {
	json, err := ioutil.ReadFile("../testdata/keystore.json")

	if err != nil {
		t.Fatal(err)
	}

	client := NewClient("http://120.77.208.222:8888")

	wallet, err := eth.OpenWallet(json, "test")

	if err != nil {
		t.Fatal(err)
	}

	balance, err := client.GetBalance(wallet.Address())

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug("balance :", balance)

	balance, err = client.GetTokenBalance(wallet.Address(), "0x07a1e67129305b8a99c86a481681032c68f9f0f9")

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug("token balance :", balance)

	// transfer token

	tokenTrans, err := client.GetTokenTransfer(
		"0x07a1e67129305b8a99c86a481681032c68f9f0f9",
		"0xa81e19b1b13981225fd2fae0bea7d1dc4017cb1b",
		"0x200000000")

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug("token transfer :", tokenTrans.Contract, "###", tokenTrans.Data)

	nonce, err := client.GetNonce(wallet.Address())

	if err != nil {
		t.Fatal(err)
	}

	gasPrice, err := client.GetGasPrice()

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug("gas prices :", gasPrice)

	tx, err := wallet.TransferToken(
		nonce.Uint64(),
		gasPrice,
		"0x07a1e67129305b8a99c86a481681032c68f9f0f9",
		[]byte(tokenTrans.Data))

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug(fmt.Sprintf("0x%s", hex.EncodeToString(tx)))

	// txID, err := client.CommitTx(fmt.Sprintf("0x%s", hex.EncodeToString(tx)))

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wallet.Debug("txId", txID)

	// nonce, err = client.GetNonce(wallet.Address())

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// tx, err = wallet.TransferCurrency(
	// 	nonce.Uint64(),
	// 	gasPrice,
	// 	"0xa81e19b1b13981225fd2fae0bea7d1dc4017cb1b",
	// 	"0xe8d4a5100000000")

	// txID, err = client.CommitTx(fmt.Sprintf("0x%s", hex.EncodeToString(tx)))

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// wallet.Debug("txId", txID)
}

func TestNewWallet(t *testing.T) {
	wallet, err := eth.NewWallet()

	if err != nil {
		t.Fatal(err)
	}

	data, err := wallet.Encrypt("test")

	if err != nil {
		t.Fatal(err)
	}

	wallet.Debug(string(data))
}

func TestMnemonic(t *testing.T) {

	for i := 0; i < 10; i++ {
		wallet, _ := eth.NewWallet()

		mn, err := wallet.Mnemonic()

		if err != nil {
			t.Fatal(err)
		}

		newWallet, err := eth.WalletFromMnemonic(mn)

		if err != nil {
			t.Fatal(err)
		}

		newmn, err := newWallet.Mnemonic()

		if err != nil || newmn != mn {
			t.Fatalf("check WalletFromMnemonic -- failed")
		}

		wallet.Debug(mn)

		key, err := eth.OpenWalletMnemonic(mn).ToKey()

		wallet.Debug(hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))

		_, err = eth.WalletFromPrivateKey(crypto.FromECDSA(key.PrivateKey))

		if err != nil {
			t.Fatal(err)
		}
	}

}
