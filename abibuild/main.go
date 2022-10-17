package main

import (
	_ "embed"
	"encoding/hex"
	"fmt"
)

//go:embed amm/swap.abi
var swap []byte

// //go:embed lp/swap_into.abi
// var swap_into []byte

// //go:embed test/close.abi
// var testclose []byte

// //go:embed test/open.abi
// var testopen []byte

// //go:embed test/create.abi
// var testcreate []byte

// //go:embed online/close.abi
// var onlineclose []byte

// //go:embed online/open.abi
// var onlineopen []byte

// //go:embed online/create.abi
// var onlinecreate []byte

// //go:embed mc/deposit.abi
// var deposit []byte

// //go:embed mc/initialize.abi
// var initialize []byte

// //go:embed mc/issue.abi
// var issue []byte

// //go:embed mc/register_withdraw.abi
// var register_withdraw []byte

// //go:embed mc/register.abi
// var register []byte

// //go:embed mc/set_pause.abi
// var set_pause []byte

// //go:embed mc/transfer.abi
// var transfer []byte

// //go:embed mc/withdraw.abi
// var withdraw []byte

func main() {
	// fmt.Printf("deposit %s\n", hex.EncodeToString(deposit))
	// fmt.Printf("initialize %s\n", hex.EncodeToString(initialize))
	// fmt.Printf("issue %s\n", hex.EncodeToString(issue))
	// fmt.Printf("register_withdraw %s\n", hex.EncodeToString(register_withdraw))
	// fmt.Printf("register %s\n", hex.EncodeToString(register))
	// fmt.Printf("set_pause %s\n", hex.EncodeToString(set_pause))
	// fmt.Printf("transfer %s\n", hex.EncodeToString(transfer))
	// fmt.Printf("withdraw %s\n", hex.EncodeToString(withdraw))

	// fmt.Printf("%s\n", hex.EncodeToString(onlinecreate))

	// println("===")
	// fmt.Printf("%s\n", hex.EncodeToString(onlineopen))

	// println("===")
	// fmt.Printf("%s\n", hex.EncodeToString(onlineclose))

	// fmt.Printf("%s\n", hex.EncodeToString(testcreate))
	// println("===")

	// fmt.Printf("%s\n", hex.EncodeToString(testopen))
	// println("===")

	// fmt.Printf("%s\n", hex.EncodeToString(testclose))

	fmt.Printf("%s\n", hex.EncodeToString(swap))
	// fmt.Printf("%s\n", hex.EncodeToString(swap_into))

}

// {
// 	"arguments": [
// 	  "1000000",
// 	  "20977"
// 	],
// 	"function": "0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::scripts::swap",
// 	"type": "entry_function_payload",
// 	"type_arguments": [
// 	  "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT",
// 	  "0xb4d7b2466d211c1f4629e8340bb1a9e75e7f8fb38cc145c54c5c9f9d5017a318::coins_extended::USDC",
// 	  "0x4e9fce03284c0ce0b86c88dd5a46f050cad2f4f33c4cdd29d98f501868558c81::curves::Uncorrelated"
// 	]
//   }
