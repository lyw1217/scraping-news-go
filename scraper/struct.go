package scraper

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
