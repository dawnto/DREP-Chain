package bft

import (
	"github.com/drep-project/binary"
	"github.com/drep-project/DREP-Chain/chain"
	"github.com/drep-project/DREP-Chain/crypto/secp256k1"
	"github.com/drep-project/DREP-Chain/crypto/secp256k1/schnorr"
	"github.com/drep-project/DREP-Chain/crypto/sha3"
	types2 "github.com/drep-project/DREP-Chain/pkgs/consensus/types"
	"github.com/drep-project/DREP-Chain/types"
)

type BlockMultiSigValidator struct {
	Producers types2.ProducerSet
}

func (blockMultiSigValidator *BlockMultiSigValidator) VerifyHeader(header, parent *types.BlockHeader) error {
	// check multisig
	// leader
	return nil
}

func (blockMultiSigValidator *BlockMultiSigValidator) VerifyBody(block *types.Block) error {
	participators := []*secp256k1.PublicKey{}
	multiSig := &MultiSignature{}
	err := binary.Unmarshal(block.Proof.Evidence, multiSig)
	if err != nil {
		return err
	}

	//Non outgoing node, only accept incoming block
	if len(blockMultiSigValidator.Producers) == 0 {
		return nil
	}

	for index, val := range multiSig.Bitmap {
		if val == 1 {
			producer := blockMultiSigValidator.Producers[index]
			participators = append(participators, producer.Pubkey)
		}
	}
	msg := block.AsSignMessage()
	sigmaPk := schnorr.CombinePubkeys(participators)

	if !schnorr.Verify(sigmaPk, sha3.Keccak256(msg), multiSig.Sig.R, multiSig.Sig.S) {
		return ErrMultiSig
	}
	return nil
}

func (blockMultiSigValidator *BlockMultiSigValidator) ExecuteBlock(context *chain.BlockExecuteContext) error {
	multiSig := &MultiSignature{}
	err := binary.Unmarshal(context.Block.Proof.Evidence, multiSig)
	if err != nil {
		return nil
	}
	AccumulateRewards(context.Db, multiSig, blockMultiSigValidator.Producers, context.GasFee)
	return nil
}
