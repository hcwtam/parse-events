package main

import (
	"sync"
	"time"
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

	for year := 1990; year < 2022; year++ {
		for _, country := range countries {
			wg.Add(1)
			go func(y int, c string) {
				fetchAndParse(y, c, ch)
				wg.Done()
			}(year, country)
		}

		time.Sleep(5 * time.Second)
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
