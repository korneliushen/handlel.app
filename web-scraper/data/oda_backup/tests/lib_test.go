package tests

import (
	"os"
	"testing"

	w "github.com/korneliushen/handlel.app/web-scraper/scraper"
	"github.com/stretchr/testify/assert"
)

// det er 24 elementer per side, så antall produkter blir delt på 24 for å få produkter
func TestGetPageCount(t *testing.T) {
	tests := []struct {
		name           string
		pageCountSlice []string
		expected       int
	}{
		{"1 produkt -> 1 side", []string{"1", "1"}, 1},
		{"72 produkter -> 3 sider", []string{"7", "2", "7", "2"}, 3},
		{"25 produkter -> 2 sider (runder opp)", []string{"2", "5", "2", "5"}, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := w.GetPageCount(test.pageCountSlice)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestWriteData(t *testing.T) {
	test := struct {
		name string
		data w.Kategorier
		path string
	}{
		name: "Write and read file",
		data: w.Kategorier{},
		path: "../data_test.json",
	}

	t.Run(test.name, func(t *testing.T) {
		err := w.WriteData(test.data, test.path)
		if err != nil {
			t.Log("Error writing to json", err)
			t.Fail()
		}
		assert.FileExists(t, test.path)
		os.Remove(test.path)
	})
}
