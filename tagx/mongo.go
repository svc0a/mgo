package tagx

import (
	"fmt"
	"github.com/svc0a/reflect2"
	"reflect"
	"strings"
)

type mongo struct {
	cache map[string]string
}

func Mongo() Client {
	b := &mongo{
		cache: map[string]string{},
	}
	return b
}

func (b *mongo) Register(in reflect.Type) Client {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in, "", "")
	return b
}

func (b *mongo) Register2(in reflect2.Type) Client {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in.Type1(), "", "")
	return b
}

func (b *mongo) register(objType reflect.Type, kPrefix, vPrefix string) {
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

func (b *mongo) getTagByStructField(structField reflect.StructField) string {
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

func (b *mongo) Export() map[string]string {
	for k, v := range b.cache {
		if v == "" {
			delete(b.cache, k)
		}
	}
	return b.cache
}

func (b *mongo) Format(value string) string {
	if value == "_id" {
		return value
	}
	return ToLowerCase(value)
}

func (b *mongo) GetTag(tag string) string {
	return getTagValue(tag, "bson")
}
