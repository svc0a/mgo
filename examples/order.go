package examples

import (
	"github.com/svc0a/mgo/examples/types"
)

type Order struct {
	types.Entity `bson:",inline"` // gewgwegew
	Amount       float64          `json:"amount" bson:"amount" gorm:"amount"` // gewgwegew
}

type Order2 struct {
	types.Entity `bson:",inline"` // gewgwegew
	Amount       float64          `json:"amount" bson:"amount"` // gewgwegew
}

// dsgergergergresdsdfds @qlGenerated
var orderFields = struct {
	source  Order
	Version string

	Amount           string
	ID               string
	CreatedAt        string
	UpdatedAt        string
	Entity_UpdatedAt string
	Entity_Version   string
	Entity_ID        string
	Entity_CreatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount", Entity_Version: "EntityVersion", Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_UpdatedAt: "EntityUpdatedAt"}

// Entitsegewgwegewgewgwegwegy1Fields @qlGenerated
var Entity1Fields = struct {
	source types.Entity1
	ID     string

	CreatedAt string
	UpdatedAt string
	Version   string
}{Version: "version", ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @qlGenerated
var Entity2Fields = struct {
	source    types.Entity1
	CreatedAt string

	UpdatedAt string
	Version   string
	ID        string
}{CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", ID: "_id"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @qlGenerated
var order2Fields = struct {
	source    Order
	UpdatedAt string

	Version          string
	Amount           string
	ID               string
	CreatedAt        string
	Entity_ID        string
	Entity_CreatedAt string
	Entity_UpdatedAt string
	Entity_Version   string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount", Entity_CreatedAt: "EntityCreatedAt", Entity_UpdatedAt: "EntityUpdatedAt", Entity_Version: "EntityVersion", Entity_ID: "EntityID"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var order3Fields = struct {
	source           Order
	Entity_UpdatedAt string
	Entity_Version   string
	Amount           string
	Entity_ID        string
	Entity_CreatedAt string
}{Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_UpdatedAt: "EntityUpdatedAt", Entity_Version: "EntityVersion", Amount: "Amount"}
