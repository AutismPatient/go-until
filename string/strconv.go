package string2

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

type IStringHelper interface {
	Decode([]byte) []byte
}

// Decode base64 解码
func Decode(raw []byte) []byte {
	var buf bytes.Buffer
	var decoded = make([]byte, 215)
	buf.Write(raw)
	decoder := base64.NewDecoder(base64.StdEncoding, &buf)
	decoder.Read(decoded)
	return decoded
}

// Encode base64 编码
func Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	encoder.Write(raw)
	encoder.Close()
	return encoded.Bytes()
}

// ToIntArray 将[]string --> []int64
func ToIntArray(s []string) []int64 {
	var arr []int64
	if len(s) > 0 {
		for _, v := range s {
			num, err := strconv.Atoi(v)
			if err != nil {
				return arr
			}
			arr = append(arr, int64(num))
		}
		return arr
	}
	return []int64{}
}

// 随机TOKEN
func RandToken(offset int64) string {
	var unix = time.Now().Unix()
	if offset != 0 {
		unix += offset
	}
	str := []byte(fmt.Sprintf("%s%d%s", "MZY", unix, "WMY"))
	newToken := md5.Sum(str)
	return fmt.Sprintf("%X", newToken)
}

// md5 sum
func MD5Sum(str string) string {
	sumStr := md5.Sum([]byte(str))
	return fmt.Sprintf("%X", sumStr)
}

// 判断字符串是否含有中文字符
func EqChineseChar(str string) (eq bool) {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// 判断非空
func EqEmpty(filter ...string) bool {
	for _, v := range filter {
		if v == "" {
			return false
		}
	}
	return true
}

// 随机生成指定位数的大写字母和数字的组合
func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+^%$#!"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成文件hash
func CreateFileHash(reader io.Reader) (str string) {
	newHash := sha256.New()
	_, err := io.Copy(newHash, reader)
	if err == nil {
		str = hex.EncodeToString(newHash.Sum(nil))
	}
	return
}
