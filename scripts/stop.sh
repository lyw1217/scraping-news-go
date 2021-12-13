#!/bin/bash

APP="scraping-news-go"
DIR_PATH="/home/leeyw/Documents/github/${APP}"

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
echo "----------------------------------"
echo "         [ SCRAPER STOP ]         "
echo "                           with go"
echo "----------------------------------"
echo ""

sleep 0.5

echo "> 현재 구동중인 애플리케이션 pid 확인"
echo ""

CURRENT_PID=$(pgrep -f ${CMD})

echo "  pid: $CURRENT_PID"
if [ -z "$CURRENT_PID" ]; then
    echo "> 현재 구동중인 애플리케이션이 없으므로 종료하지 않습니다."
else
    echo "> kill -15 $CURRENT_PID"
    kill -15 $CURRENT_PID
	sleep 0.5
	
	for cnt in $(seq ${WAIT_TIME})
	do
		CURRENT_PID=$(pgrep -f ${CMD})
		if [ -z "$CURRENT_PID" ]; then
			echo "> 어플리케이션 종료 성공!"
			echo ""
			break
		fi
		sleep 0.5
		
		if [ ${cnt} == ${WAIT_TIME} ]; then
			echo "> 어플리케이션 종료 실패.. 다시 시도하세요"
			echo ""
			exit 1
		fi
	done
fi

echo ""
echo "----------------------------------"
echo "    [ SCRAPER STOP COMPLETE ]     "
echo "----------------------------------"
echo ""

echo ""
exit 0
