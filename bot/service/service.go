package service

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/than-os/sent-dante/bot/dbo"
	"github.com/than-os/sent-dante/constants"
	"github.com/than-os/sent-dante/models"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	NODE_ID = ""
	db dbo.SentinelBot
)

func StartHandle(b *telebot.Bot, nodes []models.TONNode) (string, func(*telebot.Message)) {

	//utils.RemoveExpiredUsers()

	fnc := func(m *telebot.Message) {
		//t := time.Now().Add(time.Minute)
		//db.AddUserData(m.Sender.Username, "timestamp-", t.Format(time.RFC3339))
		if ExistingUser(m.Sender.Username) {
			b.Send(m.Sender, "you already have a node assigned to your username. Please use /mynode to access it")
			return
		}
		b.Send(m.Sender, "Hey, "+m.Sender.Username+`. Welcome to the Sentinel Socks5 Bot for Telegram.
Please share your Ethereum wallet address that will be used for payments to this bot.`)
		// b.Send(m.Sender, "1. 100 MB")
		// b.Send(m.Sender, "2. 500 MB")
		// b.Send(m.Sender, "3. 1 GB")
	}

	fnc2 := func(m *telebot.Message) {
		log.Println("update: ", m.Text)

		NodeId := regexp.MustCompile(`(n[0-9])\w*`)
		nodeId := func() string {
			if NodeId.MatchString(m.Text) {
				return m.Text
			}
			return ""
		}
		switch m.Text {
		case validWalletAddr(m):
			_, err := db.AddUserData(m.Sender.Username, "walletAddr", m.Text)
			if err != nil {
				b.Send(m.Sender, "could not store your wallet")
				return
			}
			b.Send(m.Sender, "got your wallet")
		case nodeId():
			//uri := "https://t.me/socks?server=" + nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].Password

			b.Send(m.Sender, "Please deposit 5 SENTS to xxxx wallet and submit tx hash here", inlineButton("connect to sentinel network", "http://google.com"))
		default:
			b.Send(m.Sender, "here are few options for you. 1) suicide")
		}
	}

	_ = fnc2
	return "/start", fnc
	//return "/start", fnc
}

func ListHandle(b *telebot.Bot, nodes []models.TONNode) (string, func(message *telebot.Message)) {
	fnc := func(m *telebot.Message) {
		if !m.Private() {
			return
		}
		for idx, node := range nodes {
			uri := "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
			//geoLocation, err := utils.GetGeoLocation(node.IPAddr)
			//if err != nil {
			//	b.Send(m.Sender, "Error Occurred: "+err.Error())
			//	return
			//}
			b.Send(m.Sender, strconv.Itoa(idx+1), &telebot.ReplyMarkup{
				ResizeReplyKeyboard: true,
				InlineKeyboard:      inlineButton("Connect To "+node.Username, uri),
				//InlineKeyboard: inlineKeys,
				Selective: true,
			})
		}
	}

	return "/list", fnc
}

func getNodes() (Nodes []models.TONNode) {
	resp, err := http.Get("https://ton.sentinelgroup.io/all")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&Nodes); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return Nodes
}

func TestHandle(b *telebot.Bot) {
	nodes := getNodes()
	b.Handle(telebot.OnText, func(m *telebot.Message) {
		addr := func(*telebot.Message) string {
			if common.IsHexAddress(m.Text) {
				return m.Text
			}
			return ""
		}(m)
		log.Println("addr: ", addr)

		log.Println("update: ", m.Text)
		if ExistingUser(m.Sender.Username) {
			b.Send(m.Sender, "you already have a node assigned to your username. Please use /mynode to access it")
			return
		}
		switch m.Text {
		case checkWalletAddress(m):
			handleWalletAddress(b, m)
		case "10 Days", "30 Days", "90 Days":
			handleBandwidth(b, m, nodes)
		case nodeToInt(m):
			handleNodeId(b, m, nodes)
		case validHex(b, m):
			handleTxHash(b, m, nodes)
		case addr:
			log.Println("inside")
			_, err := db.AddUserData(m.Sender.Username, "walletAddr", m.Text)
			if err != nil {
				b.Send(m.Sender, "could not store your wallet")
				return
			}
			b.Send(m.Sender, "got your wallet")
		default:
			b.Send(m.Sender, "here are few options for you."  +
				"\n1. Contact Sentinel Dev Team @SentinelNodeNetwork" +
				"\n2. Take a deep breath and retry again by sending /start")
		}
	})
}

func checkWalletAddress(m *telebot.Message) string {
	if common.IsHexAddress(m.Text) {
		return m.Text
		//return true
	}
	return ""
}

func AuthUserList(b *telebot.Bot) (string, func(*telebot.Message)) {

	route := "/mynode"
	fnc := func(m *telebot.Message) {
		resp, err := db.CheckUserOptions(m.Sender.Username, "auth")
		if err != nil {
			b.Send(m.Sender, "you are not authorised to access the Sentinel SOCKS5 node.")
			return
		}
		ok, err := strconv.ParseBool(fmt.Sprintf("%s", resp))
		if err != nil {
			b.Send(m.Sender, "you are not authorised to access the Sentinel SOCKS5 node.")
			return
		}
		if ok {
			u, err := db.CheckUserOptions(m.Sender.Username, "uri")
			if err != nil {
				b.Send(m.Sender, "unable to get the node attached to the user")
				return
			}
			uri := fmt.Sprintf("%s", u)
			b.Send(m.Sender, constants.WarningMsg)
			//b.Send(m.Sender, fmt.Sprintf("%s", r))
			b.Send(m.Sender, "here's your node. Enjoy unrestricted internet on Telegram from @sentinel_co", &telebot.ReplyMarkup{
				InlineKeyboard: inlineButton("Sentinel SOCKS5 Node", uri),
				ResizeReplyKeyboard: true,
			})
		}

	}


	return route, fnc
}

func inlineButton(text, url string) [][]telebot.InlineButton {
	return [][]telebot.InlineButton{
		{
			telebot.InlineButton{
				Text: text,
				URL:  url,
			},
		},
	}
}

