FROM golang:1.25.3 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o ai-parser ./cmd/worker/

RUN CGO_ENABLED=0 GOOS=linux go build -o wait-for-model ./cmd/wait-for-model.go

FROM alpine:latest AS parser

RUN addgroup -S parser && adduser -S parser -G parser

WORKDIR /parser

COPY --from=build /build/ai-parser ./ai-parser
COPY --from=build /build/config /config
COPY --from=build /build/wait-for-model ./wait-for-model

RUN chmod +x ./ai-parser ./wait-for-model

USER parser

ENTRYPOINT [ "./wait-for-model"]