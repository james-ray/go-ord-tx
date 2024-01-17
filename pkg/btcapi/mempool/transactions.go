package mempool

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	hcwire "github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func (c *MempoolClient) GetRawTransaction(txHash *chainhash.Hash) (*wire.MsgTx, error) {
	res, err := c.request(http.MethodGet, fmt.Sprintf("/tx/%s/raw", txHash.String()), nil)
	if err != nil {
		return nil, err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	if err := tx.Deserialize(bytes.NewReader(res)); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *MempoolClient) BroadcastTx(tx *wire.MsgTx) (*chainhash.Hash, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPost, fmt.Sprintf("/%s", c.broadcastTxEndPoint), strings.NewReader(hex.EncodeToString(buf.Bytes())))
	if err != nil {
		return nil, err
	}

	txHash, err := chainhash.NewHashFromStr(string(res))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse tx hash, %s", string(res)))
	}
	return txHash, nil
}

func (c *MempoolClient) BroadcastTxForHc(tx *hcwire.MsgTx) (*chainhash.Hash, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, err
	}

	resp, err := c.requestWithAuth(http.MethodPost, c.broadcastTxEndPoint, hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return nil, fmt.Errorf("json.unmarshal err: %v, resp=%s", err, resp)
	}
	txHash, err := chainhash.NewHashFromStr(res["result"].(string))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse tx hash, %v", res))
	}
	return txHash, nil
}
