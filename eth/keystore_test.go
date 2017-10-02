package eth

import (
	"bytes"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestCreateKeyStore(t *testing.T) {
	json, err := ioutil.ReadFile("testdata/keystore.json")

	if err != nil {
		t.Fatal(err)
	}

	key, err := keystore.DecryptKey(json, "test")

	if err != nil {
		t.Fatal(err)
	}

	logger.Debug(key.Address.Hex())

	tx := types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), big.NewInt(0), big.NewInt(0),
		nil,
	)

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, key.PrivateKey)

	if err != nil {
		t.Fatal(err)
	}

	writer := &bytes.Buffer{}

	err = signedTx.EncodeRLP(writer)

	if err != nil {
		t.Fatal(err)
	}

	logger.DebugF("%x", writer.Bytes())
}
