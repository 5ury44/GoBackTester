package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { //read header
			log.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {

					break
				}
				log.Fatal(err)

			}
			ch <- rec
		}
	}()
	return
}

func main() {

	csvIn, err := os.Open("EURUSD-2022-11.csv") // import csv into space
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(len(processCSV(csvIn)))
}
