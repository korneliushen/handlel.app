package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	scrape()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}
