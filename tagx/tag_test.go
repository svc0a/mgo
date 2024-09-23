package tagx

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	tagI1, err := GenerateV1("../")
	if err != nil {
		log.Fatal(err)
		return
	}
	//err = tagI.Generate()
	//if err != nil {
	//	logrus.Fatal(err)
	//	return
	//}
	log.Println(tagI1)
}
