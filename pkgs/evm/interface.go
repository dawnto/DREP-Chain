package evm

import (
	"github.com/drep-project/drep-chain/app"
	"github.com/drep-project/drep-chain/pkgs/evm/vm"
	chainTypes "github.com/drep-project/drep-chain/chain/types"
	"math/big"
)

type Vm interface {
	app.Service
	Eval( *vm.State, *chainTypes.Transaction, *chainTypes.BlockHeader, ChainContext, uint64, *big.Int) (ret []byte, failed bool, err error)
}