FROM golang:1.14.4-alpine

RUN apk add --no-cache git gcc g++

WORKDIR $GOPATH/src/gangsta-mock

COPY . .

RUN go get gangsta-mock
RUN go build .
RUN go build -buildmode=plugin -o plugins/generic.so plugins/generic.go

EXPOSE 8080

CMD ["./gangsta-mock"]
