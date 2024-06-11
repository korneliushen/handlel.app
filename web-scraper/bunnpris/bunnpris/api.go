package bunnpris

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/korneliushen/handlel.app/scraper/model"
	"golang.org/x/net/html"
)

// handler responser fra api-en, har med message, status kode og data
// data kan være av type *html.Node eller en JSON type (har ikke definert ennå)
type Response struct {
	Message    string
	StatusCode int
	Data       ResponseData
}

// typene til ApiRes.data
type ResponseData struct {
	// TODO: lage method til ResponseData som parser html
	HTML *html.Node
	// TODO: gjøre dette til en generic med json responses som kan komme
	JSON []byte
}

// lager en Error() method for ApiRes, som gir en formatert versjon av
// message og status code
func (data Response) Error() string {
	return fmt.Sprintf(
		"api error: %s, with status code: %d", data.Message, data.StatusCode)
}

// lager en isError() method som returnerer true eller false basert på
// status code. dette gjør at vi kan kjøre res.IsError() { "handle error" }
func (data Response) IsError() bool {
	if data.StatusCode != http.StatusOK {
		return true
	}
	return false
}

// et array med denne dataen returneres fra ParseHTML
// senere skal det kunne være generic hva som returneres

// TODO: kan bare kjøre om Data.HTML eksisterer?
// TODO: gjøre det dynamic (generics kanskje) sånn at man kan velge
// TODO: hva man vil returnere (element, type osv.) (nå er det bare href og title)
// legger til en method til Response struct som parser html og returnerer
// alle instanser av et element
func (data Response) GetCategories() Categories {
	var categories Categories

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// init variabler som brukes til å holde dataen som blir funnet
			var className, link, name, id string

			// mapper over alle attributter elementet har
			// om attr er itemgroup link og har en href, blir linken til href-en
			// lagret som href
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "itemgroup-link" {
					className = attr.Val
				}
				if attr.Key == "href" {
					href := attr.Val
					link = strings.Split(href, "&")[0]
					// Henter itemgrpno som brukes som id
					// Link henter noen ting den ikke skal, så gjør en sjekk om at det er
					// minst to elementer så den ikke panicer
					groupNo := strings.Split(link, "itemgrpno=")
					if len(groupNo) == 2 {
						id = groupNo[1]
					}
				}
			}

			// om classname og href ble funnet, hentes dataen fra span-en som er
			// child element (må kjøre FirstChild.NextSibling fordi det første er
			// ::before). må også kjøre Firstchild etter nextsibling for å få teksten
			// og ikke bare navnet på html-elementet
			if className != "" && link != "" {
				if node.FirstChild != nil {
					name = node.FirstChild.NextSibling.FirstChild.Data
					// Av en eller annen grunn så er FRUKT/GRNT (Pris pr. stk) annerledes
					// enn alle andre kategorier og vanlig endpoint funker ikke, så har
					// hard-coda inn en løsning her
					if strings.Contains(name, "Pris pr. stk") {
						link = "/itemgrouplist.aspx?grpnm=FRUKT/GR%D8NT%20(Pris%20pr.%20stk)%20&deptno=16"
					}
				}
				categories = append(categories, Category{Id: id, Name: name, Link: link})
			}
		}

		// går gjennom html-noden og så går til sibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	// kjører crawler på base-noden
	crawler(data.Data.HTML)

	return categories
}

type BunnprisProducts []string

