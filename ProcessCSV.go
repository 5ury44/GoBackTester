package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func processCSV(csvIn string) [][3]float64 {

	rc, err := os.Open(csvIn) // import csv into space
	if err != nil {
		log.Fatal(err)
	}

	processedArr := make([][3]float64, 0)
	r := csv.NewReader(rc)

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		time := strings.Split(rec[1], ":")

		floatHour, err := strconv.ParseFloat(time[0][9:], 64)
		floatDay, err3 := strconv.ParseFloat((time[0][4:])[2:], 64)
		floatMinute, err1 := strconv.ParseFloat(time[1], 64)
		floatSecond, err2 := strconv.ParseFloat(time[2], 64)
		if err != nil && err1 != nil && err2 != nil && err3 != nil {
			log.Fatal(err)
		}

		floatBid, err := strconv.ParseFloat(rec[2], 64)
		floatAsk, err1 := strconv.ParseFloat(rec[3], 64)
		if err != nil && err1 != nil {
			log.Fatal(err)
		}

		keyTime := 86400*floatDay + 3600*floatHour + 60*floatMinute + floatSecond
		valBidAsk := [3]float64{keyTime, floatBid, floatAsk}
		processedArr = append(processedArr, valBidAsk)
	}
	return processedArr
}
