package main

import (
	"fmt"

	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
)

func main() {
	// lcs.RegisterEnum(
	// 	(*txbuilder.TypeTag)(nil),

	// 	// TypeTagBool{},
	// 	// TypeTagU8{},
	// 	// TypeTagU64{},
	// 	// TypeTagU128{},
	// 	// TypeTagAddress{},
	// 	// TypeTagSigner{},
	// 	// TypeTagVector{},
	// 	txbuilder.TypeTagStruct{},
	// )

	moduleName, err := txbuilder.NewModuleIdFromString("0x1::managed_coin")
	if err != nil {
		panic(err)
	}
	typeTag, err := txbuilder.NewTypeTagStructFromString("0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT")
	if err != nil {
		panic(err)
	}
	payload := txbuilder.TransactionPayload(txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "register",
		TyArgs:       []txbuilder.TypeTag{*typeTag},
	})
	payloadLcsByte, err := lcs.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// data := "00000000000000000000000000000000000000000000000000000000000000010c6d616e616765645f636f696e087265676973746572010743417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b905636f696e7304555344540000"
	decodepayload := txbuilder.TransactionPayloadEntryFunction{}
	// payloadData, err := hex.DecodeString(data)
	// if err != nil {
	// 	panic(err)
	// }
	if err := lcs.Unmarshal(payloadLcsByte, &decodepayload); err != nil {
		panic(err)
	}
	fmt.Printf("%v", decodepayload)
}
