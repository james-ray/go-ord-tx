package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	hcchaincfg "github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcd/btcec/v2"
	"go-ord-tx/hdwallet"
	"go-ord-tx/pkg/btcapi/mempool"
	"log"
	"math/big"
)

// base58编码
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

var HcParams = &hcchaincfg.MainNetParams

const (
	DogeMinVal      = 100000    //doge utxo 最小val
	DogeDeployFee   = 100000000 //部署 1个Doge 币
	DogeMintFee     = 10000000  //铸币 0.1 个 Doge 币
	DogeTransferFee = 100000    //转账 0.001个 Doge 币
)

type OutPutItem struct {
	TxHash   string `json:"tx_hash" form:"tx_hash"`
	Vout     uint32 `json:"vout" form:"vout"`
	Value    int64  `json:"value" form:"value"`
	Pkscript string `json:"pkscript" form:"pkscript"`
}

type SignInput struct {
	PrivateKey    string //私钥
	Decimal       int
	FeeCoin       string
	Node          string
	LgrSeq        int64
	SecondFee     int64
	Satoshi       int64
	CommitSatoshi int64
	Coin          string // 主链币
	Symbol        string // symbol
	Amount        int64  //转账数量
	LargeAmount   string // 转账数量
	Change        int64  //找零数量
	Fee           int64  //交易费用
	GasLimit      int64  // gas数量
	GasPrice      int64  // gas价格
	Type          string //交易类型 //xtz TYPE = branch
	SrcAddr       string //转账地址
	DestAddr      string //接受地址
	ContractAddr  string //合约地址
	Sequence      int64  // 序列号
	Memo          string //交易备注
	Inputs        []byte //Vin构造
	Params        []byte //预留字段
	OutPutItem    []byte
}

type SignTxHashInput struct {
	Coin      string // 主链币
	Symbol    string // symbol
	Signature string //签名数据
	TxHash    string //交易Hash
	TxRawHex  string // 原始交易RawHex
	Params    []byte //预留字段
}

type DrcAddressUnit struct {
	Address string `json:"address" form:"address"` //接收地址
	Value   int64  `json:"value" form:"value"`     //转账数量
	TaprootParam
}

type Drc20Param struct {
	Amount        string `json:"amount"`
	Token         string `json:"token"`
	Op            string `json:"op"`
	CommitFeeRate int64  `json:"commitFeeRate"`
	FeeRate       int64  `json:"feeRate"`

	P    string `json:"p"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`

	Max            string `json:"max"`
	MintNeedAmount int64
	Lim            string `json:"lim"`
	Dec            string `json:"dec"`
}

type TxHashResult struct {
	ResCode  int    // 0 失败 1 成功
	Coin     string // 主链币
	Symbol   string // symbol币种
	RawTX    string //签名后的数据
	TxHash   string //交易Hash
	TxRawHex string //交易TxRawHex
	ErrMsg   string // 失败原因(便于排查问题,不是必定返回)
	ErrCode  int    //错误码(暂时保留)

	Params []byte //预留字段
}

type SignResult struct {
	Res     int    // 0 失败 1 成功
	Coin    string // 主链币
	Symbol  string // symbol币种
	RawTX   string //签名后的数据
	Txs     []*wire.MsgTx
	TxHash  string // 交易Hash
	ErrMsg  string // 失败原因(便于排查问题,不是必定返回)
	ErrCode int    //错误码
	Params  []byte //预留字段
}

type FeeResult struct {
	ResCode         int    // 0 失败 1 成功
	ErrMsg          string // 失败原因(便于排查问题,不是必定返回)
	ErrCode         int    //错误码
	Fee             int64  // 成功: 返回TxID
	Coin            string // 币种
	Symbol          string // symbol
	FeeIndex        int64
	Change          int64
	Bytes           int64
	Params          []byte //预留字段
	NetworkFee      int64
	SystemFee       int64
	ValidUntilBlock int64
	SecondFee       int64
}

