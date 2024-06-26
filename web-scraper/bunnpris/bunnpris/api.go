package bunnpris

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

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
	HTML *html.Node
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

// Legger til en method til Response struct som parser html og returnerer
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

func (productLinks BunnprisProducts) FetchProductPages(ctx context.Context, token string, products *model.Products) {
	// Flere threads ellers tar det sånn 1 time å kjøre
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, link := range productLinks {
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

			if err := res.GetProductData(products, link); err != nil {
				fmt.Printf("Error getting product data from link %s: %v\n", link, err)
			}
		}(link)
	}

	wg.Wait()
}

func (data Response) GetProductData(products *model.Products, link string) error {
	// Lager en instanse av BaseProduct som data legges til i når det blir funnet
	product := &model.Product{
		Store:   "bunnpris",
		BaseUrl: BASE_URL,	
    // Legger til et notat som sier at produktet er fra bunnpris og ikke har 
    // kategori. Dette hjelper om noen skal inn i databasen og for queries.
    Notes: "ingen_kategori_bunnpris",
    // Setter kategory og underkategori til empty string. De ble automatisk til
    // personlige artikler før av en eller annen grunn.
    Category: "",
    SubCategory: "",
	}

  // Initializer variabler for pris, så de ikke blir reassigna inni findData
  // hele tiden. Prisene som blir funnet blir lagt til i en instans av
  // model.Price og appenda til product.Prices
  price := model.Price{Url: BASE_URL + link, Store: "bunnpris"}
  nutritionalContent := model.NutritionalContent{}

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
        findData(node, attr, product, &price, &nutritionalContent)
			}
		}

		// går gjennom html-noden og så går til sibling
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	// kjører crawler på base-noden
	crawler(data.Data.HTML)

  product.Prices = append(product.Prices, price)
  product.NutritionalContent = &nutritionalContent

  baseImageLink := BASE_URL + product.Images.Medium
	product.Images.Small = strings.Replace(baseImageLink, "_m", "_s", 1)
	product.Images.Large = strings.Replace(baseImageLink, "_m", "_f", 1)
  // BaseImageLink er medium (_m) versjon av bildet, så bare assigner ImageLinkLarge
  // til baseImageLink
	product.Images.Medium = baseImageLink


	if product.Prices[0].Price == 0 {
		product.Prices[0].Price = product.Prices[0].OriginalPrice
	}
	if product.Prices[0].UnitPrice == 0 {
		product.Prices[0].UnitPrice = product.Prices[0].OriginalPrice
	}

	// Om originalPrice og price ikke er like er det salg så
	// onSale settes til true
	if product.Prices[0].Price != product.Prices[0].Price {
		product.OnSale = true
	}

	// Legger til link som slug
	product.Slug = sanitizeData(link)

	// Gjør noen ekstra checks for å populate fields i databasen
	*products = append(*products, *product)

	return nil
}

func findData(node *html.Node, attr html.Attribute, product *model.Product, price *model.Price, nutritionalContent *model.NutritionalContent) {
  // Instanse av Price

  switch attr.Val {
  case "form1":
    for _, attr := range node.Attr {
      if attr.Key == "action" {
        action := strings.Split(attr.Val, "itemno=")
        if len(action) == 0 {
          fmt.Println("NO DATA EAN FOUND FOR PRODUCT")
          return
        }
        product.Id = sanitizeData(strings.Split(action[1], "&")[0])
      }
    }

  case "titleName":
    for _, attr := range node.Attr {
      if attr.Key == "title" {
        product.Title = sanitizeData(attr.Val)
        fmt.Println("Getting data for", product.Title)
      }
    }

  // Til nå har jeg funnet to ulike måter bilder kan vises på
  case "zoomLens":
    for _, attr := range node.FirstChild.Attr {
      if attr.Key == "src" {
        product.Images.Medium = sanitizeData(attr.Val)
      }
    }
  case "itemDetailImg":
    for _, attr := range node.Attr {
      if attr.Key == "src" {
        product.Images.Medium = sanitizeData(attr.Val)
      }
    }

  case "item-cotext7":
    if node.FirstChild != nil {
      if node.FirstChild.FirstChild != nil {
        product.SubTitle = sanitizeData(node.FirstChild.FirstChild.Data)
      }
    }

  case "lblItemDesc":
    if node.FirstChild != nil {
      if node.NextSibling != nil {
        product.Description = sanitizeData(node.FirstChild.NextSibling.Data)
      }
    }

  // Henter prisen. Henter først originalPrice. Etter det sjekkes
  // price (campaign price). Om det ikke eksisterer, blir price satt til
  // originalPrice
  case "lblSalesPrice":
    for _, attr := range node.Attr {
      if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
        priceString := attr.Val
        priceFloat, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
        if err != nil {
          fmt.Printf("Couldnt convert string to float: %s", err.Error())
          continue
        }
        price.OriginalPrice = priceFloat
      }
    }
  case "lblCampaignPrice":
    for _, attr := range node.Attr {
      if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
        priceString := attr.Val
        priceFloat, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
        if err != nil {
          fmt.Printf("Couldnt convert string to float: %s", err.Error())
          continue
        }
        price.Price = priceFloat
      }
    }

  // Henter unittype og unitprice
  case "clearfix div-Weightitem":
    unitPriceNode := node.FirstChild.NextSibling
    for _, attr := range unitPriceNode.Attr {
      if attr.Key == "data-dnprice" || attr.Key == "data-dnprice-decimal" {
        priceString := attr.Val
        priceFloat, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
        if err != nil {
          fmt.Printf("Couldnt convert string to float: %s", err.Error())
          continue
        }
        price.UnitPrice = priceFloat
      }
    }
    product.UnitType = sanitizeData(unitPriceNode.NextSibling.Data)

  // Næringsinnhold
  case "Nutrition":
    // jeg beklager denne koden eksisterer
    tableBody := node.FirstChild.NextSibling.NextSibling.FirstChild

    // Mapper over alle children til tablebody, og legger til verdier i 
    // nutritionalContent basert op key valuen i html-elementet
    for child := tableBody.FirstChild; child != nil; child = child.NextSibling {
	    v := reflect.ValueOf(nutritionalContent).Elem()

      key := strings.TrimSpace(child.FirstChild.FirstChild.Data)
      // På bunnpris heter det karbohydrat, så må gjøre det til Karbohydrater
      if key == "Karbohydrat" {
        key = "Karbohydrater"
      }
      field := v.FieldByName(key)

      if field.IsValid() {
        value := strings.TrimSpace(child.FirstChild.NextSibling.FirstChild.Data)
        field.SetString(value)
      }
		}
  }
}

// Function to clean non-UTF-8 data
func sanitizeData(str string) string {
	if utf8.ValidString(str) {
		return str
	}
	return strings.ToValidUTF8(str, "")
}
