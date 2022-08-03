package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"scraping-news/config"
	"scraping-news/util"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func (p *PageInfo_t) AddLinks(t string, u string) []LinkInfo_t {
	// https://stackoverflow.com/questions/18042439/go-append-to-slice-in-struct
	// https://stackoverflow.com/questions/34329441/golang-struct-array-values-not-appending-in-loop
	p.Links = append(p.Links, LinkInfo_t{
		Title: t,
		Url:   u,
	})
	return p.Links
}

// url로 HTTP GET 요청하여 http.Response 객체 반환
func requestGetDocument(url string) (*http.Response, error) {
	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err, "Err, Failed to NewRequest()")
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0") // 안티 크롤링 회피
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Err, Failed to Get Request")
		return nil, err
	}

	return resp, err
}

// scrap maekyung M.S.G on the date as parameters
func GetMaekyungMSG(d_month int, d_day int) (int, string) {
	list_url := MkMSGUrl

	p := &PageInfo_t{
		Links:    make([]LinkInfo_t, 0, 10),
		Contents: make([]string, 0, 10),
	}

	// 1. 매세지 첫 page 목록 조회
	resp, err := requestGetDocument(list_url)
	if err != nil {
		log.Error(err, "Err, Failed to Get Request")
		return resp.StatusCode, err.Error()
	}
	defer resp.Body.Close()

	if p.StatusCode = resp.StatusCode; p.StatusCode == 200 {
		// HTML Read
		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Error(err, "Err. Failed to NewDocumentFromReader()")
			return p.StatusCode, err.Error()
		}

		// 파싱
		container := html.Find("dt.tit") // returns a new Selection ojbect
		//container := html.Find("div.list_area2")

		// 하이퍼링크 순회 및 저장
		container.Each(func(i int, s *goquery.Selection) {
			//title := util.ConvEuckrToUtf8(s.Find("a").Text())
			title := util.TransEuckrToUtf8(s.Find("a").Text())
			href, ok := s.Find("a").Attr("href")
			if !ok {
				log.Info(ok, "Err. No Exist href in", p.Links[i].Title)
				return
			}
			//fmt.Printf("title = %s, href = %s\n", title, href)
			p.AddLinks(title, href)
		})

		// 날짜에 맞는 article 확인
		for _, lnk := range p.Links {
			if strings.Contains(lnk.Title, strconv.Itoa(d_month)) &&
				strings.Contains(lnk.Title, strconv.Itoa(d_day)) {

				resp_link, err := requestGetDocument(lnk.Url)
				if err != nil {
					log.Error(err, "Err, Failed to Get Request")
					break
				}
				defer resp_link.Body.Close()

				if p.StatusCode = resp_link.StatusCode; p.StatusCode == 200 {
					html, err := goquery.NewDocumentFromReader(resp_link.Body)
					if err != nil {
						log.Error(err, "Err. Failed to NewDocumentFromReader()")
						break
					}

					container := html.Find("#content > div.content_left > div.view_txt")

					// view_txt 순회하면서 문자열에 추가
					container.Each(func(i int, s *goquery.Selection) {

						//content := util.ConvEuckrToUtf8(s.Text())
						content := util.TransEuckrToUtf8(s.Text())
						if content == "" {
							log.Error("Err. Failed to convert content : ", s.Text())
							return
						}

						p.Contents = append(p.Contents, content)
					})

					if len(p.Contents) == 0 {
						log.Error("Err. Failed to get Contents")
						return p.StatusCode, "Err. Failed to get Contents"
					}

					return p.StatusCode, strings.Join(p.Contents, "")

				} else {
					log.Error("Err. Failed to get M.S.G Article.")
					return p.StatusCode, "Err. Failed to get M.S.G Article."
				}
			}
		}

		return p.StatusCode, fmt.Sprintf("No article on %d-%d", d_month, d_day)
	}

	log.Error("Err. Failed to get the M.S.G list.")
	return resp.StatusCode, err.Error()
}

