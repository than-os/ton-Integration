package service

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/than-os/sent-dante/models"
	"github.com/than-os/sent-dante/utils"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	NODE_ID = ""
)

func StartHandle(b *telebot.Bot, nodes []models.TONNode) (string, func(*telebot.Message)) {

	fnc := func(m *telebot.Message) {
		b.Send(m.Sender, "Hey, " + m.Sender.Username + `. Welcome to the Sentinel Socks5 Bot for Telegram.
Please select an option <number> in the format of <1> for 100 MB describing how Much Bandwidth Do you need?`)
		b.Send(m.Sender, "1. 100 MB")
		b.Send(m.Sender, "2. 500 MB")
		b.Send(m.Sender, "3. 1 GB")
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
		case "1":
			for idx, node := range nodes {
				node.WalletAddress = "0xCeb5bC384012f0EEbeE119d82A24925C47714fe3"
				//uri := "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
				geo, err := utils.GetGeoLocation(node.IPAddr)

				if err != nil {
					b.Send(m.Sender, "error: ", err.Error())
					return
				}
				b.Send(m.Sender, strconv.Itoa(idx+1) + ") " + geo.Country + "hello \n world" )
				//&telebot.ReplyMarkup{
				//	ResizeReplyKeyboard: true,
				//	InlineKeyboard: [][]telebot.InlineButton{
				//		{
				//			telebot.InlineButton{
				//				//Unique: "sad_moon",
				//				Text: "Connect To " + geo.Country,
				//				URL:  uri,
				//				Data: node.Username,
				//			},
				//		},
				//	},
				//	//InlineKeyboard: inlineKeys,
				//	Selective: true,
				//}
				//)
			}
			return
		case "2":
			b.Send(m.Sender, "you got 500 mb")
			return
		case "3":
			b.Send(m.Sender, "you got 1 gb")
			return
		case nodeId() :
			//uri := "https://t.me/socks?server=" + nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].Password

			b.Send(m.Sender, "Please deposit 5 SENTS to xxxx wallet and submit tx hash here", &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: "click to connect",
							URL: "https://myetherwallet.com",
						},
					},
				},
			})
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
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							//Unique: "sad_moon",
							Text: "Connect To " + node.Username,
							URL:  uri,
							Data: node.Username,
						},
					},
				},
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
		log.Println("update: ", m.Text)
		rgx := regexp.MustCompile(`(n[0-9])\w*`)
		nodeId := func() string {
			if rgx.MatchString(m.Text) {
				return m.Text
			}
			return ""
		}

		validHex := func() string {
			_, err := hexutil.Decode(m.Text)
			if err != nil {
				b.Send(m.Sender, "invalid tx hash")
				return ""
			}
			return m.Text

		}
		switch m.Text {
		case "1", "2", "3":
			switch m.Text {
			case "1":
				b.Send(m.Sender, "You have requested for 100 MB for data via Sentinel SOCKS5 Proxy for Telegram")
			case "2":
				b.Send(m.Sender, "You have requested for 500 MB for data via Sentinel SOCKS5 Proxy for Telegram")
			case "3":
				b.Send(m.Sender, "You have requested for 1 GB for data via Sentinel SOCKS5 Proxy for Telegram")
			default:
				b.Send(m.Sender, "You've selected an invalid option")
			}
			b.Send(m.Sender, `Please select a node ID from the list below and reply in the format of 
			<n1> for Node 1, <n2> for Node 2`)
			for idx, node := range nodes {
				//uri := "https://t.me/socks?server=" + node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
				geo, err := utils.GetGeoLocation(node.IPAddr)

				if err != nil {
					b.Send(m.Sender, "error: ", err.Error())
					return
				}
				b.Send(m.Sender, strconv.Itoa(idx+1) + ". Location: " + geo.Country + "\n " + "User:" + node.Username + "\n " + "Node Wallet: " + "0xCeb5bC384012f0EEbeE119d82A24925C47714fe3" )
			}
			return
		case "begin":
			b.Send(m.Sender, "Hey, " + m.Sender.Username + " welcome to the Sentinel Socks5 Bot. How Much Bandwidth Do you need?")
			b.Send(m.Sender, "1. 100 MB")
			b.Send(m.Sender, "2. 500 MB")
			b.Send(m.Sender, "3. 1 GB")
			return
		case nodeId():
			NODE_ID = m.Text
			log.Println("in nodeid case", m.Text)
			b.Send(m.Sender, "please send 5 SENTS to the following address and submit the tx hash here: ")
			b.Send(m.Sender, "0xCeb5bC384012f0EEbeE119d82A24925C47714fe3")
		case validHex():
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
			uri := "https://t.me/socks?server=" + nodes[nodeIdx].IPAddr + "&port=" + strconv.Itoa(nodes[nodeIdx].Port) + "&user=" + nodes[nodeIdx].Username + "&pass=" + nodes[nodeIdx].Password

			b.Send(m.Sender, "Thanks for submitting the TX-HASH. We're validating it")
			b.Send(m.Sender, "Congratulations!! please click the button below to connect to the sentinel dVPN node", &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: nodes[nodeIdx].Username,
							URL: uri,

						},
					},
				},
			})
			//txInfo, err := processTxHash("0xb274ad6de167b6ab8833050a07ef7de09df8b271edb3851dddd2ec65a45a4d5f")
			//if err != nil {
			//	b.Send(m.Sender, "error occurred: ", err)
			//}
			//uri := "https://t.me/socks?server=" + nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].
		default:
			b.Send(m.Sender, "here are few options for you. \n1. Contact Sentinel Dev Team @SentinelNodeNetwork\n2. Take a deep breath and retry again by sending /start")
		}
	})
}