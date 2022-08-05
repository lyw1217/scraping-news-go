package scraper

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"scraping-news/config"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

var VilageInfo []VilageInfo_t = make([]VilageInfo_t, 0)
var RssVilageInfo []FcstZone_t = make([]FcstZone_t, 0)

func LoadFcstZoneCode(fileName string) error {
	var v []FcstZone_t

	path, _ := filepath.Abs(fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		log.Error(err, "Err. Failed to Open file", path)
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&v)
	if err != nil {
		log.Error(err, "Err. Failed Decode Json")
		return err
	}

	RssVilageInfo = append(RssVilageInfo, v...)

	return nil
}

func LoadVilageInfo(fileName string, sheetName string) error {
	var v VilageInfo_t

	path, _ := filepath.Abs(fileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
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
func KeywordVilageSearch(k string, v []VilageInfo_t, k_count int) []VilageInfo_t {
	result := make([]VilageInfo_t, 0)
	for _, val := range v {
		if strings.HasPrefix(val.Step3, k) || strings.HasPrefix(val.Step2, k) || strings.HasPrefix(val.Step1, k) {
			result = append(result, val)
		}
	}

	// 전체 글자 탐색 실패한 경우 앞의 두 글자만으로 탐색 (서현동 같은 경우 서현1동/서현2동 으로 구분되어 탐색하지 못함)
	if len(result) == 0 && k_count <= 1 {
		k_sub := k[:6] // utf8 한글 한 글자당 3개
		fmt.Printf("excel 전체(%s) 탐색 실패, (%s)로 재탐색\n", k, k_sub)
		for _, val := range v {
			if strings.HasPrefix(val.Step3, k_sub) || strings.HasPrefix(val.Step2, k_sub) || strings.HasPrefix(val.Step1, k_sub) {
				result = append(result, val)
			}
		}
	}

	return result
}

func KeywordVilageSearchRSS(k string, v []FcstZone_t, k_count int) []VilageInfo_t {

	result := make([]VilageInfo_t, 0)
	for _, val := range v {
		if strings.HasPrefix(val.Step3, k) || strings.HasPrefix(val.Step2, k) || strings.HasPrefix(val.Step1, k) {
			result = append(result,
				VilageInfo_t{
					Step1: val.Step1,
					Step2: val.Step2,
					Step3: val.Step3,
					X:     val.X,
					Y:     val.Y,
				},
			)
		}
	}

	// 전체 글자 탐색 실패한 경우 앞의 두 글자만으로 탐색 (서현동 같은 경우 서현1동/서현2동 으로 구분되어 탐색하지 못함)
	if len(result) == 0 && k_count <= 1 {
		k_sub := k[:6] // utf8 한글 한 글자당 3개
		fmt.Printf("rss 전체(%s) 탐색 실패, (%s)로 재탐색\n", k, k_sub)
		for _, val := range v {
			if strings.HasPrefix(val.Step3, k_sub) || strings.HasPrefix(val.Step2, k_sub) || strings.HasPrefix(val.Step1, k_sub) {
				result = append(result,
					VilageInfo_t{
						Step1: val.Step1,
						Step2: val.Step2,
						Step3: val.Step3,
						X:     val.X,
						Y:     val.Y,
					},
				)
			}
		}
	}

	return result
}

// TODO 탐색 속도 개선 필요, 정확한 도시 탐색 개선 필요
/* keyword : string | flag : 0(RSS Search), 1(Excel Search) */
func SearchVilage(keyword []string, flag int) []VilageInfo_t {
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
				if flag == 0 {
					result = KeywordVilageSearchRSS(decodedKey, RssVilageInfo, len(keyword))
				} else {
					result = KeywordVilageSearch(decodedKey, VilageInfo, len(keyword))
				}
			} else {
				if flag == 0 {
					result = KeywordVilageSearchRSS(decodedKey, RssVilageInfo, len(keyword))
				} else {
					result = KeywordVilageSearch(decodedKey, VilageInfo, len(keyword))
				}
			}
		}
	} else {
		decodedKey, err := url.QueryUnescape(keyword[0])
		if err != nil {
			log.Error(err, "Err. Failed to QueryUnescape.")
			return nil
		}

		if len(decodedKey) >= 6 {
			if flag == 0 {
				result = KeywordVilageSearchRSS(decodedKey, RssVilageInfo, len(keyword))
			} else {
				result = KeywordVilageSearch(decodedKey, VilageInfo, len(keyword))
			}
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
	var v []VilageInfo_t
	v = SearchVilage(keyword, 0)
	fmt.Printf("RSS 검색 결과 총 %d개의 %s가 검색되었습니다.\n", len(v), keyword)

	if len(v) <= 0 {
		// RSS에 없는 경우, 단기예보 엑셀에서 다시 조회
		v = SearchVilage(keyword, 1)
		if len(v) <= 0 {
			fmt.Printf("키워드 %s 검색 결과 없음.\n", keyword)
			return nil, nil
		}
		fmt.Printf("Excel 검색 결과 총 %d개의 %s가 검색되었습니다.\n", len(v), keyword)
	}

	var r ReqVilageFcst_t

	if len(v[0].X) > 0 {
		r.Nx = v[0].X
		r.Ny = v[0].Y
	} else {
		// RSS에 X,Y좌표가 없는 경우, 단기예보 엑셀에서 다시 조회
		v = SearchVilage(keyword, 1)
		if len(v) <= 0 {
			fmt.Printf("키워드 %s 검색 결과 없음.\n", keyword)
			return nil, nil
		}

		fmt.Printf("Excel 검색 결과 총 %d개의 %s가 검색되었습니다.\n", len(v), keyword)

		r.Nx = v[0].X
		r.Ny = v[0].Y
	}
	r.ServiceKey = config.Keys.Fcst.Encoding_key
	r.NumOfRows = "120" // default : 10 | 12개 항목당 1시간
	r.PageNo = "1"      // default : 1
	r.DataType = "JSON"
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

func getEachZoneCode(u string, pList FcstZone_t, step int) ([]FcstZone_t, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Error(err, "Err, Failed to NewRequest()")
		return nil, err
	}

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

	parse_resp := []FcstZone_t{}
	err = json.Unmarshal(body, &parse_resp)
	if err != nil {
		log.Error("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Error("syntax error at byte offset %d", e.Offset)
		}
		log.Error("response: %q", body)
		return nil, err
	}

	switch step {
	case 1:
		for i := range parse_resp {
			parse_resp[i].Step1 = parse_resp[i].Name
		}
	case 2:
		for i := range parse_resp {
			if len(pList.Step1) > 0 {
				parse_resp[i].Step1 = pList.Step1
			}
			parse_resp[i].Step2 = parse_resp[i].Name
		}
	case 3:
		for i := range parse_resp {
			if len(pList.Step1) > 0 {
				parse_resp[i].Step1 = pList.Step1
			}
			if len(pList.Step2) > 0 {
				parse_resp[i].Step2 = pList.Step2
			}
			parse_resp[i].Step3 = parse_resp[i].Name
		}
	}

	return parse_resp, nil
}

func GetZoneCode() error {
	fileName := "./spec/FcstZoneCode.json"
	err := LoadFcstZoneCode(fileName)
	if err != nil {
		log.Error(err)
	}

	// TOP ZONE
	topList, err := getEachZoneCode(TopURL, FcstZone_t{}, 1)
	if err != nil {
		log.Error(err, "Err. Failed to getEachZoneCode")
		return err
	}

	// MID ZONE
	mdlList := make([]FcstZone_t, 0)
	for _, v := range topList {
		qry_url := fmt.Sprintf(MdlURL, v.Code)
		tmpList, err := getEachZoneCode(qry_url, v, 2)
		if err != nil {
			log.Error(err, "Err. Failed to getEachZoneCode")
			return err
		}
		mdlList = append(mdlList, tmpList...)
		time.Sleep(time.Millisecond * 2)
	}

	// LEAF ZONE
	leadList := make([]FcstZone_t, 0)
	for _, v := range mdlList {
		qry_url := fmt.Sprintf(LeafURL, v.Code)
		tmpList, err := getEachZoneCode(qry_url, v, 3)
		if err != nil {
			log.Error(err, "Err. Failed to getEachZoneCode")
			return err
		}
		leadList = append(leadList, tmpList...)
		time.Sleep(time.Millisecond * 2)
	}

	topList = append(topList, mdlList...)
	topList = append(topList, leadList...)

	f, err := os.Create("./spec/FcstZoneCode.json")
	if err != nil {
		log.Error(err, "Err. Failed to os.Create")
		return err
	}
	defer f.Close()
	data, err := json.Marshal(topList)
	if err != nil {
		log.Error(err, "Err. Failed to Marshal")
		return err
	}

	n, err := f.Write(data)
	if n != len(data) || err != nil {
		log.Error(err, "Err. Failed to write file")
		return err
	}

	fileName = "./spec/FcstZoneCode.json"
	err = LoadFcstZoneCode(fileName)
	if err != nil {
		log.Error(err)
	}

	return nil
}

// 중기예보
func GetMidtermFcst(mid string) (*ResMidFcst_t, error) {

	req, err := http.NewRequest("GET", MidTermFcstURL, nil)
	if err != nil {
		log.Error(err, "Err, Failed to NewRequest()")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("stnId", MidTermStnIds[mid])

	req.URL.RawQuery = q.Encode()

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

	parse_resp := ResMidFcst_t{}
	err = xml.Unmarshal([]byte(body), &parse_resp)
	if err != nil {
		log.Error("error decoding response: %v", err)
		if e, ok := err.(*xml.SyntaxError); ok {
			log.Error("syntax error.", e.Error())
		}
		log.Error("response: %q", body)
		return nil, err
	}

	return &parse_resp, nil
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
