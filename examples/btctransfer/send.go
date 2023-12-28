package main

import (
	"bytes"
	"encoding/hex"
	"go-ord-tx/pkg/btcapi"

	mempoolApi "go-ord-tx/pkg/btcapi/mempool"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/mempool"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type blockchainClient struct {
	rpcClient    *rpcclient.Client
	btcApiClient btcapi.BTCAPIClient
}

const (
	defaultSequenceNum    = wire.MaxTxInSequenceNum - 10
	defaultRevealOutValue = int64(546) // 500 sat, ord default 10000
	MaxStandardTxWeight   = blockchain.MaxBlockWeight / 10
)

func BuildCommitTx(commitTxOutPointList []*wire.OutPoint, from, destination string, transferAmount, commitFeeRate int64) (*wire.MsgTx, *txscript.MultiPrevOutFetcher, error) {
	totalSenderAmount := btcutil.Amount(0)
	tx := wire.NewMsgTx(wire.TxVersion)
	commitTxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	// var changePkScript *[]byte
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempoolApi.NewClient(netParams, "tx")
	changeAmount := btcutil.Amount(0)
	for i := range commitTxOutPointList {
		txOut, err := getTxOutByOutPoint(btcApiClient, commitTxOutPointList[i])
		if err != nil {
			return nil, commitTxPrevOutputFetcher, err
		}
		commitTxPrevOutputFetcher.AddPrevOut(*commitTxOutPointList[i], txOut)
		// if changePkScript == nil { // first sender as change address
		// 	changePkScript = &txOut.PkScript
		// }
		in := wire.NewTxIn(commitTxOutPointList[i], nil, nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(txOut.Value)

		fee := btcutil.Amount(mempool.GetTxVirtualSize(btcutil.NewTx(tx))) * btcutil.Amount(commitFeeRate)
		changeAmount = totalSenderAmount - btcutil.Amount(transferAmount) - fee
		if changeAmount >= 0 {
			break
			// tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
		} else {
			// tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]

		}
	}
	if changeAmount < 0 {
		feeWithoutChange := btcutil.Amount(mempool.GetTxVirtualSize(btcutil.NewTx(tx))) * btcutil.Amount(commitFeeRate)
		if totalSenderAmount-btcutil.Amount(transferAmount)-feeWithoutChange < 0 {
			return nil, nil, errors.New("insufficient balance")
		}
	}
	receiver, err := btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}
	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, nil, err
	}
	if changeAmount > 0 {
		fromAddress, err := btcutil.DecodeAddress(from, &chaincfg.MainNetParams)
		if err != nil {
			return nil, nil, err
		}
		scriptPubKey, err := txscript.PayToAddrScript(fromAddress)
		if err != nil {
			return nil, nil, err
		}
		changeStr := int64(changeAmount)
		// amountInt64, err := strconv.ParseInt(changeStr, 10, 64)

		out := wire.NewTxOut(changeStr, scriptPubKey)
		tx.AddTxOut(out)
	}
	out := wire.NewTxOut(transferAmount, scriptPubKey)
	tx.AddTxOut(out)

	return tx, commitTxPrevOutputFetcher, nil
}

func BuildSendAllTx(commitTxOutPointList []*wire.OutPoint, from, destination string, commitFeeRate int64) (*wire.MsgTx, *txscript.MultiPrevOutFetcher, error) {
	totalSenderAmount := btcutil.Amount(0)
	tx := wire.NewMsgTx(wire.TxVersion)
	commitTxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	// var changePkScript *[]byte
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempoolApi.NewClient(netParams, "tx")
	transferAmount := btcutil.Amount(0)
	for i := range commitTxOutPointList {
		txOut, err := getTxOutByOutPoint(btcApiClient, commitTxOutPointList[i])
		if err != nil {
			return nil, commitTxPrevOutputFetcher, err
		}
		commitTxPrevOutputFetcher.AddPrevOut(*commitTxOutPointList[i], txOut)
		// if changePkScript == nil { // first sender as change address
		// 	changePkScript = &txOut.PkScript
		// }
		in := wire.NewTxIn(commitTxOutPointList[i], nil, nil)
		in.Sequence = defaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(txOut.Value)
	}
	fee := btcutil.Amount(mempool.GetTxVirtualSize(btcutil.NewTx(tx))) * btcutil.Amount(commitFeeRate)
	transferAmount = totalSenderAmount - fee
	if transferAmount < 0 {
		return nil, nil, errors.New("insufficient balance")
	}
	receiver, err := btcutil.DecodeAddress(destination, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}
	scriptPubKey, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, nil, err
	}
	out := wire.NewTxOut(int64(transferAmount), scriptPubKey)
	tx.AddTxOut(out)

	// tx.AddTxOut(wire.NewTxOut(0, *changePkScript))

	return tx, commitTxPrevOutputFetcher, nil

}

func getTxOutByOutPoint(btcApiClient *mempoolApi.MempoolClient, outPoint *wire.OutPoint) (*wire.TxOut, error) {
	var txOut *wire.TxOut

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

func signCommitTx(privateKey *btcec.PrivateKey, commitTx *wire.MsgTx, commitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher) (*wire.MsgTx, error) {

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

func signCommitTx1(privateKey *btcec.PrivateKey, commitTx *wire.MsgTx, commitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher) (*wire.MsgTx, error) {

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

func getTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func GetCommitTxHex(commitTx *wire.MsgTx) (string, error) {
	return getTxHex(commitTx)
}

func SendRawTransaction(btcApiClient *mempoolApi.MempoolClient, tx *wire.MsgTx) (*chainhash.Hash, error) {
	return btcApiClient.BroadcastTx(tx)
}
