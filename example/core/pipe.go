package core

import (
	bc "github.com/tendermint/tendermint2/blockchain"
	"github.com/tendermint/tendermint2/consensus"
	mempl "github.com/tendermint/tendermint2/mempool"
	"github.com/tendermint/tendermint2/p2p"
)

var blockStore *bc.BlockStore
var consensusState *consensus.ConsensusState
var mempoolReactor *mempl.MempoolReactor
var p2pSwitch *p2p.Switch

func SetPipeBlockStore(bs *bc.BlockStore) {
	blockStore = bs
}

func SetPipeConsensusState(cs *consensus.ConsensusState) {
	consensusState = cs
}

func SetPipeMempoolReactor(mr *mempl.MempoolReactor) {
	mempoolReactor = mr
}

func SetPipeSwitch(sw *p2p.Switch) {
	p2pSwitch = sw
}
