package scraper

import "encoding/xml"

type LinkInfo_t struct {
	Title string
	Url   string
}

type PageInfo_t struct {
	StatusCode int
	Links      []LinkInfo_t
	Contents   []string
}

type VilageInfo_t struct {
	Step1 string `json:"step1"`
	Step2 string `json:"step2"`
	Step3 string `json:"step3"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

type ReqVilageFcst_t struct {
	ServiceKey string // 인증키
	NumOfRows  string // 한 페이지 결과 수
	PageNo     string // 페이지 번호
	DataType   string // 응답자료형식
	Base_date  string // 발표일자
	Base_time  string // 발표시각
	Nx         string // 예보지점 X 좌표
	Ny         string // 예보지점 Y 좌표
	Name       string // 예보지점 이름
}

type FcstHeader_t struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

type FcstItem_t struct {
	BaseDate  string `json:"baseDate"`  // 발표일자
	BaseTime  string `json:"baseTime"`  // 발표시각
	FcstDate  string `json:"fcstDate"`  // 예보일자
	FcstTime  string `json:"fcstTime"`  // 예보시각
	Category  string `json:"category"`  // 자료구분문자
	FcstValue string `json:"fcstValue"` // 예보 값
	Nx        int    `json:"nx"`        // 예보지점 X 좌표
	Ny        int    `json:"ny"`        // 예보지점 Y 좌표
}

type FcstItems_t struct {
	Item []FcstItem_t `json:"item"`
}

type FcstBody_t struct {
	DataType   string      `json:"dataType"`
	Items      FcstItems_t `json:"items"`
	PageNo     int         `json:"pageNo"`
	NumOfRows  int         `json:"numOfRows"`
	TotalCount int         `json:"totalCount"`
	Name       string      `json:"name"`
}

type FcstResp_t struct {
	Header FcstHeader_t `json:"header"`
	Body   FcstBody_t   `json:"body"`
}

type ResVilageFcst_t struct {
	Response FcstResp_t `json:"response"`
}

type FcstZone_t struct {
	Step1 string `json:"step1"`
	Step2 string `json:"step2"`
	Step3 string `json:"step3"`
	Code  string `json:"code"`
	Name  string `json:"value"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

// made by 'https://www.onlinetool.io/xmltogo/'
// 불필요 파싱요소 주석처리
type ResMidFcst_t struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		//Title       string `xml:"title"`
		//Link        string `xml:"link"`
		//Description string `xml:"description"`
		//Language    string `xml:"language"`
		//Generator   string `xml:"generator"`
		PubDate string `xml:"pubDate"`
		Item    struct {
			Text string `xml:",chardata"`
			//Author      string `xml:"author"`
			//Category    string `xml:"category"`
			Title string `xml:"title"`
			//Link        string `xml:"link"`
			//Guid        string `xml:"guid"`
			Description struct {
				Text   string `xml:",chardata"`
				Header struct {
					Text  string `xml:",chardata"`
					Title string `xml:"title"`
					Tm    string `xml:"tm"`
					Wf    string `xml:"wf"`
				} `xml:"header"`
				/*
					Body struct {
						Text     string `xml:",chardata"`
						Location []struct {
							Text     string `xml:",chardata"`
							WlVer    string `xml:"wl_ver,attr"`
							Province string `xml:"province"`
							City     string `xml:"city"`
							Data     []struct {
								Text        string `xml:",chardata"`
								Mode        string `xml:"mode"`
								TmEf        string `xml:"tmEf"`
								Wf          string `xml:"wf"`
								Tmn         string `xml:"tmn"`
								Tmx         string `xml:"tmx"`
								Reliability string `xml:"reliability"`
								RnSt        string `xml:"rnSt"`
							} `xml:"data"`
						} `xml:"location"`
					} `xml:"body"`
				*/
			} `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type ResRoman_t struct {
	AResult []struct {
		SFirstName string `json:"sFirstName"`
		AItems     []struct {
			Name  string `json:"name"`
			Score string `json:"score"`
		} `json:"aItems"`
	} `json:"aResult"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type ResLangCode_t struct {
	LangCode     string `json:"langCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

type ResPapago_t struct {
	Message struct {
		Type    string `json:"@type"`
		Service string `json:"@service"`
		Version string `json:"@version"`
		Result  struct {
			SrcLangType    string `json:"srcLangType"`
			TarLangType    string `json:"tarLangType"`
			TranslatedText string `json:"translatedText"`
		} `json:"result"`
	} `json:"message"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}