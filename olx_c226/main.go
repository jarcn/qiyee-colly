package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
)

var mycate = `2000001:89 2000002:657 2000008:18 2000010:577 2000012:16 2000013:156 2000015:51 2000016:183 2000019:11 2000024:49 2000025:52 2000029:30 2000030:149 2000031:678 2000032:637 2000034:29 4000001:212 4000002:39 4000003:292 4000004:369 4000006:50 4000007:11 4000008:12 4000009:9 4000013:9 4000014:11 4000015:9 4000016:14 4000018:875 4000020:787 4000021:158 4000022:74 4000023:14 4000024:479 4000025:8 4000027:6 4000029:776 4000031:721 4000073:8 4000075:11 4000076:593 4000077:9 4000078:36 4000079:782 4000080:831 4000180:7 4000181:7 4000182:6 4000184:47 4000185:14 4000186:9 4000187:7 4000188:8 4000192:29 4000193:9 4000194:7 4000195:6 4000198:13 4000200:5 4000202:407 4000208:6 4000210:8 4000212:72 4000213:16 4000498:7 4000499:48 5000466:587 5000467:336 5000468:254 5000469:235 5000470:203 5000471:91 5000472:139 5000473:157 5000482:100 5000483:117 5000484:151 5000485:131 5000486:85 5000487:79 5000488:101 5000489:71 5000490:215 5000491:113 5000502:29 5000503:255 5000504:67 5000505:178 5000506:781 5000507:169 5002557:48 5002559:60 5002561:66 5002562:136 5002563:34 5002566:35 5002568:26 5002569:99 5002572:106 5002573:45 5002574:49 5002577:70 5002578:42 5002579:56 5002580:24 5002581:37 5002582:48 5002583:45 5002584:35 5002585:69`

var over []string
var mycateMap map[string]int
var cate map[string]int64
var _categoryUrl = "https://www.olx.co.id/api/relevance/v2/search?category=226&facet_limit=100&location=%s&location_facet_limit=20&page=0&platform=web-desktop"
var getDetail bool
var isGetPhone bool

