package tagx

import (
	"fmt"
	"github.com/svc0a/reflect2"
	"reflect"
	"strings"
)

type postgre struct {
	cache map[string]string
}

func Postgre() Client {
	b := &postgre{
		cache: map[string]string{},
	}
	return b
}

func (b *postgre) Register(in reflect.Type) Client {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in, "", "")
	return b
}

func (b *postgre) Register2(in reflect2.Type) Client {
	if b.cache == nil {
		b.cache = map[string]string{}
	}
	b.register(in.Type1(), "", "")
	return b
}

func (b *postgre) register(objType reflect.Type, kPrefix, vPrefix string) {
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		k := field.Name
		v := b.getTagByStructField(field)

		// 忽略字段处理，如果返回值为空，直接跳过
		if v == "" {
			continue
		}

		// 处理 gorm:"embedded" 或 gorm:"embeddedPrefix:xxx"
		if embedded, embeddedPrefix := b.isEmbedded(field); embedded {
			// 如果有 embeddedPrefix，设置前缀
			if embeddedPrefix != "" {
				kPrefix = fmt.Sprintf("%s%s", kPrefix, embeddedPrefix)
			} else {
				kPrefix = k
			}
			// 递归处理嵌套结构体
			b.register(field.Type, kPrefix, vPrefix)
			continue
		}

		// 如果有前缀，处理前缀的拼接
		if kPrefix != "" && vPrefix != "" {
			k = fmt.Sprintf("%s_%s", kPrefix, field.Name)
			v = fmt.Sprintf("%s%s", vPrefix, v)
		}

		// 如果是结构体，递归处理嵌套字段
		if field.Type.Kind() == reflect.Struct {
			b.register(field.Type, k, v)
		} else {
			b.cache[k] = v
		}
	}
}

func (b *postgre) getTagByStructField(structField reflect.StructField) string {
	// 获取 gorm 标签
	tag1 := structField.Tag.Get("gorm")

	// 忽略字段标签处理
	if tag1 == "-" {
		return ""
	}

	// 如果标签为空，则返回字段名称
	if tag1 == "" {
		return structField.Name
	}

	// 解析 gorm 标签并处理 column:<column_name>
	tagParts := strings.Split(tag1, ";")
	for _, part := range tagParts {
		// 处理 `column:<column_name>` 标签
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}

	// 如果没有 column 标签，返回字段名称
	return structField.Name
}

// 判断字段是否为 embedded 或者带 embeddedPrefix 标签
func (b *postgre) isEmbedded(structField reflect.StructField) (bool, string) {
	tag1 := structField.Tag.Get("gorm")

	// 检查是否包含 embedded 或 embeddedPrefix 标签
	if strings.Contains(tag1, "embedded") {
		// 检查 embeddedPrefix 标签
		tagParts := strings.Split(tag1, ";")
		for _, part := range tagParts {
			if strings.HasPrefix(part, "embeddedPrefix:") {
				return true, strings.TrimPrefix(part, "embeddedPrefix:")
			}
		}
		return true, ""
	}

	return false, ""
}

func (b *postgre) Export() map[string]string {
	for k, v := range b.cache {
		if v == "" {
			delete(b.cache, k)
		}
	}
	return b.cache
}
