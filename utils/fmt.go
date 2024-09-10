package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ParseInt64(val any) int64 {
	switch v := val.(type) {
	case string:
		r, _ := strconv.ParseInt(v, 10, 64)
		return r
	case int64:
		return v
	}
	return 0
}
func ParseInt8(val any) int8 {
	switch v := val.(type) {
	case string:
		r, _ := strconv.ParseInt(v, 10, 8)
		return int8(r)
	case int:
		return int8(v)
	case int64:
		return int8(v)
	case int8:
		return v
	}
	return 0
}
func ParseFloat64(val any) float64 {
	switch v := val.(type) {
	case string:
		r, _ := strconv.ParseFloat(v, 64)
		return r
	case float64:
		return v
	}
	return 0
}

// DecodeJSON Define
func DecodeJSON(jStr string, target any) error {
	d := json.NewDecoder(strings.NewReader(jStr))
	d.UseNumber()
	err := d.Decode(target)
	if err != nil {
		return err
	}
	return nil
}

// EncodeJSON Define
func EncodeJSON(val any) string {
	if val == nil {
		return ""
	}
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func StructToMap(obj interface{}) (map[string]any, error) {
	result := make(map[string]any)

	// 获取obj的反射对象
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		// 如果obj是个指针，则需要解引用
		objValue = objValue.Elem()
	}

	// 确保objValue是结构体类型
	if objValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got: %s", objValue.Kind())
	}

	// 遍历结构体的所有字段
	for i := 0; i < objValue.NumField(); i++ {
		fieldName := objValue.Type().Field(i).Name
		fieldValue := objValue.Field(i).Interface()

		// 通过json标签获取真正的字段名
		jsonTag := objValue.Type().Field(i).Tag.Get("json")
		if jsonTag != "" {
			fieldName = jsonTag
		}

		result[fieldName] = fieldValue
	}

	return result, nil
}

func IpToInt(ip string) int64 {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0
	}

	var sum int64 = 0
	for _, part := range parts {
		sum *= 255
		num, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return 0
		}
		sum += num
	}

	return sum
}

func CIDRToString(s string) (bool, string) {
	ip, _, err := net.ParseCIDR(s)
	if err != nil {
		return false, s
	}
	if ip.To4() == nil {
		return false, s
	}
	return true, ip.String() // 去除CIDR部分
}

func GetTimeNow() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	NewNow := now.In(loc).Format("2006-01-02 15:04:05")
	NewDt, _ := time.ParseInLocation("2006-01-02 15:04:05", NewNow, loc)
	return NewDt // 去除CIDR部分
}
func TrimString(raw string) string {
	return strings.Trim(strings.TrimSpace(strings.Trim(strings.TrimSpace(raw), `"`)), `"`)
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() != nil
}

func Intersect[T any](a, b []T) []T {
	m := make(map[any]bool)
	n := make([]T, 0)

	// 将切片a的元素设置为map的key
	for _, v := range a {
		m[v] = true
	}

	// 遍历切片b，如果元素在map中存在，则加入结果切片
	for _, v := range b {
		if m[v] {
			n = append(n, v)
		}
	}

	return n
}

func PrintPid() {
	pid := os.Getpid()
	fmt.Printf("当前进程的PID: %d\n", pid)
}

func PrintGid() {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	fmt.Printf("当前goroutine的ID: %d\n", n)
}
