package main

import (
	"encoding/hex"
	"go-ord-tx/pkg/btcapi/mempool"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempool.NewClient(netParams, "tx")

	privateKeyHex := "ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddb"
	privateKeyByte, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyByte)

	log.Printf("new priviate key %s \n", privateKeyHex)

	//taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())), netParams)
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addressPubkeyHash, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatal(err)
	}

	pay2pkHashAddr := addressPubkeyHash.EncodeAddress()
	//log.Printf("new taproot address %s \n", taprootAddress.EncodeAddress())
	log.Printf("new addressPubkeyHash address %s \n", pay2pkHashAddr)
	unspentList, err := btcApiClient.ListUnspent(addressPubkeyHash)
	commitTxOutPointList := make([]*wire.OutPoint, 0)

	if err != nil {
		log.Fatalf("list unspent err %v", err)
	}
	for i := range unspentList {
		for j := i + 1; j < len(unspentList); j++ {
			if unspentList[i].Output.Value > unspentList[j].Output.Value { //转账，为了转移铭文，优先选小的，即包含铭文的utxo
				temp := unspentList[i]
				unspentList[i] = unspentList[j]
				unspentList[j] = temp
			}
		}
	}
	for i := range unspentList {
		log.Printf("unspentList hash %s \n", unspentList[i].Outpoint.Hash.String())
		log.Printf("unspentList value %d \n", unspentList[i].Output.Value)
		commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)

	}
	msgTx, prevOutFetcher, err := BuildCommitTx(commitTxOutPointList, "1BFCiAMrVMnVYaFDiEcyP6PvqG8kvVsgNC", "1AkPTypiuVm4ztRkuG16rSeu5nfY2ETRYP", 200000, 80)
	//msgTx, prevOutFetcher, err := BuildSendAllTx(commitTxOutPointList, "bc1pghwgycut83fufa8rlm80pt5a9xt0vnew5cr7z8ceshwf5efavnxset9gk6", "bc1pl0c4elws88jynnn5s9ll2wtc9yu05ndlae8va4se49jpsqa4stxs9ljhc4", 80)
	if err != nil {
		log.Fatalf("list unspent err %v", err)
	}
	msgTx, err = signCommitTx1(privateKey, msgTx, prevOutFetcher)
	if err != nil {
		log.Fatalf("list unspent err %v", err)
	}
	hash, err := SendRawTransaction(btcApiClient, msgTx)
	if err != nil {
		log.Fatalf("list unspent err %v", err)
	}
	log.Printf("txhash %s \n", hash.String())

}
