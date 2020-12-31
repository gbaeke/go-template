# argument for Go version
ARG GO_VERSION=1.15.6

# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build

# git required for go mod
RUN apk add --no-cache git

# certs
RUN apk --no-cache add ca-certificates

# add a user here because addgroup and adduser are not available in scratch
RUN addgroup -S myapp \
    && adduser -S -u 10000 -g myapp myapp

# Working directory will be created if it does not exist
WORKDIR /src

# We use go modules; copy go.mod and go.sum
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import code
COPY ./ ./

# Run tests
RUN CGO_ENABLED=0 go test -timeout 30s -v github.com/gbaeke/go-template/pkg/api

# Build the executable
RUN CGO_ENABLED=0 go build \
	-installsuffix 'static' \
	-o /app ./cmd/app

# STAGE 2: build the container to run
FROM scratch AS final

# add maintainer label
LABEL maintainer="gbaeke"

# copy compiled app
COPY --from=build /app /app

# copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy users from builder which contains myapp user
COPY --from=0 /etc/passwd /etc/passwd

USER myapp

# run binary
ENTRYPOINT ["/app"]