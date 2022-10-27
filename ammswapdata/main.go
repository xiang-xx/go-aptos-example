package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/coming-chat/go-aptos/aptosclient"
)

const RpcUrl = "https://aptos.coming.chat"

type CoinInfo struct {
	Address  string
	Name     string
	Symbol   string
	Decimals int
}

type Connection struct {
	coins [2]CoinInfo
}

const PairPrefix = "0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::liquidity_pool::LiquidityPool"
const poolAddress = "0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"
const poolName = "lp"

var symbolToLogo map[string]string

func init() {
	symbolToLogo = make(map[string]string)
	symbolToLogo["BNB"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_BNB.png"
	symbolToLogo["BUSD"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_BUSD.png"
	symbolToLogo["ETH"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_ETH.png"
	symbolToLogo["USDC"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_USDC.png"
	symbolToLogo["USDT"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_usdt.png"
	symbolToLogo["BTC"] = "https://coming-website.s3.us-east-2.amazonaws.com/icon_xbtc_30.png"
	symbolToLogo["DAI"] = "https://coming-website.s3.us-east-2.amazonaws.com/DAI.png"
}

func main() {
	client, err := aptosclient.Dial(context.Background(), RpcUrl)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAccountResources(poolAddress)
	if err != nil {
		panic(err)
	}

	coinMap := make(map[string]CoinInfo)
	connections := make([]Connection, 0)

	for _, item := range res {
		if strings.HasPrefix(item.Type, PairPrefix) {
			coinInfos, err := getCoinsByPairAddress(item.Type, client)
			if err != nil {
				println(err)
			}

			for _, coinInfo := range coinInfos {
				coinMap[coinInfo.Address] = coinInfo
			}

			// create connection
			connections = append(connections, Connection{
				coinInfos,
			})
		}
	}

	network := "aptos_testnet"

	// 输出 connection create sql
	conSql := `insert into "public"."connection" ("chain_a", "token_a", "chain_b", "token_b", "route", "status", "ext")
	values`
	for _, conn := range connections {
		conSql += fmt.Sprintf(`
		('%s', '%s', '%s', '%s', 3, 1, '{"poolInfos":[{"address":"%s", "name":"%s"}]}'),`,
			network, conn.coins[0].Name, network, conn.coins[1].Name, poolAddress, poolName)
	}
	conSql = conSql[:len(conSql)-1] // 去掉最后一个 ,
	conSql += ";"
	ioutil.WriteFile("connections.sql", []byte(conSql), 0664)

	// 输出 token insert, chain_token insert
	// 	insert into "public"."token" ("chain_name", "asset_id", "address", "symbol", "decimals", "name", "logo_uri", "status", "price_usd", "val_order")
	// values
	// ('aptos_testnet', '-1', '0x1::aptos_coin::AptosCoin', 'APT', 8, 'APTOS_TEST', 'https://move-china-oss.oss-cn-hangzhou.aliyuncs.com/images/2022/09/13/e4316971835bcfcc5737504b2b52cafd.png', 1, 0, 0),
	tokenSql := `insert into "public"."token" ("chain_name", "asset_id", "address", "symbol", "decimals", "name", "logo_uri", "status", "price_usd", "val_order")
	values`

	// 	insert into "public"."chain_tokens" ("name", "symbol", "address", "chain_id", "chain_name", "logo_uri", "status", "decimals", "asset_id", "price_usd", "source", "value_order")
	// values
	//
	chainTokenSql := `insert into "public"."chain_tokens" ("name", "symbol", "address", "chain_id", "chain_name", "logo_uri", "status", "decimals", "asset_id", "price_usd", "source", "value_order")
	values`
	for _, token := range coinMap {
		if token.Address == "0x1::aptos_coin::AptosCoin" {
			continue
		}
		tokenLogo := GetTokenLogo(token.Symbol)
		tokenSql += fmt.Sprintf(`
			('%s', '-2', '%s', '%s', %d, '%s', '%s', 1, 0, 0),`,
			network, token.Address, token.Symbol, token.Decimals, token.Name, tokenLogo)

		chainTokenSql += fmt.Sprintf(`
		('%s', '%s', '%s', 0, '%s', '%s', 1, %d, -2, 0, 'ammswap', 0),`,
			token.Name,
			token.Symbol,
			token.Address,
			network,
			tokenLogo,
			token.Decimals)
	}
	tokenSql = tokenSql[:len(tokenSql)-1] // 去掉最后一个 ,
	tokenSql += ";"

	chainTokenSql = chainTokenSql[:len(chainTokenSql)-1] // 去掉最后一个 ,
	chainTokenSql += ";"

	ioutil.WriteFile("token.sql", []byte(tokenSql), 0664)

	ioutil.WriteFile("chain_token.sql", []byte(chainTokenSql), 0664)
}

func getCoinsByPairAddress(address string, client *aptosclient.RestClient) ([2]CoinInfo, error) {
	pool := strings.TrimPrefix(address, PairPrefix)
	pool = pool[1 : len(pool)-1] // 去掉 <>
	idx := strings.Index(pool, ", ")
	coin1 := pool[:idx]
	nextPoolStr := pool[idx+2:]

	idx = strings.Index(nextPoolStr, ", ")
	coin2 := ""
	if idx == -1 {
		coin2 = nextPoolStr
	} else {
		coin2 = nextPoolStr[:idx]
	}

	name, symbol, decimals := coinInfo(coin1, client)
	coin1Info := CoinInfo{
		Address:  coin1,
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
	name, symbol, decimals = coinInfo(coin2, client)
	coin2Info := CoinInfo{
		Address:  coin2,
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
	return [2]CoinInfo{coin1Info, coin2Info}, nil
}

func coinInfo(tokenAddress string, client *aptosclient.RestClient) (string, string, int) {
	firstIndex := strings.Index(tokenAddress, "::")

	address := tokenAddress[:firstIndex]

	res, err := client.GetAccountResource(address, fmt.Sprintf("0x1::coin::CoinInfo<%s>", tokenAddress), 0)
	if err != nil {
		panic(err)
	}

	name := res.Data["name"].(string)
	name = fixName(tokenAddress, name)

	return name, res.Data["symbol"].(string), int(res.Data["decimals"].(float64))
}

func fixName(tokenAddress string, name string) string {
	switch tokenAddress {
	case "0x870723e9a8f6d07c350e79d63655de673fb24d0695c702f479c201ab7b055f41::Coins::OmniUSDT":
		return "OmniUSDT"
	case "0xd415c5143d4f9752e462ab3476c567fdc0e2f0fb02f779d333e819c0e8624ea8::Coins::XBTC":
		return "XBTC(chainx)"
	case "0x870723e9a8f6d07c350e79d63655de673fb24d0695c702f479c201ab7b055f41::Coins::OmniXBTC":
		return "OmniXBTC"
	case "0xd415c5143d4f9752e462ab3476c567fdc0e2f0fb02f779d333e819c0e8624ea8::Coins::USDT":
		return "USDT(chainx)"
	case "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344::testnet_coins::TestBTC":
		return "testBTC"
	case "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344::testnet_coins::TestUSDC":
		return "testUSDC"
	}

	switch name {
	case "Tether":
		return "USDT"
	case "Aptos Coin":
		return "APTOS_TEST"
	}
	return name
}

func GetTokenLogo(symbol string) string {
	for t, v := range symbolToLogo {
		if strings.Contains(symbol, t) {
			return v
		}
	}
	return ""
}
