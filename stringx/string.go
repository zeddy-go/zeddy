package stringx

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func IsMobile(str string) bool {
	reg := `^1[3456789]\d{9}$`
	r, _ := regexp.Compile(reg)
	return r.MatchString(str)
}

func IsTrue(str string) bool {
	if str != "" && str != "0" && str != "false" {
		return true
	}
	return false
}

func RandString(length int) string {
	s := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = s[r.Intn(len(s))]
	}
	return string(bytes)
}

func RandNumeric(len int) string {
	numeric := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	r := rand.New(rand.NewSource(time.Now().Unix()))

	var sb strings.Builder

	for i := 0; i < len; i++ {
		sb.WriteByte(numeric[r.Intn(10)])
	}

	return sb.String()
}

func Contains(target any, str string) bool {
	switch x := target.(type) {
	case string:
		return strings.Contains(str, x)
	case []string:
		for _, s := range x {
			if strings.Contains(str, s) {
				return true
			}
		}
	}
	return false
}
