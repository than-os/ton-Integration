package utils

import (
	"encoding/json"
	"github.com/than-os/sent-dante/constants"
	"github.com/than-os/sent-dante/models"
	"log"
	"net/http"
	"os"
)

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

func processTxHash(txHash string) (models.TONNode, models.TXDetails, error) {
	var node models.TONNode
	var txDetails models.TXDetails
	resp, err := http.Get(constants.TX_BY_HASH + txHash + "&apikey=" + os.Getenv("ETH_SCAN_API_KEY"))
	if err != nil {
		log.Println(err)
		return node, txDetails, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&txDetails); err != nil {
		log.Println(err)
		return node, txDetails, err
	}


	return node, txDetails, err
}


//func main() {
//	//d, e := GetGeoLocation("185.222.24.146")
//	//if e != nil {
//	//	fmt.Println(e.Error())
//	//}
//	//fmt.Println(d.RegionName)
//	node, txInfo, err := processTxHash("0x4be1e6c88c7b3554ee9a9aab9590d511a9fe9213a2db10e1685355d9c2fc1421")
//	if err != nil {
//		log.Println("error: ", err)
//		return
//	}
//
//	_=node
//
//	fmt.Printf("FROM: %s To: %s Amount: %s", )
//
//	}