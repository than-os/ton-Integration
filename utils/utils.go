package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type GeoLocation struct {
	As string `json:"as"`
	City string `json:"city"`
	Country string `json:"country"`
	CountryCode string `json:"countryCode"`
	Isp string `json:"isp"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Query string `json:"query"`
	Region string `json:"region"`
	RegionName string `json:"regionName"`
	//Status string `json:"status"`
	//TimeZone string `json:"timezone"`
	//Zip string `json:"zip"`
}

func GetGeoLocation(ipAddr string) (GeoLocation, error) {
	gl := GeoLocation{}
	resp, err := http.Get("http://ip-api.com/json/" + ipAddr)
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

//func main() {
//	d, e := GetGeoLocation("185.222.24.146")
//	if e != nil {
//		fmt.Println(e.Error())
//	}
//	fmt.Println(d.RegionName)
//
//	}