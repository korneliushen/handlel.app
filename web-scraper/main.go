package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	getProductInfo("", 0)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}
