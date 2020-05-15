package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var start string
var end string

func init() {
	flag.StringVar(&start, "start", "1000", "start count")
	flag.StringVar(&end, "end", "2000", "end count")
}

func main() {
	flag.Parse()
	url := "http://localhost:4000/githubStar"
	for i := 1; i <= 100; i++ {
		timer := time.NewTimer(100 * time.Second)
		var post = "{" +
			"\"rangeStart\":\"" + start + "\"," +
			"\"rangeEnd\":\"" + end + "\"," +
			"\"pageStart\":\"" + strconv.Itoa(i) + "\"," +
			"\"pageEnd\":\"" + strconv.Itoa(i) + "\"" +
			"}"
		var jsonStr = []byte(post)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("create request failed")
		}
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("send request failed")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll((resp.Body))
		if resp.StatusCode == 200 {
			fmt.Printf("range %s - %s number %i is successfull", start, end, i)
		} else {
			fmt.Println("response is not 200")
			fmt.Println(body)
		}
		<-timer.C
		fmt.Println("开是请求下一页")
	}
}
