# BUILD GO BINARIES
FROM golang:1.19-alpine AS builder
# Install Git
RUN apk update && apk add --no-cache git ca-certificates
# Copy In Source Code
WORKDIR /go/src/app
COPY ./listener .
COPY go.mod go.sum ./
# Install Dependencies
RUN go mod download && \
    go mod verify
# Build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o /go/bin/water-listener

# SCRATCH IMAGE
FROM scratch
COPY --from=builder /go/bin/water-listener /go/bin/water-listener
COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
                    /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/go/bin/water-listener"]
