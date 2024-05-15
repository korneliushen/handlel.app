package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	getProducts()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}
