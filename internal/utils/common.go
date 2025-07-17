package utils

import (
	"encoding/json"
	"log"
)

func ToJson(data any) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("parse json fail:", err)
		return ""
	}
	return string(result)
}
