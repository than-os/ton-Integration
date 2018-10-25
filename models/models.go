package models

import "github.com/globalsign/mgo/bson"

type TONNode struct {
	ID 				bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	WalletAddress 	string 			`json:"walletAddress" bson:"walletAddress"`
	Country 		string 			`json:"country" bson:"country"`
	IPAddr 			string 			`json:"ipAddr" bson:"ipAddr"`
	Port 			int 			`json:"port" bson:"port"`
	Username 		string 			`json:"userName" bson:"userName"`
	Password 		string 			`json:"password" bson:"password"`
}

type TXDetails struct {
	Result  struct {
		BlockHash 			string `json:"blockHash"`
		BlockNumber 		string `json:"blockNumber"`
		From 				string `json:"from"`
		Gas 				string `json:"gas"`
		GasPrice 			string `json:"gasPrice"`
		Hash 				string `json:"hash"`
		Input 				string `json:"input"`
		Nonce 				string `json:"nonce"`
		To 					string `json:"to"`
		TransactionIndex 	string `json:"transactionIndex"`
		Value 				string `json:"value"`
		V					string `json:"v"`
		R 					string `json:"r"`
		S 					string `json:"s"`
	}	`json:"result"`
}