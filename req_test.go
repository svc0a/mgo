package main

import (
	"github.com/svc0a/mgo/gen"
	"github.com/svc0a/mgo/gen/dao/walletAdjustmentReq"
	"testing"
)

func TestData(t *testing.T) {
	t.Log(walletAdjustmentReq.New().CommonReq().TraceId().Bson)
	t.Log(walletAdjustmentReq.New().TransactionId().Bson)
	walletAdjustmentReq.New().TransactionId().Eq("test")
	filter := gen.Query().
		Where(walletAdjustmentReq.New().TransactionId().Eq("test")).
		Where(walletAdjustmentReq.New().CommonReq().TraceId().Eq("test2")).
		Build()
	t.Log(filter)
}
