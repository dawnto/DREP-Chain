package service

import (
    "math/big"
    chainType "github.com/drep-project/drep-chain/chain/types"
    "github.com/drep-project/drep-chain/database"
    "github.com/drep-project/drep-chain/crypto"
    "errors"
    "encoding/hex"
    txType "github.com/drep-project/drep-chain/transaction/types"
)


type ChainApi struct {
    chainService *ChainService
    dbService *database.DatabaseService `service:"database"`
}

func (chain *ChainApi) GetBlock(height int64) *chainType.Block {
    if height < 0 {
        return nil
    }
    return chain.dbService.GetBlock(height)
}

func (chain *ChainApi) GetMaxHeight() int64 {
    return chain.dbService.GetMaxHeight()
}

func (chain *ChainApi) GetBalance(addr crypto.CommonAddress) *big.Int{
    return chain.dbService.GetBalance(&addr, true)
}

func (chain *ChainApi) GetNonce(addr crypto.CommonAddress) int64 {
    return chain.dbService.GetNonce(&addr, true)
}

func (chain *ChainApi) GetPreviousBlockHash() string {
    bytes := chain.dbService.GetPreviousBlockHash()
    return "0x" + string(bytes)
}

func (chain *ChainApi) GetReputation(addr crypto.CommonAddress) *big.Int {
    return chain.dbService.GetReputation(&addr, true)
}

func (chain *ChainApi) GetTransactionsFromBlock(height int64) []*txType.Transaction {
    block := chain.dbService.GetBlock(height)
    return block.Data.TxList
}

func (chain *ChainApi) GetTransactionByBlockHeightAndIndex(height int64, index int) *txType.Transaction{
    block := chain.dbService.GetBlock(height)
    if index > len(block.Data.TxList) {
        return nil
    }
    return block.Data.TxList[index]
}

func (chain *ChainApi) GetTransactionCountByBlockHeight(height int64) int {
    block := chain.dbService.GetBlock(height)
    return len(block.Data.TxList)
}

func (chain *ChainApi) SendRawTransaction(tx *txType.Transaction) (string, error){
    //bytes := []byte(raw)
    //tx := &chainType.Transaction{}
    //json.Unmarshal(bytes, tx)

    can := false
    switch tx.Type() {
    case txType.TransferType:
        can, _, _, _, _ = chain.canExecute(tx, txType.TransferGas, nil)
    case txType.CreateContractType:
        can, _, _, _, _ = chain.canExecute(tx, nil, txType.CreateContractGas)
    case txType.CallContractType:
        can, _, _, _, _ = chain.canExecute(tx,nil, txType.CallContractGas)
    }

    if !can {
        return "", errors.New("error: can not execute this transaction")
    }

    err := chain.chainService.transactionPool.AddTransaction(tx)
    if err != nil {
        return "", err
    }

    chain.chainService.P2pServer.Broadcast(tx)

    hash, err:= tx.TxHash()
    encodedHash := hex.EncodeToString(hash)
    res := "0x" + string(encodedHash)
    return res, err
}

func (chain *ChainApi) canExecute(tx *txType.Transaction, gasFloor, gasCap *big.Int) (canExecute bool, addr crypto.CommonAddress, balance, gasLimit, gasPrice *big.Int) {
    chain.chainService.DatabaseService.BeginTransaction()
    addr = *tx.From()
    balance = chain.chainService.DatabaseService.GetBalance(&addr, true)
    nonce :=  chain.chainService.DatabaseService.GetNonce(&addr,true) + 1
    chain.chainService.DatabaseService.PutNonce(&addr, nonce,true)

    if nonce != tx.Nonce() {
       return
    }
    if gasFloor != nil {
        amountFloor := new(big.Int).Mul(gasFloor, tx.GasPrice())
        if tx.GasLimit().Cmp(gasFloor) < 0 || amountFloor.Cmp(balance) > 0 {
            return
        }
    }
    if gasCap != nil {
        amountCap := new(big.Int).Mul(gasCap, tx.GasPrice())
        if amountCap.Cmp(balance) > 0 {
            return
        }
    }

    canExecute = true
    chain.chainService.DatabaseService.Discard()
    return
}