FROM golang:1.23-alpine3.19 AS build

WORKDIR /url-shortener

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go test -coverprofile=coverage.out ./... && \
    go tool cover -func=coverage.out
RUN go build -o ./cmd/server/url-shortener ./cmd/server


FROM alpine:3.19

WORKDIR /url-shortener

COPY --from=build /url-shortener/cmd/server/url-shortener ./cmd/server/url-shortener

COPY --from=build /url-shortener/config.json ./config.json

EXPOSE 8080
CMD ["sh", "-c", "./cmd/server/url-shortener >> /var/log/url-shortener/url-shortener.log 2>&1"]