package bunnpris

import (
	"context"
	"fmt"

	"github.com/korneliushen/handlel.app/scraper/model"
)

func GetProducts(apiProducts *model.ApiProducts, ctx context.Context, token, endpoint string) error {
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

	res.GetProducts(apiProducts)

	return nil
}
