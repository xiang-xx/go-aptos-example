package main

import (
	_ "embed"
	"encoding/hex"
	"fmt"
)

//go:embed test/close.abi
var testclose []byte

//go:embed test/open.abi
var testopen []byte

//go:embed test/create.abi
var testcreate []byte

//go:embed online/close.abi
var onlineclose []byte

//go:embed online/open.abi
var onlineopen []byte

//go:embed online/create.abi
var onlinecreate []byte

func main() {
	fmt.Printf("%s\n", hex.EncodeToString(testcreate))
	fmt.Printf("%s\n", hex.EncodeToString(onlinecreate))

	println("===")
	fmt.Printf("%s\n", hex.EncodeToString(testopen))
	fmt.Printf("%s\n", hex.EncodeToString(onlineopen))

	println("===")
	fmt.Printf("%s\n", hex.EncodeToString(testclose))
	fmt.Printf("%s\n", hex.EncodeToString(onlineclose))
}
