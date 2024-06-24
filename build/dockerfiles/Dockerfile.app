FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /build

# Copy and download dependency using go mod
ADD ./src/go.* /build/
RUN go mod download

# Copy sources to build container
ADD ./src /build/

# Build the app
RUN go build -a -tags musl -o /build/app
######################################
FROM alpine:3
LABEL AUTHOR="Lena Ramisch (Linuxcode)"

# install curl for healthcheck
RUN apk --no-cache add curl

# Essentials
RUN apk add -U tzdata
ENV TZ=Europe/Berlin
RUN cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime

USER nobody
COPY --from=builder --chown=nobody /build/app /custom/app
ENTRYPOINT [ "/custom/app" ]
