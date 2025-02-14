FROM golang:1.23-alpine3.19 AS build

# Define current working directory
WORKDIR /url-shortener

# Download modules to local cache so we can skip re-
# downloading on consecutive docker build commands
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add sources
COPY . .

RUN go build -o ./cmd/server/url-shortener ./cmd/server


FROM alpine:3.19

WORKDIR /url-shortener

COPY --from=build /url-shortener/cmd/server/url-shortener ./cmd/server/url-shortener

COPY --from=build /url-shortener/config.json ./config.json

EXPOSE 8080
ENTRYPOINT ["./cmd/server/url-shortener"]
