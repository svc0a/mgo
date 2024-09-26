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
	source Order
	ID     string

	CreatedAt string
	UpdatedAt string
	Version   string
	Amount    string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// Entitsegewgwegewgewgwegwegy1Fields @qlGenerated
var Entity1Fields = struct {
	source types.Entity1
	Amount string

	ID        string
	CreatedAt string
	UpdatedAt string
	Version   string
}{Version: "version", Amount: "amount", ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt"}

// Entity1Fsdsefsdfdffffwegewgwgh4h43h3h43h43hields @qlGenerated
var Entity2Fields = struct {
	source types.Entity1
	Amount string

	ID        string
	CreatedAt string
	UpdatedAt string
	Version   string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// ordesegwegweghwehweeghwegewgewgwesgewgerFields @qlGenerated
var order2Fields = struct {
	source Order
	ID     string

	CreatedAt string
	UpdatedAt string
	Version   string
	Amount    string
}{Amount: "amount", ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var order3Fields = struct {
	source    Order
	UpdatedAt string
	Version   string
	Amount    string
	ID        string
	CreatedAt string
}{ID: "_id", CreatedAt: "createdAt", UpdatedAt: "updatedAt", Version: "version", Amount: "amount"}

// orderFfdbdhergefifhfhfhbefjcbcjdksihdields @qlGenerated
var userFields = struct {
	source                 User
	UsernameUpdateTimes    string
	VipLevel               string
	Order_Entity_UpdatedAt string
	ID                     string
	CreatedAt              string
	UpdatedAt              string
	Version                string
	BalanceVersion         string
	Online                 string
	Order_Entity_ID        string
	Order_Entity_CreatedAt string
	Order_Entity_Version   string
	Order_Amount           string
}{CreatedAt: "createdAt", UpdatedAt: "updatedAt", Order_Entity_ID: "order._id", Order_Entity_CreatedAt: "order.createdAt", Order_Amount: "order.amount", Order_Entity_Version: "order.version", VipLevel: "vipLevel", Order_Entity_UpdatedAt: "order.updatedAt", UsernameUpdateTimes: "userNameUpdateTimes", Version: "version", BalanceVersion: "balanceVersion", Online: "online", ID: "_id"}
