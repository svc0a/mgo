package bsonx

import (
	"fmt"
	"github.com/svc0a/mgo/tagx"
	"github.com/svc0a/reflect2"
	"reflect"
	"strings"
)

type impl struct {
	cache map[string]string
}

func Define() tagx.Service {
	b := &impl{
		cache: map[string]string{},
	}
	return b
}

func (b *impl) Register(in reflect.Type) tagx.Service {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in, "", "")
	return b
}

func (b *impl) Register2(in reflect2.Type) tagx.Service {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in.Type1(), "", "")
	return b
}

func (b *impl) register(objType reflect.Type, kPrefix, vPrefix string) {
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

func (b *impl) getTagByStructField(structField reflect.StructField) string {
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

func (b *impl) Export() map[string]string {
	for k, v := range b.cache {
		if v == "" {
			delete(b.cache, k)
		}
	}
	return b.cache
}
