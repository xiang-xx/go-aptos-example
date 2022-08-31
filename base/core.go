package base

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	txBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/the729/lcs"
)

const (
	testNetUrl = "https://fullnode.devnet.aptoslabs.com"
	faucetUrl  = "https://faucet.devnet.aptoslabs.com"

	AptosCoinType = "0x1::aptos_coin::AptosCoin"
)

var c *aptosclient.RestClient

func GetEnvAccount() (*aptosaccount.Account, error) {
	return aptosaccount.NewAccountWithMnemonic(os.Getenv("mnemonic"))
}

func GetAddress(account *aptosaccount.Account) string {
	return "0x" + hex.EncodeToString(account.AuthKey[:])
}

func RandomAccount() *aptosaccount.Account {
	seeds := make([]byte, 32)
	rand.Read(seeds)
	return aptosaccount.NewAccount(seeds)
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func FaucetFundAccount(address string, amount uint64) (err error) {
	_, err = aptosclient.FaucetFundAccount(address, amount, faucetUrl)
	return
}

func GetClient() *aptosclient.RestClient {
	if c != nil {
		return c
	}
	var err error
	c, err = aptosclient.Dial(context.Background(), testNetUrl)
	PanicError(err)
	return c
}

func WaitTxSuccess(txHash string) {
	for {
		client := GetClient()
		tx, err := client.GetTransactionByHash(txHash)
		if err != nil && strings.Contains(err.Error(), "not found") {
			time.Sleep(time.Second * 2)
			continue
		}
		PanicError(err)

		if tx.Type == aptostypes.TypePendingTransaction {
			time.Sleep(time.Second * 2)
			continue
		}

		if !tx.Success {
			PanicError(errors.New("tx failed " + txHash))
		}
		return
	}
}

func Transfer(account *aptosaccount.Account, toAddress string, amount uint64, coinType string) *aptostypes.Transaction {
	moduleName, err := txBuilder.NewModuleIdFromString("0x1::coin")
	PanicError(err)
	token, err := txBuilder.NewTypeTagStructFromString(coinType)
	PanicError(err)
	toAddr, err := txBuilder.NewAccountAddressFromHex(toAddress)
	PanicError(err)

	toAmountBytes, _ := lcs.Marshal(amount)
	payload := txBuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "transfer",
		TyArgs:       []txBuilder.TypeTag{*token},
		Args: [][]byte{
			toAddr[:], toAmountBytes,
		},
	}
	return SignAndSendEntryFunction(account, payload)
}

func SignAndSendEntryFunction(account *aptosaccount.Account, payload txBuilder.TransactionPayloadEntryFunction) *aptostypes.Transaction {
	client := GetClient()
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		panic(err)
	}
	accountData, err := client.GetAccount(GetAddress(account))
	PanicError(err)

	txn := &txBuilder.RawTransaction{
		Sender:                  account.AuthKey,
		SequenceNumber:          accountData.SequenceNumber,
		Payload:                 payload,
		MaxGasAmount:            3000,
		GasUnitPrice:            1,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}

	signedTxn, err := txBuilder.GenerateBCSTransaction(account, txn)
	PanicError(err)

	newTxn, err := client.SubmitSignedBCSTransaction(signedTxn)
	PanicError(err)

	return newTxn
}
