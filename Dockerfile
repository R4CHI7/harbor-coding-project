FROM golang:1.21

COPY go.mod go.sum /go/src/
WORKDIR /go/src

# RUN go mod download
COPY . /go/src/

RUN go build -o calendly .

EXPOSE 8080 8080
ENTRYPOINT ["./calendly"]