package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MapToStruct(input map[string]string, output interface{}) error {
	outputValue := reflect.ValueOf(output)

	// Output bir pointer olmalı
	if outputValue.Kind() != reflect.Ptr || outputValue.IsNil() {
		return fmt.Errorf("output must be a non-nil pointer")
	}

	// Struct'ı al
	structValue := outputValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct")
	}

	structType := structValue.Type()

	// Her field için
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Özel field ismi tag'den oku veya field adını kullan
		fieldName := fieldType.Tag.Get("redis")
		if fieldName == "" {
			fieldName = fieldType.Name
		}

		// Map'te bu isimle değer var mı?
		val, ok := input[fieldName]
		if !ok {
			continue // Bu field için map'te değer yok, geç
		}

		// Değeri field tipine uygun şekilde dönüştür
		switch field.Kind() {
		case reflect.String:
			field.SetString(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fieldType.Type == reflect.TypeOf(time.Duration(0)) {
				// Duration için özel dönüşüm
				if dur, err := time.ParseDuration(val); err == nil {
					field.SetInt(int64(dur))
				}
			} else {
				// Normal integer
				if i, err := strconv.ParseInt(val, 10, 64); err == nil {
					field.SetInt(i)
				}
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if i, err := strconv.ParseUint(val, 10, 64); err == nil {
				field.SetUint(i)
			}

		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				field.SetFloat(f)
			}

		case reflect.Bool:
			if b, err := strconv.ParseBool(val); err == nil {
				field.SetBool(b)
			}
		}
	}

	return nil
}

func StructToMap(src interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(src)
	if !v.IsValid() {
		return nil, errors.New("nil value")
	}
	// pointer ise indir
	v = reflect.Indirect(v)
	t := v.Type()

	// map[string]any / map[string]string direkt
	if v.Kind() == reflect.Map {
		if t.Key().Kind() != reflect.String {
			return nil, errors.New("map key must be string")
		}
		out := make(map[string]any, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			out[iter.Key().String()] = iter.Value().Interface()
		}
		return out, nil
	}

	// struct
	if v.Kind() != reflect.Struct {
		return nil, errors.New("expected struct or map")
	}

	out := make(map[string]any, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		// export edilmeyen alanları atla
		if sf.PkgPath != "" {
			continue
		}
		name, omitEmpty := fieldName(sf)
		if name == "-" || name == "" {
			continue
		}
		fv := v.Field(i)
		if omitEmpty && isZero(fv) {
			continue
		}
		out[name] = fv.Interface()
	}
	return out, nil
}

// Tag önceliği: redis > json > bson
func fieldName(sf reflect.StructField) (name string, omitEmpty bool) {
	for _, tagKey := range []string{"redis", "json", "bson"} {
		if tag, ok := sf.Tag.Lookup(tagKey); ok {
			parts := strings.Split(tag, ",")
			n := parts[0]
			if n != "" {
				if len(parts) > 1 {
					for _, opt := range parts[1:] {
						if opt == "omitempty" {
							return n, true
						}
					}
				}
				return n, false
			}
			// explicit boş ise pas geç, bir sonraki tag'e bak
		}
	}
	return sf.Name, false
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	case reflect.Struct:
		// time.Time gibi tipler için:
		z := reflect.Zero(v.Type()).Interface()
		return reflect.DeepEqual(v.Interface(), z)
	default:
		return false
	}
}