// scrap hankyung issue today on the date as a parameters
func GetHankyungIssueToday(d_month int, d_day int) (int, string) {
	list_url := HkIssueTodayUrl

	p := &PageInfo_t{Contents: make([]string, 0, 10)}

	// 1. Issue Today 조회
	resp, err := requestGetDocument(list_url)
	if err != nil {
		log.Error(err, "Err, Failed to Get Request")
		return resp.StatusCode, err.Error()
	}
	defer resp.Body.Close()

	if p.StatusCode = resp.StatusCode; p.StatusCode == 200 {
		// HTML Read
		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Error(err, "Err. Failed to NewDocumentFromReader()")
			return p.StatusCode, err.Error()
		}

		// 파싱
		//container := html.Find("td.stb-text-box")
		container := html.Find("div.stb-text-box")

		// text-box 순회하면서 문자열에 추가
		container.Each(func(i int, s *goquery.Selection) {
			//content := util.ConvEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))
			content := util.TransEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))
			if content == "" {
				log.Error("Err. Failed to convert content : ", s.Text())
				return
			}

			// https://stackoverflow.com/questions/65533097/replace-nbsp-or-0xao-with-space-in-a-string
			//content := TransEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))
			/*
				if !(strings.Contains(content, "카카오톡으로 공유하세요")) {
					p.Contents = append(p.Contents, content)
				}
			*/
			p.Contents = append(p.Contents, content)
		})

		if len(p.Contents) == 0 {
			log.Error("Err. Failed to get Contents")
			return p.StatusCode, "Err. Failed to get Contents"
		}
		t_date := strings.Split(p.Contents[0], ".")
		if len(t_date) >= 3 {
			t_year, _ := strconv.Atoi(strings.TrimSpace(t_date[0]))
			t_month, _ := strconv.Atoi(strings.TrimSpace(t_date[1]))
			t_day, _ := strconv.Atoi(strings.TrimSpace(t_date[2]))

			if d_month == t_month && d_day == t_day {
				return p.StatusCode, strings.Join(p.Contents, "\r\n")
			} else {
				return p.StatusCode, fmt.Sprintf("No article on %d-%d-%d", t_year, d_month, d_day)
			}
		}
	}

	log.Error("Err. Failed to get the Issue Today.")
	return resp.StatusCode, err.Error()
}

// scrap Quick News on the date as parameters
func GetQuickNews(d_month int, d_day int) (int, string) {
	list_url := QuicknewsUrl

	p := &PageInfo_t{Contents: make([]string, 0, 10)}

	// 1. Issue Today 조회
	resp, err := requestGetDocument(list_url)
	if err != nil {
		log.Error(err, "Err, Failed to Get Request")
		return resp.StatusCode, err.Error()
	}
	defer resp.Body.Close()

	if p.StatusCode = resp.StatusCode; p.StatusCode == 200 {
		// HTML Read
		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Error(err, "Err. Failed to NewDocumentFromReader()")
			return p.StatusCode, err.Error()
		}

		// 파싱
		container := html.Find("pre#news_0")

		contents := strings.Split(strings.TrimSpace(container.Text()), "\n")

		if strings.Contains(contents[0], strconv.Itoa(d_month)+"월") && strings.Contains(contents[0], strconv.Itoa(d_day)+"일") {
			// 주요 경제지표 및 코인가격 삭제
			var i int
			for i = range contents {
				if strings.Contains(contents[i], "----") {
					break
				}
				p.Contents = append(p.Contents, contents[i])
			}
			return p.StatusCode, strings.Join(p.Contents, "\r\n")
		} else {
			return p.StatusCode, fmt.Sprintf("No article on %d-%d", d_month, d_day)
		}
	}

	log.Error("Err. Failed to get the Quick News.")
	return resp.StatusCode, err.Error()
}

var sysClose bool

