package bunnpris

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
)

const BASE_URL = "https://nettbutikk.bunnpris.no"

// TODO: trenger bare noen ting her og ikke alt (conditional det var den ene videoen jeg husker det)

func POST(ctx context.Context, token, endpoint string, reqBody io.Reader, contentType string) Response {
	// lager en url som requests skal sendes til ved å kombinere base url og
	// endpoint vi får som arg
	apiUrl := BASE_URL + endpoint
	// gjør klar requesten med NewRequest som tar inn method, url og body.
	// body trengs ikke så er satt til nil
	req, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		return Response{Message: "Error preparing request: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

	// legger til Content-Type: text/html; charset=iso-8859-1
	// det var det postman brukte når den fikk valid responses
	req.Header.Add("Content-Type", contentType)

	// lager jar med cookies
	// jar har en SetCookies funksjon som tar inn en url med type *url.URL,
	// så parser url-en fra string
	// jar tar inn et array med *http.Cookie som andre argument, som har Name
	// og Value
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return Response{Message: "Error preparing cookies: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

	// parser url-en til *url.URL
	parsedUrl, err := url.Parse(BASE_URL)
	if err != nil {
		return Response{Message: "Error parsing url: " + err.Error(),
			StatusCode: http.StatusBadRequest}
	}

	// lager et array med cookies med en cookie som er session id og legger til
	// cookiesa i jar
	cookies := []*http.Cookie{{Name: "ASP.NET_SessionId", Value: token}}
	jar.SetCookies(parsedUrl, cookies)

	// lager en http client med jar så cookies blir sendt med requesten
	// og kjører requesten
	client := &http.Client{Jar: jar}

	// gjør en post request til /Itemgroups.aspx
	res, err := client.Do(req)
	if err != nil {
		return Response{Message: "Error during request: " + err.Error(), StatusCode: 500}
	}
	// kan kjøre defer close nå som res har en verdi
	defer res.Body.Close()

	// leser body responsen og lagrer i body variabel, dette blir brukt
	// om status koden ikke er 200
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{Message: "Error reading body: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

	// om statusen ikke er OK, returneres body til responsen ApiRes.message
	if res.StatusCode != http.StatusOK {
		return Response{Message: "Received error from API: " + string(body),
			StatusCode: res.StatusCode}
	}

	// responses er i node, så parser bodyen til *node.Node
	node, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return Response{Message: "Error parsing html: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

	// returnerer dataen
	return Response{Message: "Success",
		StatusCode: http.StatusOK,
		Data:       ResponseData{HTML: node, JSON: body}}
}
