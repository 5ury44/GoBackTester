package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func whichCSV(dateTime1 time.Time, dateTime2 time.Time, currencies string) {

	cms := exec.Command("python", "/Users/surycuh/GolandProjects/GoOrderBook1/PythonGetFiles/getFiles.py", "--startDate", strconv.Itoa(int(dateTime1.Month()))+
		"+"+strconv.Itoa(dateTime1.Year()), "--endDate", strconv.Itoa(int(dateTime2.Month()))+"+"+
		strconv.Itoa(dateTime2.Year()), "--currency", currencies)
	err := cms.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(cms)

}
