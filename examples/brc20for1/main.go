package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"go-ord-tx/internal/ord"
	"go-ord-tx/pkg/btcapi/mempool"
	"log"
	"math/big"
	"os"
)

// base58编码
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

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

func generatePrivateKey(hexprivatekey string, compressed bool) []byte {
	versionstr := ""
	//判断是否对应的是压缩的公钥，如果是，需要在后面加上0x01这个字节。同时任何的私钥，我们需要在前方0x80的字节
	if compressed {
		versionstr = "80" + hexprivatekey + "01"
	} else {
		versionstr = "80" + hexprivatekey
	}
	//字符串转化为16进制的字节
	privatekey, _ := hex.DecodeString(versionstr)
	//通过 double hash 计算checksum.checksum他是两次hash256以后的前4个字节。
	firsthash := sha256.Sum256(privatekey)

	secondhash := sha256.Sum256(firsthash[:])

	checksum := secondhash[:4]

	//拼接
	result := append(privatekey, checksum...)

	//最后进行base58的编码
	base58result := Base58Encode(result)
	return base58result
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

func main() {
	netParams := &chaincfg.MainNetParams
	btcApiClient := mempool.NewClient(netParams, "tx")

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory, %v", err)
	}
	filePath := fmt.Sprintf("%s/examples/brc20for1/mint.txt", workingDir)
	// if file size too max will return sendrawtransaction RPC error: {"code":-26,"message":"tx-size"}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}

	//contentType := http.DetectContentType(fileContent)
	contentType := "text/plain;charset=utf-8"
	log.Printf("file contentType %s", contentType)

	//utxoPrivateKeyHex := "L5GhimdTfeBg3Cxnde3A8EGFe4JVc4ySvq3HFK3KyGgsEBSu2Cy3"
	utxoPrivateKeyHex := "ce7a1c81bae8e83d8f6dbc988cb3da509fab7a77b2a64512d28ac37bb713eddb"

	destination := "1BFCiAMrVMnVYaFDiEcyP6PvqG8kvVsgNC"

	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)

	{
		utxoPrivateKeyBytes, err := hex.DecodeString(utxoPrivateKeyHex)
		if err != nil {
			log.Fatal(err)
		}
		//utxoPrivateKeyBytes := getPrivateKeyfromWIF(utxoPrivateKeyHex)
		utxoPrivateKey, _ := btcec.PrivKeyFromBytes(utxoPrivateKeyBytes)

		pubKeyHash := btcutil.Hash160(utxoPrivateKey.PubKey().SerializeCompressed())
		pay2pubkeyHashAddress, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("pay2pubkeyHashAddress, " + pay2pubkeyHashAddress.String())
		unspentList, err := btcApiClient.ListUnspent(pay2pubkeyHashAddress)

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

	tool, err := ord.NewInscriptionToolWithBtcApiClient1(netParams, btcApiClient, &request)
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
