package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"go-ord-tx/internal/ord"
	"go-ord-tx/pkg/btcapi/mempool"
	"log"
	"os"
)

func main() {
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempool.NewClient(netParams, "tx")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory, %v", err)
	}
	filePath := fmt.Sprintf("%s/examples/brc20/transfer.txt", workingDir)
	// if file size too max will return sendrawtransaction RPC error: {"code":-26,"message":"tx-size"}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}

	//contentType := http.DetectContentType(fileContent)
	contentType := "text/plain;charset=utf-8"
	log.Printf("file contentType %s", contentType)

	utxoPrivateKeyHex := "d143445ddd5bc927b3f629c50bef0cd15a9698f7f4b5d58b98f34bafaaa287e0"
	destination := "1AkPTypiuVm4ztRkuG16rSeu5nfY2ETRYP"

	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)

	{
		utxoPrivateKeyBytes, err := hex.DecodeString(utxoPrivateKeyHex)
		if err != nil {
			log.Fatal(err)
		}
		utxoPrivateKey, _ := btcec.PrivKeyFromBytes(utxoPrivateKeyBytes)

		utxoTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(utxoPrivateKey.PubKey())), netParams)
		if err != nil {
			log.Fatal(err)
		}

		unspentList, err := btcApiClient.ListUnspent(utxoTaprootAddress)

		if err != nil {
			log.Fatalf("list unspent err %v", err)
		}
		for i := range unspentList {
			for j := i + 1; j < len(unspentList); j++ {
				if unspentList[i].Output.Value < unspentList[j].Output.Value { //brc20铸造，优先选大的，即不含铭文的utxo
					temp := unspentList[i]
					unspentList[i] = unspentList[j]
					unspentList[j] = temp
				}
			}
		}
		for i := range unspentList {
			commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
			commitTxPrivateKeyList = append(commitTxPrivateKeyList, utxoPrivateKey)
		}
	}

	request := ord.InscriptionRequest{
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          80,
		FeeRate:                100,
		DataList: []ord.InscriptionData{
			{
				ContentType: contentType,
				Body:        fileContent,
				Destination: destination,
			},
		},
		SingleRevealTxOnly: false,
	}

	tool, err := ord.NewInscriptionToolWithBtcApiClient(netParams, btcApiClient, &request)
	if err != nil {
		log.Fatalf("Failed to create inscription tool: %v", err)
	}
	commitTxHash, revealTxHashList, inscriptions, fees, err := tool.Inscribe()
	if err != nil {
		log.Fatalf("send tx errr, %v", err)
	}
	log.Println("commitTxHash, " + commitTxHash.String())
	for i := range revealTxHashList {
		log.Println("revealTxHash, " + revealTxHashList[i].String())
	}
	for i := range inscriptions {
		log.Println("inscription, " + inscriptions[i])
	}
	log.Println("fees: ", fees)
}
