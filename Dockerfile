FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum conf.yaml ./
RUN go mod download
COPY . .
RUN ls -l

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /dps
EXPOSE 8080

# Run
CMD ["/dps"]