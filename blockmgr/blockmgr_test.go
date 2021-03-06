package blockmgr

import (
	"crypto/rand"
	"github.com/drep-project/binary"
	"github.com/drep-project/DREP-Chain/chain"
	"github.com/drep-project/DREP-Chain/crypto"
	"github.com/drep-project/DREP-Chain/crypto/sha3"
	"github.com/drep-project/DREP-Chain/database"
	"github.com/drep-project/DREP-Chain/types"
	"testing"
	"time"
)

func generatorChain(t *testing.T, db *database.Database) []*types.Block {

	privKey, _ := crypto.GenerateKey(rand.Reader)
	ds := database.NewDatabaseService(db)

	producer := types.Producers{Pubkey: privKey.PubKey(), IP: "127.0.0.1"}
	chain.DefaultChainConfig.Producers = append(chain.DefaultChainConfig.Producers, producer)

	cs := chain.NewChainService(chain.DefaultChainConfig, ds)

	bm := NewBlockMgr(DefaultChainConfig, "./", cs, &p2pServiceMock{})

	blks := make([]*types.Block, 0)

	for i := 0; i < 10; i++ {
		block, _, err := bm.GenerateTemplate(db, privKey.PubKey())
		if err != nil {
			t.Fatal("generate block err:", err)
			return nil
		}

		msg, err := binary.Marshal(block)
		if err != nil {
			t.Fatal("generate block err:", err)
			return nil
		}

		sig, err := privKey.Sign(sha3.Keccak256(msg))
		if err != nil {
			t.Fatal("generate block err:", err)
			return nil
		}

		multiSig := &types.MultiSignature{Sig: *sig, Bitmap: []byte{1}}

		block.MultiSig = multiSig
		db = bm.ChainService.GetDatabaseService().BeginTransaction()
		gp := new(chain.GasPool).AddGas(block.Header.GasLimit.Uint64())
		//process transaction
		_, gasFee, err := bm.ChainService.BlockValidator().ExecuteBlock(db, block, gp)
		if err != nil {
			log.WithField("ExecuteBlock", err).Debug("multySigVerify")
			return nil
		}
		err = bm.ChainService.AccumulateRewards(db, block, gasFee)
		if err != nil {
			return nil
		}

		block.Header.StateRoot = db.GetStateRoot()

		_, _, err = bm.ChainService.ProcessBlock(block)
		if err != nil {
			t.Fatal("process block err:", err)
			return nil
		}

		blks = append(blks, block)
		time.Sleep(time.Second)
	}

	return blks
}
