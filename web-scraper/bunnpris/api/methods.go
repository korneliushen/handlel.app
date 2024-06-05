package api

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
)

// TODO: bruke riktige status codes

const BASE_URL = "https://nettbutikk.bunnpris.no/"

func POST(token, endpoint string) Response {
	apiUrl := BASE_URL + endpoint
	// gjør klar requesten med NewRequest som tar inn method, url og body.
	// body trengs ikke så er satt til nil
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return Response{Message: "Error preparing request: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

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
	parsedUrl, err := url.Parse(apiUrl)
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
	res, err := client.Do(req)
	if err != nil {
		return Response{Message: "Error executing request: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}
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

	// responses er i html, så parser bodyen til *html.Node
	html, err := html.Parse(res.Body)
	if err != nil {
		return Response{Message: "Error parsing html: " + err.Error(),
			StatusCode: http.StatusInternalServerError}
	}

	// returnerer dataen
	return Response{Message: "Success",
		StatusCode: http.StatusOK,
		Data:       ResponseData{HTML: html}}
}
