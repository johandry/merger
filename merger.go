package merger

import (
	"log"
	"strings"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

// FieldSeparator separates the fields of a struct when defining paramete names
const FieldSeparator = "__"

// Merge merges the given map and optional structs into the dst structure
func Merge(dst interface{}, srcMap map[string]string, srcs ...interface{}) error {
	if err := mergeMap(dst, srcMap); err != nil {
		return err
	}

	return MergeStruct(dst, srcs...)
}

// MergeMap merges the given maps into the dst structure
func MergeMap(dst interface{}, srcMaps ...map[string]string) error {
	for i := range srcMaps {
		srcMap := srcMaps[len(srcMaps)-i-1]

		if err := mergeMap(dst, srcMap); err != nil {
			return err
		}
	}

	return nil
}

func mergeMap(dst interface{}, srcMap map[string]string) error {
	m := TransformMap(srcMap)

	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &dst,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(m); err != nil {
		return err
	}

	return nil
}

// MergeStruct merges the given structs into the dst structure
func MergeStruct(dst interface{}, srcs ...interface{}) error {
	for _, src := range srcs {
		if err := mergo.Merge(dst, src); err != nil {
			return err
		}
	}
	return nil
}

// TransformMap transform a map of string values to interface{} values
func TransformMap(srcMap map[string]string) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	for k, v := range srcMap {
		var i interface{}
		if isSlice(v) {
			i = transformToSlice(v)
		} else {
			i = v
		}
		if isStructField(k) {
			m = transformToStructField(m, k, i)
		} else if isJSONStruct(v) {
			m = transformToStruct(m, k, i)
		} else {
			m[k] = i
		}
	}

	return m
}

func isSlice(v string) bool {
	return strings.Contains(v, ",")
}
func transformToSlice(v string) []string {
	v = strings.Trim(v, "[ ]")
	strs := strings.Split(v, ",")
	values := make([]string, 0)
	for _, str := range strs {
		str = strings.Trim(str, " ")
		str = strings.Trim(str, "'")
		values = append(values, str)
	}
	return values
}

func isStructField(k string) bool {
	return strings.Contains(k, FieldSeparator)
}
func transformToStructField(m map[string]interface{}, k string, v interface{}) map[string]interface{} {
	if strings.Contains(k, FieldSeparator) {
		keys := strings.Split(k, FieldSeparator)
		k0 := keys[0]
		r := strings.Join(keys[1:], FieldSeparator)
		log.Println(k0, r)

		if _, ok := m[k0]; !ok {
			m[k0] = make(map[string]interface{}, 0)
		}

		m[k0] = transformToStruct(m[k0].(map[string]interface{}), r, v)
		return m
	}
	m[k] = v
	return m
}

func isJSONStruct(v string) bool {
	return strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")
}
func transformToStruct(m map[string]interface{}, k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
