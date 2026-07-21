FROM golang:1.25.3 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o scraper ./cmd/worker/

FROM alpine:3.21

RUN addgroup -S scraper && adduser -S scraper -G scraper

WORKDIR /app

COPY --from=build /build/scraper ./scraper
COPY --from=build /build/config /config

USER scraper

ENTRYPOINT ["./scraper"]
