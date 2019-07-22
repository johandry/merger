package merger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const defaultTagName = "json"

// TransformMap transform a map of string values to interface{} values
func TransformMap(srcMap map[string]string) map[string]interface{} {
	m := make(map[string]interface{}, 0)
	for k, v := range srcMap {
		var i interface{}
		switch {
		case isSlice(v):
			i = transformToSlice(v)
		case isJSONStruct(v):
			i = transformJSONToStruct(v)
		default:
			i = v
		}

		if isStructField(k) {
			m = transformToStructField(m, k, i)
		} else {
			if m[k] != nil && isMap(i) {
				m[k] = mergeTwoMaps(m[k].(map[string]interface{}), i.(map[string]interface{}), false)
			} else {
				m[k] = i
			}
		}
	}

	return m
}

func isSlice(v string) bool {
	return (strings.Contains(v, ",") || (strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]"))) && !isJSONStruct(v)
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
	if isStructField(k) {
		keys := strings.Split(k, FieldSeparator)
		k0 := keys[0]
		r := strings.Join(keys[1:], FieldSeparator)

		if _, ok := m[k0]; !ok {
			m[k0] = make(map[string]interface{}, 0)
		}

		m[k0] = transformToStructField(m[k0].(map[string]interface{}), r, v)
		return m
	}

	if _, ok := m[k]; !ok {
		m[k] = v
		return m
	}

	if m[k] != nil && isMap(v) {
		m[k] = mergeTwoMaps(m[k].(map[string]interface{}), v.(map[string]interface{}), false)
		return m
	}

	m[k] = v
	return m
}

func isJSONStruct(v string) bool {
	return strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")
}
func transformJSONToStruct(v string) map[string]interface{} {
	var m map[string]interface{}
	// Ignore error, if it's not a valid JSON return an empty map
	json.Unmarshal([]byte(v), &m)
	return m
}

func isMap(v interface{}) bool {
	t := reflect.TypeOf(v).String()
	return t == "map[string]interface {}" || t == "map[string]string"
}

func mergeTwoMaps(dstMap, srcMap map[string]interface{}, overwrite bool) map[string]interface{} {
	for k := range srcMap {
		// dstMap[k] doesn't have a value
		if _, ok := dstMap[k]; !ok {
			dstMap[k] = srcMap[k]
			continue
		}
		// dstMap[k] and srcMap[k] are both maps, merge them
		if isMap(dstMap[k]) && isMap(srcMap[k]) {
			dstMap[k] = mergeTwoMaps(dstMap[k].(map[string]interface{}), srcMap[k].(map[string]interface{}), overwrite)
			continue
		}
		// dstMap[k] or srcMap[k] is a maps and the other is not,
		// or both are not maps, force the assigment if overwrite
		if overwrite {
			dstMap[k] = srcMap[k]
		}
	}

	return dstMap
}

// TransformToMap returns the given interface (has to be a struct) as a set of
// variables + values map. Useful to get the environment variables or parameters
// of a given struct before merge it with other struct
func TransformToMap(v interface{}, tagNames ...string) (map[string]string, error) {
	m := map[string]string{}

	if v == nil {
		return m, nil
	}

	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return m, fmt.Errorf("invalid value, it's not a struct pointer, it's a %s. %v", ptrRef.Kind().String(), ptrRef)
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return m, fmt.Errorf("invalid value, it's not a struct, it's a %s. %v", ref.Kind().String(), ref)
	}

	if len(tagNames) == 0 {
		tagNames = []string{defaultTagName}
	}

	return parseStruct("", ref, m, tagNames), nil
}

func parseStruct(parent string, val reflect.Value, m map[string]string, tagNames []string) map[string]string {
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		refTypeField := valType.Field(i)
		name, ignore := getName(refTypeField, tagNames)
		if ignore {
			continue
		}

		if len(parent) != 0 {
			name = parent + FieldSeparator + name
		}

		// Because all the names are lower case. Case does not matter
		name = strings.ToLower(name)
		valField := val.Field(i)
		m = appendTo(name, m, valField, tagNames)
	}

	return m
}

func getName(field reflect.StructField, tagNames []string) (name string, ignore bool) {
	// Try the given keys, in order of importance ...
	for _, tagName := range tagNames {
		var ok bool
		if name, ok = field.Tag.Lookup(tagName); ok && len(name) != 0 && name != "-" {
			// if this tag is not found or contain an empty or `-` value, try the next tag
			// otherwise, end the loop bc the name have been found
			return name, false
		}
		if name == "-" {
			// If the tag name is `-` it will be ignored
			return "", true
		}
	}

	// If any tag is found, or they are empty and are not `-`, return the field name
	return field.Name, false
}

func appendTo(name string, m map[string]string, v reflect.Value, tagNames []string) map[string]string {
	if !v.CanInterface() {
		return m
	}

	val := reflect.ValueOf(v.Interface())
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		m = parseStruct(name, val, m, tagNames)
	case reflect.Map:
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			keyStr = strings.Replace(keyStr, " ", "_", -1)
			n := name + FieldSeparator + keyStr
			m = appendTo(n, m, val.MapIndex(key), tagNames)
		}
	case reflect.Slice, reflect.Array:
		switch val.Type().Elem().Kind() {
		case reflect.Struct, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Array:
			// TODO: At this time only a slice/array of simple type are possible to map
			// return m, fmt.Errorf("cannot map a slice/array of struct, map or slice/array")
		default:
			list := "["
			if len := val.Len(); len > 0 {
				var i int
				for i = 0; i < len-1; i++ {
					list = list + fmt.Sprintf("%v", val.Index(i).Interface()) + ", "
				}
				list = list + fmt.Sprintf("%v", val.Index(i).Interface())
			}
			list = list + "]"
			m[name] = list
		}
	default:
		value := val.Interface()
		m[name] = fmt.Sprintf("%v", value)
	}

	return m
}
