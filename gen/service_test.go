package gen

import (
	"github.com/sirupsen/logrus"
	_ "github.com/svc0a/mgo/examples"
	"testing"
)

func TestGen(t *testing.T) {
	err := Define(WithDir("../"), WithPostgre()).Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
}
