ARG BUILD_IMAGE=golang:1.17-alpine3.15
ARG BASE_IMAGE=alpine:3.15

FROM ${BUILD_IMAGE} AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN ["mkdir", "-p", "/build"]

WORKDIR /build

COPY . .

RUN ["go", "mod", "vendor"]

RUN ["go", "build", "-o", "/tmp/scraper", "main.go"]
RUN ["chmod", "+x", "/tmp/scraper"]

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)

COPY --chown=0:0 --from=builder /build /build
COPY --chown=0:0 --from=builder /tmp/scraper /build/

WORKDIR /build

RUN ["mkdir", "-p", "log"]
CMD ["./scraper"]