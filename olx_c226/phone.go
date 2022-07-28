package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

// const fileName = "/home/meng/文档/myshare/colly/olx/record_test.csv"

const fileName = "olx_inc_0727.csv"

// 获取手机号链接
var _phoneUrl = "https://www.olx.co.id/api/users/%s"

var id = 16

func getPhone(ch chan []string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	i := 0
	for {
		i++
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read row error:%v\n", err)
		}
		if i < 1 {
			continue
		}
		// if len(row) == 19 {
		row = append(row, reqPhone(row[16]))
		// }
		ch <- row
	}
}

var gCurCookieJar *cookiejar.Jar

func init() {
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)

}

var first bool

var cookie = `locationPath=[{"id":4000030,"name":"Jakarta Selatan","type":"CITY","longitude":106.80105,"latitude":-6.2609,"parentId":2000007}]; laquesis=pan-59312@a#pan-59740@a#pan-60601@b#pan-67471@b; lqstatus=1658195958; G_ENABLED_IDPS=google; user=j:{"id":"122507225","name":"wang test"}; lf=1; t=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJhY2Nlc3NUb2tlbiIsImlzTmV3VXNlciI6ZmFsc2UsImlhdCI6MTY1ODkxNDg1MSwiZXhwIjoxNjU4OTE1NzUxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.iouXFzGOdRInBTB0VMX65WKleWKAlSCIVuIcOsxXhlZ9zFwOfvp9sKGc-Y9bM7djl4mb82GNUHQlc5uFT01uhaU0A4Ae71c6ElAAb0-RoryVAlSAdTl8pA-RgO0EAcx4UBVutV6OIowYzeR6eF4kwzHsMevw9kIdELpIGnCY3rsSnmUHjEB44QCP18PsIolA_uz9XqpeJuEUsJMQHGogD7Oivcxk1qAuVvXENzTZdMtlcpVVmdAsXIjvFJIu1QpzYVyFypIKOLwMeLKd31eBsm6t3C0FPKooXL1RIj9g9ph46Bz56rVTMgkQJ3vp5iE-ca7loFq8Y1zrNB-_KxzBww; rt=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJyZWZyZXNoVG9rZW4iLCJyZWFjdGl2YXRlZCI6ZmFsc2UsImlhdCI6MTY1ODkxNDg1MSwiZXhwIjoxNjY2OTUwMDUxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.QGRA8ZW5xepMz9imh1H65MAdRSKg2mxu3QVPFSlq5tbV5Lxiy9Lkor12WZ8XDweuTTbulCXaSHeqDw6jv5kMFL6fE17VMU04Bwpi-pb2WFz88yDUlGmlP-q8dqsbZeEhN5-sH1ZUMEszKKhM0yQT8iYG7m651rwMtmI3QMLLKeozG3gx17IUfVMDVST98E0SWx3oqJ79aR9epJfk7OcRLJzL7RgZXQpAWNAftRSowyVGlI6zrNRP3B1SwEZDHAhnl2L2_-9T5kMivHsOmTptZecghTCd8Q2cb_TSkS9ur-_Dd70wV0lZBQMwEC3djHm1ec_6Tbt5XZoNQ4hVjpvyGg; ct=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJ0b2tlblR5cGUiOiJjaGF0VG9rZW4iLCJ2ZXJzaW9uIjoiMSIsImlhdCI6MTY1ODkxNDg1MSwiZXhwIjoxNjU5MDAxMjUxLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.rL9POVJGlDvWekT0RwA-PKaoh7uco-gzyG38YTkiImb35ezsUE4BWeD-iRiJAgwE84DAIC6ZlN2ab72SeX7C5tpUAdlvb8x0U1bOaXt_rs75hE6fOLv19TIe7BfuJhIHz6fF5wu9zS9TfyBRPribbGZCm5WIvxsdgWoY4oLp81T0ctl_L73mwdYwSTkOYbMQ_JLD3-1qTXYe-ALPmHWotx55sQ0hFc-nSETXlKmyA6pPRw7fApiM0Qfkin64o_b4vlgjscu_sFpOPElHnI70VxbKlInRjNBdbopUrgV1y_2cLQ3jEVgXzSC46bCJpzUN5DoFb5LsvP6ZdGDqBD88-g; nt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTg5MTQ4NTEsImJyYW5kIjoib2x4IiwiY291bnRyeUNvZGUiOiJpZCIsInVzZXJJZCI6MTIyNTA3MjI1fQ.S9smG-THaHhKWyFqv3brCaemG606b0W8sA5aDw-K320; ak_bmsc=F7D46913FE79C6CBC9A90C45B67BD04D~000000000000000000000000000000~YAAQtOtGaDaM6x+CAQAAz8oHPxBgBz6i4fz/3tJtUe55NjLYczQOk+pLZ57UJ/8I4k/Bagr4eRXqC3ekBvM6dnz3doJflKsMb1NQ8caokm0qeKfi/v7FW5kr92lvFbhVqY1HsQNwEDQlAgqsj0TjeNSFkbSEWICw0iK3DP5dXOIPXBDtfBw2apkSO9drM7WI4hvubykRj9JejaQ0MjrdWL3gae6Tw+erUpknwoQQoaWFb/3Ti+kqkClVUKb7SzEzNjW878EOgo7OKnxcjhuS2nvcd0cirbm6ov7+lANrCOy/wQkVSKGVkKyKO1ECYQymcNX6ExyAxQVKM1qwp3iTBbLjn4djQ0k68Ye3I1FPw0HKbuBx3wddfff5c9EZLwt1q2ylYx+lhGAFYHo=; bm_sz=83F801263C5CC1C448500D01E908D6BB~YAAQtOtGaDiM6x+CAQAAz8oHPxBaYXVr7hQFxvkQe0gdpl/uuz79KuqR4vkue5We1LPoklrARXR7d8oNxk2wmJJjHPpkfyd7UCSfi9kRRTlgYU8ltTLPzrSoKfTZw9njLnlWoKncXe0il/ZFjO+izxp9bTdYbg2CZ00QVBpafMgj56WkJ3QsxG4Zi/JSio30WjWAcGcAnFJ4Mx+V4TuJyoQzoTtkGQLjPe8PFLXwHLdbMeHXaAXDMUMU7bwpFT2JWbcJzkIIxoYbDkFC33F1kkK48Eb+Jw/GaygrOBDNCetpXQ==~4534840~4469570; bm_mi=C76C0A2A8F15FCC767BB2590A881F77D~YAAQtOtGaGiO6x+CAQAA1N4HPxDqY+laJd9oVEkRDF+evZK+rEYeAR5xPZOvdeyhopXSQR9j1IJXczmBBuVArifE9NL6OmVS4nq6Uz0AoMbuf93HlCxpsZPS1Zme3RoXsnLT4oQ1CM49Coc8AMP+XeEdwAttf61OMVdOvzYinHBFqzsp8SBAt+CA3XuYCGLB0dyafGEIFj92lqR/1s0KRZksxkNanSPHobWB8Tv+tTehnPyVofGgTDcK7pcPQphOMQ2TV3s2DPZlGl2n9qz0ppbuhe5RTdajfYg+HDhIgd29xR2cFvFELraP8SlUO8P1BBnOFvLCuWxjPg/biVnbMMUdVT6+tIEzBtQFatbPFNlBFZuTy0kn1r3jJA==~1; _abck=EBED42F344620BC816C09B5808A2C2E3~-1~YAAQtOtGaEyR6x+CAQAAYvcHPwhsYVSg/+6H6CXHTV6JJ7MhstNEzcCLxk/rRL+LFkqlkbPJH5MO3UPzEUymGeoA5FnUJSHLrE5zYqCKIsuZAhMSxvAs8yZyYAATeDH36nEkB44PQ/s4Iphq5K48sDgqL9YtJoo5WeTltELi2g0GPoXNGDn7XhFA+XImDYsU/Lrr9jhdNWbZ0CuorQE9xYSzt6cGPfoK87++HGZdLT9P5MRaYPa6980mL1zzGu4JCyNi6laN6UeIeKKN+vmFAYsQbeFrSLdtlZe5t6Gm5yaL+6Pm4aNRmUrjHGCPBG/Fe/8SoX5pAFueHOGjuHHU4MY3OKZFc5xGcXp7cu0mP3QyvBuWWY5bnU2EvVp5bnbAW/nPQSwpfjxWIQ/6qoCnpw4ejJ+8oxw=~-1~-1~-1; reply_restriction=false; kyc_reply_count=0; bm_sv=81E1064253EF00E5D8F03DEB7BBE4640~YAAQtOtGaCWS6x+CAQAAxvwHPxA801s/Jh6lfi7hXrDb2Da59ZX1IKpXJUkLMPxgSOv1q208bn9iCl10OKBe8oTYODDhm2V/JExmvDb/wuuwSJ/Gqb1yVksYqKDzLhHJpOosCJaW4ToLwRn7QWzB+nXJglMD4KwibPXQ04xpcIw5gyaogKNQSdqQbSlDQWCbxF6HP0BqEraV24tzFf817ZTYFI9KeHsoJ2JyGheofmVdOhm0wtl7XjFA04URWgc=~1`

func reqPhone(userID string) string {
	client := http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("GET", fmt.Sprintf(_phoneUrl, userID), nil) //GET大写

	req.Header.Set("cookie", cookie)
	// req.Header.Set("cookie", os.Getenv("OLX_COOKIE"))

	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("resp read err:%+v\n", err)
		return ""
	}

	//全局保存
	fmt.Printf("data:%s\n", data)
	if resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(time.Second * 3)
		panic(err)
	}
	// 提取手机号
	return gjson.Get(string(data), "data.phone").String()
}
