package tagx

import (
	"github.com/sirupsen/logrus"
	"github.com/svc0a/mgo/examples"
	"github.com/svc0a/mgo/examples/types"
	"log"
	"testing"
)

func TestDefine(t *testing.T) {
	//tagI1 := Define("../").WithModel(examples.Order{}).WithModel(examples.User{}).WithModel(types.Entity1{})
	//tagI1 := Define("../").WithModel(examples.Order{}).WithModel(types.Entity1{}).WithModel(examples.User{})
	tagI1 := Define("../").WithModel(examples.Order{}).WithModel(types.Entity1{}).WithModel(examples.User{})
	//err = tagI.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	tagI1.Generate()
	log.Println(tagI1)
}

func TestGenerate(t *testing.T) {
	//tagI1 := Define("../").WithModel(examples.Order{}).WithModel(examples.User{}).WithModel(types.Entity1{})
	//tagI1 := Define("../").WithModel(examples.Order{}).WithModel(types.Entity1{}).WithModel(examples.User{})
	//err = tagI.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	//err := tagI1.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	logrus.Info("generate ok")
}
