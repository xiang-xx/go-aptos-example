package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/coming-chat/go-aptos/aptosclient"
)

const host = "https://aptos-mainnet.pontem.network"

const typeFormat1 = "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::liquidity_pool::LiquidityPool<%s, %s, 0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::curves::Uncorrelated>"
const typeFormat2 = "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::liquidity_pool::LiquidityPool<%s, %s, 0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::curves::Stable>"

const poolAddress = "0x5a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"

//go:embed coins.json
var coindata []byte

type CoinInfo struct {
	TokenType TokenType `json:"token_type"`
}

type TokenType struct {
	Type string `json:"type"`
}

func main() {
	client, err := aptosclient.Dial(context.Background(), host)
	panicError(err)

	coins := make([]CoinInfo, 0)
	err = json.Unmarshal(coindata, &coins)
	panicError(err)

	rts := []string{}
	l := sync.RWMutex{}

	wg := sync.WaitGroup{}
	for i, aa := range coins {
		fmt.Printf("outer %d, total %d\n", i, len(coins))
		wg.Add(1)
		go func(a CoinInfo) {
			defer wg.Done()
			for j, b := range coins {
				fmt.Printf("- inner %d, total %d\n", j, len(coins))
				if a.TokenType.Type == b.TokenType.Type {
					continue
				}

				resource1 := fmt.Sprintf(typeFormat1, a.TokenType.Type, b.TokenType.Type)
				// resource2 := fmt.Sprintf(typeFormat2, a.TokenType.Type, b.TokenType.Type)

				_, e := client.GetAccountResource(poolAddress, resource1, 0)
				if e == nil {
					l.Lock()
					rts = append(rts, resource1)
					l.Unlock()
				}

				// _, e = client.GetAccountResource(poolAddress, resource2, 0)
				// if e == nil {
				// 	l.Lock()
				// 	rts = append(rts, resource2)
				// 	l.Unlock()
				// }
			}
		}(aa)
	}
	wg.Wait()

	sqls := make([]string, 0)
	// insert into "public"."token" ("chain_name", "asset_id", "address", "symbol", "decimals", "name", "logo_uri", "status", "price_usd", "val_order")
	// value;
	sqlFormat := `insert into "public"."aptos_resource_types" ("pool_name", "pool_address", "resource_type") values ('%s', '%s', '%s');`
	for _, r := range rts {
		sqls = append(sqls, fmt.Sprintf(sqlFormat, "pontem", poolAddress, r))
	}
	ioutil.WriteFile("resource_types.txt", []byte(strings.Join(rts, "\n")), 0644)
	ioutil.WriteFile("resource_types.sql", []byte(strings.Join(sqls, "\n")), 0644)
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
