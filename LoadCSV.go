package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func whichCSV(dateTime1 time.Time, dateTime2 time.Time, currencies string) {
	cmd := exec.Command("python3 PythonGetFiles/getFiles.py -startDate " + strconv.Itoa(int(dateTime1.Month())) +
		"+" + strconv.Itoa(dateTime1.Year()) + " -endDate " + strconv.Itoa(int(dateTime2.Month())) + "+" +
		strconv.Itoa(int(dateTime2.Month())) + " -currency " + currencies)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
