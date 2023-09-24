package main

import (
	"main/controller"
	"main/smscode"
	"os"
	"strconv"
	"strings"
)

func main() {
	var user smscode.Information
	_, err := os.Stat("./user.txt")
	if err == nil {
		temp, err := os.ReadFile("./user.txt")
		if err != nil {
			panic(err)
		}
		pos := strings.LastIndex(string(temp), "#")
		if pos == -1 {
			pos = 0
		}
		idx := smscode.GetIndex(pos, temp)
		if idx[0] != -1 {
			user.SMSNumber = string(temp[idx[0]+1 : idx[1]])
			user.SendTimes, err = strconv.Atoi(string(temp[idx[1]+1 : idx[2]]))
			if err != nil {
				panic(err)
			}
			user.Date = string(temp[idx[2]+1 : idx[2]+6])

		}
	} else {
		_, err = os.Create("user.txt")
		if err != nil {
			panic(err)
		}
	}
	controller.Start(&user)
}
