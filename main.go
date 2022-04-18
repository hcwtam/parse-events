package main

import (
	"sync"
)

type Data struct {
	year    int
	country string
	content []byte
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan Data)
	countries := getCountries()

	println("Generating data...")

	for _, country := range countries {
		for year := 2000; year < 2022; year++ {
			wg.Add(1)
			go func(y int, c string) {
				fetchAndParse(y, c, ch)
				wg.Done()
			}(year, country)
		}
	}

	go func() {
		wg.Wait()
		println("Completed.")
		close(ch)
	}()

	for data := range ch {
		d := parse(data)
		createHTML(d)
	}

}
