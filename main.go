package main

import (
	"./controller"
	"./smscode"
	"os"
	"strconv"
)

func main() {
	var user smscode.Information
	_, err := os.Stat("./user.txt")
	if err == nil {
		temp, err := os.ReadFile("./user.txt")
		if err != nil {
			panic(err)
		}
		user.SMSNumber = string(temp[:len(temp)-1])
		user.SendTimes, err = strconv.Atoi(string(temp[len(temp)-1]))
		if err != nil {
			panic(err)
		}
	}
	controller.Start(user)
}
