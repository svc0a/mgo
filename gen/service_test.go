package gen

import (
	"github.com/sirupsen/logrus"
	_ "github.com/svc0a/mgo/examples"
	"github.com/svc0a/mgo/tagx/bsonx"
	"testing"
)

func TestGen(t *testing.T) {
	err := Define("../", bsonx.Define()).Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
}
