package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
			<1> for Node 1, <2> for Node 2 and so on...`)
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
	resp, err := db.CheckUserOptions(m.Sender.Username, "node")
	if err != nil {
		log.Println("error in handleTxHash: ", err)
		b.Send(m.Sender, "could not get user info")
		return
	}
	respToStr := fmt.Sprintf("%s", resp)
	strToInt, err := strconv.Atoi(respToStr)
	if err != nil {
		log.Println("error in ASCII to INT: ", err)
		b.Send(m.Sender, "ASCII to INT conversion error")
		return
	}

	idx := strToInt - 1
	//this wallet should be user wallet
	wallet, err := db.CheckUserOptions(m.Sender.Username, "wallet")
	if err != nil {
		b.Send(m.Sender, "error while getting user wallet. Please attach your wallet by sharing your ETH wallet with the bot")
		return
	}
	if findTxByHash(m.Text, fmt.Sprintf("%s", wallet)) {
		log.Println("line 1")
		//var n models.TONNode
		//log.Println("nodes here: ", nodes[strToInt-1])
		//for idx, node := range nodes {
		//	log.Println(node)
		//	if string(idx) == string(strToInt -1) {
		//		log.Println("hello hello")
		//		n = node
		//		log.Println("found node: ", node)
		//		break
		//	}
		//}

		//log.Println("node here: ", n)
		uri := "https://t.me/socks?server=" + nodes[idx].IPAddr + "&port=" + strconv.Itoa(nodes[idx].Port) + "&user=" + nodes[idx].Username + "&pass=" + nodes[idx].Password
		//log.Println("line 2")

		_, err := db.AddUserData(m.Sender.Username, "uri", uri)
		//log.Println("line 3")

		if err != nil {
		//	log.Println("line 4")
			log.Println("error in adding user data: ", err)
			b.Send(m.Sender, "error in adding user details")
			return
		}
		log.Println("line 5")

		b.Send(m.Sender, "Thanks for submitting the TX-HASH. We're validating it")
		userWallet, err := db.CheckUserOptions(m.Sender.Username, "wallet")
		if err != nil {
			b.Send(m.Sender, "error while fetching user wallet address. in case you have not attached your wallet address, please share your wallet address again.")
			return
		}
		log.Printf("your wallet address: %s", userWallet)

		b.Send(m.Sender, "Congratulations!! please click the button below to connect to the sentinel dVPN node", &telebot.ReplyMarkup{
			InlineKeyboard: inlineButton(nodes[idx].Username, uri),
		})
		return
	}

	b.Send(m.Sender, "something wrong with tx hash")

}

func handleNodeId(b *telebot.Bot, m *telebot.Message, nodes []models.TONNode) {
	NODE_ID = m.Text
	idx, _ := strconv.Atoi(m.Text)
	if idx > len(nodes) {
		b.Send(m.Sender, "invalid node id")
		return
	}
	_, err := db.AddUserData(m.Sender.Username, "node", m.Text)
	if err != nil {
		b.Send(m.Sender, "could not store user info")
		return
	}

	log.Println("in nodeId case", m.Text)
	b.Send(m.Sender, "please send 2 SENTS to the following address and submit the tx hash here: ")
	b.Send(m.Sender, "0xCeb5bC384012f0EEbeE119d82A24925C47714fe3") //should be node wallet address
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
	var body models.TxReceiptList
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Println(err)
		return false
	}
	for _, txReceipt := range body.Results {
		if txReceipt.TransactionHash == txHash {
			log.Println("response decoded: ", body.Results[0])
			//topic[1] from
			//topic[2] to
			okWallet := txReceipt.Topics[1] == walletAddr
			okRecipient := txReceipt.Topics[2] == os.Getenv("nodeWallet")
			amount := hex2int(txReceipt.Data)
			if okWallet && okRecipient && amount == uint64(2) {
				return true
			}
			log.Println("\n\n money: ", hex2int(txReceipt.Data))

		}
	}

	return false
}

func handleWalletAddress(b *telebot.Bot, m *telebot.Message) {
	replyButtons := [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{
				Text: "100 MB",
			},
		},
		{
			telebot.ReplyButton{
				Text: "500 MB",
			},
		},
		{
			telebot.ReplyButton{
				Text: "1 GB",
			},
		},
	}
	ok := common.IsHexAddress(m.Text)
	if ok {
		_ , err := db.AddUserData(m.Sender.Username, "wallet" , m.Text)
		if err != nil {
			b.Send(m.Sender, "error while storing user eth address")
			return
		}
		b.Send(m.Sender, "Attached the ETH wallet to user successfully")
		b.Send(m.Sender, `Please select how much bandwidth you need by clicking on one of the buttons below: `, &telebot.ReplyMarkup{
			ReplyKeyboard:       replyButtons,
			ResizeReplyKeyboard: true,
			OneTimeKeyboard: true,
			ReplyKeyboardRemove: true,
		})
		return
	}
}

func hex2int(hexStr string) uint64 {
// remove 0x suffix if found in the input string
cleaned := strings.Replace(hexStr, "0x", "", -1)

// base 16 for hexadecimal
result, _ := strconv.ParseUint(cleaned, 16, 64)
return uint64(result)
}