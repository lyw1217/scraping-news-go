package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scraping-news/config"
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

func GetPapagoRomanization(c *gin.Context, query string) (*ResRoman_t, error) {
	req, err := http.NewRequest("GET", RomanizationURL, nil)
	if err != nil {
		log.Error(err, "Err, Failed to NewRequest()")
		return nil, err
	}

	q := req.URL.Query()
	q.Add("query", query)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Naver-Client-Id", config.Keys.Naver.ClientId)
	req.Header.Set("X-Naver-Client-Secret", config.Keys.Naver.ClientSecret)

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

	parse_resp := ResRoman_t{}
	err = json.Unmarshal([]byte(body), &parse_resp)
	if err != nil {
		log.Error("error decoding response: ", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Error("syntax error.", e.Error())
		}
		log.Error("response: ", string(body))
		return nil, err
	}

	return &parse_resp, nil
}