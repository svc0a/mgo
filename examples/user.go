package examples

import "github.com/svc0a/mgo/examples/types"

// User @generated sql keys mapping
type User struct {
	types.Entity        `bson:",inline"`
	BalanceVersion      int64 `json:"balanceVersion" bson:"balanceVersion"`
	UsernameUpdateTimes int   `json:"userNameUpdateTimes" bson:"userNameUpdateTimes"`
	Online              bool  `json:"online" bson:"online"`
	VipLevel            int   `json:"vipLevel" bson:"vipLevel"`
}

// fields @generated sql keys mapping1
//var fields = struct {
//	ID              string
//	Balance         string
//	Balance_Balance string
//}{ID: "_id", Balance: "balance", Balance_Balance: "balance.balance"}

// Gender @generated sql keys mapping
type Gender string

// Order grewhwrehw

// Order @generated sql keys mapping
type Order struct {
	types.Entity `bson:",inline"`
	Amount       float64 `json:"amount" bson:"amount"`
}
