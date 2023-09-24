package controller

import (
	"crypto/md5"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"main/smscode"
)

type Recorder struct {
	Sendtime time.Time
	Runtime  time.Time
	Date     string
}

func Start(user smscode.Information) {

	var recorder Recorder
	recorder.Runtime = time.Now()        //记录运行的时间
	recorder.Sendtime = recorder.Runtime //发送时间设定默认值
	recorder.Date = GetDate()
	err := os.WriteFile("sendtime.txt", []byte(recorder.Date), 0666) //运行时自动储存运行时间
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("[0-9]{11}")
	flag := re.MatchString(user.SMSNumber)
	if !flag { //检查文件中储存的电话格式是否正确
		user.SMSNumber = smscode.GetSMSNumber() //如无法自动登录，提示输入手机号
		user.SendTimes = 0
	}

	menu(user, recorder)
}

func menu(user smscode.Information, recorder Recorder) {
	var selection string
	for {
		temp, err := os.ReadFile("./sendtime.txt") //读取文件中储存的上次发送时间
		if err == nil {
			if string(temp) != GetDate() { //与现在不是同一天，则发送次数归零
				user.SendTimes = 0
			}
		}
		fmt.Printf("\n当前手机号:%s\n", user.SMSNumber)
		fmt.Print("1:登录 2:发送验证码 3:修改手机号\n")
		fmt.Scan(&selection)

		switch selection {
		case "1":
			var Input string
			fmt.Print("输入验证码登录\n")
			fmt.Scan(&Input)
			if md5.Sum([]byte(Input)) == user.SMSCode && time.Since(recorder.Sendtime) <= 5*time.Minute { //检查输入的验证码是否正确或者过期
				fmt.Print("登录成功，保存登录信息\n")
				os.WriteFile("user.txt", []byte(user.SMSNumber+strconv.FormatInt(int64(user.SendTimes), 10)), 0666)
				return
			} else {
				fmt.Print("验证码错误或过期\n")
			}
		case "2":
			if user.SendTimes < 5 && (time.Since(recorder.Sendtime) >= 1*time.Minute || recorder.Sendtime == recorder.Runtime) { //检查发送次数是否达到上限或过于频繁
				fmt.Printf("发送验证码到%s\n", user.SMSNumber)
				user.SMSCode, recorder.Sendtime = smscode.SendSMSCode()
				user.SendTimes++
				os.WriteFile("user.txt", []byte(user.SMSNumber+strconv.FormatInt(int64(user.SendTimes), 10)), 0666) //储存发送次数
				recorder.Date = GetDate()
				os.WriteFile("sendtime.txt", []byte(recorder.Date), 0666) //储存发送时间
			} else if user.SendTimes >= 5 {
				fmt.Print("当日发送次数过多\n")
			} else {
				fmt.Print("1分钟只能发送一次\n")
			}
		case "3":
			user.SMSNumber = smscode.GetSMSNumber()
			user.SendTimes = 0
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
