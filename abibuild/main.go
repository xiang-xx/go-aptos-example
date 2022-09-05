package main

import (
	_ "embed"
	"encoding/hex"
	"fmt"
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
}
