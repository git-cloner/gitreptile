package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Range struct {
	RangeStart string `json:"rangeStart" binding:required`
	RangeEnd   string `json:"rangeEnd" binding:required`
	PageStart  string `json:"pageStart" binding:required`
	PageEnd    string `json:"pageEnd" binding:required`
}

func main() {
	router := gin.Default()
	router.Use(Cors())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	router.GET("/list", func(c *gin.Context) {
		c.JSON(200, c.QueryArray("media"))
	})
	router.POST("/githubStar", Reptile)
	router.Run(":4000")
}

func Reptile(c *gin.Context) {
	var json Range
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startCount := json.RangeStart
	endCount := json.RangeEnd
	pageStart := json.PageStart
	pageEnd := json.PageEnd
	start := time.Now()
	var result []string

	startPage, _ := strconv.Atoi(pageStart)
	endPage, _ := strconv.Atoi(pageEnd)
	for i := startPage; i <= endPage; i++ {
		timer := time.NewTimer(5 * time.Second)
		url := "https://github.com/search?p=" + strconv.Itoa(i) + "&q=stars%3A" + startCount + ".." + endCount + "&type=Repositories"
		tempSlice := parseUrls(url)
		var isNull = tempSlice != nil
		if isNull {
			result = append(result, tempSlice...)
		}
		<-timer.C
		fmt.Println("时间到")
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)

	// 写入文件
	f, err := os.OpenFile("file.txt", os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println("打开文件失败")
	} else {
		for _, value := range result {
			_, err = f.Write([]byte(value))
			if err != nil {
				fmt.Println("写入失败")
			}
		}
	}
	defer f.Close()

	// 返回响应
	c.JSON(200, gin.H{
		"start":  startCount,
		"end":    endCount,
		"result": result,
	})
}

func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return string(body)
}

func parseUrls(url string) []string {
	var result []string
	body := fetch(url)
	body = strings.Replace(body, "\n", "", -1)
	rp := regexp.MustCompile(`<div class="f4 text-normal">(.*?)</div>`)
	idRe := regexp.MustCompile(`<a class="v-align-middle" data-hydro-click=".+" data-hydro-click-hmac=".+" href=".+">(.*?)</a>`)
	items := rp.FindAllStringSubmatch(body, -1)
	for _, item := range items {
		result = append(result, "https://github.com/"+idRe.FindStringSubmatch(item[1])[1]+".git\n")
	}
	return result
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
