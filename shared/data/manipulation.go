package data

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func Contains[T comparable](elems []T, v T, fn func(value T, element T) bool) bool {
	if fn == nil {
		fn = func(value T, element T) bool {
			return value == element
		}
	}
	for _, s := range elems {
		if fn(v, s) {
			return true
		}
	}
	return false
}

func ConvertToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}

	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "float64":
		return value.(float64)
	case "float32":
		return float64(value.(float32))
	case "int":
		return float64(value.(int))
	case "int32":
		return float64(value.(int32))
	case "int64":
		return float64(value.(int64))
	case "string":
		f, _ := strconv.ParseFloat(value.(string), 64)
		return f
	default:
		return 0
	}
}

func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ByteToMB(b int64) int64 {
	return b / 1024 / 1024
}

func MBToByte(mb int) int {
	return mb * 1024 * 1024
}

func HTMLSanitize(str string) string {
	// TODO: Add more sanitize
	re := regexp.MustCompile(`<|>|'|"`)
	return re.ReplaceAllString(str, "")
}

func ConvertToAnySlice[T any](input []T) []any {
	result := make([]any, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}

func ToUpperTurkish(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case 'ı':
			result.WriteRune('I')
		case 'i':
			result.WriteRune('İ')
		case 'ş':
			result.WriteRune('Ş')
		case 'ç':
			result.WriteRune('Ç')
		case 'ü':
			result.WriteRune('Ü')
		case 'ğ':
			result.WriteRune('Ğ')
		case 'ö':
			result.WriteRune('Ö')
		default:
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}
