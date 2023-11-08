FROM golang:latest

ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build -o link-shortener ./cmd/link-shortener/main.go

CMD ["./link-shortener"]