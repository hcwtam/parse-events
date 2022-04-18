package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getCountries() (countries []string) {
	countries = make([]string, 0)

	file, err := os.Open("./countries.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Text() != "" {

			countries = append(countries, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func parse(data Data) Data {

	content := string(data.content)

	// replace weird dash character with the correct one
	content = strings.ReplaceAll(content, "–", "-")
	// replace href with full wikipedia url, and add target="_blank"
	content = strings.ReplaceAll(content, "href=\"/wiki/", "target=\"_blank\" href=\"https://en.wikipedia.org/wiki/")

	reader := strings.NewReader(content)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Printf("cannot parse file for \"%d in %s\": %s\n", data.year, data.country, err)
	}

	// Remove the wikipedia "[edit]" elements
	doc.Find(".mw-editsection").Remove()
	// Remove citation superscript links
	doc.Find(".reference").Each(func(i int, s *goquery.Selection) {
		text := s.Find("a").Text()
		s.SetText(text)
		s.Find("a").Remove()
	})
	doc.Find(".mw-cite-backlink").Remove()

	c, err := doc.Find("body").Find(".mw-parser-output").Html()
	if err != nil {
		fmt.Printf("cannot parse file for \"%d in %s\": %s\n", data.year, data.country, err)
	}

	data.content = []byte(strings.Split(c, "<!--")[0]) // Remove comment at the end of html

	return data
}
