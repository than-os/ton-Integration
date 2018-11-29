package models

import "github.com/globalsign/mgo/bson"

type TONNode struct {
	ID 				bson.ObjectId 	`json:"id" bson:"_id,omitempty"`
	WalletAddress 	string 			`json:"walletAddr" bson:"walletAddress"`
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

type TxReceiptList struct {
	Results []TxReceipt `json:"result"`
	Status string `json:"status"`
}

type TxReceipt struct {
	Address string `json:"address"`
	Topics []string `json:"topics"`
	Data string `json:"data"`
	BlockNumber string `json:"blockNumber"`
	Timestamp string `json:"timestamp"`
	GasPrice string `json:"gasPrice"`
	GasUser string `json:"gasUsed"`
	LogIndex string `json:"logIndex"`
	TransactionHash string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
}

type KeepAliveRequest struct {
	Status 		string `json:"status"`
	NodeIpAddr string `json:"nodeIPAddr"`
}

type KeepAliveResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type ConfigData struct {
	Token 			string `json:"token"`
	EncMethod 		string `json:"enc_method"`
	AccountAddr 	string `json:"account_addr"`
	PricePerGB 		float64 `json:"price_per_gb"`
	TONPrice 		float64 `json:"tonPrice"`
}

type AddUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RemoveUser struct {
	Username string `json:"username"`
}

type UserRes struct {
	Message string `json:"messages"`
}

type ExpiredUsers struct {
	Key string
	Value string
}

type RegisterRequest struct {
	IPAddr 			string `json:"ipAddr"`
	Port			string `json:"port"`
	Username		string `json:"userName"`
	Password 		string `json:"password"`
	WalletAddr 		string  `json:"walletAddr"`
}