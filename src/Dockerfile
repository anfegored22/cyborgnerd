# Simple multi-stage Dockerfile for the Go app
# Build stage
FROM golang:1.22-alpine AS build
WORKDIR /src

# Cache go dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary (embeds templates/static via go:embed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app .

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=build /out/app /app/app

# The app listens on port 8080 per src/main.go
EXPOSE 8080

CMD ["/app/app"]

