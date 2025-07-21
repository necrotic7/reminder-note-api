package utils

import (
	"encoding/json"
	"log"
	"reflect"
)

func ToJson(data any) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("parse json fail:", err)
		return ""
	}
	return string(result)
}

func IsNill(d ...any) bool {
	for i := range d {
		v := reflect.ValueOf(d[i])
		if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
			return true
		}
	}
	return false
}
