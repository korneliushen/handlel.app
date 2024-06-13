package bunnpris

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

func (products BunnprisProducts) FetchProductPages(ctx context.Context, token string, apiProducts *model.ApiProducts) {
	// Flere threads ellers tar det sånn 1 time å kjøre
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, link := range products {
		wg.Add(1)
		sem <- struct{}{}

		go func(link string) {
			defer wg.Done()
			defer func() { <-sem }()

			res := POST(ctx, token, link, nil, "text/html; charset=us-ascii")
			if res.IsError() {
				// Om det er en error, print den og gå videre til neste produkt
				fmt.Printf("Error getting data from link %s: %v\n", link, res)
			}

			res.GetProductData(apiProducts, link)
		}(link)
	}

	wg.Wait()
}

// TODO FOR MEG IMRGN: alt av logikken for at dette skal kjøre smooth er ferdig
// så må bare legge til dataen fra produkt siden til ApiProduct instansen

func (data Response) GetProductData(apiProducts *model.ApiProducts, link string) {
	// Lager en instanse av ApiProduct som data legges til i når det blir funnet
	product := model.ApiProduct{Store: "bunnpris", BaseUrl: "https://nettbutikk.bunnpris.no"}
	// Definerer en funksjon som går gjennom base noden
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// Mapper over alle attributter elementet har
			// Om attr sin value er products-container, kjøres en ny funksjon
			// på alle child elements
			for _, attr := range node.Attr {
				// Switch statement som sjekker verdien til attributten
				// Om den har values som passer til elementer med data vi vil
				// ha, lagres dataen i product (instansen av Product)
				switch attr.Val {
				case "form1":
					for _, attr := range node.Attr {
						if attr.Key == "action" {
							product.Data.Ean = strings.Split(strings.Split(attr.Val, "itemno=")[1], "&")[0]
						}
					}
				case "titleName":
					for _, attr := range node.Attr {
						if attr.Key == "title" {
							product.Data.Title = attr.Val
							fmt.Println("Getting data for", product.Data.Title)
						}
					}

				// Henter prisen. Henter først originalPrice. Etter det sjekkes
				// price (campaign price). Om det ikke eksisterer, blir price satt til
				// originalPrice
				case "lblSalesPrice":
					for _, attr := range node.Attr {
						if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
							price := attr.Val
							priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
							if err != nil {
								fmt.Printf("Couldnt convert string to float: %s", err.Error())
								continue
							}
							product.Data.OriginalPrice = priceFloat
						}
					}
				case "lblCampaignPrice":
					for _, attr := range node.Attr {
						if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
							price := attr.Val
							priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
							if err != nil {
								fmt.Printf("Couldnt convert string to float: %s", err.Error())
								continue
							}
							product.Data.Price = priceFloat
						}
					}

				// Henter unittype og unitprice
				case "clearfix div-Weightitem":
					unitPrice := node.FirstChild.NextSibling
					for _, attr := range unitPrice.Attr {
						if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
							price := attr.Val
							priceFloat, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
							if err != nil {
								fmt.Printf("Couldnt convert string to float: %s", err.Error())
								continue
							}
							product.Data.ComparePricePerUnit = priceFloat
						}
					}
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

	if product.Data.Price == 0 {
		product.Data.Price = product.Data.OriginalPrice
	}
	if product.Data.ComparePricePerUnit == 0 {
		product.Data.ComparePricePerUnit = product.Data.OriginalPrice
	}

  // Legger til link som slug
  product.Data.Slug = link

	// Gjør noen ekstra checks for å populate fields i databasen
	*apiProducts = append(*apiProducts, product)
}
