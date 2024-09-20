package tagx

import (
	"fmt"
	"github.com/svc0a/reflect2"
	"reflect"
	"strings"
)

type BsonTag interface {
	Export() map[string]string
}

type bsonTag struct {
	cache map[string]string
}

func defineByType(in reflect2.Type) BsonTag {
	b := &bsonTag{
		cache: map[string]string{},
	}
	b.register(in.Type1())
	return b
}

func (b *bsonTag) register(objType reflect.Type, prefix ...string) {
	for i := 0; i < objType.NumField(); i++ {
		namePrefix1 := ""
		valPrefix1 := ""
		if len(prefix) != 0 && prefix[0] != "" {
			namePrefix1 = fmt.Sprintf("%s_", prefix[0])
			valPrefix1 = fmt.Sprintf("%s.", prefix[0])
		}
		k := fmt.Sprintf("%s%s", namePrefix1, objType.Field(i).Name)
		v := fmt.Sprintf("%s%s", valPrefix1, b.getTagByStructField(objType.Field(i)))
		b.cache[k] = v
		if objType.Field(i).Type.Kind() == reflect.Struct {
			b.register(objType.Field(i).Type, v)
		}
	}
}

func (b *bsonTag) getTagByStructField(structField reflect.StructField) string {
	bsonTag := structField.Tag.Get("bson")
	if bsonTag == "" {
		return structField.Name
	}
	if bsonTag == "-" {
		return ""
	}
	if strings.Contains(bsonTag, ",inline") {
		// 处理 ",inline" 标签的情况
		return ""
	}
	return strings.Split(bsonTag, ",")[0]
}

func (b *bsonTag) Export() map[string]string {
	return b.cache
}
