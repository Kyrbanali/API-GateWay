FROM golang:1.24 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM alpine:3.21

COPY --from=builder /app/app /

ENTRYPOINT ["/app"]