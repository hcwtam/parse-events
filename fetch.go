package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type SectionsResponse struct {
	Parse struct {
		Sections []map[string]interface{}
	}
}

type WikiResponse struct {
	Parse struct {
		Text map[string]string
	}
}

func getJson(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func fetchAndParse(year int, country string, ch chan<- Data) {
	baseUrl := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=parse&format=json&page=%d_in_%s&disabletoc=1", year, url.QueryEscape(strings.ReplaceAll(country, " ", "_")))

	// Fetch section index for section "Event"
	r := new(SectionsResponse)
	url := baseUrl + "&prop=sections"
	err := getJson(url, r)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("cannot fetch section index for \"%d in %s\"\n", year, country)
		return
	}
	index := ""
	for _, value := range r.Parse.Sections {
		if value["line"] == "Events" || value["line"] == "Events by month" || value["line"] == "Monthly events" {
			index = value["index"].(string)
		}
	}
	if index == "" {
		fmt.Printf("Content does not exist for \"%d in %s\"\n", year, country)
		return
	}

	// Fetch actual content
	res := new(WikiResponse)
	url = baseUrl + fmt.Sprintf("&prop=text&section=%s", index)
	err = getJson(url, res)

	if err != nil {
		fmt.Printf("cannot fetch \"%d in %s\"\n", year, country)
		return
	}

	data := Data{
		year:    year,
		country: country,
		content: []byte(res.Parse.Text["*"]),
	}

	ch <- data
}

func createHTML(data Data) {
	name := fmt.Sprintf("data/%d_in_%s.html", data.year, strings.ReplaceAll(data.country, " ", "_"))
	fo, err := os.Create(name)
	if err != nil {
		fmt.Printf("cannot create file for \"%d in %s\": %s\n", data.year, data.country, err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Printf("cannot close file for \"%d in %s\": %s\n", data.year, data.country, err)
		}
	}()

	if _, err := fo.Write(data.content); err != nil {
		fmt.Printf("cannot write file for \"%d in %s\": %s\n", data.year, data.country, err)
	}
}