type WalletAccount struct {
	Res        int    // 0 失败 1 成功
	Address    string // 成功必定包含地址
	PublicKey  string // 公钥
	PrivateKey string // 私钥
	Seed       string //根Seed

	Coin    string //币种
	ErrMsg  string // 失败原因(便于排查问题,不是必定返回)
	ErrCode int    //错误码
	Params  string //预留字段
}

type AddressUnit struct {
	Address string `json:"address" form:"address"` //接收地址
	Value   int64  `json:"value" form:"value"`     //转账数量
	TaprootParam
}

type TaprootParam struct {
	Amount        string `json:"amount"`
	Token         string `json:"token"`
	Op            string `json:"op"`
	CommitFeeRate int64  `json:"commitFeeRate"`
	FeeRate       int64  `json:"feeRate"`

	P    string `json:"p"`
	Tick string `json:"tick"`
	Amt  string `json:"amt"`

	Max string `json:"max"`
	Lim string `json:"lim"`
	Dec string `json:"dec"`
}

type drc20 struct {
	name   string
	symbol string
	key    *hdwallet.Key
}

/*
	func (c *drc20) Fee(signIn *SignInput) (*FeeResult, error) {
		//TODO implement me
		panic("implement me")
	}

	func newDrc20(key *Key) Wallet {
		return &drc20{
			name:   "drc20",
			symbol: "drc20",
			key:    key,
		}
	}
*/
func (c *drc20) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *drc20) GetName() string {
	return c.name
}

func (c *drc20) GetSymbol() string {
	return c.symbol
}

func (c *drc20) GetKey() *hdwallet.Key {
	return c.key
}

func (c *drc20) GetAddress() (string, error) {
	return "", nil
}

func (c *drc20) CreateRawTransaction(signIn *SignInput) (*SignResult, error) {
	return &SignResult{
		Res: 0,
	}, nil
}
func (c *drc20) GenerateTxHash(signIn *SignInput) (*TxHashResult, error) {
	return &TxHashResult{}, nil
}
func (c *drc20) SignTxHash(signIn *SignTxHashInput) (*TxHashResult, error) {
	return &TxHashResult{}, nil
}

