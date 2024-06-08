package bunnpris

import (
	"context"
	"fmt"
	"time"
)

type Categories []Category

type Category struct {
	Name          string
	Link          string
	SubCategories []string
}

func (c *Categories) Get(ctx context.Context, token string) error {
	// Implementerer en timeout, etter 10 sekunder vil funksjonen time out
	// og returnere en error (da er det enten noe feil med api-en, eller
	// internettet er så dårlig at ingenting vil funke)
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Lager en channel som holder en value med type Response, som vil bli
	// returnert fra api kallet
	reschan := make(chan Response)

	// Kjører requesten i en go func, så vi kan sjekke om den er timed out
	// eller returnert en value
	go func() {
		// gjør en post request til /Itemgroups.aspx
		res := POST(token, "Itemgroups.aspx")
		// setter dataen vi får fra api i channelen vi lagde
		reschan <- Response{
			Message:    res.Message,
			StatusCode: res.StatusCode,
			Data:       res.Data,
		}
	}()

	// Lager en variabel som responsen fra api-en blir lagret i inni for loopen
	res := Response{}

	// Variabel som sjekker om respons har kommet fra databasen
	// Om done blir gjort om til true, vil for loopen breake og res vil ha en
	// verdi
	done := false

	// Sjekker om funksjonen har tima ut, om den har det returneres en error
	// Om den ikke har tima ut, får res verdien til reschan og done blir true
	// og funksjonen breaker
	for {
		if done {
			break
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("Function timed out fetching from bunnpris api")
		case res = <-reschan:
			done = true
		}
	}

	// returnerer en error om ingen html ble returnert fra post requesten
	if res.Data.HTML == nil {
		return fmt.Errorf("No html returned")
	}

	// parser html-en som returneres fra post requesten
	categories := res.ParseHTML()

	// parsehtml returnerer en verdi med type Categories, så vi assigner
	// categories (c) funksjonen ble kjørt på til categories vi får fra
	// parseHTML
	*c = categories

	return nil
}
