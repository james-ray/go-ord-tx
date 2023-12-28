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
	"net/http"
	"os"
)

func main() {
	netParams := &chaincfg.SigNetParams
	btcApiClient := mempool.NewClient(netParams, "tx")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory, %v", err)
	}
	filePath := fmt.Sprintf("%s/examples/inscribefile/1.jpeg", workingDir)
	// if file size too max will return sendrawtransaction RPC error: {"code":-26,"message":"tx-size"}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}

	contentType := http.DetectContentType(fileContent)
	log.Printf("file contentType %s", contentType)

	utxoPrivateKeyHex := "ec6a38dd9e5b47ca1c479b86aa00f936c2c984ef7b55cf566aed792cb8cf71ee"
	destination := "tb1pz2g7y93ttzndw5sxlyqrsftqyze4sawgm5zcyra0l9chmgah3vlsqp4f8s"

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
			commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
			commitTxPrivateKeyList = append(commitTxPrivateKeyList, utxoPrivateKey)
		}
	}

	request := ord.InscriptionRequest{
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          2,
		FeeRate:                1,
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
