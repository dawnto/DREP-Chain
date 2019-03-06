package service

import (
	"errors"
	"fmt"
	"github.com/drep-project/dlog"
	"github.com/drep-project/drep-chain/app"
	chainTypes "github.com/drep-project/drep-chain/chain/types"
	"github.com/drep-project/drep-chain/crypto"
	"github.com/drep-project/drep-chain/pkgs/evm/vm"
	txType "github.com/drep-project/drep-chain/transaction/types"
	"math/big"

	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

var (
	childTrans []*txType.Transaction
)

func (chainService *ChainService) ExecuteTransactions(b *chainTypes.Block) (*big.Int, error) {
	if b == nil || b.Header == nil { // || b.Data == nil || b.Data.TxList == nil {
		return nil, errors.New("error block nil or header nil")
	}
	height := chainService.DatabaseService.GetMaxHeight()
	if height+1 != b.Header.Height {
		msg := fmt.Sprintf("not corrent height CurrentHeight: %d, ReceiveHeight: %d", height, b.Header.Height)
		dlog.Error(msg)
		return nil, errors.New(msg)
	}

	chainService.DatabaseService.BeginTransaction()
	total := big.NewInt(0)
	if b.Data == nil {
		return total, nil
	}
	for _, t := range b.Data.TxList {
		_, gasFee := chainService.execute(t)

		if gasFee != nil {
			total.Add(total, gasFee)
		}
		if t.Nonce()%10 == 0 {
			fmt.Println("execute tx :", t.From().Hex(), t.Nonce(), t.Amount(), string(t.Sig()))
		}
	}

	stateRoot := chainService.DatabaseService.GetStateRoot()
	if bytes.Equal(b.Header.StateRoot, stateRoot) {
		fmt.Println()
		fmt.Println("matched ", hex.EncodeToString(b.Header.StateRoot), " vs ", hex.EncodeToString(stateRoot))
		height++
		chainService.DatabaseService.PutMaxHeight(height)
		chainService.DatabaseService.PutBlock(b)
		chainService.DatabaseService.Commit()
		fmt.Println("received block: ", true)
		fmt.Println()

		chainService.accumulateRewards(b, chainService.ChainID())
		chainService.preSync(b)
		chainService.doSync(height)
	} else {
		chainService.DatabaseService.Discard()
		fmt.Println()
		fmt.Println("not matched ", hex.EncodeToString(b.Header.StateRoot), " vs ", hex.EncodeToString(stateRoot))
		fmt.Println("received block: ", false)
		fmt.Println()
	}
	return total, nil
}

func (chainService *ChainService) preSync(block *chainTypes.Block) {
	if !chainService.isRelay && chainService.chainId != chainService.RootChain() {
		return
	}
	if childTrans == nil {
		childTrans = make([]*txType.Transaction, 0)
	}
	childTrans = append(childTrans, block.Data.TxList...)
}

func (chainService *ChainService) doSync(height int64) {
	if !chainService.isRelay || chainService.chainId == chainService.RootChain() || height%2 != 0 || height == 0 {
		return
	}
	cct := &txType.CrossChainTransaction{
		ChainId:   chainService.chainId,
		StateRoot: chainService.DatabaseService.GetStateRoot(),
		Trans:     childTrans,
	}
	data, err := json.Marshal(cct)
	if err != nil {
		return
	}
	values := url.Values{}
	values.Add("data", string(data))
	body := values.Encode()
	urlStr := "http://localhost:" + strconv.Itoa(chainService.Config.RemotePort) + "/SyncChildChain?" + body
	http.Get(urlStr)
	childTrans = nil
}

func (chainService *ChainService) execute(t *txType.Transaction) (gasUsed, gasFee *big.Int) {

	switch t.Type() {
	case txType.TransferType:
		return chainService.executeTransferTransaction(t)
	case txType.CreateContractType:
		return chainService.executeCreateContractTransaction(t)
	case txType.CallContractType:
		return chainService.executeCallContractTransaction(t)
		//case CrossChainType:
		//   return chainService.executeCrossChainTransaction(t)
	}
	return nil, nil
}

func (chainService *ChainService) canExecute(tx *txType.Transaction, gasFloor, gasCap *big.Int) (canExecute bool, addr crypto.CommonAddress, balance, gasLimit, gasPrice *big.Int) {
	addr = *tx.From()
	balance = chainService.DatabaseService.GetBalance(&addr, true)
	nonce := chainService.DatabaseService.GetNonce(&addr, true) + 1
	chainService.DatabaseService.PutNonce(&addr, nonce, true)
	gasPrice = tx.GasPrice()

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
	return
}

func (chainService *ChainService) deduct(addr crypto.CommonAddress, chainId app.ChainIdType, balance, gasFee *big.Int) (leftBalance, actualFee *big.Int) {
	leftBalance = new(big.Int).Sub(balance, gasFee)
	actualFee = new(big.Int).Set(gasFee)
	if leftBalance.Sign() < 0 {
		actualFee = new(big.Int).Set(balance)
		leftBalance = new(big.Int)
	}
	chainService.DatabaseService.PutBalance(&addr, leftBalance, true)
	return leftBalance, actualFee
}

func (chainService *ChainService) executeTransferTransaction(t *txType.Transaction) (gasUsed *big.Int, gasFee *big.Int) {
	var (
		can               bool
		addr              crypto.CommonAddress
		balance, gasPrice *big.Int
	)

	gasUsed, gasFee = new(big.Int), new(big.Int)
	can, addr, balance, _, gasPrice = chainService.canExecute(t, txType.TransferGas, nil)
	if !can {
		return
	}

	gasUsed = new(big.Int).Set(txType.TransferGas)
	gasFee = new(big.Int).Mul(gasUsed, gasPrice)
	balance, gasFee = chainService.deduct(addr, t.ChainId(), balance, gasFee)
	if balance.Cmp(t.Amount()) >= 0 {
		balance = new(big.Int).Sub(balance, t.Amount())
		balanceTo := chainService.DatabaseService.GetBalance(t.To(), true)
		balanceTo = new(big.Int).Add(balanceTo, t.Amount())
		chainService.DatabaseService.PutBalance(&addr, balance, true)
		chainService.DatabaseService.PutBalance(t.To(), balanceTo, true)
	}
	return
}

func (chainService *ChainService) executeCreateContractTransaction(t *txType.Transaction) (gasUsed *big.Int, gasFee *big.Int) {
	var (
		can                         bool
		addr                        crypto.CommonAddress
		balance, gasLimit, gasPrice *big.Int
	)
	gasUsed, gasFee = new(big.Int), new(big.Int)
	can, addr, _, gasLimit, gasPrice = chainService.canExecute(t, nil, txType.CreateContractGas)
	if !can {
		return
	}

	evm := vm.NewEVM(chainService.DatabaseService)
	returnGas, _ := chainService.VmService.ApplyTransaction(evm, t)
	gasUsed = new(big.Int).Sub(gasLimit, new(big.Int).SetUint64(returnGas))
	gasFee = new(big.Int).Mul(gasUsed, gasPrice)
	balance = chainService.DatabaseService.GetBalance(&addr, true)
	_, gasFee = chainService.deduct(addr, t.ChainId(), balance, gasFee)
	return
}

func (chainService *ChainService) executeCallContractTransaction(t *txType.Transaction) (gasUsed *big.Int, gasFee *big.Int) {
	var (
		can                         bool
		addr                        crypto.CommonAddress
		balance, gasLimit, gasPrice *big.Int
	)

	gasUsed, gasFee = new(big.Int), new(big.Int)
	can, addr, _, gasLimit, gasPrice = chainService.canExecute(t, nil, txType.CallContractGas)
	if !can {
		return
	}

	evm := vm.NewEVM(chainService.DatabaseService)
	returnGas, _ := chainService.VmService.ApplyTransaction(evm, t)
	gasUsed = new(big.Int).Sub(gasLimit, new(big.Int).SetUint64(returnGas))
	gasFee = new(big.Int).Mul(gasUsed, gasPrice)
	balance = chainService.DatabaseService.GetBalance(&addr, true)
	_, gasFee = chainService.deduct(addr, t.ChainId(), balance, gasFee)
	return
}

//func (chainService *ChainService) executeCrossChainTransaction(t *chainTypes.Transaction) (gasUsed *big.Int, gasFee *big.Int) {
//    var (
//        can bool
//        addr crypto.CommonAddress
//        balance, gasPrice *big.Int
//    )
//
//    gasUsed, gasFee = new(big.Int), new(big.Int)
//    can, addr,  _, _, gasPrice = chainService.canExecute(t, nil, CrossChainGas)
//    if !can {
//        return new(big.Int), new(big.Int)
//    }
//
//    cct := &chainTypes.CrossChainTransaction{}
//    err := json.Unmarshal(t.Data.Data, cct)
//    if err != nil {
//        fmt.Println("err: ", err)
//        return new(big.Int), new(big.Int)
//    }
//
//    gasSum := new(big.Int)
//    for _, tx := range cct.Trans {
//       if tx.Data.Type == CrossChainType {
//           continue
//       }
//       g, _ := chainService.execute(tx)
//       gasSum = new(big.Int).Add(gasSum, g)
//    }
//
//    if !bytes.Equal(chainService.databaseService.GetStateRoot(), cct.StateRoot) {
//       //subDt.Discard()
//    } else {
//        amountSum := new(big.Int).Mul(gasSum, gasPrice)
//        balance = chainService.databaseService.GetBalance(addr, t.Data.ChainId, true)
//        if balance.Cmp(amountSum) >= 0 {
//            gasUsed = new(big.Int).Set(gasSum)
//            gasFee = new(big.Int).Set(amountSum)
//            _, gasFee = chainService.deduct(addr, t.Data.ChainId, balance, gasFee)
//            //subDt.Commit()
//        } else {
//            //subDt.Discard()
//        }
//    }
//    return
//}

//func preExecuteCrossChainTransaction(dt database.Transactional, t *chainTypes.Transaction) (gasUsed, gasFee *big.Int) {
//    var (
//        can bool
//        addr crypto.CommonAddress
//        balance, gasPrice *big.Int
//    )
//
//    gasUsed, gasFee = new(big.Int), new(big.Int)
//    subDt := dt.BeginTransaction()
//    can, addr,  _, _, gasPrice = canExecute(subDt, t, nil, CrossChainGas)
//    if !can {
//        return new(big.Int), new(big.Int)
//    }
//
//    cct := &chainTypes.CrossChainTransaction{}
//    err := json.Unmarshal(t.Data.Data, &cct)
//    if err != nil {
//        return new(big.Int), new(big.Int)
//    }
//
//    gasSum := new(big.Int)
//    for _, tx := range cct.Trans {
//        if tx.Data.Type == CrossChainType {
//            continue
//        }
//        g, _ := execute(subDt, tx)
//        gasSum = new(big.Int).Add(gasSum, g)
//    }
//
//    cct.StateRoot = subDt.GetChainStateRoot(database.ChildCHAIN)
//    t.Data.Data, _ = json.Marshal(cct)
//
//    amountSum := new(big.Int).Mul(gasSum, gasPrice)
//    balance = database.GetBalance(addr, t.Data.ChainId)
//    if balance.Cmp(amountSum) >= 0 {
//        gasUsed = new(big.Int).Set(gasSum)
//        gasFee = new(big.Int).Set(amountSum)
//        _, gasFee = deduct(subDt, addr, t.Data.ChainId, balance, gasFee)
//        subDt.Commit()
//    } else {
//        subDt.Discard()
//    }
//
//    return
//}