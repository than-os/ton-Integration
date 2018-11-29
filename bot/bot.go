package main

import (
	"encoding/json"
	"github.com/jasonlvhit/gocron"
	"github.com/than-os/sent-dante/utils"
	"log"
	"net/http"
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
		log.Fatalf("could not read ENV VARS. shutting down gracefully \n%v", err)
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
		Token:  "632780635:AAFLAW2XDppFLTxm1qXKqN614zFEQyjV-HU", //@et_socks_bot
		Poller: &telebot.LongPoller{Timeout: time.Second * 10},
	})

	if e != nil {
		log.Fatal(e)
	}

	s := gocron.NewScheduler()
	go func() {
		s.Every(3).Hours().Do(utils.RemoveExpiredUsers)
		//s.Every(3).Seconds().Do(dbo.FindAll)
		color.Red("%s", "running the job")
		<-s.Start()
	}()

	//
	//updates := b.Updates
	//
	//for u := range updates {
	//	fmt.Println("updates: ", u.Message)
	//}
	//b.Handle(service.ListHandle(b, Nodes))
	b.Handle(service.StartHandle(b, Nodes))
	b.Handle(service.AuthUserList(b))
	service.TestHandle(b)
	color.Green("%s", "started Telegram Bot API successfully")
	b.Start()
}
