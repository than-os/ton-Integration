package dbo

import (
	"github.com/globalsign/mgo"
	"github.com/than-os/sent-dante/models"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type TON struct {
	Server string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "nodes"
)

func (t *TON) NewSession()  {
	session, err := mgo.Dial(t.Server)
	if err != nil {
		log.Fatal("failed to establish a connection with the database")
	}
	db = session.DB(t.Database)
}

func (t *TON) FindAllTonNodes() ([]models.TONNode, error) {
	var nodes []models.TONNode
	err := db.C(COLLECTION).Find(nil).All(&nodes)
	//defer db.Session.Close()
	return nodes, err
}

func (t *TON) FindTonNodeByID(IPAddr string) (models.TONNode, error) {
	var node models.TONNode
	q := bson.M{"ipAddr": IPAddr}
	err := db.C(COLLECTION).Find(q).One(&node)
	//defer db.Session.Close()
	return node, err
}

func (t *TON) RegisterTonNode(node models.TONNode) error {
	err := db.C(COLLECTION).Insert(&node)
	//defer db.Session.Close()
	return err
}

func (t *TON) RemoveTonNode(IPAddr string) error {
	q := bson.M{"ipAddr": IPAddr}
	err := db.C(COLLECTION).Remove(q)
	//defer db.Session.Close()
	return err
}