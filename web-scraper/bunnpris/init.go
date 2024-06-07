package main

import (
	"bunnpris/test/bunnpris"
	"encoding/json"
	"fmt"
	"time"
)

func run() {
	token := "wlli2e4n52sjbydqekk2nnki"

	var categories bunnpris.Categories
	if err := categories.Get(token); err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}

	json, _ := json.MarshalIndent(categories, "", "  ")
	fmt.Println(string(json))
}

func Init() {
	start := time.Now()

	run()

	fmt.Println("Elapsed:", time.Now().Sub(start))
}
