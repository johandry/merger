package merger

import (
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
