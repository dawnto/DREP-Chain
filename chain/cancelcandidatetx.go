package chain

import (
	"encoding/json"
	"github.com/drep-project/drep-chain/types"
)

/**********************un register ********************/

type CancelCandidateTxSelector struct {
}

func (StakeTxSelector *CancelCandidateTxSelector) Select(tx *types.Transaction) bool {
	return tx.Type() == types.CancelCandidateType
}

var (
	_ = (ITransactionSelector)((*StakeTxSelector)(nil))
	_ = (ITransactionValidator)((*CancelCandidateTransactionProcessor)(nil))
)

type CancelCandidateTransactionProcessor struct {
}

func (processor *CancelCandidateTransactionProcessor) ExecuteTransaction(context *ExecuteTransactionContext) ([]byte, bool, []*types.Log, error) {
	from := context.From()
	tx := context.Tx()
	stakeStore := context.TrieStore()

	detail, err := stakeStore.CancelCandidateCredit(from, tx.Amount(), context.header.Height)
	if err != nil {
		return nil, false, nil, err
	}

	logs := make([]*types.Log, 0, 1)
	data, _ := json.Marshal(detail)

	log := types.Log{TxType: tx.Type(), TxHash: *tx.TxHash(), Data: data, Height: context.header.Height, TxIndex: 0}
	logs = append(logs, &log)

	err = stakeStore.PutNonce(from, tx.Nonce()+1)
	if err != nil {
		return nil, false, nil, err
	}

	return nil, true, logs, err
}
