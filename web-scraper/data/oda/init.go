package oda

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/korneliushen/handlel.app/scraper/model"
)

//TODO: split into separate files

type ProductsResp struct {
	Blocks []ProductsRespBlock `json:"blocks"`
}

type ProductsRespBlock struct {
	Id string `json:"id"`
	Items []ProductsRespBlockItem `json:"items"`
}

type ProductsRespBlockItem struct {
	Data ProductsRespBlockItemData `json:"data"`
}

type ProductsRespBlockItemData struct {
	Id string `json:"id"`
	Title string `json:"title"`
}

//TODO: actual error messages
func Fetch(products *model.Products) {
	resp, err := http.Get("https://oda.com/tienda-web-api/v1/app-components/products/")
	if err != nil {
		fmt.Println("Error in app-components/products get req")
		return
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 || err != nil {
		fmt.Println("Error in getting products")
		return
	}
	var productsResp ProductsResp
	err = json.Unmarshal(body, &productsResp)
	if err != nil {
		fmt.Println("unmarshal json fail")
		return
	}
	fmt.Println(productsResp)
}

