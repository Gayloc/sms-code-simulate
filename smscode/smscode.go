package smscode

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type Information struct {
	SMSNumber string
	SMSCode   [16]byte
	SendTimes int
}

func SendSMSCode() ([16]byte, time.Time) {
	SMSCode := ""
	temp := "abcdefghijklmnopqrstuvwxyz1234567890" //验证码字符
	for i := 0; i < 6; i++ {                       //循环次数决定验证码长度
		SMSCode += string([]byte(temp)[rand.Intn(len(temp))])
	}
	fmt.Print("验证码是:" + SMSCode + "  有效时间5分钟\n")
	return md5.Sum([]byte(SMSCode)), time.Now()
}

func GetSMSNumber() string {
	var SMSNumber string
	for {
		fmt.Print("输入电话号码:")
		fmt.Scan(&SMSNumber)
		re := regexp.MustCompile("[^0-9]")
		flag := re.MatchString(SMSNumber) //检查输入格式是否正确
		if !flag && len(SMSNumber) == 11 {
			break
		} else {
			fmt.Print("输入不正确\n")
		}
	}
	return SMSNumber
}
