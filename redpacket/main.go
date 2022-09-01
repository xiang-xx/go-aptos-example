package main

import (
	"encoding/json"
	"go-aptos-example/base"
	"os"

	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/aptos"
)

func main() {
	account, err := aptos.NewAccountWithMnemonic(os.Getenv("mnemonic"))
	base.PanicError(err)

	chain := aptos.NewChainWithRestUrl(base.TestNetUrl)
	payload := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      "0x9fecf7c6ad1fc0e5337c6d64443cda47b41f61b556a698193646ce0b8917cbe1::red_packet::open",
		TypeArguments: []string{},
		Arguments: []interface{}{
			"1",
			[]string{"0x9fecf7c6ad1fc0e5337c6d64443cda47b41f61b556a698193646ce0b8917cbe1", "0x9fecf7c6ad1fc0e5337c6d64443cda47b41f61b556a698193646ce0b8917cbe1"},
			[]string{"1000", "8750"},
		},
	}
	data, err := json.Marshal(payload)
	base.PanicError(err)
	hash, err := chain.SubmitTransactionPayload(account, data)
	base.PanicError(err)
	base.WaitTxSuccess(hash)
}
