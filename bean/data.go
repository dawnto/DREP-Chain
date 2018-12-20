package bean

import (
    "encoding/hex"
    "BlockChainTest/mycrypto"
    "math/big"
    "encoding/json"
    "BlockChainTest/accounts"
    "BlockChainTest/config"
)

type BlockHeader struct {
    ChainId              config.ChainIdType
    Version              uint32
    PreviousHash         []byte
    GasLimit             *big.Int
    GasUsed              *big.Int
    Height               uint64
    Timestamp            int64
    StateRoot            []byte
    MerkleRoot           []byte
    TxHashes             [][]byte
    LeaderPubKey         *mycrypto.Point
    MinorPubKeys         []*mycrypto.Point
}

type BlockData struct {
    TxCount              uint32
    TxList               []*Transaction
}

type Block struct {
    Header               *BlockHeader
    Data                 *BlockData
    MultiSig             *MultiSignature
}

type TransactionData struct {
    Version              uint32
    Nonce                uint64
    Type                 uint32
    To                   string
    ChainId              config.ChainIdType
    DestChain            config.ChainIdType
    Amount               *big.Int
    GasPrice             *big.Int
    GasLimit             *big.Int
    Timestamp            int64
    Data                 []byte
    PubKey               *mycrypto.Point
}

type Transaction struct {
    Data                 *TransactionData
    Sig                  *mycrypto.Signature
}

type CrossChainTransaction struct {
    ChainId   config.ChainIdType
    StateRoot []byte
    Trans     []*Transaction
}

type Log struct {
    Address      accounts.CommonAddress
    ChainId      config.ChainIdType
    TxHash       []byte
    Topics       [][]byte
    Data         []byte
}

type MultiSignature struct {
    Sig                  *mycrypto.Signature
    Bitmap               []byte
}

func (tx *Transaction) TxId() (string, error) {
    b, err := json.Marshal(tx.Data)
    if err != nil {
        return "", err
    }
    id := hex.EncodeToString(mycrypto.Hash256(b))
    return id, nil
}

func (tx *Transaction) TxHash() ([]byte, error) {
    b, err := json.Marshal(tx)
    if err != nil {
        return nil, err
    }
    h := mycrypto.Hash256(b)
    return h, nil
}

func (tx *Transaction) TxSig(prvKey *mycrypto.PrivateKey) (*mycrypto.Signature, error) {
    b, err := json.Marshal(tx.Data)
    if err != nil {
        return nil, err
    }
    return mycrypto.Sign(prvKey, b)
}

func (tx *Transaction) GetGasUsed() *big.Int {
    return new(big.Int).SetInt64(int64(100))
}

func (tx *Transaction) GetGas() *big.Int {
    gasQuantity := tx.GetGasUsed()
    gasPrice := new(big.Int).Set(tx.Data.GasPrice)
    gasUsed := new(big.Int).Mul(gasQuantity, gasPrice)
    return gasUsed
}

func (block *Block) BlockHash() ([]byte, error) {
    b, err := json.Marshal(block.Header)
    if err != nil {
        return nil, err
    }
    return mycrypto.Hash256(b), nil
}

func (block *Block) BlockHashHex() (string, error) {
    h, err := block.BlockHash()
    if err != nil {
        return "", err
    }
    return "0x" + hex.EncodeToString(h), nil
}

func (block *Block) TxHashes() []string {
    th := make([]string, len(block.Header.TxHashes))
    for i, hash := range block.Header.TxHashes {
        th[i] = "0x" + hex.EncodeToString(hash)
    }
    return th
}