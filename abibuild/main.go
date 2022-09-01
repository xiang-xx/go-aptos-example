package main

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"go-aptos-example/base"

	builder "github.com/coming-chat/go-aptos/transaction_builder"
)

//go:embed close.abi
var close []byte

//go:embed open.abi
var open []byte

//go:embed create.abi
var create []byte

func main() {
	fmt.Printf("%s\n", hex.EncodeToString(create))
	fmt.Printf("%s\n", hex.EncodeToString(open))
	fmt.Printf("%s\n", hex.EncodeToString(close))
	account, err := base.GetEnvAccount()
	base.PanicError(err)

	abi, err := builder.NewTransactionBuilderABI([][]byte{create, open, close}, &builder.ABIBuilderConfig{
		Sender: account.AuthKey,
	})
	if err != nil {
		panic(err)
	}
	println(abi)
}
