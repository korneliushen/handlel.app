package bunnpris

import (
	"context"
	"fmt"
)

type Products []Product

type Product struct {
	Gtin      string
	Name      string
	Link      string
	Category  string
	Price     float64
	ImageLink string
}

func (p *Products) Get(ctx context.Context, token, endpoint string) error {
	// Gjør en post request til Itemgroups.aspx endpoint i bunnpris api
	// post tar inn en ctx av type context.Context som brukes til
	// time ut funksjonen om den tar for lang tid
	res := POST(ctx, token, endpoint)
	if res.IsError() {
		return fmt.Errorf(res.Error())
	}

	// returnerer en error om ingen html ble returnert fra post requesten
	if res.Data.HTML == nil {
		return fmt.Errorf("No html returned")
	}

	products := res.GetProducts()

	*p = products

	return nil
}
