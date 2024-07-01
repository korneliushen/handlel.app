FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV BUNNPRIS_TOKEN=u1bqkpp1au5zelihvzrnhp1s

RUN CGO_ENABLED=0 GOOS=linux go build -o /scraper

CMD ["/scraper"]
