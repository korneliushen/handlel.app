## INFO OM PROGRAMMET / GUIDE OM NY BUTIKK SKAL LEGGES TIL

NOTE: om nye produkter blir lagt til blir det litt goofy med kategorier, så må se hva vi skal gjøre med det

### krav

Du trenger en .env fil i root directory for å kjøre programmet

Den trenger en NEON_URL (neon connection string, må være pooled connection) og en ALGOLIA_SECRET (algolia read og write api key)

### hvordan koden fungerer

main.go har en variabel som heter apiProducts som er et array med typen ApiProducts.

I main.go hentes data fra butikkene (for nå ngdata og bunnpris), hvor produkter appendes til apiProducts arrayet.

Når alle produkter er hentet, kjøres en formatering til et format databasen vil ha, og så sendes dataen til databasen og til algolia indexen.


### legge til data fra nye butikker

Lag et nytt directory her og gi det navnet på butikken du scraper.

Sett opp scraperen hvordan du vil, men lag en funksjon som kan importes og kjøres av main.go

Denne funksjonen må ta inn et argument med type *model.ApiProducts. Produktene du får fra scraperen din skal appendes til dette arrayet.

Kjør funksjonen i main.go og gi apiProducts som argument.

Om dataen du henter er formatert som model.ApiProducts, skal det funke av seg selv og dataen skal legges til.


#### Fil struktur
```js
~/web-scraper/
├── algolia // funksjoner som interacter med algolia
├── neon // funksjoner som interacter med neon
├── lib // dir for funksjoner som skal brukes flere steder og har ganske basic funksjonalitet ig
├── model // dir for felles structs (produkter og kategorier)
├── ngdata/bunnpris/oda/osv. // scrapere av ulike butikker
├── .gitignore // no leaking
├── Dockerfile // docker (funker)
├── README.md // det du leser rn
├── go.mod
├── go.sum
└── main.go // entry
```
