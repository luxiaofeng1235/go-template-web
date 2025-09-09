package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gogf/gf/v2/net/ghttp"
	uuid "github.com/satori/go.uuid"
)

// GetRandomString 生成指定长度的随机字符串
func GetRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为种子
		return fmt.Sprintf("%d", time.Now().UnixNano())[:length]
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// GetNowUnix 获取当前时间戳
func GetNowUnix() int64 {
	return time.Now().Unix()
}

// StrToUnix 字符串时间转Unix时间戳
func StrToUnix(str string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// StrToInt 字符串转int
func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// StrToInt64 字符串转int64
func StrToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// FormatMobileStar 手机号中间4位替换为*号
func FormatMobileStar(mobile string) string {
	if len(mobile) <= 10 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}

// GetRequestIP 获取请求的IP地址（GoFrame版本）
func GetRequestIP(r *ghttp.Request) string {
	clientIP := r.GetClientIp()
	if clientIP == "::1" {
		clientIP = "127.0.0.1"
	}
	return clientIP
}

// GetLocalIP 获取本机IP地址
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		// 检查ip地址判断是否回环地址
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}

		if ipAddr.IP.To4() != nil {
			return ipAddr.IP.To4().String(), nil
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// FirstUpper 首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	regRuler := "^1([356789][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\\d{8}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}

// CheckEmail 验证邮箱
func CheckEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// StructToMap JSON序列化方式将结构体转map
func StructToMap(stuObj interface{}) (map[string]interface{}, error) {
	strRet, err := json.Marshal(stuObj)
	if err != nil {
		return nil, err
	}
	var mRet map[string]interface{}
	err1 := json.Unmarshal(strRet, &mRet)
	if err1 != nil {
		return nil, err1
	}
	return mRet, nil
}

// SubStr 截取字符串，支持多字节字符
func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

// Union 求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 求交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// Difference 求差集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

// GetUUID 获取uuid
func GetUUID() uuid.UUID {
	return uuid.NewV4()
}

// GetRandomUsername 获取随机用户名
func GetRandomUsername() string {
	uuidStr := GetUUID().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	shortUsername := uuidStr[:9]
	return fmt.Sprintf("%d%s", time.Now().Unix()%10, shortUsername)
}

// FormatDecimal 格式化小数
func FormatDecimal(amount string) string {
	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return amount
	}

	// 判断是否有小数部分
	hasDecimal := value != float64(int(value))

	if hasDecimal {
		// 保留3位小数
		formatted := fmt.Sprintf("%.3f", value)
		return formatted
	}

	// 没有小数部分，返回整数形式
	return strconv.Itoa(int(value))
}

// FormatInt64 字符串转int64
func FormatInt64(str string) (num int64) {
	var err error
	num, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	return
}

// GetWanFormatted 格式化数字为万、亿单位
func GetWanFormatted(hits int64) string {
	hitsFloat := float64(hits)
	absHits := math.Abs(hitsFloat)

	switch {
	case absHits >= 1e8:
		tempStr := fmt.Sprintf("%.2f", hitsFloat/1e8)
		return fmt.Sprintf("%v 亿", trimZero(tempStr))
	case absHits >= 1e4:
		tempStr := fmt.Sprintf("%.2f", hitsFloat/1e4)
		return fmt.Sprintf("%v 万", trimZero(tempStr))
	default:
		return fmt.Sprintf("%d", hits)
	}
}

// trimZero 去除小数点后的0
func trimZero(tempStr string) string {
	tempStr = strings.TrimRight(tempStr, "0")
	tempStr = strings.TrimRight(tempStr, ".")
	return tempStr
}

// GetWords 获取字符串字符数
func GetWords(text string) (words int) {
	words = utf8.RuneCountInString(text)
	return
}

// CapitalizeFirstChar 把第一个字母转换大写
func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// GetRequestHeaderByName 根据特定的name获取对应的header信息
func GetRequestHeaderByName(r *ghttp.Request, name string) string {
	return r.Header.Get(name)
}

// SlicePage 切片分页
func SlicePage(page, pageSize int, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		if page > 1 {
			return
		}
		return sliceStart, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize
	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

// IsValidIP 验证IP地址是否有效
func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// JSONString 将对象转换为JSON字符串
func JSONString(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
