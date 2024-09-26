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
	source    Order
	Amount    string
	CreatedAt string
	ID        string
	UpdatedAt string
	Version   string
}{Amount: "amount", CreatedAt: "createdAt", ID: "_id", UpdatedAt: "updatedAt", Version: "version"}

// Entitsegewgwegewgewgwegwegy1Fields @qlGenerated
var Entity1Fields = struct {
	source    types.Entity1
	Amount    string
	CreatedAt string
	ID        string
	UpdatedAt string
	Version   string
}{Amount: "amount", CreatedAt: "createdAt", ID: "_id", UpdatedAt: "updatedAt", Version: "version"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @qlGenerated
var Entity2Fields = struct {
	source    types.Entity1
	Amount    string
	CreatedAt string
	ID        string
	UpdatedAt string
	Version   string
}{Amount: "amount", CreatedAt: "createdAt", ID: "_id", UpdatedAt: "updatedAt", Version: "version"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @qlGenerated
var order2Fields = struct {
	source    Order
	Amount    string
	CreatedAt string
	ID        string
	UpdatedAt string
	Version   string
}{Amount: "amount", CreatedAt: "createdAt", ID: "_id", UpdatedAt: "updatedAt", Version: "version"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var order3Fields = struct {
	source    Order
	Amount    string
	CreatedAt string
	ID        string
	UpdatedAt string
	Version   string
}{Amount: "amount", CreatedAt: "createdAt", ID: "_id", UpdatedAt: "updatedAt", Version: "version"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var userFields = struct {
	source                 User
	Amount                 string
	BalanceVersion         string
	CreatedAt              string
	ID                     string
	Online                 string
	Order_Amount           string
	Order_Entity_CreatedAt string
	Order_Entity_ID        string
	Order_Entity_UpdatedAt string
	Order_Entity_Version   string
	UpdatedAt              string
	UsernameUpdateTimes    string
	Version                string
	VipLevel               string
}{Amount: "amount", BalanceVersion: "balanceVersion", CreatedAt: "createdAt", ID: "_id", Online: "online", Order_Amount: "order.amount", Order_Entity_CreatedAt: "order.createdAt", Order_Entity_ID: "order._id", Order_Entity_UpdatedAt: "order.updatedAt", Order_Entity_Version: "order.version", UpdatedAt: "updatedAt", UsernameUpdateTimes: "userNameUpdateTimes", Version: "version", VipLevel: "vipLevel"}
