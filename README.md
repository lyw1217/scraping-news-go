# scraping-news-go

go 기반 뉴스 스크래핑

## 웹 스크래핑(Web Scraping)

웹에 존재하는 데이터를 추출하는 기술

## 개발 환경

- CentOS Linux release 7.9.2009 (Core)
- go version go1.17.3 linux/amd64

## 기능 설명

- 매일경제의 `매경이 전하는 세상의 지식(매-세-지)` [링크](https://www.mk.co.kr/premium/series/20007/)
- 한국경제의 `한경 Issue Today` [링크](https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb=)
...

매일 업데이트 되는 연재 기사들을 스크래핑해서 slack 채널에 메시지를 전송한다.