// secondFee 351 * fee
func (c *drc20) SignRawTransaction(signIn *SignInput) (*SignResult, error) {

	//1. 先生成一个地址
	var vins []OutPutItem
	var vouts []OutPutItem

	/*	a, err := hcutil.DecodeWIF(signIn.PrivateKey)
		if err != nil {
			return nil, err
		}*/
	var centerAddress string
	//priv := a.PrivKey
	priv, err := hex.DecodeString(signIn.PrivateKey)
	if err != nil {
		return nil, err
	}
	a, pub := chainec.Secp256k1.PrivKeyFromBytes(priv)
	json.Unmarshal(signIn.Inputs, &vins)
	fmt.Println("vins : ", vins)

	var param AddressUnit
	var redeemScript []byte
	var transferValue int64
	var platformValue int64
	if signIn.Type == "drc20" {
		//err := json.Unmarshal(signIn.Params, &param)
		//if err != nil {
		//	return nil, err
		//}
		var addressUnit []AddressUnit
		err = json.Unmarshal(signIn.OutPutItem, &addressUnit)
		if err != nil {
			return nil, err
		}
		if len(addressUnit) > 0 {
			param = addressUnit[0]
		}
		signIn.DestAddr = param.Address
		var jsonData string
		//{"dec":18,"lim":"10","max":"5000","op":"deploy","p":"drc-20","tick":"DAMN"}
		switch param.Op {
		case "deploy": //部署 100个Doge 币
			platformValue = DogeDeployFee
			//添加第一个out 代币地址 最小额度
			transferValue = platformValue + DogeMinVal
			jsonData = fmt.Sprintf(`{"p":"drc-20","op":"%s","tick":"%s","max":"%s","lim":"%s"}`, param.Op, param.Tick, param.Max, param.Lim)
		case "mint": //铸币 0.5 个 Doge 币
			platformValue = DogeMintFee
			//添加第一个out 代币地址 最小额度
			transferValue = platformValue + DogeMinVal
			jsonData = fmt.Sprintf(`{"p":"drc-20","op":"%s","tick":"%s","amt":"%s"}`, param.Op, param.Tick, param.Amt)
		case "transfer": //转账 0.001个 Doge 币
			transferValue = DogeTransferFee
			jsonData = fmt.Sprintf(`{"p":"drc-20","op":"%s","tick":"%s","amt":"%s"}`, param.Op, param.Tick, param.Amt)

		default:
			return nil, errors.New("not support operation")
		}

		signIn.SecondFee = signIn.Fee
		builder := txscript.NewScriptBuilder()
		builder.AddOp(txscript.OP_1).AddData(pub.SerializeCompressed()).AddOp(txscript.OP_1)
		builder.AddOp(txscript.OP_CHECKMULTISIG)
		builder.AddData([]byte("ord")).AddData([]byte("text/plain;charset=utf-8")).AddData([]byte(jsonData))
		builder.AddOp(txscript.OP_DROP).AddOp(txscript.OP_DROP).AddOp(txscript.OP_DROP).AddOp(txscript.OP_DROP)
		redeemScript, err = builder.Script()
		if err != nil {
			return nil, err
		}

		fmt.Println("the jsonData is ", jsonData)
		addr_new, err := hcutil.NewAddressScriptHash(redeemScript, HcParams)
		if err != nil {
			return nil, err
		}
		fmt.Println("the address is ", addr_new.String())
		centerAddress = addr_new.String()
	}
	rawTx := ""
	//2.向这个地址打钱的交易
	centerAddress1, err := hcutil.DecodeAddress(centerAddress)
	if err != nil {
		return nil, err
	}
	mtxOne := wire.NewMsgTx()
	centerAddressScript, err := txscript.PayToAddrScript(centerAddress1)
	if err != nil {
		return nil, err
	}

	//添加第二笔交易的手续费
	transferValue = transferValue + signIn.SecondFee
	fmt.Println("the second fee is", signIn.SecondFee)

	outputToCenter := &wire.TxOut{
		Value:    transferValue,
		PkScript: centerAddressScript,
	}

	mtxOne.AddTxOut(outputToCenter)

	srcAddr, err := hcutil.DecodeAddress(signIn.SrcAddr)
	if err != nil {
		return nil, err
	}

	srcAddrScript, err := txscript.PayToAddrScript(srcAddr)

	if err != nil {
		return nil, err
	}
	destAddr, err := hcutil.DecodeAddress(signIn.DestAddr)
	if err != nil {
		return nil, err
	}
	destAddrScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return nil, err
	}
	var spendValue int64 = 0
	for _, input := range vins {
		txHash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return nil, fmt.Errorf("txid error")
		}
		prevOut := wire.NewOutPoint(txHash, input.Vout, 0)
		txIn := wire.NewTxIn(prevOut, []byte{})
		mtxOne.AddTxIn(txIn)
		vouts = append(vouts, input)
		spendValue = spendValue + input.Value
		if spendValue >= transferValue+signIn.Fee { //the second fee is already included in transferValue
			break
		}
	}

	//添加找零
	if spendValue > transferValue+signIn.Fee {
		changeOutput := &wire.TxOut{
			Value:    spendValue - transferValue - signIn.Fee, //change for the first tx
			PkScript: srcAddrScript,
		}
		mtxOne.AddTxOut(changeOutput)
	}

	for i, input := range vouts {

		//var setSignatureScript []byte
		txInPkScript, err := hex.DecodeString(input.Pkscript)
		if err != nil {
			return nil, err
		}
		// 获取vin的script的类型
		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(0, txInPkScript, HcParams)
		if err != nil {
			return nil, err
		}
		switch scriptClass {
		case txscript.PubKeyHashTy:
			script, err := txscript.SignatureScript(
				mtxOne,
				i,
				txInPkScript,
				txscript.SigHashAll,
				a,
				true,
			)

			if err != nil {
				return nil, err
			}
			mtxOne.TxIn[i].SignatureScript = script

		default:
			return nil, errors.New("unsupport script")
		}
	}

	needSpendHash := mtxOne.TxHash()

	fmt.Println("the needSpendHash is ", needSpendHash.String())
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, mtxOne.SerializeSize()))
	if err := mtxOne.Serialize(buf); err != nil {
		return nil, err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	fmt.Printf("commitTx size:%v\r\n", mtxOne.SerializeSize())

	rawTx += txHex
	rawTx += ","

	//3. 将这个地址打的钱 消费到指定地址上去
	payToAddressStr := "Hsao5VNXao7DZwDBHdkLwfN7dL8tEsTxkUX"
	payToAddressAddr, err := hcutil.DecodeAddress(payToAddressStr)

	if err != nil {
		return nil, err
	}

	payToAddrScript, err := txscript.PayToAddrScript(payToAddressAddr)
	if err != nil {
		return nil, err
	}

	//合约脚本
	mtxTwo := wire.NewMsgTx()

	//添加out给接收方 固定0.001 个 Doge
	outputToDest := &wire.TxOut{
		Value:    DogeMinVal,
		PkScript: destAddrScript,
	}
	mtxTwo.AddTxOut(outputToDest)

	if platformValue > 0 {
		//给项目方
		outPutToFee := &wire.TxOut{
			Value:    platformValue,
			PkScript: payToAddrScript,
		}
		mtxTwo.AddTxOut(outPutToFee)
	}

	prevOut1 := wire.NewOutPoint(&needSpendHash, 0, 0) // outputToCenter is the 0th vout
	txIn1 := wire.NewTxIn(prevOut1, []byte{})
	mtxTwo.AddTxIn(txIn1)
	sig, err := txscript.RawTxInSignature(mtxTwo, 0, redeemScript, txscript.SigHashAll, a)
	if err != nil {
		return nil, err
	}
	signature := txscript.NewScriptBuilder()
	signature.AddOp(txscript.OP_10).AddOp(txscript.OP_FALSE).AddData(sig)
	signature.AddData(redeemScript)
	signatureScript, err := signature.Script()
	if err != nil {
		return nil, err
	}
	mtxTwo.TxIn[0].SignatureScript = signatureScript

	buf2 := bytes.NewBuffer(make([]byte, 0, mtxTwo.SerializeSize()))
	if err := mtxTwo.Serialize(buf2); err != nil {
		return nil, err
	}

	txHex2 := hex.EncodeToString(buf2.Bytes())
	fmt.Printf("revealTx size:%v\r\n", mtxTwo.SerializeSize())
	//fmt.Println(mtxTwo.SerializeSize())
	rawTx += txHex2

	txs := make([]*wire.MsgTx, 2)
	txs[0] = mtxOne
	txs[1] = mtxTwo
	return &SignResult{
		Res:   1,
		RawTX: rawTx,
		Txs:   txs,
	}, nil
}

