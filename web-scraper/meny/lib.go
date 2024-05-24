package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteData(data any, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0666)
	if err != nil {
		return err
	}

	fmt.Println("Data written: ", string(jsonData))

	return nil
}
