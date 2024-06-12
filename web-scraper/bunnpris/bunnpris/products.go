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

func (products *BunnprisProducts) GetProducts(apiProducts *model.ApiProducts, ctx context.Context, token, id string) error {
	// Gjør en først request for å få mengden produkter i kategorien i tillegg
	// til å få data fra første page
	body := genBody(1, id)
	res := POST(ctx, token, "/productwebmethod.aspx/LoadItemGroupListControl", bytes.NewBuffer(body), "application/json; charset=iso-8859-1")
	if res.IsError() {
		return fmt.Errorf(res.Error())
	}

	// returnerer en error om ingen json ble returnert fra post requesten
	if res.Data.JSON == nil {
		return fmt.Errorf("No json returned\n")
	}

	// Unmarshaler json vi får fra post requesten til formatet definert i
	// CategoryRes structet
	var categoryData CategoryRes
	err := json.Unmarshal(res.Data.JSON, &categoryData)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	totalPages := getPageCount(categoryData.Data.TotalProducts)

	// For hver side sendes en request til api-en med en body med som bruker i+1
	// for å hente data fra siden vi er på
	for i := range totalPages {
		body := genBody(i+1, id)
		res := POST(ctx, token, "/productwebmethod.aspx/LoadItemGroupListControl", bytes.NewBuffer(body), "application/json; charset=iso-8859-1")
		if res.IsError() {
			return err
		}

		// Unmarshaler dataen vi får og henter ut Output (som er html data i en
		// string)
		var responseData CategoryRes
		err := json.Unmarshal(res.Data.JSON, &responseData)
		if err != nil {
			return err
		}

		// Parser html-en fra string til *html.Node og kjører GetProductLinks
		// I GetProductLinks legges linken til i products arrayet
		html, err := html.Parse(strings.NewReader(responseData.Data.Output))
		if err != nil {
			return err
		}
		products.GetProductLinks(html)
	}

	return nil
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

func (products *BunnprisProducts) GetProductLinks(data *html.Node) {
	// Definerer en funksjon som går gjennom base noden
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// Mapper over alle attributter elementet har
			// Om attr sin value er er lblName, legges slug (linken til produktet)
			// til i et BunnprisProducts arrayet (array med strings)
			for _, attr := range node.Attr {
				if attr.Val == "lblName" {
					*products = append(*products, strings.Split(strings.Split(node.Parent.Attr[1].Val, ".no")[1], "&grpnm")[0])
				}
			}
		}

		// går gjennom html-noden og så går til sibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	// kjører crawler på base-noden
	crawler(data)
}
