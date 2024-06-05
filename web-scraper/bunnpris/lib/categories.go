package lib

import (
	"bunnpris/test/api"
	"errors"
	"fmt"
)

const endpoint = "https://nettbutikk.bunnpris.no/Itemgroups.aspx"

func getCategories(token string) error {
	res := api.POST(token, "Itemgroups.aspx")
	if res.IsError() {
		return errors.New(res.Error())
	}

	categories := res.ParseHTML()

	fmt.Printf("%+v", categories)

	return nil
}
