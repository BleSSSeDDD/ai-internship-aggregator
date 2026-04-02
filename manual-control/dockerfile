FROM golang:1.25.3 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o manual-control ./cmd/server/

FROM alpine:latest AS manual-control

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone

RUN addgroup -S manual && adduser -S manual -G manual

WORKDIR /app

COPY --from=build /build/manual-control ./manual-control

COPY --from=build /build/static ./static
COPY --from=build /build/templates ./templates

RUN chmod +x ./manual-control

RUN chown -R manual:manual /app

USER manual

ENV KAFKA_BROKERS=kafka:9092
ENV PORT=2228

EXPOSE 2228

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:2228/health || exit 1

CMD ["./manual-control"]