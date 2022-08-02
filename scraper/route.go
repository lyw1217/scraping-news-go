package scraper

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func weatherQuery(c *gin.Context) {

	// Query : Keyword, 검색 키워드(지역명)
	k1 := c.Query("k1")
	k2 := c.Query("k2")
	k3 := c.Query("k3")

	if len(k1) == 0 && len(k2) == 0 && len(k3) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
		return
	}

	keywords := make([]string, 0)

	if len(k1) > 0 {
		keywords = append(keywords, k1)
	}
	if len(k2) > 0 {
		keywords = append(keywords, k2)
	}
	if len(k3) > 0 {
		keywords = append(keywords, k3)
	}

	resp, err := GetVilageFcstInfo(keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Internal Server Error",
		})
		return
	}
	// 키워드 조회 결과 없음
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"reason": "Not Found",
		})
		return
	}

	// Query : Period
	p := c.Query("p")

	if len(p) > 0 {
		period, err := strconv.Atoi(p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"reason": "Query p is not integer.",
			})
			return
		}
		// 최대 24시간 조회 제한
		if period > 24 {
			period = 24
		}

		type Fcst_t struct {
			FcstDate string `json:"fcstDate"`
			FcstTime string `json:"fcstTime"`
			Category string `json:"category"`
			Value    string `json:"fcstValue"`
		}

		map_resp := map[int][]Fcst_t{}
		contents := 0 // 1시간당 12개 item

		for i := 0; i < period; i++ {
			for _, v := range resp.Response.Body.Items.Item[contents : contents+12] {
				cat, val := ParseCode(v.Category, v.FcstValue)
				map_resp[i] = append(map_resp[i], Fcst_t{v.FcstDate, v.FcstTime, cat, val})
			}
			contents += 12
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"contents": map_resp,
		})
	} else {
		// 조회된 전체 기간 (default : 24h)
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"contents": resp.Response.Body.Items.Item,
		})
	}
}

func articleQuery(c *gin.Context) {
	d_month := int(time.Now().Month())
	d_day := int(time.Now().Day())

	paper := c.Query("paper") // shortcut for c.Request.URL.Query().get("paper")

	if len(paper) > 0 {
		switch paper {
		case "hankyung":
			StatusCode, contents := GetHankyungIssueToday(d_month, d_day)
			if StatusCode != http.StatusOK {
				log.Error("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
				break
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": contents,
			})

		case "maekyung":
			StatusCode, contents := GetMaekyungMSG(d_month, d_day)
			if StatusCode != http.StatusOK {
				log.Error("Err. news.GetMaekyungMSG, StatusCode :", StatusCode)
				break
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": contents,
			})

		case "quicknews":
			StatusCode, contents := GetQuickNews(d_month, d_day)
			if StatusCode != http.StatusOK {
				log.Error("Err. news.GetQuickNews, StatusCode :", StatusCode)
				break
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": contents,
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"reason": "Bad Request",
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
	}
}

func initRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(
			http.StatusOK,
			"Hello World. GOSCRAPER!",
		)
	})

	r.GET("/maekyung", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, MkMSGUrl)
	})

	r.GET("/hankyung", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, HkIssueTodayUrl)
	})

	r.GET("/quicknews", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, QuicknewsUrl)
	})

	r.GET("/article", articleQuery)

	r.GET("/weather", weatherQuery)

	return r
}

func InitHandler() {

	// Initialize the routes
	routeHttp := initRoutes()

	// HTTP
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Wrong Value of environment : $PORT = '", port, "'")
		os.Exit(1)
	}
	err := routeHttp.Run(":" + port)
	if err != nil {
		log.Error("routeHttp.Run: ", err)
	}
}
