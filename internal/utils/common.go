package utils

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/zivwu/reminder-note-api/internal/config"
	"github.com/zivwu/reminder-note-api/internal/consts"
)

func ToJson(data any) string {
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("parse json fail:", err)
		return ""
	}
	return string(result)
}

func IsEmpty(values ...any) bool {
	for i := range values {
		v := values[i]
		if v == nil {
			return true
		}

		val := reflect.ValueOf(v)

		// interface 包 nil
		if val.Kind() == reflect.Interface || val.Kind() == reflect.Pointer {
			if val.IsNil() {
				return true
			}
		}

		switch val.Kind() {
		case reflect.String:
			return val.Len() == 0
		case reflect.Bool:
			return !val.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return val.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return val.Uint() == 0
		case reflect.Float32, reflect.Float64:
			return val.Float() == 0
		case reflect.Slice, reflect.Array, reflect.Map:
			return val.Len() == 0
		case reflect.Struct:
			// 空 struct: 每個欄位都空才算空
			for i := 0; i < val.NumField(); i++ {
				if !IsEmpty(val.Field(i).Interface()) {
					return false
				}
			}
			return true
		}
	}
	return false
}

// 泛用 struct 轉換 function
func StructConvert[T any](src any) (*T, error) {
	var dst T
	err := copier.Copy(&dst, src)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func GenToken(claims jwt.MapClaims) (string, error) {
	// 產生token
	claims["exp"] = time.Now().Add(consts.TokenExpireTime).Unix() // 過期時間
	claims["iat"] = time.Now().Unix()                             // 發行時間

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽名
	signedToken, err := token.SignedString([]byte(config.Env.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
