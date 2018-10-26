package dbo

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type SentinelBot struct {
	Server   string
	Database string
}

//var db *mgo.Database
var ldb *leveldb.DB

//func (b *SentinelBot) NewSession() {
//	session, err := mgo.Dial(b.Server)
//	if err != nil {
//		log.Fatal("failed to establish a connection with the database")
//	}
//	db = session.DB(t.Database)
//}

func (b *SentinelBot) NewLevelDB() {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		log.Fatal(err)
	}
	ldb = db
}

func (b *SentinelBot) AddUserData(userName, key, value string) (string, error) {

	k := []byte(userName+"-"+key)
	v := []byte(value)
	err := ldb.Put(k, v, nil)
	// defer ldb.Close()
	if err != nil {
		log.Println("error in add: ", err)
		return "", err
	}
	return "Updated Successfully", nil
}

func (b *SentinelBot) CheckUserOptions(userName, key string) ([]byte, error) {

	resp, err := ldb.Get([]byte(userName+"-"+key), nil)
	return resp, err
}
