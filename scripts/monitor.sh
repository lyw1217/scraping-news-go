#!/bin/bash

APP="scraping-news-go"
DIR_PATH="/home/ubuntu/Documents/github/${APP}"

GOBIN=$DIR_PATH/bin

LOG_PATH="${DIR_PATH}/log"
LOG_NAME="nohup.log"

MOD="scraping-news"
EXE="${GOBIN}/${MOD}"
CMD="GOCRAPER"
CMD_GO="/usr/local/go/bin/go"

export GOBIN

WAIT_TIME=7

# sudo check
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; exit; fi

echo ""
echo " --------------------------------------"
echo "          [ GOCRAPER MONITOR ]         "
echo "                                with go"
echo " --------------------------------------"
echo ""

echo " > 현재 구동중인 애플리케이션 pid 확인"

CURRENT_PID=$(pgrep -f ${CMD})

echo "   pid: ${CURRENT_PID}"
echo ""

if [ -z "${CURRENT_PID}" ]; then
    echo " > 현재 구동중인 애플리케이션이 없음"
	echo ""
else
	echo " > 정상 실행중..!"
	echo " > 모니터링 종료"
	exit 0
fi

exec ${DIR_PATH}/scripts/start.sh
exit 0