func (c *drc20) GetWalletAccountFromWif() (*WalletAccount, error) {
	wif := c.GetKey().Wif
	if len(wif) > 0 {
		btcWif, err := hcutil.DecodeWIF(wif)
		if err != nil {
			fmt.Println("Wif err : ", err.Error())
			return nil, err
		}
		//netID: btc 128 ltc 176
		isHc := btcWif.IsForNet(HcParams)
		if isHc == false {
			return nil, errors.New("key type error")
		}
		pk := btcWif.SerializePubKey()
		fmt.Println("pk : ", hex.EncodeToString(pk))
		address, err := hcutil.NewAddressPubKeyHash(hcutil.Hash160(pk), HcParams, chainec.ECTypeSecp256k1)
		if err != nil {
			fmt.Println("Wif err : ", err.Error())
			return nil, err

		}
		btcAddress := address.EncodeAddress()

		return &WalletAccount{
			Res:        1,
			PrivateKey: wif,
			PublicKey:  hex.EncodeToString(pk),
			Address:    btcAddress,
		}, nil
	}
	return &WalletAccount{
		Res: 0,
	}, nil
}

func (c *drc20) GetWalletAccount() *WalletAccount {
	if c.GetKey().Extended == nil {
		return &WalletAccount{
			Res: 0,
		}
	}
	// fmt.Println("Mnemonic = ", c.GetKey().Mnemonic)

	// fmt.Println("Seed = ", c.GetKey().Seed)
	address, err := c.GetAddress()
	if err != nil {
		return &WalletAccount{
			Res:    0,
			ErrMsg: err.Error(),
		}
	}
	pWif, err1 := c.GetKey().PrivateWIF(true)
	if err1 != nil {
		return &WalletAccount{
			Res:    0,
			ErrMsg: err1.Error(),
		}
	}
	publicKey := c.GetKey().PublicHex(true)

	return &WalletAccount{
		Res:        1,
		PrivateKey: pWif,
		Address:    address,
		PublicKey:  publicKey,
		Seed:       c.GetKey().Seed,
	}
}

