package scraper

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func weatherKeyword(c *gin.Context, keywords []string) {

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

	// Query : Period, 0: 오늘, 1: 내일
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

		if period > 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"reason": "Query p is invalid.",
			})
			return
		}

		type Fcst_t struct {
			FcstDate string `json:"fcstDate"`
			FcstTime string `json:"fcstTime"`
			Category string `json:"category"`
			Value    string `json:"fcstValue"`
		}

		map_resp := []Fcst_t{}
		item_cnt := 0

		now := time.Now()
		then := now.AddDate(0, 0, period)

		then_date := fmt.Sprintf("%04d%02d%02d", then.Year(), then.Month(), then.Day())
		then_hour := fmt.Sprintf("%02d00", then.Hour())

		for i, v := range resp.Response.Body.Items.Item {
			if v.FcstDate == then_date {
				if v.FcstTime == then_hour {
					item_cnt = i
					break
				}
			}
		}

		for _, v := range resp.Response.Body.Items.Item[item_cnt : item_cnt+12] {
			cat, val := ParseCode(v.Category, v.FcstValue)
			map_resp = append(map_resp, Fcst_t{v.FcstDate, v.FcstTime, cat, val})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"contents": map_resp,
			"name":     resp.Response.Body.Name,
		})
	} else {
		// 조회된 전체 기간 (default : 24h)
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"contents": resp.Response.Body.Items.Item,
			"name":     resp.Response.Body.Name,
		})
	}
}

func weatherMidterm(c *gin.Context, mid string) {

	_, exists := MidTermStnIds[mid]
	if !exists {
		// KEY ERROR!
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Query mid is invalid.",
		})
		return
	}

	resp, err := GetMidtermFcst(mid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"contents": fmt.Sprintf("%s%s", resp.Channel.Item.Title, resp.Channel.Item.Description.Header.Wf),
	})
}

func weatherQuery(c *gin.Context) {

	// Query : [mid], midterm forecast, 중기예보 RSS(지역명)
	mid := c.Query("mid")

	// Query : [k1, k2, k3], Keyword, 검색 키워드(지역명)
	k1 := c.Query("k1")
	k2 := c.Query("k2")
	k3 := c.Query("k3")

	if len(k1) > 0 || len(k2) > 0 || len(k3) > 0 {
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
		weatherKeyword(c, keywords)

	} else if len(mid) > 0 {
		weatherMidterm(c, mid)

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
		return
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
