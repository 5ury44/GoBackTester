package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type instant struct {
	date time.Time
	ask  float64
	bid  float64
}

var currencyPairs = []string{"audjpy", "audnzd", "audusd", "cadjpy", "chfjpy", "eurchf", "eurgbp", "eurjpy", "eurpln",
	"eurusd", "gbpjpy", "gbpusd", "nzdusd", "usdcad", "usdchf", "usdjpy", "usdmxn", "usdrub", "usdtry", "usdzar"}

func currencyMap(currs []string) map[string][]instant {
	mapCurr := make(map[string][]instant)
	var pairs []int
	for i, str1 := range currs {
		for _, str2 := range currs[i+1:] {
			for i, str3 := range currencyPairs {
				if strings.Contains(str3, str2) && strings.Contains(str3, str1) {
					pairs = append(pairs, i)
				}
			}
		}
	}

	files, err := ioutil.ReadDir("PythonGetFiles/files")
	if err != nil {
		fmt.Println(err)
	}

	for _, index := range pairs {
		processedArr := make([]instant, 0)
		for _, file := range files {
			if strings.Contains(file.Name(), strings.ToUpper(currencyPairs[index])) {
				processedArr = append(processedArr, processCSV("PythonGetFiles/files/"+file.Name())...) //todo see if works
			}
		}
		mapCurr[currencyPairs[index]] = processedArr
	}

	return mapCurr
}

func processCSV(csvIn string) []instant {
	processedArr := make([]instant, 0)
	rc, err := os.Open(csvIn) // import csv into space
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(rc)

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		split := strings.Split(rec[1], ":")

		intYear, err4 := strconv.Atoi(csvIn[7:10])
		intMonth, err5 := strconv.Atoi(csvIn[12:13])
		floatHour, err := strconv.Atoi(split[0][9:])
		floatDay, err3 := strconv.Atoi((split[0][4:])[2:])
		floatMinute, err1 := strconv.Atoi(split[1])
		floatSecond, err2 := strconv.ParseFloat(split[2], 64)
		if err != nil && err1 != nil && err2 != nil && err3 != nil && err4 != nil && err5 != nil {
			log.Fatal(err)
		}
		date := time.Date(intYear, time.Month(intMonth), floatDay, floatHour, floatMinute, 0, int(floatSecond*1000000000), time.UTC)

		floatBid, err := strconv.ParseFloat(rec[2], 64)
		floatAsk, err1 := strconv.ParseFloat(rec[3], 64)
		if err != nil && err1 != nil {
			log.Fatal(err)
		}
		in := instant{date, floatAsk, floatBid}

		processedArr = append(processedArr, in)
	}
	return processedArr
}
