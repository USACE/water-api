#############
# API Docs
#############
FROM node:19 AS apidoc-builder

WORKDIR /app
COPY apidoc.yml ./apidoc.yml
RUN npm install -yg redoc-cli \
    && redoc-cli build apidoc.yml

#############
# API
#############
FROM golang:1.19.3-alpine AS builder
# Install Git
RUN apk update && apk add --no-cache git ca-certificates
# Copy In Source Code
WORKDIR /go/src/app
COPY api .
# Copy sql into subdir for flyaway
COPY sql ./sql

# Install Dependencies
RUN go get -d -v
# Build
RUN go get -d -v \
  && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o /go/bin/water-api

#############
# SCRATCH
#############
FROM scratch
COPY --from=builder /go/bin/water-api /go/bin/water-api
COPY --from=builder /go/src/app/sql /sql/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
                    /etc/ssl/certs/ca-certificates.crt
COPY --from=apidoc-builder /app/redoc-static.html /apidoc.html
COPY --from=apidoc-builder /app/apidoc.yml /apidoc.yml
VOLUME [ "/sql" ]
ENTRYPOINT ["/go/bin/water-api"]