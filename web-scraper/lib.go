// funksjoner flyttet hit for å gjøre scraper mere readable
// ikke så viktig å forstå disse for å forstå scraper

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// om det bare er en side så vil tall være 1 siffer, så antallsider blir bare satt til 1
// hvis det er flere sider, kjør en scuffed alg for å få tallet
func getPageCount(pageCountSlice []string) int {
	if len(pageCountSlice) <= 2 {
		return 1
	}

	pageCountTextLength := len(pageCountSlice) / 2
	var s []string
	for i := 0; i < pageCountTextLength; i++ {
		s = append(s, pageCountSlice[i])
	}
	str := strings.Join(s, "")
	amountOfProducts, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 1
	}
	return int(math.Ceil(float64(amountOfProducts / 24)))
}

// ingredienser er satt opp annerledes avhengig av hvor mange ingredienser det her, så det accountes for her
func getIngredients(e *colly.HTMLElement, value string) string {
	if value != "" {
		// om det er få ingredienser
		return e.ChildText("div > div p")
	}

	// om det er flere ingredienser (da havner også Ingredienser i value, så det fjernes)
	return strings.Split(e.ChildText("div > div span"), "Ingredienser")[1]
}

// Sammenligner key til innholdet (navnet), med et field i Innhold structen
// om den finner en key som matcher en field, legges det til i instansen av Innhold
// alt med "n" foran seg, er for næringsinnhold
func setFieldValue(contents *Innhold, nContents *Næringsinnhold, key string, value string) {
	v := reflect.ValueOf(contents).Elem()

	// denne koden gjør at alt med hvorav funker som det skal
	key = strings.Title(key)
	key = strings.ReplaceAll(key, " ", "")

	// finner fieldet som matcher key
	field := v.FieldByName(key)

	if !field.IsValid() || !field.CanSet() {
		nV := reflect.ValueOf(nContents).Elem()
		nField := nV.FieldByName(key)

		if !nField.IsValid() || !nField.CanSet() {
			fmt.Printf("Cannot set field %s\n", key)
			return
		}

		nField.SetString(value)
		return
	}

	field.SetString(value)
}

func writeData(data Kategorier) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return
	}

	err = os.WriteFile("./data.json", jsonData, 0666)
	if err != nil {
		fmt.Printf("Error writing json data to file %v\n", err)
		return
	}

	fmt.Println("Data written: ", string(jsonData))
}
