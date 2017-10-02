package btc

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/goany/slf4go"
)

// UTXO .
type UTXO struct {
	Address       string  `json:"address"`
	TxID          string  `json:"txid"`
	VOut          uint32  `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      float64 `json:"satoshis"`
	Height        float64 `json:"height"`
	Confirmations float64 `json:"confirmations"`
}

// Transaction btc transaction object
type transaction struct {
	slf4go.Logger
	inputs     []UTXO // input utxos
	payto      btcutil.Address
	amount     btcutil.Amount
	paychange  btcutil.Address
	payFeeRate btcutil.Amount
	privatekey *btcec.PrivateKey
	txIn       map[*wire.TxIn]UTXO
	compressed bool
}

func txFrom(inputs []UTXO) *transaction {
	return &transaction{
		inputs: inputs,
		Logger: slf4go.Get("tx"),
		txIn:   make(map[*wire.TxIn]UTXO),
	}
}

func (trans *transaction) to(addr btcutil.Address, amount btcutil.Amount) *transaction {
	trans.payto = addr
	trans.amount = amount
	return trans
}

func (trans *transaction) change(addr btcutil.Address) *transaction {
	trans.paychange = addr
	return trans
}

func (trans *transaction) feeRate(fee btcutil.Amount) *transaction {
	trans.payFeeRate = fee
	return trans
}

func (trans *transaction) sign(privatekey *btcec.PrivateKey, compressed bool) *transaction {
	trans.privatekey = privatekey
	trans.compressed = compressed
	return trans
}

func (trans *transaction) done() (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)

	trans.Debug(trans.payto.EncodeAddress(), hex.EncodeToString(trans.payto.ScriptAddress()))

	addrScript, err := txscript.PayToAddrScript(trans.payto)

	if err != nil {
		return nil, err
	}

	tx.AddTxOut(wire.NewTxOut(int64(trans.amount), addrScript))

	err = trans.calcChange(tx)

	if err != nil {
		return nil, err
	}

	for _, txin := range tx.TxIn {
		utxo := trans.txIn[txin]

		pkScript, err := hex.DecodeString(utxo.ScriptPubKey)

		if err != nil {
			return nil, err
		}

		sigScript, err := txscript.SignatureScript(tx, len(tx.TxIn)-1, pkScript, txscript.SigHashNone, trans.privatekey, trans.compressed)

		if err != nil {
			return nil, err
		}

		txin.SignatureScript = sigScript
	}

	return tx, nil
}

func (trans *transaction) calcChange(tx *wire.MsgTx) error {
	const (
		// spendSize is the largest number of bytes of a sigScript
		// which spends a p2pkh output: OP_DATA_73 <sig> OP_DATA_33 <pubkey>
		spendSize = 1 + 73 + 1 + 33
	)

	var (
		amtSelected btcutil.Amount
		txSize      int
	)

	for _, utxo := range trans.inputs {

		amtSelected += btcutil.Amount(utxo.Satoshis)

		hash, err := chainhash.NewHashFromStr(utxo.TxID)

		if err != nil {
			return err
		}

		outPoint := wire.OutPoint{
			Hash:  *hash,
			Index: utxo.VOut,
		}

		trans.Debug(trans.privatekey)

		txin := wire.NewTxIn(&outPoint, nil)

		tx.AddTxIn(txin)

		trans.txIn[txin] = utxo

		txSize = tx.SerializeSize() + spendSize*len(tx.TxIn)

		reqFee := btcutil.Amount(txSize * int(trans.payFeeRate))
		if amtSelected-reqFee < trans.amount {
			continue
		}

		changeVal := amtSelected - trans.amount - reqFee
		if changeVal > 0 {
			pkScript, err := txscript.PayToAddrScript(trans.paychange)
			if err != nil {
				return err
			}
			changeOutput := &wire.TxOut{
				Value:    int64(changeVal),
				PkScript: pkScript,
			}
			tx.AddTxOut(changeOutput)
		}

		trans.Debug("reqFee: ", reqFee, " change: ", changeVal)

		return nil
	}

	return fmt.Errorf("not enough funds for coin selection")
}
