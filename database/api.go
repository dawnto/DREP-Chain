package database

import (
    "BlockChainTest/bean"
    "math/big"
    "BlockChainTest/mycrypto"
    "strconv"
    "encoding/json"
    "BlockChainTest/accounts"
    "BlockChainTest/config"
    "encoding/hex"
    "github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
    db *Database
)

func InitDataBase(config *config.NodeConfig){
    db = NewDatabase(config)
}

func GetItr() iterator.Iterator {
    return db.db.NewIterator(nil, nil)
}

func BeginTransaction() Transactional {
    return db.BeginTransaction()
}

func GetBlock(height uint64) *bean.Block {
    key := mycrypto.Hash256([]byte("block_" + strconv.FormatUint(height, 10)))
    value := db.Get(key)
    block := &bean.Block{}
    err := json.Unmarshal(value, block)
    if err != nil {
        return nil
    }
    return block
}

func PutBlock(block *bean.Block) error {
    key := mycrypto.Hash256([]byte("block_" + strconv.FormatUint(block.Header.Height, 10)))
    value, err := json.Marshal(block)
    if err != nil {
        return err
    }
    return db.PutOutState(config.Hex2ChainId(config.GetConfig().ChainId), key, value)
}

func GetBlocksFrom(start, size uint64) []*bean.Block {
    var (
        currentBlock =&bean.Block{}
        height = start
        blocks = make([]*bean.Block, 0)
    )
    for currentBlock != nil && (height < start + size || size == -1)  {
        currentBlock = GetBlock(height)
        if currentBlock != nil {
            blocks = append(blocks, currentBlock)
        }
        height += 1
    }
    return blocks
}

func GetAllBlocks() []*bean.Block {
    return GetBlocksFrom(uint64(0), uint64(-1))
}

func GetMostRecentBlocks(n uint64) []*bean.Block {
    height := GetMaxHeight()
    if n <= height {
        return GetBlocksFrom(height - n, n)
    } else {
        return GetBlocksFrom(0, height)
    }
}

func GetHighestBlock() *bean.Block {
    maxHeight := GetMaxHeight()
    return GetBlock(maxHeight)
}

func GetMaxHeight() uint64 {
    key := mycrypto.Hash256([]byte("max_height"))
    value := db.Get(key)
    if value == nil {
        return -1
    } else {
        return new(big.Int).SetBytes(value).Uint64()
    }
}

func PutMaxHeight(height int64) error {
    key := mycrypto.Hash256([]byte("max_height"))
    value := new(big.Int).SetInt64(height).Bytes()
    return db.PutOutState(config.Hex2ChainId(config.GetConfig().ChainId), key, value)
}

func GetStorage(addr accounts.CommonAddress, chainId config.ChainIdType) *accounts.Storage {
    key := mycrypto.Hash256([]byte("storage_" + addr.Hex() + chainId.Hex()))
    value := db.Get(key)
    storage := &accounts.Storage{}
    json.Unmarshal(value, storage)
    return storage
}

func PutStorage(t Transactional, addr accounts.CommonAddress, chainId config.ChainIdType, storage *accounts.Storage) error {
    key := mycrypto.Hash256([]byte("storage_" + addr.Hex() + chainId.Hex()))
    value, err := json.Marshal(storage)
    if err != nil {
        return err
    }
    return t.Put(chainId, key, value)
}

func GetBalance(addr accounts.CommonAddress, chainId config.ChainIdType) *big.Int {
    storage := GetStorage(addr, chainId)
    if storage.Balance == nil {
        return new(big.Int)
    }
    return storage.Balance
}

func PutBalance(t Transactional, addr accounts.CommonAddress, chainId config.ChainIdType, balance *big.Int) error {
    storage := GetStorage(addr, chainId)
    storage.Balance = balance
    return PutStorage(t, addr, chainId, storage)
}

func GetNonce(addr accounts.CommonAddress, chainId config.ChainIdType) uint64 {
    storage := GetStorage(addr, chainId)
    return storage.Nonce
}

func PutNonce(t Transactional, addr accounts.CommonAddress, chainId config.ChainIdType, nonce uint64) error {
    storage := GetStorage(addr, chainId)
    storage.Nonce = nonce
    return PutStorage(t, addr, chainId, storage)
}

func GetByteCode(addr accounts.CommonAddress, chainId config.ChainIdType) []byte {
    storage := GetStorage(addr, chainId)
    return storage.ByteCode
}

func PutByteCode(t Transactional, addr accounts.CommonAddress, chainId config.ChainIdType, byteCode []byte) error {
    storage := GetStorage(addr, chainId)
    storage.ByteCode = byteCode
    storage.CodeHash = accounts.GetByteCodeHash(byteCode)
    return PutStorage(t, addr, chainId, storage)
}

func GetCodeHash(addr accounts.CommonAddress, chainId config.ChainIdType) accounts.Hash {
    storage := GetStorage(addr, chainId)
    return storage.CodeHash
}

func GetReputation(addr accounts.CommonAddress, chainId config.ChainIdType) *big.Int {
    storage := GetStorage(addr, chainId)
    if storage.Balance == nil {
        return new(big.Int)
    }
    return storage.Reputation
}

func PutReputation(t Transactional, addr accounts.CommonAddress, chainId config.ChainIdType, reputation *big.Int) error {
    storage := GetStorage(addr, chainId)
    storage.Reputation = reputation
    return PutStorage(t, addr, chainId, storage)
}

func GetLogs(txHash []byte, chainId config.ChainIdType) []*bean.Log {
    key := mycrypto.Hash256([]byte("logs_" + hex.EncodeToString(txHash) + chainId.Hex()))
    value := db.Get(key)
    var logs []*bean.Log
    err := json.Unmarshal(value, &logs)
    if err != nil {
        return make([]*bean.Log, 0)
    }
    return logs
}

func PutLogs(logs []*bean.Log, txHash []byte, chainId config.ChainIdType) error {
    key := mycrypto.Hash256([]byte("logs_" + hex.EncodeToString(txHash) + chainId.Hex()))
    value, err := json.Marshal(logs)
    if err != nil {
        return err
    }
    db.PutOutState(chainId, key, value)
    return nil
}

func AddLog(log *bean.Log) error {
    logs := GetLogs(log.TxHash, log.ChainId)
    logs = append(logs, log)
    return PutLogs(logs, log.TxHash, log.ChainId)
}