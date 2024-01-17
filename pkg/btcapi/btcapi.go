package btcapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	hcWire "github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type UnspentOutput struct {
	Outpoint *wire.OutPoint
	Output   *wire.TxOut
}

type UnspentOutputForHc struct {
	Outpoint *hcWire.OutPoint
	Output   *hcWire.TxOut
}

type BTCAPIClient interface {
	GetRawTransaction(txHash *chainhash.Hash) (*wire.MsgTx, error)
	BroadcastTx(tx *wire.MsgTx) (*chainhash.Hash, error)
	ListUnspent(address btcutil.Address) ([]*UnspentOutput, error)
}

func Request(method, baseURL, subPath string, requestBody io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", baseURL, subPath)
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	return body, nil
}

func RequestWithAuth(method, baseURL, subPath string, requestBody string) ([]byte, error) {
	url := baseURL
	j, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  subPath,
		"params":  []interface{}{requestBody},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(j))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Accept", "application/json")D
	// Configure basic access authorization.
	req.SetBasicAuth("admin", "123")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	return body, nil
}
