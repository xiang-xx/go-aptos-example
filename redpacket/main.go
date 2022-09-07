package main

import (
	"encoding/hex"
	"fmt"
	"go-aptos-example/base"
	"os"
	"time"

	"github.com/coming-chat/go-aptos/aptosaccount"
	transactionbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/go-red-packet/redpacket"
	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/the729/lcs"
)

const (
	createABI = "0106637265617465b6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb0a7265645f7061636b657400000205636f756e74020d746f74616c5f62616c616e636502"
	openABI   = "01046f70656eb6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb0a7265645f7061636b6574000003026964020e6c75636b795f6163636f756e747306040862616c616e6365730602"
	closeABI  = "0105636c6f7365b6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb0a7265645f7061636b657400000102696402"
)

func main() {
	// account, err := aptos.NewAccountWithMnemonic(os.Getenv("mnemonic"))
	// base.PanicError(err)

	chain := aptos.NewChainWithRestUrl(base.TestNetUrl)
	contractAddress := os.Getenv("redpacket")
	// contract := redpacket.NewAptosRedPacketContract(chain, contractAddress)

	// abi
	abiBytes := make([][]byte, 0)
	abiStrs := []string{createABI, openABI, closeABI}
	for _, s := range abiStrs {
		bs, err := hex.DecodeString(s)
		base.PanicError(err)
		abiBytes = append(abiBytes, bs)
	}
	redpacketAbi, err := transactionbuilder.NewTransactionBuilderABI(abiBytes)
	base.PanicError(err)

	account, err := aptosaccount.NewAccountWithMnemonic(os.Getenv("mnemonic"))
	base.PanicError(err)
	create(account, chain, redpacketAbi, contractAddress)

	// rid := 1
	// open(account, chain, contract, int64(rid))
	// close(account, chain, contract, int64(rid))
}

func create(account *aptosaccount.Account, chain *aptos.Chain, abi *transactionbuilder.TransactionBuilderABI, contractAddress string) {
	functionName := contractAddress + "::red_packet::create"
	// 0xb6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb::red_packet::create
	payloadAbi, err := abi.BuildTransactionPayload(
		functionName,
		[]string{},
		[]any{
			uint64(5), uint64(100000),
		},
	)

	bs, err := lcs.Marshal(payloadAbi)
	base.PanicError(err)
	fmt.Printf("%s\b", hex.EncodeToString(bs))
	pabi := transactionbuilder.TransactionPayloadEntryFunction{}
	lcs.Unmarshal(bs, &pabi)
	fmt.Printf("%v\n", pabi)

	base.PanicError(err)
	client, err := chain.GetClient()
	base.PanicError(err)
	accountData, err := client.GetAccount("0x" + hex.EncodeToString(account.AuthKey[:]))
	base.PanicError(err)
	ledgerInfo, err := client.LedgerInfo()
	base.PanicError(err)
	txAbi := &transactionbuilder.RawTransaction{
		Sender:                  account.AuthKey,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            20000,
		GasUnitPrice:            1,
		Payload:                 pabi,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}

	signedTxn, err := transactionbuilder.GenerateBCSTransaction(account, txAbi)
	base.PanicError(err)

	tx, err := client.SubmitSignedBCSTransaction(signedTxn)
	base.PanicError(err)
	txHash := tx.Hash
	time.Sleep(time.Second * 2)
	txDetail, err := chain.FetchTransactionDetail(txHash)
	if err != nil {
		panic(err)
	}

	println(txHash)
	println(txDetail.Status)
}

func open(account *aptos.Account, chain *aptos.Chain, contract redpacket.RedPacketContract, id int64) {
	action, err := redpacket.NewRedPacketActionOpen(id, []string{
		account.Address(),
		account.Address(),
	}, []string{
		"5000",
		"5000",
	})
	base.PanicError(err)
	// 使用合约对象发送 action 交易到链上
	txHash, err := contract.SendTransaction(account, action)
	if err != nil {
		panic(err)
	}
	txDetail, err := chain.FetchTransactionDetail(txHash)
	if err != nil {
		panic(err)
	}

	println(txHash)
	println(txDetail.Status)
}

func close(account *aptos.Account, chain *aptos.Chain, contract redpacket.RedPacketContract, id int64) {
	action, err := redpacket.NewRedPacketActionClose(id, account.Address())
	base.PanicError(err)
	// 使用合约对象发送 action 交易到链上
	txHash, err := contract.SendTransaction(account, action)
	if err != nil {
		panic(err)
	}
	txDetail, err := chain.FetchTransactionDetail(txHash)
	if err != nil {
		panic(err)
	}

	println(txHash)
	println(txDetail.Status)
}

func getAuthKey(account *aptos.Account) transactionbuilder.AccountAddress {
	key, err := hex.DecodeString(account.Address()[2:])
	base.PanicError(err)
	var a transactionbuilder.AccountAddress
	copy(a[:], key)
	return a
}
