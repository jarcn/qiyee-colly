package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	max := 240000
	step := 12
	name := fmt.Sprintf("record_%d.csv", max)

	// service
	service := newCsvRecordService()
	go service.run(name)

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second, // keepAlive 超时时间
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 空闲连接超时
		TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 10 * time.Second,
	})
	c.SetRequestTimeout(time.Second * 30)

	// Find and visit all links
	c.OnHTML(".job-card .footer center a.lightblue", func(e *colly.HTMLElement) {
		u := e.Attr("data-url")
		if u != "" {
			// fmt.Printf("dizhi:%+v\n", e.Attr("data-url"))
			e.Request.Visit(u)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting:%s\n", r.URL)
	})
	//detail_detail_job
	c.OnHTML("#detail_job", func(e *colly.HTMLElement) {
		res := make([]string, 0, 9)
		_url := strings.Trim(e.Request.URL.String(), " ")
		name := strings.Trim(e.DOM.Find("#title-comprof").Text(), " ")
		// fmt.Printf("title:%s\n", title)
		website, _ := e.DOM.Find("#detail-comprof a").Attr("href")
		// address:2 lable:4
		var address, lable string
		e.DOM.Find("#detail-comprof").Contents().Each(func(i int, s *goquery.Selection) {
			if i == 2 {
				address = strings.Trim(s.Text(), " ")
			}
			if i == 4 {
				lable = strings.Trim(s.Text(), " ")
				return
			}
		})
		desc := strings.Trim(e.DOM.Find("#comp-detail .jobdesc .desc").Text(), " ")
		detail_address1 := strings.Trim(e.DOM.Find(".detail div:nth-of-type(2)").Text(), " ")
		detail_address2 := strings.Trim(e.DOM.Find(".detail div:nth-of-type(3)").Text(), " ")
		detail_address3 := strings.Trim(e.DOM.Find(".detail div:nth-of-type(4)").Text(), " ")
		res = append(res, _url, name, address, lable, website, desc, detail_address1, detail_address2, detail_address3)
		service.ch <- res
		// fmt.Printf("url:%s,name:%s,address:%s,lable:%s,website:%s,desc:%s,detail_address1:%s,detail_address2:%s,detail_address3:%s\n", _url, name, address, lable, website, desc, detail_address1, detail_address2, detail_address3)
	})
	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL:%s,err:%+v\n", r.Request.URL, err)
	})

	for k := step; k <= max; k += step {
		c.Visit(fmt.Sprintf("https://www.topkarir.com/lowongan/%d/%d?group=0", step, k))
	}

	// c.Visit("https://www.topkarir.com/lowongan/12/132?group=0")
	time.Sleep(time.Second * 3)

	// c.Visit("https://www.topkarir.com/lowongan/detil/jmc-it-consultant-sales-animasi-studio")
}

// 处理用户上报日志服务
type CsvRecordService struct {
	// 通道大小
	len int
	// 缓存消息的通道
	ch chan []string
}

func newCsvRecordService() *CsvRecordService {
	return &CsvRecordService{
		len: 128,
		ch:  make(chan []string, 128),
	}
}

func (c *CsvRecordService) run(name string) {
	f, err := os.Create(name)
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	header := []string{"url", "name", "address", "lable", "website", "desc", "detail_address1", "detail_address2", "detail_address3"}
	w.Write(header)
	for {
		select {
		case ch := <-c.ch:
			w.Write(ch)
		case <-time.After(time.Second * 2):
			w.Flush()
		}
	}
}
