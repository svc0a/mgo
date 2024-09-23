package examples

import "github.com/svc0a/mgo/examples/types"

// Order @generated sql keys mapping
type Order struct {
	types.Entity `bson:",inline"`
	Amount       float64 `json:"amount" bson:"amount"`
}
