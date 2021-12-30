package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"push-sdk-go/logger"
	"reflect"
	"sort"
	"strings"
	"time"
)

/****
   apiSign  参考微信签名规则
 @see https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_3
**/
func MakeApiSignature(sendParamObj interface{}, apiKeyByte []byte) string {
	apiKey := base64.StdEncoding.EncodeToString(apiKeyByte)
	str := getFieldString(sendParamObj)
	if str == "" {
		return ""
	}
	stringA := fmt.Sprintf("%s&%s=%s", str, "apiKey", apiKey)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(stringA))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))
	return strings.ToUpper(sign)
}

// 获取随机码
func GetNonceStr() (string, error) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b), nil
}

/**
  获取请求体对象字段值并按照ascii码排序，生成&key=v&key=v
*/
func getFieldString(sendParamEntity interface{}) string {
	m := reflect.TypeOf(sendParamEntity)
	v := reflect.ValueOf(sendParamEntity)
	if m.Kind() != reflect.Struct {
		logger.Warnw("m.Kind()", errors.New("Check type error not Struct"), "m.Kind()", m.Kind())
		return ""
	}
	var tagName string
	numField := m.NumField()
	w := make([]string, numField)
	numFieldCount := 0
	for i := 0; i < numField; i++ {
		fieldName := m.Field(i).Name
		tags := strings.Split(string(m.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		} else {
			tagName = m.Field(i).Name
		}
		if tagName == "xml" {
			continue
		}
		fieldValue := v.FieldByName(fieldName).Interface()

		if fieldValue != "" {
			if strings.Contains(tagName, "omitempty") {
				tagName = strings.Split(tagName, ",")[0]
			}
			s := fmt.Sprintf("%s=%v", tagName, fieldValue)
			w[numFieldCount] = s
			numFieldCount++
		}
	}
	if numFieldCount == 0 {
		return ""
	}
	w = w[:numFieldCount]
	sort.Strings(w)
	return strings.Join(w, "&")
}
