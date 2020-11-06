FROM golang:1.15

WORKDIR /go/src

COPY . .

RUN GOOS=linux go build coupom.go

EXPOSE 9092

ENTRYPOINT [ "./coupom" ]