package scraper

const (
	MkMSGUrl        string = "https://www.mk.co.kr/premium/series/20007/"
	HkIssueTodayUrl string = "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="
	QuicknewsUrl    string = "https://quicknews.co.kr/"
	HostName        string = "http://mumeog.site:9090/"
	VilageFcstUrl   string = "http://apis.data.go.kr/1360000/VilageFcstInfoService_2.0/getVilageFcst" // 기상청_단기예보 ((구)_동네예보) 조회서비스

	// 기상청 동네 zone code
	TopURL  string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/top.json.txt"
	MdlURL  string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/mdl.%s.json.txt"
	LeafURL string = "http://www.kma.go.kr/DFSROOT/POINT/DATA/leaf.%s.json.txt"
)
