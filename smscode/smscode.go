package smscode

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 结构体
type Information struct {
	Data
	SMSNumber string
	SMSCode   [16]byte
	SendTimes int
	Date      string
}
type Recorder struct {
	Sendtime time.Time
	Runtime  time.Time
}

// 接口
type Data interface {
	Save(Information)
}

func (i Information) Save(user Information) {
	temp, err := os.ReadFile("user.txt")
	if err != err {
		panic(err)
	}

	idx := strings.LastIndex(string(temp), user.SMSNumber)

	if idx != -1 {
		file, err := os.OpenFile("user.txt", os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteAt([]byte("#"+user.SMSNumber+"-"+strconv.FormatInt(int64(user.SendTimes), 10)+"-"+user.Date+"\n"), int64(idx-1))
	} else {
		file, err := os.OpenFile("user.txt", os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		write.WriteString("#" + user.SMSNumber + "-" + strconv.FormatInt(int64(user.SendTimes), 10) + "-" + user.Date + "\n")
		write.Flush()
	}

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

func GetSMSNumber(user *Information, recorder *Recorder) {
	user.SMSNumber = ""
	for {

		fmt.Print("输入电话号码:")
		fmt.Scan(&user.SMSNumber)
		re := regexp.MustCompile("^[0-9]{11}$")
		flag := re.MatchString(user.SMSNumber) //检查输入格式是否正确
		if flag {
			break
		} else {
			fmt.Print("输入不正确\n")
		}
	}
	temp, err := os.ReadFile("./user.txt")
	if err != nil {
		panic(err)
	}
	pos := strings.LastIndex(string(temp), user.SMSNumber)
	if pos != -1 {
		idx := GetIndex(pos-1, temp)
		if idx[0] != -1 {
			user.Date = string(temp[idx[2]+1 : idx[2]+8])
			if user.Date == GetDate() {
				user.SendTimes, err = strconv.Atoi(string(temp[idx[1]+1 : idx[2]]))
				if err != nil {
					panic(err)
				}
			} else {
				user.SendTimes = 0
			}
		}
		recorder.Runtime = time.Now()        //记录运行的时间
		recorder.Sendtime = recorder.Runtime //发送时间设定默认值
	} else {
		user.Date = GetDate()
		user.SendTimes = 0
		recorder.Runtime = time.Now()        //记录运行的时间
		recorder.Sendtime = recorder.Runtime //发送时间设定默认值
		user.Save(*user)
	}
}

func GetDate() string {
	year, month, day := time.Now().Date()
	date := strconv.FormatInt(int64(year), 10) + strconv.FormatInt(int64(month), 10) + strconv.FormatInt(int64(day), 10)
	return date
}

func GetIndex(start int, content []byte) [3]int {
	var result [3]int
	idx := start
	result[0] = strings.LastIndex(string(content[start:]), "#") + start
	if result[0] == -1 {
		return result
	}
	for i, j := idx, 1; i < idx+21; i++ {
		if content[i] == '-' {
			idx++
			if j >= 3 {
				result[0] = -1
				break
			}
			result[j] = i
			j++
		}

	}
	return result
}
