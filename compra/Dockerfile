FROM golang:1.15

WORKDIR /go/src

COPY . .

RUN GOOS=linux go build compra.go

EXPOSE 9090

ENTRYPOINT [ "./compra" ]