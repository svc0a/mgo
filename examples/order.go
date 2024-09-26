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
	source           Order
	Entity_Version   string
	Amount           string
	Entity_ID        string
	Entity_CreatedAt string
	Entity_UpdatedAt string
}{Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_UpdatedAt: "EntityUpdatedAt", Entity_Version: "EntityVersion", Amount: "Amount"}

// Entitsegewgwegewgewgwegwegy1Fields @qlGenerated
var Entity1Fields = struct {
	source           types.Entity1
	Entity_CreatedAt string
	Entity_Version   string
	ID               string
	Entity_ID        string
	Amount           string
	CreatedAt        string
	UpdatedAt        string
	Version          string
	Entity_UpdatedAt string
}{Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_Version: "EntityVersion", ID: "ID", Version: "Version", Entity_UpdatedAt: "EntityUpdatedAt", Amount: "Amount", CreatedAt: "CreatedAt", UpdatedAt: "UpdatedAt"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @qlGenerated
var Entity2Fields = struct {
	source           types.Entity1
	CreatedAt        string
	UpdatedAt        string
	Version          string
	Entity_UpdatedAt string
	Amount           string
	Entity_Version   string
	ID               string
	Entity_ID        string
	Entity_CreatedAt string
}{Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_Version: "EntityVersion", ID: "ID", Entity_UpdatedAt: "EntityUpdatedAt", Amount: "Amount", CreatedAt: "CreatedAt", UpdatedAt: "UpdatedAt", Version: "Version"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @qlGenerated
var order2Fields = struct {
	source           Order
	Entity_UpdatedAt string
	Amount           string
	CreatedAt        string
	UpdatedAt        string
	Version          string
	Entity_ID        string
	Entity_CreatedAt string
	Entity_Version   string
	ID               string
}{Entity_ID: "EntityID", Entity_CreatedAt: "EntityCreatedAt", Entity_Version: "EntityVersion", ID: "ID", Version: "Version", Entity_UpdatedAt: "EntityUpdatedAt", Amount: "Amount", CreatedAt: "CreatedAt", UpdatedAt: "UpdatedAt"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var order3Fields = struct {
	source           Order
	Entity_Version   string
	ID               string
	Entity_ID        string
	Entity_CreatedAt string
	CreatedAt        string
	UpdatedAt        string
	Version          string
	Entity_UpdatedAt string
	Amount           string
}{Amount: "Amount", CreatedAt: "CreatedAt", UpdatedAt: "UpdatedAt", Version: "Version", Entity_UpdatedAt: "EntityUpdatedAt",

	Entity_CreatedAt: "EntityCreatedAt", Entity_Version: "EntityVersion", ID: "ID", Entity_ID: "EntityID"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var userFields = struct {
	source                 User
	Order_Entity_CreatedAt string
	Entity_UpdatedAt       string
	Amount                 string
	UpdatedAt              string
	Online                 string
	Entity_ID              string
	Entity_Version         string
	VipLevel               string
	Order_Amount           string
	CreatedAt              string
	Version                string
	UsernameUpdateTimes    string
	Order_Entity_UpdatedAt string
	Order_Entity_Version   string
	Entity_CreatedAt       string
	ID                     string
	BalanceVersion         string
	Order_Entity_ID        string
}{UpdatedAt: "UpdatedAt", Online: "Online", Order_Entity_CreatedAt: "OrderEntityCreatedAt", Entity_UpdatedAt: "EntityUpdatedAt", Amount: "Amount", VipLevel: "VipLevel", Entity_ID: "EntityID", Entity_Version: "EntityVersion", UsernameUpdateTimes: "UsernameUpdateTimes", Order_Entity_UpdatedAt: "OrderEntityUpdatedAt", Order_Amount: "OrderAmount", CreatedAt: "CreatedAt", Version: "Version", BalanceVersion: "BalanceVersion", Order_Entity_ID: "OrderEntityID", Order_Entity_Version: "OrderEntityVersion", Entity_CreatedAt: "EntityCreatedAt", ID: "ID"}
