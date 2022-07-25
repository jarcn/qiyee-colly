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
	max := 10000
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

	//配置反爬策略(设置ua和refer扩展)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// Find and visit all links
	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		data := e.DOM.Text()
		// fmt.Println(a)
		res := make([]string, 0, 24)
		_url := e.Request.URL.String()
		vacancyBySlug := gjson.Get(data, "props.pageProps.initialState.detail.vacancyBySlug").String()
		id := gjson.Get(vacancyBySlug, "id").String()
		positionName := gjson.Get(vacancyBySlug, "positionName").String()
		companyID := gjson.Get(vacancyBySlug, "company.id").String()
		companyName := gjson.Get(vacancyBySlug, "company.name").String()
		postedDate := gjson.Get(vacancyBySlug, "postedDate").String()
		postedDateStr := gjson.Get(vacancyBySlug, "postedDateStr").String()
		requirementStr := gjson.Get(vacancyBySlug, "requirementStr").String()
		educationLevel := gjson.Get(vacancyBySlug, "educationLevel").String()
		gender := gjson.Get(vacancyBySlug, "gender").String()
		maxAge := gjson.Get(vacancyBySlug, "maxAge").String()
		minExperience := gjson.Get(vacancyBySlug, "minExperience").String()
		province := gjson.Get(vacancyBySlug, "province.name").String()
		city := gjson.Get(vacancyBySlug, "city.name").String()
		typeStr := gjson.Get(vacancyBySlug, "typeStr").String()
		locationSite := gjson.Get(vacancyBySlug, "LocationSite").String()
		vacancyCount := gjson.Get(vacancyBySlug, "vacancyCount").String()
		salaryLowerBoundStr := gjson.Get(vacancyBySlug, "salaryLowerBoundStr").String()
		salaryUpperBoundStr := gjson.Get(vacancyBySlug, "salaryUpperBoundStr").String()
		description := gjson.Get(vacancyBySlug, "description").String()
		workingDayStr := gjson.Get(vacancyBySlug, "workingDayStr").String()
		workingHourStr := gjson.Get(vacancyBySlug, "workingHourStr").String()
		companyDescription := gjson.Get(vacancyBySlug, "company.description").String()
		contactWeblink := gjson.Get(vacancyBySlug, "company.contactWeblink").String()
		res = append(res, _url, id, positionName, companyID, companyName, postedDate, postedDateStr, requirementStr,
			educationLevel, gender, maxAge, minExperience, province, city, typeStr, locationSite, vacancyCount,
			salaryLowerBoundStr, salaryUpperBoundStr, description, workingDayStr, workingHourStr, companyDescription, contactWeblink)
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

	for k := 1; k <= max; k++ {
		req(c, k)
	}
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
	header := []string{"_url", "id", "positionName", "companyID", "companyName", "postedDate", "postedDateStr", "requirementStr",
		"educationLevel", "gender", "maxAge", "minExperience", "province", "city", "typeStr", "locationSite", "vacancyCount",
		"salaryLowerBoundStr", "salaryUpperBoundStr", "description", "workingDayStr", "workingHourStr", "companyDescription", "contactWeblink"}
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

const payload = `{"operationName":"getSearchJob","variables":{"filter":{"page":%d,"limit":5},"provinceName":"","cityName":"","districtName":"","keyword":"","categories":[],"filters":[]},"query":"fragment JobHome on Vacancy {\n  city {\n    name\n    __typename\n  }\n  company {\n    name\n    __typename\n  }\n  deferredLinkWeb\n  educationLevelStr\n  genderStr\n  id\n  isBookmarked\n  isClosed\n  isPublished\n  LocationSite: locationSiteStr\n  maxAge\n  minExperience\n  positionName\n  postedDate: createdAt\n  postedDateStr\n  province {\n    name\n    __typename\n  }\n  requirementStr\n  salaryLowerBound\n  salaryUpperBound\n  salaryLowerBoundStr\n  salaryUpperBoundStr\n  skillStr\n  slug\n  typeStr\n  vacancyCount: vacancyCountStr\n  __typename\n}\n\nquery getSearchJob($filter: CommonFilter, $provinceName: String, $cityName: String, $districtName: String, $keyword: String, $categories: [SearchCategoryFilter], $filters: [SearchCategoryFilter]) {\n  vacanciesV2(\n    filter: $filter\n    provinceName: $provinceName\n    cityName: $cityName\n    districtName: $districtName\n    keyword: $keyword\n    categories: $categories\n    filters: $filters\n  ) {\n    page\n    elements\n    provinceName\n    provinceId\n    cityName\n    cityId\n    list {\n      ...JobHome\n      __typename\n    }\n    __typename\n  }\n}"}`
const _url = `https://cons-gql.kitalulus.com/graphql`

func req(c *colly.Collector, page int) {
	p := fmt.Sprintf(payload, page)
	b := strings.NewReader(p)
	fmt.Printf("req page:%d\n", page)
	resp, err := http.Post(_url, "application/json", b)
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
	_totalPage := gjson.Get(data, "data.vacanciesV2.page").Int()
	if _totalPage > 0 {
		totalPage = int(_totalPage)
	}
	pageUrl := gjson.Get(data, "data.vacanciesV2.list").Array()
	for _, v := range pageUrl {
		detailUrl := gjson.Get(v.Raw, "deferredLinkWeb").String()
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
