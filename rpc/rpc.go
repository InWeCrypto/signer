package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/dghubble/sling"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/goany/btc"
	"github.com/goany/slf4go"
)

// Client the fbg RESTful api wrapper client
type Client struct {
	slf4go.Logger
	client *sling.Sling
}

type nonceRequest struct {
	Address string `json:"address"`
}

type nonceResponse struct {
	Count string `json:"count"`
}

type gasPriceResp struct {
	GasPrice string `json:"gasPrice"`
}

type balanceResp struct {
	Value string `json:"value"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type tokenRequest struct {
	Address  string `json:"address"`
	Contract string `json:"contract"`
}

type tokenNameRequest struct {
	Contract string `json:"contract"`
}

type tokenTransferRequest struct {
	Contract string `json:"contract"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

// TokenTransfer .
type TokenTransfer struct {
	Contract string `json:"contract"`
	Data     string `json:"data"`
}

type rawTx struct {
	Data string `json:"data"`
}

type rawTxResp struct {
	TxHash string `json:"txHash"`
}

type btcRawTx struct {
	Data string `json:"rawtx"`
}

type btcFeeRate struct {
	Blocks int64 `json:"nbBlocks"`
}

type btcFeeRateResp struct {
	Fee string `json:"satoshi"`
}

// NewClient cretae new RESTful api client
func NewClient(base string) *Client {
	return &Client{
		Logger: slf4go.Get("ethclient"),
		client: sling.New().Base(base),
	}
}

// GetNonce get transaction nonce for input address
func (client *Client) GetNonce(address string) (*big.Int, error) {

	response := &nonceResponse{}
	errorResp := &errorResponse{}

	_, err := client.client.Post("/eth/getTransactionCount").BodyJSON(&nonceRequest{
		Address: address,
	}).Receive(response, errorResp)

	if err != nil {
		return nil, err
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response :%s", string(jsontext))

	var count hexutil.Big

	count.UnmarshalText([]byte(response.Count))

	return count.ToInt(), nil
}

// GetGasPrice get transaction gace price
func (client *Client) GetGasPrice() (*big.Int, error) {

	response := &gasPriceResp{}

	resp, err := client.client.Get("/eth/getGasPrice").ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, errors.New(resp.Status)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	var count hexutil.Big

	count.UnmarshalText([]byte(response.GasPrice))

	return count.ToInt(), nil
}

// GetBalance get account balance
func (client *Client) GetBalance(address string) (*big.Int, error) {

	response := &balanceResp{}

	resp, err := client.client.Post("/eth/getBalance").BodyJSON(&nonceRequest{
		Address: address,
	}).ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, errors.New(resp.Status)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	var count hexutil.Big

	count.UnmarshalText([]byte(response.Value))

	return count.ToInt(), nil
}

// GetTokenBalance get account balance
func (client *Client) GetTokenBalance(address string, token string) (*big.Int, error) {

	response := &balanceResp{}

	resp, err := client.client.Post("/eth/tokens/balanceOf").BodyJSON(&tokenRequest{
		Address:  address,
		Contract: token,
	}).ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, errors.New(resp.Status)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	var count hexutil.Big

	var numbersRegExp = regexp.MustCompile("0x[0]*")

	response.Value = numbersRegExp.ReplaceAllString(response.Value, "0x")

	if response.Value == "0x" {
		response.Value = "0x0"
	}

	err = count.UnmarshalText([]byte(response.Value))

	return count.ToInt(), err
}

// GetTokenSupply get token supply
func (client *Client) GetTokenSupply(token string) (*big.Int, error) {

	response := &balanceResp{}

	resp, err := client.client.Post("/eth/tokens/totalSupply").BodyJSON(&tokenNameRequest{
		Contract: token,
	}).ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, errors.New(resp.Status)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	var count hexutil.Big

	var numbersRegExp = regexp.MustCompile("0x[0]*")

	response.Value = numbersRegExp.ReplaceAllString(response.Value, "0x")

	err = count.UnmarshalText([]byte(response.Value))

	return count.ToInt(), err
}

// GetTokenTransfer get token supply
func (client *Client) GetTokenTransfer(token string, to string, value string) (*TokenTransfer, error) {

	response := &TokenTransfer{}

	resp, err := client.client.Post("/eth/tokens/transferABI").BodyJSON(&tokenTransferRequest{
		Contract: token,
		To:       to,
		Value:    value,
	}).ReceiveSuccess(response)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, errors.New(resp.Status)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	return response, nil
}

// CommitTx commit raw tx
func (client *Client) CommitTx(data string) (string, error) {

	response := &rawTxResp{}
	errorResp := &errorResponse{}

	resp, err := client.client.Post("/eth/sendRawTransaction").BodyJSON(&rawTx{
		Data: data,
	}).Receive(response, errorResp)

	if err != nil {
		return "", err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return "", fmt.Errorf("(%s) %s", resp.Status, errorResp.Message)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	return response.TxHash, nil
}

// UTXO get address's utxo list
func (client *Client) UTXO(address string) ([]btc.UTXO, error) {
	var response []btc.UTXO
	errorResp := &errorResponse{}

	resp, err := client.client.Post("/btc/getUtxo").BodyJSON(&nonceRequest{
		Address: address,
	}).Receive(&response, errorResp)

	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return nil, fmt.Errorf("(%s) %s", resp.Status, errorResp.Message)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	return response, nil
}

// BTCSend .
func (client *Client) BTCSend(rawTx string) error {
	errorResp := &errorResponse{}

	response := new(interface{})

	resp, err := client.client.Post("/btc/send").BodyJSON(&btcRawTx{
		Data: rawTx,
	}).Receive(response, errorResp)

	if err != nil {
		return err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return fmt.Errorf("(%s) %s", resp.Status, errorResp.Message)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	return nil
}

// BTCFeeRate .
func (client *Client) BTCFeeRate(nbBlocks int64) (int64, error) {
	errorResp := &errorResponse{}

	response := &btcFeeRateResp{}

	resp, err := client.client.Post("/btc/estimatefee").BodyJSON(&btcFeeRate{
		Blocks: nbBlocks,
	}).Receive(response, errorResp)

	if err != nil {
		return 0, err
	}

	if code := resp.StatusCode; 200 < code || code > 299 {
		return 0, fmt.Errorf("(%s) %s", resp.Status, errorResp.Message)
	}

	jsontext, _ := json.Marshal(response)

	client.DebugF("response(%d) :%s", resp.StatusCode, string(jsontext))

	var count hexutil.Big

	var numbersRegExp = regexp.MustCompile("0x[0]*")

	response.Fee = numbersRegExp.ReplaceAllString(response.Fee, "0x")

	err = count.UnmarshalText([]byte(response.Fee))

	return count.ToInt().Int64(), err
}
