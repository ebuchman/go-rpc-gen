// File generated by github.com/ebuchman/rpc-gen

package rpc

import (
	"fmt"
	"github.com/ebuchman/go-rpc-gen/example/core"
	"github.com/tendermint/tendermint/account"
	"github.com/tendermint/tendermint/binary"
	"github.com/tendermint/tendermint/rpc"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"net/http"
)

type Client interface {
	BlockchainInfo(minHeight uint, maxHeight uint) (*core.ResponseBlockchainInfo, error)
	BroadcastTx(tx types.Tx) (*core.ResponseBroadcastTx, error)
	GenPrivAccount() (*core.ResponseGenPrivAccount, error)
	GetAccount(address []byte) (*core.ResponseGetAccount, error)
	GetBlock(height uint) (*core.ResponseGetBlock, error)
	ListAccounts() (*core.ResponseListAccounts, error)
	ListValidators() (*core.ResponseListValidators, error)
	NetInfo() (*core.ResponseNetInfo, error)
	SignTx(tx types.Tx, privAccounts []*account.PrivAccount) (*core.ResponseSignTx, error)
	Status() (*core.ResponseStatus, error)
}

func (c *ClientHTTP) BlockchainInfo(minHeight uint, maxHeight uint) (*core.ResponseBlockchainInfo, error) {
	values, err := argsToURLValues([]string{"minHeight", "maxHeight"}, minHeight, maxHeight)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"blockchain_info", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseBlockchainInfo
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) BroadcastTx(tx types.Tx) (*core.ResponseBroadcastTx, error) {
	values, err := argsToURLValues([]string{"tx"}, tx)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"broadcast_tx", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseBroadcastTx
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) GenPrivAccount() (*core.ResponseGenPrivAccount, error) {
	values, err := argsToURLValues(nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"gen_priv_account", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGenPrivAccount
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) GetAccount(address []byte) (*core.ResponseGetAccount, error) {
	values, err := argsToURLValues([]string{"address"}, address)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"get_account", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGetAccount
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) GetBlock(height uint) (*core.ResponseGetBlock, error) {
	values, err := argsToURLValues([]string{"height"}, height)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"get_block", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGetBlock
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) ListAccounts() (*core.ResponseListAccounts, error) {
	values, err := argsToURLValues(nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"list_accounts", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseListAccounts
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) ListValidators() (*core.ResponseListValidators, error) {
	values, err := argsToURLValues(nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"list_validators", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseListValidators
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) NetInfo() (*core.ResponseNetInfo, error) {
	values, err := argsToURLValues(nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"net_info", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseNetInfo
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) SignTx(tx types.Tx, privAccounts []*account.PrivAccount) (*core.ResponseSignTx, error) {
	values, err := argsToURLValues([]string{"tx", "privAccounts"}, tx, privAccounts)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"sign_tx", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseSignTx
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientHTTP) Status() (*core.ResponseStatus, error) {
	values, err := argsToURLValues(nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(c.addr+"status", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseStatus
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) BlockchainInfo(minHeight uint, maxHeight uint) (*core.ResponseBlockchainInfo, error) {
	params, err := binaryWriter(minHeight, maxHeight)
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "blockchain_info",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseBlockchainInfo
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) BroadcastTx(tx types.Tx) (*core.ResponseBroadcastTx, error) {
	params, err := binaryWriter(tx)
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "broadcast_tx",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseBroadcastTx
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) GenPrivAccount() (*core.ResponseGenPrivAccount, error) {
	params, err := binaryWriter()
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "gen_priv_account",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGenPrivAccount
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) GetAccount(address []byte) (*core.ResponseGetAccount, error) {
	params, err := binaryWriter(address)
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "get_account",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGetAccount
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) GetBlock(height uint) (*core.ResponseGetBlock, error) {
	params, err := binaryWriter(height)
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "get_block",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseGetBlock
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) ListAccounts() (*core.ResponseListAccounts, error) {
	params, err := binaryWriter()
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "list_accounts",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseListAccounts
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) ListValidators() (*core.ResponseListValidators, error) {
	params, err := binaryWriter()
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "list_validators",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseListValidators
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) NetInfo() (*core.ResponseNetInfo, error) {
	params, err := binaryWriter()
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "net_info",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseNetInfo
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) SignTx(tx types.Tx, privAccounts []*account.PrivAccount) (*core.ResponseSignTx, error) {
	params, err := binaryWriter(tx, privAccounts)
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "sign_tx",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseSignTx
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}

func (c *ClientJSON) Status() (*core.ResponseStatus, error) {
	params, err := binaryWriter()
	if err != nil {
		return nil, err
	}
	s := rpc.RPCRequest{
		JSONRPC: "2.0",
		Method:  "status",
		Params:  params,
		Id:      0,
	}
	body, err := c.requestResponse(s)
	if err != nil {
		return nil, err
	}
	var status struct {
		Status string
		Data   *core.ResponseStatus
		Error  string
	}
	binary.ReadJSON(&status, body, &err)
	if err != nil {
		return nil, err
	}
	if status.Error != "" {
		return nil, fmt.Errorf(status.Error)
	}
	return status.Data, nil
}
