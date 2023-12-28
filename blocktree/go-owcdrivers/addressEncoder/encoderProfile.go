package addressEncoder

var (
	BTCAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	XMRAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	OntAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	XRPAlphabet        = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	ZECAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	BTCBech32Alphabet  = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	BTMBech32Alphabet  = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	LTCAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	LTCBech32Alphabet  = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	BCHLegacyAlphabet  = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	BCHCashAlphabet    = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	XTZAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	HCAlphabet         = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	QTUMAlphabet       = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	DCRDAlphabet       = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	NASAlphabet        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	TRONAlphabet       = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	VSYSAlphabet       = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	ATOMBech32Alphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

//func (at *AddressType) Prefix() []byte {
//	return at.Prefix
//}

var (
	//BTC stuff
	BTC_mainnetAddressP2PKH         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x00}, nil}
	BTC_mainnetAddressP2SH          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x05}, nil}
	BTC_mainnetAddressBech32V0      = AddressType{"bech32", BTCBech32Alphabet, "bc", "h160", 20, []byte{0}, nil}
	BTC_mainnetPrivateWIF           = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 32, []byte{0x80}, nil}
	BTC_mainnetPrivateWIFCompressed = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 32, []byte{0x80}, []byte{0x01}}
	BTC_mainnetPublicBIP32          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x88, 0xB2, 0x1E}, nil}
	BTC_mainnetPrivateBIP32         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x88, 0xAD, 0xE4}, nil}
	BTC_testnetAddressP2PKH         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x6F}, nil}
	BTC_testnetAddressP2SH          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0xC4}, nil}
	BTC_testnetAddressBech32V0      = AddressType{"bech32", BTCBech32Alphabet, "tb", "h160", 20, []byte{0}, nil}
	BTC_testnetPrivateWIF           = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 32, []byte{0xEF}, nil}
	BTC_testnetPrivateWIFCompressed = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 32, []byte{0xEF}, []byte{0x01}}
	BTC_testnetPublicBIP32          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x35, 0x87, 0xCF}, nil}
	BTC_testnetPrivateBIP32         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x35, 0x83, 0x94}, nil}

	//XMR stuff
	XMR_mainnetPublicAddress           = AddressType{"XMR", XMRAlphabet, "keccak256", "", 64, []byte{0x12}, nil}
	XMR_mainnetPublicSubAddress        = AddressType{"XMR", XMRAlphabet, "keccak256", "", 64, []byte{0x2A}, nil}
	XMR_mainnetPublicIntegratedAddress = AddressType{"XMR", XMRAlphabet, "keccak256", "payID", 72, []byte{0x13}, nil}
	XMR_testnetPublicAddress           = AddressType{"XMR", XMRAlphabet, "keccak256", "", 64, []byte{0x35}, nil}
	XMR_testnetPublicSubAddress        = AddressType{"XMR", XMRAlphabet, "keccak256", "", 64, []byte{0x3f}, nil}
	XMR_testnetPublicIntegratedAddress = AddressType{"XMR", XMRAlphabet, "keccak256", "payID", 72, []byte{0x36}, nil}

	//DASH stuff
	DASH_mainnetAddressP2PKH = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x4c}, nil}
	//DOGE stuff
	DOGE_singleSignAddressP2PKH = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x1e}, nil}
	DOGE_multiSignAddressP2PKH  = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x16}, nil}
	//ONT stuff
	ONT_Address = AddressType{"base58", OntAlphabet, "doubleSHA256", "h160", 20, []byte{0x17}, nil}
	//XRP stuff
	XRP_Address = AddressType{"base58", XRPAlphabet, "doubleSHA256", "h160", 20, []byte{0x00}, nil}
	//BTM stuff
	BTM_mainnetAddressBech32V0 = AddressType{"bech32", BTCBech32Alphabet, "bm", "h160", 20, []byte{0}, nil}
	BTM_testnetAddressBech32V0 = AddressType{"bech32", BTCBech32Alphabet, "tm", "h160", 20, []byte{0}, nil}
	//ZEC stuff
	ZEC_mainnet_t_AddressP2PKH = AddressType{"base58", ZECAlphabet, "doubleSHA256", "h160", 20, []byte{0x1C, 0xB8}, nil}
	ZEC_mainnet_t_AddressP2SH  = AddressType{"base58", ZECAlphabet, "doubleSHA256", "h160", 20, []byte{0x1C, 0xBD}, nil}
	ZEC_testnet_t_AddressP2PKH = AddressType{"base58", ZECAlphabet, "doubleSHA256", "h160", 20, []byte{0x1D, 0x25}, nil}
	ZEC_testnet_t_AddressP2SH  = AddressType{"base58", ZECAlphabet, "doubleSHA256", "h160", 20, []byte{0x1C, 0xBA}, nil}

	//LTC stuff
	LTC_mainnetAddressP2PKH         = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x30}, nil}
	LTC_mainnetAddressP2SH          = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x05}, nil}
	LTC_mainnetAddressP2SH2         = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x32}, nil}
	LTC_mainnetAddressBech32V0      = AddressType{"bech32", LTCBech32Alphabet, "ltc", "h160", 20, []byte{0}, nil}
	LTC_mainnetPrivateWIF           = AddressType{"base58", LTCAlphabet, "doubleSHA256", "", 32, []byte{0xB0}, nil}
	LTC_mainnetPrivateWIFCompressed = AddressType{"base58", LTCAlphabet, "doubleSHA256", "", 32, []byte{0xB0}, []byte{0x01}}
	LTC_mainnetPublicBIP32          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x88, 0xB2, 0x1E}, nil}
	LTC_mainnetPrivateBIP32         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x88, 0xAD, 0xE4}, nil}
	LTC_testnetAddressP2PKH         = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x6F}, nil}
	LTC_testnetAddressP2SH          = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0xC4}, nil}
	LTC_testnetAddressP2SH2         = AddressType{"base58", LTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x3A}, nil}
	LTC_testnetAddressBech32V0      = AddressType{"bech32", LTCBech32Alphabet, "tltc", "h160", 20, []byte{0}, nil}
	LTC_testnetPrivateWIF           = AddressType{"base58", LTCAlphabet, "doubleSHA256", "", 32, []byte{0xEF}, nil}
	LTC_testnetPrivateWIFCompressed = AddressType{"base58", LTCAlphabet, "doubleSHA256", "", 32, []byte{0xEF}, []byte{0x01}}
	LTC_testnetPublicBIP32          = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x35, 0x87, 0xCF}, nil}
	LTC_testnetPrivateBIP32         = AddressType{"base58", BTCAlphabet, "doubleSHA256", "", 74, []byte{0x04, 0x35, 0x83, 0x94}, nil}

	// BSV stuff
	BSV_mainnetAddressP2PKH = AddressType{"base58", BTCAlphabet, "doubleSHA256", "h160", 20, []byte{0x00}, nil}

	//BCH stuff
	BCH_mainnetAddressLegacy = AddressType{"base58", BCHLegacyAlphabet, "doubleSHA256", "h160", 20, []byte{0x00}, nil}
	BCH_mainnetAddressCash   = AddressType{"base32PolyMod", BCHCashAlphabet, "bitcoincash", "h160", 21, nil, nil}

	//ETH stuff
	ETH_mainnetPublicAddress = AddressType{"eip55", "", "", "keccak256", 32, nil, nil}

	//DCRD stuff
	DCRD_mainnetAddressP2PKH      = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x07, 0x3f}, nil} //PubKeyHashAddrID, stars with Ds
	DCRD_mainnetAddressP2PK       = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x13, 0x86}, nil} //PubKeyAddrID,stars with Dk
	DCRD_mainnetAddressPKHEdwards = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x07, 0x1f}, nil} //PKHEdwardsAddrID,starts with De
	DCRD_mainnetAddressPKHSchnorr = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x07, 0x01}, nil} //PKHSchnorrAddrID,starts with DS
	DCRD_mainnetAddressP2SH       = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x07, 0x1a}, nil} //ScriptHashAddrID,starts with Dc
	DCRD_mainnetAddressPrivate    = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x22, 0xde}, nil} // PrivateKeyID, starts with Pm

	DCRD_testnetAddressP2PKH        = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0f, 0x21}, nil} //PubKeyHashAddrID,starts with Ts
	DCRD_testnetAddressP2PK         = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x28, 0xf7}, nil} //PubKeyAddrID, starts with Tk
	DCRD_testnetAddressPKHEdwards   = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0f, 0x01}, nil} //PKHEdwardsAddrID,starts with Te
	DCRD_testnetAddressP2PKHSchnorr = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0xe3}, nil} //PKHSchnorrAddrID,starts with TS
	DCRD_testnetAddressP2SH         = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0xfc}, nil} //ScriptHashAddrID,starts with Tc
	DCRD_testnetAddressPrivate      = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x23, 0x0e}, nil} //PrivateKeyID,starts with Pt

	DCRD_simnetAddressP2PKH      = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0x91}, nil} //PubKeyHashAddrID,starts with Ss
	DCRD_simnetAddressP2PK       = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x27, 0x6f}, nil} //PubKeyAddrID,starts with Sk
	DCRD_simnetAddressPKHEdwards = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0x71}, nil} //PKHEdwardsAddrID,starts with Se
	DCRD_simnetAddressPKHSchnorr = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0x53}, nil} //PKHSchnorrAddrID,starts with SS
	DCRD_simnetAddressP2SH       = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x0e, 0x6c}, nil} //ScriptHashAddrID,starts with Sc
	DCRD_simnetAddressPrivate    = AddressType{"base58", DCRDAlphabet, "doubleBlake256", "ripemd160", 20, []byte{0x23, 0x07}, nil} //PrivateKeyID, starts with Ps

	//TRON stuff
	TRON_mainnetAddress = AddressType{"base58", TRONAlphabet, "doubleSHA256", "keccak256_last_twenty", 20, []byte{0x41}, nil}
	TRON_testnetAddress = AddressType{"base58", TRONAlphabet, "doubleSHA256", "keccak256_last_twenty", 20, []byte{0xa0}, nil}

	//	OKChain_mainnetAddress = AddressType{"bech32", BTCBech32Alphabet, "okexchain", "keccak256_last_twenty", 20, nil, nil}
	//OKChain_mainnetAddress = AddressType{"bech32", BTCBech32Alphabet, "okexchain", "keccak256_last_twenty", 20, nil, nil}
)

type AddressType struct {
	EncodeType   string //编码类型
	Alphabet     string //码表
	ChecksumType string //checksum类型(Prefix string when encode type is base32PolyMod)
	HashType     string //地址hash类型，传入数据为公钥时起效
	HashLen      int    //编码前的数据长度
	Prefix       []byte //数据前面的填充
	Suffix       []byte //数据后面的填充
}
