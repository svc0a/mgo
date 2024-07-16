package bsontag

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"unicode"
)

var cache = make(map[any]string)

func Get(field any) string {
	if reflect.TypeOf(field).Kind() != reflect.Ptr {
		logrus.Fatal("field must be a pointer to a struct")
	}
	s, ok := cache[field]
	if !ok {
		logrus.Fatal("field not registered")
	}
	return s
}

func Register[T any]() *T {
	var obj T
	objVal := reflect.ValueOf(&obj)
	if objVal.Kind() != reflect.Ptr {
		logrus.Fatal("obj must be a pointer to a struct")
	}
	objVal = objVal.Elem()
	register(objVal)
	return &obj
}

func register(objVal reflect.Value, prefix ...string) {
	objType := objVal.Type()
	for i := 0; i < objVal.NumField(); i++ {
		prefix1 := ""
		if len(prefix) != 0 && prefix[0] != "" {
			prefix1 = fmt.Sprintf("%s.", prefix[0])
		}
		cache[objVal.Field(i).Addr().Interface()] = fmt.Sprintf("%s%s", prefix1, getTagByStructField(objType.Field(i)))
		if objVal.Field(i).Type().Kind() == reflect.Struct {
			register(objVal.Field(i), cache[objVal.Field(i).Addr().Interface()])
		}
	}
}

func getTagByStructField(structField reflect.StructField) string {
	bsonTag := structField.Tag.Get("bson")
	if bsonTag == "" {
		return lowerCamelCase(structField.Name)
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

// lowerCamelCase 将字符串转换为小写开头的驼峰形式
func lowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
