package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func processCSV(rc io.Reader) map[float64][2]float64 {

	processedMap := make(map[float64][2]float64)
	r := csv.NewReader(rc)

	// handle header
	rec, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		rec, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		time := strings.Split(rec[1], ":")

		floatHour, err := strconv.ParseFloat(time[0][9:], 64)
		floatMinute, err1 := strconv.ParseFloat(time[1], 64)
		floatSecond, err2 := strconv.ParseFloat(time[2], 64)
		if err != nil && err1 != nil && err2 != nil {
			log.Fatal(err)
		}

		floatBid, err := strconv.ParseFloat(rec[2], 64)
		floatAsk, err1 := strconv.ParseFloat(rec[3], 64)
		if err != nil && err1 != nil {
			log.Fatal(err)
		}

		keyTime := 3600*floatHour + 60*floatMinute + floatSecond
		valBidAsk := [2]float64{floatBid, floatAsk}
		processedMap[keyTime] = valBidAsk

	}
	return processedMap
}

func main() {

	csvIn, err := os.Open("EURUSD-2022-11.csv") // import csv into space
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(len(processCSV(csvIn)))
}