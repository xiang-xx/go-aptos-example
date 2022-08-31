package main

import (
	"fmt"
	"go-aptos-example/base"
)

func main() {
	alice := base.RandomAccount()
	err := base.FaucetFundAccount(base.GetAddress(alice), 50000)
	base.PanicError(err)

	bob := base.RandomAccount()
	err = base.FaucetFundAccount(base.GetAddress(bob), 0)
	base.PanicError(err)

	client := base.GetClient()
	aliceCoin, err := client.GetAccountResource(base.GetAddress(alice), "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", 0)
	base.PanicError(err)
	fmt.Printf("alice: %s\n", aliceCoin.Data["coin"].(map[string]interface{})["value"])

	bobCoin, err := client.GetAccountResource(base.GetAddress(bob), "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", 0)
	base.PanicError(err)
	fmt.Printf("bob: %s\n", bobCoin.Data["coin"].(map[string]interface{})["value"])

	tx := base.Transfer(alice, base.GetAddress(bob), 1000, base.AptosCoinType)
	fmt.Printf("txhash: %s\n", tx.Hash)
	base.WaitTxSuccess(tx.Hash)

	aliceCoin, err = client.GetAccountResource(base.GetAddress(alice), "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", 0)
	base.PanicError(err)
	fmt.Printf("alice: %s\n", aliceCoin.Data["coin"].(map[string]interface{})["value"])

	bobCoin, err = client.GetAccountResource(base.GetAddress(bob), "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", 0)
	base.PanicError(err)
	fmt.Printf("bob: %s\n", bobCoin.Data["coin"].(map[string]interface{})["value"])
}
