package evm

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/drep-project/drep-chain/app"
	"github.com/drep-project/drep-chain/chain/types"
	"github.com/drep-project/drep-chain/database"
	"github.com/drep-project/dlog"
	"github.com/drep-project/drep-chain/pkgs/evm/vm"
	"gopkg.in/urfave/cli.v1"
	"math/big"
)

var (
	DefaultEvmConfig = &vm.VMConfig{}
)


type EvmService struct {
	Config *vm.VMConfig
	DatabaseService *database.DatabaseService `service:"database"`
}

func (evmService *EvmService) Name() string {
	return "vm"
}

func (evmService *EvmService) Api() []app.API {
	return []app.API{}
}

func (evmService *EvmService) CommandFlags() ([]cli.Command, []cli.Flag) {
	return nil, []cli.Flag{}
}

func (evmService *EvmService) P2pMessages() map[int]interface{} {
	return map[int]interface{}{}
}

func (evmService *EvmService) Init(executeContext *app.ExecuteContext) error {
	evmService.Config = DefaultEvmConfig
	err := executeContext.UnmashalConfig(evmService.Name(), evmService.Config)
	if err != nil {
		return err
	}
	return nil
}

func (evmService *EvmService)  Start(executeContext *app.ExecuteContext) error {
	return nil
}

func (evmService *EvmService)  Stop(executeContext *app.ExecuteContext) error{
	return nil
}

func (evmService *EvmService)  Receive(context actor.Context) { }

func (evmService *EvmService)  Eval(state *vm.State,tx *types.Transaction, header *types.BlockHeader, bc ChainContext, gas uint64, value *big.Int) (ret []byte, failed bool, err error)  {
	sender := tx.From()
	contractCreation := tx.To() == nil|| tx.To().IsEmpty()

	// Create a new context to be used in the EVM environment
	context := NewEVMContext(tx, header, bc)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, state, evmService.Config)
	var (
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	if contractCreation {
		ret, _, gas, vmerr = vmenv.Create(*sender, tx.Data.Data, gas, value)
	} else {
		// Increment the nonce for the next transaction
		state.SetNonce(tx.From(), state.GetNonce(sender)+1)
		ret, gas, vmerr = vmenv.Call(*sender, *tx.To(), vmenv.ChainId, tx.Data.Data, gas, value)
	}
	if vmerr != nil {
		dlog.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == vm.ErrInsufficientBalance {
			return nil, false, vmerr
		}
	}
	return ret, vmerr != nil, err
}