func GetProductLinks(data *html.Node, products *BunnprisProducts) {
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

func (products BunnprisProducts) FetchProductPages(ctx context.Context, token string, apiProducts *model.ApiProducts) {
	for _, link := range products {
		res := POST(ctx, token, link, nil, "text/html; charset=us-ascii")
		if res.IsError() {
			// Om det er en error, print den og gå videre til neste produkt
			fmt.Println(res.Error())
			continue
		}

		res.GetProductData(apiProducts)
	}
}

func (data Response) GetProductData(apiProducts *model.ApiProducts) {
	// Definerer en funksjon som går gjennom base noden
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// BaseUrl er satt til en empty string fordi linker til produktet har
			// allerede hele linken
			product := model.ApiProduct{Store: "bunnpris", BaseUrl: ""}

			// Mapper over alle attributter elementet har
			// Om attr sin value er products-container, kjøres en ny funksjon
			// på alle child elements
			for _, attr := range node.Attr {
				// Switch statement som sjekker verdien til attributten
				// Om den har values som passer til elementer med data vi vil
				// ha, lagres dataen i product (instansen av Product)
				switch attr.Val {
				case "titleName":
					product.Data.Title = node.FirstChild.Data
					fmt.Println(node.Attr[2].Val)
				case "productImage":
					product.Data.ImageLink = node.FirstChild.Attr[5].Val
				case "lblName":
					product.Data.Slug = node.Parent.Attr[1].Val
					// Henter gtin fra linken (henter fra itemno search paramen
					// ved å splitte 2 ganger)
					product.Data.Ean = strings.Split(strings.Split(product.Data.Slug, "itemno=")[1], "&")[0]
					fmt.Println(product.Data.Title, product.Data.Ean)
				case "ContentPlaceHolder1_ucItemGroupProduct_rptItemGroupProd_lblSalesPrice_0":
					for _, attr := range node.Attr {
						if attr.Key == "data-dnprice" {
							price := attr.Val
							priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
							if err != nil {
								fmt.Printf("Couldnt convert string to float: %s", err.Error())
								continue
							}
							product.Data.OriginalPrice = priceFloat
						}
					}
				case "ContentPlaceHolder1_ucItemGroupProduct_rptItemGroupProd_lblCampaignPrice_0":
					for _, attr := range node.Attr {
						if attr.Key == "data-dnprice" {
							price := attr.Val
							if price == "" {
								product.Data.Price = product.Data.OriginalPrice
								continue
							}
							priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
							if err != nil {
								fmt.Printf("Couldnt convert string to float: %s", err.Error())
								continue
							}
							product.Data.Price = priceFloat
						}
					}
				}
			}
			*apiProducts = append(*apiProducts, product)
		}

		// går gjennom html-noden og så går til sibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	// kjører crawler på base-noden
	crawler(data.Data.HTML)
}

func (data Response) OLDGETPRODUCTS(apiProducts *model.ApiProducts) {
	// Definerer en funksjon som går gjennom base noden
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// BaseUrl er satt til en empty string fordi linker til produktet har
			// allerede hele linken
			product := model.ApiProduct{Store: "bunnpris", BaseUrl: ""}

			// Mapper over alle attributter elementet har
			// Om attr sin value er products-container, kjøres en ny funksjon
			// på alle child elements
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "products-container" {
					// Definerer ny funksjon som henter data om produktet
					var traverseChildren func(*html.Node)
					traverseChildren = func(child *html.Node) {
						if child.Type == html.ElementNode {
							for _, attr := range child.Attr {
								// Switch statement som sjekker verdien til attributten
								// Om den har values som passer til elementer med data vi vil
								// ha, lagres dataen i product (instansen av Product)
								switch attr.Val {
								case "productImage":
									product.Data.ImageLink = child.FirstChild.Attr[5].Val
								case "lblName":
									product.Data.Title = child.FirstChild.Data
									product.Data.Slug = child.Parent.Attr[1].Val
									// Henter gtin fra linken (henter fra itemno search paramen
									// ved å splitte 2 ganger)
									product.Data.Ean = strings.Split(strings.Split(product.Data.Slug, "itemno=")[1], "&")[0]
									fmt.Println(product.Data.Title, product.Data.Ean)
								case "ContentPlaceHolder1_ucItemGroupProduct_rptItemGroupProd_lblSalesPrice_0":
									for _, attr := range child.Attr {
										if attr.Key == "data-dnprice" {
											price := attr.Val
											priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
											if err != nil {
												fmt.Printf("Couldnt convert string to float: %s", err.Error())
												continue
											}
											product.Data.OriginalPrice = priceFloat
										}
									}
								case "ContentPlaceHolder1_ucItemGroupProduct_rptItemGroupProd_lblCampaignPrice_0":
									for _, attr := range child.Attr {
										if attr.Key == "data-dnprice" {
											price := attr.Val
											if price == "" {
												product.Data.Price = product.Data.OriginalPrice
												continue
											}
											priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
											if err != nil {
												fmt.Printf("Couldnt convert string to float: %s", err.Error())
												continue
											}
											product.Data.Price = priceFloat
										}
									}
								}
							}
						}

						// Kjører funksjonen for hvert child for å gå gjennom hele tre-et
						for c := child.FirstChild; c != nil; c = c.NextSibling {
							traverseChildren(c)
						}
					}

					// Kjører traverseChildren når en produkt er funnet
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						traverseChildren(child)
					}

					fmt.Printf("%+v", product)

					// Legger til produktet som har blitt funnet i products arrayet
					*apiProducts = append(*apiProducts, product)
				}
			}
		}

		// går gjennom html-noden og så går til sibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	// kjører crawler på base-noden
	crawler(data.Data.HTML)
}
