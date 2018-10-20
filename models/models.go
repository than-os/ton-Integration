package models

import "github.com/globalsign/mgo/bson"

type TONNode struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	IPAddr string `json:"ipAddr" bson:"ipAddr"`
	Port int `json:"port" bson:"port"`
	Username string `json:"userName" bson:"userName"`
	Password string `json:"password" bson:"password"`
}
