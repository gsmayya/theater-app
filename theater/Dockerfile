FROM golang:1.24-alpine

# Install MySQL client for database initialization
RUN apk add --no-cache mysql-client

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN rm -f theater 
RUN rm -f main
RUN CGO_ENABLED=0 GOOS=linux go build -o theater .

# Make scripts executable
RUN chmod +x scripts/*.sh
RUN chmod +x run.sh

EXPOSE 8080
CMD ["./run.sh"]
