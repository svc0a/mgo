package tagx

import (
	"github.com/svc0a/mgo/examples"
	"github.com/svc0a/mgo/examples/types"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	//tagI1 := Define("../").WithModel(examples.Order{}).WithModel(examples.User{}).WithModel(types.Entity1{})
	tagI1 := Define("../").WithModel(examples.Order{}).WithModel(types.Entity1{}).WithModel(examples.User{})
	//err = tagI.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	log.Println(tagI1)
}
