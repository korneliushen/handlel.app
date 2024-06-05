package api

import (
	"fmt"
	"net/http"

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
type ParsedHTML struct {
	href  string
	title string
}

// TODO: kan bare kjøre om Data.HTML eksisterer?
// TODO: gjøre det dynamic (generics kanskje) sånn at man kan velge
// hva man vil returnere (element, type osv.) (nå er det bare href og title)
// legger til en method til Response struct som parser html og returnerer
// alle instanser av et element
func (data Response) ParseHTML() []ParsedHTML {
	var foundItems []ParsedHTML

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		fmt.Println(node.Data)
		if node.Type == html.ElementNode &&
			len(node.Attr) > 0 &&
			node.Attr[0].Val == "itemgroup-link" {
			href := node.Attr[1].Val
			title := node.FirstChild.Data
			foundItems = append(foundItems, ParsedHTML{href, title})
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(data.Data.HTML)

	return foundItems
}
