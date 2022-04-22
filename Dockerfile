FROM golang:1.18-alpine

WORKDIR $GOPATH/src
COPY . .

ENV GO111MODULE=auto

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main -v cmd/main.go

ENTRYPOINT /go/src/main