// success txhash
// https://explorer.viawallet.com/doge/tx/fd5293454d63015bb22ef9d84055b0124a5155c5884f2f663d3dcfdb54f719b1
// https://explorer.viawallet.com/doge/tx/3a70515f0bfe1359a87d19165592564c6b155f86fed510e57ff90a92b7c0d8c2
// 3f6496a5f40f00936af2e39130dba68815b58fd7d5f6362bf68e09fc13c04d99
// {"dec":18,"lim":"10","max":"5000","op":"deploy","p":"drc-20","tick":"DAMN"}
// 77835ee83c655f7b82f1e864f424d3df2dbd7020088fb3b5ac107c9bcb7b9753
// 83c0665369dd588c9770ccc0a0ffa080b3c01413b9bb248e0666dcf64eaa9778
func main() {
	initParams()
	btcApiClient := mempool.NewHcClient()

	//contentType := http.DetectContentType(fileContent)
	contentType := "text/plain;charset=utf-8"
	log.Printf("file contentType %s", contentType)

	//utxoPrivateKeyWif := "QNykotPn5HQs8GkWP1ZrnP3uZornUbL4eQFB64FrmEvxhdiNp8h4"
	utxoPrivateKeyHex := "ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddb" //HsQt7gFhtoiSKq8UwH6Sa9i6VA2ivm5jS6H
	//"ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddc"    HsMpSkPvW39iGhUgayHu3v7p8SWoVAfyMFR
	//"ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddd"    Hsao5VNXao7DZwDBHdkLwfN7dL8tEsTxkUX
	//"ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713edde"    HsSDD4XmJ6YDoUMHqX3db7Wur1jEyfRExX5
	//"ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddf"    HsLfRbdKLBaCmJUTN5YysnSasxvaKEudxPG

	utxoPrivateKeyBytes, err := hex.DecodeString(utxoPrivateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	//utxoPrivateKeyBytes := getPrivateKeyfromWIF(utxoPrivateKeyWif)
	utxoPrivateKey, _ := btcec.PrivKeyFromBytes(utxoPrivateKeyBytes)

	pubKeyHash := hcutil.Hash160(utxoPrivateKey.PubKey().SerializeCompressed())
	pay2pubkeyHashAddress, err := hcutil.NewAddressPubKeyHash(pubKeyHash, HcParams, chainec.ECTypeSecp256k1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("pay2pubkeyHashAddress, " + pay2pubkeyHashAddress.String())
	unspentList, err := btcApiClient.ListUnspentForHc(pay2pubkeyHashAddress)

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
	totalAmount := int64(0)
	outputs := make([]OutPutItem, 0)
	for i := range unspentList {
		item1 := OutPutItem{
			TxHash:   unspentList[i].Outpoint.Hash.String(),
			Value:    unspentList[i].Output.Value,
			Vout:     unspentList[i].Outpoint.Index,
			Pkscript: hex.EncodeToString(unspentList[i].Output.PkScript),
		}
		outputs = append(outputs, item1)
		totalAmount += unspentList[i].Output.Value
	}

	param := TaprootParam{
		Amt:           "100",
		Tick:          "DRC2",
		Op:            "deploy",
		CommitFeeRate: 18,
		FeeRate:       19,
	}

	// fmt.Println("ltc outputs: ", outputs)

	jsonInputs, err := json.Marshal(outputs)
	if err != nil {
		//log.Fatal("Cannot encode to JSON ", err)
		fmt.Println("outputs err: ", err.Error())

	}
	// fmt.Println("ltc outputs: ", jsonInputs)

	marshal, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err)
	}
	//[{"p":"drc-20","op":"transfer","tick":"DRC2","amt":"900","address":"D92dNTMsSWJRBoUPjv2n8LMozrEFZ1NNRJ"}]
	packs := make([]AddressUnit, 1)
	packs[0] = AddressUnit{
		Address: pay2pubkeyHashAddress.String(),
		TaprootParam: TaprootParam{
			P:    "drc-20",
			Amt:  "100",
			Tick: "DRC1",
			Op:   "deploy",
			Max:  "10000000000000",
			Lim:  "100000000",
			Dec:  "8",
		},
	}
	marshall2, err := json.Marshal(packs)
	if err != nil {
		fmt.Println(err)
	}
	signInput := &SignInput{
		Coin:       "drc20",
		Symbol:     "drc20",
		PrivateKey: utxoPrivateKeyHex,
		SrcAddr:    pay2pubkeyHashAddress.String(),
		DestAddr:   pay2pubkeyHashAddress.String(),
		Fee:        20000000,
		//SecondFee:  20000000,
		Amount:     totalAmount,
		Change:     0,
		Inputs:     jsonInputs,
		Type:       "drc20",
		Params:     marshal,
		OutPutItem: marshall2,
	}
	drc20v := drc20{
		name:   "p1p1",
		symbol: "p1p1",
		key:    nil,
	}
	tranferResult, err := drc20v.SignRawTransaction(signInput) //fmt.Println(tranferResult.RawTX)
	if err != nil {
		panic(err)
	}
	fmt.Println("rawTx--->", tranferResult.RawTX)

	/*txHash, err := btcApiClient.BroadcastTx(tranferResult.Txs[0])
	if err != nil {
		panic(err)
	}
	fmt.Println("txHash1 --->", txHash)

	txHash, err = btcApiClient.BroadcastTx(tranferResult.Txs[1])
	if err != nil {
		panic(err)
	}
	fmt.Println("txHash2 --->", txHash)*/
	//https://www.oklink.com/doge/tx/8d15218e949db0077a327ba2d9abb061dc1cdea04be3548e0860c6bdf3b61f23
}

