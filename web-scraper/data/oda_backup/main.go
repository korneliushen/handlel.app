package main

import (
	"fmt"
	"time"

	"github.com/korneliushen/handlel.app/web-scraper/scraper"
)

func main() {
	start := time.Now()
	scraper.Scrape()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}