func main() {
	name := fmt.Sprintf("olx_inc_%s.csv", time.Now().Format("20060102"))
	// service
	service := newCsvRecordService()
	go service.run(name)
	// 获取手机号
	// isGetPhone = true
	// getPhone(service.ch)
	// return
	reg1 := regexp.MustCompile(`window.__APP = (.*)*;`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return
	}
	mycateMap = make(map[string]int)
	var tt int
	_c := strings.Split(mycate, " ")
	for _, v := range _c {
		a := strings.Split(v, ":")
		_i, _ := strconv.Atoi(a[1])
		mycateMap[a[0]] = _i
		tt += _i
	}
	fmt.Println("total:", tt)
	getDetail = true
	cate = make(map[string]int64)

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second, // 超时时间
			KeepAlive: 60 * time.Second, // keepAlive 超时时间
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 空闲连接超时
		TLSHandshakeTimeout:   30 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 30 * time.Second,
	})
	c.SetRequestTimeout(time.Second * 60)

	//配置反爬策略(设置ua和refer扩展)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		res := make([]string, 0, 20)
		// 获取id
		_url := e.Request.URL.String()
		isDetail := strings.Contains(_url, "iid")
		if isDetail {
			e.ForEach("script", func(i int, h *colly.HTMLElement) {
				id := e.Request.URL.String()[strings.LastIndex(_url, "-")+1:]
				if i == 3 {
					data := h.Text
					title := gjson.Get(data, "title").String()
					name := gjson.Get(data, "name").String()
					employmentType := gjson.Get(data, "employmentType").String()
					salaryCurrency := gjson.Get(data, "salaryCurrency").String()
					salaryMinValue := gjson.Get(data, "baseSalary.minValue").String()
					salaryMaxValue := gjson.Get(data, "baseSalary.maxValue").String()
					image := gjson.Get(data, "image").String()
					description := gjson.Get(data, "description").String()
					addressRegion := gjson.Get(data, "jobLocation.address.addressRegion").String()
					addressLocality := gjson.Get(data, "jobLocation.address.addressLocality").String()
					address := gjson.Get(data, "jobLocation.address.name").String()
					datePosted := gjson.Get(data, "datePosted").String()
					validThrough := gjson.Get(data, "validThrough").String()
					res = append(res, id, _url, title, name, employmentType, salaryCurrency,
						salaryMinValue, salaryMaxValue, image, description, addressRegion, addressLocality, address, datePosted, validThrough)

				}
				if i == 5 {
					data2 := strings.Trim(strings.Trim(strings.ReplaceAll(h.Text, "window.__APP = ", " "), " "), ";")
					data2 = strings.ReplaceAll(data2, "props", `"props"`)
					data2 = strings.ReplaceAll(data2, "states", `"states"`)
					data2 = strings.ReplaceAll(data2, "config", `"config"`)
					data2 = strings.ReplaceAll(data2, "translations", `"translations"`)
					images := gjson.Get(data2, "states.items.elements."+id+".images").String()
					userID := gjson.Get(data2, "states.items.elements."+id+".user_id").String()
					userCreatedAt := gjson.Get(data2, "states.users.elements."+userID+".created_at").String()
					userName := gjson.Get(data2, "states.users.elements."+userID+".name").String()
					res = append(res, images, userID, userName, userCreatedAt)
				}
			})
			res = append(res, e.DOM.Find("label._10wN3 span").Text())
			if len(res) == 20 {
				service.ch <- res
			}
		}
	})

	// Find and visit all links
	// c.OnHTML("script[data-rh]", func(e *colly.HTMLElement) {
	// 	//
	// 	data := e.Text
	// 	res := make([]string, 0, 14)
	// 	_url := e.Request.URL.String()
	// 	title := gjson.Get(data, "title").String()
	// 	// name := gjson.Get(data, "name").String()
	// 	employmentType := gjson.Get(data, "employmentType").String()
	// 	salaryCurrency := gjson.Get(data, "salaryCurrency").String()
	// 	salaryMinValue := gjson.Get(data, "baseSalary.minValue").String()
	// 	salaryMaxValue := gjson.Get(data, "baseSalary.maxValue").String()
	// 	image := gjson.Get(data, "image").String()
	// 	description := gjson.Get(data, "description").String()
	// 	addressRegion := gjson.Get(data, "jobLocation.address.addressRegion").String()
	// 	addressLocality := gjson.Get(data, "jobLocation.address.addressLocality").String()
	// 	address := gjson.Get(data, "jobLocation.address.name").String()
	// 	datePosted := gjson.Get(data, "datePosted").String()
	// 	validThrough := gjson.Get(data, "validThrough").String()
	// 	res = append(res, _url, title, employmentType, salaryCurrency,
	// 		salaryMinValue, salaryMaxValue, image, description, addressRegion, addressLocality, address, datePosted, validThrough)
	// 	service.ch <- res
	// })

	// jobsearch-ViewJobLayout-jobDisplay

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting:%s\n", r.URL)
		// time.Sleep(time.Second)
		// if r.URL.String() == "https://www.jobstreet.co.id/en/job-search/job-vacancy/1/" {
		// 	panic(r.URL)
		// }
	})
	c.OnResponse(func(r *colly.Response) {

		if r.StatusCode == http.StatusOK {
			data := string(r.Body)
			// 判断是否大于0并小于1000，如果是就直接解析，否则循环子级别
			if !getDetail {
				location := gjson.Get(data, "metadata.filters.2.values").Array()
				f(c, location)
				fmt.Printf("category:%+v\n", cate)
			} else {
				// 获取并提取链接
				_data := gjson.Get(data, "data").Array()
				for _, v := range _data {
					_id := gjson.Get(v.Raw, "id").String()
					_title := gjson.Get(v.Raw, "title").String()
					_title = strings.ReplaceAll(strings.ToLower(_title), " ", "-")
					_url := "https://www.olx.co.id/item/" + _title + "-iid-" + _id
					c.Visit(_url)
				}
				_nextUrl := gjson.Get(data, "metadata.next_page_url").String()
				c.Visit(_nextUrl)
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL:%s,err:%+v\n", r.Request.URL, err)
		if r.StatusCode != http.StatusNotFound {
			r.Request.Retry()
		}
	})

	// 爬取全部
	// _detailUrl := "https://www.olx.co.id/api/relevance/v2/search?category=226&facet_limit=100&location=%s&location_facet_limit=20&page=%d&platform=web-desktop"
	// for k, v := range mycateMap {
	// 	tp := v / 20
	// 	for i := 0; i <= tp; i++ {
	// 		c.Visit(fmt.Sprintf(_detailUrl, k, i))
	// 	}
	// 	over = append(over, k)
	// }

	// c.Visit("https://www.olx.co.id/item/gratis-ijazah-sma-gampang-kerja-ikuti-paket-c-maksimal-20-tahun-iid-866209908")
	// https://www.olx.co.id/api/relevance/v2/search?category=226&facet_limit=100&location=1000001&location_facet_limit=20&page=1&platform=web-desktop
	c.Visit("https://www.olx.co.id/api/relevance/v2/search?category=226&facet_limit=100&location=1000001&location_facet_limit=20&page=0&platform=web-desktop")
	// 分类地址
	// c.Visit(fmt.Sprintf("https://id.indeed.com/lowongan-kerja?q=%s&start=%d&filter=0&vjk=bddaa9c3d03e3959", "dibutuhkan%20segera", 60))
	// c.Visit("https://id.indeed.com/lihat-lowongan-kerja?cmp=CV.-Heloklin-Indonesia&t=Customer+Care&jk=7ab3564df29c0cae&vjs=3")
	time.Sleep(time.Second * 3)
}

func f(c *colly.Collector, arr []gjson.Result) {
	for _, v := range arr {
		_id := gjson.Get(v.Raw, "id").String()
		_count := gjson.Get(v.Raw, "count").Int()
		_child := gjson.Get(v.Raw, "children").Array()
		fmt.Printf("id:%s,count:%d\n", _id, _count)
		if _count >= 1000 && len(_child) == 0 && _id[:1] != "5" {
			// 递归
			c.Visit(fmt.Sprintf(_categoryUrl, _id))
			continue
		}
		// 加入
		if len(_child) > 0 {
			f(c, _child)
		} else {
			cate[_id] = _count
		}
	}
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

	header := []string{"id", "_url", "title", "name", "employmentType", "salaryCurrency",
		"salaryMinValue", "salaryMaxValue", "image", "description", "addressRegion", "addressLocality",
		"address", "datePosted", "validThrough",
		"images", "userID", "userName", "userCreatedAt", "highlight"}
	if isGetPhone {
		header = append(header, "phone")
	}
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
