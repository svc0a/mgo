package gen

import (
	"fmt"
	"github.com/sirupsen/logrus"
	_ "github.com/svc0a/mgo/examples"
	"testing"
)

func TestGen(t *testing.T) {
	tagI1 := Define("../")
	err := tagI1.Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Sprintf("%v", tagI1)
}
