package controller

import (
	"crypto/md5"
	"fmt"
	"main/smscode"
	"time"
)

func Start(user *smscode.Information) {

	var recorder smscode.Recorder
	recorder.Runtime = time.Now()        //记录运行的时间
	recorder.Sendtime = recorder.Runtime //发送时间设定默认值

	user.GetSMSNumber(&recorder) //提示输入手机号

	menu(user, &recorder)
}

func menu(user *smscode.Information, recorder *smscode.Recorder) {
	var data smscode.Data = user //声明并初始化接口变量
	var selection string
	user.Date = smscode.GetDate()
	data.Save()
	for {

		if user.Date != smscode.GetDate() { //与现在不是同一天，则发送次数归零
			user.SendTimes = 0
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
				data.LogIn()
				return
			} else {
				fmt.Print("验证码错误或过期\n")
			}
		case "2":
			if user.SendTimes < 5 && (time.Since(recorder.Sendtime) >= 1*time.Minute || recorder.Sendtime == recorder.Runtime) { //检查发送次数是否达到上限或过于频繁
				fmt.Printf("发送验证码到%s\n", user.SMSNumber)
				user.SMSCode, recorder.Sendtime = smscode.SendSMSCode()
				user.SendTimes++
				user.Date = smscode.GetDate()
				data.Save() //储存发送次数
			} else if user.SendTimes >= 5 {
				fmt.Print("当日发送次数过多\n")
			} else {
				fmt.Print("1分钟只能发送一次\n")
			}
		case "3":
			user.GetSMSNumber(recorder)

		default:
			fmt.Print("输入不正确\n")
		}

	}

}
