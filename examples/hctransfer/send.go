package main

import (
	"bytes"
	"encoding/hex"
	"go-ord-tx/pkg/btcapi"

	mempoolApi "go-ord-tx/pkg/btcapi/mempool"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/mempool"

	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/txscript"
	hcWire "github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
)

type blockchainClient struct {
	rpcClient    *rpcclient.Client
	btcApiClient btcapi.BTCAPIClient
}

const (
	defaultSequenceNum    = hcWire.MaxTxInSequenceNum - 10
	defaultRevealOutValue = int64(546) // 500 sat, ord default 10000
	MaxStandardTxWeight   = blockchain.MaxBlockWeight / 10
)

func BuildCommitTx(commitTxOutPointList []*hcWire.OutPoint, from, destination string, transferAmount, commitFeeRate int64) (*hcWire.MsgTx, error) {
	totalSenderAmount := hcutil.Amount(0)
	tx := hcWire.NewMsgTx()
	// var changePkScript *[]byte
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempoolApi.NewClient(netParams, "tx")
	changeAmount := hcutil.Amount(0)
	for i := range commitTxOutPointList {
		txOut, err := getTxOutByOutPoint(btcApiClient, commitTxOutPointList[i])
		if err != nil {
			return nil, err
		}
		// if changePkScript == nil { // first sender as change address
		// 	changePkScript = &txOut.PkScript
		// }
		in := hcWire.NewTxIn(commitTxOutPointList[i], nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += hcutil.Amount(txOut.Value)

		fee := hcutil.Amount(20000000)
		changeAmount = totalSenderAmount - hcutil.Amount(transferAmount) - fee
		if changeAmount >= 0 {
			break
			// tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
		} else {
			// tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]

		}
	}
	if changeAmount < 0 {
		feeWithoutChange := hcutil.Amount(20000000)
		if totalSenderAmount-hcutil.Amount(transferAmount)-feeWithoutChange < 0 {
			return nil, errors.New("insufficient balance")
		}
	}
	receiver, err := hcutil.DecodeAddress(destination)
	if err != nil {
		return nil, err
	}
	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, err
	}
	if changeAmount > 0 {
		fromAddress, err := hcutil.DecodeAddress(from, &chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}
		scriptPubKey, err := txscript.PayToAddrScript(fromAddress)
		if err != nil {
			return nil, err
		}
		changeStr := int64(changeAmount)
		// amountInt64, err := strconv.ParseInt(changeStr, 10, 64)

		out := hcWire.NewTxOut(changeStr, scriptPubKey)
		tx.AddTxOut(out)
	}
	out := hcWire.NewTxOut(transferAmount, scriptPubKey)
	tx.AddTxOut(out)

	return tx, nil
}

func BuildSendAllTx(commitTxOutPointList []*hcWire.OutPoint, from, destination string, commitFeeRate int64) (*hcWire.MsgTx, *txscript.MultiPrevOutFetcher, error) {
	totalSenderAmount := hcutil.Amount(0)
	tx := hcWire.NewMsgTx(hcWire.TxVersion)
	commitTxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	// var changePkScript *[]byte
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempoolApi.NewClient(netParams, "tx")
	transferAmount := hcutil.Amount(0)
	for i := range commitTxOutPointList {
		txOut, err := getTxOutByOutPoint(btcApiClient, commitTxOutPointList[i])
		if err != nil {
			return nil, commitTxPrevOutputFetcher, err
		}
		commitTxPrevOutputFetcher.AddPrevOut(*commitTxOutPointList[i], txOut)
		// if changePkScript == nil { // first sender as change address
		// 	changePkScript = &txOut.PkScript
		// }
		in := hcWire.NewTxIn(commitTxOutPointList[i], nil, nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += hcutil.Amount(txOut.Value)
	}
	fee := hcutil.Amount(mempool.GetTxVirtualSize(hcutil.NewTx(tx))) * hcutil.Amount(commitFeeRate)
	transferAmount = totalSenderAmount - fee
	if transferAmount < 0 {
		return nil, nil, errors.New("insufficient balance")
	}
	receiver, err := hcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}
	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, nil, err
	}
	out := hcWire.NewTxOut(int64(transferAmount), scriptPubKey)
	tx.AddTxOut(out)

	// tx.AddTxOut(hcWire.NewTxOut(0, *changePkScript))

	return tx, commitTxPrevOutputFetcher, nil

}

func getTxOutByOutPoint(btcApiClient *mempoolApi.MempoolClient, outPoint *hcWire.OutPoint) (*hcWire.TxOut, error) {
	var txOut *hcWire.TxOut

	tx, err := btcApiClient.GetRawTransaction(&outPoint.Hash)
	if err != nil {
		return nil, err
	}
	if int(outPoint.Index) >= len(tx.TxOut) {
		return nil, errors.New("err out point")
	}
	txOut = tx.TxOut[outPoint.Index]
	return txOut, nil
}

func signCommitTx(privateKey *btcec.PrivateKey, commitTx *hcWire.MsgTx, commitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher) (*hcWire.MsgTx, error) {

	for i := range commitTx.TxIn {
		txOut := commitTxPrevOutputFetcher.FetchPrevOutput(commitTx.TxIn[i].PreviousOutPoint)
		// signatureScript, err := txscript.SignatureScript(commitTx, i, txOut.PkScript, txscript.SigHashAll, privateKey, true)
		witness, err := txscript.TaprootWitnessSignature(commitTx, txscript.NewTxSigHashes(commitTx, commitTxPrevOutputFetcher),
			i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, privateKey) // i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, privateKey)
		if err != nil {
			return nil, err
		}
		// commitTx.TxIn[i].SignatureScript = signatureScript
		commitTx.TxIn[i].Witness = witness

		// signatureScript[i] = witness
	}
	// for i := range signatureScript {
	// 	// commitTx.TxIn[i].Witness = witnessList[i]
	// }
	return commitTx, nil
}

func signCommitTx1(privateKey *btcec.PrivateKey, commitTx *hcWire.MsgTx, commitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher) (*hcWire.MsgTx, error) {

	for i := range commitTx.TxIn {
		txOut := commitTxPrevOutputFetcher.FetchPrevOutput(commitTx.TxIn[i].PreviousOutPoint)
		signatureScript, err := txscript.SignatureScript(commitTx, i, txOut.PkScript, txscript.SigHashAll, privateKey, true)
		if err != nil {
			return nil, err
		}
		commitTx.TxIn[i].SignatureScript = signatureScript
		//commitTx.TxIn[i].Witness = witness

		// signatureScript[i] = witness
	}
	// for i := range signatureScript {
	// 	// commitTx.TxIn[i].Witness = witnessList[i]
	// }
	return commitTx, nil
}

func getTxHex(tx *hcWire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func GetCommitTxHex(commitTx *hcWire.MsgTx) (string, error) {
	return getTxHex(commitTx)
}

func SendRawTransaction(btcApiClient *mempoolApi.MempoolClient, tx *hcWire.MsgTx) (*chainhash.Hash, error) {
	return btcApiClient.BroadcastTx(tx)
}
