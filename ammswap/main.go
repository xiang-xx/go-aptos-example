package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go-aptos-example/base"
	"math/big"
	"strconv"

	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	transactionbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/omnibtc/go-ammswap/ammswap"
	"github.com/shopspring/decimal"
)

const (
	swapAbiFormat = "010473776170%s09696e746572666163650002017801790208636f696e5f76616c0210636f696e5f6f75745f6d696e5f76616c02"
	// swapIntoAbi
	// 0109737761705f696e746f%s0773637269707473ab012053776170206d6178696d756d20636f696e2060586020666f7220657861637420636f696e206059602e0a202a2060636f696e5f76616c5f6d617860202d20686f77206d756368206f6620636f696e73206058602063616e206265207573656420746f206765742060596020636f696e2e0a202a2060636f696e5f6f757460202d20686f77206d756368206f6620636f696e73206059602073686f756c642062652072657475726e65642e0301780179056375727665020c636f696e5f76616c5f6d61780208636f696e5f6f757402

	APTOS = "0x1::aptos_coin::AptosCoin"
	USDC  = "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344::testnet_coins::TestUSDC"
	BTC   = "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344::testnet_coins::TestBTC"
	Pool  = ""

	scriptAddress     = "0xf69f9ec8348a803e2822c9f90950121130539f2a426dfb86e82d67e3613e6d6b" //
	scriptPoolAddress = "0xf69f9ec8348a803e2822c9f90950121130539f2a426dfb86e82d67e3613e6d6b"
	poolAddress       = "0xe98445b5e7489d1a4afee94940ca4c40e1f6c87a59c3b392e4744614af209de4" // 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::lp::LP<CoinA, CoinB>
)
const upperhex = "0123456789ABCDEF"

var (
	address2Coin map[string]ammswap.Coin
)

func init() {
	address2Coin = make(map[string]ammswap.Coin)
	address2Coin[APTOS] = ammswap.Coin{
		Decimals: 8,
		Name:     "APTOS",
		Symbol:   "APTOS",
		TokenType: ammswap.TokenType{
			Address:    "0x1",
			Module:     "aptos_coin",
			StructName: "AptosCoin",
		},
	}
	address2Coin[USDC] = ammswap.Coin{
		Decimals: 8,
		Symbol:   "testUSDC",
		Name:     "USD Coin",
		TokenType: ammswap.TokenType{
			Address:    "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344",
			Module:     "testnet_coins",
			StructName: "TestUSDC",
		},
	}
	address2Coin[BTC] = ammswap.Coin{
		Decimals: 8,
		Symbol:   "testBTC",
		Name:     "Bitcoin",
		TokenType: ammswap.TokenType{
			Address:    "0xcb0b45f3b49a6ab957facd2029ee0cd6720bb12907877d2f499946a7fd8f8344",
			Module:     "testnet_coins",
			StructName: "TestBTC",
		},
	}

	lcs.RegisterEnum(
		(*transactionbuilder.TypeTag)(nil),

		transactionbuilder.TypeTagBool{},
		transactionbuilder.TypeTagU8{},
		transactionbuilder.TypeTagU64{},
		transactionbuilder.TypeTagU128{},
		transactionbuilder.TypeTagAddress{},
		transactionbuilder.TypeTagSigner{},
		transactionbuilder.TypeTagVector{},
		transactionbuilder.TypeTagStruct{},
	)
}

func main() {
	account, err := base.GetEnvAptosAccount()
	println(account.Address())
	base.PanicError(err)

	chain := base.GetChain()

	// 构造交易，预估得到的 coin，执行 swap，查看交易详情
	swap(account, chain, APTOS, BTC, "1000000")
}

