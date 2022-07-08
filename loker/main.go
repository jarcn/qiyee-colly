package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func main() {
	max := 60
	name := fmt.Sprintf("record_%d_%d.csv", max, rand.Int31n(100))
	// service
	service := newCsvRecordService()
	go service.run(name)

	f, err := os.Open("a.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	v0 := ""
	doc.Find("select#category option").Each(func(i int, s *goquery.Selection) {
		res := make([]string, 0, 2)
		// For each item found, get the title
		if level, ok := s.Attr("class"); ok {
			title := s.Text()
			if level == "level-0" {
				v0 = title
			}
			res := append(res, v0, trimHtml(title))
			service.ch <- res
		}
	})

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
	header := []string{"category", "sub_category"}
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

const payload = `{"query":"query getJobDetail($jobId: String, $locale: String, $country: String, $candidateId: ID, $solVisitorId: String, $flight: String) {\n  jobDetail(\n    jobId: $jobId\n    locale: $locale\n    country: $country\n    candidateId: $candidateId\n    solVisitorId: $solVisitorId\n    flight: $flight\n  ) {\n    id\n    pageUrl\n    jobTitleSlug\n    applyUrl {\n      url\n      isExternal\n    }\n    isExpired\n    isConfidential\n    isClassified\n    accountNum\n    advertisementId\n    subAccount\n    showMoreJobs\n    adType\n    header {\n      banner {\n        bannerUrls {\n          large\n        }\n      }\n      salary {\n        max\n        min\n        type\n        extraInfo\n        currency\n        isVisible\n      }\n      logoUrls {\n        small\n        medium\n        large\n        normal\n      }\n      jobTitle\n      company {\n        name\n        url\n        slug\n        advertiserId\n      }\n      review {\n        rating\n        numberOfReviewer\n      }\n      expiration\n      postedDate\n      postedAt\n      isInternship\n    }\n    companyDetail {\n      companyWebsite\n      companySnapshot {\n        avgProcessTime\n        registrationNo\n        employmentAgencyPersonnelNumber\n        employmentAgencyNumber\n        telephoneNumber\n        workingHours\n        website\n        facebook\n        size\n        dressCode\n        nearbyLocations\n      }\n      companyOverview {\n        html\n      }\n      videoUrl\n      companyPhotos {\n        caption\n        url\n      }\n    }\n    jobDetail {\n      summary\n      jobDescription {\n        html\n      }\n      jobRequirement {\n        careerLevel\n        yearsOfExperience\n        qualification\n        fieldOfStudy\n        industryValue {\n          value\n          label\n        }\n        skills\n        employmentType\n        languages\n        postedDate\n        closingDate\n        jobFunctionValue {\n          code\n          name\n          children {\n            code\n            name\n          }\n        }\n        benefits\n      }\n      whyJoinUs\n    }\n    location {\n      location\n      locationId\n      omnitureLocationId\n    }\n    sourceCountry\n  }\n}\n","variables":{"jobId":"%s","country":"id","locale":"en","candidateId":"","solVisitorId":""}}`
const _url = `https://xapi.supercharge-srp.co/job-search/graphql?country=id&isSmartSearch=true`

func req(s *CsvRecordService, id string) {
	p := fmt.Sprintf(payload, id)
	b := strings.NewReader(p)
	fmt.Printf("req id:%s\n", id)
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
	res := make([]string, 0, 23)
	_id := gjson.Get(data, "data.jobDetail.id").String()
	pageUrl := gjson.Get(data, "data.jobDetail.pageUrl").String()
	logoUrl := gjson.Get(data, "data.jobDetail.header.logoUrls.normal").String()
	salary := gjson.Get(data, "data.jobDetail.header.salary").String()
	jobTitle := gjson.Get(data, "data.jobDetail.header.jobTitle").String()
	companyName := gjson.Get(data, "data.jobDetail.header.company.name").String()
	postedAt := gjson.Get(data, "data.jobDetail.header.postedAt").String()
	companyWebsite := gjson.Get(data, "data.jobDetail.companyDetail.companyWebsite").String()
	companySize := gjson.Get(data, "data.jobDetail.companyDetail.companySnapshot.size").String()
	avgProcessTime := gjson.Get(data, "data.jobDetail.companyDetail.companySnapshot.avgProcessTime").String()
	companyOverview := trimHtml(gjson.Get(data, "data.jobDetail.companyDetail.companyOverview.html").String())
	companyTelephoneNumber := gjson.Get(data, "data.jobDetail.companyDetail.companySnapshot.telephoneNumber").String()
	companyNearbyLocations := gjson.Get(data, "data.jobDetail.companyDetail.companySnapshot.nearbyLocations").String()
	jobDescription := trimHtml(gjson.Get(data, "data.jobDetail.jobDetail.jobDescription.html").String())
	jobCareerLeveln := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.careerLevel").String()
	jobYearsOfExperience := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.yearsOfExperience").String()
	jobQualification := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.qualification").String()
	jobFieldOfStudy := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.fieldOfStudy").String()
	jobSkills := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.skills").String()
	jobEmploymentType := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.employmentType").String()
	jobLanguages := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.languages").String()
	jobClosingDate := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.closingDate").String()
	jobFunctionValue := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.jobFunctionValue").String()
	category1 := gjson.Get(jobFunctionValue, "0.name").String()
	category2 := gjson.Get(jobFunctionValue, "1.name").String()
	jobBenefits := gjson.Get(data, "data.jobDetail.jobDetail.jobRequirement.benefits").String()
	location := gjson.Get(data, "data.jobDetail.location").String()
	locationStr := gjson.Get(location, "0.location").String()
	sourceCountry := gjson.Get(data, "data.jobDetail.sourceCountry").String()

	res = append(res, _id, pageUrl, logoUrl, salary, jobTitle, companyName, postedAt, companyWebsite, companySize, avgProcessTime, companyOverview, companyTelephoneNumber, companyNearbyLocations,
		jobDescription, jobCareerLeveln, jobYearsOfExperience, jobQualification, jobFieldOfStudy, jobSkills, jobEmploymentType,
		jobLanguages, jobClosingDate, jobFunctionValue, category1, category2, jobBenefits, location, locationStr, sourceCountry)
	s.ch <- res

	// time.Sleep(time.Second * 1)
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
