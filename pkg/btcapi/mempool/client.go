package mempool

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"go-ord-tx/hdwallet"
	"go-ord-tx/pkg/btcapi"
	"io"
	"log"
)

type MempoolClient struct {
	baseURL             string
	rpcURL              string
	broadcastTxEndPoint string
}

func NewClient(netParams *chaincfg.Params, broadcastTxEndPoint string) *MempoolClient {
	baseURL := ""
	if netParams.Net == wire.MainNet {
		baseURL = "https://mempool.space/api"
	} else if netParams.Net == wire.TestNet3 {
		baseURL = "https://mempool.space/testnet/api"
	} else if netParams.Net == chaincfg.SigNetParams.Net {
		baseURL = "https://mempool.space/signet/api"
	} else if netParams.Net == hdwallet.DOGEParams.Net {
		baseURL = "https://dogeblocks.com/api/v2"
	} else {
		log.Fatal("mempool don't support other netParams")
	}
	return &MempoolClient{
		baseURL:             baseURL,
		broadcastTxEndPoint: broadcastTxEndPoint,
	}
}

func NewDOGEClient() *MempoolClient {
	baseURL := ""
	baseURL = "https://dogeblocks.com"
	return &MempoolClient{
		baseURL:             baseURL,
		broadcastTxEndPoint: "sendTx",
	}
}

func NewHcClient() *MempoolClient {
	baseURL := "http://8.210.235.169"
	rpcURL := "https://8.210.235.169:14009"
	return &MempoolClient{
		baseURL:             baseURL,
		rpcURL:              rpcURL,
		broadcastTxEndPoint: "sendrawtransaction",
	}
}

func (c *MempoolClient) request(method, subPath string, requestBody io.Reader) ([]byte, error) {
	return btcapi.Request(method, c.baseURL, subPath, requestBody)
}

func (c *MempoolClient) requestWithAuth(method, subPath string, requestBody string) ([]byte, error) {
	return btcapi.RequestWithAuth(method, c.rpcURL, subPath, requestBody)
}

var _ btcapi.BTCAPIClient = (*MempoolClient)(nil)
