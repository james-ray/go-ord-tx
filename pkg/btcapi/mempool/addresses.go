package mempool

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	hcChainHash "github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcutil"
	hcWire "github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"go-ord-tx/pkg/btcapi"
	"math/big"
	"net/http"
)

type UTXO struct {
	Txid   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int64  `json:"block_time"`
	} `json:"status"`
	Value int64 `json:"value"`
}

// UTXOs is a slice of UTXO
type UTXOs []UTXO

type DOGEUTXO struct {
	Txid          string `json:"txid"`
	Vout          int    `json:"vout"`
	Height        int    `json:"height"`
	Confirmations int    `json:"confirmations"`
	ScriptPubKey  string `json:"scriptPubKey"`
	Value         string `json:"value"`
}

// DOGEUTXOs is a slice of UTXO
type DOGEUTXOs []DOGEUTXO

type HCUTXO struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	Ts            int     `json:"ts"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Height        int     `json:"height"`
	Amount        float32 `json:"amount"`
	Satoshis      int     `json:"satoshis"`
	Confirmations int     `json:"confirmations"`
}

// HCUTXOs is a slice of UTXO
type HCUTXOs []HCUTXO

func (c *MempoolClient) ListUnspent(address btcutil.Address) ([]*btcapi.UnspentOutput, error) {
	res, err := c.request(http.MethodGet, fmt.Sprintf("/address/%s/utxo", address.EncodeAddress()), nil)
	if err != nil {
		return nil, err
	}

	var utxos UTXOs
	err = json.Unmarshal(res, &utxos)
	if err != nil {
		return nil, err
	}

	unspentOutputs := make([]*btcapi.UnspentOutput, 0)
	for _, utxo := range utxos {
		txHash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}
		unspentOutputs = append(unspentOutputs, &btcapi.UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(utxo.Vout)),
			Output:   wire.NewTxOut(utxo.Value, address.ScriptAddress()),
		})
	}
	return unspentOutputs, nil
}

func (c *MempoolClient) ListUnspentForDoge(address btcutil.Address) ([]*btcapi.UnspentOutput, error) {
	restEndPoint := fmt.Sprintf("/api/v2/utxo/%s?confirmed=false", address.EncodeAddress())
	res, err := c.request(http.MethodGet, restEndPoint, nil)
	if err != nil {
		return nil, err
	}

	var utxos DOGEUTXOs
	err = json.Unmarshal(res, &utxos)
	if err != nil {
		return nil, err
	}

	unspentOutputs := make([]*btcapi.UnspentOutput, 0)
	for _, utxo := range utxos {
		txHash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}
		valueInt, succ := big.NewInt(0).SetString(utxo.Value, 10)
		if !succ {
			return nil, fmt.Errorf("invaild utxo value")
		}
		scriptPubkey, err := hex.DecodeString(utxo.ScriptPubKey)
		if err != nil {
			return nil, err
		}
		unspentOutputs = append(unspentOutputs, &btcapi.UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(utxo.Vout)),
			Output:   wire.NewTxOut(valueInt.Int64(), scriptPubkey),
		})
	}
	return unspentOutputs, nil
}

func (c *MempoolClient) ListUnspentForHc(address hcutil.Address) ([]*btcapi.UnspentOutputForHc, error) {
	restEndPoint := fmt.Sprintf("/insight/api/addr/%s/utxo", address.EncodeAddress())
	res, err := c.request(http.MethodGet, restEndPoint, nil)
	if err != nil {
		return nil, err
	}

	var utxos HCUTXOs
	err = json.Unmarshal(res, &utxos)
	if err != nil {
		return nil, err
	}

	unspentOutputs := make([]*btcapi.UnspentOutputForHc, 0)
	for _, utxo := range utxos {
		txHash, err := hcChainHash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}
		scriptPubkey, err := hex.DecodeString(utxo.ScriptPubKey)
		if err != nil {
			return nil, err
		}
		unspentOutputs = append(unspentOutputs, &btcapi.UnspentOutputForHc{
			Outpoint: hcWire.NewOutPoint(txHash, uint32(utxo.Vout), 0),
			Output:   hcWire.NewTxOut(int64(utxo.Satoshis), scriptPubkey),
		})
	}
	return unspentOutputs, nil
}
