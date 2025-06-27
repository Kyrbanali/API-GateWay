FROM golang:1.24

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

CMD ["/app/app"]