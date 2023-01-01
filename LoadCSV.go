package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func whichCSV(dateTime1 time.Time, dateTime2 time.Time, currencies string) {
	current := time.Date(dateTime1.Year(), dateTime1.Month(), 1, 0, 0, 0, 0,
		dateTime1.Location())
	for current.Before(dateTime2) {
		func() {
			fmt.Println(current.Month())
			getCSV(currencies, current)
		}()
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	loginRequest := Login{
		Username: "5ury44",
		Password: "Shreya_vr12",
	}
	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		fmt.Println("Error parsing webpage:", err)
	}

	resp, err := client.Post("https://www.truefx.com/truefx-login/", "application/json",
		bytes.NewReader(requestBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	idcat, err := findIdCat(current.Month().String())
	if err != nil || idcat == -1 {
		fmt.Println("Error parsing webpage:", err)
		return
	}

	// Send a GET request to truefx
	resp, err = client.Get("https://www.truefx.com/truefx-historical-downloads/#93-" +
		strconv.Itoa(idcat) + "-top")
	if err != nil {
		fmt.Println("Error retrieving webpage:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing webpage:", err)
		return
	}

	//var fileURL string
	var found bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					var b bytes.Buffer
					html.Render(&b, n)
					fmt.Println(b.String())
					//fileURL = a.Val
					found = true
					return
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

	/*response, err = client.Get(fileURL)
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
	}*/
}

func findIdCat(month string) (int, error) {

	htmlData, err := ioutil.ReadFile("trueFX.html")
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	doc, err := html.Parse(strings.NewReader(string(htmlData)))
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	id := findElementByTitle(doc, month)
	fmt.Println(id)
	return strconv.Atoi(id)
}

func findElementByTitle(n *html.Node, title string) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "data-idcat" {
				var b bytes.Buffer
				html.Render(&b, n.FirstChild.NextSibling)
				if strings.Contains(b.String(), title) {
					fmt.Println(b.String())
					// Return the value of the idcat attribute
					return attr.Val
				}
			}
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		id := findElementByTitle(c, title)
		if id != "" {
			return id
		}
	}

	return ""

}
