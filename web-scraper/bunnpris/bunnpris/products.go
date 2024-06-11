package bunnpris

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/korneliushen/handlel.app/scraper/model"
	"golang.org/x/net/html"
)

type CategoryRes struct {
	Data Data `json:"d"`
}

type Data struct {
	TotalProducts string `json:"TotalRecords"`
	Output        string `json:"OutputString"`
}

// TODO: Trenger ikke å gjøre en ekstra request for den første siden (gidder ikke fikse nå)

func GetProducts(apiProducts *model.ApiProducts, ctx context.Context, token, id string) (BunnprisProducts, error) {
	// Gjør en først request for å få mengden produkter i kategorien i tillegg
	// til å få data fra første page
	body := genBody(1, id)
	res := POST(ctx, token, "/productwebmethod.aspx/LoadItemGroupListControl", bytes.NewBuffer(body), "application/json; charset=iso-8859-1")
	if res.IsError() {
		return BunnprisProducts{}, fmt.Errorf(res.Error())
	}

	// returnerer en error om ingen json ble returnert fra post requesten
	if res.Data.JSON == nil {
		return BunnprisProducts{}, fmt.Errorf("No json returned\n")
	}

	// Unmarshaler json vi får fra post requesten til formatet definert i
	// CategoryRes structet
	var categoryData CategoryRes
	err := json.Unmarshal(res.Data.JSON, &categoryData)
	if err != nil {
		return BunnprisProducts{}, fmt.Errorf(err.Error())
	}

	// Lager en instanse av BunnprisProducts (er bare et array av strings)
	// Denne mappes over senere for å hente data for hvert produkt
	products := &BunnprisProducts{}

	fmt.Println(categoryData.Data.TotalProducts)

	totalPages := getPageCount(categoryData.Data.TotalProducts)

	fmt.Println(totalPages)

	// For hver side sendes en request til api-en med en body med som bruker i+1
	// for å hente data fra siden vi er på
	for i := range totalPages {
		body := genBody(i+1, id)
		fmt.Println(string(body))
		res := POST(ctx, token, "/productwebmethod.aspx/LoadItemGroupListControl", bytes.NewBuffer(body), "application/json; charset=iso-8859-1")
		if res.IsError() {
			return BunnprisProducts{}, err
		}

		// Unmarshaler dataen vi får og henter ut Output (som er html data i en
		// string)
		var responseData CategoryRes
		err := json.Unmarshal(res.Data.JSON, &responseData)
		if err != nil {
			return BunnprisProducts{}, err
		}

		// Parser html-en fra string til *html.Node og kjører GetProductLinks
		// I GetProductLinks legges linken til i products arrayet
		html, err := html.Parse(strings.NewReader(responseData.Data.Output))
		if err != nil {
			return BunnprisProducts{}, err
		}
		GetProductLinks(html, products)
	}

	return *products, nil
}

func genBody(page int, id string) []byte {
	return []byte(fmt.Sprintf("%s%d%s%s%s", `{dnItemParams: '{"PageName":"itemgrouplist","PageNo":"`, page, `", "ItemGroupNo":"`, id, `"}'}`))
}

func getPageCount(totalProducts string) int {
	// Bruker TotalProducts til å kalkulere hvor mange sider kategorien har
	totalProductsInt, err := strconv.Atoi(totalProducts)
	if err != nil {
		return 0
	}
	if totalProductsInt <= 18 {
		return 1
	}
	return int(math.Ceil(float64(totalProductsInt) / 18))
}
