package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"log"
)

func main() {
	netParams := &chaincfg.MainNetParams
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyHex := hex.EncodeToString(privateKey.Serialize())
	log.Printf("new priviate key %s \n", privateKeyHex)

	taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(privateKey.PubKey())), netParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("new taproot address %s %s \n", taprootAddress, taprootAddress.EncodeAddress())

	restorePrivateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	restorePrivateKey, _ := btcec.PrivKeyFromBytes(restorePrivateKeyBytes)

	restoreTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(restorePrivateKey.PubKey())), netParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("restore taproot address %s \n", restoreTaprootAddress.EncodeAddress())

	if taprootAddress.EncodeAddress() != restoreTaprootAddress.EncodeAddress() {
		log.Fatal("restore privateKey error")
	}

	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addressPubkeyHash, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("new addressPubkeyHash address %s %s \n", addressPubkeyHash, addressPubkeyHash.EncodeAddress())
	restorePrivateKey, _ = btcec.PrivKeyFromBytes(restorePrivateKeyBytes)

	pubKeyHash = btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	restoreP2PKHashAddress, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("restore restoreP2PKHashAddress address %s \n", restoreP2PKHashAddress.EncodeAddress())
	/**
	test btc faucet
	https://signetfaucet.com/
	https://alt.signetfaucet.com/
	*/
}
