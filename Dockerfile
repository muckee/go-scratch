ARG GO_VERSION=1.20
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

# Prepare the environment
RUN apk add --no-cache coreutils \
                       git \
                       ca-certificates
 
# Add 'goserver' user
RUN addgroup -S goserver && adduser -S -u 10000 -g goserver goserver

# Set the current working directory
WORKDIR /src

# Install Go dependencies
COPY ./go.mod ./
RUN go mod download

# Copy static files from repository into the build image
COPY ./ ./
 
# TODO: Run tests
 
# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/app

# Set permissions for the executable
RUN chmod +x /app && \
    chown goserver:goserver /app
 
# STAGE 2: build the container to run
FROM scratch AS final


# DEV: Copy 'ls' from the build image into the final image
COPY --from=build /bin/ls /ls
 
# Copy CA certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 
# Copy users
COPY --from=build /etc/passwd /etc/passwd

# Copy Go executable
COPY --from=build /src/cmd/app /app

# Create user
USER goserver

 # Run the executable
ENTRYPOINT ["/app"]
