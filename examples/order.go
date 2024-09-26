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
	ID        string
	CreatedAt string
	UpdatedAt string
	Version   string
	Amount    string
}{Version: "version", Amount: "amount", ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt"}

// Entitsegewgwegewgewgwegwegy1Fields @qlGenerated
var Entity1Fields = struct {
	source    types.Entity1
	UpdatedAt string
	Version   string
	Amount    string
	ID        string
	CreatedAt string
}{UpdatedAt: "updatedAt", Version: "version", Amount: "amount", ID: "_id", CreatedAt: "createdAt"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @qlGenerated
var Entity2Fields = struct {
	source    types.Entity1
	Version   string
	Amount    string
	ID        string
	CreatedAt string
	UpdatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @qlGenerated
var order2Fields = struct {
	source    Order
	ID        string
	CreatedAt string
	UpdatedAt string
	Version   string
	Amount    string
}{UpdatedAt: "updatedAt", Version: "version", Amount: "amount", ID: "_id", CreatedAt: "createdAt"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var order3Fields = struct {
	source    Order
	ID        string
	CreatedAt string
	UpdatedAt string
	Version   string
	Amount    string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var userFields = struct {
	source                 User
	Online                 string
	Amount                 string
	UsernameUpdateTimes    string
	Order_Entity_ID        string
	Order_Entity_UpdatedAt string
	ID                     string
	UpdatedAt              string
	Version                string
	BalanceVersion         string
	Order_Amount           string
	CreatedAt              string
	VipLevel               string
	Order_Entity_CreatedAt string
	Order_Entity_Version   string
}{VipLevel: "vipLevel", Order_Entity_CreatedAt: "order.createdAt", Order_Entity_Version: "order.version", CreatedAt: "createdAt", Online: "online", UsernameUpdateTimes: "userNameUpdateTimes", Order_Entity_ID: "order._id", Order_Entity_UpdatedAt: "order.updatedAt", Amount: "amount", UpdatedAt: "updatedAt", Version: "version", BalanceVersion: "balanceVersion", Order_Amount: "order.amount", ID: "_id"}
