package bunnpris

import (
	"errors"
)

type Categories []Category

type Category struct {
	Name          string
	Link          string
	SubCategories []string
}

func (c *Categories) Get(token string) error {
	res := POST(token, "Itemgroups.aspx")
	if res.IsError() {
		return errors.New(res.Error())
	}

	categories := res.ParseHTML()

	*c = categories

	return nil
}
