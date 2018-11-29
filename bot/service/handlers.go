package service

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/than-os/sent-dante/constants"
	"github.com/than-os/sent-dante/models"
	"github.com/than-os/sent-dante/utils"
	"gopkg.in/tucnak/telebot.v2"
)

var (
	Duration = 10
	Amount = 2
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
		case "10 Days":
			b.Send(m.Sender, "you have opted for 10 days of unlimited bandwidth")
			t := time.Hour * 24 * 10
			db.AddUserData(m.Sender.Username, "timestamp-", time.Now().Add(t).Format(time.RFC3339))
		case "30 Days":
			b.Send(m.Sender, "you have opted for 30 days of unlimited bandwidth")
			Duration = 30
			Amount = 4
			t := time.Hour * 24 * 30
			db.AddUserData(m.Sender.Username, "timestamp-", time.Now().Add(t).Format(time.RFC3339))
		case "90 Days":
			b.Send(m.Sender, "you have opted for 90 days of unlimited bandwidth")
			Duration = 90
			t := time.Hour * 24 * 90
			db.AddUserData(m.Sender.Username, "timestamp-", time.Now().Add(t).Format(time.RFC3339))
			Amount = 6
		}
		b.Send(m.Sender, `Please select a node ID from the list below and reply in the format of
			1 for Node 1, 2 for Node 2 and so on...`)
		for idx, node := range nodes {
			//uri := "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
			geo, err := utils.GetGeoLocation(node.IPAddr)

			if err != nil {
				b.Send(m.Sender, "error: ", err.Error())
				return
			}
			b.Send(m.Sender, strconv.Itoa(idx+1)+". Location: "+geo.Country + "Node Wallet: "+"0xCeb5bC384012f0EEbeE119d82A24925C47714fe3")
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
	if findTxByHash(m.Text,  , m) {
		uri := "https://t.me/socks?server=" + nodes[idx].IPAddr + "&port=" + strconv.Itoa(nodes[idx].Port) + "&user=" + nodes[idx].Username + "&pass=" + nodes[idx].Password
		//log.Println("line 2")

		_, err := db.AddUserData(m.Sender.Username, "ipaddr", nodes[idx].IPAddr)
		if err != nil {
			log.Println("error in adding user data: ", err)
			b.Send(m.Sender, "error in adding user details")
			return
		}
		_, err = db.AddUserData(m.Sender.Username, "uri", uri)
		//log.Println("line 3")

		if err != nil {
		//	log.Println("line 4")
			log.Println("error in adding user data: ", err)
			b.Send(m.Sender, "error in adding user details")
			return
		}
		//log.Println("line 5")
		_ ,err = db.AddUserData(m.Sender.Username, "auth", "true")
		if err != nil {
			b.Send(m.Sender, "error while adding user to auth group. please try again")
			return
		}
		b.Send(m.Sender, "Thanks for submitting the TX-HASH. We're validating it")
		b.Send(m.Sender, "creating new user for " + m.Sender.Username + "...")
		node := nodes[idx]
		err = utils.AddUser(node.IPAddr, m.Sender.Username)
		if err != nil {
			b.Send(m.Sender, "Error while creating SOCKS5 user for "+ m.Sender.Username)
			return
		}
		pass, err := db.CheckUserOptions(m.Sender.Username, "password")
		if err != nil {
			b.Send(m.Sender, "error while getting user pass")
			return
		}
		uri = "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + m.Sender.Username + "&pass=" + fmt.Sprintf("%s", pass)

		_, err = db.AddUserData(m.Sender.Username, "ipaddr", nodes[idx].IPAddr)
		if err != nil {
			log.Println("error in adding user data: ", err)
			b.Send(m.Sender, "error in adding user details")
			return
		}
		_, err = db.AddUserData(m.Sender.Username, "uri", uri)
		if err != nil {
			b.Send(m.Sender, "error while adding user details")
			return
		}
		userWallet, err := db.CheckUserOptions(m.Sender.Username, "wallet")
		if err != nil {
			b.Send(m.Sender, constants.CheckWalletOptionsError)
			return
		}
		log.Printf("your wallet address: %s", userWallet)

		b.Send(m.Sender, constants.Success, &telebot.ReplyMarkup{
			InlineKeyboard: inlineButton(nodes[idx].Username, uri),
		})
		return
	}
	b.Send(m.Sender, "invalid TXN Hash. Please try again")
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
	b.Send(m.Sender, "please send " + strconv.Itoa(Amount) + " SENTS to the following address and submit the tx hash here: ")
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

func findTxByHash(txHash, walletAddr string, m *telebot.Message) bool {

	wallet := "0x" + constants.ZFill + strings.TrimLeft(walletAddr, "0x")
	uri := constants.TEST_SENT_URI + wallet + constants.TEST_SENT_URI2 + wallet
	//log.Println("length of the wallet := ", len(wallet))
	resp, err := http.Get(uri)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	var body models.TxReceiptList
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Println(err)
		return false
	}

	var val bool
	//log.Println("decoded: \n", body)
	w, err := db.CheckUserOptions(m.Sender.Username, "wallet")
	if err != nil {
		return false
	}
	userWallet := fmt.Sprintf("%s", w)
	for _, txReceipt := range body.Results {

		//log.Println("insane")
		if txReceipt.TransactionHash == txHash {
			//log.Println("insane22")
			nodeWallet := "0xceb5bc384012f0eebee119d82a24925c47714fe3"
			//log.Println("insane2")

			okWallet :=  strings.EqualFold(txReceipt.Topics[1], "0x" + constants.ZFill + strings.TrimLeft(userWallet, "0x"))
			//log.Println("insane3")
			color.Red("%s", txReceipt.Topics[1])
			color.Green("%s", "0x" + constants.ZFill + strings.TrimLeft(userWallet, "0x"))
			okRecipient := strings.EqualFold(txReceipt.Topics[2], "0x" + constants.ZFill + strings.TrimLeft(nodeWallet, "0x"))
			//log.Println("insane4")
			okAmount := false
			if Duration == 10 {
				okAmount = hex2int(txReceipt.Data) == uint64(200000000)
			} else if Duration == 30 {
				okAmount = hex2int(txReceipt.Data) == uint64(400000000)
			} else {
				okAmount = hex2int(txReceipt.Data) == uint64(600000000)
			}
			log.Println("insane5", okWallet, okRecipient, okAmount)

			if okWallet && okRecipient && okAmount {
				val = true
			}
			//return false
			//log.Println("\n\n money: ", hex2int(txReceipt.Data))
		}
	}
	//log.Println("came here")
	return val
}

func handleWalletAddress(b *telebot.Bot, m *telebot.Message) {
	replyButtons := [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{
				Text: "10 Days",
			},
		},
		{
			telebot.ReplyButton{
				Text: "30 Days",
			},
		},
		{
			telebot.ReplyButton{
				Text: "90 Days",
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

func ExistingUser(username string) bool {
	ok, err := db.CheckUserOptions(username, "auth")
	if err != nil {
		log.Println("error in user auth", err)
		return false
	}
	isVerified, err := strconv.ParseBool(fmt.Sprintf("%s", ok))
	if err != nil {
		log.Println("error in user auth", err)
		return false
	}
	if isVerified {
		_, err = db.CheckUserOptions(username, "uri")
		if err != nil {
			log.Println("error in user auth", err)
			return false
		}
		isVerified = true
	}
	return isVerified
}

func RemoveUser(username string) bool {
	ok := false
	err := db.RemoveUser(username)
	if err != nil {
		log.Println("error while removing user route: ", err)
		return ok
	}

	return !ok
}