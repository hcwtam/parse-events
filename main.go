package main

import (
	"fmt"
	"os"
)

type WikiResponse struct {
	Parse struct {
		Text map[string]string
	}
}

func main() {
	res := new(WikiResponse)
	err := getJson("https://en.wikipedia.org/w/api.php?action=parse&format=json&page=2000_in_the_Netherlands&prop=text&section=2&disabletoc=1", res)

	if err != nil {
		fmt.Printf("cannot fetch %s in %d\n", "The Netherlands", 2021)
	}

	fmt.Println(res.Parse.Text["*"])

	fo, err := os.Create("data/the_netherlands_2021.html")
	if err != nil {
		fmt.Printf("cannot create file for %s in %d: %s\n", "The Netherlands", 2021, err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Printf("cannot close file for %s in %d: %s\n", "The Netherlands", 2021, err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := []byte(res.Parse.Text["*"])

	if _, err := fo.Write(buf); err != nil {
		fmt.Printf("cannot write file for %s in %d: %s\n", "The Netherlands", 2021, err)
	}

}
