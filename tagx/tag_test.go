package tagx

import (
	"github.com/svc0a/mgo/examples"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	tagI1 := Define("../").WithModel(examples.Order{}).WithModel(examples.User{})
	//tagI1 := Define("../").WithModel(examples.User{}).WithModel(examples.Order{})
	//err = tagI.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	log.Println(tagI1)
}
