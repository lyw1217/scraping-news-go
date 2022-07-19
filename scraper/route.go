package scraper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func articleQuery(c *gin.Context) {
	d_month := int(time.Now().Month())
	d_day := int(time.Now().Day())

	paper := c.Query("paper") // shortcut for c.Request.URL.Query().get("paper")

	if len(paper) > 0 {
		switch paper {
		case "hankyung":
			StatusCode, contents := GetHankyungIssueToday(d_month, d_day)
			if StatusCode != 200 {
				log.Error("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
				break
			}

			c.JSON(200, gin.H{
				"status":   200,
				"contents": contents,
			})

		case "maekyung":
			StatusCode, contents := GetMaekyungMSG(d_month, d_day)
			if StatusCode != 200 {
				log.Error("Err. news.GetMaekyungMSG, StatusCode :", StatusCode)
				break
			}

			c.JSON(200, gin.H{
				"status":   200,
				"contents": contents,
			})

		case "quicknews":
			StatusCode, contents := GetQuickNews(d_month, d_day)
			if StatusCode != 200 {
				log.Error("Err. news.GetQuickNews, StatusCode :", StatusCode)
				break
			}

			c.JSON(200, gin.H{
				"status":   200,
				"contents": contents,
			})
		default:
			c.JSON(400, gin.H{
				"status": 400,
				"reason": "Bad Request",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"status": 400,
			"reason": "Bad Request",
		})
	}
}

func initRoutes() *gin.Engine {
	r := gin.Default()

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

	return r
}

func InitHandler() {

	// Initialize the routes
	routeHttp := initRoutes()

	// HTTP
	err := routeHttp.Run(":9090")
	if err != nil {
		log.Error("routeHttp.Run: ", err)
	}
}
