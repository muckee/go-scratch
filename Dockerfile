ARG GO_VERSION=1.20
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache git \
                       ca-certificates
 
# Add user here. Cannot be added in scratch
RUN [ "sh", "-c", "addgroup -S goserver && adduser -S -u 10000 -g goserver goserver" ]

# Install Go modules
WORKDIR /src
COPY ./go.mod ./
RUN go mod download

# Copy static files from repository
COPY ./ ./
 
# TODO: Run tests
 
# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/app
 
# STAGE 2: build the container to run
FROM scratch AS final

# Copy Go executable
COPY --chown=goserver:goserver --from=build /src/cmd/app /app
 
# Copy CA certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 
# Copy users
COPY --from=build /etc/passwd /etc/passwd

# Create user
USER goserver

 # Run the executable
ENTRYPOINT ["/app"]
