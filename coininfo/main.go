package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/coming-chat/go-aptos/aptosclient"
)

func main() {
	client, err := aptosclient.Dial(context.Background(), "https://fullnode.testnet.aptoslabs.com")
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 1 {
		usagePanic()
	}

	tokenAddress := os.Args[1]

	firstIndex := strings.Index(tokenAddress, "::")
	if firstIndex < 0 {
		usagePanic()
	}

	address := tokenAddress[:firstIndex]

	res, err := client.GetAccountResource(address, fmt.Sprintf("0x1::coin::CoinInfo<%s>", tokenAddress), 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("decimals: %d\n", int(res.Data["decimals"].(float64)))
	fmt.Printf("symbol: %s\n", res.Data["symbol"].(string))
	fmt.Printf("name: %s\n", res.Data["name"].(string))
}

func usagePanic() {
	panic("Usage: coininfo 0x123123::coin::Coin")
}
