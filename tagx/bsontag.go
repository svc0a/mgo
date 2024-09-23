package tagx

import (
	"fmt"
	"reflect"
	"strings"
)

var cache = map[string]map[string]string{} // map[structName]map[fieldName]string{}

type BsonTag interface {
	Export() map[string]string
}

type bsonTag struct {
	cache map[string]string
}

func DefineType[T any]() BsonTag {
	t := reflect.TypeFor[T]()
	b := &bsonTag{
		cache: map[string]string{},
	}
	b.register(t)
	cache[fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())] = b.Export()
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
