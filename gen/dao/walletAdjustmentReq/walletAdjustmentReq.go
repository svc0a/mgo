package walletAdjustmentReq

import "github.com/svc0a/mgo/gen"

func New() *walletAdjustmentReq {
	return &walletAdjustmentReq{}
}

type walletAdjustmentReq struct{}

func (m *walletAdjustmentReq) TransactionId() gen.Field[string] {
	return gen.Field[string]{
		Name: "TransactionId",
		Bson: "transactionId",
	}
}

func (m *walletAdjustmentReq) CommonReq() *commonReq {
	return &commonReq{}
}
