package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"go-aptos-example/base"
	"math/big"
	"strconv"
	"strings"

	"github.com/coming-chat/go-aptos-liquidswap/liquidswap"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	transactionbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/shopspring/decimal"
	"github.com/the729/lcs"
)

const (
	swapAbiFormat = "010473776170%s077363726970747300030866726f6d436f696e06746f436f696e076c70546f6b656e030b706f6f6c41646472657373040a66726f6d416d6f756e74020b746f416d6f756e744d696e02"

	APTOS = "0x1::aptos_coin::AptosCoin"
	USDT  = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT"
	BTC   = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::BTC"
	Pool  = ""

	scriptAddress = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9" // 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::liquidity_pool::LiquidityPool<CoinA, CoinB, Pool>
	poolAddress   = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9" // 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::lp::LP<CoinA, CoinB>
)
const upperhex = "0123456789ABCDEF"

var (
	address2Coin map[string]liquidswap.Coin
)

func init() {
	address2Coin = make(map[string]liquidswap.Coin)
	address2Coin[APTOS] = liquidswap.Coin{
		Decimals: 8,
		Symbol:   "APTOS",
	}
	address2Coin[USDT] = liquidswap.Coin{
		Decimals: 8,
		Symbol:   "USDT",
	}
	address2Coin[BTC] = liquidswap.Coin{
		Decimals: 8,
		Symbol:   "BTC",
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
	base.PanicError(err)

	chain := base.GetChain()

	// 构造交易，预估得到的 coin，执行 swap，查看交易详情
	swap(account, chain, APTOS, USDT, "100000")
}

func swap(account *aptos.Account, chain *aptos.Chain, fromCoinAddress, toCoinAddress, fromAmount string) {
	// 获取 resource
	client, err := chain.GetClient()
	base.PanicError(err)
	fromCoin := address2Coin[fromCoinAddress]
	toCoin := address2Coin[toCoinAddress]
	xAddress, yAddress := fromCoinAddress, toCoinAddress
	if !liquidswap.IsSortedSymbols(fromCoin.Symbol, toCoin.Symbol) {
		xAddress, yAddress = yAddress, xAddress
	}
	p := getPoolReserve(client, scriptAddress, poolAddress, xAddress, yAddress)
	amount, b := big.NewInt(0).SetString(fromAmount, 10)
	if !b {
		panic("invali params")
	}
	fmt.Printf("x: %s, %s\ny: %s %s\n", p.CoinXReserve, xAddress, p.CoinYReserve, yAddress)

	res := liquidswap.CalculateRates(fromCoin, toCoin, amount, "from", p)
	fmt.Printf("in %s: %s, out %s: %s\n", fromCoin.Symbol, amount.String(), toCoin.Name, res.String())

	payload, err := liquidswap.CreateTxPayload(&liquidswap.CreateTxPayloadParams{
		Script:           scriptAddress + "::scripts",
		FromCoin:         fromCoinAddress,
		ToCoin:           toCoinAddress,
		FromAmount:       amount,
		ToAmount:         res,
		InteractiveToken: "from",
		Slippage:         decimal.NewFromFloat(0.005),
		Pool: liquidswap.Pool{
			Address:       poolAddress,
			ModuleAddress: poolAddress,
			LpToken:       fmt.Sprintf("%s::lp::LP<%s,%s>", poolAddress, xAddress, yAddress),
		},
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
	var arg0 transactionbuilder.AccountAddress
	arg0Slice, err := hex.DecodeString(payload.Args[0][2:])
	base.PanicError(err)
	copy(arg0[:], arg0Slice)
	arg1, err := strconv.ParseUint(payload.Args[1], 10, 64)
	base.PanicError(err)
	arg2, err := strconv.ParseUint(payload.Args[2], 10, 64)
	base.PanicError(err)
	args := []interface{}{
		arg0,
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
	err = token.EnsureOwnerRegistedToken(account.Address(), account)
	base.PanicError(err)
}

func getPoolReserve(client *aptosclient.RestClient, scriptAddress, poolAddress, xAddress, yAddress string) liquidswap.PoolResource {
	poolResourceType := fmt.Sprintf(
		"%s::liquidity_pool::LiquidityPool<%s,%s,%s::lp::LP<%s,%s>>",
		scriptAddress,
		xAddress, // 这两个顺序与 lp 不一致
		yAddress,
		poolAddress,
		xAddress,
		yAddress,
	)
	poolResourceType = escapeTypes(poolResourceType)
	resource, err := client.GetAccountResource(scriptAddress, poolResourceType, 0)
	if err == nil {
		return resourceToPoolReserve(resource, false)
	} else {
		if !strings.HasPrefix(err.Error(), "Resource not found") {
			base.PanicError(err)
		}
	}

	poolResourceType = fmt.Sprintf(
		"%s::liquidity_pool::LiquidityPool<%s,%s,%s::lp::LP<%s,%s>>",
		scriptAddress,
		yAddress,
		xAddress, // 这两个顺序与 lp 不一致
		poolAddress,
		xAddress,
		yAddress,
	)
	poolResourceType = escapeTypes(poolResourceType)
	resource, err = client.GetAccountResource(scriptAddress, poolResourceType, 0)
	if err != nil {
		base.PanicError(err)
	}
	return resourceToPoolReserve(resource, true)
}

func resourceToPoolReserve(resource *aptostypes.AccountResource, reverse bool) liquidswap.PoolResource {
	x := resource.Data["coin_x_reserve"].(map[string]interface{})["value"].(string)
	y := resource.Data["coin_y_reserve"].(map[string]interface{})["value"].(string)
	xint, b := big.NewInt(0).SetString(x, 10)
	if !b {
		base.PanicError(errors.New("invalid reserve"))
	}
	yint, b := big.NewInt(0).SetString(y, 10)
	if !b {
		base.PanicError(errors.New("invalid reserve"))
	}
	s := liquidswap.PoolResource{
		CoinXReserve: xint,
		CoinYReserve: yint,
	}
	if reverse {
		s.CoinXReserve, s.CoinYReserve = s.CoinYReserve, s.CoinXReserve
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
