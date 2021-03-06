ARG BUILD_IMAGE=golang:1.17-alpine3.15
ARG BASE_IMAGE=alpine:3.15

FROM ${BUILD_IMAGE} AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=debug \
    PORT=30200

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/scraper .

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GIN_MODE=debug \
    PORT=30200

WORKDIR /usr/src/app

COPY --chown=0:0 --from=builder /usr/src/app /usr/src/app
COPY --chown=0:0 --from=builder /usr/local/bin/scraper /usr/local/bin/scraper

RUN mkdir -p log

EXPOSE ${PORT}

CMD ["scraper"]