func swap(account *aptos.Account, chain *aptos.Chain, fromCoinAddress, toCoinAddress, fromAmount string) {
	// 获取 resource
	client, err := chain.GetClient()
	base.PanicError(err)
	fromCoin := address2Coin[fromCoinAddress]
	toCoin := address2Coin[toCoinAddress]
	xAddress, yAddress := fromCoinAddress, toCoinAddress
	if !ammswap.IsSortedCoin(fromCoin, toCoin) {
		xAddress, yAddress = yAddress, xAddress
	}
	p := getPoolReserve(client, scriptPoolAddress, poolAddress, xAddress, yAddress)
	amount, b := big.NewInt(0).SetString(fromAmount, 10)
	if !b {
		panic("invali params")
	}
	fmt.Printf("x: %s, %s\ny: %s %s\n", p.CoinXReserve, xAddress, p.CoinYReserve, yAddress)

	res := ammswap.GetAmountOut(fromCoin, toCoin, amount, p)
	fmt.Printf("in %s: %s, out %s: %s\n", fromCoin.Symbol, amount.String(), toCoin.Name, res.String())

	payload, err := ammswap.CreateSwapPayload(&ammswap.SwapParams{
		Script:     scriptAddress + "::interface",
		FromCoin:   fromCoinAddress,
		ToCoin:     toCoinAddress,
		FromAmount: amount,
		ToAmount:   res,
		Slippage:   decimal.NewFromFloat(0.005),
	})
	base.PanicError(err)

	fmt.Printf("%v", payload)

	// abi
	abiStr := fmt.Sprintf(swapAbiFormat, scriptAddress[2:])
	swapbytes, err := hex.DecodeString(abiStr)
	base.PanicError(err)
	abiBytes := [][]byte{
		swapbytes,
	}
	abi, err := transactionbuilder.NewTransactionBuilderABI(abiBytes)
	base.PanicError(err)

	// encode args
	arg1, err := strconv.ParseUint(payload.Args[0], 10, 64)
	base.PanicError(err)
	arg2, err := strconv.ParseUint(payload.Args[1], 10, 64)
	base.PanicError(err)
	args := []interface{}{
		arg1,
		arg2,
	}

	payloadBcs, err := abi.BuildTransactionPayload(payload.Function, payload.TypeArgs, args)
	base.PanicError(err)
	bcsBytes, err := lcs.Marshal(payloadBcs)
	base.PanicError(err)

	ensureRegisterCoin(account, chain, toCoinAddress)

	hash, err := chain.SubmitTransactionPayloadBCS(account, bcsBytes)
	base.PanicError(err)
	println(hash)
}

func ensureRegisterCoin(account *aptos.Account, chain *aptos.Chain, toCoinAddress string) {
	token, err := aptos.NewToken(chain, toCoinAddress)
	base.PanicError(err)
	_, err = token.EnsureOwnerRegistedToken(account)
	base.PanicError(err)
}

func getPoolReserve(client *aptosclient.RestClient, scriptPoolAddress, poolAddress, xAddress, yAddress string) ammswap.PoolResource {
	poolResourceType := fmt.Sprintf(
		"%s::implements::LiquidityPool<%s,%s>",
		scriptPoolAddress,
		xAddress, // 这两个顺序与 lp 不一致
		yAddress,
	)
	resource, err := client.GetAccountResource(poolAddress, poolResourceType, 0)
	if err != nil {
		base.PanicError(err)
	}

	return resourceToPoolReserve(resource)
}

func resourceToPoolReserve(resource *aptostypes.AccountResource) ammswap.PoolResource {
	x := resource.Data["coin_x"].(map[string]interface{})["value"].(string)
	y := resource.Data["coin_y"].(map[string]interface{})["value"].(string)
	xint, b := big.NewInt(0).SetString(x, 10)
	if !b {
		base.PanicError(errors.New("invalid reserve"))
	}
	yint, b := big.NewInt(0).SetString(y, 10)
	if !b {
		base.PanicError(errors.New("invalid reserve"))
	}
	s := ammswap.PoolResource{
		CoinXReserve: xint,
		CoinYReserve: yint,
	}
	return s
}

func escapeTypes(s string) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			if c == ' ' {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte
	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}
	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ':
			t[j] = '+'
			j++
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	return c == '<' || c == '>'
}
