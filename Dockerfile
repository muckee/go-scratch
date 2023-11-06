ARG GO_VERSION=1.20
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

# Prepare the environment
RUN apk add --no-cache git \
                       ca-certificates

 # Create a group with a specific GID
RUN addgroup goserver --gid 10000

# Create a user with a specific UID and assign them to the group
RUN adduser goserver --uid 10000 --ingroup goserver --home /home/goserver --disabled-password --gecos ""

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
RUN chmod 770 /app && \
    chown 10000:goserver /app
 
# STAGE 2: build the container to run
FROM scratch AS final
 
# Copy CA certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 
# Copy users
COPY --from=build /etc/passwd /etc/passwd

# Copy Go executable
COPY --from=build /app /app

# Create user
USER goserver

 # Run the executable
ENTRYPOINT ["/app"]
