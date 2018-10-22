package main

import (
	"encoding/json"
	"github.com/than-os/sent-dante/models"
	"github.com/than-os/sent-dante/utils"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	uri string
	text string
	nodes []models.TONNode
)

func init() {
	resp, err := http.Get("https://ton.sentinelgroup.io/all")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("%v", nodes)
	// always respond!

	uri = "https://t.me/socks?server="+ nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].Password
}

func main() {

	b, e := telebot.NewBot(telebot.Settings{
		//Token: "632780635:AAFLAW2XDppFLTxm1qXKqN614zFEQyjV-HU", //deployed key @et_socks_bot
		//Token: "672453988:AAGMXyZXLHVU4nh2uR5PDRNpu595Dubl-Hs", //test key http://t.me/qwerty_et_bot
		Token: "713749545:AAFHjglnYnjBAc0wLZyPSTwuIOANm5C0sTw", //official_key http://t.me/Sentinel_SOCKS5_bot
		Poller: &telebot.LongPoller{ Timeout: time.Second},
		//Updates: 20,
	})

	if e != nil {
		log.Fatal(e)
	}

	inlineConnectBtn := telebot.InlineButton{
		//Unique: "sad_moon",
		Text: "Auto Select",
		URL: uri,
		//URL: "https://t.me/socks?server=185.222.24.146&port=1090&user=sentinel&pass=$entinel@12!@",
	}
	inlineListBtn := telebot.InlineButton{
		//Unique: "sad_moon",
		Text: "List",
	}
	//array of array of inline buttons
	inlineKeys := [][]telebot.InlineButton{
		{ inlineConnectBtn },
		//{ inlineListBtn },
		//{inlineConnectBtn},
		// ...
	}
	newKeys := telebot.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeys,
	}

	replyBtn := telebot.ReplyButton{
		Text: "Auto Select",
	}
	listBtn := telebot.ReplyButton{
		Text: "List",
	}
	replyKeys := [][]telebot.ReplyButton{
		{replyBtn},
		{listBtn},
	}
	//inlineKeys = append(inlineKeys, []telebot.InlineButton{inlineListBtn})
	b.Handle(&replyBtn, func(m *telebot.Message) {
		// on reply button pressed
		uri := "https://t.me/socks?server="+ nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].Password

		b.Send(m.Sender, "Here is your auto selected SENT-TON Node", &telebot.ReplyMarkup {
			OneTimeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{ telebot.InlineButton{
					URL: uri,
					Text: "Click To Connect",
				} },
			},
		})

			//b.Respond(c, &telebot.CallbackResponse{
			//	URL: uri,
			//	Text: "well well well",
			//	ShowAlert: true,
			//})
	})
	b.Handle(&listBtn, func(m *telebot.Message) {
		// on reply button pressed
		for idx, node := range nodes {
			uri := "https://t.me/socks?server="+ node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
			geoLocation, err := utils.GetGeoLocation(node.IPAddr)
			if err != nil {
				b.Send(m.Sender, "Error Occurred: " + err.Error())
				return
			}
			b.Send(m.Sender, "SENT-TON Node #" + strconv.Itoa(idx + 1), &telebot.ReplyMarkup{
				ResizeReplyKeyboard: true,
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							//Unique: "sad_moon",
							Text: "Connect To " + geoLocation.RegionName,
							URL: uri,
							Data: node.Username,
						},
					},
				},
				//InlineKeyboard: inlineKeys,
				Selective: true,
			})
		}
	})

	b.Handle(&inlineListBtn, func(c *telebot.Callback) {
		b.Send(c.Sender, "hello world")
	})

	b.Handle(&inlineConnectBtn, func(c *telebot.Callback) {
		// on inline button pressed (callback!)

		resp, err := http.Get("https://ton.sentinelgroup.io/all")
		if err != nil {
			b.Respond(c, &telebot.CallbackResponse{
				Text: err.Error(),
				ShowAlert: true,
			})
			log.Fatal(err)
		}
		if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
			b.Respond(c, &telebot.CallbackResponse{
				Text: err.Error(),
				ShowAlert: true,
			})
			log.Fatal(err)
		}
		defer resp.Body.Close()
		log.Printf("%v", nodes)
		// always respond!

		uri = "https://t.me/socks?server="+ nodes[0].IPAddr + "&port=" + strconv.Itoa(nodes[0].Port) + "&user=" + nodes[0].Username + "&pass=" + nodes[0].Password
		inlineConnectBtn.Text = nodes[0].IPAddr
		inlineConnectBtn.URL = uri

		b.Respond(c, &telebot.CallbackResponse{
			Text: uri,
			//URL: "https://t.me/socks?server=185.222.24.146&port=1090&user=sentinel&pass=$entinel@12!",
			ShowAlert:true,
		})
	})
	_ = newKeys

	b.Handle("/list", func(m *telebot.Message) {
		if !m.Private() {
			return
		}
		for idx, node := range nodes {
			uri := "https://t.me/socks?server="+ node.IPAddr + "&port=" + strconv.Itoa(node.Port) + "&user=" + node.Username + "&pass=" + node.Password
			geoLocation, err := utils.GetGeoLocation(node.IPAddr)
			if err != nil {
				b.Send(m.Sender, "Error Occurred: " + err.Error())
				return
			}
			b.Send(m.Sender, "SENT-TON Node #" + strconv.Itoa(idx + 1), &telebot.ReplyMarkup{
				ResizeReplyKeyboard: true,
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							//Unique: "sad_moon",
							Text: "Connect To " + geoLocation.RegionName,
							URL: uri,
							Data: node.Username,
						},
					},
				},
				//InlineKeyboard: inlineKeys,
				Selective: true,
			})
		}
	})

	b.Handle("/start", func(m *telebot.Message) {
		if !m.Private() {
			return
		}

		msg := "List of all Telegram Open Network Nodes in Sentinel dNetwork"
		b.Send(m.Sender, msg, &telebot.ReplyMarkup{
			ReplyKeyboard: replyKeys,
			ResizeReplyKeyboard: true,
		})
	})

	b.Start()
}