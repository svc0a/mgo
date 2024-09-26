package gen

import (
	"github.com/sirupsen/logrus"
	_ "github.com/svc0a/mgo/examples"
	"github.com/svc0a/mgo/tagx/pgx"
	"testing"
)

func TestGen(t *testing.T) {
	err := Define("../", pgx.Define()).Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
}
