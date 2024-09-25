package examples

import (
	"github.com/svc0a/mgo/examples/types"
)

type Order struct {
	types.Entity `bson:",inline"` // gewgwegew
	Amount       float64          `json:"amount" bson:"amount"` // gewgwegew
}

type Order2 struct {
	types.Entity `bson:",inline"` // gewgwegew
	Amount       float64          `json:"amount" bson:"amount"` // gewgwegew
}

// dsgergergergresdsdfds @mongoGenerated
var orderFields = struct {
	source  Order
	Version string

	Amount    string
	ID        string
	CreatedAt string
	UpdatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// Entitsegewgwegewgewgwegwegy1Fields @mongoGenerated
var Entity1Fields = struct {
	source types.Entity1
	ID     string

	CreatedAt string
	UpdatedAt string
	Version   string
}{Version: "version", ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @mongoGenerated
var Entity2Fields = struct {
	source    types.Entity1
	CreatedAt string

	UpdatedAt string
	Version   string
	ID        string
}{CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", ID: "_id"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @mongoGenerated
var order2Fields = struct {
	source    Order
	UpdatedAt string

	Version   string
	Amount    string
	ID        string
	CreatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @mongoGenerated
var order3Fields = struct {
	source    Order
	Version   string
	Amount    string
	ID        string
	CreatedAt string
	UpdatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}
