package hdwallet

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// default options
var (
	DefaultParams = &chaincfg.MainNetParams

	// master key options
	DefaultPassword = ""
	DefaultLanguage = ""
	DefaultWIF      = ""
	DefaultMnemonic = ""
	DefaultNet      = "mainnet"

	DefaultSeed = []byte{}

	// child key options
	DefaultPurpose      = ZeroQuote + 44
	DefaultCoinType     = BTC
	DefaultAccount      = ZeroQuote
	DefaultChange       = Zero
	DefaultAddressIndex = Zero
)

// Option of key
type Option func(*Options)

// Options of key
type Options struct {
	Params *chaincfg.Params

	// master key options
	Mnemonic string
	Password string
	Language string
	Seed     []byte
	WIF      string
	Net      string
	// child key options
	Purpose      uint32
	CoinType     uint32
	Account      uint32
	Change       uint32
	AddressIndex uint32
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		Params:       DefaultParams,
		Password:     DefaultPassword,
		Language:     DefaultLanguage,
		Purpose:      DefaultPurpose,
		CoinType:     DefaultCoinType,
		Account:      DefaultAccount,
		Change:       DefaultChange,
		AddressIndex: DefaultAddressIndex,
		Mnemonic:     DefaultMnemonic,
		Net:          DefaultNet,
		WIF:          DefaultWIF,
		Seed:         DefaultSeed,
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

// GetPath return path in bip44 style
func (o *Options) GetPath() []uint32 {
	return []uint32{
		o.Purpose,
		o.CoinType,
		o.Account,
		o.Change,
		o.AddressIndex,
	}
}
func (o *Options) GetQTUMPath() []uint32 {
	return []uint32{
		o.CoinType,
		DefaultAccount,
		DefaultAccount,
	}
}

func (o *Options) GetBTMPath() []uint32 {
	return []uint32{
		o.CoinType,
		DefaultChange,
		DefaultChange,
		DefaultAddressIndex,
	}
}

// Params set to options
func Params(p *chaincfg.Params) Option {
	return func(o *Options) {
		o.Params = p
	}
}

// net set to options
func Net(m string) Option {
	return func(o *Options) {
		o.Net = m
	}
}

// Mnemonic set to options
func Mnemonic(m string) Option {
	return func(o *Options) {
		o.Mnemonic = m
	}
}

// WIF set to options
func WIF(wif string) Option {
	return func(o *Options) {
		o.WIF = wif
	}
}

// Password set to options
func Password(p string) Option {
	return func(o *Options) {
		o.Password = p
	}
}

// Language set to options
func Language(l string) Option {
	return func(o *Options) {
		o.Language = l
	}
}

// Seed set to options
func Seed(s []byte) Option {
	return func(o *Options) {
		o.Seed = s
	}
}

// Purpose set to options
func Purpose(p uint32) Option {
	return func(o *Options) {
		o.Purpose = p
	}
}

// CoinType set to options
func CoinType(c uint32) Option {
	return func(o *Options) {
		o.CoinType = c
	}
}

// Account set to options
func Account(a uint32) Option {
	return func(o *Options) {
		o.Account = a
	}
}

// Change set to options
func Change(c uint32) Option {
	return func(o *Options) {
		o.Change = c
	}
}

// AddressIndex set to options
func AddressIndex(a uint32) Option {
	return func(o *Options) {
		o.AddressIndex = a
	}
}
