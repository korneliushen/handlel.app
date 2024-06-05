package lib

import (
	"fmt"
	"time"
)

func run() {
	token := "wlli2e4n52sjbydqekk2nnki"

	err := getCategories(token)
	if err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}
}

func Init() {
	start := time.Now()

	run()

	fmt.Println("Elapsed:", time.Now().Sub(start))
}
