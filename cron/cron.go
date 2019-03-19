package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jasonlvhit/gocron"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/than-os/sent-dante/models"
	"log"
	"math"
	"strings"
	"time"
)

var ldb *leveldb.DB

func init() {
	//func (b *SentinelBot) NewLevelDB() {
		db, err := leveldb.OpenFile("../bot/db", nil)
		if err != nil {
			log.Fatal(err)
		}
		color.Green("%s", "connected to db successfully")
		ldb = db
	//}

}
func main() {
	//for {
			s := gocron.NewScheduler()
			s.Every(2).Seconds().Do(RemoveExpiredUsers)
			color.Cyan("%s", "it came here too")
			<-s.Start()
		//time.Sleep(time.Second * 5)
	//}
}

func RemoveExpiredUsers() {
	//for {
	usersWithTimestamp := Iterate()
	today := time.Now()
	for _, user := range usersWithTimestamp {
		log.Println("user with timestamp: ", user.Value)
		userExpiryTime, err := time.Parse(time.RFC3339, user.Value)
		if err != nil {
			log.Println("error in parsing time: ", err)
			break
			//return err.Error()
		}
		log.Println("user key: ", fmt.Sprintf("%s", user.Key))
		if math.Signbit(userExpiryTime.Sub(today).Hours()) {
			log.Println("came insidem no luck", userExpiryTime.Sub(today).Hours())
			username := strings.TrimLeft(fmt.Sprintf("%s", user.Key), "timestamp")
			log.Printf("this one: %s and this two: %s", username[2:], user.Key )
			err := RemoveUser(username[2:])
			if err != nil {
				log.Println("this error: ",err)
				break
				//return err.Error()
			}
			log.Println("this error: ",err)

			//return return"deleted"
		}

		//log.Println("time left for user: ", diff.Hours())

	}
	//return "deleted"
	//time.Sleep(time.Minute)
	//}
	//log.Println("userswithtimestamp: ", usersWithTimestamp )
}



func Iterate() []models.ExpiredUsers {
	itr := ldb.NewIterator(util.BytesPrefix([]byte("timestamp")), nil)

	var usersWithTimestamp []models.ExpiredUsers
	for itr.Next() {
		//today := time.Now().UTC()
		log.Println("this guy: ", itr.Value())
		usersWithTimestamp = append(usersWithTimestamp, models.ExpiredUsers{
			Key: fmt.Sprintf("%s", itr.Key()), Value: fmt.Sprintf("%s", itr.Value()),
		})
	}
	itr.Release()
	err := itr.Error()
	if err != nil {
		return usersWithTimestamp
	}

	return usersWithTimestamp
}

func RemoveUser(username string) error {
	log.Println("maybe? \n", username)
	err := ldb.Delete([]byte("timestamp--"+username), nil)
	err = ldb.Delete([]byte("auth-"+username), nil)
	err = ldb.Delete([]byte("node-"+username), nil)
	err = ldb.Delete([]byte("password-"+username), nil)
	err  = ldb.Delete([]byte("bw-"+username), nil)
	err  = ldb.Delete([]byte("uri-"+username), nil)
	return err
}