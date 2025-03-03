package bunnpris

import (
	"encoding/json"
	"os"
	"time"
)

type Token struct {
  Store   string `json:"store"`
	Created int64  `json:"created"`
	Expiry  int64  `json:"expiry"`
	Value   string `json:"value"`
}

// Leser enten fra en lokal valid token, eller lager en ny en og lagrer den
// lokalt. Tokens kan være invalid ved å være expired (4-6 timer tror jeg)
func ReadToken() (Token, error) {
	// Sjekker først om det eksisterer en fil som heter token.json.
	// Om det ikke gjør det, lages en med en token
	if _, err := os.Stat("./token.json"); err != nil {
		if _, err := NewToken(); err != nil {
			return Token{}, err
		}
	}

	// Leser token.json filen
	token, err := os.ReadFile("./token.json")
	if err != nil {
		return Token{}, err
	}

	// Unmarshaler tokenen fra json filen inn i en variabel med type Token
	var newToken Token
	if err = json.Unmarshal(token, &newToken); err != nil {
		return newToken, err
	}

	// Sjekker om tokenen er valid
  // Om den er invalid, lages en ny token
	if !newToken.Valid() {
		if newToken, err = NewToken(); err != nil {
			return newToken, err
		}

	}

	return newToken, nil
}

func NewToken() (Token, error) {
	// Setter expiry til å være 5 timer fra da tokenen ble laget
	newToken := Token{
    Store:   "bunnpris",
		Created: time.Now().Unix(),
		Expiry:  time.Now().Add(5 * time.Hour).Unix(),
		Value:   os.Getenv("BUNNPRIS_TOKEN"),
	}

	byteToken, err := json.MarshalIndent(newToken, "", "  ")
	if err != nil {
		return Token{}, err
	}

	if err = os.WriteFile("./token.json", byteToken, 0666); err != nil {
		return Token{}, err
	}

	return newToken, nil
}

// Sjekker om tokenen er valid
func (token *Token) Valid() bool {
	// Sjekker om det er en expiration satt på tokenen
	if token.Created == 0 {
		return false
	}

	if token.Value == "" {
		return false
	}

	// En if condition for å sjekke og returnere tidlig om tokenen er invalid
	if token.Expiry < time.Now().Unix() {
		return false
	}

	return true
}
