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

var pairs = make([]int, 0)

func findPairs(currs []string) {

	for i, str1 := range currs {
		for _, str2 := range currs[i+1:] {
			for j, str3 := range currencyPairs {
				if strings.Contains(str3, str2) && strings.Contains(str3, str1) {
					pairs = append(pairs, j)
				}
			}
		}
	}
}

func currencyMap() map[string][]instant {
	mapCurr := make(map[string][]instant)

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

		const layout = "20060102 15:04:05.000"
		date, err := time.Parse(layout, rec[1])

		//date := time.Date(intYear, time.Month(11+0*intMonth), floatDay, floatHour, floatMinute, floatSecond, 1000000*floatNs, time.UTC)

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
