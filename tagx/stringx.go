package tagx

import (
	"strings"
	"unicode"
)

// 将字符串转为小写驼峰
func ToLowerCase(s string) string {
	var words []string
	start := 0
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			words = append(words, s[start:i])
			start = i
		}
	}
	// 加入最后一个词
	words = append(words, s[start:])

	// 转为小写驼峰
	for i := range words {
		words[i] = strings.ToLower(words[i])
		if i > 0 {
			words[i] = capitalize(words[i])
		}
	}
	resp := strings.Join(words, "")
	if resp == "" {
		return ""
	}
	runes := []rune(resp)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// 将字符串的首字母大写
func capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	s := string(runes)
	return s
}

// 将字符串转为大写驼峰
func ToUpperCamel(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || unicode.IsSpace(r)
	})
	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, "")
}

// 将字符串转为下划线分割
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// 将字符串转为中划线分割
func ToKebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '-')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// getTagValue 解析标签字符串，提取指定键的值（如 json:"value"）
func getTagValue(tag, key string) string {
	// 去除反引号，标签形如 `json:"name" db:"value"`
	tag = strings.Trim(tag, "`")

	// 使用空格分隔不同的键值对
	parts := strings.Split(tag, " ")

	for _, part := range parts {
		// 查找以指定键开头的部分，如 json:"name"
		if strings.HasPrefix(part, key+":") {
			// 取出键的值，去除引号
			tagValue := strings.TrimPrefix(part, key+":")
			return strings.Trim(tagValue, `"`)
		}
	}
	return ""
}
