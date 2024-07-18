package walletAdjustmentReq

import "github.com/svc0a/mgo/gen"

type commonReq struct{}

func (m *commonReq) TraceId() gen.Field[string] {
	return gen.Field[string]{
		Name: "TraceId",
		Bson: "commonReq.traceId",
	}
}