// start scraping
func StartScraping() {
	log.Error("< SCRAPER > Start Scraping Routine - Started ......")

	c := config.Config

	for !sysClose {
		d_month := int(time.Now().Month())
		d_day := int(time.Now().Day())
		d_hour := int(time.Now().Hour())
		d_min := int(time.Now().Minute())

		for i, media := range c.Media {
			if c.SendHour == d_hour && c.SendMin <= d_min {
				if media.Flag {
					switch media.Name {
					// 한국 경제
					case "hankyung":
						StatusCode, contents := GetHankyungIssueToday(d_month, d_day)
						if StatusCode != 200 {
							log.Error("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.SendMessageToSlack("한국경제 Issue Today", contents); err != nil {
							log.Error("Err. slack.SendMessageToSlack")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.KakaoSendToMe(media.Name, contents, HostName+media.Name); err != nil {
							log.Error("Err. KakaoSendToMe")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Info(contents)
						c.Media[i].Flag = false

					// 매일 경제
					case "maekyung":
						StatusCode, contents := GetMaekyungMSG(d_month, d_day)
						if StatusCode != 200 {
							log.Error("Err. news.GetMaekyungMSG, StatusCode :", StatusCode)
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.SendMessageToSlack("매일경제 매.세.지", contents); err != nil {
							log.Error("Err. slack.SendMessageToSlack")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.KakaoSendToMe(media.Name, contents, HostName+media.Name); err != nil {
							log.Error("Err. KakaoSendToMe")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Info(contents)
						c.Media[i].Flag = false

					// 간추린 아침뉴스
					case "quicknews":
						StatusCode, contents := GetQuickNews(d_month, d_day)
						if StatusCode != 200 {
							log.Error("Err. news.GetQuickNews, StatusCode :", StatusCode)
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.SendMessageToSlack("간추린 아침뉴스", contents); err != nil {
							log.Error("Err. slack.SendMessageToSlack")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.KakaoSendToMe(media.Name, contents, HostName+media.Name); err != nil {
							log.Error("Err. KakaoSendToMe")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Info(contents)
						c.Media[i].Flag = false
					default:
						log.Error("Err. Wrong Key")
					}
				}
			} else if c.SendHour != d_hour {
				config.ResetConfig(&c.Media[i])
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	log.Error("< SCRAPER > Exit StartScraping Routine ...")
}

var VilageInfo []VilageInfo_t = make([]VilageInfo_t, 0)

func LoadVilageInfo(fileName string, sheetName string) error {
	var v VilageInfo_t

	path, _ := filepath.Abs(fileName)
	f, err := excelize.OpenFile(path)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, row := range rows {
		for i, colCell := range row {
			switch i {
			case 2: // 1단계
				v.Step1 = colCell
			case 3: // 2단계
				v.Step2 = colCell
			case 4: // 3단계
				v.Step3 = colCell
			case 5: // 격자 X
				v.X = colCell
			case 6: // 격자 Y
				v.Y = colCell
			default:
				continue
			}
		}
		VilageInfo = append(VilageInfo, v)
	}

	return err
}

// 도시 정보 키워드 순차 탐색
func KeywordVilageSearch(k string, v []VilageInfo_t) []VilageInfo_t {
	result := make([]VilageInfo_t, 0)
	for _, val := range v {
		if strings.Contains(val.Step3, k) || strings.Contains(val.Step2, k) || strings.Contains(val.Step1, k) {
			result = append(result, val)
		}
	}
	return result
}

// TODO 탐색 속도 개선 필요, 정확한 도시 탐색 개선 필요
func SearchVilage(keyword []string) []VilageInfo_t {
	result := make([]VilageInfo_t, 0)

	if len(keyword) > 1 {
		for i, k := range keyword {
			// 첫 번째 전체 탐색, 두 번째부터는 결과 내 탐색
			decodedKey, err := url.QueryUnescape(k)
			if err != nil {
				log.Error(err, "Err. Failed to QueryUnescape.")
				return nil
			}

			if i == 0 {
				// 키워드 글자수 2개 이상인 경우 앞의 두 글자만으로 탐색 (서현동 같은 경우 서현1동/서현2동 으로 구분되어 탐색하지 못함)

				if len(decodedKey) > 6 {
					result = KeywordVilageSearch(decodedKey[:6], VilageInfo)
				} else {
					result = KeywordVilageSearch(decodedKey, VilageInfo)
				}
			} else {
				if len(decodedKey) > 6 {
					result = KeywordVilageSearch(decodedKey[:6], result)
				} else {
					result = KeywordVilageSearch(decodedKey, result)
				}
			}
		}
	} else {
		decodedKey, err := url.QueryUnescape(keyword[0])
		if err != nil {
			log.Error(err, "Err. Failed to QueryUnescape.")
			return nil
		}

		if len(decodedKey) > 6 {
			result = KeywordVilageSearch(decodedKey[:6], VilageInfo)
		} else {
			result = KeywordVilageSearch(decodedKey, VilageInfo)
		}
	}

	return result
}

func requestFcstApi(url string, r ReqVilageFcst_t) (*ResVilageFcst_t, error) {
	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err, "Err, Failed to NewRequest()")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("numOfRows", r.NumOfRows)
	q.Add("pageNo", r.PageNo)
	q.Add("dataType", r.DataType)
	q.Add("base_date", r.Base_date)
	q.Add("base_time", r.Base_time)
	q.Add("nx", r.Nx)
	q.Add("ny", r.Ny)
	//q.Add("serviceKey", r.ServiceKey)

	// '%'가 포함된 문자열을 q.Encode() 시 '%' 문자열이 escape되어서 값이 달라짐
	req.URL.RawQuery = q.Encode() + fmt.Sprintf("&serviceKey=%s", r.ServiceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Err, Failed to Get Request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "Err, Failed to ReadAll")
		return nil, err
	}

	parse_resp := ResVilageFcst_t{}
	err = json.Unmarshal(body, &parse_resp)
	if err != nil {
		log.Error("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Error("syntax error at byte offset %d", e.Offset)
		}
		log.Error("response: %q", body)
		return nil, err
	}

	parse_resp.Response.Body.Name = r.Name

	return &parse_resp, err
}

/*
- Base_time : 0200, 0500, 0800, 1100, 1400, 1700, 2000, 2300 (1일 8회)
- API 제공 시간(~이후) : 02:10, 05:10, 08:10, 11:10, 14:10, 17:10, 20:10, 23:10
*/
func calcBaseDateTime(r *ReqVilageFcst_t) {
	baseTime := []int{2, 5, 8, 11, 14, 17, 20, 23}
	now := time.Now()
	now_year := now.Year()
	now_month := now.Month()
	now_day := now.Day()
	now_hour := time.Now().Hour()
	now_minute := time.Now().Minute()

	// 02시 10분 이전, 전날 마지막 정보 조회
	if now_hour <= baseTime[0] {
		if now_minute <= 10 {
			then := now.AddDate(0, 0, -1)
			r.Base_date = fmt.Sprintf("%04d%02d%02d", then.Year(), then.Month(), int(then.Day()))
			r.Base_time = fmt.Sprintf("%02d00", baseTime[len(baseTime)-1])
			return
		}
	}

	// 02시 10분 이후
	for i, t := range baseTime {
		if now_hour <= t {
			if now_minute < 10 {
				// 제공시간 전이므로 이전 시간 정보 조회
				r.Base_date = fmt.Sprintf("%04d%02d%02d", now_year, now_month, now_day)
				r.Base_time = fmt.Sprintf("%02d00", baseTime[i-2])
				return
			} else {
				// 제공시간 이후
				r.Base_date = fmt.Sprintf("%04d%02d%02d", now_year, now_month, now_day)
				r.Base_time = fmt.Sprintf("%02d00", baseTime[i-1])
				return
			}
		}
	}

	// default (혹시나)
	r.Base_date = fmt.Sprintf("%04d%02d%02d", now_year, now_month, now_day)
	r.Base_time = "0200"
}

/*
	POP	강수확률		%
	PTY	강수형태		코드값
	PCP	1시간 강수량	범주 (1 mm)
	REH	습도			%
	SNO	1시간 신적설	범주(1 cm)
	SKY	하늘상태		코드값
	TMP	1시간 기온		℃
	TMN	일 최저기온		℃
	TMX	일 최고기온		℃
	UUU	풍속(동서성분)	m/s
	VVV	풍속(남북성분)	m/s
	WAV	파고			M
	VEC	풍향			deg
	WSD	풍속			m/s
*/
func ParseCode(category string, value string) (string, string) {

	switch category {
	case "POP": // 강수 확률
		cat := "강수 확률"
		return cat, fmt.Sprintf("%s %%", value)
	case "PTY": // 강수 형태
		cat := "강수 형태"
		switch value {
		case "0":
			return cat, "없음"
		case "1":
			return cat, "비"
		case "2":
			return cat, "비/눈"
		case "3":
			return cat, "눈"
		case "4":
			return cat, "소나기"
		default:
			return cat, "UNKNOWN_VALUE"
		}
	case "PCP": // 강수량
		cat := "강수량"
		if value == "-" || value == "null" || value == "0" || value == "" {
			return cat, "강수없음"
		}
		return cat, value
	case "REH": // 습도
		cat := "습도"
		return cat, fmt.Sprintf("%s%%", value)
	case "SNO": // 1시간 신적설
		cat := "1시간 신적설"
		if value == "-" || value == "null" || value == "0" || value == "" {
			return cat, "강수없음"
		}
		return cat, value
	case "SKY": // 하늘상태
		cat := "하늘 상태"
		switch value {
		case "1":
			return cat, "맑음"
		case "2":
			return cat, "구름조금"
		case "3":
			return cat, "구름많음"
		case "4":
			return cat, "흐림"
		default:
			return cat, "UNKNOWN_VALUE"
		}
	case "TMP": // 1시간 기온
		cat := "1시간 기온"
		return cat, fmt.Sprintf("%s℃", value)
	case "TMN": // 일 최저 기온
		cat := "일 최저 기온"
		return cat, fmt.Sprintf("%s℃", value)
	case "TMX": // 일 최고 기온
		cat := "일 최고 기온"
		return cat, fmt.Sprintf("%s℃", value)
	case "UUU": // 동서바람성분
		cat := "동서바람성분"
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			log.Error(err, "Err. ParseFloat, value = ", value)
			return cat, "VALUE_ERR"
		}
		if v < 0.0 {
			return cat, fmt.Sprintf("서풍%sm/s", value)
		} else {
			return cat, fmt.Sprintf("동풍%sm/s", value)
		}
	case "VVV": // 남북바람성분
		cat := "남북바람성분"
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			log.Error(err, "Err. ParseFloat, value = ", value)
			return cat, "VALUE_ERR"
		}
		if v < 0.0 {
			return cat, fmt.Sprintf("남풍%sm/s", value)
		} else {
			return cat, fmt.Sprintf("북풍%sm/s", value)
		}
	case "WAV": // 파고
		cat := "파고"
		return cat, fmt.Sprintf("%s미터", value)
	case "VEC": // 풍향
		cat := "풍향"
		return cat, fmt.Sprintf("%sdeg", value)
	case "WSD": // 풍속
		cat := "풍속"
		return cat, fmt.Sprintf("%sm/s", value)
	default:
		log.Error("Unknown Category. cat =", category, "val =", value)
		return "UNKNOWN_CATEGORY", "UNKNOWN_VALUE"
	}
}

func GetVilageFcstInfo(keyword []string) (*ResVilageFcst_t, error) {
	// 키워드 검색
	v := SearchVilage(keyword)
	if len(v) <= 0 {
		fmt.Printf("키워드 %s 검색 결과 없음.\n", keyword)
		return nil, nil
	}
	//fmt.Println(v)

	fmt.Printf("검색 결과 총 %d개의 %s가 검색되었습니다.\n", len(v), keyword)

	var r ReqVilageFcst_t
	r.ServiceKey = config.Keys.Fcst.Encoding_key
	r.NumOfRows = "360" // default : 10 | 12개 항목당 1시간
	r.PageNo = "1"      // default : 1
	r.DataType = "JSON"
	r.Nx = v[0].X
	r.Ny = v[0].Y
	r.Name = fmt.Sprintf("%s %s %s", v[0].Step1, v[0].Step2, v[0].Step3)
	calcBaseDateTime(&r)

	resp, err := requestFcstApi(VilageFcstUrl, r)
	if err != nil {
		log.Error(err, "Failed to requestFcstApi")
	}

	if resp.Response.Header.ResultCode != "00" {
		log.Error(err, "ResultCode(", resp.Response.Header.ResultCode, ") is not Normal. ")
	}

	return resp, err
}

func init() {
	// https://www.data.go.kr/data/15084084/openapi.do , 기상청_단기예보 ((구)_동네예보) 조회서비스
	fileName := "./spec/기상청41_단기예보 조회서비스_오픈API활용가이드_최종/기상청41_단기예보 조회서비스_오픈API활용가이드_격자_위경도(20220103).xlsx"
	sheetName := "최종 업데이트 파일_20220103"
	err := LoadVilageInfo(fileName, sheetName)
	if err != nil {
		log.Error(err)
	}
}
