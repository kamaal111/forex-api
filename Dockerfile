# App builder
FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /go/src/github.com/kamaal111/forex-api/

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download -x && go mod verify

# Copy source and build
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-w -s" -trimpath -v -o /go/bin/forex-api .

# Build a smaller image with the minimum required things to run.
FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/forex-api /go/bin/forex-api
# Run the forex-api binary.
EXPOSE 8000
ENTRYPOINT ["/go/bin/forex-api"]
