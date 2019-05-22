package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

//https://gophp.io/parsing-websites-with-golang-and-colly/
//https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/

func main() {
	result, err := os.Create("result.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()

	words, err := os.Open("words.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer words.Close()

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		word, stars := getStars(scanner.Text())
		fmt.Println(word, stars)
		result.WriteString(word + " " + strconv.Itoa(stars) + "\n")
	}
}

var host = "https://www.macmillandictionary.com/"
var agent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"

func getStars(w string) (word string, stars int) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	c.UserAgent = agent
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", host)
		r.Headers.Set("Referer", host)
		r.Headers.Set("Accept-Encoding", "gzip,deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	})

	// c.OnResponse(func(resp *colly.Response) {
	// 	//fmt.Println(string(resp.Body))
	// 	//print("this is response.")
	// })

	c.OnHTML("div#headwordleft", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildText(".BASE-FORM"))
		w := e.ChildText(".BASE")
		s := e.DOM.Find(".stars_grp .icon_star").Length()
		word = w
		stars = s
	})
	c.OnError(func(resp *colly.Response, err error) {
		fmt.Println(err)
	})

	url := "https://www.macmillandictionary.com/dictionary/american/" + w
	c.Visit(url)
	return
}
