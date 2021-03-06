// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"github.com/drep-project/binary"
	"github.com/drep-project/DREP-Chain/common"
	"github.com/drep-project/DREP-Chain/crypto"
	"github.com/drep-project/DREP-Chain/crypto/sha3"
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)
)

// Receipt represents the results of a transaction.
type Receipt struct {
	// Consensus fields
	PostState         []byte
	Status            uint64
	CumulativeGasUsed uint64
	Logs              []*Log
	Bloom             Bloom

	// Implementation fields (don't reorder!)
	TxHash          crypto.Hash
	ContractAddress crypto.CommonAddress
	GasUsed         uint64

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash   crypto.Hash
	BlockNumber uint64
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(root []byte, failed bool, cumulativeGasUsed uint64) *Receipt {
	r := &Receipt{PostState: common.CopyBytes(root), CumulativeGasUsed: cumulativeGasUsed}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

//func (r *Receipt) MarshalBinary() ([]byte, error) {
//	return binary.Marshal(r)
//}
//
//func (r *Receipt) UnmarshalBinary(data []byte) error {
//	return binary.Unmarshal(data, r)
//}

func (r *Receipt) ReceiptHash() *crypto.Hash {
	b, _ := binary.Marshal(r)
	h := sha3.Keccak256(b)
	hash := &crypto.Hash{}
	hash.SetBytes(h)
	return hash
}

// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }
