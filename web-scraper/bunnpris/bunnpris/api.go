package bunnpris

import (
	"fmt"
	"net/http"
	"net/url"

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
	JSON interface{}
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
func (data Response) ParseHTML() Categories {
	var categories Categories

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// sjekker om node-en er en ElementNode
		if node.Type == html.ElementNode {
			// init variabler som brukes til å holde dataen som blir funnet
			var className, href, name string

			// mapper over alle attributter elementet har
			// om attr er itemgroup link og har en href, blir linken til href-en
			// lagret som href
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "itemgroup-link" {
					className = attr.Val
				}
				if attr.Key == "href" {
					// bruker queryescape for å kunne bruke linken i en url
					href = url.QueryEscape(attr.Val)
				}
			}

			// om classname og href ble funnet, hentes dataen fra span-en som er
			// child element (må kjøre FirstChild.NextSibling fordi det første er
			// ::before). må også kjøre Firstchild etter nextsibling for å få teksten
			// og ikke bare navnet på html-elementet
			if className != "" && href != "" {
				if node.FirstChild != nil {
					name = node.FirstChild.NextSibling.FirstChild.Data
				}
				categories = append(categories, Category{Name: name, Link: href})
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
