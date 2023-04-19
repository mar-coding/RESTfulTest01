FROM golang:1.18.1-alpine AS builder

LABEL maintainer = "marcoding78@gmail.com"

RUN apk --no-cache add ca-certificates git
WORKDIR /build

# Fetch dependencies
COPY movieCatcherApp/go.mod movieCatcherApp/go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build

# Create final image
FROM alpine
WORKDIR /
COPY --from=builder /build/myapp .
EXPOSE 8080