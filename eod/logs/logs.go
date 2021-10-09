package logs

import "os"

var DiscordLogs *os.File
var MysqLogs *os.File
var DataFile *os.File

func InitEoDLogs() {
	var err error
	DataFile, err = os.OpenFile("createlogs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	DiscordLogs, err = os.OpenFile("discordlogs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	MysqLogs, err = os.Create("mysqlogs.txt")
	if err != nil {
		panic(err)
	}
}
