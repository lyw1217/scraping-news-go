ARG ARCH=
ARG BUILD_IMAGE=${ARCH}golang:1.18-alpine3.15
ARG BASE_IMAGE=${ARCH}alpine:3.15

FROM ${BUILD_IMAGE} AS builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /usr/local/bin/scraper .

WORKDIR /usr/src
RUN tar -cvf app.tar ./app

FROM ${BASE_IMAGE}
LABEL AUTHOR Youngwoo Lee (mvl100d@gmail.com)

ENV GIN_MODE=debug \
	SCRAP_HOME=/usr/src/app \
    PORT=30200 \
    TZ=Asia/Seoul
    
RUN apk --no-cache add tzdata && \
	cp /usr/share/zoneinfo/$TZ /etc/localtime && \
	echo $TZ > /etc/timezone \
	apk del tzdata

WORKDIR /usr/src

COPY --chown=0:0 --from=builder /usr/src/app.tar		/usr/src
COPY --chown=0:0 --from=builder /usr/src/app/run.sh		/usr/src
COPY --chown=0:0 --from=builder /usr/local/bin/scraper	/usr/local/bin/scraper

EXPOSE ${PORT}

ENTRYPOINT ["/usr/src/run.sh"]
