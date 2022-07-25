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

const fileName = "olx_inc_0725.csv"

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

var cookie = `locationPath=[{"id":4000030,"name":"Jakarta Selatan","type":"CITY","longitude":106.80105,"latitude":-6.2609,"parentId":2000007}]; laquesis=pan-59312@a#pan-59740@a#pan-60601@b#pan-67471@b; lqstatus=1658195958; G_ENABLED_IDPS=google; user=j:{"id":"122507225","name":"wang test"}; lf=1; reply_restriction=false; kyc_reply_count=0; t=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJhY2Nlc3NUb2tlbiIsImlzTmV3VXNlciI6ZmFsc2UsImlhdCI6MTY1ODczMTk0NiwiZXhwIjoxNjU4NzMyODQ2LCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.c6nWwNTyo0C4_XTudixOSackn6u_Qw6Y8h1BHhv5mQItbxMDpfnBkgyS2vG4xU0GEW-Fbg2CUAwnIugCdqP4tMU1d8SHRGZpqWpBadwTF9ergh-rS-N9_U755vlAsYskwERiwSC-pNhjdBgncX91-eCI0Q3KwtGqdc5xEQ1RlTXb54vMMN0EJbpLMgLO9HNXoQEvMY0vrSmRa4sW1RPaUW_z1RWSBlbWETz0yH28ij-9GGPK5qj23y_UIwLU7YMQZU0RIgDjO-p-C6dIy8AsYxJOVhJ4V8JKTqXHn5si8FCwzsYdl8vveTrkZixtzcp_zxIfDIsAlkB7sZ-lUNefZg; rt=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJncmFudFR5cGUiOiJncGx1cyIsImNsaWVudFR5cGUiOiJ3ZWIiLCJ0b2tlblR5cGUiOiJyZWZyZXNoVG9rZW4iLCJyZWFjdGl2YXRlZCI6ZmFsc2UsImlhdCI6MTY1ODczMTk0NiwiZXhwIjoxNjY2NzY3MTQ2LCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.H4yRarvDMIj190m5V58-Kkq0zaTpE6bkkpMKcnWLegNIVFO3axzk2KMp29gg3ntCMxu3Bw7aGlYuZgug1QpY5iIuGviuK6EkrqkN0gyCwTtWwBD3orONg4jhYXMQohbJ0PohqrcH89yCXZNbdEIHzMLv-Z0BwxfhtyhO7DirwdHIAIQtt3HMITVWiuB1Xuvo0O4jqc9AB2w9My4N71sfW9j2RUcHoSsOj3MUJ7YX2DahaGcDusjI9e0PkVRiid93hhe4l8TBMZuw70hLWN4VBqmCZXPzBIKIqNtE7_sviF_BdbM9QsjCda3ru2Loca3T3C8kgEwa9mbaLbiqGDxttA; ct=eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6ImViT21QTmlrIn0.eyJ0b2tlblR5cGUiOiJjaGF0VG9rZW4iLCJ2ZXJzaW9uIjoiMSIsImlhdCI6MTY1ODczMTk0NiwiZXhwIjoxNjU4ODE4MzQ2LCJhdWQiOiJvbHhpZCIsImlzcyI6Im9seCIsInN1YiI6IjEyMjUwNzIyNSIsImp0aSI6ImNiNWRmODFlZDIwMjU3YTQyOGQ4ZDMyZDNiYzAwNWJlMjc5MWMwYzYifQ.qxeU3AJYPdY1nlbY49kZDp9cH_srMpmg3fAOis-EJfPRS4yue6N3caTWuebsArivU10HRTWb8JR-dSpU5uXoe03XDVw20ANyljVrJ00ZRSAr9VE9AmS8iCpRZY8VGaLi98F6FSwynjAW9O0FH8lXpGLpyRKqxtdedEWds3rhzciQlZ5PXLDUcgkGQ7olqhpV995GKah4RlSIk4RLdlrvdrhCM0gOHkCHKNRSBkm-P_jt-mwzwS5Q9o4F8V-wNALQbqWJ1d5AW0uW041WmGy35tc87RaLrW5L6M96rvqaItyqfc4Eg6EQ-3q3rLxYPQmcjU5F_WqFnTBdglkRzqRhKg; nt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTg3MzE5NDYsImJyYW5kIjoib2x4IiwiY291bnRyeUNvZGUiOiJpZCIsInVzZXJJZCI6MTIyNTA3MjI1fQ.un1s6KOpRE0Q1EMFNLaX4NW4OF1Ayqaei4gnCIJ7B1s; _abck=EBED42F344620BC816C09B5808A2C2E3~0~YAAQjPxkX1ePjh+CAQAAXeMgNAg5hxQL2AGPfFxGpVhLfVAe3IAF6R3h176pYm28x6niWJHOiXsvbeZtlP/U3GsPvt6rp2WuYX9Q9w63ibFinjO9rqilcDWaQaO0Ply/idGlF4tPCdx5iKWNjjDmGsm5j14MeVgVXw+ZvUlGZr1zTTqdDUlpaooZixo4gtLAvvS4T2Jb/t9PQbULSdvsWZM6Ghrpbisk58hms65poYqbmvccV+aR53vChcccscpHd3HfDW22yvpqF9eFMt1n8wTHlCNBUCne8HYcXMzPzulF7pSdbtjyXEweezFlilNg67x42UD7r3Iu0EicpwQqCHRAVKSPoEylwdceA/sZNr0to9NeHQ5bNg1XgMhGh/paklqF9uyn2jTxSMbcoDOVS22ePeHCfA==~-1~-1~-1; ak_bmsc=7B24706E5FFDD53BE4F33A99946A680D~000000000000000000000000000000~YAAQjPxkX1iPjh+CAQAAXeMgNBAysfwX/1E66SUjpr11cWYLzK/3JJhsdAP9eUONT+xseDI5H5x7kK/auzzdFi/YMlVT0BWi7yUfHFN1bC3/FCqixguOQUPbGEaLgBl/P1qERc+hAMX6Ef77srIuLxJINudKDwQuD3QcTAjN0kUQC4CwgJqB+8yZKH65yWvaVAm8d5qof/2RKjPZJe3HFb34sJnfOcKwOwpc50Ba2o3hXu0+77GBXWLNvU25LijEVCGGCdLHM0BPS8d4gIAIPB5bKPtPonvPEEd7Op/N5Cyjj3lnb2WJC8mppF98ZBbCSI0TNZoVa3hwZCkFLp61uqau2LK+cXXAfT83HDBs9tYYF5DNgFhOg/yyrM5a2tuOxLYaWXe+0grw8Og=; bm_sz=3A299347B7EDDD7AB5BE21564E889515~YAAQjPxkX1qPjh+CAQAAXeMgNBB55/6t4EhY82gY+iqmvYXH0moK94RBWU1r3GUtiLisqVCRZEPGUuPvuOapqlXcVRo4CpDp2FwQPQzGvKRZIFcD+nqK41FkIANkDp92Vp+9t10sn3M8F3Cx53xoPgB3i3zw3j2yItDRPQbws1PQUBSOnUB1bUimgaSoZp9WBPon5dcV5Gr/3VY7BrW8Z8TqMTBbVsx8TgwzImyiz9J2/qlCDuZWUUCc7XCgFGBIpCZ3jrP6imrOtMJ0+YxyQxjG9rMn9cBY5fqbMiJhZuINzw==~3749936~3420469; bm_mi=DBF177490C1CD3210D3115D6C8C847A9~YAAQjPxkX7mPjh+CAQAANYshNBAvUykdbJ0Vc3ts/Fj9U9Ncz5oW3UFVpeGxiySGLVm4Mj+tO3Z9o+Ms8WNikYmo8lZqRg9wVrCqE3Yi6abuy/dt1k5hih32AS5jAHr4gRggacaK7eYaEpMLA3RF5/Wvx5ppU0aTf1d8nj2Pe1r/NCF9H/BoR3mAkvxAtvM+jHR4QXkzdEz11ZpTayhJexwv2m53QaWWbWBgU7TUjcV5nKOIBDG6DomKq9E9/4npFt+4FVY7ILfkyFSaLUjddpukoRreMpI08vWzY+1s+yiqxcV4YWW1kbdF0qPbbfbZxnytMvEm2JMGXicRXDqgHjQuElWEpERxQs2uFNLjSOu9nuPQ1JIjsOAPu67qAA7DNifR4taW817jnQ==~1; bm_sv=1B0DD136DC2B5C7DD8A890AC9EC16E16~YAAQjPxkX7uPjh+CAQAAeIshNBAg+g+l2XLmEs6nGC2NQEv35x28QquT3D4LcZlGSwfaKylTNG9hAZFYcC4mOukDAM74WClw/hC2uEe1UAPZw0Zpl7uWt/W4FcV6SteL5BkrQY1IOTctjP7yZyLbzeTFbyJOcFUjavZ4O2xdIgomg9bM7Dzr9oUXD/Dg+TttZEN9Qq9LqkQT2KMoopym29TVm5T2R5C3tyORAKeRaxTKcX8XQtOioRVlIgGgnas=~1`

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
