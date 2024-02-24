package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var t = rand.New(rand.NewSource(time.Now().Unix()))

// 随机字符串
func RandString() string {
	upLength := 12

	chars := "QWERTYUOPASDFGHJKZXCVBNMqwertyuiopasdfghjkzxcvbnm123456789"
	randStr := ""
	for i := 0; i < upLength; i++ {
		randStr += string(chars[t.Intn(len(chars))])
	}
	return randStr
}

// 随机字符串
func RandStringIdent() string {
	upLength := 20

	chars := "QWERTYUOPASDFGHJKZXCVBNMqwertyuiopasdfghjkzxcvbnm123456789"
	randStr := ""
	for i := 0; i < upLength; i++ {
		randStr += string(chars[t.Intn(len(chars))])
	}
	return randStr
}

// 输出错误消息
func EchoError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// 输出错误消息后退出程序
func EchoExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 获取证书字符串
func GetCertString(certPath string) (certString string) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(data)
}

// 输出条件消息，返回一个字符串
func InputMsg(outMsg string) string {
	fmt.Println(outMsg)
	inputMsg := ""
	fmt.Scanln(&inputMsg)
	return inputMsg
}
