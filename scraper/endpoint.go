package scraper

// URLs
const (
	MkMSGUrl        string = "https://www.mk.co.kr/premium/series/20007/"
	HkIssueTodayUrl string = "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="
	QuicknewsUrl    string = "https://quicknews.co.kr/"
	HostName        string = "http://mumeog.site/?redirect="
	VilageFcstUrl   string = "http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getVilageFcst" // 기상청_단기예보 ((구)_동네예보) 조회서비스

	// 기상청 동네 zone code
	TopURL  string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/top.json.txt"
	MdlURL  string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/mdl.%s.json.txt"
	LeafURL string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/leaf.%s.json.txt"

	// 기상청 중기예보 RSS
	MidTermFcstURL string = "https://www.kma.go.kr/weather/forecast/mid-term-rss3.jsp"

	// 네이버 로마자 한글 변환
	RomanizationURL string = "https://openapi.naver.com/v1/krdict/romanization"
	// 네이버 파파고 번역
	PapagoURL string = "https://openapi.naver.com/v1/papago/n2mt"
	// 네이버 언어 감지
	DetectLangURL string = "https://openapi.naver.com/v1/papago/detectLangs"
)

var MidTermStnIds = map[string]string{
	"전국": "108",

	"서울": "109",
	"경기": "109",

	"강원": "105",

	"충북": "131",

	"대전": "133",
	"세종": "133",
	"충남": "133",

	"전북": "146",

	"광주": "156",
	"전남": "156",

	"대구": "143",
	"경북": "143",

	"부산": "159",
	"울산": "159",
	"경남": "159",

	"제주": "184",
}
