package dbo

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/than-os/sent-dante/models"
	"log"
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

	k := []byte(key+"-"+userName)
	v := []byte(value)
	err := ldb.Put(k, v, nil)
	// defer ldb.Close()
	if err != nil {
		log.Println("error in add: ", err)
		return "", err
	}
	return "Updated Successfully", nil
}

func (b *SentinelBot) Iterate() []models.ExpiredUsers {
	itr := ldb.NewIterator(util.BytesPrefix([]byte("timestamp")), nil)

	var usersWithTimestamp []models.ExpiredUsers
	for itr.Next() {
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
func (b *SentinelBot) CheckUserOptions(userName, key string) ([]byte, error) {

	resp, err := ldb.Get([]byte(key+"-"+userName), nil)
	return resp, err
}

//func FindAll() {
//	itr := ldb.NewIterator(nil, nil)
//
//	for itr.Next() {
//		log.Printf("alluser Key %s\n", itr.Key())
//		log.Printf("alluser value %s\n", itr.Value())
//	}
//
//	itr.Release()
//	err := itr.Error()
//
//	log.Println("erorr in all user", err)
//
//}

func (b *SentinelBot) RemoveUser(username string) error {
	err := ldb.Delete([]byte("timestamp--"+username), nil)
	err = ldb.Delete([]byte("auth-"+username), nil)
	err = ldb.Delete([]byte("node-"+username), nil)
	err = ldb.Delete([]byte("password-"+username), nil)
	err  = ldb.Delete([]byte("bw-"+username), nil)
	err  = ldb.Delete([]byte("uri-"+username), nil)
	return err
}