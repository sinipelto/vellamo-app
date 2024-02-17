FROM golang:1.22.0-alpine

WORKDIR /app

RUN go build -o './vellamo' .

ENTRYPOINT [ "./vellamo" ]
