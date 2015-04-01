package core

import (
	"github.com/tendermint/tendermint2/account"
	sm "github.com/tendermint/tendermint2/state"
	"github.com/tendermint/tendermint2/types"
)

type ResponseGenPrivAccount struct {
	PrivAccount *account.PrivAccount
}

type ResponseGetAccount struct {
	Account *account.Account
}

type ResponseListAccounts struct {
	BlockHeight uint
	Accounts    []*account.Account
}

type ResponseBlockchainInfo struct {
	LastHeight uint
	BlockMetas []*types.BlockMeta
}

type ResponseGetBlock struct {
	BlockMeta *types.BlockMeta
	Block     *types.Block
}

// curl -H 'content-type: text/plain;' http://127.0.0.1:8888/submit_tx?tx=...
type ResponseBroadcastTx struct {
	Receipt Receipt
}

type ResponseStatus struct {
	GenesisHash       []byte
	Network           string
	LatestBlockHash   []byte
	LatestBlockHeight uint
	LatestBlockTime   int64 // nano
}

type ResponseNetInfo struct {
	NumPeers  int
	Listening bool
	Network   string
}

type ResponseSignTx struct {
	Tx types.Tx
}

type ResponseListValidators struct {
	BlockHeight         uint
	BondedValidators    []*sm.Validator
	UnbondingValidators []*sm.Validator
}
