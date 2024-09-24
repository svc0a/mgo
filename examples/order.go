package examples

import (
	"github.com/svc0a/mgo/examples/types"
)

type Order struct {
	types.Entity `bson:",inline"`
	Amount       float64 `json:"amount" bson:"amount"`
}

// dsgergergergresdsdfds @mongoGenerated
var orderFields = struct {
	source    Order
	ID        string
	Amount    string
	CreatedAt string
	UpdatedAt string

	Version string
}{
	ID: "OrderID", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount",
}

// Entitsegewgwegewgewgwegwegy1Fields @mongoGenerated
var Entity1Fields = struct {
	source    types.Entity1
	ID        string `json:"id" bson:"_id"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string

	Version string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @mongoGenerated
var Entity2Fields = struct {
	source    types.Entity1
	ID        string `json:"id" bson:"_id"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string

	Version string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @mongoGenerated
var order2Fields = struct {
	source    Order
	ID        string
	Amount    string
	UpdatedAt string
	Version   string

	CreatedAt string
}{
	ID: "OrderID", Amount: "amount", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version",
}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @mongoGenerated
var order3Fields = struct {
	source    Order
	ID        string
	Amount    string
	CreatedAt string
	UpdatedAt string
	Version   string
}{
	ID: "OrderID", Amount: "amount", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version",
}
