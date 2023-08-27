ARG GO_VERSION=1.20
ARG USER_NAME=goserver
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache git \
                       ca-certificates
 
# Add user here. Cannot be added in scratch
RUN addgroup -S ${USER_NAME} \
    && adduser -S -u 10000 -g ${USER_NAME} ${USER_NAME}

# Install Go modules
WORKDIR /src
COPY ./go.mod ./
RUN go mod download

COPY ./ ./
 
Run tests
RUN CGO_ENABLED=0 go test -timeout 30s -v github.com/gbaeke/go-template/pkg/api
 
# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/app
 
# STAGE 2: build the container to run
FROM scratch AS final

COPY --from=build /src/cmd/app /app
 
# copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
 
copy users
COPY --from=build /etc/passwd /etc/passwd

USER ${USER_NAME}
 
ENTRYPOINT ["/app"]
