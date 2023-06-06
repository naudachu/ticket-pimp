FROM golang:1.20-alpine

WORKDIR $GOPATH/src/ticket-creator/

COPY . .

RUN go mod download
RUN go mod verify

RUN cd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/app cmd/main.go

CMD ["app"]