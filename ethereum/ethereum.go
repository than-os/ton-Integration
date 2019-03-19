package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {


	client, err := ethclient.Dial("https://rinkeby.infura.io")
	//client, err := ethclient.Dial("/home/user/.ethereum/geth.ipc")
	if err != nil {
		log.Println("error while connecting to eth backend")
		return
	}

	_ = client
	log.Println("connection established:")

	log.Println(GetAllTransactions(client))

	txHash := common.HexToHash("0xd882ba5070f976ee6fc6a6ff62c6c74e348493d1e035e60bb9d926ffdf80077e")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)

	if err != nil {
		println(err)
		return
	}

	if isPending {
		println("transaction status pending: ", isPending)
		return
	}

	d, _ := tx.MarshalJSON()
	log.Printf("tx details: \n%s",d )


}

func GetBalance(client *ethclient.Client) (string, error) {
	account := common.HexToAddress("0x0abdb71f2bdf4523b83770cd6dadae1b8e5e8f57")
	bal, err := client.BalanceAt(context.Background(), account, nil)

	log.Println("balance: ", bal.String())
	return bal.String(), err
}

func GetAllTransactions(client *ethclient.Client) (string) {

	//block, err := client.BlockByNumber(context.Background(), nil)
	//if err != nil {
	//	return err.Error()
	//}
	return "sdkjfn"
	//hexutil
	//return block.Transaction()
	//for _, b := range block.tr {
	//
	//}
}