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

const fileName = "/home/meng/文档/myshare/colly/olx/record_test.csv"

// const fileName = "record_01.csv"

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
		if i < 8389 {
			continue
		}
		if len(row) == 19 {
			row = append(row, reqPhone(row[16]))
		}
		ch <- row
	}
}

var gCurCookieJar *cookiejar.Jar

func init() {
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)

}

var first bool

var cookie = `laquesis=pan-59312@a#pan-59740@a#pan-60601@b#pan-67471@b; lqstatus=1657185091; G_ENABLED_IDPS=google; locationPath=[{"id":4000030,"name":"Jakarta Selatan","type":"CITY","longitude":106.80105,"latitude":-6.2609,"parentId":2000007}]; reply_restriction=false; kyc_reply_count=0; G_AUTHUSER_H=0; lf=1; user=j:{"id":"122415481","name":"meng mytest"}; bm_sz=FF52FE1DB2BB42AFEE8BEDD9936812BA~YAAQz/xkXzwYDfOBAQAAsW8Q9hCdch4bNWG+ZoasezzmaJH+m/MoVCg19vBPLYVzBKRm4TfKV4JTltn17YqkV3zIsk5RSmUUq1cqjg35096sEdHxW3ZqOUWXV18ZI0TRbfxjIj8TYxbF0/+tAo6pN2S7bE3t0SpaQczYgP8+nVkeCysevBcsOg9KZdT25sM4sXTOtBvx1r3U7Y/Bf4xAo+Mrp2gOf5q6r/QFpHrxfM/mWdUyk7iCFMoX7t4T/UoXYhC6gVQe0pI1bqQmGh7Lxu0PbJYS5o/h5XY160vV4lbktFayJ1pdwx0ZpILNmySnWwf7H1IlbG15bQ==~3487032~3422275; bm_mi=65C3589D8B723386A365EDF7CA773A22~YAAQjPxkXwqMzeyBAQAAJOud9hAFqQG6L4mJht3c9DR4J6ikvtwg413jS3tBwWsNEFJrjJMgQaAQUCferOeQ5+1vyu5UptH6bDcrtNN84COhy07dMt83p4+nYJ2ZK0KSlJGzCTOlgc0AwyRTElSoFhqbJxMQF0+MP0PMPR0orOdIgh8/EWbVZkgLy8ExozZT9BJ/PJJxvkLaEpJvdwd3hAIMJLflSoHH2GURyV1wEQqQT6hiX+5Bbu2/Zpx09KbmIyQDMF4Qhui6obsKBnauoADVWGzzK5o8QmYzE5bd/U7jBcQ3IekX80bD7wWl0w8vWlCksCujbFdk5AEgWDebCRqfwUOhnDJOLBbPcrInWcP+goHz2ag/XDw=~1; bm_sv=5142F5FE9CD3E008FA63D217260C77CF~YAAQjPxkX06OzeyBAQAAyR2e9hCc7cwRa4o+8Ka2En8kHpJ3WbuAHFUsV9/SaP6shL3G3N6lZIkhpGgtk9aXmbvrHjIZzrDhGVCeRCWrrBH1qLLpj1wv70K5PtPAeGn5jefQ+9bb2+lwKJGQS2kZciCvbIRu1pxxPInFS4Xx7pQ5uK8QXMp9LA+WmNp73wDVnhZpfLL0/evL2FaHmwlhFgSp65DPc2Yn+sv0H3cYZCi31xQmhBorFsKufvLfYHPU~1; ak_bmsc=CAA50830D15EFF87E93F4FFD850785A4~000000000000000000000000000000~YAAQjPxkXw2QzeyBAQAAuEGe9hBIo3HCpFEoc+VR1FoRr/8rzvjYZpm4I159lyNq96tZGuJNJ9K+yK9MQYRTW5X8h1VJWZQSuCLx0vlNXk0VZFhiSUbUrIaaQHTzOjQ/HUV/4zyTUugbdjiSY3qZTmS0gELdm3P1bIx5Wx+OAp3cExGghwzTJX4d5TamJIYlesICcXH3X15WbaV5T0y3Q0ccu1OTJyvMl8iswCdGRgPX+EmpxDM/TgRLPkg86wnjgwkPFYoh4ET2cSqiQGR63sECu6pT6uJCVbVvUCxJpfVPh6RPQ324SbXNbLjcVJfM5j35UQ9wLpYTCWbUJVW23+fO1AZ7stXrA0OtYXATWpZ94e7LBdSNfcqu3ec0RE2A8o+zTYjokTdpjQ+M/NUtyH2l5JX1wNFFQCIQw1YuJ4xN4WxZH+6DhAcYJ/fIW68UpksuzP7x02qpIAmZDouBKsc3wGDuGrj+vD4Yhg==; t=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJhY2Nlc3NUb2tlbiIsImlzTmV3VXNlciI6ZmFsc2UsImlhdCI6MTY1NzcwMzU2MCwiZXhwIjoxNjU3NzA0NDYwLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjQxNTQ4MSIsImp0aSI6IjkyNDE0ZjUzZDlmZjQxMWRkNWVhZGZjOWZmMThjYzRhYmVlYzdiMzAifQ.CYs7aq9DgFio0obmtmaq68HoOBlkm2vFEukEtrzATUO6q8fRe86b3nhXH49vOqUWnl78FeroABhbybHqQ2xN_kR0rM-ssKXVDo_hzd_Tnor9RKkkoLA6enidg715es1mZA1oeq9rfaGP5IdgwLtjORsmeNq6pLGGCKKVCST8QcXQQZf4EJsriEB1HL0VVTzJQ1snVeeW9uoZmGOAdmbf9PFAun_iwk5VWsHzu9vbK-5043HBAkdaAbtdM-AVs_My27t4UMykLFo2PwCrPw1-0fLxPEh1lk4kCnptX_owZByho6ieeLDLNqTzWfNyHewJdJnx2tXJ8x1QNXCOAzjKSA; rt=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJyZWZyZXNoVG9rZW4iLCJyZWFjdGl2YXRlZCI6ZmFsc2UsImlhdCI6MTY1NzcwMzU2MCwiZXhwIjoxNjY1NzM4NzYwLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjQxNTQ4MSIsImp0aSI6IjkyNDE0ZjUzZDlmZjQxMWRkNWVhZGZjOWZmMThjYzRhYmVlYzdiMzAifQ.GfxHJhaQn9zxT-EtjK485HbdLQiJpczsI48JeEy__wcx7YXj9JEFkzAf4zCW4L1QJK6sTx80lq5zSPzbsIsX_BWB_Q5gLCe3Lkx73PiJ4N_g83nvyLRgm8EC9Y6LQiUtSSzDbnC0xzDpH5q7W2ldh73LasVANNyLJKgX5afQpGLKHcwd74BRfiPiz5edE1Z8jmh31rdVVGVBI0D1t9l-KkzWEUCh35Vj4VyrfPUUjYBIyayIcn4oXePddqCu8WVOPdOKA8TSxBzB4JaLxpwAwCCwEAEYqJS_47Y_21i7Xqw0_iFZlrXdrZIUjk61pwBRKjJQA9cCTSwz1ByT84-7YA; ct=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJ0b2tlblR5cGUiOiJjaGF0VG9rZW4iLCJ2ZXJzaW9uIjoiMSIsImlhdCI6MTY1NzcwMzU2MCwiZXhwIjoxNjU3Nzg5OTYwLCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjQxNTQ4MSIsImp0aSI6IjkyNDE0ZjUzZDlmZjQxMWRkNWVhZGZjOWZmMThjYzRhYmVlYzdiMzAifQ.eBS-mQq3vnNbpuj-Tmdq8Zm7FzY5bzBhtAPhLFFifIMoQUdubSC5iR5UowNdu424c_ghunMGNBk4EzQRXTQwuZ-Eb3qpOWaFjKiNkfvmiD7aSM2Ytl5iy8Kt9LB1y2wAMD3vV9BUj1_Rx0K-P_ZkwHSGPnCoP0dQ7suzN0oYvUV36BvKx9r6GCWte548BbXMnO5suCvAfbnUnxt_CjMrPOoBYvVSBHbkczIFUdYAXvojQyqklf_XlZqneJKIt13ZQejA-CONWknDxtlz3bZ2lto__t5zmE3PmLPxgXzhFJOVDAaKVbMOAw4mGb7CuVrQGGQ3mdBZYUyvCRauwQqaXA; nt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTc3MDM1NjAsImJyYW5kIjoib2x4IiwiY291bnRyeUNvZGUiOiJpZCIsInVzZXJJZCI6MTIyNDE1NDgxfQ.fch9TGemGP5fO4fnI2m4o4XKZoi7wInwuPwtgS9q7mY; _abck=2E28F73D2FBF2F1E22FF1258E5D1EF6B~-1~YAAQjPxkX/460OyBAQAAQhbV9giula7qHV/kef2qI4g2IzjBSAwkYDIvuJXX4cVVVl3Fk88txKnrCID72TCG1xd9XwDxWvZ5uRN0XGmhYXM9j4CFSsdJMUoqgxfjgdtxYhvYK2u7eq5/dWttkV7lsMxEI3eMk1Ov+33MmQAJUrCbKXrskVoYH4hA96lav4HRhdpTXqHTSqVuS1MBTvfIN+5lIQDO0iuRMKSLdM+EOiKg/hCFEnr0iaAEP6poEW9OUhay706EG93gEwd7hmyBVQa3dsDiXOTjg3pXtshlheyO5eHSsGl8qBPkMcqsPDgAiWWxDtIAfV7bwUVzZbPUkyW/hGI3XhkbSGILvf3OLhbZT1ve9IW29QbAnwER~-1~-1~-1`

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
