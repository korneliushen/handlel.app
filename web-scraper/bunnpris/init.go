package main

import (
	"bunnpris/test/bunnpris"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func run() {
	ctx := context.Background()

	token := "wlli2e4n52sjbydqekk2nnki"

	var categories bunnpris.Categories
	if err := categories.Get(ctx, token); err != nil {
		fmt.Printf("Error getting categories: %v\n", err)
		return
	}

	json, _ := json.MarshalIndent(categories, "", "  ")
	fmt.Println(string(json))
}

func Init() {
	start := time.Now()

	run()

	fmt.Println("Elapsed:", time.Since(start))
}
