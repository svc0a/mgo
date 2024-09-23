package examples

import "github.com/svc0a/mgo/examples/types"

type Order struct {
	types.Entity `bson:",inline"`
	Amount       float64 `json:"amount" bson:"amount"`
}

// @generatedMongo keys mapping
var orderFields = struct {
	source Order
	//generated fields
	ID     string
	Amount string
}{}
