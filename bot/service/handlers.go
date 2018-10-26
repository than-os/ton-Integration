package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/than-os/sent-dante/constants"
	"github.com/than-os/sent-dante/models"
	"github.com/than-os/sent-dante/utils"
	"gopkg.in/tucnak/telebot.v2"
)

func handleBandwidth(b *telebot.Bot, m *telebot.Message, nodes []models.TONNode) {
	resp, err := db.CheckUserOptions(m.Sender.Username, "bw")
	if err != nil {
		log.Println("error in checkUserOptions", err.Error())

		_, err := db.AddUserData(m.Sender.Username, "bw", m.Text)
		if err != nil {
			log.Println("error in adduser")
		}
		switch m.Text {
		case "100 MB":
			b.Send(m.Sender, "you have opted for 100MB bandwidth")
		case "500 MB":
			b.Send(m.Sender, "you have opted for 500MB bandwidth")
		case "1 GB":
			b.Send(m.Sender, "you have opted for 1GB bandwidth")
		}
		b.Send(m.Sender, `Please select a node ID from the list below and reply in the format of
			<#1> for Node 1, <#2> for Node 2`)
		for idx, node := range nodes {
			//uri := "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
			geo, err := utils.GetGeoLocation(node.IPAddr)

			if err != nil {
				b.Send(m.Sender, "error: ", err.Error())
				return
			}
			b.Send(m.Sender, strconv.Itoa(idx+1)+". Location: "+geo.Country+"\n "+"User:"+node.Username+"\n "+"Node Wallet: "+"0xCeb5bC384012f0EEbeE119d82A24925C47714fe3")
		}
		log.Println("here insdie")
		return
	}
	log.Println("here")
	nodeIdx, err := strconv.ParseInt(fmt.Sprintf("%s", resp), 10, 64)
	if err != nil {
		log.Println(err)
	}

	log.Println("below: ", nodes)
	var n models.TONNode
	for i := 0; i < len(nodes); i++ {
		if i == int(nodeIdx) {
			n = nodes[i]
			return
		}
	}
	uri := "https://t.me/socks?server=" + n.IPAddr + "&port=" + strconv.Itoa(n.Port) + "&user=" + n.Username + "&pass=" + n.Password
	log.Println("url: => ", uri, "\n", n)
	b.Send(m.Sender, "you have already selected : Node "+fmt.Sprintf("%s", resp), &telebot.ReplyMarkup{
		InlineKeyboard: inlineButton("sentinel node", uri),
	})
	log.Println("much ", nodeIdx)
}

func handleTxHash(b *telebot.Bot, m *telebot.Message, nodes []models.TONNode) {

	val, err := strconv.Atoi(NODE_ID[1:])
	if err != nil {
		b.Send(m.Sender, "Could not read NODE ID")
		return
	}
	log.Println("in validhex case", m.Text, val-1)
	nodeIdx := val - 1
	if nodeIdx > len(nodes) {
		b.Send(m.Sender, "invalid node id")
		return
	}
	if findTxByHash(m.Text, "0x6b6df9e25f7bf2e363ec1a52a7da4c4a64f5769e") {
		uri := "https://t.me/socks?server=" + nodes[nodeIdx].IPAddr + "&port=" + strconv.Itoa(nodes[nodeIdx].Port) + "&user=" + nodes[nodeIdx].Username + "&pass=" + nodes[nodeIdx].Password

		b.Send(m.Sender, "Thanks for submitting the TX-HASH. We're validating it")
		b.Send(m.Sender, "Congratulations!! please click the button below to connect to the sentinel dVPN node", &telebot.ReplyMarkup{
			InlineKeyboard: inlineButton(nodes[nodeIdx].Username, uri),
		})
	}
}

func handleNodeId(b *telebot.Bot, m *telebot.Message, nodes []models.TONNode) {
	NODE_ID = m.Text
	idx, _ := strconv.Atoi(m.Text)
	if idx > len(nodes) {
		b.Send(m.Sender, "invalid node id")
		return
	}
	log.Println("in nodeId case", m.Text)
	b.Send(m.Sender, "please send 5 SENTS to the following address and submit the tx hash here: ")
	b.Send(m.Sender, "0xCeb5bC384012f0EEbeE119d82A24925C47714fe3")
}

func validHex(b *telebot.Bot, m *telebot.Message) string {
	_, err := hexutil.Decode(m.Text)
	if err != nil {
		b.Send(m.Sender, "invalid tx hash")
		return ""
	}
	return m.Text
}

func validWalletAddr(m *telebot.Message) string {
	log.Println("isWallet: ")
	isWallet := common.IsHexAddress(m.Text)
	log.Println("isWallet: ", isWallet)
	if isWallet {
		return m.Text
	}
	return ""
}

func nodeToInt(m *telebot.Message) string {

	_, err := strconv.Atoi(m.Text)
	if err != nil {
		return ""
	}

	return m.Text
}

func nodeId(m *telebot.Message) string {
	rgx := regexp.MustCompile(`([0-9])\w*`)
	if rgx.MatchString(m.Text) {
		return m.Text
	}
	//b.Send(m.Sender, "Invalid Node Id. Please use correct format as <n1> for Node 1 or Node Doesn't exists")
	return ""
}

func findTxByHash(txHash, walletAddr string) bool {
	//w, err := hexutil.Decode(walletAddr)
	//if err != nil {
	//      log.Println(err)
	//      return false
	//}
	wallet := "0x" + constants.ZFill + strings.TrimLeft(walletAddr, "0x")
	uri := constants.TEST_SENT_URI + wallet + constants.TEST_SENT_URI2 + wallet
	log.Println("length of the wallet := ", len(wallet))
	resp, err := http.Get(uri)
	if err != nil {
		log.Println(err)
		return false
	}
	var body interface{}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Println(err)
		return false
	}
	log.Println("response decoded: ", body)

	return true
}
