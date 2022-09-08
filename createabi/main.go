package main

import (
	"encoding/hex"
	"go-aptos-example/base"

	transactionbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/the729/lcs"
)

func main() {
	addressStr := "43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9"
	address, err := hex.DecodeString(addressStr)
	var a transactionbuilder.AccountAddress
	lcs.RegisterEnum(
		(*transactionbuilder.TypeTag)(nil),

		transactionbuilder.TypeTagBool{},
		transactionbuilder.TypeTagU8{},
		transactionbuilder.TypeTagU64{},
		transactionbuilder.TypeTagU128{},
		transactionbuilder.TypeTagAddress{},
		transactionbuilder.TypeTagSigner{},
		transactionbuilder.TypeTagVector{},
		transactionbuilder.TypeTagStruct{},
	)
	// lcs.RegisterEnum(
	// 	(*transactionbuilder.ScriptABI)(nil),

	// 	transactionbuilder.TransactionScriptABI{},
	// 	transactionbuilder.EntryFunctionABI{},
	// )
	copy(a[:], address)
	base.PanicError(err)
	abi := transactionbuilder.EntryFunctionABI{
		Name: "swap",
		ModuleName: transactionbuilder.ModuleId{
			Address: a,
			Name:    "scripts",
		},
		Doc: "",
		TyArgs: []transactionbuilder.TypeArgumentABI{
			{
				Name: "fromCoin",
			},
			{
				Name: "toCoin",
			},
			{
				Name: "lpToken",
			},
		},
		Args: []transactionbuilder.ArgumentABI{
			{
				Name:    "poolAddress",
				TypeTag: transactionbuilder.TypeTagAddress{},
			},
			{
				Name:    "fromAmount",
				TypeTag: transactionbuilder.TypeTagU64{},
			},
			{
				Name:    "toAmountMin",
				TypeTag: transactionbuilder.TypeTagU64{},
			},
		},
	}
	bs, err := lcs.Marshal(transactionbuilder.ScriptABI(abi))
	bs = append([]byte{1}, bs...)
	base.PanicError(err)
	println(hex.EncodeToString(bs))

	abiBytes := [][]byte{
		bs,
	}
	res, err := transactionbuilder.NewTransactionBuilderABI(abiBytes)
	base.PanicError(err)
	println(res)
}
