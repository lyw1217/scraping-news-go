#!/bin/bash

APP="scraping-news-go"
DIR_PATH="/home/leeyw/Documents/github/${APP}"

GOPATH=$DIR_PATH
GOBIN=$DIR_PATH/bin

LOG_PATH="${DIR_PATH}/log"
LOG_NAME="nohup.log"

MOD="scraping"
EXE="${GOBIN}/${MOD}"
CMD="GOCRAPER"

export GOPATH
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

echo " 현재 구동중인 어플리케이션 pid: ${CURRENT_PID}"

if [ -z "${CURRENT_PID}" ]; then
    echo " > 현재 구동중인 애플리케이션이 없음"
	echo ""
else
	echo " > 정상 실행중..!"
	echo " > 모니터링 종료"
	exit 0
fi

exec ${DIR_PATH}/start
exit 0