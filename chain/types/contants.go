package types

import "math/big"

type TxType  uint64
const(
	TransferType TxType = iota
	MinerType
	CreateContractType
	CallContractType
	CrossChainType
)


var (
	DefaultGasPrice    *big.Int
	TransferGas              = big.NewInt(60000)
	MinerGas                 = big.NewInt(20000)
	CreateContractGas        = big.NewInt(1000000)
	CallContractGas          = big.NewInt(10000000)
	CrossChainGas            = big.NewInt(10000000)
	GasTable = map[TxType]*big.Int{}
)

func init() {
	DefaultGasPrice, _ = new(big.Int).SetString("20000000000", 10)
	GasTable[TransferType] = TransferGas
	GasTable[MinerType] = MinerGas
	GasTable[CreateContractType] = CreateContractGas
	GasTable[CallContractType] = CallContractGas
	GasTable[CrossChainType] = CrossChainGas
}
