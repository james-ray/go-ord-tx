package hdwallet

import (
	"encoding/hex"
	"go-ord-tx/blocktree/go-owcdrivers/addressEncoder"
	"go-ord-tx/tss/crypto"
	"go-ord-tx/tss/tss"
	"math/big"
	"strings"
)

func GenerateAddress(pubkeyHex string, coin string) string {
	point := strings.Split(pubkeyHex, "|")
	if len(point) == 2 {
		pkX, _ := new(big.Int).SetString(point[0], 16)
		pkY, _ := new(big.Int).SetString(point[1], 16)
		point, err := crypto.NewECPoint(tss.S256(), pkX, pkY)
		if err != nil {
			return ""
		}
		pubkeyHex = hex.EncodeToString(point.Marshal(tss.S256()))
	}

	coinType := strings.ToUpper(coin)
	if len(pubkeyHex) < 1 {
		return ""
	}
	pk, _ := hex.DecodeString(pubkeyHex)
	switch coinType {
	case "BTC":
		return addressEncoder.AddressEncode(pk, addressEncoder.BTC_mainnetAddressP2PKH)
	case "BCH":
		return addressEncoder.AddressEncode(pk, addressEncoder.BCH_mainnetAddressCash)
	case "LTC":
		return addressEncoder.AddressEncode(pk, addressEncoder.LTC_mainnetAddressP2PKH)
	case "DOGE":
		return addressEncoder.AddressEncode(pk, addressEncoder.DOGE_singleSignAddressP2PKH)
	case "DASH":
		return addressEncoder.AddressEncode(pk, addressEncoder.DASH_mainnetAddressP2PKH)
	case "BSC":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)
	case "ETH":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)
	case "POLYGON":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)

	case "HECO":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)

	case "ARBITRUM":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)

	case "OPTIMISM":
		return "0x" + addressEncoder.AddressEncode(pk, addressEncoder.ETH_mainnetPublicAddress)

	default:
		return ""
	}
}

func CheckCoinAddress(address string, coin string) bool {
	coinType := strings.ToUpper(coin)
	if len(address) < 1 {
		return false
	}
	switch coinType {

	case "BCH":
		if address[0] != '1' {
			if address[0] != 'b' {
				address = "bitcoincash:" + address
			}
		}
		result, err := addressEncoder.AddressCheck(address, coinType)
		if err != nil {
			return false
		}
		return result
	case "BSC":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "FANTOM":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "POLYGON":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "ZKSYNC":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "HECO":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "ARBITRUM":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "OPTIMISM":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "AVAX-C":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	case "HSC":
		result, err := addressEncoder.AddressCheck(address, "ETH")
		if err != nil {
			return false
		}
		return result
	default:
		result, err := addressEncoder.AddressCheck(address, coinType)
		if err != nil {
			return false
		}
		return result
	}

}
