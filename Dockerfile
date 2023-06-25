ARG GO_VERSION=1.20
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache git \
                       ca-certificates
 
# Add user here. Cannot be added in scratch
RUN addgroup -S myapp \
    && adduser -S -u 10000 -g myapp myapp

# Install Go modules
WORKDIR /src
COPY ./go.mod ./
RUN go mod download

COPY ./ ./
 
# Run tests
# RUN CGO_ENABLED=0 go test -timeout 30s -v github.com/gbaeke/go-template/pkg/api
 
# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/app
 
# STAGE 2: build the container to run
FROM scratch AS final

COPY --from=build /app /app
 
# copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 
# copy users from builder (use from=0 for illustration purposes)
COPY --from=0 /etc/passwd /etc/passwd
 
USER myapp
 
ENTRYPOINT ["/app"]
