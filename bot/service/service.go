package service

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/than-os/sent-dante/bot/dbo"
	"github.com/than-os/sent-dante/models"
	"gopkg.in/tucnak/telebot.v2"
)

var (
	NODE_ID = ""
	db      dbo.SentinelBot
)

func StartHandle(b *telebot.Bot, nodes []models.TONNode) (string, func(*telebot.Message)) {

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

	fnc := func(m *telebot.Message) {
		b.Send(m.Sender, "Hey, "+m.Sender.Username+`. Welcome to the Sentinel Socks5 Bot for Telegram.
Please select an option <number> in the format of <1> for 100 MB describing how Much Bandwidth Do you need?`, &telebot.ReplyMarkup{
			ReplyKeyboard:       replyButtons,
			ResizeReplyKeyboard: true,
			OneTimeKeyboard: true,
			ReplyKeyboardRemove: true,
		})
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
		switch m.Text {
		case "100 MB", "500 MB", "1 GB":
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
			b.Send(m.Sender, `here are few options for you. \n1. Contact Sentinel Dev Team @SentinelNodeNetwork
\n2. Take a deep breath and retry again by sending /start`)
		}
	})
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
