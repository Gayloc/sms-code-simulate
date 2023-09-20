package controller

import (
	"crypto/md5"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"../smscode"
)

func Start(user smscode.Information) {
	var selection int
	var sendtime, runtime time.Time
	runtime = time.Now()
	sendtime = runtime

	re := regexp.MustCompile("[^0-9]")
	flag := re.MatchString(user.SMSNumber)
	if flag || !(len(user.SMSNumber) == 11) {
		user.SMSNumber = smscode.GetSMSNumber()
		user.SendTimes = 0
	}

	for {
		temp, err := os.ReadFile("./sendtime.txt")
		if err == nil {
			if string(temp) != GetDate() {
				user.SendTimes = 0
			}
		}
		selection = -1
		fmt.Printf("\n当前手机号:%s\n", user.SMSNumber)
		fmt.Print("1:登录 2:发送验证码 3:修改手机号\n")
		fmt.Scan(&selection)

		switch selection {
		case 1:
			var Input string
			fmt.Print("输入验证码登录\n")
			fmt.Scan(&Input)
			if md5.Sum([]byte(Input)) == user.SMSCode && time.Since(sendtime) <= 5*time.Minute {
				fmt.Print("登录成功，保存登录信息\n")
				os.WriteFile("user.txt", []byte(user.SMSNumber+strconv.FormatInt(int64(user.SendTimes), 10)), 0666)
				return
			} else {
				fmt.Print("验证码错误或过期\n")
			}
		case 2:
			if user.SendTimes < 5 && (time.Since(sendtime) >= 1*time.Minute || sendtime == runtime) {
				fmt.Printf("发送验证码到%s\n", user.SMSNumber)
				user.SMSCode, sendtime = smscode.SendSMSCode()
				user.SendTimes++
				os.WriteFile("user.txt", []byte(user.SMSNumber+strconv.FormatInt(int64(user.SendTimes), 10)), 0666)
				date := GetDate()
				os.WriteFile("sendtime.txt", []byte(date), 0666)
			} else if user.SendTimes >= 5 {
				fmt.Print("当日发送次数过多\n")
			} else {
				fmt.Print("1分钟只能发送一次\n")
			}
		case 3:
			user.SMSNumber = smscode.GetSMSNumber()
		case -1:
		default:
			fmt.Print("输入不正确\n")
		}
	}
}

func GetDate() string {
	year, month, day := time.Now().Date()
	date := strconv.FormatInt(int64(year), 10) + strconv.FormatInt(int64(month), 10) + strconv.FormatInt(int64(day), 10)
	return date
}
