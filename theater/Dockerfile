FROM golang:1.24-alpine

WORKDIR /app

# later go.sum
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o theater .
RUN go test -v ./...

EXPOSE 8080
USER nonroot:nonroot
CMD ["./theater"]
