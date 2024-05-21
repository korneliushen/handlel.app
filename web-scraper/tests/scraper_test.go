package tests

import (
	"testing"

	w "github.com/korneliushen/handlel.app/web-scraper/scraper"
	"github.com/stretchr/testify/assert"
)

func TestGetProductInfo(t *testing.T) {
	tests := []struct {
		name          string
		link          string
		cursor        int
		underCategory *w.Underkategori
	}{
		{"Scrape data fra underkategori Øl cursor 2", "https://oda.com/no/categories/60-drikke/1237-ol-og-alkoholfritt/?filters=", 2, &w.Underkategori{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := w.GetProductInfo(test.link, test.cursor, test.underCategory); err != nil {
				t.Log("Error getting data for link: ", test.link, "and cursor: ", test.cursor)
				t.Fail()
			}
			assert.IsType(t, w.Underkategori{}, *test.underCategory, "underCategory should be of type Underkategori")
		})
	}

}
