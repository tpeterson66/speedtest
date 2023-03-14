FROM golang:alpine3.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go .

RUN go build -o /speedtest

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /speedtest /speedtest

USER nonroot:nonroot

ENTRYPOINT ["/speedtest"]
