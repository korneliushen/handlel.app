package bunnpris

import (
	"context"
	"fmt"
)

type Categories []Category

type Category struct {
	Name          string
	Link          string
	SubCategories []string
}

func (c *Categories) Get(ctx context.Context, token string) error {
	// Gjør en post request til Itemgroups.aspx endpoint i bunnpris api
	// post tar inn en ctx av type context.Context som brukes til
	// time ut funksjonen om den tar for lang tid
	res := POST(ctx, token, "/itemgroups.aspx")
	if res.IsError() {
		return fmt.Errorf(res.Error())
	}

	// returnerer en error om ingen html ble returnert fra post requesten
	if res.Data.HTML == nil {
		return fmt.Errorf("No html returned")
	}

	// parser html-en som returneres fra post requesten
	categories := res.GetCategories()

	// parsehtml returnerer en verdi med type Categories, så vi assigner
	// categories (c) funksjonen ble kjørt på til categories vi får fra
	// parseHTML
	*c = categories

	return nil
}
