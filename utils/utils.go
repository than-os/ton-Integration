package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/than-os/sent-dante/bot/dbo"
	"github.com/than-os/sent-dante/models"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var db dbo.SentinelBot


type GeoLocation struct {
	//As string `json:"as"`
	City string `json:"city_name"`
	Country string `json:"country_name"`
	CountryCode string `json:"country_code"`
	//Isp string `json:"isp"`
	//Lat float64 `json:"lat"`
	//Lon float64 `json:"lon"`
	//Query string `json:"query"`
	//Region string `json:"region"`
	RegionName string `json:"region_name"`
	//Status string `json:"status"`
	//TimeZone string `json:"timezone"`
	//Zip string `json:"zip"`
}

func GetGeoLocation(ipAddr string) (GeoLocation, error) {
	gl := GeoLocation{}
	resp, err := http.Get("https://ipleak.net/json/" + ipAddr)
	if err != nil {
		log.Println("Error occurred while fetching GeoLocation: ", err.Error())
		return gl, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&gl); err != nil {
		log.Println("Error occurred Decoding Response Body: ", err.Error())
		return gl, err
	}

	return gl, err
}

var src = rand.NewSource(time.Now().UnixNano())
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
func StrongPassword(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func AddUser(ipAddr, userName string) error {

	err := DeleteUser(userName, ipAddr)
	if err != nil {
		log.Println("error while deleting user: ", err)
		return err
	}
	uri := "http://"+ipAddr+":30002/user"
	password := StrongPassword(6)
	_, err = db.AddUserData(userName, "password", password)
	if err != nil {
		log.Println("error while storing password: ", err)
		return err
	}

	log.Println("whats going on? ", ipAddr, uri)
	req := models.AddUser{
		Username: userName,
		Password: password,
	}
	b, e := json.Marshal(req)
	if e != nil {
		log.Println("error in marshal: ", e)
		return e
	}
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(b))
	if err  != nil {
		log.Println("error in post request: ", uri)
		return err
	}
	var res models.UserRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("error on decoding resp: ", err)
		return err
	}
	log.Println("success response: ", res)
	return err
}

func DeleteUser(username, ipAddr string) error {
	client := &http.Client{}

	uri := "http://"+ipAddr+":30002/user"
	body := models.RemoveUser{
		Username: username,
	}

	b, e := json.Marshal(body)
	if e != nil {
		log.Println("error in marshal: ", e)
		return e
	}
	// Create request
	req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("respbody: ",respBody)
	return err
}

func RemoveExpiredUsers() {
	//for {
		usersWithTimestamp := db.Iterate()
		today := time.Now()
		for _, user := range usersWithTimestamp {
			//log.Println("user with timestamp: ", user.Value)
			userExpiryTime, err := time.Parse(time.RFC3339, user.Value)
			if err != nil {
				log.Println("error in parsing time: ", err)
				break
				//return err.Error()
			}
			//log.Println("user key: ", fmt.Sprintf("%s", user.Key))
			if math.Signbit(userExpiryTime.Sub(today).Hours()) {
				username := strings.TrimLeft(fmt.Sprintf("%s", user.Key), "timestamp")
				ip, err := db.CheckUserOptions(username[2:], "ipaddr")
				if err != nil {
					log.Println("error while getting node ip")
					return
				}
				err = DeleteUser(username[2:], fmt.Sprintf("%s", ip))
				if err != nil {
					log.Println("error in deleting SOCKS5 user")
					return
				}
				err = db.RemoveUser(username[2:])
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

func RemoveUserJob() {
	s := gocron.NewScheduler()
	s.Every(3).Hours().Do(RemoveExpiredUsers)
	<-s.Start()
}