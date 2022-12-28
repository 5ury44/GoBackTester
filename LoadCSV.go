package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	idcat, err := findIdCat(current.Month().String())
	if err != nil {
		fmt.Println("Error parsing webpage:", err)
		return
	}

	// Send a GET request to truefx
	response, err := client.Get("https://www.truefx.com/truefx-historical-downloads/#93-" +
		strconv.Itoa(idcat) + "-top")
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

func findIdCat(month string) (int, error) {

	resp, err := http.Get("https://www.truefx.com/truefx-historical-downloads/")
	if err != nil {
		fmt.Println(err)
		return -1, nil
	}
	defer resp.Body.Close()

	htmlBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return -1, nil
	}

	tokenizer := html.NewTokenizer(strings.NewReader(string(htmlBody)))

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			t := tokenizer.Token()
			if t.Data == "div" {
				for _, attr := range t.Attr {
					if attr.Key == "title" && attr.Val == month {
						for _, attr := range t.Attr {
							if attr.Key == "data-idcat" {
								fmt.Println(attr.Val)
								return strconv.Atoi(attr.Val)
							}
						}
					}
				}
			}
		}
	}
	return -1, nil
}
