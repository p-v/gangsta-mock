FROM golang:1.14.4-alpine

RUN apk add --no-cache git gcc musl-dev

WORKDIR $GOPATH/src/gangsta-mock

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build .
RUN go build -buildmode=plugin -o plugins/callback.so plugins/callback.go
RUN go build -buildmode=plugin -o plugins/handler.so plugins/handler.go

EXPOSE 8080

CMD ["./gangsta-mock", "server", "start"]
