package ethmobile

import (
	"bytes"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/goany/slf4go"
)

var logger = slf4go.Get("test")

// TestTrans .
func TestTrans() {
	json, _ := ioutil.ReadFile("testdata/keystore.json")

	key, _ := keystore.DecryptKey(json, "test")

	logger.Debug(key.Address.Hex())

	tx := types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), big.NewInt(0), big.NewInt(0),
		nil,
	)

	signedTx, _ := types.SignTx(tx, types.HomesteadSigner{}, key.PrivateKey)

	writer := &bytes.Buffer{}

	signedTx.EncodeRLP(writer)

	logger.DebugF("%x", writer.Bytes())
}
