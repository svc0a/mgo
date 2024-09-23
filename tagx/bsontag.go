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
	b.register(in.Type1(), "", "")
	return b
}

func (b *bsonTag) register(objType reflect.Type, kPrefix, vPrefix string) {
	for i := 0; i < objType.NumField(); i++ {
		k := objType.Field(i).Name
		v := b.getTagByStructField(objType.Field(i))
		if kPrefix != "" && vPrefix != "" {
			k = fmt.Sprintf("%s_%s", kPrefix, objType.Field(i).Name)
			v = fmt.Sprintf("%s.%s", vPrefix, b.getTagByStructField(objType.Field(i)))
			v = strings.ReplaceAll(v, "..", ".")
		}
		if objType.Field(i).Type.Kind() == reflect.Struct {
			b.register(objType.Field(i).Type, k, v)
		} else {
			b.cache[k] = v
		}
	}
}

func (b *bsonTag) getTagByStructField(structField reflect.StructField) string {
	bsonTag1 := structField.Tag.Get("bson")
	if bsonTag1 == "" {
		return structField.Name
	}
	if bsonTag1 == "-" {
		return ""
	}
	if strings.Contains(bsonTag1, ",inline") {
		// 处理 ",inline" 标签的情况
		return ""
	}
	return strings.Split(bsonTag1, ",")[0]
}

func (b *bsonTag) Export() map[string]string {
	for k, v := range b.cache {
		if v == "" {
			delete(b.cache, k)
		}
	}
	return b.cache
}
