FROM golang:1.15.8-alpine3.13

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/lucky
COPY . .

RUN apk -U upgrade && apk --no-cache add ca-certificates git gcc
RUN go mod download
RUN go build -o deck_api cmd/api/deck.go

CMD ["./deck_api"]

EXPOSE 3000