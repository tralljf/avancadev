FROM golang:1.15

WORKDIR /go/src

COPY . .

RUN GOOS=linux go build pagamento.go

ENTRYPOINT [ "./pagamento" ]