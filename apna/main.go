package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
)

var totalPage = 10000

func main() {
	insert()
	return
	// max := 10000
	name := fmt.Sprintf("apna_%s.csv", time.Now().Format("20060102"))
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

	//配置反爬策略(设置ua和refer扩展)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Find and visit all links
	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		data := e.DOM.Text()
		// fmt.Println(data)
		// return
		res := make([]string, 0, 24)
		_url := e.Request.URL.String()
		job := gjson.Get(data, "props.pageProps.job").String()
		id := gjson.Get(job, "id").String()
		title := gjson.Get(job, "title").String()
		num_openings := gjson.Get(data, "props.pageProps.jobApplicationData.num_openings").String()
		interested_count := gjson.Get(data, "props.pageProps.jobApplicationData.count").String()
		category := gjson.Get(job, "category").String()
		company_id := gjson.Get(job, "organization.id").String()
		company_name := gjson.Get(job, "organization.name").String()
		address_area := gjson.Get(job, "location_name").String()
		min_salary := gjson.Get(job, "min_salary").String()
		max_salary := gjson.Get(job, "max_salary").String()
		shift := gjson.Get(job, "shift").String()
		is_part_time := gjson.Get(job, "is_part_time").String()
		is_wfh := gjson.Get(job, "is_wfh").String()
		created_on := gjson.Get(job, "created_on").String()
		expiry := gjson.Get(job, "expiry").String()
		education := gjson.Get(job, "education").String()
		english := gjson.Get(job, "english").String()
		min_experience := gjson.Get(job, "min_experience").String()
		max_experience := gjson.Get(job, "max_experience").String()
		experience_level := gjson.Get(job, "experience_level").String()
		gender := gjson.Get(job, "gender").String()
		address := gjson.Get(job, "company_address.line_1").String()
		job_description := trimHtml(gjson.Get(data, "props.pageProps.jobDescription").String())
		last_updated := gjson.Get(job, "last_updated").String()

		res = append(res, _url, id, title, num_openings, interested_count, category, company_id, company_name, address_area, min_salary,
			max_salary, shift, is_part_time, is_wfh, created_on, expiry, education, english, min_experience,
			max_experience, experience_level, gender, address, job_description, last_updated)
		service.ch <- res

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting:%s\n", r.URL)
		// if r.URL.String() == "https://www.jobstreet.co.id/en/job-search/job-vacancy/1/" {
		// 	panic(r.URL)
		// }
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL:%s,err:%+v\n", r.Request.URL, err)
	})

	for k := 1; k <= 500; k++ {
		req(c, k)
	}
	// req(c, 3)
	// c.Visit(fmt.Sprintf("https://www.jobstreet.co.id/en/job-search/job-vacancy/%d/", 2))
	time.Sleep(time.Second * 3)
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
	header := []string{"_url", "id", "title", "num_openings", "interested_count", "category", "company_id", "company_name", "address_area", "min_salary",
		"max_salary", "shift", "is_part_time", "is_wfh", "created_on", "expiry", "education", "english", "min_experience",
		"max_experience", "experience_level", "gender", "address", "job_description", "last_updated"}
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

const _url = `https://apna.co/_next/data/daFCSXO_8Y8piRQdWC59U/jobs.json?page=%d`

func req(c *colly.Collector, page int) {
	reqUrl := fmt.Sprintf(_url, page)
	fmt.Printf("req page:%d\n", page)
	resp, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := string(body)
	_totalPage := gjson.Get(data, "pageProps.totalJobCount").Int()
	if _totalPage > 0 {
		totalPage = int(_totalPage)
	}
	pageUrl := gjson.Get(data, "pageProps.jobs").Array()
	for _, v := range pageUrl {
		detailUrl := gjson.Get(v.Raw, "public_url").String()
		c.Visit(detailUrl)
		// fmt.Println(detailUrl)
	}

}

func trimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