// 检查checkWIF是否有效
func checkWIF(wifprivate string) bool {
	rawdata := []byte(wifprivate)
	//包含了80、私钥、checksum
	base58decodedata := Base58Decode(rawdata)

	fmt.Printf("base58decodedata：%x\n", base58decodedata)
	length := len(base58decodedata)

	if length < 37 {
		fmt.Printf("长度小于37，一定有问题")
		return false
	}

	private := base58decodedata[:(length - 4)]
	//得到检查码
	//fmt.Printf("private：%x\n",private)
	firstsha := sha256.Sum256(private)

	secondsha := sha256.Sum256(firstsha[:])

	checksum := secondsha[:4]
	//fmt.Printf("%x\n",checksum)
	//得到原始的检查码
	orignchecksum := base58decodedata[(length - 4):]
	//	fmt.Printf("%x\n",orignchecksum)

	//[]byte对比
	if bytes.Compare(checksum, orignchecksum) == 0 {
		return true
	}

	return false

}

// 通过wif格式的私钥，得到原始的私钥。
func getPrivateKeyfromWIF(wifprivate string) []byte {
	if checkWIF(wifprivate) {
		rawdata := []byte(wifprivate)
		//包含了80、私钥、checksum
		base58decodedata := Base58Decode(rawdata)
		//私钥一共32个字节，排除了0x80
		return base58decodedata[1:33]
	}
	return []byte{}

}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b == '1' {
			zeroBytes++
		} else {
			break
		}
	}

	payload := input[zeroBytes:]

	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b) //反推出余数

		result.Mul(result, big.NewInt(58)) //之前的结果乘以58

		result.Add(result, big.NewInt(int64(charIndex))) //加上这个余数

	}

	decoded := result.Bytes()

	decoded = append(bytes.Repeat([]byte{0x00}, zeroBytes), decoded...)
	return decoded
}

func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)

	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 对x取余数
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)

	for _, b := range input {

		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result

}

// 字节数组的反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
