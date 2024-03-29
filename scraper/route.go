package scraper

import (
	b64 "encoding/base64"
	"net/http"
	"os"
	"scraping-news/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func papagoQuery(c *gin.Context) {
	auth := c.Query("auth")
	if !checkAuth(auth) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"reason": "Unauthorized API Key",
		})
		return
	}

	// 언어 감지
	// Query : text, 어떤 언어인지 확인할 텍스트
	text := c.Query("text")

	if len(text) > 0 || len(text) <= 5000 {
		langCode, err := PostDetectLangs(text)
		if err != nil {
			log.Error(err, " Err. failed to PostDetectLangs")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}

		// 파파고 번역
		if len(langCode) > 0 {
			resp, err := PostPapagoTrans(langCode, text)
			if err != nil {
				log.Error(err, " Err. failed to PostPapagoTrans")
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": http.StatusInternalServerError,
					"reason": "Internal Server Error",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": resp.Message.Result.TranslatedText,
			})
			return
		} else {
			log.Error("Err. failed to GetDetectLangs")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
		return
	}
}

func romanizationQuery(c *gin.Context) {
	auth := c.Query("auth")
	if !checkAuth(auth) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"reason": "Unauthorized API Key",
		})
		return
	}

	// Query : query, 로마자로 변환할 한글 이름
	query := c.Query("query")
	if len(query) > 0 {
		resp, err := GetPapagoRomanization(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"reason": "Internal Server Error",
			})
			return
		}

		if len(resp.AResult) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": resp.AResult[0].AItems,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"reason": "Bad Request",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"reason": "Bad Request",
		})
		return
	}
}

func weatherQuery(c *gin.Context) {
	auth := c.Query("auth")
	if !checkAuth(auth) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"reason": "Unauthorized API Key",
		})
		return
	}

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
	auth := c.Query("auth")
	if !checkAuth(auth) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"reason": "Unauthorized API Key",
		})
		return
	}

	d_month := int(time.Now().Month())
	d_day := int(time.Now().Day())

	paper := c.Query("paper") // shortcut for c.Request.URL.Query().get("paper")

	if len(paper) > 0 {
		switch paper {
		case "hankyung":
			StatusCode, contents := GetHankyungIssueToday(d_month, d_day)
			if StatusCode != http.StatusOK {
				log.Error("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
				c.JSON(StatusCode, gin.H{
					"status": StatusCode,
					"reason": contents,
				})
				break
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   http.StatusOK,
				"contents": contents,
			})

		case "maekyung":
			StatusCode, contents, _ := GetMaekyungMSG(d_month, d_day)
			if StatusCode != http.StatusOK {
				log.Error("Err. news.GetMaekyungMSG, StatusCode :", StatusCode)
				c.JSON(StatusCode, gin.H{
					"status": StatusCode,
					"reason": contents,
				})
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
				c.JSON(StatusCode, gin.H{
					"status": StatusCode,
					"reason": contents,
				})
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

func checkAuth(key string) bool {
	sDec, _ := b64.StdEncoding.DecodeString(key)
	if strings.Compare(strings.Trim(string(sDec), " "), config.Keys.Newyo.Apikey) == 0 {
		return true
	} else {
		return false
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

	r.GET("/romanization", romanizationQuery)

	r.POST("/papago", papagoQuery)

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
