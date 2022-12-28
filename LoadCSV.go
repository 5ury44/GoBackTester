package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func whichCSV(dateTime1 time.Time, dateTime2 time.Time, currencies string) {
	current := time.Date(dateTime1.Year(), dateTime1.Month(), 0, 0, 0, 0, 0,
		dateTime1.Location())
	for current.Before(dateTime2) {
		getCSV(currencies, current)
		current = current.AddDate(0, 1, 0)
	}

}

func getCSV(currencies string, current time.Time) {
	if _, err := os.Stat("files"); os.IsNotExist(err) {
		if err := os.Mkdir("files", 0755); err != nil {
			fmt.Printf("Error creating 'files' directory: %v", err)
			return
		}
	}

	client := &http.Client{}

	// Send a GET request to truefx
	response, err := client.Get("https://www.truefx.com/truefx-historical-downloads/#93-98-april")
	if err != nil {
		fmt.Println("Error retrieving webpage:", err)
		return
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("Error parsing webpage:", err)
		return
	}

	var fileURL string
	var found bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "title" && a.Val == currencies+"-"+strconv.Itoa(current.Year())+"-"+
					strconv.Itoa(current.Day()) {
					for _, a := range n.Attr {
						if a.Key == "href" {
							fileURL = a.Val
							found = true
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if !found {
		fmt.Println("File not found")
		return
	}

	response, err = client.Get(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer response.Body.Close()

	fileName := filepath.Base(fileURL)
	file, err := os.Create(filepath.Join("files", fileName))
	if err != nil {
		fmt.Printf("error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("error writing to file: %v\n", err)
		return
	}
}
