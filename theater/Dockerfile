FROM golang:1.24-alpine

WORKDIR /app

# later go.sum
COPY go.mod ./
RUN go mod download

COPY . .
RUN rm -f theater 
RUN rm -f main
RUN CGO_ENABLED=0 GOOS=linux go build -o theater .

EXPOSE 8080
CMD ["./run.sh"]
