package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/than-os/sent-dante/bot/dbo"
	"github.com/than-os/sent-dante/bot/service"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
	"github.com/than-os/sent-dante/models"
	"gopkg.in/tucnak/telebot.v2"
)

var ldb dbo.SentinelBot
var (
	uri   string
	text  string
	Nodes []models.TONNode
)

func init() {
	ldb.Database = "SentinelBot"
	ldb.Server = "localhost"
	ldb.NewLevelDB()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not read ENV VARS. Shutting Down!!! \n%v", err)
	}
	resp, err := http.Get("https://ton.sentinelgroup.io/all")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&Nodes); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	uri = "https://t.me/socks?server=" + Nodes[0].IPAddr + "&port=" + strconv.Itoa(Nodes[0].Port) + "&user=" + Nodes[0].Username + "&pass=" + Nodes[0].Password
}

func main() {

	b, e := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TEST_API_KEY2"), //@et_socks_bot
		Poller: &telebot.LongPoller{Timeout: time.Second * 10},
	})

	if e != nil {
		log.Fatal(e)
	}

	inlineConnectBtn := telebot.InlineButton{
		Text: "Auto Select",
		URL:  uri,
	}
	inlineListBtn := telebot.InlineButton{
		Text: "List",
	}

	//replyBtn := telebot.ReplyButton{
	//	Text: "Auto Select",
	//}
	//listBtn := telebot.ReplyButton{
	//	Text: "List",
	//}
	//replyKeys := [][]telebot.ReplyButton{
	//	{replyBtn},
	//	{listBtn},
	//}

	b.Handle(&inlineListBtn, func(c *telebot.Callback) {
		b.Send(c.Sender, "hello world")
	})

	b.Handle(&inlineConnectBtn, func(c *telebot.Callback) {
		resp, err := http.Get("https://ton.sentinelgroup.io/all")
		if err != nil {
			b.Respond(c, &telebot.CallbackResponse{
				Text:      err.Error(),
				ShowAlert: true,
			})
			log.Fatal(err)
		}
		if err := json.NewDecoder(resp.Body).Decode(&Nodes); err != nil {
			b.Respond(c, &telebot.CallbackResponse{
				Text:      err.Error(),
				ShowAlert: true,
			})
			log.Fatal(err)
		}
		defer resp.Body.Close()

		uri = "https://t.me/socks?server=" + Nodes[0].IPAddr + "&port=" + strconv.Itoa(Nodes[0].Port) + "&user=" + Nodes[0].Username + "&pass=" + Nodes[0].Password
		inlineConnectBtn.Text = Nodes[0].IPAddr
		inlineConnectBtn.URL = uri

		b.Respond(c, &telebot.CallbackResponse{
			Text:      uri,
			ShowAlert: true,
		})
	})
	//
	//updates := b.Updates
	//
	//for u := range updates {
	//	fmt.Println("updates: ", u.Message)
	//}

	b.Handle(service.ListHandle(b, Nodes))
	b.Handle(service.StartHandle(b, Nodes))

	color.Green("%s", "started Telegram Bot API successfully")
	service.TestHandle(b)
	b.Start()
}
