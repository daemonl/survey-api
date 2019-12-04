FROM golang:1.13

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app .

CMD ["/app"]
