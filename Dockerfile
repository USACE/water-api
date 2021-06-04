FROM golang:1.16.4-alpine AS builder
# Install Git
RUN apk update && apk add --no-cache git
# Copy In Source Code
WORKDIR /go/src/app
COPY . .

# Install Dependencies
RUN go get -d -v
# Build
RUN go get -d -v \
  && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o /go/bin/water-api

# SCRATCH IMAGE
FROM scratch
COPY --from=builder /go/bin/water-api /go/bin/water-api
ENTRYPOINT ["/go/bin/water-api"]